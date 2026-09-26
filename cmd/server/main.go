package main

import (
	"github.com/gin-gonic/gin"

	"gchat/internal"
)

func main() {
	router := gin.Default()

	internal.ServeStatic(router)

	router.Run()
}
