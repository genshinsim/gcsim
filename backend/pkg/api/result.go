package api

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type ResultStore interface {
	ResultReader
	Create(data []byte, ctx context.Context) (string, error)
	Update(id string, data []byte, ctx context.Context) error
	SetTTL(id string, ctx context.Context) error
	Delete(id string, ctx context.Context) error
	Random(ctx context.Context) (string, error)
}

type ResultReader interface {
	Read(id string, ctx context.Context) ([]byte, uint64, error)
}

var ErrKeyNotFound = errors.New("key does not exist")

const DefaultTLL = 24 * 14

func (s *Server) decryptHash(ciphertext, key []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
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

func (s *Server) CreateShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			s.Log.Errorw("unexpected error reading request body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !json.Valid(data) {
			s.Log.Info("request is not valid json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		str := r.Header.Get("X-GCSIM-SHARE-AUTH")
		if str == "" {
			s.Log.Infow("create share request failed - no hash received", "header", r.Header)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		s.Log.Infow("create share request received", "hash", str)

		//check if from valid source
		//valid key is in the form of id:hash
		id, hashStr, ok := strings.Cut(str, ":")
		if !ok {
			s.Log.Infow("create share request failed - invalid hash (no id:hash separation)")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		//hashStr is a hexstring
		hash, err := base64.StdEncoding.DecodeString(hashStr)
		if err != nil {
			s.Log.Infow("create share request failed - hash not valid hex string")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		key, ok := s.cfg.AESDecryptionKeys[id]
		if !ok {
			s.Log.Infow("create share request failed - id does not exist", "id", id)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		h := sha256.New()
		h.Write(data)
		bs := h.Sum(nil)

		dh, err := s.decryptHash(hash, key)
		if err != nil {
			s.Log.Infow("create share request failed; error decrypting", "err", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if !bytes.Equal(bs, dh) {
			s.Log.Infow("create share request failed; hash not equal", "computed_sha256_hex_string", hex.EncodeToString(bs), "decrypted_hex_string", hex.EncodeToString(dh))
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		uuid, err := s.cfg.ResultStore.Create(data, context.WithValue(r.Context(), TTLContextKey, DefaultTLL))

		if err != nil {
			s.Log.Errorw("unexpected error saving result", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.Log.Infow("create share request success", "key", uuid)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(uuid))
	}
}

func (s *Server) GetShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "share-key")

		share, ttl, err := s.cfg.ResultStore.Read(key, r.Context())
		switch err {
		case nil:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-gcsim-ttl", strconv.FormatUint(ttl, 10))
			w.WriteHeader(http.StatusOK)
			w.Write(share)
		case ErrKeyNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			s.Log.Errorw("unexpected error getting share", "err", err)
		}

	}
}

func (s *Server) GetRandomShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		share, err := s.cfg.ResultStore.Random(r.Context())
		switch err {
		case nil:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(share))
		case ErrKeyNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			s.Log.Errorw("unexpected error getting share", "err", err)
		}

	}
}
