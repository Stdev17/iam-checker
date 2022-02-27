package main

import (
    "context"
    "log"
    "time"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMProfile struct {
    accessKeyId string
    userName string
    createdDate time.Time
}

func main() {
    _, err := FetchIAM()
    if err != nil {
        //
        return
    }

    return
    //for _, obj := range results {
    //    log.Print(obj.(string))
    //}
}

// FetchIAM connects to AWS SDK API with an IAM access key pair and fetch the IAM access key data from its account.
func FetchIAM() ([]IAMProfile, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    client := iam.NewFromConfig(cfg)
    maxItems := int32(1000)
    input := &iam.ListAccessKeysInput{
        MaxItems: aws.Int32(maxItems),
    }

    fetched, err := client.ListAccessKeys(context.TODO(), input)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    var IAMs []IAMProfile

    for _, key := range fetched.AccessKeyMetadata {
        tmp := IAMProfile{
            accessKeyId: *key.AccessKeyId,
            userName: string(*key.UserName),
            createdDate: *key.CreateDate,
        }
        IAMs = append(IAMs, tmp)
    }

    return IAMs, nil
}

// CheckProfileExpired filters out the valid profiles from the originally fetched data.
func CheckProfileExpired(hour time.Duration, given []IAMProfile) ([]IAMProfile) {
    var filtered []IAMProfile
    for _, val := range given {
        if val.createdDate.Add(hour).Before(time.Now()) {
            filtered = append(filtered, val)
        }
    }

    return filtered
}