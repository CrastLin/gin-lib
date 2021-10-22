package library

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strings"
	"crast/library/upload"
)

/**
 @auth CrastGin
 @date 2020-10
 */
type UploadFactory struct {
	driver string
	config *upload.Option
}

type FactoryUploader interface {
	SaveMultipart(form *multipart.Form) bool
	Save(file *multipart.FileHeader) bool
	GetError() string
}

// get upload client
func Upload(driver string, opt *upload.Option) FactoryUploader {
	// init base save path
	basePath := GetConfig("upload.basePath").MustString("../runtime/upload")
	if opt.RootPath != "" {
		opt.RootPath = strings.Trim(opt.RootPath, "/") + "/"
	}
	opt.RootPath = strings.TrimSuffix(basePath, "/") + "/" + opt.RootPath
	// init max size limit
	if opt.MaxSize <= 0 {
		opt.MaxSize = GetConfig("upload.maxSize").MustInt64(324288)
	}
	// init allow ext limit
	if len(opt.Exts) == 0 {
		extConfig := GetConfig("upload.allowExt").MustString("jpg,jpeg,png,gif,bmp")
		opt.Exts = strings.Split(extConfig, ",")
	}
	client := &UploadFactory{driver: driver, config: opt}
	return client
}

// get upload error
func (u *UploadFactory) GetError() string {
	return u.config.ErrText
}

// upload multipart files
func (u *UploadFactory) SaveMultipart(form *multipart.Form) bool {
	for _, files := range form.File {
		for _, file := range files {
			result := u.Save(file)
			if !result {
				return false
			}
		}
	}
	return true
}

// upload file
func (u *UploadFactory) Save(file *multipart.FileHeader) bool {
	if file == nil {
		u.config.ErrText = "没有上传的图片"
		return false
	}

	var driver upload.Uploader
	switch strings.ToLower(u.driver) {
	case "local":
		driver = upload.InitLocal(u.config)
		break
	case "ftp":
		driver = upload.InitFtp(u.config)
		break
	case "aliyun":
		driver = upload.InitAliYun(u.config.AliYunConfig)
		break
	default:
		u.config.ErrText = fmt.Sprintf("driver %s is not defined", u.driver)
		return false
	}
	// check allow size
	if file.Size > u.config.MaxSize {
		u.config.ErrText = "文件超出允许上传的大小"
	}

	src, err := file.Open()

	if err != nil {
		u.config.ErrText = err.Error()
		return false
	}

	defer src.Close()

	sources, err := ioutil.ReadAll(src)
	if err != nil {
		u.config.ErrText = err.Error()
		return false
	}

	fileExt := upload.GetFileType(sources[:10])

	//check allow mines
	fileMine := file.Header.Get("Content-Type")
	if len(u.config.Mines) > 0 {
		isAllowMine := false
		for _, allowMine := range u.config.Mines {
			if allowMine == fileMine {
				isAllowMine = true
				break
			}
		}
		if !isAllowMine {
			u.config.ErrText = fmt.Sprintf("不允许上传MINE类型: %s", fileMine)
			return false
		}
	}
	// check allow file ext
	splitMine := strings.Split(fileMine, "/")
	fileExt = splitMine[0]
	if len(splitMine) > 0 {
		fileExt = splitMine[1]
	}
	if len(u.config.Exts) > 0 {
		isAllowExt := false
		for _, ext := range u.config.Exts {
			if fileExt == ext {
				isAllowExt = true
				break
			}
		}
		if !isAllowExt {
			u.config.ErrText = fmt.Sprintf("不允许上传文件类型: %s", fileMine)
			return false
		}
	}

	// check root path
	if !driver.CheckPath() {
		u.config.ErrText = driver.GetError()
		return false
	}

	return true
}
