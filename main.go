package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

type IAMProfile struct {
	UserName    string    `json:"userName,string"`
	AccessKeyId string    `json:"accessKeyId,string"`
	CreatedDate time.Time `json:"createdDate,string"`
}

func main() {

	//err = godotenv.Load(filepath.Join(cwd, ".env"))
	//if err != nil {
	//	log.Fatalf("Error loading .env file")
	//	return
	//}

	lifetime := os.Getenv("LIFETIME")

	var elapsedTime int
	var err error

	if lifetime != "" {
		elapsedTime, err = strconv.Atoi(lifetime)
		if err != nil {
			log.Fatal(err)
			return
		}
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
