package assets

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi"
)

func (s *Server) handleGetData(t AssetType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		s.logger.Info("request for data", "t", t.String(), "key", key)

		// special assets are handled separately
		if _, ok := specialAssets[key]; ok {
			s.logger.Info("key is special", "key", key)
			data, err := staticAssets.ReadFile(fmt.Sprintf("static/special/%v.png", key))
			if err != nil {
				s.logger.Warn("unexpected special key but asset not found", "key", key, "err", err)
				s.handleNotFound(w)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}

		// try cache first
		data, err := s.loadFromCache(t, key)
		switch {
		case data != nil && err == nil:
			// load successful
			s.logger.Info("loaded from cache ok", "t", t.String(), "key", key)
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		case err != nil:
			s.logger.Info("unexpected error trying to read from cache", "t", t.String(), "key", key)
		case data == nil:
			s.logger.Info("cache data not found", "t", t.String(), "key", key)
		}

		if assetName, ok := s.assetNameMapping[t][key]; ok {
			hosts := s.hosts[t]
			// else try external 1 at a time
			for i, v := range hosts {
				joinedURL := v.JoinPath(fmt.Sprintf("/%v.png", assetName))
				s.logger.Info("trying external image source", "host", v.String(), "try", i, "key", key, "full_path", joinedURL.String())
				data, err := s.proxyImageRequest(joinedURL)
				if err != nil {
					s.logger.Info("error getting", "err", err, "host", v.String(), "try", i, "key", key, "full_path", joinedURL.String())
					continue
				}
				s.logger.Info("received image from source ok", "host", v.String())
				s.saveToCache(t, key, data)
				// found ok, save to cache and end request
				w.WriteHeader(http.StatusOK)
				w.Write(data)
				return
			}
		}
		// if we reached here then all external failed so we should serve default question mark image
		s.handleNotFound(w)
	}
}

func (s *Server) handleNotFound(w http.ResponseWriter) {
	// should return /static/misc/default.png
	data, err := staticAssets.ReadFile("static/misc/default.png")
	if err != nil {
		s.logger.Warn("error reading default.png", "err", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (s *Server) saveToCache(t AssetType, key string, data []byte) {
	fp := path.Join(s.cacheDir, fmt.Sprintf("/%v/%v.png", t.String(), key))
	f, err := os.Create(fp)
	if err != nil {
		s.logger.Warn("error writing to cache", "err", err)
	}
	_, err = f.Write(data)
	if err != nil {
		s.logger.Warn("error writing to cache", "err", err)
	}
	f.Close()
}

func (s *Server) loadFromCache(t AssetType, key string) ([]byte, error) {
	fp := path.Join(s.cacheDir, fmt.Sprintf("/%v/%v.png", t.String(), key))
	_, err := os.Stat(fp)
	if err == nil {
		return os.ReadFile(fp)
	} else if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected error checking file: %w", err)
}

func (s *Server) proxyImageRequest(path *url.URL) ([]byte, error) {
	resp, err := s.httpClient.Get(path.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("response is not an image: %v", contentType)
	}

	return io.ReadAll(resp.Body)
}

func (s *Server) handleOnlineCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
