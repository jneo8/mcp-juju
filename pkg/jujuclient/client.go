package jujuclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/juju/api"
	applicationapi "github.com/juju/juju/api/client/application"
	apiclient "github.com/juju/juju/api/client/client"
	cloudapi "github.com/juju/juju/api/client/cloud"
	"github.com/juju/juju/api/client/modelmanager"
	jujucloud "github.com/juju/juju/cloud"
	"github.com/juju/juju/cmd/juju/application"
	"github.com/juju/juju/cmd/juju/common"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/juju"
	jclient "github.com/juju/juju/jujuclient"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/names/v6"
	"github.com/rs/zerolog/log"
)

// AddModelOptions contains configuration for creating a new model
type AddModelOptions struct {
	Name        string            // Required: model name
	Controller  string            // Optional: controller name (uses current if empty)
	Owner       string            // Optional: model owner (uses current user if empty)
	Credential  string            // Optional: credential name (auto-detect if empty)
	CloudRegion string            // Optional: cloud/region (uses controller default if empty)
	Config      map[string]string // Optional: model configuration
	NoSwitch    bool              // Optional: don't switch to new model after creation
}

type Client interface {
	// Basic information
	GetControllers() (Controllers, error)
	GetModels(controllerName string) (Models, error)
	GetStatus(ctx context.Context, controllerName, modelName string, includeStorage bool) (Status, error)

	// Application configuration
	GetApplicationConfig(ctx context.Context, controllerName, modelName, appName string) (ApplicationConfig, error)
	SetApplicationConfig(ctx context.Context, controllerName, modelName, appName string, settings map[string]string) error

	// Model management
	AddModel(ctx context.Context, opts AddModelOptions) error
}

type client struct {
	clientStore jclient.ClientStore
	logger      logger.Logger
}

func NewClient() (Client, error) {
	clientStore := jclient.NewFileClientStore()
	return &client{
		clientStore: clientStore,
		logger:      NewLoggerWrapper(),
	}, nil
}

func (c *client) GetControllers() (Controllers, error) {
	allControllers, err := c.clientStore.AllControllers()
	if err != nil {
		return Controllers{}, err
	}
	currentController, err := c.clientStore.CurrentController()
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			log.Debug().Msg("CurrentController not found")
		} else {
			return Controllers{}, err
		}
	}
	return Controllers{ControllerDetails: allControllers, Current: currentController}, nil
}

func (c *client) GetModels(currentController string) (Models, error) {
	if currentController == "" {
		controllers, err := c.GetControllers()
		if err != nil {
			return Models{}, err
		}
		currentController = controllers.Current
	}

	allModels, err := c.clientStore.AllModels(currentController)
	if err != nil {
		return Models{}, err
	}
	currentModel, err := c.clientStore.CurrentModel(currentController)
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			log.Debug().Msg("CurrentModel not found")
		} else {
			return Models{}, err
		}
	}
	return Models{ModelDetails: allModels, Current: currentModel}, nil
}

func (c *client) GetStatus(ctx context.Context, controllerName, modelName string, includeStorage bool) (Status, error) {
	if controllerName == "" {
		currentController, err := c.clientStore.CurrentController()
		if err != nil {
			return Status{}, err
		}
		controllerName = currentController
	}
	if modelName == "" {
		currentModel, err := c.clientStore.CurrentModel(controllerName)
		if err != nil {
			return Status{}, err
		}
		modelName = currentModel
	}

	apiClient, err := c.getAPIClient(ctx, controllerName, modelName)
	if err != nil {
		return Status{}, err
	}
	fullStatus, err := apiClient.Status(ctx, &apiclient.StatusArgs{Patterns: []string{}, IncludeStorage: includeStorage})
	if err != nil {
		return Status{}, err
	}
	return Status{FullStatus: *fullStatus}, nil
}

func (c *client) GetApplicationConfig(ctx context.Context, controllerName, modelName, appName string) (ApplicationConfig, error) {
	if controllerName == "" {
		currentController, err := c.clientStore.CurrentController()
		if err != nil {
			return ApplicationConfig{}, err
		}
		controllerName = currentController
	}
	if modelName == "" {
		currentModel, err := c.clientStore.CurrentModel(controllerName)
		if err != nil {
			return ApplicationConfig{}, err
		}
		modelName = currentModel
	}
	appAPI, err := c.getApplicationAPI(ctx, controllerName, modelName)
	if err != nil {
		return ApplicationConfig{}, err
	}
	results, err := appAPI.Get(ctx, appName)
	if err != nil {
		return ApplicationConfig{}, err
	}
	var appConfig ApplicationConfig
	appConfig.Application = results.Application
	appConfig.Charm = results.Charm
	appConfig.Settings = results.CharmConfig
	if len(results.ApplicationConfig) > 0 {
		appConfig.ApplicationConfig = results.ApplicationConfig
	}
	return appConfig, nil
}

func (c *client) SetApplicationConfig(ctx context.Context, controllerName, modelName, appName string, settings map[string]string) error {
	if controllerName == "" {
		currentController, err := c.clientStore.CurrentController()
		if err != nil {
			return err
		}
		controllerName = currentController
	}
	if modelName == "" {
		currentModel, err := c.clientStore.CurrentModel(controllerName)
		if err != nil {
			return err
		}
		modelName = currentModel
	}
	appAPI, err := c.getApplicationAPI(ctx, controllerName, modelName)
	if err != nil {
		return err
	}
	return appAPI.SetConfig(ctx, appName, "", settings)
}

func (c *client) getApplicationAPI(ctx context.Context, controllerName, modelName string) (application.ApplicationAPI, error) {
	root, err := c.getAPIConn(ctx, controllerName, modelName)
	if err != nil {
		return nil, err
	}
	return applicationapi.NewClient(root), nil
}

func (c *client) getModelManagerClient(ctx context.Context, controllerName, modelName string) (*modelmanager.Client, error) {
	root, err := c.getAPIConn(ctx, controllerName, modelName)
	if err != nil {
		return nil, err
	}
	return modelmanager.NewClient(root), nil
}

func (c *client) getCloudClient(ctx context.Context, controllerName string) (*cloudapi.Client, error) {
	root, err := c.getAPIConn(ctx, controllerName, "")
	if err != nil {
		return nil, err
	}
	return cloudapi.NewClient(root), nil
}

func (c *client) getAPIClient(ctx context.Context, controllerName string, modelName string) (*apiclient.Client, error) {
	root, err := c.getAPIConn(ctx, controllerName, modelName)
	if err != nil {
		return nil, err
	}
	return apiclient.NewClient(root, c.logger), nil
}

func (c *client) getAPIConn(ctx context.Context, controllerName string, modelName string) (api.Connection, error) {
	accountDetails, err := c.clientStore.AccountDetails(controllerName)
	if err != nil {
		return nil, err
	}
	newAPIConnectionParams, err := c.GetNewAPIConnectionParams(
		c.clientStore,
		controllerName,
		modelName,
		accountDetails,
	)
	if err != nil {
		return nil, err
	}
	apiRoot, err := juju.NewAPIConnection(ctx, newAPIConnectionParams)
	if err != nil {
		return nil, err
	}
	return apiRoot, nil
}

func (c *client) GetNewAPIConnectionParams(
	store jclient.ClientStore,
	controllerName, modelName string,
	accountDetails *jclient.AccountDetails,
) (juju.NewAPIConnectionParams, error) {
	var modelUUID string
	if modelName != "" {
		modelDetails, err := store.ModelByName(controllerName, modelName)
		if err != nil {
			return juju.NewAPIConnectionParams{}, err
		}
		modelUUID = modelDetails.ModelUUID
	}
	dialOpts := api.DefaultDialOpts()
	return juju.NewAPIConnectionParams{
		ControllerStore: store,
		ControllerName:  controllerName,
		AccountDetails:  accountDetails,
		ModelUUID:       modelUUID,
		DialOpts:        dialOpts,
		OpenAPI:         modelcmd.OpenAPIFuncWithMacaroons(api.Open, store, controllerName),
	}, nil
}

func (c *client) AddModel(ctx context.Context, opts AddModelOptions) error {
	if opts.Name == "" {
		return errors.New("model name is required")
	}

	// Resolve controller
	controller, err := c.resolveController(opts.Controller)
	if err != nil {
		return err
	}

	// Check for existing model
	if exists, err := c.CheckForExistingModel(controller, opts.Name); err != nil {
		log.Warn().Err(err).Msg("Could not check for existing model")
	} else if exists {
		return errors.Errorf("model '%s' already exists", opts.Name)
	}

	// Resolve model owner
	modelOwner, err := c.resolveModelOwner(controller, opts.Owner)
	if err != nil {
		return err
	}

	// Resolve cloud and region
	cloudInfo, err := c.resolveCloudRegion(ctx, controller, opts.CloudRegion)
	if err != nil {
		return err
	}

	// Find credential
	credential, err := c.resolveCredential(ctx, controller, opts.Credential, cloudInfo, modelOwner)
	if err != nil {
		return err
	}

	// Create the model
	model, err := c.createModel(ctx, controller, opts.Name, modelOwner, cloudInfo, credential, opts.Config)
	if err != nil {
		return err
	}

	// Update client store
	if err := c.updateClientStore(controller, opts.Name, modelOwner, model, !opts.NoSwitch); err != nil {
		return err
	}

	log.Info().Str("model", opts.Name).Str("uuid", model.UUID).Msg("Model created successfully")
	return nil
}

// Helper types for cloud and credential information
type cloudInfo struct {
	tag    names.CloudTag
	cloud  jujucloud.Cloud
	region string
}

type credentialInfo struct {
	credential *jujucloud.Credential
	tag        names.CloudCredentialTag
	region     string
}

// resolveController gets the controller name, using current if not specified
func (c *client) resolveController(controllerName string) (string, error) {
	if controllerName != "" {
		return controllerName, nil
	}
	
	controllers, err := c.GetControllers()
	if err != nil {
		return "", err
	}
	return controllers.Current, nil
}

// resolveModelOwner determines the model owner, using current user if not specified
func (c *client) resolveModelOwner(controller, owner string) (string, error) {
	accountDetails, err := c.clientStore.AccountDetails(controller)
	if err != nil {
		return "", err
	}
	
	if owner == "" {
		return accountDetails.User, nil
	}
	
	if !names.IsValidUser(owner) {
		return "", errors.Errorf("%q is not a valid user name", owner)
	}
	return names.NewUserTag(owner).Id(), nil
}

// resolveCloudRegion determines cloud and region information
func (c *client) resolveCloudRegion(ctx context.Context, controller, cloudRegion string) (*cloudInfo, error) {
	cloudClient, err := c.getCloudClient(ctx, controller)
	if err != nil {
		return nil, err
	}
	
	adapter := &cloudClientAdapter{client: cloudClient, ctx: ctx}
	
	var cloudTag names.CloudTag
	var cloud jujucloud.Cloud
	var region string
	
	if cloudRegion != "" {
		cloudTag, cloud, region, err = GetCloudRegion(adapter, cloudRegion)
		if err != nil {
			return nil, errors.Annotate(err, "failed to get cloud region")
		}
	} else {
		cloudTag, cloud, err = MaybeGetControllerCloud(adapter)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	
	return &cloudInfo{
		tag:    cloudTag,
		cloud:  cloud,
		region: region,
	}, nil
}

// resolveCredential finds the appropriate credential
func (c *client) resolveCredential(ctx context.Context, controller, credentialName string, cloud *cloudInfo, modelOwner string) (*credentialInfo, error) {
	cloudClient, err := c.getCloudClient(ctx, controller)
	if err != nil {
		return nil, err
	}
	
	adapter := &cloudClientAdapter{client: cloudClient, ctx: ctx}
	
	credentialObj, credentialTag, credentialRegion, err := FindCredential(adapter, c.clientStore, credentialName, &findCredentialParams{
		cloudTag:    cloud.tag,
		cloudRegion: cloud.region,
		cloud:       cloud.cloud,
		modelOwner:  modelOwner,
	})
	if err != nil {
		if credentialName != "" {
			return nil, errors.Errorf("credential '%s' not found for cloud '%s'", credentialName, cloud.tag.Id())
		}
		return nil, errors.Errorf("no credential found for cloud '%s'. Use 'juju add-credential %s' to add one", cloud.tag.Id(), cloud.tag.Id())
	}
	
	// Upload credential if needed
	if credentialName != "" && credentialObj != nil {
		if err := adapter.AddCredential(credentialTag.String(), *credentialObj); err != nil {
			return nil, errors.Annotate(err, "failed to upload credential")
		}
	}
	
	// Use credential region if no cloud region specified
	finalRegion := cloud.region
	if finalRegion == "" {
		finalRegion = credentialRegion
	}
	
	return &credentialInfo{
		credential: credentialObj,
		tag:        credentialTag,
		region:     finalRegion,
	}, nil
}

// createModel performs the actual model creation
func (c *client) createModel(ctx context.Context, controller, modelName, modelOwner string, cloud *cloudInfo, credential *credentialInfo, config map[string]string) (*modelmanager.Model, error) {
	// Convert config to interface map
	attrs := make(map[string]interface{})
	for k, v := range config {
		attrs[k] = v
	}
	
	modelManagerClient, err := c.getModelManagerClient(ctx, controller, "")
	if err != nil {
		return nil, err
	}
	
	model, err := modelManagerClient.CreateModel(ctx, modelName, modelOwner, cloud.tag.Id(), credential.region, credential.tag, attrs)
	if err != nil {
		return nil, err
	}
	
	return model, nil
}

// updateClientStore updates the local client store with model information
func (c *client) updateClientStore(controller, modelName, modelOwner string, model *modelmanager.Model, switchToModel bool) error {
	accountDetails, err := c.clientStore.AccountDetails(controller)
	if err != nil {
		return err
	}
	
	// Only update store if current user owns the model
	if modelOwner != accountDetails.User {
		return nil
	}
	
	details := jclient.ModelDetails{
		ModelUUID: model.UUID,
		ModelType: model.Type,
	}
	
	// Create qualified model name for client store
	ownerName := names.NewUserTag(modelOwner).Name()
	qualifiedModelName := fmt.Sprintf("%s/%s", ownerName, modelName)
	
	if err := c.clientStore.UpdateModel(controller, qualifiedModelName, details); err != nil {
		return errors.Trace(err)
	}
	
	if switchToModel {
		if err := c.clientStore.SetCurrentModel(controller, qualifiedModelName); err != nil {
			return errors.Trace(err)
		}
	}
	
	return nil
}


func (c *client) getCloudRegion(ctx context.Context, cloudClient *cloudapi.Client, cloudRegionStr string) (cloudTag names.CloudTag, cloud jujucloud.Cloud, cloudRegion string, err error) {
	var cloudName string

	sep := strings.IndexRune(cloudRegionStr, '/')
	if sep >= 0 {
		// Cloud and region specified (e.g., "aws/us-east-1")
		cloudName, cloudRegion = cloudRegionStr[:sep], cloudRegionStr[sep+1:]
		if !names.IsValidCloud(cloudName) {
			err = errors.NotValidf("cloud name %q", cloudName)
			return
		}
		cloudTag = names.NewCloudTag(cloudName)
		if cloud, err = cloudClient.Cloud(ctx, cloudTag); err != nil {
			if params.IsCodeNotFound(err) {
				clouds, err2 := cloudClient.Clouds(ctx)
				if err2 != nil {
					err = errors.Trace(err2)
					return
				}
				err = c.unsupportedCloudOrRegionError(clouds, cloudRegionStr)
			}
			return
		}

		// Validate region exists in the cloud
		if _, err = jujucloud.RegionByName(cloud.Regions, cloudRegion); err != nil {
			clouds, err2 := cloudClient.Clouds(ctx)
			if err2 != nil {
				err = errors.Trace(err2)
				return
			}
			err = c.unsupportedCloudOrRegionError(clouds, cloudRegionStr)
			return
		}
	} else {
		// Only cloud or only region specified - need to determine which
		if names.IsValidCloud(cloudRegionStr) {
			// It's a valid cloud name
			cloudName = cloudRegionStr
			cloudTag = names.NewCloudTag(cloudName)
			if cloud, err = cloudClient.Cloud(ctx, cloudTag); err != nil {
				if params.IsCodeNotFound(err) {
					clouds, err2 := cloudClient.Clouds(ctx)
					if err2 != nil {
						err = errors.Trace(err2)
						return
					}
					err = c.unsupportedCloudOrRegionError(clouds, cloudRegionStr)
				}
				return
			}
		} else {
			// Assume it's a region in the default controller cloud
			controllerCloud, err2 := c.maybeGetControllerCloud(ctx, cloudClient)
			if err2 != nil {
				err = errors.Trace(err2)
				return
			}
			cloud = controllerCloud.Cloud
			cloudTag = controllerCloud.CloudTag
			cloudRegion = cloudRegionStr

			// Validate region exists in the cloud
			if _, err = jujucloud.RegionByName(cloud.Regions, cloudRegion); err != nil {
				clouds, err2 := cloudClient.Clouds(ctx)
				if err2 != nil {
					err = errors.Trace(err2)
					return
				}
				err = c.unsupportedCloudOrRegionError(clouds, cloudRegionStr)
				return
			}
		}
	}

	return cloudTag, cloud, cloudRegion, nil
}

type controllerCloudInfo struct {
	Cloud    jujucloud.Cloud
	CloudTag names.CloudTag
}

func (c *client) maybeGetControllerCloud(ctx context.Context, cloudClient *cloudapi.Client) (controllerCloudInfo, error) {
	clouds, err := cloudClient.Clouds(ctx)
	if err != nil {
		return controllerCloudInfo{}, errors.Trace(err)
	}

	if len(clouds) == 0 {
		return controllerCloudInfo{}, errors.New("no clouds available")
	}

	if len(clouds) > 1 {
		return controllerCloudInfo{}, errors.New("multiple clouds available, please specify cloud/region")
	}

	// Single cloud available - use it as default
	for cloudTag, cloud := range clouds {
		return controllerCloudInfo{
			Cloud:    cloud,
			CloudTag: cloudTag,
		}, nil
	}

	return controllerCloudInfo{}, errors.New("no clouds available")
}

func (c *client) unsupportedCloudOrRegionError(clouds map[names.CloudTag]jujucloud.Cloud, cloudRegion string) error {
	var availableClouds []string
	var availableRegions []string

	for cloudTag, cloud := range clouds {
		cloudName := cloudTag.Id()
		availableClouds = append(availableClouds, cloudName)
		for _, region := range cloud.Regions {
			availableRegions = append(availableRegions, fmt.Sprintf("%s/%s", cloudName, region.Name))
		}
	}

	msg := fmt.Sprintf("cloud or region %q not found", cloudRegion)
	if len(availableClouds) > 0 {
		msg += fmt.Sprintf("\nAvailable clouds: %s", strings.Join(availableClouds, ", "))
	}
	if len(availableRegions) > 0 {
		msg += fmt.Sprintf("\nAvailable cloud/regions: %s", strings.Join(availableRegions, ", "))
	}

	return errors.New(msg)
}

// GetCloudRegion is a standalone function refactored from juju's cmd/juju/controller/addmodel.go
// It parses and validates a cloud/region string, returning the cloud tag, cloud info, and region name.
func GetCloudRegion(cloudClient CloudAPI, cloudRegionStr string) (cloudTag names.CloudTag, cloud jujucloud.Cloud, cloudRegion string, err error) {
	fail := func(err error) (names.CloudTag, jujucloud.Cloud, string, error) {
		return names.CloudTag{}, jujucloud.Cloud{}, "", err
	}

	var cloudName string
	sep := strings.IndexRune(cloudRegionStr, '/')
	if sep >= 0 {
		// User specified "cloud/region".
		cloudName, cloudRegion = cloudRegionStr[:sep], cloudRegionStr[sep+1:]
		if !names.IsValidCloud(cloudName) {
			return fail(errors.NotValidf("cloud name %q", cloudName))
		}
		cloudTag = names.NewCloudTag(cloudName)
		if cloud, err = cloudClient.Cloud(cloudTag); err != nil {
			return fail(errors.Trace(err))
		}
	} else {
		// User specified "cloud" or "region". We'll try first
		// for cloud (check if it's a valid cloud name, and
		// whether there is a cloud by that name), and then
		// as a region within the default cloud.
		if names.IsValidCloud(cloudRegionStr) {
			cloudName = cloudRegionStr
		} else {
			cloudRegion = cloudRegionStr
		}
		if cloudName != "" {
			cloudTag = names.NewCloudTag(cloudName)
			cloud, err = cloudClient.Cloud(cloudTag)
			if params.IsCodeNotFound(err) {
				// No such cloud with the specified name,
				// so we'll try the name as a region in
				// the default cloud.
				cloudRegion, cloudName = cloudName, ""
			} else if err != nil {
				return fail(errors.Trace(err))
			}
		}
		if cloudName == "" {
			cloudTag, cloud, err = MaybeGetControllerCloud(cloudClient)
			if err != nil {
				return fail(errors.Trace(err))
			}
		}
	}
	if cloudRegion != "" {
		// A region has been specified, make sure it exists.
		if _, err := jujucloud.RegionByName(cloud.Regions, cloudRegion); err != nil {
			if cloudRegion == cloudRegionStr {
				// The string is not in the format cloud/region,
				// so we should tell that the user that it is
				// neither a cloud nor a region in the
				// controller's cloud.
				clouds, err := cloudClient.Clouds()
				if err != nil {
					return fail(errors.Annotate(err, "querying supported clouds"))
				}
				return fail(UnsupportedCloudOrRegionError(clouds, cloudRegionStr))
			}
			return fail(errors.Trace(err))
		}
	}
	return cloudTag, cloud, cloudRegion, nil
}

// CloudAPI interface for cloud operations
type CloudAPI interface {
	Clouds() (map[names.CloudTag]jujucloud.Cloud, error)
	Cloud(tag names.CloudTag) (jujucloud.Cloud, error)
	UserCredentials(user names.UserTag, cloud names.CloudTag) ([]names.CloudCredentialTag, error)
	AddCredential(tag string, credential jujucloud.Credential) error
}

// cloudClientAdapter adapts *cloudapi.Client to the CloudAPI interface
type cloudClientAdapter struct {
	client *cloudapi.Client
	ctx    context.Context
}

func (a *cloudClientAdapter) Clouds() (map[names.CloudTag]jujucloud.Cloud, error) {
	return a.client.Clouds(a.ctx)
}

func (a *cloudClientAdapter) Cloud(tag names.CloudTag) (jujucloud.Cloud, error) {
	return a.client.Cloud(a.ctx, tag)
}

func (a *cloudClientAdapter) UserCredentials(user names.UserTag, cloud names.CloudTag) ([]names.CloudCredentialTag, error) {
	return a.client.UserCredentials(a.ctx, user, cloud)
}

func (a *cloudClientAdapter) AddCredential(tag string, credential jujucloud.Credential) error {
	return a.client.AddCredential(a.ctx, tag, credential)
}

// MaybeGetControllerCloud gets the controller's default cloud
func MaybeGetControllerCloud(cloudClient CloudAPI) (names.CloudTag, jujucloud.Cloud, error) {
	clouds, err := cloudClient.Clouds()
	if err != nil {
		return names.CloudTag{}, jujucloud.Cloud{}, errors.Trace(err)
	}
	if len(clouds) != 1 {
		return names.CloudTag{}, jujucloud.Cloud{}, UnsupportedCloudOrRegionError(clouds, "")
	}
	for cloudTag, cloud := range clouds {
		return cloudTag, cloud, nil
	}
	panic("unreachable")
}

// UnsupportedCloudOrRegionError creates an error for unsupported cloud/region combinations
func UnsupportedCloudOrRegionError(clouds map[names.CloudTag]jujucloud.Cloud, cloudRegion string) error {
	return errors.New("cloud or region not supported")
}

// findCredentialParams holds parameters for credential finding
type findCredentialParams struct {
	cloudTag    names.CloudTag
	cloud       jujucloud.Cloud
	cloudRegion string
	modelOwner  string
}

// FindCredential finds a suitable credential to use for the new model.
// The credential will first be searched for locally and then on the
// controller. If a credential is found locally then it's value will be
// returned as the first return value. If it is found on the controller
// this will be nil as there is no need to upload it in that case.
func FindCredential(cloudClient CloudAPI, clientStore jclient.ClientStore, credentialName string, p *findCredentialParams) (_ *jujucloud.Credential, _ names.CloudCredentialTag, cloudRegion string, _ error) {
	if credentialName == "" {
		return findUnspecifiedCredential(cloudClient, clientStore, p)
	}
	return findSpecifiedCredential(cloudClient, clientStore, credentialName, p)
}

func findUnspecifiedCredential(cloudClient CloudAPI, clientStore jclient.ClientStore, p *findCredentialParams) (_ *jujucloud.Credential, _ names.CloudCredentialTag, cloudRegion string, _ error) {
	fail := func(err error) (*jujucloud.Credential, names.CloudCredentialTag, string, error) {
		return nil, names.CloudCredentialTag{}, "", err
	}
	
	log.Debug().Str("cloud", p.cloudTag.Id()).Strs("authTypes", authTypesToStrings(p.cloud.AuthTypes)).Msg("Finding unspecified credential")
	
	// If the user has not specified a credential, and the cloud advertises
	// itself as supporting the "empty" auth-type, then return immediately.
	for _, authType := range p.cloud.AuthTypes {
		if authType == jujucloud.EmptyAuthType {
			log.Debug().Str("cloud", p.cloudTag.Id()).Msg("Cloud supports empty auth, no credential needed")
			return nil, names.CloudCredentialTag{}, p.cloudRegion, nil
		}
	}

	// No credential has been specified, so see if there is one already on the controller we can use.
	modelOwnerTag := names.NewUserTag(p.modelOwner)
	credentialTags, err := cloudClient.UserCredentials(modelOwnerTag, p.cloudTag)
	if err != nil {
		log.Debug().Err(err).Str("cloud", p.cloudTag.Id()).Str("user", p.modelOwner).Msg("Failed to get user credentials from controller")
		// Don't fail immediately, try to find local credentials
	} else {
		var credentialTag names.CloudCredentialTag
		if len(credentialTags) == 1 {
			credentialTag = credentialTags[0]
		}

		if (credentialTag != names.CloudCredentialTag{}) {
			log.Debug().Str("credential", credentialTag.Id()).Msg("Found controller credential")
			// If the controller already has a credential, see if
			// there is a local version that has an associated
			// region.
			credential, _, cloudRegion, err := findLocalCredential(clientStore, p, credentialTag.Name())
			if errors.IsNotFound(err) {
				// No local credential; use the region
				// specified by the user, if any.
				cloudRegion = p.cloudRegion
			} else if err != nil {
				log.Debug().Err(err).Str("credential", credentialTag.Id()).Msg("Error finding local version of controller credential")
			}
			// If there is a credential in the controller use it even if we don't have a local version.
			return credential, credentialTag, cloudRegion, nil
		}
		
		if len(credentialTags) > 1 {
			log.Debug().Int("count", len(credentialTags)).Msg("Multiple credentials available on controller")
		}
	}
	
	// There is not a default credential on the controller (either
	// there are no credentials, or there is more than one). Look for
	// a local credential we might use.
	credential, credentialName, cloudRegion, err := findLocalCredential(clientStore, p, "")
	if err != nil {
		log.Debug().Err(err).Str("cloud", p.cloudTag.Id()).Msg("No local credential found")
		
		// If cloud supports empty auth as a fallback, use it
		for _, authType := range p.cloud.AuthTypes {
			if authType == jujucloud.EmptyAuthType {
				log.Debug().Str("cloud", p.cloudTag.Id()).Msg("Falling back to empty auth")
				return nil, names.CloudCredentialTag{}, p.cloudRegion, nil
			}
		}
		
		return fail(errors.Trace(err))
	}
	
	// We've got a local credential to use.
	credentialTag, err := common.ResolveCloudCredentialTag(
		modelOwnerTag, p.cloudTag, credentialName,
	)
	if err != nil {
		return fail(errors.Trace(err))
	}
	
	log.Debug().Str("credential", credentialName).Str("cloud", p.cloudTag.Id()).Msg("Using local credential")
	return credential, credentialTag, cloudRegion, nil
}

func authTypesToStrings(authTypes []jujucloud.AuthType) []string {
	result := make([]string, len(authTypes))
	for i, authType := range authTypes {
		result[i] = string(authType)
	}
	return result
}

func findSpecifiedCredential(cloudClient CloudAPI, clientStore jclient.ClientStore, credentialName string, p *findCredentialParams) (_ *jujucloud.Credential, _ names.CloudCredentialTag, cloudRegion string, _ error) {
	fail := func(err error) (*jujucloud.Credential, names.CloudCredentialTag, string, error) {
		return nil, names.CloudCredentialTag{}, "", err
	}
	// Look for a local credential with the specified name
	credential, credentialNameFound, cloudRegion, err := findLocalCredential(clientStore, p, credentialName)
	if err != nil && !errors.IsNotFound(err) {
		return fail(errors.Trace(err))
	}
	if credential != nil {
		// We found a local credential with the specified name.
		modelOwnerTag := names.NewUserTag(p.modelOwner)
		credentialTag, err := common.ResolveCloudCredentialTag(
			modelOwnerTag, p.cloudTag, credentialNameFound,
		)
		if err != nil {
			return fail(errors.Trace(err))
		}
		return credential, credentialTag, cloudRegion, nil
	}

	// There was no local credential with that name, check the controller
	modelOwnerTag := names.NewUserTag(p.modelOwner)
	credentialTags, err := cloudClient.UserCredentials(modelOwnerTag, p.cloudTag)
	if err != nil {
		return fail(errors.Trace(err))
	}
	credentialTag, err := common.ResolveCloudCredentialTag(
		modelOwnerTag, p.cloudTag, credentialName,
	)
	if err != nil {
		return fail(errors.Trace(err))
	}
	credentialId := credentialTag.Id()
	for _, tag := range credentialTags {
		if tag.Id() != credentialId {
			continue
		}
		// Found it on controller, no need to upload
		return nil, credentialTag, "", nil
	}
	// Cannot find a credential with the correct name
	return fail(errors.NotFoundf("credential '%s'", credentialName))
}

func findLocalCredential(clientStore jclient.ClientStore, p *findCredentialParams, name string) (_ *jujucloud.Credential, credentialName, cloudRegion string, _ error) {
	fail := func(err error) (*jujucloud.Credential, string, string, error) {
		return nil, "", "", err
	}
	
	// First try to get credentials directly from the client store
	credentials, err := clientStore.AllCredentials()
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get credentials from client store")
		return fail(errors.Trace(err))
	}
	
	cloudName := p.cloudTag.Id()
	if cloudCreds, exists := credentials[cloudName]; exists {
		if name != "" {
			// Looking for a specific credential
			if cred, found := cloudCreds.AuthCredentials[name]; found {
				log.Debug().Str("credential", name).Str("cloud", cloudName).Msg("Found specific credential")
				return &cred, name, p.cloudRegion, nil
			}
		} else {
			// Looking for any credential for this cloud
			for credName, cred := range cloudCreds.AuthCredentials {
				log.Debug().Str("credential", credName).Str("cloud", cloudName).Msg("Found credential")
				return &cred, credName, p.cloudRegion, nil
			}
		}
	}
	
	// If no credentials found in store, try provider detection
	providerRegistry := environs.GlobalProviderRegistry()
	provider, err := providerRegistry.Provider(p.cloud.Type)
	if err != nil {
		log.Debug().Err(err).Str("cloudType", p.cloud.Type).Msg("Failed to get provider")
		return fail(errors.Trace(err))
	}
	
	credential, credentialName, cloudRegion, _, err := common.GetOrDetectCredential(
		nil, clientStore, provider, modelcmd.GetCredentialsParams{
			Cloud:          p.cloud,
			CloudRegion:    p.cloudRegion,
			CredentialName: name,
		},
	)
	if err == nil {
		log.Debug().Str("credential", credentialName).Str("cloud", cloudName).Msg("Detected credential")
		return credential, credentialName, cloudRegion, nil
	}
	
	log.Debug().Err(err).Str("cloud", cloudName).Str("requestedCredential", name).Msg("Credential detection failed")
	
	switch errors.Cause(err) {
	case modelcmd.ErrMultipleCredentials:
		return fail(errors.New("more than one credential is available"))
	case common.ErrMultipleDetectedCredentials:
		return fail(errors.New("more than one credential detected"))
	}
	return fail(errors.Trace(err))
}

// CheckForExistingModel checks if a model with the given name already exists
func (c *client) CheckForExistingModel(currentController, modelName string) (bool, error) {
	models, err := c.GetModels(currentController)
	if err != nil {
		return false, err
	}
	
	// Check both qualified and unqualified names
	for modelKey := range models.ModelDetails {
		// Check direct match
		if modelKey == modelName {
			return true, nil
		}
		// Check if the model name matches the unqualified part
		if strings.Contains(modelKey, "/") {
			parts := strings.Split(modelKey, "/")
			if len(parts) == 2 && parts[1] == modelName {
				return true, nil
			}
		}
	}
	return false, nil
}

