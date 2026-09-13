package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

type FetchFunc func(root, name string) ([]byte, error)

func FetchAny(list ...FetchFunc) FetchFunc {
	return func(root, name string) ([]byte, error) {
		var data []byte
		err := errors.New("fetchany: requires at least one argument")
		for _, fn := range list {
			data, err = fn(root, name)
			if err == nil {
				return data, nil
			}
		}
		return data, err
	}
}

func UseCache(invalidate bool, fetch FetchFunc) FetchFunc {
	return func(root, name string) ([]byte, error) {
		if invalidate {
			_ = cacheRoot.Remove(name)
		} else {
			data, err := cacheRoot.ReadFile(name)
			if err == nil {
				return data, nil
			}
		}

		data, err := fetch(root, name)
		if err == nil {
			_ = cacheRoot.MkdirAll(path.Dir(name), 0o755)
			_ = cacheRoot.WriteFile(name, data, 0o644)
		}
		return data, err
	}
}

func FetchLocal(dir string) FetchFunc {
	return func(_, name string) ([]byte, error) {
		return os.ReadFile(filepath.Join(dir, name))
	}
}

func FetchHTTP(base string) FetchFunc {
	return func(_, name string) ([]byte, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		url, err := url.JoinPath(base, name)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(resp.Status)
		}
		return data, err
	}
}
