package main

import (
	"github.com/gin-gonic/gin"
)

func Tokenize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		var realTokenForContext string

		if err != nil {
			/// Токена нет, генерим новый
			encryptedToken, realToken, err := Generate()

			if err == nil {
				setTokenCookie(ctx, encryptedToken)
				realTokenForContext = realToken
			}
		} else {
			realToken, err := GetRealTokenIfValid(token)

			if err != nil {
				// Токен не валидный, так что генерим новый
				encryptedToken, realToken, err := Generate()

				if err != nil {
					setTokenCookie(ctx, encryptedToken)
					realTokenForContext = realToken
				}
			} else {
				realTokenForContext = realToken
			}
		}

		ctx.Set("token", realTokenForContext)
		ctx.Next()
	}
}

func setTokenCookie(ctx *gin.Context, encryptedToken string) {
	ctx.SetCookie("token", encryptedToken, 3600, "/", "localhost", false, true)
}
