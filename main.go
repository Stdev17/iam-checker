package main

import (
    "context"
    "fmt"
    "log"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/iam"
)

func main() {
    err := FetchIAM()
    if err != nil {
        //
        return
    }
    //for _, obj := range results {
    //    log.Print(obj.(string))
    //}
}

// FetchIAM connects to AWS SDK API with an IAM access key pair and fetch the IAM access key data from its account.
func FetchIAM() (error) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
        return err
    }

    client := iam.NewFromConfig(cfg)
    maxItems := int32(1000)
    input := &iam.ListAccessKeysInput{
        MaxItems: aws.Int32(maxItems),
    }

    fetched, err := client.ListAccessKeys(context.TODO(), input)
    if err != nil {
        log.Fatal(err)
        return err
    }

    for _, key := range fetched.AccessKeyMetadata {
        fmt.Println("Status for access key " + *key.AccessKeyId + ": " + string(*key.UserName) + ", " + (*key.CreateDate).String())
    }

    return nil
}