package gweb

import (
	"log"
	"time"
)

// Logger as a middleware example
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time

		log.Printf("%s\t[%d]\t%s in %s\n", c.Method, c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
