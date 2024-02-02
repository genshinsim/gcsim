package api

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (s *Server) decryptHash(ciphertext, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		s.Log.Warnw("decryptHash: error creating AES cipher", "err", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		s.Log.Warnw("decryptHash: error creating GCM", "err", err)
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		s.Log.Warnw("decryptHash: ciphertext < nonce size", "ciphertext", ciphertext)
		return nil, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		s.Log.Warnw("decryptHash: error decrypting ciphertext", "err", err)
		return nil, err
	}
	return plaintext, nil
}

func (s *Server) validateSigning(data []byte, str string) error {
	// check if from valid source
	// valid key is in the form of id:hash
	id, hashStr, ok := strings.Cut(str, ":")
	if !ok {
		return errors.New("no id:hash separation")
	}

	// hashStr is a hexstring
	hash, err := base64.StdEncoding.DecodeString(hashStr)
	if err != nil {
		return errors.New("hash not base64 string")
	}

	key, ok := s.cfg.AESDecryptionKeys[id]
	if !ok {
		return errors.New("id does not exist")
	}

	var res map[string]interface{}
	json.Unmarshal(data, &res)
	d, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("error marshaling: %w", err)
	}

	h := sha256.New()
	h.Write(d)
	bs := h.Sum(nil)

	dh, err := s.decryptHash(hash, key)
	if err != nil {
		return fmt.Errorf("error decrypting: %w", err)
	}

	if !bytes.Equal(bs, dh) {
		s.Log.Infow("create share request failed; hash not equal", "computed_sha256_hex_string", hex.EncodeToString(bs), "decrypted_hex_string", hex.EncodeToString(dh))
		return errors.New("bytes do not match")
	}

	s.Log.Infow("hash validation ok", "id_used", id)

	return nil
}
