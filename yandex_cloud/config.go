package yandex_cloud

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"gopkg.in/ini.v1"
)

var (
	bucket             string
	awsAccessKeyId     string
	awsSecretAccessKey string
	region             string
	endpointUrl        string
	client             *s3.Client
)

func isEmpty(str string) bool {
	return str == ""
}

func getConfLocation() string {
	home, err := os.UserHomeDir()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	location := home + "/.config/cloudphoto/cloudphotorc/config.ini"
	return location
}

func LoadFromConfigFile() {
	cfg, err := ini.Load(getConfLocation())
	if err != nil {
		println("Failed to read config file")
		os.Exit(1)
	}
	section := cfg.Section("DEFAULT")
	bucket = section.Key("bucket").String()
	awsAccessKeyId = section.Key("aws_access_key_id").String()
	awsSecretAccessKey = section.Key("aws_secret_access_key").String()
	region = section.Key("region").String()
	endpointUrl = section.Key("endpoint_url").String()

	if isEmpty(bucket) || isEmpty(awsAccessKeyId) || isEmpty(awsSecretAccessKey) || isEmpty(region) || isEmpty(endpointUrl) {
		println("Not all credentials are given")
		os.Exit(1)
	}
}

func CreateConfigFile(bucketName, keyId, accessKey string) error {
	cfg := ini.Empty()
	cfg.Section("DEFAULT").Key("bucket").SetValue(bucketName)
	cfg.Section("DEFAULT").Key("aws_access_key_id").SetValue(keyId)
	cfg.Section("DEFAULT").Key("aws_secret_access_key").SetValue(accessKey)
	cfg.Section("DEFAULT").Key("region").SetValue("ru-central1")
	cfg.Section("DEFAULT").Key("endpoint_url").SetValue("https://storage.yandexcloud.net")

	err := os.Remove(getConfLocation())
	if err != nil {
		return err
	}

	err = cfg.SaveTo(getConfLocation())
	if err != nil {
		return err
	}

	return nil
}

func Init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyId, awsSecretAccessKey, "")),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == s3.ServiceID && region == "ru-central1" {
					return aws.Endpoint{
						PartitionID:   "yc",
						URL:           endpointUrl,
						SigningRegion: "ru-central1",
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			})),
	)
	if err != nil {
		println("Init error")
		os.Exit(1)
	}

	// Создаем клиента для доступа к хранилищу S3
	client = s3.NewFromConfig(cfg)
}
