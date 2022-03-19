package main

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *GzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func Compress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.GetHeader("Accept-Encoding"), "gzip") {
			gz, err := gzip.NewWriterLevel(ctx.Writer, gzip.BestSpeed)

			if err != nil {
				return
			}

			defer gz.Close()

			ctx.Header("Content-Encoding", "gzip")
			ctx.Writer = &GzipWriter{ctx.Writer, gz}
		}

		ctx.Next()
	}
}

func Unpack() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.GetHeader("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(ctx.Request.Body)

			if err != nil {
				ctx.Status(http.StatusInternalServerError)

				return
			}

			defer gz.Close()
			ctx.Request.Body = gz
		}

		ctx.Next()
	}
}

func Tokenize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")

		if err != nil {
			setTokenCookie(ctx)
		} else {
			if !IsValid(token) {
				setTokenCookie(ctx)
			}
		}

		ctx.Next()
	}
}

func setTokenCookie(ctx *gin.Context) {
	token, err := Generate()

	if err != nil {
		return
	}

	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
}
