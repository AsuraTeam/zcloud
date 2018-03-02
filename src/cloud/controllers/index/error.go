package index

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error403() {
	c.TplName = "error/403.html"
}

func (c *ErrorController) Error404() {
	c.TplName = "error/404.html"
}

func (c *ErrorController) Error500() {
	c.TplName = "error/500.html"
}

func (c *ErrorController) Error503() {
	c.TplName = "error/503.html"
}

func (c *ErrorController) Error400() {
	c.TplName = "error/400.html"
}