package linodego_test

import (
	"testing"
)

func TestListInstances(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, teardown := createTestClient(t, "fixtures/TestListInstances")
	defer teardown()

	linodes, err := client.ListInstances(nil)
	if err != nil {
		t.Errorf("Error listing instances, expected struct, got error %v", err)
	}
	if len(linodes) == 0 {
		t.Errorf("Expected a list of instances, but got %v", linodes)
	}
}

func TestGetInstance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, teardown := createTestClient(t, "fixtures/TestGetInstance")
	defer teardown()

	instance, err := client.GetInstance(TestInstanceID)
	if err != nil {
		t.Errorf("Error getting instance TestInstanceID, expected *LinodeInstance, got error %v", err)
	}
	if instance.Specs.Disk <= 0 {
		t.Errorf("Error in instance TestInstanceID spec for disk size, %v", instance.Specs)
	}
}

func TestListInstanceDisks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, teardown := createTestClient(t, "fixtures/TestListInstanceDisks")
	defer teardown()

	disks, err := client.ListInstanceDisks(TestInstanceID, nil)
	if err != nil {
		t.Errorf("Error listing instance disks, expected struct, got error %v", err)
	}
	if len(disks) == 0 {
		t.Errorf("Expected a list of instance disks, but got %v", disks)
	}
}

func TestListInstanceConfigs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, teardown := createTestClient(t, "fixtures/TestListInstanceConfigs")
	defer teardown()

	configs, err := client.ListInstanceConfigs(TestInstanceID, nil)
	if err != nil {
		t.Errorf("Error listing instance configs, expected struct, got error %v", err)
	}
	if len(configs) == 0 {
		t.Errorf("Expected a list of instance configs, but got %v", configs)
	}
}

func TestListInstanceVolumes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, teardown := createTestClient(t, "fixtures/TestListInstanceVolumes")
	defer teardown()

	volumes, err := client.ListInstanceVolumes(TestInstanceID, nil)
	if err != nil {
		t.Errorf("Error listing instance volumes, expected struct, got error %v", err)
	}
	if len(volumes) == 0 {
		t.Errorf("Expected an list of instance volumes, but got %v", volumes)
	}
}