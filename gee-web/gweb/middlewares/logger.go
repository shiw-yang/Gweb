package middlewares

import (
	"gweb"
	"log"
	"time"
)

// Logger as a middleware example
func Logger() gweb.HandlerFunc {
	return func(c *gweb.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Println("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
