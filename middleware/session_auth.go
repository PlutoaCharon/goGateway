package middleware

import (
	"errors"
	"github.com/PlutoaCharon/goGateway/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get(public.AdminInfoSessionKey) == nil {
			ResponseError(c, 200, errors.New("管理端未登陆"))
			c.Abort()
			return
		}
		c.Next()
	}
}
