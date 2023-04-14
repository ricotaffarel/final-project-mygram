package middlewares

import (
	"final-project-mygram/database"
	"final-project-mygram/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		// productId, err := strconv.Atoi(c.Param("productId"))
		// if productId == 0 {

		// } else if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		// 		"error":   "Bad Request",
		// 		"messege": "Invalid parameter",
		// 	})
		// 	return
		// }

		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))
		User := models.User{}

		err := db.First(&User, "id = ?", userID).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusLocked, gin.H{
				"error":   "Data not found",
				"messege": "Data doesn't exist",
			})
			return
		}

		c.Next()
	}
}
