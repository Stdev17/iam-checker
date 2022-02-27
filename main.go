package main

import (
    "context"
    "log"
    "time"
    "bufio"
    "io"
    "os"
    "strconv"
    "encoding/json"
    "github.com/joho/godotenv"
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

    err := godotenv.Load(".env")
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

    filtered := CheckProfileExpired(time.Duration(time.Hour * time.Duration(elapsedTime)), fetched)

    if SaveTargetIAMProfiles(filtered) != nil {
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
                UserName: string(*key.UserName),
                CreatedDate: *key.CreateDate,
            }
            IAMs = append(IAMs, tmp)
        }
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