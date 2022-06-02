package main

import (
	"fmt"
	"gee-web/pojo"
	"gweb"
	"gweb/middlewares"
	"html/template"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gweb.HandlerFunc {
	return func(c *gweb.Context) {
		// Start timer
		t := time.Now()
		// if a server err occurred
		if c.Params == nil {
			c.Status(http.StatusInternalServerError)
			c.Fail(http.StatusInternalServerError, "Internal Server Error")
		}
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", http.StatusOK, c.Req.RequestURI, time.Since(t))
	}
}

// FormatAsDate is an overwriting func.
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return fmt.Sprintf("%d-%02d-%02d\t%02d:%02d:%02d", year, month, day, hour, min, sec)
}

// gweb 启动入口
func main() {
	r := gweb.New()
	r.Use(middlewares.Logger()) // global middleware
	// use Template func
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	// define
	stu1 := &pojo.Student{Name: "shiwyang", Age: 21}
	stu2 := &pojo.Student{Name: "Jack", Age: 50}
	r.GET("/", func(c *gweb.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gweb.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gweb.H{
			"title":  "gweb",
			"stuArr": [2]*pojo.Student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *gweb.Context) {
		year, month, day := time.Now().Date()
		hour, min, sec := time.Now().Clock()
		c.HTML(http.StatusOK, "custom_func.tmpl", gweb.H{
			"title": "gweb",
			"now":   time.Date(year, month, day, hour, min, sec, 0, time.UTC),
		})
	})
	// use group func
	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
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
