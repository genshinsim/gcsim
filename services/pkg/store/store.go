package store

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Simulation struct {
	Metadata    string `json:"metadata"`
	ViewerFile  string `json:"viewer_file"`
	IsPermanent bool   `json:"is_permanent"`
}

type Result struct {
	Config string `json:"config_file"`
}

func (s *Simulation) DecodeViewer() (*Result, error) {
	//base64 zlib encoded string
	z, err := base64.StdEncoding.DecodeString(s.ViewerFile)
	if err != nil {
		return nil, err
	}

	//decompress
	reader, err := zlib.NewReader(bytes.NewReader(z))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var target Result
	err = json.Unmarshal(b, &target)
	return &target, err
}

type SimStore interface {
	Fetch(url string) (Simulation, error)
}

type DBEntry struct {
	Key          string `json:"simulation_key"`
	GitHash      string `json:"git_hash"`
	HashedConfig string `json:"config_hash"`
	Author       int64  `json:"author"`
	AuthorString string `json:"author_string,omitempty"`
	Description  string `json:"sim_description"`
}

type SimInfo struct {
	Key         string `json:"simulation_key"`
	Description string `json:"sim_description"`
}

func (d *DBEntry) ConvertConfig() error {
	var b bytes.Buffer
	zw := zlib.NewWriter(&b)
	if _, err := zw.Write([]byte(d.HashedConfig)); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}

	d.HashedConfig = base64.StdEncoding.EncodeToString(b.Bytes())

	return nil
}

type SimDBStore interface {
	Add(entry DBEntry) (int64, error)
	Replace(key string, entry DBEntry) (int64, error)
	List(char string) ([]SimInfo, error)
	Delete(key string) (int64, error)
}

type PostgRESTStore struct {
	URL    string
	client *http.Client
}

func NewPostgRESTStore(url string) *PostgRESTStore {
	return &PostgRESTStore{
		URL:    url,
		client: &http.Client{},
	}
}

func (b *PostgRESTStore) postRPCRequestReturnInt(jsonStr []byte, url string) (int64, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := b.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	msg, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading body: %v", err)
	}

	if r.StatusCode != 200 {
		return 0, fmt.Errorf("bad status code %v msg %v", r.StatusCode, string(msg))
	}

	//body should be db key
	id, err := strconv.ParseInt(string(msg), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing result id: %v", err)
	}

	return id, nil
}

func (b *PostgRESTStore) Add(entry DBEntry) (int64, error) {
	url := fmt.Sprintf(`%v/rpc/add_db_sim`, b.URL)
	entry.AuthorString = ""
	err := entry.ConvertConfig()
	if err != nil {
		return 0, err
	}
	jsonStr, err := json.Marshal(entry)
	if err != nil {
		return 0, err
	}
	return b.postRPCRequestReturnInt(jsonStr, url)
}

func (b *PostgRESTStore) Delete(key string) (int64, error) {
	url := fmt.Sprintf(`%v/rpc/delete_from_db`, b.URL)
	var data = struct {
		Key string `json:"simulation_key"`
	}{
		Key: key,
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return b.postRPCRequestReturnInt(jsonStr, url)
}

func (b *PostgRESTStore) Replace(key string, entry DBEntry) (int64, error) {
	url := fmt.Sprintf(`%v/rpc/replace_db_sim`, b.URL)
	entry.AuthorString = ""
	err := entry.ConvertConfig()
	if err != nil {
		return 0, err
	}
	var data = struct {
		DBEntry
		OldKey string `json:"old_key"`
	}{
		DBEntry: entry,
		OldKey:  key,
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return b.postRPCRequestReturnInt(jsonStr, url)
}

func (b *PostgRESTStore) Fetch(key string) (Simulation, error) {
	url := fmt.Sprintf(`%v/active_sim?simulation_key=eq.%v`, b.URL, key)
	resp, err := http.Get(url)
	if err != nil {
		return Simulation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return Simulation{}, err
		} else {
			return Simulation{}, fmt.Errorf("bad status code %v msg %v", resp.StatusCode, string(msg))
		}
	}

	var result []Simulation

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Simulation{}, err
	}

	if len(result) == 0 {
		return Simulation{}, fmt.Errorf("unexpected result length is 0")
	}

	return result[0], nil
}

func (b *PostgRESTStore) List(key string) ([]SimInfo, error) {
	url := fmt.Sprintf(`%v/db_sims_by_avatar?avatar_name=eq.%v`, b.URL, key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf("bad status code %v msg %v", resp.StatusCode, string(msg))
		}
	}

	var result []SimInfo

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("unexpected result length is 0")
	}

	return result, nil
}
