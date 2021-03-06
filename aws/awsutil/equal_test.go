package awsutil_test

import (
	"testing"

	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/aws"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/aws/awsutil"
	"github.com/stretchr/testify/assert"
)

func TestDeepEqual(t *testing.T) {
	cases := []struct {
		a, b  interface{}
		equal bool
	}{
		{"a", "a", true},
		{"a", "b", false},
		{"a", aws.String(""), false},
		{"a", nil, false},
		{"a", aws.String("a"), true},
		{(*bool)(nil), (*bool)(nil), true},
		{(*bool)(nil), (*string)(nil), false},
		{nil, nil, true},
	}

	for i, c := range cases {
		assert.Equal(t, c.equal, awsutil.DeepEqual(c.a, c.b), "%d, a:%v b:%v, %t", i, c.a, c.b, c.equal)
	}
}
