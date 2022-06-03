package gweb

import (
	"log"
	"time"
)

// TODO: add color support
const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
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
