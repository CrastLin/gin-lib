package library

/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
	"sync"
	"time"
	"crast/util"
)

type Logger interface {
	InitFile(filePath string) (*os.File, error)
	Write(content string, args ...string) error
}

type Log struct {
	rootPath     string
	fullFileName string
}

var (
	logInstance *Log
	logOnce     sync.Once
	logSyncLock sync.Mutex
)

// auto create directory and open file stream
func (lo *Log) InitFile(filePath string) (*os.File, error) {
	rootPath := GetConfig("log.rootPath", "app").MustString("../")
	relPath := "general/"
	if filePath, ok := util.Ternary(filePath != "", filePath, "").(string); ok {
		relPath = filePath + "/"
	}
	date := time.Now().Format("2006-01-02")
	dateList := strings.Split(date, "-")
	directory := rootPath + "/" + relPath + dateList[0] + dateList[1]
	if _, err := os.Stat(directory); err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("check path :%s stat error :%v", directory, err))
		}
		// create directory
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("create directory :%s error :%v", directory, err))
		}
	}
	lo.fullFileName = directory + "/" + dateList[2] + ".log"

	return os.OpenFile(lo.fullFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// write error level logs
func (lo *Log) Write(content string, args ...string) error {
	filePath, _ := util.Ternary(args[0] != "", args[0], "general").(string)
	logSyncLock.Lock()
	defer logSyncLock.Unlock()

	file, err := lo.InitFile(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	if content == "" {
		return errors.New("log content is empty")
	}
	prefix := fmt.Sprintf("[%v|%s]", time.Now().Format("2006-01-02 15:04:05"), ClientIp)
	write := bufio.NewWriter(file)
	_, err = write.Write([]byte(fmt.Sprintf("%s %s \r\n", prefix, content)))
	if err != nil {
		return err
	}
	if err := write.Flush(); err != nil {
		return err
	}
	return nil
}

// singleton log
func LogInstance() Logger {
	logOnce.Do(func() {
		logInstance = &Log{}
	})
	return logInstance
}

// write diy directory log
func LogWrite(content string, args ...string) error {
	return LogInstance().Write(content, args...)
}

// write slice data log with json format
func LogRecord(content gin.H, args ...string) error {
	logContent, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return LogInstance().Write(string(logContent), args...)
}

// write notice level log
func LogNotice(content string) error {
	data := []string{"Notice"}
	return LogInstance().Write(content, data...)
}

// write error level log
func LogError(content string) error {
	data := []string{"Error"}
	return LogInstance().Write(content, data...)
}

// write warning level log
func LogWarning(content string) error {
	data := []string{"Warning"}
	return LogInstance().Write(content, data...)
}
