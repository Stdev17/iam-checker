package main

import (
    "context"
    "log"
    "time"
    "bufio"
    "io"
    "os"
    "encoding/json"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMProfile struct {
    UserName string `json:"userName, string"`
    AccessKeyId string `json:"accessKeyId, string"`
    CreatedDate time.Time `json:"createdDate, string"`
}

func main() {
    _, err := FetchIAM()
    if err != nil {
        log.Fatal(err)
        return
    }

    return
}

// FetchIAM connects to AWS SDK API with an IAM access key pair and fetch the IAM access key data from its account.
func FetchIAM() ([]IAMProfile, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        return nil, err
    }

    client := iam.NewFromConfig(cfg)
    maxItems := int32(1000)
    input := &iam.ListAccessKeysInput{
        MaxItems: aws.Int32(maxItems),
    }

    fetched, err := client.ListAccessKeys(context.TODO(), input)
    if err != nil {
        return nil, err
    }

    var IAMs []IAMProfile

    for _, key := range fetched.AccessKeyMetadata {
        tmp := IAMProfile{
            AccessKeyId: *key.AccessKeyId,
            UserName: string(*key.UserName),
            CreatedDate: *key.CreateDate,
        }
        IAMs = append(IAMs, tmp)
    }

    return IAMs, nil
}

// CheckProfileExpired filters out the valid profiles from the originally fetched data.
func CheckProfileExpired(hour time.Duration, given []IAMProfile) ([]IAMProfile) {
    var filtered []IAMProfile
    for _, val := range given {
        if val.CreatedDate.Add(hour).Before(time.Now()) {
            filtered = append(filtered, val)
        }
    }

    return filtered
}

// SaveTargetIAMProfiles saves a file which contains UserName, Access Key ID, Lifetime of the target IAM profiles.
func SaveTargetIAMProfiles(given []IAMProfile) (error) {
    // marshaling phase
    buf, _ := json.MarshalIndent(given, "", "    ")

    // write phase
    result, err := os.Create("/Users/leta/Github/iam-checker/" + string([]byte(time.Now().String())[:19]) +".json")
    if err != nil {
        return err
    }
    defer result.Close()

    w := bufio.NewWriter(result)
    for _, b := range buf {
        err := w.WriteByte(b)
        if err == io.EOF {
            break
        }
    }
    w.Flush()

    return nil
}