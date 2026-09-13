package jujuadapter

import "github.com/mark3labs/mcp-go/mcp"

// toolHints describes the MCP tool annotations for a command.
//
// The MCP specification treats an un-annotated tool as non read-only,
// destructive, non-idempotent and open-world, so every command is classified
// explicitly here. A unit test enforces full coverage.
type toolHints struct {
	readOnly    bool
	destructive bool
	idempotent  bool
	openWorld   bool
}

// Classification helpers. Names describe the effect the command has on the
// Juju controller or the local client store.
var (
	// hintReadOnly: only reads state; safe to call repeatedly.
	hintReadOnly = toolHints{readOnly: true, destructive: false, idempotent: true, openWorld: false}
	// hintReadOnlyOpenWorld: reads from an external service such as Charmhub.
	hintReadOnlyOpenWorld = toolHints{readOnly: true, destructive: false, idempotent: true, openWorld: true}
	// hintReadWrite: reads or updates settings depending on arguments; re-running
	// with the same arguments has no additional effect.
	hintReadWrite = toolHints{readOnly: false, destructive: false, idempotent: true, openWorld: false}
	// hintAdditive: creates or extends state without deleting anything.
	hintAdditive = toolHints{readOnly: false, destructive: false, idempotent: false, openWorld: false}
	// hintAdditiveIdempotent: additive and safe to repeat with the same arguments.
	hintAdditiveIdempotent = toolHints{readOnly: false, destructive: false, idempotent: true, openWorld: false}
	// hintAdditiveOpenWorld: additive and talks to an external service.
	hintAdditiveOpenWorld = toolHints{readOnly: false, destructive: false, idempotent: false, openWorld: true}
	// hintDestructive: removes or replaces state, or upgrades software.
	hintDestructive = toolHints{readOnly: false, destructive: true, idempotent: false, openWorld: false}
	// hintDestructiveIdempotent: destructive, but repeating it changes nothing more.
	hintDestructiveIdempotent = toolHints{readOnly: false, destructive: true, idempotent: true, openWorld: false}
	// hintDestructiveOpenWorld: runs arbitrary commands or reaches arbitrary hosts.
	hintDestructiveOpenWorld = toolHints{readOnly: false, destructive: true, idempotent: false, openWorld: true}
)

// commandHints maps every command ID to its annotations.
var commandHints = map[JujuCommandID]toolHints{
	// Reporting
	CmdVersion:       hintReadOnly,
	CmdStatus:        hintReadOnly,
	CmdShowStatusLog: hintReadOnly,
	CmdSwitch:        hintReadWrite,

	// Creation
	CmdBootstrap: hintAdditive,
	CmdIntegrate: hintAdditive,

	// Cross model relations
	CmdOffer:           hintAdditive,
	CmdRemoveOffer:     hintDestructive,
	CmdShowOffer:       hintReadOnly,
	CmdOffers:          hintReadOnly,
	CmdFindOffers:      hintReadOnly,
	CmdConsume:         hintAdditive,
	CmdSuspendRelation: hintDestructiveIdempotent,
	CmdResumeRelation:  hintAdditiveIdempotent,

	// Firewall rules
	CmdSetFirewallRule: hintAdditiveIdempotent,
	CmdFirewallRules:   hintReadOnly,

	// Destruction
	CmdRemoveRelation:    hintDestructive,
	CmdRemoveApplication: hintDestructive,
	CmdRemoveUnit:        hintDestructive,
	CmdRemoveSaas:        hintDestructive,

	// Error resolution and debugging
	CmdExec:       hintDestructiveOpenWorld,
	CmdScp:        hintDestructiveOpenWorld,
	CmdSsh:        hintDestructiveOpenWorld,
	CmdResolved:   hintAdditiveIdempotent,
	CmdDebugLog:   hintReadOnly,
	CmdDebugHooks: hintDestructiveOpenWorld,
	CmdDebugCode:  hintDestructiveOpenWorld,

	// Configuration
	CmdConstraints:       hintReadOnly,
	CmdSetConstraints:    hintAdditiveIdempotent,
	CmdSyncAgentBinary:   hintAdditiveOpenWorld,
	CmdUpgradeModel:      hintDestructive,
	CmdUpgradeController: hintDestructive,
	CmdRefresh:           hintDestructive,
	CmdBind:              hintAdditiveIdempotent,

	// Charm tool help
	CmdHelpHookCommands:   hintReadOnly,
	CmdHelpActionCommands: hintReadOnly,

	// Backups
	CmdCreateBackup:   hintAdditive,
	CmdDownloadBackup: hintAdditive,

	// SSH keys
	CmdAddSshKey:    hintAdditiveIdempotent,
	CmdRemoveSshKey: hintDestructive,
	CmdImportSshKey: hintAdditiveOpenWorld,
	CmdSshKeys:      hintReadOnly,

	// Users and access
	CmdAddUser:            hintAdditive,
	CmdChangeUserPassword: hintAdditive,
	CmdShowUser:           hintReadOnly,
	CmdUsers:              hintReadOnly,
	CmdEnableUser:         hintAdditiveIdempotent,
	CmdDisableUser:        hintDestructiveIdempotent,
	CmdLogin:              hintAdditive,
	CmdLogout:             hintDestructiveIdempotent,
	CmdRemoveUser:         hintDestructive,
	CmdWhoami:             hintReadOnly,

	// Machines
	CmdAddMachine:     hintAdditive,
	CmdRemoveMachine:  hintDestructive,
	CmdMachines:       hintReadOnly,
	CmdShowMachine:    hintReadOnly,
	CmdUpgradeMachine: hintDestructive,

	// Models
	CmdModelConfig:         hintReadWrite,
	CmdModelDefaults:       hintReadWrite,
	CmdModelConstraints:    hintReadOnly,
	CmdSetModelConstraints: hintAdditiveIdempotent,
	CmdRetryProvisioning:   hintAdditiveIdempotent,
	CmdDestroyModel:        hintDestructive,
	CmdGrant:               hintAdditiveIdempotent,
	CmdRevoke:              hintDestructiveIdempotent,
	CmdShowModel:           hintReadOnly,
	CmdSetCredential:       hintAdditiveIdempotent,
	CmdMigrate:             hintDestructive,
	CmdExportBundle:        hintReadOnly,

	// Actions
	CmdActions:       hintReadOnly,
	CmdShowAction:    hintReadOnly,
	CmdCancelTask:    hintDestructive,
	CmdRun:           hintDestructive,
	CmdOperations:    hintReadOnly,
	CmdShowOperation: hintReadOnly,
	CmdShowTask:      hintReadOnly,

	// Controller HA
	CmdEnableHa: hintAdditiveIdempotent,

	// Applications
	CmdAddUnit:            hintAdditive,
	CmdConfig:             hintReadWrite,
	CmdDeploy:             hintAdditive,
	CmdExpose:             hintAdditiveIdempotent,
	CmdUnexpose:           hintDestructiveIdempotent,
	CmdDiffBundle:         hintReadOnly,
	CmdShowApplication:    hintReadOnly,
	CmdShowUnit:           hintReadOnly,
	CmdSetApplicationBase: hintAdditiveIdempotent,

	// Operation protection
	CmdDisableCommand:   hintDestructiveIdempotent,
	CmdDisabledCommands: hintReadOnly,
	CmdEnableCommand:    hintAdditiveIdempotent,

	// Storage
	CmdAddStorage:        hintAdditive,
	CmdStorage:           hintReadOnly,
	CmdCreateStoragePool: hintAdditive,
	CmdStoragePools:      hintReadOnly,
	CmdRemoveStoragePool: hintDestructive,
	CmdUpdateStoragePool: hintAdditiveIdempotent,
	CmdShowStorage:       hintReadOnly,
	CmdRemoveStorage:     hintDestructive,
	CmdDetachStorage:     hintDestructive,
	CmdAttachStorage:     hintAdditive,
	CmdImportFilesystem:  hintAdditive,

	// Spaces and subnets
	CmdAddSpace:     hintAdditive,
	CmdSpaces:       hintReadOnly,
	CmdMoveToSpace:  hintAdditiveIdempotent,
	CmdReloadSpaces: hintAdditiveIdempotent,
	CmdShowSpace:    hintReadOnly,
	CmdRemoveSpace:  hintDestructive,
	CmdRenameSpace:  hintAdditive,
	CmdSubnets:      hintReadOnly,

	// Controllers
	CmdAddModel:                hintAdditive,
	CmdDestroyController:       hintDestructive,
	CmdModels:                  hintReadOnly,
	CmdKillController:          hintDestructive,
	CmdControllers:             hintReadOnly,
	CmdRegister:                hintAdditiveOpenWorld,
	CmdUnregister:              hintDestructive,
	CmdEnableDestroyController: hintAdditiveIdempotent,
	CmdShowController:          hintReadOnly,
	CmdControllerConfig:        hintReadWrite,

	// Clouds and credentials
	CmdUpdateCloud:         hintAdditiveIdempotent,
	CmdUpdatePublicClouds:  hintAdditiveOpenWorld,
	CmdClouds:              hintReadOnly,
	CmdRegions:             hintReadOnly,
	CmdShowCloud:           hintReadOnly,
	CmdAddCloud:            hintAdditive,
	CmdRemoveCloud:         hintDestructive,
	CmdCredentials:         hintReadOnly,
	CmdAutoloadCredentials: hintAdditive,
	CmdDefaultRegion:       hintReadWrite,
	CmdDefaultCredential:   hintReadWrite,
	CmdAddCredential:       hintAdditive,
	CmdRemoveCredential:    hintDestructive,
	CmdUpdateCredential:    hintAdditiveIdempotent,
	CmdShowCredential:      hintReadOnly,
	CmdGrantCloud:          hintAdditiveIdempotent,
	CmdRevokeCloud:         hintDestructiveIdempotent,

	// Kubernetes clouds
	CmdAddK8s:    hintAdditive,
	CmdUpdateK8s: hintAdditiveIdempotent,
	CmdRemoveK8s: hintDestructive,

	// Scaling and trust
	CmdScaleApplication: hintDestructiveIdempotent,
	CmdTrust:            hintAdditiveIdempotent,

	// Dashboard
	CmdDashboard: hintReadOnlyOpenWorld,

	// Resources
	CmdAttachResource: hintAdditive,
	CmdResources:      hintReadOnly,
	CmdCharmResources: hintReadOnlyOpenWorld,

	// Payloads
	CmdPayloads: hintReadOnly,

	// Charmhub
	CmdInfo:     hintReadOnlyOpenWorld,
	CmdFind:     hintReadOnlyOpenWorld,
	CmdDownload: hintAdditiveOpenWorld,

	// Secrets
	CmdSecrets:      hintReadOnly,
	CmdShowSecret:   hintReadOnly,
	CmdAddSecret:    hintAdditive,
	CmdUpdateSecret: hintAdditive,
	CmdRemoveSecret: hintDestructive,
	CmdGrantSecret:  hintAdditiveIdempotent,
	CmdRevokeSecret: hintDestructiveIdempotent,

	// Secret backends
	CmdSecretBackends:      hintReadOnly,
	CmdAddSecretBackend:    hintAdditive,
	CmdUpdateSecretBackend: hintAdditiveIdempotent,
	CmdRemoveSecretBackend: hintDestructive,
	CmdShowSecretBackend:   hintReadOnly,

	// Wait for
	CmdWaitFor: hintReadOnly,
}

// hintsFor returns the annotations for a command. Unknown commands fall back
// to the most conservative classification.
func hintsFor(id JujuCommandID) toolHints {
	if h, ok := commandHints[id]; ok {
		return h
	}
	return hintDestructiveOpenWorld
}

// toolAnnotation converts hints into the MCP annotation structure.
func (h toolHints) toolAnnotation(title string) mcp.ToolAnnotation {
	return mcp.ToolAnnotation{
		Title:           title,
		ReadOnlyHint:    mcp.ToBoolPtr(h.readOnly),
		DestructiveHint: mcp.ToBoolPtr(h.destructive),
		IdempotentHint:  mcp.ToBoolPtr(h.idempotent),
		OpenWorldHint:   mcp.ToBoolPtr(h.openWorld),
	}
}
