package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func bootstrap(router *gin.Engine) {
	fmt.Println("Bootstrap - Starting app...")

	application := resolveGameController()
	mapUrlsToControllers(router, application)

	fmt.Println("Bootstrap - Application is up")
}
