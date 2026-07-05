# OpenCode Go Kit

A Go API for interacting with the OpenCode API.

## Requirements

- **Go**: 1.26.4+ (to compile, run tests, and generate client bindings)


## Installation

To integrate the OpenCode Go Kit into an existing Go project, run:

```bash
go get github.com/dakolli/opencode-go-kit
```

## Quick Start Example

Here is a simple example demonstrating how to configure the client, instantiate the API, 
and call a typed API method (e.g., `AppAgents`) returning Go types:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dakolli/opencode-go-kit/pkg/api"
	"github.com/dakolli/opencode-go-kit/pkg/client"
)

func main() {
	ctx := context.Background()

	// 1. Create client configuration (either pass values explicitly or let them fallback to environment variables)
	config, err := api.NewClientCFG(
		"http://localhost:4096", // OpenCode Server URL
		"opencode",              // Username
		"pass1",                 // Password
	)
	if err != nil {
		log.Fatalf("failed to initialize configuration: %v", err)
	}

	// 2. Instantiate the clean, typed API wrapper
	opencodeAPI, err := api.NewAPI(config)
	if err != nil {
		log.Fatalf("failed to initialize API: %v", err)
	}

	// 3. Invoke a typed API method with full parameter support
	// Our wrapper automatically unwraps sum-type interfaces into clean, native slice types
	agents, err := opencodeAPI.AppAgents(ctx, client.AppAgentsParams{})
	if err != nil {
		log.Fatalf("API invocation failed: %v", err)
	}

	fmt.Printf("Retrieved %d AI agents from OpenCode:\n", len(agents))
	for _, agent := range agents {
		fmt.Printf(" - %s: %s\n", agent.Name, agent.Description.Value)
	}
}
```



## API Endpoint Coverage

<!-- COVERAGE_START -->
[![API Coverage](https://img.shields.io/badge/Coverage-100.00%25-brightgreen)](#)

We have wrapped **83 out of 83** (100.00%) OpenAPI client methods in our clean, typed API wrapper layer.

### Covered Endpoints

| | |
| --- | --- |
| - [x] [`AppAgents`](https://dakolli.github.io/opencode-go-kit/api/#func-api-appagents) | - [x] [`SessionDiff`](https://dakolli.github.io/opencode-go-kit/api/#func-api-sessiondiff) |
| - [x] [`AppSkills`](https://dakolli.github.io/opencode-go-kit/api/#func-api-appskills) | - [x] [`SessionList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-sessionlist) |
| - [x] [`CommandList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-commandlist) | - [x] [`SessionShare`](https://dakolli.github.io/opencode-go-kit/api/#func-api-sessionshare) |
| - [x] [`ConfigProviders`](https://dakolli.github.io/opencode-go-kit/api/#func-api-configproviders) | - [x] [`SessionUnshare`](https://dakolli.github.io/opencode-go-kit/api/#func-api-sessionunshare) |
| - [x] [`ExperimentalCapabilitiesGet`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalcapabilitiesget) | - [x] [`SyncStart`](https://dakolli.github.io/opencode-go-kit/api/#func-api-syncstart) |
| - [x] [`ExperimentalConsoleGet`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalconsoleget) | - [x] [`TuiClearPrompt`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuiclearprompt) |
| - [x] [`ExperimentalConsoleListOrgs`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalconsolelistorgs) | - [x] [`TuiControlNext`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuicontrolnext) |
| - [x] [`ExperimentalConsoleSwitchOrg`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalconsoleswitchorg) | - [x] [`TuiControlResponse`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuicontrolresponse) |
| - [x] [`ExperimentalProjectCopyGenerateName`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalprojectcopygeneratename) | - [x] [`TuiOpenHelp`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuiopenhelp) |
| - [x] [`ExperimentalResourceList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalresourcelist) | - [x] [`TuiOpenModels`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuiopenmodels) |
| - [x] [`ExperimentalSessionList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalsessionlist) | - [x] [`TuiOpenSessions`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuiopensessions) |
| - [x] [`ExperimentalWorkspaceAdapterList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalworkspaceadapterlist) | - [x] [`TuiOpenThemes`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuiopenthemes) |
| - [x] [`ExperimentalWorkspaceStatus`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalworkspacestatus) | - [x] [`TuiShowToast`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuishowtoast) |
| - [x] [`ExperimentalWorkspaceSyncList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-experimentalworkspacesynclist) | - [x] [`TuiSubmitPrompt`](https://dakolli.github.io/opencode-go-kit/api/#func-api-tuisubmitprompt) |
| - [x] [`FileList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-filelist) | - [x] [`V2CommandList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2commandlist) |
| - [x] [`FileRead`](https://dakolli.github.io/opencode-go-kit/api/#func-api-fileread) | - [x] [`V2CredentialRemove`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2credentialremove) |
| - [x] [`FileStatus`](https://dakolli.github.io/opencode-go-kit/api/#func-api-filestatus) | - [x] [`V2CredentialUpdate`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2credentialupdate) |
| - [x] [`FindFiles`](https://dakolli.github.io/opencode-go-kit/api/#func-api-findfiles) | - [x] [`V2FsFind`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2fsfind) |
| - [x] [`FindSymbols`](https://dakolli.github.io/opencode-go-kit/api/#func-api-findsymbols) | - [x] [`V2FsList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2fslist) |
| - [x] [`FindText`](https://dakolli.github.io/opencode-go-kit/api/#func-api-findtext) | - [x] [`V2FsRead`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2fsread) |
| - [x] [`FormatterStatus`](https://dakolli.github.io/opencode-go-kit/api/#func-api-formatterstatus) | - [x] [`V2HealthGet`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2healthget) |
| - [x] [`GlobalDispose`](https://dakolli.github.io/opencode-go-kit/api/#func-api-globaldispose) | - [x] [`V2IntegrationAttemptCancel`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2integrationattemptcancel) |
| - [x] [`GlobalHealth`](https://dakolli.github.io/opencode-go-kit/api/#func-api-globalhealth) | - [x] [`V2LocationGet`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2locationget) |
| - [x] [`InstanceDispose`](https://dakolli.github.io/opencode-go-kit/api/#func-api-instancedispose) | - [x] [`V2PermissionRequestList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2permissionrequestlist) |
| - [x] [`LspStatus`](https://dakolli.github.io/opencode-go-kit/api/#func-api-lspstatus) | - [x] [`V2PermissionSavedList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2permissionsavedlist) |
| - [x] [`McpAuthRemove`](https://dakolli.github.io/opencode-go-kit/api/#func-api-mcpauthremove) | - [x] [`V2PermissionSavedRemove`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2permissionsavedremove) |
| - [x] [`McpConnect`](https://dakolli.github.io/opencode-go-kit/api/#func-api-mcpconnect) | - [x] [`V2PtyConnect`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2ptyconnect) |
| - [x] [`McpDisconnect`](https://dakolli.github.io/opencode-go-kit/api/#func-api-mcpdisconnect) | - [x] [`V2PtyConnectToken`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2ptyconnecttoken) |
| - [x] [`PathGet`](https://dakolli.github.io/opencode-go-kit/api/#func-api-pathget) | - [x] [`V2PtyCreate`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2ptycreate) |
| - [x] [`PermissionList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-permissionlist) | - [x] [`V2PtyGet`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2ptyget) |
| - [x] [`ProjectCurrent`](https://dakolli.github.io/opencode-go-kit/api/#func-api-projectcurrent) | - [x] [`V2PtyList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2ptylist) |
| - [x] [`ProjectDirectories`](https://dakolli.github.io/opencode-go-kit/api/#func-api-projectdirectories) | - [x] [`V2PtyRemove`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2ptyremove) |
| - [x] [`ProjectInitGit`](https://dakolli.github.io/opencode-go-kit/api/#func-api-projectinitgit) | - [x] [`V2PtyUpdate`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2ptyupdate) |
| - [x] [`ProjectList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-projectlist) | - [x] [`V2QuestionRequestList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2questionrequestlist) |
| - [x] [`ProviderList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-providerlist) | - [x] [`V2SessionActive`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2sessionactive) |
| - [x] [`PtyConnect`](https://dakolli.github.io/opencode-go-kit/api/#func-api-ptyconnect) | - [x] [`V2SessionCreate`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2sessioncreate) |
| - [x] [`PtyConnectToken`](https://dakolli.github.io/opencode-go-kit/api/#func-api-ptyconnecttoken) | - [x] [`V2SkillList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-v2skilllist) |
| - [x] [`PtyGet`](https://dakolli.github.io/opencode-go-kit/api/#func-api-ptyget) | - [x] [`VcsDiff`](https://dakolli.github.io/opencode-go-kit/api/#func-api-vcsdiff) |
| - [x] [`PtyList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-ptylist) | - [x] [`VcsDiffRaw`](https://dakolli.github.io/opencode-go-kit/api/#func-api-vcsdiffraw) |
| - [x] [`PtyRemove`](https://dakolli.github.io/opencode-go-kit/api/#func-api-ptyremove) | - [x] [`VcsGet`](https://dakolli.github.io/opencode-go-kit/api/#func-api-vcsget) |
| - [x] [`PtyShells`](https://dakolli.github.io/opencode-go-kit/api/#func-api-ptyshells) | - [x] [`VcsStatus`](https://dakolli.github.io/opencode-go-kit/api/#func-api-vcsstatus) |
| - [x] [`QuestionList`](https://dakolli.github.io/opencode-go-kit/api/#func-api-questionlist) |  |

<!-- COVERAGE_END -->
