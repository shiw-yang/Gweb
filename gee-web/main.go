package main

import (
	"fmt"
	"gweb"
	"net/http"
)

// gweb 启动入口
func main() {
	r := gweb.New()
	r.GET("/", func(c *gweb.Context) {
		c.HTML(http.StatusOK, "<h1>hello GWeb</h1>")
	})

	r.GET("/hello", func(c *gweb.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gweb.Context) {
		// expect /hello/ysw
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.POST("/login", func(c *gweb.Context) {
		c.Json(http.StatusOK, gweb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/assets/*filepath", func(c *gweb.Context) {
		c.Json(http.StatusOK, gweb.H{"filepath": c.Param("filepath")})
	})

	err := r.Run(":9999")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
