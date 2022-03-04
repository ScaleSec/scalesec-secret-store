// TODO:

// [] Write some place for example
// [] Read some place for example
// [] verify crashing logger

// Credit where Credit is due.  This code is based on the
// Hashicorp Vault Guide:  https://github.com/hashicorp/vault-guides.git

// ********************************************************************************
// NOTE On Logging
// Vault loads the plugin in twice.  Once for initiation the second through the framework.
// Depending on the phase you need to logout information differently. This is why you see two ways
// of loggin in this code base.
// Init: you use the hclogger
// Framework: you use the backend logger
// If you mix them up then you will not see the messages and/or vault will see the
// messages and think the pluging did not load correctly
// ********************************************************************************

package scalesecSecretStore

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const scalesecSecretStoreBackendHelp = `
The ScaleSec Secrets Store Vault Plugin is a functional example of how a vault secret backend work.
Add your business logic to the plugin as needed to crate your own custom secret backend.
`

// Hashicorp: "github.com/hashicorp/vault/sdk/framework"
//
// Backend is an implementation of logical.Backend that allows
// the implementer to code a backend using a much more programmer-friendly
// framework that handles a lot of the routing and validation for you.
//
// This is recommended over implementing logical.Backend directly.
// backend wraps the backend framework and adds a map for storing key value pairs
type scalesecSecretStoreBackend struct {
	*framework.Backend

	// You can add additional vars here that you want to keep for the running
	// of the plugin .. Like configuration arguments.
	pluginName string
}

var _ logical.Factory = Factory

// Factory configures and returns Mock backends
// Goal: Take in the configurations and configure our backend.
// Return: our newly configured backend or and error object
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	b, err := newBackend()
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	//  https://pkg.go.dev/github.com/hashicorp/vault/sdk@v0.3.0/logical#BackendConfig
	// TODO: THESE - SHOW - VERIFY IF THEY ARE CRASHING PLUGIN
	//	conf.Logger.Debug("scalesecSecretStore.Factory:-> Enter")
	//	conf.Logger.Debug("scalesecSecretStore.Factory:-> conf.BackendUUID: %s", conf.BackendUUID)
	//	conf.Logger.Debug("scalesecSecretStore.Factory:-> conf.Config: %v", conf.Config)

	if err := b.Setup(ctx, conf); err != nil {
		conf.Logger.Debug("scalesecSecretStore.Factory:-> b.Setup error: %s", err)
		return nil, err
	}

	//	conf.Logger.Debug("scalesecSecretStore.Factory:-> Leaving")
	return b, nil
}

// Constructor to inital our Backend structure so Vault plugin framework cancall out functions
// on our Backend.

func newBackend() (*scalesecSecretStoreBackend, error) {

	// TODO: THESE - SHOW - VERIFY IF THEY ARE CRASHING PLUGIN
	//	hclogger.Default().Info("scalesecSecretStore:newBackend(): -> Enter")

	b := &scalesecSecretStoreBackend{
		// if you have additional vars to the backend structure you would init them here
		pluginName: "scalesecSecretStore",
	}

	b.Backend = &framework.Backend{
		// Set the plugin help string
		Help: strings.TrimSpace(scalesecSecretStoreBackendHelp),
		// Tell vault the type of this plugin.  We have 2 choices:
		// 1 TypeLogical    = Secret Store Backend
		// 2 TypeCredential = Authorization Backend
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			b.paths(),
		),
	}

	// TODO: THESE - SHOW - VERIFY IF THEY ARE CRASHING PLUGIN
	//	hclogger.Default().Info("scalesecSecretStore:newBackend(): -> Leaving")
	return b, nil
}

// setup the mapping between our functions and the hashicorp framework so it knows how to call
// the functions that we have implemented based on the command given to vault.

// It gets the path from the command and operation that was requested.
// It provides pointer to our functions for vault to call to perfrom internal operations
func (b *scalesecSecretStoreBackend) paths() []*framework.Path {
	// TODO: THESE - SHOW - VERIFY IF THEY ARE CRASHING PLUGIN
	//	hclogger.Default().Info("scalesecSecretStore.paths(): -> Enter")
	//	hclogger.Default().Info("scalesecSecretStore.paths(): -> Leaving")

	// TODO: change this to return to a var then return the var so we can move the leaving message
	return []*framework.Path{
		{
			//
			// getting the path
			//

			Pattern: framework.MatchAllRegex("path"),

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "Specifies the path of the secret.",
				},
			},

			//
			// mapping the operational request:  Read; Write; Create; Delete
			//

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRead,
					Summary:  "Retrieve the secret from the map.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Creates the secret at the specified location.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleList,
					Summary:  "Lists the secret at the specified location.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDelete,
					Summary:  "Deletes the secret at the specified location.",
				},
			},

			//
			// Provide pointers to our implemented functions for interall operations
			// https://github.com/hashicorp/vault/blob/main/sdk/framework/backend.go
			//

			ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

//
//
//
//

// ============================================================================================
// handleExistenceCheck: Check your secret Store to see if an secret exist
//
// GOAL: 	Replace this logic with your logic to determine if a secret exists or not
// RETURN:  bool  := Return True/False  True = Exist  False = Does Not Exist
//			error := Error message if there is an error in your processing - nil if you have no error
// ============================================================================================
func (b *scalesecSecretStoreBackend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	b.Logger().Debug("scalesecSecretStore.handleExistenceCheck:-> Enter")
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleExistenceCheck:-> *logical.Request: %v", req))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleExistenceCheck:-> *framework.FieldData: %v", data))

	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleExistenceCheck:-> req.Data: %s", req.Data))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleExistenceCheck:-> req.Path: %s", req.Path))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleExistenceCheck:-> req.MountPoint: %s", req.MountPoint))

	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****
	// ***** Start Replace with your logic to determin if the secret exists:

	// Read from the local storage to see if the secret exists
	out, err := req.Storage.Get(ctx, req.Path)

	if err != nil {
		b.Logger().Debug("scalesecSecretStore.handleExistenceCheck:-> Leaving with error")
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	b.Logger().Debug("scalesecSecretStore.handleExistenceCheck:-> Leaving")

	// ***** End Replace of your logic: Return Boolean (True if Exist or False if it does not); Error or nil
	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****

	return out != nil, nil
}

// ============================================================================================
// handleRead: Read from your secret store
//
// GOAL:  	Read from the secret store useing the request information and path
// 		  	Build a response with the secret data that was stored and return or return error
// Return:
// 			*logical.Response := Response with the Key Value pairs of secret data that was stored or nil if there is an error
// 			error := Error with details if the data was not able to be read or nil if success.
// ============================================================================================

func (b *scalesecSecretStoreBackend) handleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Debug("scalesecSecretStore.handleRead:-> Enter")
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleRead:-> *logical.Request: %v", req))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleRead:-> *framework.FieldData: %v", data))

	// TODO: Do we need this
	//	if req.ClientToken == "" {
	//		b.Logger().Debug("scalesecSecretStore.handleRead:-> Leaving with error")
	//		return nil, fmt.Errorf("client token empty")
	//	}

	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****
	// **** Start your read logic

	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleRead:-> req.Data: %s", req.Data))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleRead:-> req.Path: %s", req.Path))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleRead:-> req.MountPoint: %s", req.MountPoint))

	path := data.Get("path").(string)
	// read from local storage
	//	out, readerr := req.Storage.Get(ctx, req.Path)
	//	if readerr != nil {
	//		return nil, fmt.Errorf("error getting reading secret path: %s error: %s", req.Path, readerr)
	//	}

	// Read data from the storage backend based on the path provided
	// the data we are reading should come back as a json string
	var rawData map[string]interface{}
	// Example hardcoded data
	fetchedData := []byte(`{"secretkey":"secretValue", "secretPath":"` + path + `"}`)
	//fetchedData := []byte(fmt.Sprintf("{\"%s\":\"%v\", \"secretPath\":\"%s\"}", out.Key, out.Value, path))

	// Check to see if we have data that should be returned.
	if fetchedData == nil {
		b.Logger().Debug("scalesecSecretStore.handleRead:-> Leaving error message in response")
		resp := logical.ErrorResponse("No value at Mount:%v Path:%v", req.MountPoint, path)
		return resp, nil
	}

	// Take the data and load  the rawData interface so it can go into the response
	err := jsonutil.DecodeJSON(fetchedData, &rawData)
	if err != nil {
		// use the HCP errwrap class to create and return an error message
		b.Logger().Debug("scalesecSecretStore.handleRead:-> Leaving with error")
		return nil, fmt.Errorf("json decoding failed: %w", err)
	}

	// ***** End - Your Read logic
	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****

	// Generate the json response
	resp := &logical.Response{
		Data: rawData,
	}

	b.Logger().Debug("scalesecSecretStore.handleRead:-> Leaving Resp with data")
	return resp, nil
}

// ============================================================================================
// handleWrite: Write to your secret store. If the secret exists then you should over write it for an update
//
// GOAL:  	Write or update the secret to your secret store
// Return:
// 			*logical.Response := Response that is null as write does not have response
// 			error := Error with details if the write failed or nil if success.
// ============================================================================================

func (b *scalesecSecretStoreBackend) handleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Debug("scalesecSecretStore.handleWrite:-> Enter")
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleWrite:-> req.ClientToken: %s", req.ClientToken))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleWrite:-> *logical.Reqeust: %v", req))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleWrite:-> *framework.FieldData: %v", data))

	// We don't need this
	// check to see if we have a client token on the request
	//	if req.ClientToken == "" {
	//		b.Logger().Debug("scalesecSecretStore.handleWrite:-> Leaving with error")
	//		return nil, fmt.Errorf("client token empty")
	//	}

	// Check to make sure that we have data to actually store
	if len(req.Data) == 0 {
		b.Logger().Debug("scalesecSecretStore.handleWrite:-> Leaving with error")
		return nil, fmt.Errorf("data must be provided to store in secret")
	}

	// TODO: Is this the path we want or the req.Path
	path := data.Get("path").(string)
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleWrite:-> Path: %s", path))

	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****
	// ***** Start - YOUR STORAGE WRITE LOGIC

	// Store the secert data in the storage backend for the specified path

	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleWrite:-> req.Data: %s", req.Data))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleWrite:-> req.Path: %s", req.Path))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleWrite:-> req.MountPoint: %s", req.MountPoint))

	// ** THIS WRITE LOGIC DOES NOT WORK >> NEED SOME THING TO WRITE TO .. MAYBE OS ???

	// Read from the local storage to get it's location

	//*	out, readerr := req.Storage.Get(ctx, req.Path)
	//*	if readerr != nil {
	//*		return nil, fmt.Errorf("error getting storage location for path: %s error: %s", req.Path, readerr)
	//*	}
	//*	if out == nil {
	//*		b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleWrite:-> logical storeage is nil"))
	//*		return nil, fmt.Errorf("scalesecSecretStore.handleWrite: logical storeage is nil")
	//*	}
	//*
	//*	// write to the local storage
	//*	for key, value := range req.Data {
	//*		out.Key = key
	//*		out.Value = value.([]byte)
	//*		fmt.Printf("key: %s value: %v\n", key, value)
	//*		err := req.Storage.Put(ctx, out)
	//*		if err != nil {
	//*			return nil, fmt.Errorf("error writing %s error %s", key, err)
	//*		}
	//*	}

	// ***** End - YOUR STORAGE WRITE LOGIC
	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****

	// return nil logical.Response and nil error for success
	b.Logger().Debug("scalesecSecretStore.handleWrite:-> Leaving")
	return nil, nil
}

// ============================================================================================
// handleDelete: Delete from your secret store.
//
// GOAL:  	Delete from your secret store
// Return:
// 			Response that is nil
// 			Error with details if the delete failed or nil if success.
// ============================================================================================

func (b *scalesecSecretStoreBackend) handleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Debug("scalesecSecretStore.handleDelete:-> Enter")
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleDelete:-> *logical.Request: %v", req))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleDelete:-> *framework.FieldData: %v", data))

	if req.ClientToken == "" {
		b.Logger().Debug("scalesecSecretStore.handleDelete:-> Leaving with error")
		return nil, fmt.Errorf("client token empty")
	}

	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****
	// ***** Start your delete Logic

	// EXAMPLE DELETE LOGIC

	// Remove entry form the storage backend for the specified path

	// ***** End your Delete Logic
	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****

	// return nil logical.Response and nil error for success

	b.Logger().Debug("scalesecSecretStore.handleDelete:-> Leaving")
	return nil, nil
}

// ============================================================================================
// handleList: keys in a secret store.
//
// GOAL:  	Provide a string array list of secret keys stored
// Return:
// 			*logical.Response := Response that is not nil with a list of keys
// 			error : = Error with details if the write failed or nil if success.
// ============================================================================================

func (b *scalesecSecretStoreBackend) handleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Debug("scalesecSecretStore.handleList:-> Enter")
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleList:-> *logical.Request: %v", req))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleList:-> *framework.FieldData: %v", data))

	// TODO: Do we need this
	//	if req.ClientToken == "" {
	//		b.Logger().Debug("scalesecSecretStore.handleRead:-> Leaving with error")
	//		return nil, fmt.Errorf("client token empty")
	//	}

	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****
	// **** Start your List logic

	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleList:-> req.Data: %s", req.Data))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleList:-> req.Path: %s", req.Path))
	b.Logger().Debug(fmt.Sprintf("scalesecSecretStore.handleList:-> req.MountPoint: %s", req.MountPoint))

	path := data.Get("path").(string)

	// Example hardcoded data
	fetchedData := []string{}
	fetchedData = append(fetchedData, "key1")
	fetchedData = append(fetchedData, "key2")

	// Check to see if we have data that should be returned.
	if fetchedData == nil {
		b.Logger().Debug("scalesecSecretStore.handleList:-> Leaving error message in response")
		resp := logical.ErrorResponse("No value at Mount:%v Path:%v", req.MountPoint, path)
		return resp, nil
	}

	// Take the data and load into the response
	resp := logical.ListResponse(fetchedData)

	// ***** End - Your List logic
	// ***** ***** ***** ***** ***** ***** ***** ***** ***** ***** *****

	b.Logger().Debug("scalesecSecretStore.handleRead:-> Leaving Resp with data")
	return resp, nil

}
