package utils

import (
	"fmt"
	"log"
	"os"
)

func EnsureFolderExist(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println(fmt.Sprintf("%s does not exist. creating folder ...", path))
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
