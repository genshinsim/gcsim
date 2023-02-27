package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	sha1ver string
)

func main() {
	//steps:
	// 1. check github for latest release hash to check if out of date
	// 2. ask backend server for compute work to do
	// 3. run sim
	// 4. post result to callback url
	var err error
	log.Println("Checking for new version")
	err = checkHashOK()
	if err != nil {
		log.Panicf("Error encountered: %v", err)
	}
	log.Println("Version ok. Getting latest job")
	var w work
	for {
		w, err = getWork()
		//blank key means no more work
		if w.Key == "" {
			break
		}
	}

}

type latestInfo struct {
	TagName string `json:"tag_name"`
}

type latestHash struct {
	Object struct {
		SHA string `json:"sha"`
	} `json:"object"`
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func checkHashOK() error {
	var err error
	const latestUrl = `https://api.github.com/repos/genshinsim/gcsim/releases/latest`
	const tagUrl = `https://api.github.com/repos/genshinsim/gcsim/git/ref/tags/`

	var info latestInfo
	err = getJson(latestUrl, &info)
	if err != nil {
		return err
	}

	log.Printf("Latest tag found: %v\n", info.TagName)

	var tagInfo latestHash
	err = getJson(tagUrl+info.TagName, &tagInfo)
	if err != nil {
		return err
	}

	log.Printf("Latest hash found: %v\n", tagInfo.Object.SHA)

	if sha1ver != tagInfo.Object.SHA {
		//TODO: add hash check here; for now assume it's fine
		return nil
	}

	return nil
}

type work struct {
	Key    string `json:"key"`
	Config string `json:"config"`
}

func getWork() (work, error) {
	const url = `https://simimpact.app/api/db/work`
	var w work
	err := getJson(url, &w)
	if err != nil {
		return work{}, err
	}

	return w, nil
}
