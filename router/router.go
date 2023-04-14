package router

import (
	"final-project-mygram/controllers"
	"final-project-mygram/middlewares"
	"net/http"

	_ "final-project-mygram/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"
)

// @title Final project Api
// @version 1.1
// @description This is a sample service for managing books
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email rico@gmail.com
// @lisence.name Apache 2.0
// @lisence.url https://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func StartApp() *gin.Engine {
	r := gin.Default()

	// READ
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	user := r.Group("/user")
	{
		// POST
		user.POST("/register", controllers.UserRegister)
		// POST
		user.POST("/login", controllers.UserLogin)
	}

	userRouter := r.Group("/users")
	{
		userRouter.Use(middlewares.Authentication())
		userPhoto := r.Group("/photo")
		userPhoto.Use(middlewares.UserAuthorization())
		// POST
		userPhoto.POST("/create", controllers.CreatePhoto)
		// UPDATE
		userPhoto.PUT("/update/:photoId", controllers.UpdatedPhoto)
		// READ
		userPhoto.GET("/view", controllers.ViewPhoto)
		// DELETE
		userPhoto.DELETE("/delete/:photoId", controllers.DeletedPhoto)

		userSocialMedia := r.Group("/socialmedia")
		userSocialMedia.Use(middlewares.UserAuthorization())
		// POST
		userSocialMedia.POST("/create", controllers.CreateSocialMedia)
		// UPDATE
		userSocialMedia.PUT("/update/:socialMediaId", controllers.UpdatedSocialMedia)
		// READ
		userSocialMedia.GET("/view", controllers.ViewSocialMedia)
		// DELETE
		userSocialMedia.DELETE("/delete/:socialMediaId", controllers.DeletedSocialMedia)

		userComment := r.Group("/comment")
		userComment.Use(middlewares.UserAuthorization())
		// POST
		userComment.POST("/create", controllers.CreateComment)
		// UPDATE
		userComment.PUT("/update/:commentId", controllers.UpdatedComment)
		// READ
		userComment.GET("/view", controllers.ViewComment)
		// DELETE
		userComment.DELETE("/delete/:commentId", controllers.DeletedComment)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
