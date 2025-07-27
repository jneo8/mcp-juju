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
	"github.com/juju/juju/cmd/juju/commands"
	"github.com/juju/juju/cmd/juju/waitfor"
	"github.com/juju/cmd/v3"
)

type CommandFactory interface {
	GetCommand(name string) (Command, error)
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

func (c *commandFactory) GetCommand(name string) (Command, error) {
	var jujuCmd cmd.Command
	
	switch name {
	// Reporting commands
	case "status":
		jujuCmd = status.NewStatusCommand()
	case "status-history":
		jujuCmd = status.NewStatusHistoryCommand()
	
	// Application commands
	case "add-unit":
		jujuCmd = application.NewAddUnitCommand()
	case "config":
		jujuCmd = application.NewConfigCommand()
	case "deploy":
		jujuCmd = application.NewDeployCommand()
	case "expose":
		jujuCmd = application.NewExposeCommand()
	case "unexpose":
		jujuCmd = application.NewUnexposeCommand()
	case "get-constraints":
		jujuCmd = application.NewApplicationGetConstraintsCommand()
	case "set-constraints":
		jujuCmd = application.NewApplicationSetConstraintsCommand()
	case "diff-bundle":
		jujuCmd = application.NewDiffBundleCommand()
	case "show-application":
		jujuCmd = application.NewShowApplicationCommand()
	case "show-unit":
		jujuCmd = application.NewShowUnitCommand()
	case "refresh":
		jujuCmd = application.NewRefreshCommand()
	case "bind":
		jujuCmd = application.NewBindCommand()
	case "scale-application":
		jujuCmd = application.NewScaleApplicationCommand()
	case "trust":
		jujuCmd = application.NewTrustCommand()
	case "add-relation":
		jujuCmd = application.NewAddRelationCommand()
	case "remove-relation":
		jujuCmd = application.NewRemoveRelationCommand()
	case "remove-application":
		jujuCmd = application.NewRemoveApplicationCommand()
	case "remove-unit":
		jujuCmd = application.NewRemoveUnitCommand()
	case "remove-saas":
		jujuCmd = application.NewRemoveSaasCommand()
	case "consume":
		jujuCmd = application.NewConsumeCommand()
	case "suspend-relation":
		jujuCmd = application.NewSuspendRelationCommand()
	case "resume-relation":
		jujuCmd = application.NewResumeRelationCommand()
	case "resolved":
		jujuCmd = application.NewResolvedCommand()

	// Machine commands
	case "add-machine":
		jujuCmd = machine.NewAddCommand()
	case "remove-machine":
		jujuCmd = machine.NewRemoveCommand()
	case "machines":
		jujuCmd = machine.NewListMachinesCommand()
	case "show-machine":
		jujuCmd = machine.NewShowMachineCommand()

	// Model commands
	case "model-config":
		jujuCmd = model.NewConfigCommand()
	case "model-defaults":
		jujuCmd = model.NewDefaultsCommand()
	case "retry-provisioning":
		jujuCmd = model.NewRetryProvisioningCommand()
	case "destroy-model":
		jujuCmd = model.NewDestroyCommand()
	case "grant":
		jujuCmd = model.NewGrantCommand()
	case "revoke":
		jujuCmd = model.NewRevokeCommand()
	case "show-model":
		jujuCmd = model.NewShowCommand()
	case "model-credential":
		jujuCmd = model.NewModelCredentialCommand()
	case "export-bundle":
		jujuCmd = model.NewExportBundleCommand()
	case "migrate":
		var err error
		jujuCmd, err = commands.NewCommandByName("migrate")
		if err != nil {
			return nil, err
		}

	// Commands from main commands package
	case "bootstrap":
		var err error
		jujuCmd, err = commands.NewCommandByName("bootstrap")
		if err != nil {
			return nil, err
		}
	case "switch":
		var err error
		jujuCmd, err = commands.NewCommandByName("switch")
		if err != nil {
			return nil, err
		}
	case "version":
		var err error
		jujuCmd, err = commands.NewCommandByName("version")
		if err != nil {
			return nil, err
		}
	case "sync-agent-binary":
		var err error
		jujuCmd, err = commands.NewCommandByName("sync-agent-binary")
		if err != nil {
			return nil, err
		}
	case "upgrade-model":
		var err error
		jujuCmd, err = commands.NewCommandByName("upgrade-model")
		if err != nil {
			return nil, err
		}
	case "upgrade-controller":
		var err error
		jujuCmd, err = commands.NewCommandByName("upgrade-controller")
		if err != nil {
			return nil, err
		}
	case "help-hooks":
		var err error
		jujuCmd, err = commands.NewCommandByName("help-hooks")
		if err != nil {
			return nil, err
		}
	case "help-actions":
		var err error
		jujuCmd, err = commands.NewCommandByName("help-actions")
		if err != nil {
			return nil, err
		}
	case "debug-log":
		var err error
		jujuCmd, err = commands.NewCommandByName("debug-log")
		if err != nil {
			return nil, err
		}

	// Controller commands
	case "add-model":
		jujuCmd = controller.NewAddModelCommand()
	case "destroy-controller":
		jujuCmd = controller.NewDestroyCommand()
	case "models":
		jujuCmd = controller.NewListModelsCommand()
	case "kill-controller":
		jujuCmd = controller.NewKillCommand()
	case "controllers":
		jujuCmd = controller.NewListControllersCommand()
	case "register":
		jujuCmd = controller.NewRegisterCommand()
	case "unregister":
		jujuCmd = controller.NewUnregisterCommand(jujuclient.NewFileClientStore())
	case "enable-destroy-controller":
		jujuCmd = controller.NewEnableDestroyControllerCommand()
	case "show-controller":
		jujuCmd = controller.NewShowControllerCommand()
	case "controller-config":
		jujuCmd = controller.NewConfigCommand()

	// Action commands
	case "actions":
		jujuCmd = action.NewListCommand()
	case "show-action":
		jujuCmd = action.NewShowCommand()
	case "cancel-action":
		jujuCmd = action.NewCancelCommand()
	case "run":
		jujuCmd = action.NewRunCommand()
	case "operations":
		jujuCmd = action.NewListOperationsCommand()
	case "show-operation":
		jujuCmd = action.NewShowOperationCommand()
	case "show-task":
		jujuCmd = action.NewShowTaskCommand()

	// Cloud and credential commands
	case "update-cloud":
		jujuCmd = cloud.NewUpdateCloudCommand(&cloudToCommandAdaptor{})
	case "update-public-clouds":
		jujuCmd = cloud.NewUpdatePublicCloudsCommand()
	case "clouds":
		jujuCmd = cloud.NewListCloudsCommand()
	case "regions":
		jujuCmd = cloud.NewListRegionsCommand()
	case "show-cloud":
		jujuCmd = cloud.NewShowCloudCommand()
	case "add-cloud":
		jujuCmd = cloud.NewAddCloudCommand(&cloudToCommandAdaptor{})
	case "remove-cloud":
		jujuCmd = cloud.NewRemoveCloudCommand()
	case "credentials":
		jujuCmd = cloud.NewListCredentialsCommand()
	case "detect-credentials":
		jujuCmd = cloud.NewDetectCredentialsCommand()
	case "set-default-region":
		jujuCmd = cloud.NewSetDefaultRegionCommand()
	case "set-default-credential":
		jujuCmd = cloud.NewSetDefaultCredentialCommand()
	case "add-credential":
		jujuCmd = cloud.NewAddCredentialCommand()
	case "remove-credential":
		jujuCmd = cloud.NewRemoveCredentialCommand()
	case "update-credential":
		jujuCmd = cloud.NewUpdateCredentialCommand()
	case "show-credential":
		jujuCmd = cloud.NewShowCredentialCommand()
	case "grant-cloud":
		jujuCmd = model.NewGrantCloudCommand()
	case "revoke-cloud":
		jujuCmd = model.NewRevokeCloudCommand()

	// User commands
	case "add-user":
		jujuCmd = user.NewAddCommand()
	case "change-password":
		jujuCmd = user.NewChangePasswordCommand()
	case "show-user":
		jujuCmd = user.NewShowUserCommand()
	case "users":
		jujuCmd = user.NewListCommand()
	case "enable-user":
		jujuCmd = user.NewEnableCommand()
	case "disable-user":
		jujuCmd = user.NewDisableCommand()
	case "login":
		jujuCmd = user.NewLoginCommand()
	case "logout":
		jujuCmd = user.NewLogoutCommand()
	case "remove-user":
		jujuCmd = user.NewRemoveCommand()
	case "whoami":
		jujuCmd = user.NewWhoAmICommand()

	// Storage commands
	case "add-storage":
		jujuCmd = storage.NewAddCommand()
	case "storage":
		jujuCmd = storage.NewListCommand()
	case "create-storage-pool":
		jujuCmd = storage.NewPoolCreateCommand()
	case "storage-pools":
		jujuCmd = storage.NewPoolListCommand()
	case "remove-storage-pool":
		jujuCmd = storage.NewPoolRemoveCommand()
	case "update-storage-pool":
		jujuCmd = storage.NewPoolUpdateCommand()
	case "show-storage":
		jujuCmd = storage.NewShowCommand()
	case "remove-storage":
		jujuCmd = storage.NewRemoveStorageCommandWithAPI()
	case "detach-storage":
		jujuCmd = storage.NewDetachStorageCommandWithAPI()
	case "attach-storage":
		jujuCmd = storage.NewAttachStorageCommandWithAPI()
	case "import-filesystem":
		jujuCmd = storage.NewImportFilesystemCommand(storage.NewStorageImporter, nil)

	// Space commands
	case "add-space":
		jujuCmd = space.NewAddCommand()
	case "spaces":
		jujuCmd = space.NewListCommand()
	case "move-to-space":
		jujuCmd = space.NewMoveCommand()
	case "reload-spaces":
		jujuCmd = space.NewReloadCommand()
	case "show-space":
		jujuCmd = space.NewShowSpaceCommand()
	case "remove-space":
		jujuCmd = space.NewRemoveCommand()
	case "rename-space":
		jujuCmd = space.NewRenameCommand()

	// Subnet commands
	case "subnets":
		jujuCmd = subnet.NewListCommand()

	// SSH and debugging commands
	case "exec":
		jujuCmd = action.NewExecCommand(nil)
	case "scp":
		jujuCmd = ssh.NewSCPCommand(nil, ssh.DefaultSSHRetryStrategy, ssh.DefaultSSHPublicKeyRetryStrategy)
	case "ssh":
		jujuCmd = ssh.NewSSHCommand(nil, nil, ssh.DefaultSSHRetryStrategy, ssh.DefaultSSHPublicKeyRetryStrategy)
	case "debug-hooks":
		jujuCmd = ssh.NewDebugHooksCommand(nil, ssh.DefaultSSHRetryStrategy, ssh.DefaultSSHPublicKeyRetryStrategy)
	case "debug-code":
		jujuCmd = ssh.NewDebugCodeCommand(nil, ssh.DefaultSSHRetryStrategy, ssh.DefaultSSHPublicKeyRetryStrategy)

	// SSH keys commands
	case "add-ssh-key":
		jujuCmd = commands.NewAddKeysCommand()
	case "remove-ssh-key":
		jujuCmd = commands.NewRemoveKeysCommand()
	case "import-ssh-key":
		jujuCmd = commands.NewImportKeysCommand()
	case "ssh-keys":
		jujuCmd = commands.NewListKeysCommand()

	// Backup commands
	case "create-backup":
		jujuCmd = backups.NewCreateCommand()
	case "download-backup":
		jujuCmd = backups.NewDownloadCommand()

	// Block commands
	case "disable-command":
		jujuCmd = block.NewDisableCommand()
	case "disabled-commands":
		jujuCmd = block.NewListCommand()
	case "enable-command":
		jujuCmd = block.NewEnableCommand()

	// Firewall commands
	case "set-firewall-rule":
		jujuCmd = firewall.NewSetFirewallRuleCommand()
	case "list-firewall-rules":
		jujuCmd = firewall.NewListFirewallRulesCommand()

	// Cross model commands
	case "offer":
		jujuCmd = crossmodel.NewOfferCommand()
	case "remove-offer":
		jujuCmd = crossmodel.NewRemoveOfferCommand()
	case "show-offered-endpoint":
		jujuCmd = crossmodel.NewShowOfferedEndpointCommand()
	case "list-endpoints":
		jujuCmd = crossmodel.NewListEndpointsCommand()
	case "find-endpoints":
		jujuCmd = crossmodel.NewFindEndpointsCommand()

	// CAAS commands
	case "add-k8s":
		jujuCmd = caas.NewAddCAASCommand(&cloudToCommandAdaptor{})
	case "update-k8s":
		jujuCmd = caas.NewUpdateCAASCommand(&cloudToCommandAdaptor{})
	case "remove-k8s":
		jujuCmd = caas.NewRemoveCAASCommand(&cloudToCommandAdaptor{})

	// Dashboard commands
	case "dashboard":
		jujuCmd = dashboard.NewDashboardCommand()

	// Resource commands
	case "attach-resource":
		jujuCmd = resource.NewUploadCommand()
	case "resources":
		jujuCmd = resource.NewListCommand()
	case "charm-resources":
		jujuCmd = resource.NewCharmResourcesCommand()

	// CharmHub commands
	case "info":
		jujuCmd = charmhub.NewInfoCommand()
	case "find":
		jujuCmd = charmhub.NewFindCommand()
	case "download":
		jujuCmd = charmhub.NewDownloadCommand()

	// Secrets commands
	case "secrets":
		jujuCmd = secrets.NewListSecretsCommand()
	case "show-secret":
		jujuCmd = secrets.NewShowSecretsCommand()
	case "add-secret":
		jujuCmd = secrets.NewAddSecretCommand()
	case "update-secret":
		jujuCmd = secrets.NewUpdateSecretCommand()
	case "remove-secret":
		jujuCmd = secrets.NewRemoveSecretCommand()
	case "grant-secret":
		jujuCmd = secrets.NewGrantSecretCommand()
	case "revoke-secret":
		jujuCmd = secrets.NewRevokeSecretCommand()

	// Secret backends commands
	case "secret-backends":
		jujuCmd = secretbackends.NewListSecretBackendsCommand()
	case "add-secret-backend":
		jujuCmd = secretbackends.NewAddSecretBackendCommand()
	case "update-secret-backend":
		jujuCmd = secretbackends.NewUpdateSecretBackendCommand()
	case "remove-secret-backend":
		jujuCmd = secretbackends.NewRemoveSecretBackendCommand()
	case "show-secret-backend":
		jujuCmd = secretbackends.NewShowSecretBackendCommand()
	// case "model-secret-backend":
	//	jujuCmd = secretbackends.NewModelSecretBackendCommand()

	// Wait for commands
	case "wait-for":
		jujuCmd = waitfor.NewWaitForCommand()
	
	// Missing commands that should be added
	case "upgrade-machine":
		jujuCmd = machine.NewUpgradeMachineCommand()
	case "enable-ha":
		var err error
		jujuCmd, err = commands.NewCommandByName("enable-ha")
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown command: %s", name)
	}
	
	if jujuCmd == nil {
		return nil, fmt.Errorf("failed to create command: %s", name)
	}

	return &command{
		cmd:  jujuCmd,
		info: jujuCmd.Info(),
		t:    time.Now(),
	}, nil
}
