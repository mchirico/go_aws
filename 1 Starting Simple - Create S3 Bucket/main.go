package main

import (
	"context"
	"flag"
	"fmt"
	"lesson1/client"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var deleteBucket *bool
var listBuckets *bool

func Options() *string {
	bucketPtr := flag.String("bucket", "", "unique bucket name")
	deleteBucket = flag.Bool("delete", false, "delete flag")
	listBuckets = flag.Bool("list", false, "list buckets")
	flag.Parse()
	if *bucketPtr == "" && ! *listBuckets{
		log.Fatalf(`Need to enter a bucket name.
	
		 -bucket=bucketsample
	`)
	}
	return bucketPtr
}

func CreateBucket(cfg aws.Config, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	client := s3.NewFromConfig(cfg)
	result, err := client.CreateBucket(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteBucket(cfg aws.Config, input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	client := s3.NewFromConfig(cfg)
	result, err := client.DeleteBucket(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ListBuckets(cfg aws.Config, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	client := s3.NewFromConfig(cfg)
	result, err := client.ListBuckets(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func RunDelete(bucket *string) {
	input := &s3.DeleteBucketInput{
		Bucket: bucket,
	}
	result, err := DeleteBucket(client.Config(), input)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	fmt.Println(result)
	return
}

func RunCreate(bucket *string) {
	input := &s3.CreateBucketInput{
		Bucket: bucket,
	}
	result, err := CreateBucket(client.Config(), input)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	fmt.Printf("result: %s", result)
}

func RunList() {
	input := &s3.ListBucketsInput{}
	result, err := ListBuckets(client.Config(), input)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	for _,v := range result.Buckets {

		fmt.Printf("%s, \t\t%s\n", *v.Name,*v.CreationDate)
	}
}



func main() {
	bucket := Options()
	if *deleteBucket {
		RunDelete(bucket)
		return
	}
	if *listBuckets {
		RunList()
		return
	}
	RunCreate(bucket)
}
