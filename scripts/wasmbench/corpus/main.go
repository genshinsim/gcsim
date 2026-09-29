// Command corpus turns a dump of the gcsim DB (one protojson db.Entry per .json file, in any
// directory layout) into a directory of config files the wasmbench corpus tools can run.
//
//	corpus -in dump/ -out corpus/
//
// It writes:
//   - out/configs/<id>.txt: the config of every entry whose config text is not a duplicate of an
//     earlier entry's (entries are taken in id order)
//   - out/entries.tsv: one row per parsed entry, with the metadata used to stratify samples:
//     the entry's folder (e.g. approved), is_db_valid, sim mode, target count, character names
//     (sorted, comma separated), mean sim duration and, for duplicates, the id whose config is
//     used
//   - out/unparsed.tsv: files that aren't a db.Entry, with the error
//
// See ../corpus-extract.sh, which also unpacks an archive of the dump.
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
)

type entry struct {
	folder string
	e      *db.Entry
}

func main() {
	in := flag.String("in", "", "directory holding the dump's .json files (searched recursively)")
	out := flag.String("out", "", "output directory")
	flag.Parse()
	if *in == "" || *out == "" {
		log.Fatal("usage: corpus -in dump/ -out corpus/")
	}
	if err := os.MkdirAll(filepath.Join(*out, "configs"), 0o755); err != nil {
		log.Fatal(err)
	}

	var entries []entry
	var unparsed []string
	total := 0
	opts := protojson.UnmarshalOptions{DiscardUnknown: true}
	err := filepath.WalkDir(*in, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") {
			return err
		}
		total++
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		e := &db.Entry{}
		if err := opts.Unmarshal(raw, e); err != nil {
			unparsed = append(unparsed, fmt.Sprintf("%s\t%v", path, err))
			return nil
		}
		if e.GetId() == "" || strings.TrimSpace(e.GetConfig()) == "" {
			unparsed = append(unparsed, fmt.Sprintf("%s\tno id or config", path))
			return nil
		}
		rel, _ := filepath.Rel(*in, filepath.Dir(path))
		entries = append(entries, entry{folder: filepath.Base(rel), e: e})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	slices.SortFunc(entries, func(a, b entry) int { return strings.Compare(a.e.GetId(), b.e.GetId()) })

	var rows strings.Builder
	rows.WriteString("id\tfolder\tdb_valid\tmode\ttargets\tnchars\tchars\tdur_mean\tdup_of\n")
	firstByHash := map[[32]byte]string{}
	unique := 0
	for _, x := range entries {
		e := x.e
		sum := sha256.Sum256([]byte(e.GetConfig()))
		dupOf := firstByHash[sum]
		if dupOf == "" {
			firstByHash[sum] = e.GetId()
			unique++
			if err := os.WriteFile(filepath.Join(*out, "configs", e.GetId()+".txt"), []byte(e.GetConfig()), 0o644); err != nil {
				log.Fatal(err)
			}
		}
		s := e.GetSummary()
		chars := slices.Clone(s.GetCharNames())
		slices.Sort(chars)
		fmt.Fprintf(&rows, "%s\t%s\t%t\t%s\t%d\t%d\t%s\t%s\t%s\n",
			e.GetId(), x.folder, e.GetIsDbValid(), s.GetMode(), s.GetTargetCount(), len(chars),
			strings.Join(chars, ","), strconv.FormatFloat(s.GetSimDuration().GetMean(), 'f', 1, 64), dupOf)
	}
	if err := os.WriteFile(filepath.Join(*out, "entries.tsv"), []byte(rows.String()), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(*out, "unparsed.tsv"), []byte(strings.Join(unparsed, "\n")), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d files, %d parsed as db.Entry, %d not, %d distinct configs written to %s\n",
		total, len(entries), len(unparsed), unique, filepath.Join(*out, "configs"))
}
