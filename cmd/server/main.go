package main

import (
	"github.com/gin-gonic/gin"

	"gchat/internal"
	"gchat/internal/ws"
)

func main() {
	hub := ws.NewHub()
	router := gin.Default()

	internal.ServeStatic(router)

	router.GET("/ws", hub.ServeWS)

	router.Run()
}
