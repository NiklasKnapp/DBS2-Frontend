package main

import (
	"de.stuttgart.hft/DBS2-Frontend/routes"
	"github.com/gin-gonic/gin"
)

//Main
func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*.html")

	routes.RegisterRoutes(r)

	r.Run(":8081")
}
