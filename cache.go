package mydictionary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// CacheItemStruct : item struct
type CacheItemStruct struct {
	QueryString  string   `json:"queryString"`
	Word         string   `json:"word"`
	Definition   []string `json:"definition"`
	Status       string   `json:"status"`
	CreationTime int64    `json:"creationTime"`
}

// CacheStruct : cache struct
type CacheStruct struct {
	path         string
	shelfLifeDay int64
	Content      []CacheItemStruct `json:"content"`
}

// Read : read cache
func (cache *CacheStruct) Read(path string, shelfLifeDay int64) (err error) {
	var (
		buf       []byte
		watershed int64
	)
	// set
	cache.path = path
	cache.shelfLifeDay = shelfLifeDay
	// read
	buf, err = ioutil.ReadFile(cache.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}
	err = json.Unmarshal(buf, cache)
	if err != nil {
		return
	}
	// remove item before watershed
	if cache.shelfLifeDay > 0 {
		watershed = time.Now().Unix() - cache.shelfLifeDay*24*60*60
		for i := len(cache.Content) - 1; i >= 0; i-- {
			if cache.Content[i].CreationTime < watershed {
				cache.Content = cache.Content[i+1:]
				break
			}
		}
	}
	return
}

// Query : query item in cache
func (cache *CacheStruct) Query(queryString string) (item CacheItemStruct, err error) {
	for i := 0; i < len(cache.Content); i++ {
		if strings.Compare(cache.Content[i].QueryString, queryString) == 0 {
			item = cache.Content[i]
			return
		}
	}
	err = fmt.Errorf("non-existent")
	return
}

// Add : add item to cache
func (cache *CacheStruct) Add(item CacheItemStruct) {
	cache.Content = append(cache.Content, item)
}

// Write : write cache
func (cache *CacheStruct) Write() (information string, err error) {
	var (
		buf  []byte
		path string
	)
	// write
	buf, err = json.MarshalIndent(cache, "", "\t")
	if err != nil {
		return
	}
	os.Remove(path)
	err = ioutil.WriteFile(cache.path, buf, 0644)
	if err != nil {
		return
	}
	// output
	information = fmt.Sprintf("Cache \"%s\" has been updated.\n\n", cache.path)
	return
}
