package library
/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"sync"
	"time"
)

type Redis struct {
	*redis.Client
}

var (
	redisInstance []*Redis
	redisOnce     []sync.Once
)

// get singleton redis
func RedisInstance(selectDb ...int) *Redis {
	db := 16
	if len(selectDb) > 0 && selectDb[0] >= 0 {
		db = selectDb[0]
	}
	if redisOnce == nil {
		redisOnce = make([]sync.Once, 18)
		redisInstance = make([]*Redis, 18)
	}
	index := db
	if redisInstance[index] == nil {
		redisOnce[index] = sync.Once{}
	}
	redisOnce[index].Do(func() {
		config := SourceConfig("redis")
		host := config.Get("host").MustString("127.0.0.1")
		port := config.Get("port").MustInt(6379)
		password := config.Get("password").MustString("")
		timeout := config.Get("timeout").MustInt(0)

		if db == 16 {
			db = config.Get("db").MustInt(0)
		}
		redisInstance[index] = &Redis{}
		redisInstance[index].Client = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Password:     password,
			DB:           db,
			WriteTimeout: time.Duration(timeout) * time.Second,
		})
		_, err := redisInstance[index].Ping().Result()
		if err != nil {
			_ = LogError(fmt.Sprintf("connect redis server was failed:%v", err.Error()))
		}
	})
	return redisInstance[index]
}

// close redis connect
func RedisClose(selectDb ...int) {
	db := 16
	if len(selectDb) > 0 && selectDb[0] >= 0 {
		db = selectDb[0]
	}
	if redisInstance[db] != nil {
		_ = redisInstance[db].Close()
		redisInstance[db] = nil
	}
	fmt.Println("close redis client:" + strconv.Itoa(db))
}

// make key with prefix
func RedisKey(key string, prefix string) string {
	return prefix + "_" + key
}
