package yandex_cloud

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"os"
	"strings"
)

func IsBucketExist() bool {
	input := &s3.HeadBucketInput{
		Bucket: &bucket,
	}

	_, err := client.HeadBucket(context.TODO(), input)
	if err != nil {
		return false
	} else {
		return true
	}
}

func CreateBucket() error {
	_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    types.BucketCannedACLPublicRead,
		//todo: set acl public in mksit funcion
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		return err
	}
	//_, err := client.PutPublicAccessBlock(context.TODO(), &s3.PutPublicAccessBlockInput{
	//	Bucket: aws.String(bucketName),
	//	PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
	//		BlockPublicAcls:       false,
	//		BlockPublicPolicy:     false,
	//		IgnorePublicAcls:      false,
	//		RestrictPublicBuckets: false,
	//	},
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}

func GetAllRootFolders() ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(""),
		Delimiter: aws.String("/"),
	}

	// Вызываем операцию ListObjectsV2 для получения списка объектов и папок
	resp, err := client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	folders := resp.CommonPrefixes

	var result []string

	for _, folder := range folders {
		s := *folder.Prefix
		result = append(result, s[:len(s)-1])
	}
	return result, nil
}

func CreateFolderIfNotExists(folderName string) {
	folderName = folderName + "/"
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(folderName),
	}
	_, err := client.HeadObject(context.TODO(), input)
	if err != nil {
		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &folderName,
		})
	}
}

func UploadImagesToAlbum(folderName string, imagePaths []string) {
	for _, filename := range imagePaths {
		uploadFile(client, bucket, filename, folderName)
	}
}

func uploadFile(client *s3.Client, bucket string, filename string, fileAlbum string) {
	file, err := os.Open(filename)
	if err != nil {
		println("Warning: Photo not sent", filename)
		return
	}
	defer file.Close()

	// Извлекаем имя файла из пути
	key := filename
	if idx := strings.LastIndex(key, "/"); idx != -1 {
		key = key[idx+1:]
	}

	key = fileAlbum + "/" + key

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		println("Warning: Photo not sent", filename)
		return
	}
}
