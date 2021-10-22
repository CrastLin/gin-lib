package library

/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	basePath   string
	encryption bool
}

type CacheInterface interface {
	Has(name string) bool
	Set(name string, value interface{}, expire ...int) error
	Get(name string) (interface{}, error)
	Destroy(name string) bool
}

var (
	cacheInstance   *Cache
	cacheOnce       sync.Once
	cacheRwSyncLock sync.RWMutex
)

// singleton mode
func CacheInstance(isEncryption ...bool) CacheInterface {
	cacheOnce.Do(func() {
		cacheInstance = &Cache{basePath: GetConfig("cache.rootPath").MustString("../runtime/data")}
	})
	cacheInstance.encryption = false
	if len(isEncryption) > 0 && isEncryption[0] {
		cacheInstance.encryption = true
	}
	return cacheInstance
}

// make cache layer path
func CacheMake(root string, filename string, level int) string {
	var directory []string
	var sep = "/"
	if level < 1 || level > 32 {
		level = 0
	}
	if level > 0 {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(filename)))
		for i := 0; i < level; i++ {
			directory = append(directory, hash[i:i+1])
		}
	}
	dir := sep
	if len(directory) > 0 {
		dir = fmt.Sprintf("%s%s%s", sep, strings.Join(directory, sep), sep)
	}
	return fmt.Sprintf("%s%s%s", root, dir, filename)
}

// check encryption
func (c *Cache) checkEncryption(key string) string {
	if c.encryption {
		keys := strings.Split(key, "/")
		last := len(keys) - 1
		keys[last] = fmt.Sprintf("%x.d", md5.Sum([]byte(strings.TrimSuffix(keys[last], ".d"))))
		return strings.Join(keys, "/")
	}
	return key
}

// check file exists
func (c *Cache) checkExists(file string) bool {
	_, err := os.Stat(c.checkEncryption(file))
	if err == nil {
		return true
	}
	return false
}

// init file stream
func (c *Cache) source(name string, flag int) (*os.File, error) {
	names := strings.Split(name, "/")
	directory := strings.TrimSuffix(c.basePath, "/")
	count := len(names)
	if count > 1 {
		directory += "/" + strings.Join(names[:count-1], "/")
	}
	file := directory + "/" + names[count-1] + ".d"
	if !c.checkExists(file) {
		if flag == os.O_CREATE|os.O_WRONLY {
			// create directory
			err := os.MkdirAll(directory, 0755)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("create directory :%s error :%v", file, err))
			}
		} else {
			return nil, errors.New("cache file is not exists")
		}
	}
	return os.OpenFile(c.checkEncryption(file), flag, 0755)
}

// check cache exists
func (c *Cache) Has(name string) bool {
	file := strings.TrimSuffix(c.basePath, "/") + "/" + name + ".d"
	exists := c.checkExists(file)
	if !exists {
		return false
	}
	data, err := c.Get(name)
	if err == nil && data != nil {
		return true
	}
	return false
}

// set cache
func (c *Cache) Set(name string, value interface{}, expire ...int) error {
	cacheRwSyncLock.Lock()
	defer cacheRwSyncLock.Unlock()
	initFile, err := c.source(name, os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return err
	}
	defer initFile.Close()

	cacheExpire := 0
	if len(expire) > 0 && expire[0] > 0 {
		cacheExpire = int(time.Now().Unix()) + expire[0]
	}
	cacheString, err := json.Marshal(map[string]interface{}{"data": value, "expire": cacheExpire,})
	if err != nil {
		return err
	}
	write := bufio.NewWriter(initFile)
	if _, err = write.Write(cacheString); err != nil {
		return err
	}
	if err = write.Flush(); err != nil {
		return err
	}
	return nil
}

// get cache
func (c *Cache) Get(name string) (interface{}, error) {
	cacheRwSyncLock.RLock()
	initFile, err := c.source(name, os.O_RDONLY)
	if err != nil {
		cacheRwSyncLock.RUnlock()
		return nil, err
	}
	info, err := initFile.Stat()
	if err != nil {
		cacheRwSyncLock.RUnlock()
		initFile.Close()
		return nil, err
	}
	// buff read file
	reader := bufio.NewReader(initFile)

	cacheValue := make([]byte, info.Size())
	_, err = reader.Read(cacheValue)
	// close lock and file sources
	cacheRwSyncLock.RUnlock()
	initFile.Close()

	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(cacheValue, &data)
	if err != nil {
		return nil, err
	}
	if expire, ok := data["expire"]; ok {
		currentTime := time.Now().Unix()
		if expireTime, ok := expire.(float64); ok && expireTime > 0 && int64(expireTime) < currentTime {
			c.Destroy(name)
			return nil, nil
		}
	}
	if result, ok := data["data"]; ok {
		return result, nil
	}
	return nil, nil
}

// destroy cache
func (c *Cache) Destroy(name string) bool {
	cacheRwSyncLock.Lock()
	defer cacheRwSyncLock.Unlock()
	file := strings.TrimSuffix(c.basePath, "/") + "/" + name + ".d"
	err := os.Remove(c.checkEncryption(file))
	if err != nil {
		return false
	}
	return true
}
