package filecache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// FileCache is the main struct for file based caches.
type FileCache struct {
	FileName string
	Keys     map[string]CacheObject
}

// CacheObject is the representation of a value in the cache.
type CacheObject struct {
	Value   string
	Expires uint64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// New returns a new file cache from the specified file.
func New(fileName string) *FileCache {
	fileName = os.ExpandEnv(fileName)
	cache := FileCache{
		FileName: fileName,
		Keys:     make(map[string]CacheObject),
	}
	cache.parse()
	return &cache
}

// Get returns an element from the cache.
func (f *FileCache) Get(key string) (CacheObject, error) {
	obj := CacheObject{}
	value, ok := f.Keys[key]
	if !ok {
		return obj, errors.New("Not found")
	}
	obj = value
	return obj, nil
}

func (f *FileCache) parse() {
	fd, err := os.Open(f.FileName)
	if err != nil {
		// File does not exist, nothing to parse.
		return
	}
	defer fd.Close()

	dat, err := ioutil.ReadAll(fd)
	check(err)

	stored := FileCache{
		FileName: f.FileName,
		Keys:     make(map[string]CacheObject),
	}

	json.Unmarshal(dat, &stored)
	if err == nil {
		f.Keys = stored.Keys
	} else {
		fmt.Println("ERROR Reading cache", err)
	}
}

// Set a cache key
func (f *FileCache) Set(key, value string) {
	obj := CacheObject{value, 0}
	f.Keys[key] = obj
}

// Save the file cache in the store file.
func (f *FileCache) Save() {
	j, _ := json.Marshal(f)

	// fmt.Println("MARSHALLED: ", j)
	ioutil.WriteFile(f.FileName, j, 0644)
}
