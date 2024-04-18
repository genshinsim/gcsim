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
)

func (s *Server) loadFromFile(file string) ([]byte, error) {
	fp := path.Join(s.cacheDir, file)
	_, err := os.Stat(fp)
	if err == nil {
		return os.ReadFile(fp)
	} else if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}

	return nil, fmt.Errorf("unexpected error checking file: %w", err)
}

func (s *Server) loadFromSpecial(key string, w http.ResponseWriter) {

}

func (s *Server) proxyImageRequest(path url.URL) ([]byte, error) {
	resp, err := s.httpClient.Get(path.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("response is not an email: %v", contentType)
	}

	return io.ReadAll(resp.Body)
}

func (s *Server) handleNotFound(w http.ResponseWriter) {

}
