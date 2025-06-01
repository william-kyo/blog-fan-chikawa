package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// print panic error
				fmt.Printf("panic recovered: %v\n%s\n", err, debug.Stack())

				// return 500 error
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "Internal Server Error",
				})

				// abort request
				c.Abort()
			}
		}()

		// continue request
		c.Next()
	}
}
