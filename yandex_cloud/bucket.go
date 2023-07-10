package yandex_cloud

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func IsBucketExist(bucketName string) bool {
	input := &s3.HeadBucketInput{
		Bucket: &bucketName,
	}

	_, err := client.HeadBucket(context.TODO(), input)
	if err != nil {
		return false
	} else {
		return true
	}
}

func CreateBucket(bucketName string) error {
	_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		ACL:    types.BucketCannedACLPublicRead,
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
