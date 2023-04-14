package controllers

import (
	"final-project-mygram/database"
	"final-project-mygram/helpers"
	"final-project-mygram/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

// UserRegister godoc
// @Summary Post details for a given Id
// @Description Register User
// @Tags User
// @Accept json
// @Produce json
// @Param models.User body models.User true "Register User"
// @Success 200 {object} models.User
// @Router /user/register [post]
func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Debug().Create(&User).Error
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusBadRequest, "Error Bad Request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("User Registered", http.StatusCreated, "Success", User)
	c.JSON(http.StatusCreated, response)
}

// UserLogin godoc
// @Summary Post details for a given Id
// @Description Register User
// @Tags User
// @Accept json
// @Produce json
// @Param models.User body models.User true "Login User"
// @Success 200 {object} models.User
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}
	password := ""

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	password = User.Password

	err := db.Debug().Where("email = ?", User.Email).Take(&User).Error

	if err != nil {
		response := helpers.APIResponse("Unauthorized", http.StatusBadRequest, "Invalid email/password", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))
	if !comparePass {
		response := helpers.APIResponse("Unauthorized", http.StatusBadRequest, "Invalid email/password", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email)
	data := map[string]interface{}{
		"token": token,
		"user":  User,
	}

	response := helpers.APIResponse("User Login", http.StatusOK, "Success", data)
	c.JSON(http.StatusOK, response)
}
