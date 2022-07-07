package bucket

import (
	"lesson1/client"
	"testing"
)

func TestBucket_RunList(t *testing.T) {
	b := NewBucket()
	b.RunList(client.Config())
}
