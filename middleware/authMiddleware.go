package middleware

import (
	helper "golang-techque/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")

		if clientToken == "" {
			c.JSON(403, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)

		if err != "" {
			c.JSON(403, gin.H{"error": "Token is not valid"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		c.Next()

	}
}
