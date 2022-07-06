package main

import (
	"flag"

	"lesson1/bucket"
	"lesson1/client"
	"log"
)

var deleteBucket *bool
var listBuckets *bool

func Options() *string {
	bucketPtr := flag.String("bucket", "", "unique bucket name")
	deleteBucket = flag.Bool("delete", false, "delete flag")
	listBuckets = flag.Bool("list", false, "list buckets")
	flag.Parse()
	if *bucketPtr == "" && !*listBuckets {
		log.Fatalf(`Need to enter a bucket name.
	
		 -bucket=bucketsample or -list
	`)
	}
	return bucketPtr
}

func main() {
	bucketName := Options()
	b := bucket.NewBucket(bucketName)
	if *deleteBucket {
		b.RunDelete(client.Config())
		return
	}
	if *listBuckets {
		b.RunList(client.Config())
		return
	}
	b.RunCreate(client.Config())
}
