package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type IAMProfile struct {
	UserName    string    `json:"userName,string"`
	AccessKeyId string    `json:"accessKeyId,string"`
	CreatedDate time.Time `json:"createdDate,string"`
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = godotenv.Load(filepath.Join(cwd, ".env"))
	if err != nil {
		log.Fatalf("Error loading .env file")
		return
	}

	elapsedTime, err := strconv.Atoi(os.Getenv("LIFETIME"))
	if err != nil {
		log.Fatal(err)
		return
	}

	fetched, err := FetchIAM()
	if err != nil {
		log.Fatal(err)
		return
	}

	filtered := CheckProfileExpired(time.Duration(time.Hour*time.Duration(elapsedTime)), fetched)

	if SaveTargetIAMProfiles(filtered) != nil {
		log.Fatal(err)
		return
	}

	return
}
