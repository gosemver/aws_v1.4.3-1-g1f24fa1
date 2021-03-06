package s3crypto_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/aws"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/aws/request"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/awstesting/unit"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/service/kms"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/service/s3"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/service/s3/s3crypto"
)

func TestDefaultConfigValues(t *testing.T) {
	sess := unit.Session.Copy(&aws.Config{
		MaxRetries:       aws.Int(0),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-west-2"),
	})
	svc := kms.New(sess)
	handler := s3crypto.NewKMSKeyGenerator(svc, "testid")

	c := s3crypto.NewEncryptionClient(sess, s3crypto.AESGCMContentCipherBuilder(handler))

	assert.NotNil(t, c)
	assert.NotNil(t, c.ContentCipherBuilder)
	assert.NotNil(t, c.SaveStrategy)
}

func TestPutObject(t *testing.T) {
	size := 1024 * 1024
	data := make([]byte, size)
	expected := bytes.Repeat([]byte{1}, size)
	generator := mockGenerator{}
	cb := mockCipherBuilder{generator}
	sess := unit.Session.Copy(&aws.Config{
		MaxRetries:       aws.Int(0),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-west-2"),
	})
	c := s3crypto.NewEncryptionClient(sess, cb)
	assert.NotNil(t, c)
	input := &s3.PutObjectInput{
		Key:    aws.String("test"),
		Bucket: aws.String("test"),
		Body:   bytes.NewReader(data),
	}
	req, _ := c.PutObjectRequest(input)
	req.Handlers.Send.Clear()
	req.Handlers.Send.PushBack(func(r *request.Request) {
		r.Error = errors.New("stop")
		r.HTTPResponse = &http.Response{
			StatusCode: 200,
		}
	})
	err := req.Send()
	assert.Equal(t, "stop", err.Error())
	b, err := ioutil.ReadAll(req.HTTPRequest.Body)
	assert.NoError(t, err)
	assert.Equal(t, expected, b)
}
