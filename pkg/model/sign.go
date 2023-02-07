package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

func (s *SimulationResult) Sign(key string) (string, error) {
	if key == "" {
		return "", nil
	}

	id, aeskey, err := decodeKey(key)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(aeskey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	hash, err := s.hash()
	if err != nil {
		return "", err
	}

	signed := gcm.Seal(nonce, nonce, hash, nil)
	return id + ":" + base64.StdEncoding.EncodeToString(signed), nil
}

func decodeKey(key string) (string, []byte, error) {
	id, hexkey, ok := strings.Cut(key, ":")
	if !ok {
		return "", nil, errors.New("invalid key")
	}

	out, err := hex.DecodeString(hexkey)
	if err != nil {
		return "", nil, errors.New("invalid key")
	}

	return id, out, nil
}

func (s *SimulationResult) hash() ([]byte, error) {
	data, err := s.MarshalJson()
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	json.Unmarshal(data, &res)
	data, _ = json.Marshal(res)

	hash := sha256.Sum256(data)
	return hash[:], nil
}
