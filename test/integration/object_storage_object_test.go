package integration

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/linode/linodego"
)

var objectStorageObjectURLExpirySeconds = 360

func putObjectStorageObject(t *testing.T, client *linodego.Client, bucket *linodego.ObjectStorageBucket, name, content string) {
	t.Helper()

	url, err := client.CreateObjectStorageObjectURL(context.TODO(), bucket.Cluster, bucket.Label, linodego.ObjectStorageObjectURLCreateOptions{
		Name:        name,
		Method:      http.MethodPut,
		ContentType: "text/plain",
		ExpiresIn:   &objectStorageObjectURLExpirySeconds,
	})
	if err != nil {
		t.Errorf("failed to get object PUT url: %s", err)
	}

	rec, teardownRecorder := testRecorder(t, "fixtures/TestObjectStorageObject_ACLConfig_Bucket_Put", testingMode, nil)
	defer teardownRecorder()

	httpClient := http.Client{Transport: rec}
	req, err := http.NewRequest(http.MethodPost, url.URL, bytes.NewReader([]byte(content)))
	if err != nil {
		t.Errorf("failed to build request: %s", err)
	}
	req.Method = http.MethodPut
	req.Header.Add("Content-Type", "text/plain")

	res, err := httpClient.Do(req)
	if err != nil {
		t.Errorf("failed to make request: %s", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("expected status code to be 200; got %d", res.StatusCode)
	}
}

func deleteObjectStorageObject(t *testing.T, client *linodego.Client, bucket *linodego.ObjectStorageBucket, name string) {
	t.Helper()

	url, err := client.CreateObjectStorageObjectURL(context.TODO(), bucket.Cluster, bucket.Label, linodego.ObjectStorageObjectURLCreateOptions{
		Name:      name,
		Method:    http.MethodDelete,
		ExpiresIn: &objectStorageObjectURLExpirySeconds,
	})
	if err != nil {
		t.Errorf("failed to get object PUT url: %s", err)
	}

	rec, teardownRecorder := testRecorder(t, "fixtures/TestObjectStorageObject_ACLConfig_Bucket_Delete", testingMode, nil)
	defer teardownRecorder()

	httpClient := http.Client{Transport: rec}
	req, err := http.NewRequest(http.MethodPost, url.URL, nil)
	if err != nil {
		t.Errorf("failed to build request: %s", err)
	}
	req.Method = http.MethodDelete

	res, err := httpClient.Do(req)
	if res.StatusCode != 204 {
		t.Errorf("expected status code to be 204; got %d", res.StatusCode)
	}
}

func TestObjectStorageObject_Smoke(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(
		t, nil, "fixtures/TestObjectStorageObject_Smoke",
		nil, nil, nil,
	)
	if err != nil {
		t.Fatalf("failed to create Object Storage Object: %s", err)
	}
	defer teardown()

	object := "test"
	putObjectStorageObject(t, client, bucket, object, "testing123")
	defer deleteObjectStorageObject(t, client, bucket, object)

	config, err := client.GetObjectStorageObjectACLConfig(context.TODO(), bucket.Cluster, bucket.Label, object)
	if err != nil {
		t.Errorf("failed to get ACL config: %s", err)
	}

	if config.ACL != "private" {
		t.Errorf("expected ACL to be private; got %s", config.ACL)
	}

	if config.ACLXML == "" {
		t.Error("expected ACL XML to be included")
	}

	configv2, err := client.GetObjectStorageObjectACLConfigV2(context.TODO(), bucket.Cluster, bucket.Label, object)
	if err != nil {
		t.Errorf("failed to get ACL config: %s", err)
	}

	if configv2.ACL == nil {
		t.Errorf("expected ACL to be private; got nil")
	}

	if configv2.ACL != nil && *configv2.ACL != "private" {
		t.Errorf("expected ACL to be private; got %s", *configv2.ACL)
	}

	content, err := client.ListObjectStorageBucketContents(context.TODO(), bucket.Cluster, bucket.Label, nil)
	if err != nil {
		t.Errorf("failed to get bucket contents: %s", err)
	}

	if content.Data[0].Name != object {
		t.Errorf("ObjectStorageBucket contents name does not match, expected %s, got %s", "test", content.Data[0].Name)
	}

	updateOpts := linodego.ObjectStorageObjectACLConfigUpdateOptions{ACL: "public-read", Name: object}
	if _, err = client.UpdateObjectStorageObjectACLConfig(context.TODO(), bucket.Cluster, bucket.Label, updateOpts); err != nil {
		t.Errorf("failed to update ACL config: %s", err)
	}

	if _, err = client.UpdateObjectStorageObjectACLConfigV2(context.TODO(), bucket.Cluster, bucket.Label, updateOpts); err != nil {
		t.Errorf("failed to update ACL config: %s", err)
	}

	config, err = client.GetObjectStorageObjectACLConfig(context.TODO(), bucket.Cluster, bucket.Label, object)
	if err != nil {
		t.Errorf("failed to get updated ACL config: %s", err)
	}

	if config.ACL != updateOpts.ACL {
		t.Errorf("expected ACL config to be %s; got %s", updateOpts.ACL, config.ACL)
	}
	if config.ACLXML == "" {
		t.Error("expected ACL XML to be included")
	}

	configv2, err = client.GetObjectStorageObjectACLConfigV2(context.TODO(), bucket.Cluster, bucket.Label, object)
	if err != nil {
		t.Errorf("failed to get ACL config: %s", err)
	}

	if configv2.ACL == nil {
		t.Errorf("expected ACL config to be %s; got nil", updateOpts.ACL)
	}

	if configv2.ACL != nil && *configv2.ACL != updateOpts.ACL {
		t.Errorf("expected ACL config to be %s; got nil", updateOpts.ACL)
	}

	if configv2.ACLXML == nil {
		t.Error("expected ACL XML to be included")
	}

	if configv2.ACLXML != nil && *configv2.ACLXML == "" {
		t.Error("expected ACL XML to be included")
	}
}
