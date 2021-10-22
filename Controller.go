package library

/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"runtime"
	"strings"
)

type Controller struct {
	Request      *Request
	TemplateData map[string]interface{}
}

// before init handle
func (c *Controller) Init(child interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ClientIp = ctx.ClientIP()
		c.Request = &Request{Request: ctx.Request}
		c.Request.Cxt = ctx
		c.TemplateData = make(map[string]interface{})
		// call initialize func
		c.Request.Static(child, "Initialize")
		controllerName := strings.Split(strings.ToLower(reflect.TypeOf(child).String()), ".")
		c.Request.Module(strings.TrimSuffix(strings.TrimPrefix(controllerName[0], "*"), "controller"))
		c.Request.Controller(controllerName[1])
		// handle all exception
		defer func() {
			if exception := recover(); exception != nil {
				errText := fmt.Sprintf("api restful has exception --> %v", exception)
				fmt.Println(errText)
				_ = LogError(errText)
			}
		}()
		// run next context
		ctx.Next()
		return
	}
}

// init controller
func (*Controller) Initialize() {

}

// success redirect
func (c *Controller) Success(message string, params ...string) {
	if c.Request.IsAjax() {
		c.Request.Cxt.JSON(200, ApiJsonData{1, message, params})
	} else {
		c.Request.Cxt.HTML(200, "success.html", ApiJsonData{1, message, params})
	}
}

// error redirect
func (c *Controller) Error(message string, params ...string) {
	if c.Request.IsAjax() {
		c.Request.Cxt.JSON(200, ApiJsonData{0, message, params})
	} else {
		c.Request.Cxt.HTML(200, "error.html", ApiJsonData{0, message, params})
	}
}

// assign template var
func (c *Controller) Assign(key interface{}, data ...interface{}) {
	switch key.(type) {
	case string:
		if len(data) >= 1 {
			c.TemplateData[key.(string)] = data[0]
		}
		break
	case map[string]interface{}:
		for k, v := range key.(map[string]interface{}) {
			c.TemplateData[k] = v
		}
	}
}

// fetch template and data
func (c *Controller) Fetch(name ...string) {
	currentAction := ""
	if len(name) == 0 {
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		action := runtime.FuncForPC(pc[0])
		fullAction := strings.Split(action.Name(), ".")
		currentAction = strings.ToLower(fullAction[len(fullAction)-1])
	} else {
		currentAction = strings.ToLower(name[0])
	}
	suffix := GetConfig("template.suffix").MustString("html")
	currentAction = strings.TrimSuffix(currentAction, "."+suffix)
	template := ""
	if strings.Contains(currentAction, "/") || strings.Contains(currentAction, ":") {
		template = strings.ReplaceAll(currentAction, ":", "/") + "." + suffix
	} else {
		template = fmt.Sprintf("%s/%s/%s.%s", c.Request.module, c.Request.controller, strings.TrimSuffix(currentAction, "."+suffix), suffix)
	}
	c.Request.Cxt.HTML(200, template, c.TemplateData)
}
