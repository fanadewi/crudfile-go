package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fanadewi/crudfile-go/app/models"
	"github.com/fanadewi/crudfile-go/config"
	cloud "github.com/fanadewi/go-cloudinary"
	"github.com/gin-gonic/gin"
)

func UploadCloudinary(c *gin.Context) {
	file := models.UploadRequest{}

	if c.ContentType() == "multipart/form-data" {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		files := form.File["files"]
		var response []interface{}
		var invalidFiles []interface{}

		for _, theFile := range files {
			mediaType := mime.TypeByExtension(filepath.Ext(theFile.Filename))
			if !isValidMediaType(mediaType) {
				invalidFiles = append(invalidFiles, map[string]interface{}{"file": theFile.Filename, "error": fmt.Sprintf("Invalid file type(only support image and PDF)")})
				break
			}
		}

		if len(invalidFiles) > 0 {
			c.JSON(http.StatusBadRequest, invalidFiles[0])
			return
		} else if len(files) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Files not found"})
			return
		}

		for _, theFile := range files {
			cloudinaryResponse := uploadFile(file, theFile, c)
			response = append(response, cloudinaryResponse)
		}

		c.JSON(http.StatusOK, response)
	} else {
		decodeFile := &models.UploadRequest{}
		err := json.NewDecoder(c.Request.Body).Decode(decodeFile)
		file = *decodeFile
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := upload(file, c)
		c.JSON(http.StatusOK, response)
	}
}

func uploadFile(file models.UploadRequest, theFile *multipart.FileHeader, c *gin.Context) interface{} {
	fileOpened, err := theFile.Open()
	if err != nil {
		return map[string]interface{}{"status": 500, "message": err.Error()}
	}

	mediaType := mime.TypeByExtension(filepath.Ext(theFile.Filename))

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, fileOpened); err != nil {
		return map[string]interface{}{"status": 500, "message": err.Error()}
	}

	base64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	file.File = fmt.Sprintf("data:%s;base64,%s", mediaType, base64)
	return upload(file, c)
}

func isValidMediaType(mediaType string) bool {
	isValid := false
	if strings.Contains(mediaType, "image") || strings.Contains(mediaType, "pdf") || strings.Contains(mediaType, "plain") {
		isValid = true
	}
	return isValid
}

func upload(file models.UploadRequest, c *gin.Context) interface{} {
	conf := config.New()
	e := cloud.CloudinaryRequest{
		File:        file.File,
		CloudName:   conf.CloudName,
		CloudKey:    conf.CloudKey,
		CloudSecret: conf.CloudSecret,
	}
	return e.Upload()
}
