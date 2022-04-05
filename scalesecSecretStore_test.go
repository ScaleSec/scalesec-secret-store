package scalesecSecretStore

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

// Path and Mount point values are set by the vault read and write commands:
// and should be consistant for all of these tests.
// vault write scalesecsecrets/test secret_key="secret_value"
// vault list scalesecsecrets/test
const MOUNT_POINT = "scalesecsecrets/"
const BACKEND_PATH = "test/"

func getBackend(t *testing.T) (logical.Backend, logical.Storage) {

	configMap := map[string]string{}
	configMap["plugin_name"] = "scalesecSecretStorePlugin"
	configMap["plugin_type"] = "secret"
	// key value set by vault secrets enable -options=config_key=config_value
	configMap["config_key"] = "config_value"

	backendConfig := &logical.BackendConfig{
		Logger:      logging.NewVaultLogger(log.Trace),
		System:      &logical.StaticSystemView{},
		StorageView: &logical.InmemStorage{},
		BackendUUID: "test",
		Config:      configMap,
	}

	backend, err := Factory(context.Background(), backendConfig)

	if err != nil {
		t.Fatalf("unable to create backend: %v", err)
	}

	return backend, backendConfig.StorageView
}

// Test the list command:
// vault list scalesecsecrets/test
func TestList(t *testing.T) {

	b, storage := getBackend(t)

	request := &logical.Request{
		Operation:   logical.ListOperation,
		Path:        BACKEND_PATH,
		MountPoint:  MOUNT_POINT,
		Storage:     storage,
		ClientToken: "test_token",
	}

	response, err := b.HandleRequest(context.Background(), request)

	assert.Containsf(t, response.Data["keys"], "key1", "Vault List response should contain 'key1' - %v", response.Data)
	assert.Containsf(t, response.Data["keys"], "key2", "Vault List response should contain 'key2' - %v", response.Data)
	assert.Nil(t, err, "Response error %s", err)
	assert.NotNil(t, response, "Response should not be null")
	b.Logger().Debug("Response Object: %v", response)

}

// Test the write command:
// vault write scalesecsecrets/test secret_key="secret_value"
func TestWrite(t *testing.T) {

	b, storage := getBackend(t)

	data := map[string]interface{}{}

	data["secret_key"] = "secret_value"

	request := &logical.Request{
		Operation:   logical.CreateOperation,
		Path:        BACKEND_PATH,
		MountPoint:  MOUNT_POINT,
		Storage:     storage,
		ClientToken: "test_token",
		Data:        data,
	}

	response, err := b.HandleRequest(context.Background(), request)
	assert.Nil(t, err, "Response error %s", err)
	assert.Nil(t, response, "Response message %v", response)
	b.Logger().Debug("Response Object: %v", response)

}

// *********************************************************
// Test both read commands:
// vault read scalesecsecrets/test
// vault read scalesecsecrets/test secret_key=key_name
// *********************************************************

// vault read scalesecsecrets/test
func TestRead(t *testing.T) {

	b, storage := getBackend(t)

	request := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        BACKEND_PATH,
		MountPoint:  MOUNT_POINT,
		Storage:     storage,
		ClientToken: "test_token",
	}

	response, err := b.HandleRequest(context.Background(), request)

	// in our read with no data we return secretPath=test/
	assert.Containsf(t, response.Data["secretPath"], BACKEND_PATH, "Vault read response should contain '%s' - %v", BACKEND_PATH, response.Data)

	assert.Nil(t, err, "Response error %s", err)
	assert.NotNil(t, response, "Response should not be null")
	b.Logger().Debug("Response Object: %v", response)

}

// vault read scalesecsecrets/test secret_key=key_name
func TestReadWithData(t *testing.T) {

	b, storage := getBackend(t)

	data := map[string]interface{}{}

	data["secret_key"] = "key_name"

	request := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        BACKEND_PATH,
		MountPoint:  MOUNT_POINT,
		Storage:     storage,
		ClientToken: "test_token",
		Data:        data,
	}

	response, err := b.HandleRequest(context.Background(), request)

	// in our read with no data we return all_secrets_keys=all_secrets_values and secretPath=test/
	assert.Containsf(t, response.Data["secretPath"], BACKEND_PATH, "Vault read response should contain '%s' - %v", BACKEND_PATH, response.Data)
	assert.Containsf(t, response.Data["all_secrets_keys"], "all_secrets_values", "Vault read response should contain all_secrets_values - %v", response.Data)

	assert.Nil(t, err, "Response error %s", err)
	assert.NotNil(t, response, "Response should not be null")
	b.Logger().Debug("Response Object: %v", response)

}

// *********************************************************
// Test both delete commands:
// vault delete scalesecsecrets/test
// vault delete scalesecsecrets/test secret_key=key_name
// *********************************************************

// vault delete scalesecsecrets/test
func TestDelete(t *testing.T) {

	b, storage := getBackend(t)

	request := &logical.Request{
		Operation:   logical.DeleteOperation,
		Path:        BACKEND_PATH,
		MountPoint:  MOUNT_POINT,
		Storage:     storage,
		ClientToken: "test_token",
	}

	response, err := b.HandleRequest(context.Background(), request)

	assert.Nil(t, err, "Response error %s", err)
	b.Logger().Debug("Response Object: %v", response)
	fmt.Printf("TestDeleteWithData Response: %v", response)

}

// vault delete scalesecsecrets/test secret_key=key_name
func TestDeleteWithData(t *testing.T) {

	b, storage := getBackend(t)

	data := map[string]interface{}{}

	data["secret_key"] = "key_name"

	request := &logical.Request{
		Operation:   logical.DeleteOperation,
		Path:        BACKEND_PATH,
		MountPoint:  MOUNT_POINT,
		Storage:     storage,
		ClientToken: "test_token",
		Data:        data,
	}

	response, err := b.HandleRequest(context.Background(), request)

	assert.Nil(t, err, "Response error %s", err)
	b.Logger().Debug("Response Object: %v", response)
	fmt.Printf("TestDeleteWithData Response: %v", response)
}
