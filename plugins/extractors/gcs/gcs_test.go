//go:build plugins
// +build plugins

package gcs

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/test/mocks"
	"github.com/goto/meteor/test/utils"
	"github.com/goto/salt/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
)

func TestInit(t *testing.T) {
	t.Run("should return error if no project_id in config", func(t *testing.T) {
		err := New(utils.Logger, nil).Init(context.TODO(), plugins.Config{
			URNScope: "test",
			RawConfig: map[string]interface{}{
				"wrong-config": "sample-project",
			},
		})

		assert.ErrorAs(t, err, &plugins.InvalidConfigError{})
	})

	t.Run("should return error if service_account_base64 config is invalid", func(t *testing.T) {
		extr := New(utils.Logger, createClient)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-gcs",
			RawConfig: map[string]interface{}{
				"project_id":             "google-project-id",
				"service_account_base64": "----", // invalid
			},
		})

		assert.ErrorContains(t, err, "decode Base64 encoded service account")
	})

	t.Run("should return no error", func(t *testing.T) {
		extr := New(utils.Logger, func(context.Context, log.Logger, Config) (stiface.Client, error) {
			return mockClient{}, nil
		})
		ctx := context.Background()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-gcs",
			RawConfig: map[string]interface{}{
				"project_id": "google-project-id",
			},
		})

		assert.NoError(t, err)
	})
}

func TestExtract(t *testing.T) {
	t.Run("should return no error", func(t *testing.T) {
		extr := New(utils.Logger, func(context.Context, log.Logger, Config) (stiface.Client, error) {
			return mockClient{
				mockBuckets: &mockBucketIterator{
					next: []storage.BucketAttrs{
						{Name: "bucket-1"},
					},
				},
				mockBucket: mockBucketHandle{
					it: &mockObjectIterator{
						next: []storage.ObjectAttrs{
							{
								Name:   "object-1",
								Bucket: "bucket-1",
								Owner:  "owner-1",
							},
						},
					},
				},
			}, nil
		})
		ctx := context.Background()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-gcs",
			RawConfig: map[string]interface{}{
				"project_id":   "google-project-id",
				"extract_blob": "true",
			},
		})

		assert.NoError(t, err)

		err = extr.Extract(context.TODO(), mocks.NewEmitter().Push)
		assert.NoError(t, err)
	})

	t.Run("should return error when bucket iterator returns error", func(t *testing.T) {
		extr := New(utils.Logger, func(context.Context, log.Logger, Config) (stiface.Client, error) {
			return mockClient{
				mockBuckets: &mockBucketIterator{
					err: errors.New("some error"),
					next: []storage.BucketAttrs{
						{Name: "bucket-1"},
					},
				},
			}, nil
		})
		ctx := context.Background()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-gcs",
			RawConfig: map[string]interface{}{
				"project_id": "google-project-id",
			},
		})

		assert.NoError(t, err)

		err = extr.Extract(context.TODO(), mocks.NewEmitter().Push)
		assert.ErrorContains(t, err, "iterate over")
	})

	t.Run("should return error when extract blob", func(t *testing.T) {
		extr := New(utils.Logger, func(context.Context, log.Logger, Config) (stiface.Client, error) {
			return mockClient{
				mockBuckets: &mockBucketIterator{
					next: []storage.BucketAttrs{
						{Name: "bucket-1"},
					},
				},
				mockBucket: mockBucketHandle{
					it: &mockObjectIterator{
						err: errors.New("some error"),
						next: []storage.ObjectAttrs{
							{
								Name:   "object-1",
								Bucket: "bucket-1",
								Owner:  "owner-1",
							},
						},
					},
				},
			}, nil
		})
		ctx := context.Background()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-gcs",
			RawConfig: map[string]interface{}{
				"project_id":   "google-project-id",
				"extract_blob": "true",
			},
		})

		assert.NoError(t, err)

		err = extr.Extract(context.TODO(), mocks.NewEmitter().Push)
		assert.ErrorContains(t, err, "extract blobs from")
	})
}

type mockClient struct {
	stiface.Client
	mockBuckets stiface.BucketIterator
	mockBucket  stiface.BucketHandle
}

type mockBucketIterator struct {
	stiface.BucketIterator
	i    int
	next []storage.BucketAttrs
	err  error
}

func (it *mockBucketIterator) Next() (*storage.BucketAttrs, error) {
	if it.i >= len(it.next) {
		return nil, iterator.Done
	}

	nextAttr := &it.next[it.i]
	it.i++
	return nextAttr, it.err
}

func (c mockClient) Buckets(ctx context.Context, projectID string) stiface.BucketIterator {
	return c.mockBuckets
}

func (c mockClient) Bucket(name string) stiface.BucketHandle {
	return c.mockBucket
}

type mockBucketHandle struct {
	stiface.BucketHandle
	it stiface.ObjectIterator
}

func (h mockBucketHandle) Objects(ctx context.Context, query *storage.Query) stiface.ObjectIterator {
	return h.it
}

type mockObjectIterator struct {
	stiface.ObjectIterator
	i    int
	next []storage.ObjectAttrs
	err  error
}

func (it *mockObjectIterator) Next() (*storage.ObjectAttrs, error) {
	if it.i >= len(it.next) {
		return nil, iterator.Done
	}

	nextAttr := &it.next[it.i]
	it.i++
	return nextAttr, it.err
}
