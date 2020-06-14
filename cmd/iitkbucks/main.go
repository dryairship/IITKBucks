package main

import (
	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/controllers"
)

func main() {
	router := gin.Default()
	controllers.SetUpRoutes(router)

	err := router.Run(":8001")
	if err != nil {
		panic(err)
	}
}
