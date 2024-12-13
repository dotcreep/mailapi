package utils

import (
	"log"
	"os"
)

func InitStorage() {
	var storagePath string
	err := os.Getenv("STORAGE")
	if err == "" {
		log.Println(err)
	}
	storagePath = os.Getenv("STORAGE")
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Println(err)
		}
	}
}
