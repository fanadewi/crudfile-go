package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/fanadewi/crudfile-go/app/models"
	"github.com/fanadewi/crudfile-go/config"
	cloud "github.com/fanadewi/go-cloudinary"
	"github.com/gin-gonic/gin"
	goCache "github.com/patrickmn/go-cache"
)

func UploadCloudinary(c *gin.Context) {
	cache := goCache.New(5*time.Minute, 10*time.Minute)
	if c.ContentType() == "multipart/form-data" {
		traceId := c.PostForm("traceId")
		cacheFound, found := cache.Get(traceId)
		if found {
			c.JSON(http.StatusOK, cacheFound)
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		httpStatus, response, err := filesUploader(form)
		if err != nil {
			c.JSON(httpStatus, gin.H{"error": err.Error()})
			return
		}

		c.Set(traceId, response)
		c.JSON(httpStatus, response)
	} else {
		decodeFile := &models.UploadRequest{}
		err := json.NewDecoder(c.Request.Body).Decode(decodeFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cacheFound, found := cache.Get(decodeFile.TraceId)
		if found {
			c.JSON(http.StatusOK, cacheFound)
			return
		}

		response := upload(decodeFile.File)
		c.Set(decodeFile.TraceId, response)
		c.JSON(http.StatusOK, response)
	}
}

func filesUploader(form *multipart.Form) (int, interface{}, error) {
	files := form.File["files"]
	if len(files) < 1 {
		return http.StatusBadRequest, nil, errors.New("File not found")
	}
	if !filesAreValid(files) {
		return http.StatusBadRequest, nil, errors.New("Only upload images and pdf!")
	}

	var response []interface{}
	for _, theFile := range files {
		cloudinaryResponse := uploadFile(theFile)
		response = append(response, cloudinaryResponse)
	}

	return http.StatusOK, response, nil
}

func uploadFile(theFile *multipart.FileHeader) interface{} {
	var file models.UploadRequest
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
	return upload(file.File)
}

func filesAreValid(files []*multipart.FileHeader) bool {
	isValid := true
	for _, theFile := range files {
		mediaType := mime.TypeByExtension(filepath.Ext(theFile.Filename))
		if !(strings.Contains(mediaType, "image") || strings.Contains(mediaType, "pdf")) {
			isValid = false
			break
		}
	}
	return isValid
}

func upload(file string) interface{} {
	conf := config.New()
	e := cloud.CloudinaryRequest{
		File:        file,
		CloudName:   conf.CloudName,
		CloudKey:    conf.CloudKey,
		CloudSecret: conf.CloudSecret,
	}
	return e.Upload()
}
