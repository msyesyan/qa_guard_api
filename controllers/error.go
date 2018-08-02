package controllers

import (
    "github.com/astaxie/beego"
)

// ErrorController handle general errors abort
type ErrorController struct {
    beego.Controller
}

// Error401 handle 401 Unauthorized error
func (c *ErrorController) Error401() {
		c.Data["json"] = "Not authorized"
		c.ServeJSON()
}
