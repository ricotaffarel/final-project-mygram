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
// @host final-project-mygram-production.up.railway.app
// @BasePath /
func StartApp() *gin.Engine {
	r := gin.Default()

	// READ
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	user := r.Group("/user")
	{
		user.POST("/register", controllers.UserRegister)
		user.POST("/login", controllers.UserLogin)
	}

	userRouter := r.Group("/users")
	{
		userRouter.Use(middlewares.Authentication())
		//Photo
		userPhoto := userRouter.Group("/photo", middlewares.UserAuthorization())
		userPhoto.POST("/create", controllers.CreatePhoto)
		userPhoto.PUT("/update/:photoId", controllers.UpdatedPhoto)
		userPhoto.GET("/view", controllers.ViewPhoto)
		userPhoto.DELETE("/delete/:photoId", controllers.DeletedPhoto)

		//Social Media
		userSocialMedia := userRouter.Group("/socialmedia", middlewares.UserAuthorization())
		userSocialMedia.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
		userSocialMedia.POST("/create", controllers.CreateSocialMedia)
		userSocialMedia.PUT("/update/:socialMediaId", controllers.UpdatedSocialMedia)
		userSocialMedia.GET("/view", controllers.ViewSocialMedia)
		userSocialMedia.DELETE("/delete/:socialMediaId", controllers.DeletedSocialMedia)

		//Comment
		userComment := userRouter.Group("/comment", middlewares.Authentication())
		userComment.POST("/create", controllers.CreateComment)
		userComment.PUT("/update/:commentId", controllers.UpdatedComment)
		userComment.GET("/view", controllers.ViewComment)
		userComment.DELETE("/delete/:commentId", controllers.DeletedComment)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
