package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lerenn/cryptellation/v1/web/ui"
)

func main() {
	router := gin.Default()
	ui.AddRoutes(router)
	router.Run(":8080")
}
