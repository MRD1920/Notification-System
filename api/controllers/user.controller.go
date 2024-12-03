package controllers

import (
	"net/http"

	"github.com/MRD1920/Notification-System/api/service"
	DB "github.com/MRD1920/Notification-System/db"
	model "github.com/MRD1920/Notification-System/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(ctx *gin.Context) {
	var newUser model.User

	err := ctx.BindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Add new UUID to the user
	newUser.Id = uuid.New()

	//Add the user to the Database
	err = service.AddUserToDB(newUser)
	if err != nil {
		if DB.ErrorCode(err) == DB.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "User already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := service.DeleteUserFromDB(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})

}
func GetUsers(ctx *gin.Context) {
	users, err := service.GetAllUsersFromDb()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := service.GetUserFromDB(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
