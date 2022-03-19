package main

import (
	"compress/gzip"
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
