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

// CreatePhoto godoc
// @Security ApiKeyAuth
// @Summary Create Photo account
// @Description Create Photo account for logged in user
// @Tags Photo
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT authorization token"
// @Param models.Photo body models.Photo true "Photo data to create"
// @Success 201 {object} models.Photo
// @Router /users/photo/create [post]
func CreatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	var User models.User
	contentType := helpers.GetContentType(c)

	Photo := models.Photo{}
	userID := uint(userData["id"].(float64))
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	if User.Role == "user" {
		Photo.UserID = userID
		err = db.Create(&Photo).Error
		Photo.User = &User
	} else {
		err = db.Create(&Photo).Error
	}

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Created photo success", http.StatusCreated, "Success", Photo)
	c.JSON(http.StatusCreated, response)
}

// Update Photo godoc
// @Summary Update Photo
// @Description Update an existing Photo
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param photoId path uint true "Photo ID"
// @Param models.Photo body models.Photo true "Photo information"
// @Success 200 {object} models.SocialMedia
// @Router /users/photo/update/{photoId} [patch]
func UpdatedPhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Photo := models.Photo{}
	User := models.User{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userID := uint(userData["id"].(float64))
	// Product.UserID = userID
	Photo.ID = uint(photoId)
	photo := db.Preload("User").First(&Photo)
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	if userID == Photo.UserID {
		err = photo.Save(&Photo).Error
	} else if User.Role == "admin" {
		err = photo.Save(&Photo).Error
	} else {
		response := helpers.APIResponse("Invalid photo id", http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Update photo success", http.StatusOK, "Success", Photo)
	c.JSON(http.StatusOK, response)
}

// View Photo godoc
// @Summary View photo
// @Description View photo data
// @Tags photo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param photoId path int true "photo ID"
// @Param models.Photo body models.Photo true "photo data"
// @Success 200 {object} models.SocialMedia
// @Router /users/photo/view [get]
func ViewPhoto(c *gin.Context) {
	db := database.GetDB()
	User := models.User{}
	Photo := []models.Photo{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	photoId, _ := strconv.Atoi(c.Query("photoId"))
	userID := uint(userData["id"].(float64))
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	fmt.Println(photoId)
	fmt.Println(User.ID)
	if photoId != 0 {
		err = db.Where("id = ?", photoId).Find(&Photo).Preload("User").Error
	} else if User.Role == "user" {
		err = db.Where("user_id = ?", userID).Find(&Photo).Preload("User").Error
	} else {
		err = db.Find(&Photo).Preload("User").Error
	}

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(Photo) == 0 {
		response := helpers.APIResponse("Data not found", http.StatusNotFound, "Error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
	response := helpers.APIResponse("Get data photo success", http.StatusOK, "Success", Photo)
	c.JSON(http.StatusOK, response)

}

// Delete Photo godoc
// @Summary Delete photo
// @Description Delete photo data
// @Tags photo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param photoId path int true "photo ID"
// @Success 200 {object} models.Photo
// @Router /users/socialmedia/delete/{photoId} [delete]
func DeletedPhoto(c *gin.Context) {
	db := database.GetDB()
	Photo := models.Photo{}
	User := models.User{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	Photo.ID = uint(photoId)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := db.First(&Photo)

	if userID == Photo.UserID {
		err = data.Delete(&Photo).Error
	} else if User.Role == "admin" {
		err = data.Delete(&Photo).Error
	} else {
		response := helpers.APIResponse("Your don't have access to delete this photo", http.StatusBadRequest, "Error Bad Request", nil)
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
