package internal

import (
	"github.com/gin-gonic/gin"
)

/*
Returns static index.html and assets like js and css
Mainly for SPA
*/
func ServeStatic(router *gin.Engine) {
	router.Static("/assets", "./static/assets")

	router.NoRoute(func(ctx *gin.Context) {
		ctx.File("./static/index.html")
	})
}
