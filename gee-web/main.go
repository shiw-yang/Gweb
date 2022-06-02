package main

import (
	"fmt"
	"gweb"
	"net/http"
)

// gweb 启动入口
func main() {
	r := gweb.New()
	r.GET("/index", func(c *gweb.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gweb.Context) {
			c.HTML(http.StatusOK, "<h1> Hello Gweb</h1>")
		})
		v1.GET("/hello/*", func(c *gweb.Context) {
			// expect /hello?name=ysw
			c.String(http.StatusOK, "hello %s,you are at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("v2")
	v2.GET("/hello/:name", func(c *gweb.Context) {
		// expect /hello/ysw
		c.String(http.StatusOK, "hello %s,you are at %s\n", c.Param("name"), c.Path)
	})
	v2.POST("/login", func(c *gweb.Context) {
		c.Json(http.StatusOK, gweb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	err := r.Run(":9999")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
