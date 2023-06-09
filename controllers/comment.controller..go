package controllers

import (
	"final-project-mygram/database"
	"final-project-mygram/helpers"
	"final-project-mygram/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// CreateComment godoc
// @Security ApiKeyAuth
// @Summary Create Comment account
// @Description Create Comment account for logged in user
// @Tags Comment
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT authorization token"
// @Param models.Comment body models.Comment true "Comment data to create"
// @Success 201 {object} models.Comment
// @Router /users/comment/create [post]
func CreateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	var User models.User
	var Photo models.Photo
	contentType := helpers.GetContentType(c)

	Comment := models.Comment{}
	userID := uint(userData["id"].(float64))
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}
	Comment.UserID = userID

	err = db.Create(&Comment).Error
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	db.Where("user_id= ?", Comment.PhotoID).Take(&Photo)
	Comment.Photo = &Photo

	response := helpers.APIResponse("Created Comment success", http.StatusCreated, "Success", Comment)
	c.JSON(http.StatusCreated, response)
}

// UpdateComment godoc
// @Summary Update Comment
// @Description Update an existing Comment
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param commentId path uint true "Comment ID"
// @Param models.Comment body models.Comment true "Comment information"
// @Success 200 {object} models.SocialMedia
// @Router /users/comment/update/{commentId} [patch]
func UpdatedComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Comment := models.Comment{}
	User := models.User{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := uint(userData["id"].(float64))
	// Product.UserID = userID
	Comment.ID = uint(commentId)
	comment := db.Preload("Photo").First(&Comment)
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	if userID == Comment.UserID {
		err = comment.Save(&Comment).Error
	} else if User.Role == "admin" {
		err = comment.Save(&Comment).Error
	} else {
		response := helpers.APIResponse("Invalid comment id", http.StatusBadRequest, "Unauthorize", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Update comment success", http.StatusOK, "Success", Comment)
	c.JSON(http.StatusOK, response)
}

// ViewComment godoc
// @Summary View Comment
// @Description View Comment data
// @Tags Comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param commentId path int true "Comment ID"
// @Param models.Comment body models.Comment true "Comment data"
// @Success 200 {object} models.SocialMedia
// @Router /users/comment/view [get]
func ViewComment(c *gin.Context) {
	db := database.GetDB()
	User := models.User{}
	Comment := []models.Comment{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	commentId, _ := strconv.Atoi(c.Query("commentId"))
	userID := uint(userData["id"].(float64))
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	fmt.Println(commentId)
	fmt.Println(User.ID)
	if commentId != 0 {
		err = db.Where("id = ?", commentId).Find(&Comment).Preload("Photo").Error
		if int(userID) != int(Comment[0].UserID) {
			response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	} else if User.Role == "user" {
		err = db.Where("user_id = ?", userID).Find(&Comment).Preload("Photo").Error
	} else {
		err = db.Find(&Comment).Preload("User").Error
	}

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(Comment) == 0 {
		response := helpers.APIResponse("Data not found", http.StatusNotFound, "Error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
	response := helpers.APIResponse("Get data photo success", http.StatusOK, "Success", Comment)
	c.JSON(http.StatusOK, response)

}

// CommentDelete godoc
// @Summary Delete details for a given Id
// @Description Update Social Media User
// @Tags Comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Query commentId true "comment id"
// @Success 200 {object} models.Comment
// @Router /users/comment/delete/{commentId} [delete]
func DeletedComment(c *gin.Context) {
	db := database.GetDB()
	Comment := models.Comment{}
	User := models.User{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	Comment.ID = uint(commentId)
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := db.First(&Comment)

	if userID == Comment.UserID {
		err = data.Delete(&Comment).Error
	} else if User.Role == "admin" {
		err = data.Delete(&Comment).Error
	} else {
		response := helpers.APIResponse("Your don't have access to delete this comment", http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Product has be deleted", http.StatusBadRequest, "Success", nil)
	c.JSON(http.StatusBadRequest, response)

}
