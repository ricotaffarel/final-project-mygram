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

// SocialMediaCreate godoc
// @Summary Post details for a given Id
// @Description Create Social Media User
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param models.SocialMedia body models.SocialMedia true "Create SocialMedia"
// @Success 200 {object} models.SocialMedia
// @Router /users/socialmedia/create [post]
func CreateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	SocialMedia := models.SocialMedia{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	User := models.User{}
	User.ID = userID

	err := db.First(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if User.Role == "user" {
		SocialMedia.UserID = userID
	}

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	err = db.Create(&SocialMedia).Error

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Created Social Media success", http.StatusCreated, "Success", SocialMedia)
	c.JSON(http.StatusCreated, response)
	// c.JSON(http.StatusAccepted, "ok")
}

// update social media
// SocialMediaUpdate godoc
// @Summary Post details for a given Id
// @Description Update Social Media User
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param socialmediaId path int true "social media id"
// @Param models.SocialMedia body models.SocialMedia true "Update SocialMedia"
// @Success 200 {object} models.SocialMedia
// @Router /users/socialmedia/update/{socialMediaId} [patch]
func UpdatedSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}
	User := models.User{}

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	userID := uint(userData["id"].(float64))
	// Product.UserID = userID
	SocialMedia.ID = uint(socialMediaId)
	socialMedia := db.Preload("User").First(&SocialMedia)
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	User.Password = ""

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	if userID == SocialMedia.UserID {
		err = socialMedia.Save(&SocialMedia).Error
	} else if User.Role == "admin" {
		err = socialMedia.Save(&SocialMedia).Error
	} else {
		response := helpers.APIResponse("Invalid social media id", http.StatusBadRequest, "Unauthorize", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Update social media success", http.StatusOK, "Success", SocialMedia)
	c.JSON(http.StatusOK, response)
}

// SocialMediaView godoc
// @Summary Post details for a given Id
// @Description View Social Media User
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Query socialmediaId true "social media id"
// @Param models.SocialMedia body models.SocialMedia true "View SocialMedia"
// @Success 200 {object} models.SocialMedia
// @Router /users/socialmedia/view [get]
func ViewSocialMedia(c *gin.Context) {
	db := database.GetDB()
	User := models.User{}
	SocialMedia := []models.SocialMedia{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	socialMediaId, _ := strconv.Atoi(c.Query("socialMediaId"))
	userID := uint(userData["id"].(float64))
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	User.Password = ""

	fmt.Println(socialMediaId)
	fmt.Println(User.ID)
	if socialMediaId != 0 {
		err = db.Where("id = ?", socialMediaId).Find(&SocialMedia).Preload("User").Error
		if int(userID) != int(SocialMedia[0].UserID){
			response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	} else if User.Role == "user" {
		err = db.Where("user_id = ?", userID).Find(&SocialMedia).Preload("User").Error
	} else {
		err = db.Find(&SocialMedia).Preload("User").Error
	}

	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(SocialMedia) == 0 {
		response := helpers.APIResponse("Data not found", http.StatusNotFound, "Error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
	response := helpers.APIResponse("Get data social media success", http.StatusOK, "Success", SocialMedia)
	c.JSON(http.StatusOK, response)

}

// SocialMediaDelete godoc
// @Summary Delete details for a given Id
// @Description Update Social Media User
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param socialmediaId path int true "social media id"
// @Success 200 {object} models.SocialMedia
// @Router /users/socialmedia/delete/{socialmediaId} [delete]
func DeletedSocialMedia(c *gin.Context) {
	db := database.GetDB()
	SocialMedia := models.SocialMedia{}
	User := models.User{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	SocialMedia.ID = uint(socialMediaId)
	userID := uint(userData["id"].(float64))
	err := db.Where("id = ?", userID).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Invalid user id", http.StatusBadRequest, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := db.First(&SocialMedia)
	if SocialMedia.Name == "" {
		response := helpers.APIResponse("id not found", http.StatusNotFound, "Unauthorized", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if userID == SocialMedia.UserID {
		err = data.Delete(&SocialMedia).Error
	} else if User.Role == "admin" {
		err = data.Delete(&SocialMedia).Error
	} else {
		response := helpers.APIResponse("Your don't have access to delete this social media", http.StatusBadRequest, "Error Bad Request", nil)
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
