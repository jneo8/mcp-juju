package jujuadapter

// JujuCommandID represents a unique command identifier
type JujuCommandID string


// Define all command IDs as constants
const (
	// Reporting commands
	CmdVersion       JujuCommandID = "version"
	CmdStatus        JujuCommandID = "status"
	CmdStatusHistory JujuCommandID = "status-history"
	CmdSwitch        JujuCommandID = "switch"

	// Creation commands
	CmdBootstrap   JujuCommandID = "bootstrap"
	CmdAddRelation JujuCommandID = "add-relation"

	// Cross model relations commands
	CmdOffer               JujuCommandID = "offer"
	CmdRemoveOffer         JujuCommandID = "remove-offer"
	CmdShowOfferedEndpoint JujuCommandID = "show-offered-endpoint"
	CmdListEndpoints       JujuCommandID = "list-endpoints"
	CmdFindEndpoints       JujuCommandID = "find-endpoints"
	CmdConsume             JujuCommandID = "consume"
	CmdSuspendRelation     JujuCommandID = "suspend-relation"
	CmdResumeRelation      JujuCommandID = "resume-relation"

	// Firewall rule commands
	CmdSetFirewallRule   JujuCommandID = "set-firewall-rule"
	CmdListFirewallRules JujuCommandID = "list-firewall-rules"

	// Destruction commands
	CmdRemoveRelation    JujuCommandID = "remove-relation"
	CmdRemoveApplication JujuCommandID = "remove-application"
	CmdRemoveUnit        JujuCommandID = "remove-unit"
	CmdRemoveSaas        JujuCommandID = "remove-saas"

	// Error resolution and debugging commands
	CmdExec       JujuCommandID = "exec"
	CmdScp        JujuCommandID = "scp"
	CmdSsh        JujuCommandID = "ssh"
	CmdResolved   JujuCommandID = "resolved"
	CmdDebugLog   JujuCommandID = "debug-log"
	CmdDebugHooks JujuCommandID = "debug-hooks"
	CmdDebugCode  JujuCommandID = "debug-code"

	// Configuration commands
	CmdGetConstraints    JujuCommandID = "get-constraints"
	CmdSetConstraints    JujuCommandID = "set-constraints"
	CmdSyncAgentBinary   JujuCommandID = "sync-agent-binary"
	CmdUpgradeModel      JujuCommandID = "upgrade-model"
	CmdUpgradeController JujuCommandID = "upgrade-controller"
	CmdRefresh           JujuCommandID = "refresh"
	CmdBind              JujuCommandID = "bind"

	// Charm tool commands
	CmdHelpHooks   JujuCommandID = "help-hooks"
	CmdHelpActions JujuCommandID = "help-actions"

	// Manage backups
	CmdCreateBackup   JujuCommandID = "create-backup"
	CmdDownloadBackup JujuCommandID = "download-backup"

	// Manage authorized ssh keys
	CmdAddSshKey    JujuCommandID = "add-ssh-key"
	CmdRemoveSshKey JujuCommandID = "remove-ssh-key"
	CmdImportSshKey JujuCommandID = "import-ssh-key"
	CmdSshKeys      JujuCommandID = "ssh-keys"

	// Manage users and access
	CmdAddUser        JujuCommandID = "add-user"
	CmdChangePassword JujuCommandID = "change-password"
	CmdShowUser       JujuCommandID = "show-user"
	CmdUsers          JujuCommandID = "users"
	CmdEnableUser     JujuCommandID = "enable-user"
	CmdDisableUser    JujuCommandID = "disable-user"
	CmdLogin          JujuCommandID = "login"
	CmdLogout         JujuCommandID = "logout"
	CmdRemoveUser     JujuCommandID = "remove-user"
	CmdWhoami         JujuCommandID = "whoami"

	// Manage machines
	CmdAddMachine     JujuCommandID = "add-machine"
	CmdRemoveMachine  JujuCommandID = "remove-machine"
	CmdMachines       JujuCommandID = "machines"
	CmdShowMachine    JujuCommandID = "show-machine"
	CmdUpgradeMachine JujuCommandID = "upgrade-machine"

	// Manage model
	CmdModelConfig       JujuCommandID = "model-config"
	CmdModelDefaults     JujuCommandID = "model-defaults"
	CmdRetryProvisioning JujuCommandID = "retry-provisioning"
	CmdDestroyModel      JujuCommandID = "destroy-model"
	CmdGrant             JujuCommandID = "grant"
	CmdRevoke            JujuCommandID = "revoke"
	CmdShowModel         JujuCommandID = "show-model"
	CmdModelCredential   JujuCommandID = "model-credential"
	CmdMigrate           JujuCommandID = "migrate"
	CmdExportBundle      JujuCommandID = "export-bundle"

	// Manage and control actions
	CmdActions       JujuCommandID = "actions"
	CmdShowAction    JujuCommandID = "show-action"
	CmdCancelAction  JujuCommandID = "cancel-action"
	CmdRun           JujuCommandID = "run"
	CmdOperations    JujuCommandID = "operations"
	CmdShowOperation JujuCommandID = "show-operation"
	CmdShowTask      JujuCommandID = "show-task"

	// Manage controller availability
	CmdEnableHa JujuCommandID = "enable-ha"

	// Manage and control applications
	CmdAddUnit         JujuCommandID = "add-unit"
	CmdConfig          JujuCommandID = "config"
	CmdDeploy          JujuCommandID = "deploy"
	CmdExpose          JujuCommandID = "expose"
	CmdUnexpose        JujuCommandID = "unexpose"
	CmdDiffBundle      JujuCommandID = "diff-bundle"
	CmdShowApplication JujuCommandID = "show-application"
	CmdShowUnit        JujuCommandID = "show-unit"

	// Operation protection commands
	CmdDisableCommand   JujuCommandID = "disable-command"
	CmdDisabledCommands JujuCommandID = "disabled-commands"
	CmdEnableCommand    JujuCommandID = "enable-command"

	// Manage storage
	CmdAddStorage        JujuCommandID = "add-storage"
	CmdStorage           JujuCommandID = "storage"
	CmdCreateStoragePool JujuCommandID = "create-storage-pool"
	CmdStoragePools      JujuCommandID = "storage-pools"
	CmdRemoveStoragePool JujuCommandID = "remove-storage-pool"
	CmdUpdateStoragePool JujuCommandID = "update-storage-pool"
	CmdShowStorage       JujuCommandID = "show-storage"
	CmdRemoveStorage     JujuCommandID = "remove-storage"
	CmdDetachStorage     JujuCommandID = "detach-storage"
	CmdAttachStorage     JujuCommandID = "attach-storage"
	CmdImportFilesystem  JujuCommandID = "import-filesystem"

	// Manage spaces
	CmdAddSpace     JujuCommandID = "add-space"
	CmdSpaces       JujuCommandID = "spaces"
	CmdMoveToSpace  JujuCommandID = "move-to-space"
	CmdReloadSpaces JujuCommandID = "reload-spaces"
	CmdShowSpace    JujuCommandID = "show-space"
	CmdRemoveSpace  JujuCommandID = "remove-space"
	CmdRenameSpace  JujuCommandID = "rename-space"

	// Manage subnets
	CmdSubnets JujuCommandID = "subnets"

	// Manage controllers
	CmdAddModel                JujuCommandID = "add-model"
	CmdDestroyController       JujuCommandID = "destroy-controller"
	CmdModels                  JujuCommandID = "models"
	CmdKillController          JujuCommandID = "kill-controller"
	CmdControllers             JujuCommandID = "controllers"
	CmdRegister                JujuCommandID = "register"
	CmdUnregister              JujuCommandID = "unregister"
	CmdEnableDestroyController JujuCommandID = "enable-destroy-controller"
	CmdShowController          JujuCommandID = "show-controller"
	CmdControllerConfig        JujuCommandID = "controller-config"

	// Manage clouds and credentials
	CmdUpdateCloud          JujuCommandID = "update-cloud"
	CmdUpdatePublicClouds   JujuCommandID = "update-public-clouds"
	CmdClouds               JujuCommandID = "clouds"
	CmdRegions              JujuCommandID = "regions"
	CmdShowCloud            JujuCommandID = "show-cloud"
	CmdAddCloud             JujuCommandID = "add-cloud"
	CmdRemoveCloud          JujuCommandID = "remove-cloud"
	CmdCredentials          JujuCommandID = "credentials"
	CmdDetectCredentials    JujuCommandID = "detect-credentials"
	CmdSetDefaultRegion     JujuCommandID = "set-default-region"
	CmdSetDefaultCredential JujuCommandID = "set-default-credential"
	CmdAddCredential        JujuCommandID = "add-credential"
	CmdRemoveCredential     JujuCommandID = "remove-credential"
	CmdUpdateCredential     JujuCommandID = "update-credential"
	CmdShowCredential       JujuCommandID = "show-credential"
	CmdGrantCloud           JujuCommandID = "grant-cloud"
	CmdRevokeCloud          JujuCommandID = "revoke-cloud"

	// CAAS commands
	CmdAddK8s           JujuCommandID = "add-k8s"
	CmdUpdateK8s        JujuCommandID = "update-k8s"
	CmdRemoveK8s        JujuCommandID = "remove-k8s"
	CmdScaleApplication JujuCommandID = "scale-application"

	// Manage Application Credential Access
	CmdTrust JujuCommandID = "trust"

	// Juju Dashboard commands
	CmdDashboard JujuCommandID = "dashboard"

	// Resource commands
	CmdAttachResource JujuCommandID = "attach-resource"
	CmdResources      JujuCommandID = "resources"
	CmdCharmResources JujuCommandID = "charm-resources"

	// CharmHub related commands
	CmdInfo     JujuCommandID = "info"
	CmdFind     JujuCommandID = "find"
	CmdDownload JujuCommandID = "download"

	// Secrets
	CmdSecrets      JujuCommandID = "secrets"
	CmdShowSecret   JujuCommandID = "show-secret"
	CmdAddSecret    JujuCommandID = "add-secret"
	CmdUpdateSecret JujuCommandID = "update-secret"
	CmdRemoveSecret JujuCommandID = "remove-secret"
	CmdGrantSecret  JujuCommandID = "grant-secret"
	CmdRevokeSecret JujuCommandID = "revoke-secret"

	// Secret backends
	CmdSecretBackends      JujuCommandID = "secret-backends"
	CmdAddSecretBackend    JujuCommandID = "add-secret-backend"
	CmdUpdateSecretBackend JujuCommandID = "update-secret-backend"
	CmdRemoveSecretBackend JujuCommandID = "remove-secret-backend"
	CmdShowSecretBackend   JujuCommandID = "show-secret-backend"

	// Payload commands
	CmdWaitFor JujuCommandID = "wait-for"
)



// GetAllCommandIDs returns all available command IDs in order
func GetAllCommandIDs() []JujuCommandID {
	return []JujuCommandID{
		// From registerCommands in exact order as juju/cmd/juju/commands/main.go
		CmdVersion,

		// Creation commands.
		CmdBootstrap,
		CmdAddRelation,

		// Cross model relations commands.
		CmdOffer,
		CmdRemoveOffer,
		CmdShowOfferedEndpoint,
		CmdListEndpoints,
		CmdFindEndpoints,
		CmdConsume,
		CmdSuspendRelation,
		CmdResumeRelation,

		// Firewall rule commands.
		CmdSetFirewallRule,
		CmdListFirewallRules,

		// Destruction commands.
		CmdRemoveRelation,
		CmdRemoveApplication,
		CmdRemoveUnit,
		CmdRemoveSaas,

		// Reporting commands.
		CmdStatus,
		CmdSwitch,
		CmdStatusHistory,

		// Error resolution and debugging commands.
		CmdExec,
		CmdScp,
		CmdSsh,
		CmdResolved,
		CmdDebugLog,
		CmdDebugHooks,
		CmdDebugCode,

		// Configuration commands.
		CmdGetConstraints,
		CmdSetConstraints,
		CmdSyncAgentBinary,
		CmdUpgradeModel,
		CmdUpgradeController,
		CmdRefresh,
		CmdBind,

		// Charm tool commands.
		CmdHelpHooks,
		CmdHelpActions,

		// Manage backups.
		CmdCreateBackup,
		CmdDownloadBackup,

		// Manage authorized ssh keys.
		CmdAddSshKey,
		CmdRemoveSshKey,
		CmdImportSshKey,
		CmdSshKeys,

		// Manage users and access
		CmdAddUser,
		CmdChangePassword,
		CmdShowUser,
		CmdUsers,
		CmdEnableUser,
		CmdDisableUser,
		CmdLogin,
		CmdLogout,
		CmdRemoveUser,
		CmdWhoami,

		// Manage machines
		CmdAddMachine,
		CmdRemoveMachine,
		CmdMachines,
		CmdShowMachine,
		CmdUpgradeMachine,

		// Manage model
		CmdModelConfig,
		CmdModelDefaults,
		CmdRetryProvisioning,
		CmdDestroyModel,
		CmdGrant,
		CmdRevoke,
		CmdShowModel,
		CmdModelCredential,

		CmdMigrate,
		CmdExportBundle,

		// Manage and control actions
		CmdActions,
		CmdShowAction,
		CmdCancelAction,
		CmdRun,
		CmdOperations,
		CmdShowOperation,
		CmdShowTask,

		// Manage controller availability
		CmdEnableHa,

		// Manage and control applications
		CmdAddUnit,
		CmdConfig,
		CmdDeploy,
		CmdExpose,
		CmdUnexpose,
		CmdDiffBundle,
		CmdShowApplication,
		CmdShowUnit,

		// Operation protection commands
		CmdDisableCommand,
		CmdDisabledCommands,
		CmdEnableCommand,

		// Manage storage
		CmdAddStorage,
		CmdStorage,
		CmdCreateStoragePool,
		CmdStoragePools,
		CmdRemoveStoragePool,
		CmdUpdateStoragePool,
		CmdShowStorage,
		CmdRemoveStorage,
		CmdDetachStorage,
		CmdAttachStorage,
		CmdImportFilesystem,

		// Manage spaces
		CmdAddSpace,
		CmdSpaces,
		CmdMoveToSpace,
		CmdReloadSpaces,
		CmdShowSpace,
		CmdRemoveSpace,
		CmdRenameSpace,

		// Manage subnets
		CmdSubnets,

		// Manage controllers
		CmdAddModel,
		CmdDestroyController,
		CmdModels,
		CmdKillController,
		CmdControllers,
		CmdRegister,
		CmdUnregister,
		CmdEnableDestroyController,
		CmdShowController,
		CmdControllerConfig,

		// Manage clouds and credentials
		CmdUpdateCloud,
		CmdUpdatePublicClouds,
		CmdClouds,
		CmdRegions,
		CmdShowCloud,
		CmdAddCloud,
		CmdRemoveCloud,
		CmdCredentials,
		CmdDetectCredentials,
		CmdSetDefaultRegion,
		CmdSetDefaultCredential,
		CmdAddCredential,
		CmdRemoveCredential,
		CmdUpdateCredential,
		CmdShowCredential,
		CmdGrantCloud,
		CmdRevokeCloud,

		// CAAS commands
		CmdAddK8s,
		CmdUpdateK8s,
		CmdRemoveK8s,
		CmdScaleApplication,

		// Manage Application Credential Access
		CmdTrust,

		// Juju Dashboard commands.
		CmdDashboard,

		// Resource commands.
		CmdAttachResource,
		CmdResources,
		CmdCharmResources,

		// CharmHub related commands
		CmdInfo,
		CmdFind,
		CmdDownload,

		// Secrets.
		CmdSecrets,
		CmdShowSecret,
		CmdAddSecret,
		CmdUpdateSecret,
		CmdRemoveSecret,
		CmdGrantSecret,
		CmdRevokeSecret,

		// Secret backends.
		CmdSecretBackends,
		CmdAddSecretBackend,
		CmdUpdateSecretBackend,
		CmdRemoveSecretBackend,
		CmdShowSecretBackend,

		// Payload commands.
		CmdWaitFor,
	}
}
