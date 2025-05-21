package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/cmd/container"
	"github.com/sanda-bunescu/ExploRO/common"
	"github.com/sanda-bunescu/ExploRO/initializers"
	"github.com/sanda-bunescu/ExploRO/routes"
)

func init() {
	initializers.LoadEnvFiles()
	initializers.FirebaseInitialization()
	initializers.ConnectToDB()
	initializers.MigrateDB()
	common.DataSeeder(initializers.Database)
}

func main() {
	// Initialize the container that handles dependencies
	c, err := container.NewContainer()
	if err != nil {
		panic(err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Register routes and inject the container (with services and repositories)
	routes.RegisterRoutes(r, c)

	// Run the Gin server
	r.Run()
}
