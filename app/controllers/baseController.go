package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/fanadewi/crudfile-go/app/models"
	"github.com/fanadewi/crudfile-go/config"
	cloud "github.com/fanadewi/go-cloudinary"
	"github.com/gin-gonic/gin"
)

func UploadCloudinary(c *gin.Context) {
	conf := config.New()
	file := &models.UploadRequest{}
	err := json.NewDecoder(c.Request.Body).Decode(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e := cloud.CloudinaryRequest{
		File:        file.File,
		CloudName:   conf.CloudName,
		CloudKey:    conf.CloudKey,
		CloudSecret: conf.CloudSecret,
	}
	c.JSON(http.StatusOK, e.Upload())
}
