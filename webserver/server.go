package webserver

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func cut_time(elapsedTime time.Duration) string{
	s := elapsedTime.String()
	sslice:=strings.Split(s, ".") 
		s = s[:len(s)-6]
	if len(sslice) > 1 {
		sslice[1] = sslice[1][:(len(sslice[1])-(len(sslice[1])-2))]
		s = strings.Join(sslice, ".")
	}
	return s
}

func Run() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	
	router.POST("/srv_process", Srv_process)
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}