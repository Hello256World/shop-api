package main

import (
	"github.com/Hello256World/shop-api/database"
	"github.com/Hello256World/shop-api/database/migrate"
	"github.com/Hello256World/shop-api/routes"
	"github.com/Hello256World/shop-api/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	database.Init()
	migrate.Init()
}

func main() {
	server := gin.Default()
	utils.Validation()
	routes.RegisterRouter(server, database.DB)
	server.Run(":8080")
}
