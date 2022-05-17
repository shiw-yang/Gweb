package main

import (
	"fmt"
	"gweb"
	"net/http"
)

// gweb 启动入口
func main() {
	r := gweb.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	err := r.Run(":9999")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
