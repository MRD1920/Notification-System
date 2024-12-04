package controllers

import (
	"net/http"

	model "github.com/MRD1920/Notification-System/models"
	"github.com/MRD1920/Notification-System/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateNotification(ctx *gin.Context) {
	var notification model.Notification
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Validate notification
	if err := utils.ValidateNotification(notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification.Id = uuid.New()
	notification.Status = "pending"

	//Save notification
	if err := utils.SaveNotification(notification); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	//Return success
	ctx.JSON(http.StatusOK, gin.H{"message": "Notification created successfully"})
}
