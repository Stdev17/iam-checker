package main

import (
	"os"
	"context"
	"path/filepath"
	"github.com/joho/godotenv"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// FetchIAM connects to AWS SDK API with an IAM access key pair and fetch the IAM access key data from its account.
func FetchIAM() ([]IAMProfile, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = godotenv.Load(filepath.Join(cwd, ".env"))
	if err != nil {
		return nil, err
	}

	//key_id := os.Getenv("AWS_ACCESS_KEY_ID")
	//log.Println(key_id)
	//key_secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	var IAMs []IAMProfile
	var iam_users []string

	client := iam.NewFromConfig(cfg)
	maxItems := int32(1000)

	// fetch IAM Users first
	input := &iam.ListUsersInput{
		MaxItems: aws.Int32(maxItems),
	}

	list_users, err := client.ListUsers(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	for _, u := range list_users.Users {
		iam_users = append(iam_users, *u.UserName)
	}

	// now we fetch IAM profiles
	for _, u := range iam_users {
		input := &iam.ListAccessKeysInput{
			MaxItems: aws.Int32(maxItems),
			UserName: aws.String(u),
		}

		fetched, err := client.ListAccessKeys(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		for _, key := range fetched.AccessKeyMetadata {
			tmp := IAMProfile{
				AccessKeyId: *key.AccessKeyId,
				UserName:    string(*key.UserName),
				CreatedDate: *key.CreateDate,
			}
			IAMs = append(IAMs, tmp)
		}
	}

	return IAMs, nil
}
