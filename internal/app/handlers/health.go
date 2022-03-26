package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type Health struct {
	DatabaseURL string
}

func (health Health) CheckDatabase(ctx *gin.Context) {
	conn, err := pgx.Connect(context.Background(), health.DatabaseURL)

	if err == nil {
		ctx.Status(http.StatusOK)
	} else {
		ctx.Status(http.StatusInternalServerError)
	}

	defer conn.Close(context.Background())

}
