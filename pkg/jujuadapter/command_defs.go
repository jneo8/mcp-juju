package jujuadapter

// JujuCommandID represents a unique command identifier
type JujuCommandID string

// CommandDefinition contains all metadata for a command
type CommandDefinition struct {
	ID                        JujuCommandID
	StdioMcpServerDisableArgs map[string]bool // Arguments disabled for stdio MCP server
	HTTPMcpServerDisableArgs  map[string]bool // Arguments disabled for HTTP MCP server
}

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

// commandRegistry maps command IDs to their definitions
var commandRegistry = map[JujuCommandID]CommandDefinition{
	// Commands with disabled arguments
	CmdBootstrap: {
		ID: CmdBootstrap,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"build-agent":      true, // Local agent building
			"metadata-source":  true, // Local file paths
			"config":           true, // Config file operations
			"model-default":    true, // Model default file operations
		},
	},
	CmdDestroyController: {
		ID: CmdDestroyController,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdDestroyModel: {
		ID: CmdDestroyModel,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdDebugLog: {
		ID: CmdDebugLog,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRemoveApplication: {
		ID: CmdRemoveApplication,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRemoveUnit: {
		ID: CmdRemoveUnit,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRemoveMachine: {
		ID: CmdRemoveMachine,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},

	// Commands without disabled args - just ID field needed
	CmdVersion: {
		ID: CmdVersion,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdStatus: {
		ID: CmdStatus,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdStatusHistory: {
		ID: CmdStatusHistory,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdSwitch: {
		ID: CmdSwitch,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdAddRelation: {
		ID: CmdAddRelation,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdOffer: {
		ID: CmdOffer,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRemoveOffer: {
		ID: CmdRemoveOffer,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdShowOfferedEndpoint: {
		ID: CmdShowOfferedEndpoint,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdListEndpoints: {
		ID: CmdListEndpoints,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdFindEndpoints: {
		ID: CmdFindEndpoints,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdConsume: {
		ID: CmdConsume,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdSuspendRelation: {
		ID: CmdSuspendRelation,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdResumeRelation: {
		ID: CmdResumeRelation,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdSetFirewallRule: {
		ID: CmdSetFirewallRule,
	},
	CmdListFirewallRules: {
		ID: CmdListFirewallRules,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdRemoveRelation: {
		ID: CmdRemoveRelation,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRemoveSaas: {
		ID: CmdRemoveSaas,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdExec: {
		ID: CmdExec,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdScp: {
		ID: CmdScp,
	},
	CmdSsh: {
		ID: CmdSsh,
	},
	CmdResolved: {
		ID: CmdResolved,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdDebugHooks: {
		ID: CmdDebugHooks,
	},
	CmdDebugCode: {
		ID: CmdDebugCode,
	},
	CmdGetConstraints: {
		ID: CmdGetConstraints,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdSetConstraints: {
		ID: CmdSetConstraints,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdSyncAgentBinary: {
		ID: CmdSyncAgentBinary,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"source":           true, // Local source directory
			"local-dir":        true, // Local destination directory
		},
	},
	CmdUpgradeModel: {
		ID: CmdUpgradeModel,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdUpgradeController: {
		ID: CmdUpgradeController,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRefresh: {
		ID: CmdRefresh,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"path":             true, // Local charm path
			"config":           true, // Config file operations
			"resource":         true, // Resource file operations
		},
	},
	CmdBind: {
		ID: CmdBind,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdHelpHooks: {
		ID: CmdHelpHooks,
	},
	CmdHelpActions: {
		ID: CmdHelpActions,
	},
	CmdCreateBackup: {
		ID: CmdCreateBackup,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"filename":         true, // Local file operations
		},
	},
	CmdDownloadBackup: {
		ID: CmdDownloadBackup,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"filename":         true, // Local file operations
		},
	},
	CmdAddSshKey: {
		ID: CmdAddSshKey,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRemoveSshKey: {
		ID: CmdRemoveSshKey,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdImportSshKey: {
		ID: CmdImportSshKey,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdSshKeys: {
		ID: CmdSshKeys,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdAddUser: {
		ID: CmdAddUser,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdChangePassword: {
		ID: CmdChangePassword,
	},
	CmdShowUser: {
		ID: CmdShowUser,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdUsers: {
		ID: CmdUsers,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdEnableUser: {
		ID: CmdEnableUser,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdDisableUser: {
		ID: CmdDisableUser,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdLogin: {
		ID: CmdLogin,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdLogout: {
		ID: CmdLogout,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRemoveUser: {
		ID: CmdRemoveUser,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdWhoami: {
		ID: CmdWhoami,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdAddMachine: {
		ID: CmdAddMachine,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"private-key":      true, // SSH private key file
			"public-key":       true, // SSH public key file
		},
	},
	CmdMachines: {
		ID: CmdMachines,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdShowMachine: {
		ID: CmdShowMachine,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdUpgradeMachine: {
		ID: CmdUpgradeMachine,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdModelConfig: {
		ID: CmdModelConfig,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"file":             true, // Config file operations
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdModelDefaults: {
		ID: CmdModelDefaults,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"file":             true, // Config file operations
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdRetryProvisioning: {
		ID: CmdRetryProvisioning,
	},
	CmdGrant: {
		ID: CmdGrant,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRevoke: {
		ID: CmdRevoke,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdShowModel: {
		ID: CmdShowModel,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdModelCredential: {
		ID: CmdModelCredential,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdMigrate: {
		ID: CmdMigrate,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdExportBundle: {
		ID: CmdExportBundle,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"filename":         true, // Local file operations
		},
	},
	CmdActions: {
		ID: CmdActions,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdShowAction: {
		ID: CmdShowAction,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdCancelAction: {
		ID: CmdCancelAction,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdRun: {
		ID: CmdRun,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"params":           true, // Parameter file operations
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdOperations: {
		ID: CmdOperations,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdShowOperation: {
		ID: CmdShowOperation,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdShowTask: {
		ID: CmdShowTask,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdEnableHa: {
		ID: CmdEnableHa,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdAddUnit: {
		ID: CmdAddUnit,
	},
	CmdConfig: {
		ID: CmdConfig,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"file":             true, // Config file operations
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdDeploy: {
		ID: CmdDeploy,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"config":           true, // Config file operations
			"resource":         true, // Resource file operations
			"overlay":          true, // Bundle overlay files
		},
	},
	CmdExpose: {
		ID: CmdExpose,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdUnexpose: {
		ID: CmdUnexpose,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdDiffBundle: {
		ID: CmdDiffBundle,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"overlay":          true, // Bundle overlay files
		},
	},
	CmdShowApplication: {
		ID: CmdShowApplication,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdShowUnit: {
		ID: CmdShowUnit,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdDisableCommand: {
		ID: CmdDisableCommand,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdDisabledCommands: {
		ID: CmdDisabledCommands,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdEnableCommand: {
		ID: CmdEnableCommand,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdAddStorage: {
		ID: CmdAddStorage,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdStorage: {
		ID: CmdStorage,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdCreateStoragePool: {
		ID: CmdCreateStoragePool,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdStoragePools: {
		ID: CmdStoragePools,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdRemoveStoragePool: {
		ID: CmdRemoveStoragePool,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdUpdateStoragePool: {
		ID: CmdUpdateStoragePool,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdShowStorage: {
		ID: CmdShowStorage,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdRemoveStorage: {
		ID: CmdRemoveStorage,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdDetachStorage: {
		ID: CmdDetachStorage,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdAttachStorage: {
		ID: CmdAttachStorage,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdImportFilesystem: {
		ID: CmdImportFilesystem,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdAddSpace: {
		ID: CmdAddSpace,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdSpaces: {
		ID: CmdSpaces,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdMoveToSpace: {
		ID: CmdMoveToSpace,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdReloadSpaces: {
		ID: CmdReloadSpaces,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdShowSpace: {
		ID: CmdShowSpace,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdRemoveSpace: {
		ID: CmdRemoveSpace,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRenameSpace: {
		ID: CmdRenameSpace,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdSubnets: {
		ID: CmdSubnets,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdAddModel: {
		ID: CmdAddModel,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"config":           true, // Config file operations
		},
	},
	CmdModels: {
		ID: CmdModels,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdKillController: {
		ID: CmdKillController,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdControllers: {
		ID: CmdControllers,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdRegister: {
		ID: CmdRegister,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdUnregister: {
		ID: CmdUnregister,
	},
	CmdEnableDestroyController: {
		ID: CmdEnableDestroyController,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdShowController: {
		ID: CmdShowController,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdControllerConfig: {
		ID: CmdControllerConfig,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"file":             true, // Config file operations
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdUpdateCloud: {
		ID: CmdUpdateCloud,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"f":                true, // Cloud definition file
		},
	},
	CmdUpdatePublicClouds: {
		ID: CmdUpdatePublicClouds,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdClouds: {
		ID: CmdClouds,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdRegions: {
		ID: CmdRegions,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdShowCloud: {
		ID: CmdShowCloud,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdAddCloud: {
		ID: CmdAddCloud,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"f":                true, // Cloud definition file
			"file":             true, // Cloud definition file
		},
	},
	CmdRemoveCloud: {
		ID: CmdRemoveCloud,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdCredentials: {
		ID: CmdCredentials,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdDetectCredentials: {
		ID: CmdDetectCredentials,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdSetDefaultRegion: {
		ID: CmdSetDefaultRegion,
	},
	CmdSetDefaultCredential: {
		ID: CmdSetDefaultCredential,
	},
	CmdAddCredential: {
		ID: CmdAddCredential,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"f":                true, // Credentials file
			"file":             true, // Credentials file
		},
	},
	CmdRemoveCredential: {
		ID: CmdRemoveCredential,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdUpdateCredential: {
		ID: CmdUpdateCredential,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"f":                true, // Credentials file
			"file":             true, // Credentials file
		},
	},
	CmdShowCredential: {
		ID: CmdShowCredential,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdGrantCloud: {
		ID: CmdGrantCloud,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRevokeCloud: {
		ID: CmdRevokeCloud,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdAddK8s: {
		ID: CmdAddK8s,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdUpdateK8s: {
		ID: CmdUpdateK8s,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"f":                true, // K8s definition file
		},
	},
	CmdRemoveK8s: {
		ID: CmdRemoveK8s,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdScaleApplication: {
		ID: CmdScaleApplication,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdTrust: {
		ID: CmdTrust,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdDashboard: {
		ID: CmdDashboard,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdAttachResource: {
		ID: CmdAttachResource,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdResources: {
		ID: CmdResources,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdCharmResources: {
		ID: CmdCharmResources,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
			"o":                true, // Output file operations
			"output":           true, // Output file operations
		},
	},
	CmdInfo: {
		ID: CmdInfo,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdFind: {
		ID: CmdFind,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdDownload: {
		ID: CmdDownload,
		HTTPMcpServerDisableArgs: map[string]bool{
			"filepath": true, // Local file path for download
		},
	},
	CmdSecrets: {
		ID: CmdSecrets,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdShowSecret: {
		ID: CmdShowSecret,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdAddSecret: {
		ID: CmdAddSecret,
		HTTPMcpServerDisableArgs: map[string]bool{
			"file": true, // Secret values file
		},
	},
	CmdUpdateSecret: {
		ID: CmdUpdateSecret,
		HTTPMcpServerDisableArgs: map[string]bool{
			"file": true, // Secret values file
		},
	},
	CmdRemoveSecret: {
		ID: CmdRemoveSecret,
	},
	CmdGrantSecret: {
		ID: CmdGrantSecret,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdRevokeSecret: {
		ID: CmdRevokeSecret,
		HTTPMcpServerDisableArgs: map[string]bool{
			"B":                true, // Browser operations on remote server
			"no-browser-login": true, // Browser operations on remote server
		},
	},
	CmdSecretBackends: {
		ID: CmdSecretBackends,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdAddSecretBackend: {
		ID: CmdAddSecretBackend,
		HTTPMcpServerDisableArgs: map[string]bool{
			"config": true, // Config file operations
		},
	},
	CmdUpdateSecretBackend: {
		ID: CmdUpdateSecretBackend,
		HTTPMcpServerDisableArgs: map[string]bool{
			"config": true, // Config file operations
		},
	},
	CmdRemoveSecretBackend: {
		ID: CmdRemoveSecretBackend,
	},
	CmdShowSecretBackend: {
		ID: CmdShowSecretBackend,
		HTTPMcpServerDisableArgs: map[string]bool{
			"o":      true, // Output file operations
			"output": true, // Output file operations
		},
	},
	CmdWaitFor: {
		ID: CmdWaitFor,
	},
}

// GetCommandDefinition returns the command definition for a given ID
func GetCommandDefinition(id JujuCommandID) (CommandDefinition, bool) {
	def, exists := commandRegistry[id]
	return def, exists
}

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
