package main

import (
	"errors"
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	var p pipeline
	var err error

	flag.StringVar(&p.root, "path", "./internal", "folder to look for character/artifact/weapon files")

	//load dm

	//walk through folder, handle each char
	err = p.walkCharacters()
	if err != nil {
		panic(err)
	}

	err = p.generateCharacters()
	if err != nil {
		panic(err)
	}
}

type pipeline struct {
	root       string
	characters []string
}

type characterConfig struct {
	PackageName string   `yaml:"package_name"`
	GenshinID   string   `yaml:"genshin_id"`
	Key         string   `yaml:"key"`
	Shortcuts   []string `yaml:"shortcuts"`
}

func (p *pipeline) load() error {

	return nil
}

func (p *pipeline) walkCharacters() error {
	err := filepath.Walk(p.root+"/characters",
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Println("error encountered while walking")
				return err
			}

			//we're only interested in finding directories
			//also we skip any root directories
			switch {
			case !info.IsDir():
				return nil
			case strings.HasSuffix(path, "/characters"):
				return nil
			}

			//we should be in a folder that's "/internal/characters/xxx"
			//try generateCharacter if

			// log.Println(path)
			p.characters = append(p.characters, path)

			return nil
		})

	if err != nil {
		return err
	}

	return nil
}

func (p *pipeline) generateCharacters() error {
	//check if config exists
	for _, path := range p.characters {
		charPath := path + "/config.yml"
		if _, err := os.Stat(charPath); errors.Is(err, os.ErrNotExist) {
			continue
		}
		err := p.generateCharacter(charPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *pipeline) generateCharacter(charPath string) error {
	data, err := os.ReadFile(charPath)
	if err != nil {
		log.Printf("error reading %v: %v\n", charPath, err)
	}
	c := characterConfig{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Printf("error parsing config %v: %v\n", charPath, err)
	}

	log.Printf("%v read ok; generating character\n", charPath)

	log.Println(c)

	return nil
}
