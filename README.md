# OpenCode Go Kit

A Go API for interacting with the OpenCode API.

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
		"http://localhost:4002", // OpenCode Server URL
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

## Requirements

- **Go**: 1.26.4+ (to compile, run tests, and generate client bindings)
- **Docker & Docker Compose**: To build and host the OpenCode agent locally (or connect to your existing opencode)

## Getting Started with Docker

- NOTE: DOCKERFILE NEEDS OPENCODE CONFIG MOUNTED, THIS IS UNIMPLEMENTED, AND CONTAINER IS NOT 
- UNOPTIMIZED FOR RUNNING AGENTS. WILL FIX SOON. AT THIS TIME, CONTAINERS MOSTLY FUNCTION TO PARSE OPENCODE 
- SERVER API FOR GO BINDING GENERATION.
- 

1. Copy the environment variables template to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Copy the Docker Compose template to `docker-compose.yml`:
   ```bash
   cp docker-compose.example.yml docker-compose.yml
   ```

3. Open `.env` and customize parameters such as `HOST_PORT`, `VOLUME_PATH`, `WORKSPACE_PASSWORD`, and `PROJECT_NAME` if desired.

## Using Standard Docker Commands (No Makefile)

 `make`, cmds below:

### 1. Standard Docker Compose Commands

- **using Compose**:
  ```bash
  docker compose up -p PROJECT_NAME -d --build
  ```

### 2. Docker CLI Commands

- **Build the Image**:
  ```bash
  docker build \
    --build-arg USER_ID=$(id -u) \
    --build-arg GROUP_ID=$(id -g) \
    --build-arg USER_NAME=opencode \
    --build-arg GROUP_NAME=opencode \
    -t opencode-go-kit:latest .
  ```

- **Run the Container**:
  ```bash
  docker run -d \
    --name opencode-agent \
    -p 127.0.0.1:4002:4096 \
    -v $(pwd)/volumes/workspace_one:/workspace \
    -e OPENCODE_SERVER_PASSWORD=pass1 \
    --restart unless-stopped \
    opencode-go-kit:latest
  ```

## Running Multiple Instances

To run multiple parallel instances of the OpenCode agent concurrently:

1. Open `docker-compose.yml`.
2. Locate the commented-out `opencode-go-kit-2` service.
3. Uncomment the block. 
4. Run `make up` to boot both containers simultaneously. Each container binds to its own unique host port and mounts its own dedicated workspace volume (e.g. `./volumes/workspace_one` and `./volumes/workspace_two`).

## Makefile Commands

`Makefile` commands:

- **Build the custom agent Docker image**:
  ```bash
  make build
  ```

- **Run via Docker Compose** (Starts the containers in background):
  ```bash
  make up
  ```

- **Run direct Docker container** (Starts a single container using standard `docker run` without Compose):
  ```bash
  make run
  ```
## Code Gen 

- **Regenerate Go client and API wrappers dynamically**:
  Checks your api coverage. If is partial coverage fetch the latest schema from the running local container to regenerate client bindings, then rebuilds the typed wrappers in the `api` package:
  
  ```bash
  make generate-api
  ```



## API Endpoint Coverage

<!-- COVERAGE_START -->
[![API Coverage](https://img.shields.io/badge/Coverage-100.00%25-brightgreen)](#)

We have wrapped **83 out of 83** (100.00%) OpenAPI client methods in our clean, typed API wrapper layer.

### Covered Endpoints

- [x] `AppAgents`
- [x] `AppSkills`
- [x] `CommandList`
- [x] `ConfigProviders`
- [x] `ExperimentalCapabilitiesGet`
- [x] `ExperimentalConsoleGet`
- [x] `ExperimentalConsoleListOrgs`
- [x] `ExperimentalConsoleSwitchOrg`
- [x] `ExperimentalProjectCopyGenerateName`
- [x] `ExperimentalResourceList`
- [x] `ExperimentalSessionList`
- [x] `ExperimentalWorkspaceAdapterList`
- [x] `ExperimentalWorkspaceStatus`
- [x] `ExperimentalWorkspaceSyncList`
- [x] `FileList`
- [x] `FileRead`
- [x] `FileStatus`
- [x] `FindFiles`
- [x] `FindSymbols`
- [x] `FindText`
- [x] `FormatterStatus`
- [x] `GlobalDispose`
- [x] `GlobalHealth`
- [x] `InstanceDispose`
- [x] `LspStatus`
- [x] `McpAuthRemove`
- [x] `McpConnect`
- [x] `McpDisconnect`
- [x] `PathGet`
- [x] `PermissionList`
- [x] `ProjectCurrent`
- [x] `ProjectDirectories`
- [x] `ProjectInitGit`
- [x] `ProjectList`
- [x] `ProviderList`
- [x] `PtyConnect`
- [x] `PtyConnectToken`
- [x] `PtyGet`
- [x] `PtyList`
- [x] `PtyRemove`
- [x] `PtyShells`
- [x] `QuestionList`
- [x] `SessionDiff`
- [x] `SessionList`
- [x] `SessionShare`
- [x] `SessionUnshare`
- [x] `SyncStart`
- [x] `TuiClearPrompt`
- [x] `TuiControlNext`
- [x] `TuiControlResponse`
- [x] `TuiOpenHelp`
- [x] `TuiOpenModels`
- [x] `TuiOpenSessions`
- [x] `TuiOpenThemes`
- [x] `TuiShowToast`
- [x] `TuiSubmitPrompt`
- [x] `V2CommandList`
- [x] `V2CredentialRemove`
- [x] `V2CredentialUpdate`
- [x] `V2FsFind`
- [x] `V2FsList`
- [x] `V2FsRead`
- [x] `V2HealthGet`
- [x] `V2IntegrationAttemptCancel`
- [x] `V2LocationGet`
- [x] `V2PermissionRequestList`
- [x] `V2PermissionSavedList`
- [x] `V2PermissionSavedRemove`
- [x] `V2PtyConnect`
- [x] `V2PtyConnectToken`
- [x] `V2PtyCreate`
- [x] `V2PtyGet`
- [x] `V2PtyList`
- [x] `V2PtyRemove`
- [x] `V2PtyUpdate`
- [x] `V2QuestionRequestList`
- [x] `V2SessionActive`
- [x] `V2SessionCreate`
- [x] `V2SkillList`
- [x] `VcsDiff`
- [x] `VcsDiffRaw`
- [x] `VcsGet`
- [x] `VcsStatus`

<!-- COVERAGE_END -->
