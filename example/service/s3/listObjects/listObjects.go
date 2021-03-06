package main

import (
	"fmt"
	"os"

	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/aws/session"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/service/s3"
)

// Lists all objects in a bucket using pagination
//
// Usage:
// go run listObjects.go <bucket>
func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := s3.New(sess)

	i := 0
	err = svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &os.Args[1],
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		fmt.Println("Page,", i)
		i++

		for _, obj := range p.Contents {
			fmt.Println("Object:", *obj.Key)
		}
		return true
	})
	if err != nil {
		fmt.Println("failed to list objects", err)
		return
	}
}
