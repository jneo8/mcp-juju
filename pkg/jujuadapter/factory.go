package jujuadapter

import (
	"fmt"
	"os"
	"time"

	cloudfile "github.com/juju/juju/cloud"
	"github.com/juju/juju/jujuclient"

	"github.com/juju/juju/cmd/juju/action"
	"github.com/juju/juju/cmd/juju/application"
	"github.com/juju/juju/cmd/juju/backups"
	"github.com/juju/juju/cmd/juju/block"
	"github.com/juju/juju/cmd/juju/caas"
	"github.com/juju/juju/cmd/juju/charmhub"
	"github.com/juju/juju/cmd/juju/cloud"
	"github.com/juju/juju/cmd/juju/commands"
	"github.com/juju/juju/cmd/juju/controller"
	"github.com/juju/juju/cmd/juju/crossmodel"
	"github.com/juju/juju/cmd/juju/dashboard"
	"github.com/juju/juju/cmd/juju/firewall"
	"github.com/juju/juju/cmd/juju/machine"
	"github.com/juju/juju/cmd/juju/model"
	"github.com/juju/juju/cmd/juju/resource"
	"github.com/juju/juju/cmd/juju/secretbackends"
	"github.com/juju/juju/cmd/juju/secrets"
	"github.com/juju/juju/cmd/juju/space"
	"github.com/juju/juju/cmd/juju/ssh"
	"github.com/juju/juju/cmd/juju/status"
	"github.com/juju/juju/cmd/juju/storage"
	"github.com/juju/juju/cmd/juju/subnet"
	"github.com/juju/juju/cmd/juju/user"
	"github.com/juju/juju/cmd/juju/waitfor"
	"github.com/juju/cmd/v3"
)

type CommandFactory interface {
	GetCommand(id JujuCommandID) (Command, error)
	GetCommandByName(name string) (Command, error)
}

type commandFactory struct{}

type cloudToCommandAdaptor struct{}

func (cloudToCommandAdaptor) ReadCloudData(path string) ([]byte, error) {
	return os.ReadFile(path)
}
func (cloudToCommandAdaptor) ParseOneCloud(data []byte) (cloudfile.Cloud, error) {
	return cloudfile.ParseOneCloud(data)
}
func (cloudToCommandAdaptor) PublicCloudMetadata(searchPaths ...string) (map[string]cloudfile.Cloud, bool, error) {
	return cloudfile.PublicCloudMetadata(searchPaths...)
}
func (cloudToCommandAdaptor) PersonalCloudMetadata() (map[string]cloudfile.Cloud, error) {
	return cloudfile.PersonalCloudMetadata()
}
func (cloudToCommandAdaptor) WritePersonalCloudMetadata(cloudsMap map[string]cloudfile.Cloud) error {
	return cloudfile.WritePersonalCloudMetadata(cloudsMap)
}

func (c *commandFactory) GetCommand(id JujuCommandID) (Command, error) {
	def, exists := GetCommandDefinition(id)
	if !exists {
		return nil, fmt.Errorf("unknown command: %s", id)
	}

	jujuCmd, err := c.createJujuCommand(id)
	if err != nil {
		return nil, err
	}

	if jujuCmd == nil {
		return nil, fmt.Errorf("failed to create command: %s", id)
	}

	return &command{
		cmd:          jujuCmd,
		info:         jujuCmd.Info(),
		t:            time.Now(),
		disabledArgs: def.HTTPMcpServerDisableArgs, // TODO: Make this configurable based on server type
	}, nil
}

func (c *commandFactory) createJujuCommand(id JujuCommandID) (cmd.Command, error) {
	cloudAdapter := &cloudToCommandAdaptor{}
	
	switch id {
	// Reporting commands
	case CmdStatus:
		return status.NewStatusCommand(), nil
	case CmdStatusHistory:
		return status.NewStatusHistoryCommand(), nil
	case CmdSwitch:
		return commands.NewCommandByName(string(id))
	case CmdVersion:
		return commands.NewCommandByName(string(id))

	// Application commands
	case CmdAddUnit:
		return application.NewAddUnitCommand(), nil
	case CmdConfig:
		return application.NewConfigCommand(), nil
	case CmdDeploy:
		return application.NewDeployCommand(), nil
	case CmdExpose:
		return application.NewExposeCommand(), nil
	case CmdUnexpose:
		return application.NewUnexposeCommand(), nil
	case CmdGetConstraints:
		return application.NewApplicationGetConstraintsCommand(), nil
	case CmdSetConstraints:
		return application.NewApplicationSetConstraintsCommand(), nil
	case CmdDiffBundle:
		return application.NewDiffBundleCommand(), nil
	case CmdShowApplication:
		return application.NewShowApplicationCommand(), nil
	case CmdShowUnit:
		return application.NewShowUnitCommand(), nil
	case CmdRefresh:
		return application.NewRefreshCommand(), nil
	case CmdBind:
		return application.NewBindCommand(), nil
	case CmdScaleApplication:
		return application.NewScaleApplicationCommand(), nil
	case CmdTrust:
		return application.NewTrustCommand(), nil
	case CmdAddRelation:
		return application.NewAddRelationCommand(), nil
	case CmdRemoveRelation:
		return application.NewRemoveRelationCommand(), nil
	case CmdRemoveApplication:
		return application.NewRemoveApplicationCommand(), nil
	case CmdRemoveUnit:
		return application.NewRemoveUnitCommand(), nil
	case CmdRemoveSaas:
		return application.NewRemoveSaasCommand(), nil
	case CmdConsume:
		return application.NewConsumeCommand(), nil
	case CmdSuspendRelation:
		return application.NewSuspendRelationCommand(), nil
	case CmdResumeRelation:
		return application.NewResumeRelationCommand(), nil
	case CmdResolved:
		return application.NewResolvedCommand(), nil

	// Machine commands
	case CmdAddMachine:
		return machine.NewAddCommand(), nil
	case CmdRemoveMachine:
		return machine.NewRemoveCommand(), nil
	case CmdMachines:
		return machine.NewListMachinesCommand(), nil
	case CmdShowMachine:
		return machine.NewShowMachineCommand(), nil
	case CmdUpgradeMachine:
		return machine.NewUpgradeMachineCommand(), nil

	// Model commands
	case CmdModelConfig:
		return model.NewConfigCommand(), nil
	case CmdModelDefaults:
		return model.NewDefaultsCommand(), nil
	case CmdRetryProvisioning:
		return model.NewRetryProvisioningCommand(), nil
	case CmdDestroyModel:
		return model.NewDestroyCommand(), nil
	case CmdGrant:
		return model.NewGrantCommand(), nil
	case CmdRevoke:
		return model.NewRevokeCommand(), nil
	case CmdShowModel:
		return model.NewShowCommand(), nil
	case CmdModelCredential:
		return model.NewModelCredentialCommand(), nil
	case CmdExportBundle:
		return model.NewExportBundleCommand(), nil
	case CmdGrantCloud:
		return model.NewGrantCloudCommand(), nil
	case CmdRevokeCloud:
		return model.NewRevokeCloudCommand(), nil

	// Commands from main commands package
	case CmdBootstrap, CmdMigrate, CmdSyncAgentBinary, CmdUpgradeModel, CmdUpgradeController,
		 CmdHelpHooks, CmdHelpActions, CmdDebugLog, CmdEnableHa:
		return commands.NewCommandByName(string(id))

	// Controller commands
	case CmdAddModel:
		return controller.NewAddModelCommand(), nil
	case CmdDestroyController:
		return controller.NewDestroyCommand(), nil
	case CmdModels:
		return controller.NewListModelsCommand(), nil
	case CmdKillController:
		return controller.NewKillCommand(), nil
	case CmdControllers:
		return controller.NewListControllersCommand(), nil
	case CmdRegister:
		return controller.NewRegisterCommand(), nil
	case CmdUnregister:
		return controller.NewUnregisterCommand(jujuclient.NewFileClientStore()), nil
	case CmdEnableDestroyController:
		return controller.NewEnableDestroyControllerCommand(), nil
	case CmdShowController:
		return controller.NewShowControllerCommand(), nil
	case CmdControllerConfig:
		return controller.NewConfigCommand(), nil

	// Action commands
	case CmdActions:
		return action.NewListCommand(), nil
	case CmdShowAction:
		return action.NewShowCommand(), nil
	case CmdCancelAction:
		return action.NewCancelCommand(), nil
	case CmdRun:
		return action.NewRunCommand(), nil
	case CmdOperations:
		return action.NewListOperationsCommand(), nil
	case CmdShowOperation:
		return action.NewShowOperationCommand(), nil
	case CmdShowTask:
		return action.NewShowTaskCommand(), nil
	case CmdExec:
		return action.NewExecCommand(nil), nil

	// Cloud and credential commands
	case CmdUpdateCloud:
		return cloud.NewUpdateCloudCommand(cloudAdapter), nil
	case CmdUpdatePublicClouds:
		return cloud.NewUpdatePublicCloudsCommand(), nil
	case CmdClouds:
		return cloud.NewListCloudsCommand(), nil
	case CmdRegions:
		return cloud.NewListRegionsCommand(), nil
	case CmdShowCloud:
		return cloud.NewShowCloudCommand(), nil
	case CmdAddCloud:
		return cloud.NewAddCloudCommand(cloudAdapter), nil
	case CmdRemoveCloud:
		return cloud.NewRemoveCloudCommand(), nil
	case CmdCredentials:
		return cloud.NewListCredentialsCommand(), nil
	case CmdDetectCredentials:
		return cloud.NewDetectCredentialsCommand(), nil
	case CmdSetDefaultRegion:
		return cloud.NewSetDefaultRegionCommand(), nil
	case CmdSetDefaultCredential:
		return cloud.NewSetDefaultCredentialCommand(), nil
	case CmdAddCredential:
		return cloud.NewAddCredentialCommand(), nil
	case CmdRemoveCredential:
		return cloud.NewRemoveCredentialCommand(), nil
	case CmdUpdateCredential:
		return cloud.NewUpdateCredentialCommand(), nil
	case CmdShowCredential:
		return cloud.NewShowCredentialCommand(), nil

	// User commands
	case CmdAddUser:
		return user.NewAddCommand(), nil
	case CmdChangePassword:
		return user.NewChangePasswordCommand(), nil
	case CmdShowUser:
		return user.NewShowUserCommand(), nil
	case CmdUsers:
		return user.NewListCommand(), nil
	case CmdEnableUser:
		return user.NewEnableCommand(), nil
	case CmdDisableUser:
		return user.NewDisableCommand(), nil
	case CmdLogin:
		return user.NewLoginCommand(), nil
	case CmdLogout:
		return user.NewLogoutCommand(), nil
	case CmdRemoveUser:
		return user.NewRemoveCommand(), nil  
	case CmdWhoami:
		return user.NewWhoAmICommand(), nil

	// Storage commands
	case CmdAddStorage:
		return storage.NewAddCommand(), nil
	case CmdStorage:
		return storage.NewListCommand(), nil
	case CmdCreateStoragePool:
		return storage.NewPoolCreateCommand(), nil
	case CmdStoragePools:
		return storage.NewPoolListCommand(), nil
	case CmdRemoveStoragePool:
		return storage.NewPoolRemoveCommand(), nil
	case CmdUpdateStoragePool:
		return storage.NewPoolUpdateCommand(), nil
	case CmdShowStorage:
		return storage.NewShowCommand(), nil
	case CmdRemoveStorage:
		return storage.NewRemoveStorageCommandWithAPI(), nil
	case CmdDetachStorage:
		return storage.NewDetachStorageCommandWithAPI(), nil
	case CmdAttachStorage:
		return storage.NewAttachStorageCommandWithAPI(), nil
	case CmdImportFilesystem:
		return storage.NewImportFilesystemCommand(storage.NewStorageImporter, nil), nil

	// Space commands
	case CmdAddSpace:
		return space.NewAddCommand(), nil
	case CmdSpaces:
		return space.NewListCommand(), nil
	case CmdMoveToSpace:
		return space.NewMoveCommand(), nil
	case CmdReloadSpaces:
		return space.NewReloadCommand(), nil
	case CmdShowSpace:
		return space.NewShowSpaceCommand(), nil
	case CmdRemoveSpace:
		return space.NewRemoveCommand(), nil
	case CmdRenameSpace:
		return space.NewRenameCommand(), nil

	// Subnet commands
	case CmdSubnets:
		return subnet.NewListCommand(), nil

	// SSH and debugging commands
	case CmdScp:
		return ssh.NewSCPCommand(nil, ssh.DefaultSSHRetryStrategy, ssh.DefaultSSHPublicKeyRetryStrategy), nil
	case CmdSsh:
		return ssh.NewSSHCommand(nil, nil, ssh.DefaultSSHRetryStrategy, ssh.DefaultSSHPublicKeyRetryStrategy), nil
	case CmdDebugHooks:
		return ssh.NewDebugHooksCommand(nil, ssh.DefaultSSHRetryStrategy, ssh.DefaultSSHPublicKeyRetryStrategy), nil
	case CmdDebugCode:
		return ssh.NewDebugCodeCommand(nil, ssh.DefaultSSHRetryStrategy, ssh.DefaultSSHPublicKeyRetryStrategy), nil

	// SSH keys commands
	case CmdAddSshKey:
		return commands.NewAddKeysCommand(), nil
	case CmdRemoveSshKey:
		return commands.NewRemoveKeysCommand(), nil
	case CmdImportSshKey:
		return commands.NewImportKeysCommand(), nil
	case CmdSshKeys:
		return commands.NewListKeysCommand(), nil

	// Backup commands
	case CmdCreateBackup:
		return backups.NewCreateCommand(), nil
	case CmdDownloadBackup:
		return backups.NewDownloadCommand(), nil

	// Block commands
	case CmdDisableCommand:
		return block.NewDisableCommand(), nil
	case CmdDisabledCommands:
		return block.NewListCommand(), nil
	case CmdEnableCommand:
		return block.NewEnableCommand(), nil

	// Firewall commands
	case CmdSetFirewallRule:
		return firewall.NewSetFirewallRuleCommand(), nil
	case CmdListFirewallRules:
		return firewall.NewListFirewallRulesCommand(), nil

	// Cross model commands
	case CmdOffer:
		return crossmodel.NewOfferCommand(), nil
	case CmdRemoveOffer:
		return crossmodel.NewRemoveOfferCommand(), nil
	case CmdShowOfferedEndpoint:
		return crossmodel.NewShowOfferedEndpointCommand(), nil
	case CmdListEndpoints:
		return crossmodel.NewListEndpointsCommand(), nil
	case CmdFindEndpoints:
		return crossmodel.NewFindEndpointsCommand(), nil

	// CAAS commands
	case CmdAddK8s:
		return caas.NewAddCAASCommand(cloudAdapter), nil
	case CmdUpdateK8s:
		return caas.NewUpdateCAASCommand(cloudAdapter), nil
	case CmdRemoveK8s:
		return caas.NewRemoveCAASCommand(cloudAdapter), nil

	// Dashboard commands
	case CmdDashboard:
		return dashboard.NewDashboardCommand(), nil

	// Resource commands
	case CmdAttachResource:
		return resource.NewUploadCommand(), nil
	case CmdResources:
		return resource.NewListCommand(), nil
	case CmdCharmResources:
		return resource.NewCharmResourcesCommand(), nil

	// CharmHub commands
	case CmdInfo:
		return charmhub.NewInfoCommand(), nil
	case CmdFind:
		return charmhub.NewFindCommand(), nil
	case CmdDownload:
		return charmhub.NewDownloadCommand(), nil

	// Secrets commands
	case CmdSecrets:
		return secrets.NewListSecretsCommand(), nil
	case CmdShowSecret:
		return secrets.NewShowSecretsCommand(), nil
	case CmdAddSecret:
		return secrets.NewAddSecretCommand(), nil
	case CmdUpdateSecret:
		return secrets.NewUpdateSecretCommand(), nil
	case CmdRemoveSecret:
		return secrets.NewRemoveSecretCommand(), nil
	case CmdGrantSecret:
		return secrets.NewGrantSecretCommand(), nil
	case CmdRevokeSecret:
		return secrets.NewRevokeSecretCommand(), nil

	// Secret backends commands
	case CmdSecretBackends:
		return secretbackends.NewListSecretBackendsCommand(), nil
	case CmdAddSecretBackend:
		return secretbackends.NewAddSecretBackendCommand(), nil
	case CmdUpdateSecretBackend:
		return secretbackends.NewUpdateSecretBackendCommand(), nil
	case CmdRemoveSecretBackend:
		return secretbackends.NewRemoveSecretBackendCommand(), nil
	case CmdShowSecretBackend:
		return secretbackends.NewShowSecretBackendCommand(), nil

	// Wait for commands
	case CmdWaitFor:
		return waitfor.NewWaitForCommand(), nil

	default:
		return nil, fmt.Errorf("unknown command: %s", id)
	}
}

func (c *commandFactory) GetCommandByName(name string) (Command, error) {
	id := JujuCommandID(name)
	return c.GetCommand(id)
}