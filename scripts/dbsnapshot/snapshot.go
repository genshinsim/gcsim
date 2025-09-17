package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/proto"
)

type snapshot struct {
	ids     []string
	results []*model.SimulationResult
}

func (s *snapshot) save(filename string) error {
	if len(s.ids) != len(s.results) {
		return fmt.Errorf("mismatched ids and results length")
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i, result := range s.results {
		// Marshal protobuf to []byte
		data, err := proto.Marshal(result)
		if err != nil {
			return err
		}

		// Convert to base64 string
		b64 := base64.StdEncoding.EncodeToString(data)

		// Write id and base64 string separated by '|'
		_, err = writer.WriteString(s.ids[i] + "|" + b64 + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func load(filename string) (*snapshot, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ids []string
	var results []*model.SimulationResult
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // Skip empty lines
		}

		// Split line into id and base64 string
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			return nil, err
		}
		id := parts[0]
		b64 := parts[1]

		// Decode base64 string to []byte
		data, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, err
		}

		// Unmarshal []byte to protobuf
		result := &model.SimulationResult{}
		err = proto.Unmarshal(data, result)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &snapshot{ids: ids, results: results}, nil
}
