package utils

import (
	"log"
	"os"
)

func InitStorage() {
	cfg, err := OpenYAML()
	if err != nil {
		log.Println(err)
	}
	directory := cfg.Config.Storage
	if directory == "" {
		log.Println("Storage path is required")
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
}
