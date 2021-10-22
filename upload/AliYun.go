package upload
/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliYun struct {
	Driver
	option *AliYunOption
	client *oss.Client
}

type AliYunOption struct {
	Key         string
	Secret      string
	Bucket      string
	Endpoint    string
	Url         string
	SaveName    string
	AllowExtSet map[string]string
}

// init aliYun driver
func InitAliYun(opt *AliYunOption) Uploader {
	aliYun := &AliYun{}
	aliYun.option = opt
	client, err := oss.New(opt.Endpoint, opt.Key, opt.Secret)
	if err != nil {
		panic(fmt.Sprintf("connect aliYun oss was fialed:%s", err.Error()))
	}
	aliYun.client = client
	return aliYun
}

// check root path
func (a *AliYun) CheckPath(path ...string) bool {
	if ok, _ := a.client.IsBucketExist(a.option.Bucket); ok {
		option := []oss.Option{
			oss.ObjectACL(oss.ACLPublicRead),
		}
		err := a.client.CreateBucket(a.option.Bucket, option...)
		if err != nil {
			a.config.ErrText = err.Error()
			return false
		}
	}
	return true
}


// make path directory
func (*AliYun) MakeDir(path string) bool {
	return false
}

// save file to path
func (*AliYun) Save(file interface{}, replace bool) bool {
	return false
}
