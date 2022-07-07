package bucket

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

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

func ItemsInBucket(cfg aws.Config, bucket string) ([]string, error) {
	client := s3.NewFromConfig(cfg)
	out := []string{}
	input := &s3.ListObjectsV2Input{
		Bucket: &bucket,
	}

	resp, err := getObjects(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error retrieving list of objects:")
		fmt.Println(err)
		return out, err
	}
	fmt.Println("\n\nObjects in " + bucket + ":")

	for _, item := range resp.Contents {
		fmt.Println("Name:          ", *item.Key)
		fmt.Println("Last modified: ", *item.LastModified)
		fmt.Println("Size:          ", item.Size)
		fmt.Println("Storage class: ", item.StorageClass)
		fmt.Println("")
		out = append(out, *item.Key)
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
	fmt.Println("")
	return out, nil

}

type s3ListObjectsAPI interface {
	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func getObjects(c context.Context, api s3ListObjectsAPI, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(c, input)
}

func (b *Bucket) download(cfg aws.Config, bucket, key string, w io.WriterAt) (int64, error) {
	client := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)

	numBytes, err := downloader.Download(context.TODO(), w, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return numBytes, err
}

func (b *Bucket) upload(cfg aws.Config, bucket, key string, r io.Reader) (*manager.UploadOutput, error) {
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    aws.String(key),
		Body:   r,
	})

	return result, err
}

func (b *Bucket) Upload(cfg aws.Config, bucket, key, fileToUpload string) (*manager.UploadOutput, error) {
	f, err := os.Open(fileToUpload)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return b.upload(cfg, bucket, key, f)

}

func (b *Bucket) Delete(cfg aws.Config) *s3.DeleteBucketOutput {
	input := &s3.DeleteBucketInput{
		Bucket: b.Name,
	}
	result, err := b.deleteBucket(cfg, input)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	return result
}

func (b *Bucket) Create(cfg aws.Config) {
	input := &s3.CreateBucketInput{
		Bucket: b.Name,
	}
	result, err := b.createBucket(cfg, input)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	fmt.Printf("result: %v", result)
}

func (b *Bucket) List(cfg aws.Config) (*s3.ListBucketsOutput, error) {
	input := &s3.ListBucketsInput{}
	result, err := b.listBuckets(cfg, input)
	if err != nil {
		return result, err
	}
	for _, v := range result.Buckets {
		fmt.Printf("%s, \t\t%s\n", *v.Name, *v.CreationDate)
	}
	return result, err
}
