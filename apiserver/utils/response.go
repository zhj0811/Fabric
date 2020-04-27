package utils

import (
	"github.com/gin-gonic/gin"
)

// Response http response
func Response(body []byte, c *gin.Context, status int) {
	//c.Writer.Header().Set("version", c.Request.Header.Get("version"))
	c.Writer.Header().Set("content-Type", c.Request.Header.Get("content-Type"))
	//c.Writer.Header().Set("trackId", c.Request.Header.Get("trackId"))
	//c.Writer.Header().Set("language", c.Request.Header.Get("language"))
	//c.Writer.Header().Set("page", string(jsonPage))
	c.Writer.WriteHeader(status)
	c.Writer.Write(body)
}
