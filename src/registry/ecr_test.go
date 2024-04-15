package registry

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryInfoFromURLSucceeds(t *testing.T) {
	cases := []struct {
		test     string
		url      string
		expected autogold.Value
	}{
		{
			test: "Url with label",
			url:  "123456789012.dkr.ecr.us-west-2.amazonaws.com/test-repo:latest",
			expected: autogold.Expect(ImageReference{
				RegistryID: "123456789012", Region: "us-west-2",
				Name: "test-repo",
				Tag:  "latest",
			}),
		},
		{
			test: "Url with digest",
			url:  "123456789012.dkr.ecr.us-west-2.amazonaws.com/test-repo@sha256:hash",
			expected: autogold.Expect(ImageReference{
				RegistryID: "123456789012", Region: "us-west-2",
				Name:   "test-repo",
				Digest: "sha256:hash",
			}),
		},
		{
			test: "Url with tag and digest",
			url:  "123456789012.dkr.ecr.us-west-2.amazonaws.com/test-repo:tagged@sha256:hash",
			expected: autogold.Expect(ImageReference{
				RegistryID: "123456789012", Region: "us-west-2",
				Name: "test-repo",
				Tag:  "tagged@sha256:hash",
			}),
		},
		{
			test: "Url without label",
			url:  "123456789012.dkr.ecr.us-west-2.amazonaws.com/test-repo",
			expected: autogold.Expect(ImageReference{
				RegistryID: "123456789012", Region: "us-west-2",
				Name: "test-repo",
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			info, err := ParseReferenceFromURL(c.url)
			require.NoError(t, err)
			c.expected.Equal(t, info)
		})
	}
}

func TestRegistryInfoFromURLFails(t *testing.T) {
	url := "123456789012.dkr.ecr.us-west-2.amazonaws.com"

	info, err := ParseReferenceFromURL(url)
	require.ErrorContains(t, err, "invalid registry URL")

	assert.Equal(t, ImageReference{}, info)
}

type mockedDescribeImageScanFindings struct {
	ECRAPI
}

func (m mockedDescribeImageScanFindings) DescribeImageScanFindings(ctx context.Context, params *ecr.DescribeImageScanFindingsInput, opts ...func(*ecr.Options)) (*ecr.DescribeImageScanFindingsOutput, error) {
	return nil, &types.ReferencedImagesNotFoundException{}
}

func TestWaitForScanFindings(t *testing.T) {
	r := &RegistryScan{
		Client:          &mockedDescribeImageScanFindings{},
		MinAttemptDelay: 1 * time.Millisecond,
		MaxAttemptDelay: 2 * time.Millisecond,
		MaxTotalDelay:   100 * time.Millisecond,
	}
	err := r.WaitForScanFindings(context.TODO(), ImageReference{})
	assert.Error(t, err)
}
