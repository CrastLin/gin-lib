package library

/**
 @auth CrastGin
 @date 2020-10
 */

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type Request struct {
	*http.Request
	Cxt        *gin.Context
	value      interface{}
	module     string
	controller string
	action     string
	lang       string
}

type Requester interface {
	Scheme() string
	Domain() string
	Query() string
	Port() string
	BaseUrl() string
	Url() string
	Module(name ...string) string
	Controller(name ...string) string
	Action() string
	Lang(name ...string) string
	Get(key string, defaultValue ...string) string
	GetArray(key string) []string
	GetMap(key string) map[string]string
	Post(key string, defaultValue ...string) string
	PostArray(key string) []string
	Input() []byte
	PostMap(key string) map[string]string
	PostJsonMap() map[string]interface{}
	PostJsonBind(T interface{}) interface{}
	Param(key string, defaultValue ...string) string
	IsAjax() bool
	IsPage() bool
	IsGet() bool
	IsPost() bool
	IsPut() bool
	IsDelete() bool
	IsOption() bool
	IsMobile() bool
	IsWeChat() bool
	IsQQ() bool
	Static(child interface{}, methodName string)
	StaticWithReturn(child interface{}, methodName string, params ...interface{}) interface{}
}

// get current scheme
func (r *Request) Scheme() string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme
}

// get current domain
func (r *Request) Domain() string {
	return fmt.Sprintf("%s://%s", r.Scheme(), r.Host)
}

// get query string
func (r *Request) Query() string {
	return r.RequestURI
}

// get port
func (r *Request) Port() string {
	var port string
	host := strings.Split(r.Host, ":")
	if len(host) == 2 {
		port = host[1]
	} else {
		if r.Scheme() == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return port
}

// get current full url
func (r *Request) Url() string {
	return r.Domain() + r.Query()
}

// get base url
func (r *Request) BaseUrl() string {
	baseUrl := strings.Split(r.Url(), "?")
	return baseUrl[0]
}

// get current module name
func (r *Request) Module(name ...string) string {
	if len(name) > 0 {
		r.module = name[0]
	}
	return r.module
}

// get current controller name
func (r *Request) Controller(name ...string) string {
	if len(name) > 0 {
		r.controller = name[0]
	}
	return r.controller
}

// get current action name
func (r *Request) Action() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	action := runtime.FuncForPC(pc[0])
	fullAction := strings.Split(action.Name(), ".")
	r.action = strings.ToLower(fullAction[len(fullAction)-1])
	return r.action
}

// get current lang mode
func (r *Request) Lang(name ...string) string {
	if len(name) > 0 {
		r.lang = name[0]
	}
	return r.lang
}

// get data
func (r *Request) Get(key string, defaultValue ...string) string {
	value := ""
	if len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	if getValue, ok := r.Cxt.Get(key); ok {
		if getValue, ok := getValue.(string); ok {
			value = getValue
		}
	}
	return value
}

// get array
func (r *Request) GetArray(key string) []string {
	return r.Cxt.QueryArray(key)
}

// get request map
func (r *Request) GetMap(key string) map[string]string {
	return r.Cxt.QueryMap(key)
}

// post data
func (r *Request) Post(key string, defaultValue ...string) string {
	value := ""
	if len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	return r.Cxt.DefaultPostForm(key, value)
}

// post form data
func (r *Request) PostArray(key string) []string {
	return r.Cxt.PostFormArray(key)
}

// post form array data
func (r *Request) PostMap(key string) map[string]string {
	return r.Cxt.PostFormMap(key)
}

// input raw data
func (r *Request) Input() []byte {
	input, _ := ioutil.ReadAll(r.Body)
	return input
}

// post json into map
func (r *Request) PostJsonMap() map[string]interface{} {
	input := r.Input()
	data := make(map[string]interface{})
	if err := json.Unmarshal(input, &data); err != nil {
		return nil
	}
	return data
}

// post json into struct
func (r *Request) PostJsonBind(T interface{}) interface{} {
	input := r.Input()
	if err := json.Unmarshal(input, &T); err != nil {
		return nil
	}
	return T
}

// get request param data
func (r *Request) Param(key string, defaultValue ...string) string {
	data := ""
	switch r.Method {
	case "POST":
		data = r.Post(key, defaultValue...)
		break
	case "GET":
		data = r.Get(key, defaultValue...)
		break
	}
	return data
}

// check request header is ajax
func (r *Request) IsAjax() bool {
	return r.Cxt.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

func (r *Request) IsPage() bool {
	return strings.Contains(r.RequestURI, ".html")
}

// check request method is get
func (r *Request) IsGet() bool {
	return r.Method == "GET"
}

// check request method is post
func (r *Request) IsPost() bool {
	return r.Method == "POST"
}

// check request method is post
func (r *Request) IsPut() bool {
	return r.Method == "PUT"
}

// check request method is post
func (r *Request) IsDelete() bool {
	return r.Method == "DELETE"
}

// check request method is post
func (r *Request) IsOption() bool {
	return r.Method == "OPTION"
}

// check client is mobile
func (r *Request) IsMobile() bool {
	userAgent := r.Cxt.GetHeader("User-Agent")
	if len(userAgent) == 0 {
		return false
	}
	isMobile := false
	keywords := []string{"Android", "iPhone", "iPod", "iPad", "Windows Phone", "MQQBrowser"}
	for i := 0; i < len(keywords); i++ {
		if strings.Contains(userAgent, keywords[i]) {
			isMobile = true
			break
		}
	}
	return isMobile
}

// check client is WeChat
func (r *Request) IsWeChat() bool {
	userAgent := r.Cxt.GetHeader("User-Agent")
	if len(userAgent) == 0 {
		return false
	}
	if strings.Contains(userAgent, "MicroMessenger") {
		return true
	}
	return false
}

// check client is QQ
func (r *Request) IsQQ() bool {
	userAgent := r.Cxt.GetHeader("User-Agent")
	if len(userAgent) == 0 {
		return false
	}
	if strings.Contains(userAgent, "V1_AND_SQ_") || strings.Contains(userAgent, "V1_IPH_SQ_") {
		return true
	}
	return false
}

// reflect call method
func (r *Request) Static(child interface{}, methodName string) {
	ref := reflect.ValueOf(child)
	method := ref.MethodByName(methodName)
	if method.IsValid() {
		_ = method.Call(make([]reflect.Value, 0))
	}
}

// reflect call method and return
func (r *Request) StaticWithReturn(child interface{}, method string, params ...interface{}) interface{} {
	ref := reflect.ValueOf(child)
	reflectMethod := ref.MethodByName(method)
	if reflectMethod.IsValid() {
		paramsLen := len(params)
		reflectParams := make([]reflect.Value, paramsLen)
		if paramsLen > 0 {
			for param := range params {
				reflectParams = append(reflectParams, reflect.ValueOf(param))
			}
		}
		result := reflectMethod.Call(reflectParams)
		return result[0].Interface()
	}
	return nil
}

type RequesterFactory interface {
	Instance(Cxt *gin.Context) Requester
}

type RequestFactory struct {
}

// instance request object
func (*RequestFactory) Instance(Cxt *gin.Context) Requester {
	factory := &Request{Request: Cxt.Request}
	factory.Cxt = Cxt
	return factory
}

// api response json data
type ApiJsonData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"result"`
}

// require api return format
func Callback(code int, message string, data ...interface{}) *ApiJsonData {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}
	return &ApiJsonData{Code: code, Message: message, Data: responseData}
}
