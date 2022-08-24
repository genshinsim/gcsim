package pygenerator

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/genshinsim/gcsim/services/pkg/embed"
	"github.com/genshinsim/gcsim/services/pkg/store"
)

type Generator struct {
	scriptPath string
}

func New(scriptPath string) embed.ImageGenerator {
	return &Generator{
		scriptPath: scriptPath,
	}
}

func (g *Generator) Generate(sim store.Simulation, filepath string) error {
	meta := sim.Metadata
	meta = strings.TrimPrefix(meta, `"`)
	meta = strings.TrimSuffix(meta, `"`)

	cmd := exec.Command(g.scriptPath, filepath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	_, err = stdin.Write([]byte(meta))
	if err != nil {
		log.Fatal(err)
	}
	stdin.Close()

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err = cmd.Run()
	if err != nil {
		return err
	}
	log.Println(stdBuffer.String())

	return nil
}
