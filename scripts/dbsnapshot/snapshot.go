package main

import (
	"bufio"
	"encoding/base64"
	"os"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/proto"
)

type snapshot struct {
	ids     []string
	results []*model.SimulationResult
}

func (s *snapshot) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, result := range s.results {
		// Marshal protobuf to []byte
		data, err := proto.Marshal(result)
		if err != nil {
			return err
		}

		// Convert to base64 string
		b64 := base64.StdEncoding.EncodeToString(data)

		// Write to file with newline
		_, err = writer.WriteString(b64 + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func Load(filename string) (*snapshot, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []*model.SimulationResult
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // Skip empty lines
		}

		// Decode base64 string to []byte
		data, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return nil, err
		}

		// Unmarshal []byte to protobuf
		result := &model.SimulationResult{}
		err = proto.Unmarshal(data, result)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &snapshot{results: results}, nil
}
