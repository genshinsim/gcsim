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
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/go-chi/chi"
)

type ShareStore interface {
	ShareReader
	Create(ctx context.Context, data *model.SimulationResult) (string, error)
	Update(ctx context.Context, id string, data *model.SimulationResult) error
	SetTTL(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Random(ctx context.Context) (string, error)
}

type ShareReader interface {
	Read(ctx context.Context, id string) (*model.SimulationResult, uint64, error)
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

func (s *Server) validateShare(data []byte, str string) error {

	//check if from valid source
	//valid key is in the form of id:hash
	id, hashStr, ok := strings.Cut(str, ":")
	if !ok {
		return errors.New("no id:hash separation")
	}

	//hashStr is a hexstring
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
	data, _ = json.Marshal(res)

	h := sha256.New()
	h.Write(data)
	bs := h.Sum(nil)

	dh, err := s.decryptHash(hash, key)
	if err != nil {
		return fmt.Errorf("error decrypting: %v", err)
	}

	if !bytes.Equal(bs, dh) {
		s.Log.Infow("create share request failed; hash not equal", "computed_sha256_hex_string", hex.EncodeToString(bs), "decrypted_hex_string", hex.EncodeToString(dh))
		return errors.New("bytes do not match")
	}

	s.Log.Infow("hash validation ok", "id_used", id)

	return nil
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

		err = s.validateShare(data, str)
		if err != nil {
			s.Log.Infow("create share request - validation failed", "err", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var res *model.SimulationResult
		err = res.UnmarshalJson(data)
		if err != nil {
			s.Log.Infow("create share request - unmarshall failed", "err", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		id, err := s.cfg.ShareStore.Create(context.WithValue(r.Context(), TTLContextKey, DefaultTLL), res)

		if err != nil {
			s.Log.Errorw("unexpected error saving result", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.Log.Infow("create share request success", "key", id)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(id))
	}
}

func (s *Server) GetShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "share-key")

		share, ttl, err := s.cfg.ShareStore.Read(r.Context(), key)
		switch err {
		case nil:
			d, err := share.MarshalJson()
			if err != nil {
				s.Log.Errorw("unexpected error marshalling to json", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-gcsim-ttl", strconv.FormatUint(ttl, 10))
			w.WriteHeader(http.StatusOK)
			w.Write(d)
		case ErrKeyNotFound:
			http.Error(w, "not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			s.Log.Errorw("unexpected error getting share", "err", err)
		}

	}
}

func (s *Server) GetRandomShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		share, err := s.cfg.ShareStore.Random(r.Context())
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
