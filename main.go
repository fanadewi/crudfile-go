package main

import (
	"fmt"
	"log"

	c "github.com/fanadewi/crudfile-go/app/controllers"
	u "github.com/fanadewi/crudfile-go/app/utils"
	"github.com/fanadewi/crudfile-go/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	u.EnsureFolderExist("tmp")
	u.EnsureFolderExist("logs")
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %s", err.Error()))
	}
	log.Println("Init success")
}

func main() {
	conf := config.New()
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Use(gin.Logger())
	r.POST("/upload", c.UploadCloudinary)

	r.Run(fmt.Sprintf(":%d", conf.Port))
}
