package routes

import (
	"de.stuttgart.hft/DBS2-Frontend/controllers"
	"github.com/gin-gonic/gin"
)

// Frondend Routes
var RegisterRoutes = func(router *gin.Engine) {
	router.GET("/", controllers.RedirectToIndex)
	router.GET("/photos", controllers.OpenPhotos)
	router.POST("/photos", controllers.OpenPhotosByTypeId)
	router.GET("/photos/all", controllers.OpenPhotos)
	router.POST("/photos/all", controllers.OpenPhotos)
	router.GET("/deletephoto/:id", controllers.DeleteSinglePhoto)
	router.GET("/rolls", controllers.OpenRolls)
	router.POST("/rolls", controllers.OpenRolls)
	router.GET("/albums", controllers.OpenAlbums)
	router.POST("/albums", controllers.OpenAlbums)
	router.POST("/createRoll", controllers.CreateRoll)
	router.POST("/createAlbum", controllers.CreateAlbum)
	router.POST("/uploadPhoto", controllers.UploadPhotos)
	router.GET("/roll/:id", controllers.OpenRollById)
	router.POST("/roll/:id", controllers.OpenRollById)
	router.GET("/album/:id", controllers.OpenAlbumById)
	router.GET("/deleteroll/:id", controllers.DeleteRollAndPhotos)
	router.GET("/deleteAlbum/:id", controllers.DeleteAlbum)
	router.GET("/deletephotopool/:id", controllers.DeletePhotoFromPool)
	router.GET("/showphoto/:uuid", controllers.ShowPhoto)
}
