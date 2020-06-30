package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/config"
	"github.com/dryairship/IITKBucks/controllers"
)

func main() {
	router := gin.Default()
	controllers.SetUpRoutes(router)
	controllers.PerformInitialization()

	err := router.Run(fmt.Sprintf(":%s", config.PORT))
	if err != nil {
		panic(err)
	}
}
