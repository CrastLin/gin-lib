package upload

/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type Local struct {
	Driver
}

// init local driver
func InitLocal(opt *Option) Uploader {
	local := &Local{Driver{config: opt}}
	return local
}

// check root path
func (l *Local) CheckPath(path ...string) bool {
	if len(path) > 0 && path[0] != "" {
		l.SaveFileName = path[0]
	} else {
		l.SaveFileName = l.config.RootPath + strings.Trim(l.config.SavePath, "/")
	}
	if _, err := os.Stat(l.SaveFileName); os.IsNotExist(err) {
		return l.MakeDir(l.SaveFileName)
	}
	return true
}

// make path directory
func (l *Local) MakeDir(path string) bool {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		l.config.ErrText = err.Error()
		return false
	}
	return true
}

// save file to path
func (l *Local) Save(file *multipart.FileHeader, replace bool) bool {
	src, err := file.Open()
	if err != nil {
		l.config.ErrText = err.Error()
		return false
	}
	defer src.Close()
	out, err := os.Create(l.SaveFileName)
	if err != nil {
		l.config.ErrText = err.Error()
		return false
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	if err != nil {
		l.config.ErrText = err.Error()
		return false
	}
	return true
}
