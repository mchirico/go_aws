package bucket

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	Name *string
}

func NewBucket(name ...string) *Bucket {
	if len(name) >= 1 {
		return &Bucket{Name: &name[0]}
	}
	return &Bucket{}
}

func (b *Bucket) SetName(name string) {
	b.Name = &name
}

func (b *Bucket) createBucket(cfg aws.Config, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	client := s3.NewFromConfig(cfg)
	result, err := client.CreateBucket(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Bucket) deleteBucket(cfg aws.Config, input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	client := s3.NewFromConfig(cfg)
	result, err := client.DeleteBucket(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Bucket) listBuckets(cfg aws.Config, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	client := s3.NewFromConfig(cfg)
	result, err := client.ListBuckets(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Bucket) RunDelete(cfg aws.Config) *s3.DeleteBucketOutput {
	input := &s3.DeleteBucketInput{
		Bucket: b.Name,
	}
	result, err := b.deleteBucket(cfg, input)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	return result
}

func (b *Bucket) RunCreate(cfg aws.Config) {
	input := &s3.CreateBucketInput{
		Bucket: b.Name,
	}
	result, err := b.createBucket(cfg, input)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	fmt.Printf("result: %v", result)
}

func (b *Bucket) RunList(cfg aws.Config) {
	input := &s3.ListBucketsInput{}
	result, err := b.listBuckets(cfg, input)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	for _, v := range result.Buckets {

		fmt.Printf("%s, \t\t%s\n", *v.Name, *v.CreationDate)
	}
}
