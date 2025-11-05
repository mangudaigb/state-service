
# State Service

## Overview
It defines entities for:

- **Context & State** — the evolving system knowledge and workspace.
- **Queries & Answers** — user requests and system responses.
- **Agents** — intelligent actors (planner, executor, diagnostics, etc.).
- **Workflows** — structured multi-step plans with dependencies.
- **Execution Graphs** — directed graphs of executable steps.
- **Artifacts** — tangible outputs like logs, queries, or remediation scripts.
- **MCP Tools** — model capability providers and their invocations.
- **Conversations & Snapshots** — persistent, replayable histories of interactions.

## Core Entities

### `Context`
Represents the **active working state** of the system during a given step.

It holds multiple contextual scopes:
- `Cognitive`: ephemeral reasoning state (temporary variables, hypotheses).
- `Workspace`: environment variables, paths, or runtime endpoints.
- `Knowledge`: persistent facts, goals, or insights gathered so far.
- Optional domain-specific maps for `Logs`, `Metrics`, `Systems`, and `Incidents`.

**Used in:** `Step.InputContext`, `Step.OutputContext`, and `Interaction.BaseContext`.

### `Query`
Represents a **user or planner request**, such as:
> “I’m not able to connect to the server.”

Includes tags and metadata for classification, and a timestamp.

**Used in:** `Interaction.BaseQuery`, `Step.Query`.

### `Answer`
Captures the **system’s response** or resolution to a `Query`.

**Used in:** `Step.Answer`.

### `Message`
A generic wrapper around both `Query` and `Answer` that enables full conversational traceability (with role and timestamp).

**Used in:** `Interaction.Messages`.

## Agents

### `Agent`
Defines a **worker entity** (e.g., Planner, Diagnostics, Remediator).

Each Agent has:
- A `Model` (e.g., GPT, Claude, local LLM)
- A `Role` (planner, diagnostics, executor)
- A `SystemPrompt` and `UserPrompt` (behavioral configuration)
- `Capabilities` (supported operations)
- `Parameters` (tunable runtime settings)
- `LastUpdatedAt` timestamp (version tracking)

**Used in:** `Step.Agent`, `Workflow.Agents`.

---

### `AgentRef`
A lightweight reference used in workflows, linking `ID` to `Role`.

---

## MCP and Tooling Layer

### `MCP` (Model Capability Provider)
Represents an external **capability source** such as:
- GitHub
- FileSystem
- Monitoring APIs
- Database connections

Each MCP can expose multiple **Tools**.

### `Tool`
Defines an **executable capability** offered by an MCP:
- Example: `ping_check`, `query_metrics`, `read_logs`
- Declares its inputs and outputs schema

### `McpToolInvocation`
Represents a **runtime invocation** of a specific tool by an agent during a step.

Includes:
- Tool metadata (`ToolName`, `MCPID`, `Category`)
- `Input` and `Output`
- `Status`, `Error`, and timestamps
- `AgentID` and `StepID` linkage

This is crucial for auditing and tracing fine-grained tool usage.

## Workflow and Planning
### `Plan` and `PlanStep`
The **Planner Agent** generates a `Plan` consisting of ordered `PlanStep`s, which define what needs to be done (e.g., "Check logs", "Run ping", "Diagnose root cause").

A `Plan` is purely declarative — the `ExecutionGraph` defines how these steps are executed.

### `ExecutionGraph`
Represents the **directed graph** of execution.

- `Nodes`: actual `Step` instances (with execution results).
- `Edges`: `Edge{From, To, Type}` defines dependencies, e.g.:
    - `depends_on` → must complete before next step.
    - `depends_on_any` → triggered if any parent step succeeds.
    - `triggers` → conditional event linkage.

This graph enables both **sequential** and **parallel** execution.

### `Workflow`
Describes a **full multi-agent process**:
- `Name` and `Description`
- Participating `Agents`
- Optional `Variables`
- `Mode` (`sequential` or `parallel`)
- Embedded `ExecutionGraph`
- Link to source `Plan`

The `Workflow` is the operational backbone of an interaction.

## Execution Model

### `Step`
The fundamental unit of execution — each **Step** represents a concrete action performed by an agent.

Contains:
- `Agent`: the agent executing the step
- `InputContext` / `OutputContext`: environment before and after execution
- `Query` / `Answer`: reasoning input and result
- `Artifacts`: tangible outputs (logs, queries, remediations)
- `CuratedTools`: detailed record of invoked tools
- `Status`, `Error`, and timestamps for observability
- `InputStepID`: linkage for dependency tracing

### `Artifact`
Represents **produced outputs** such as:
- `log_snippet` — captured log segments
- `query_result` — database or metrics query output
- `doc_summary` — documentation summaries
- `root_cause` — RCA findings
- `remediation` — actionable fixes

Artifacts enable persistence and reusability of generated data across steps.

## Interaction and Conversation Layer

### `Interaction`
Encapsulates a **single end-to-end reasoning episode** — one full workflow triggered by a base query.

Includes:
- `BaseQuery`: user or system request
- `BaseContext`: initial environment or knowledge state
- `Workflow`: the structured plan and agents
- `Steps`: executed steps with inputs, outputs, and artifacts
- `Messages`: communication trace
- `Summary`: final explanation or remediation
- Lifecycle timestamps (`CreatedAt`, `CompletedAt`)

### `Conversation`
Groups multiple `Interactions` — e.g., an SRE troubleshooting session or multi-turn conversation with several related tasks.

### `Snapshot`
Captures the **entire execution state** of a conversation for:
- Replay
- Post-mortem analysis
- Historical auditing

Includes:
- `StateRef`: full `Conversation` object
- `Timestamp`: time of snapshot

## Example Use Case

### User query:
> “I am not able to connect to a server.”

**System flow:**
1. The user query becomes a new `Interaction` with `BaseQuery`.
2. The `Planner Agent` generates a `Plan` with steps like:
    - Gather system logs
    - Run connectivity tests
    - Identify root cause
    - Suggest remediation
3. The `ExecutionGraph` defines dependencies between steps.
4. Each `Step` executes under a corresponding `Agent`, invoking MCP tools (e.g., `ping_check`, `log_fetch`).
5. Artifacts such as log snippets and remediation suggestions are stored.
6. The final `Answer` summarizes the findings (e.g., “Firewall rule blocking traffic to port 443”).
7. A `Snapshot` may be saved for audit or replay.

## Design Idea

`runtimev2` is built to support:

- **Multi-agent orchestration** — agents collaborate under structured workflows.
- **Observability** — every action, tool invocation, and artifact is traceable.
- **Reproducibility** — via context snapshots and deterministic graph execution.
- **Extensibility** — adding new agents, MCPs, or artifact types requires no schema changes.
- **Neutrality** — no direct coupling to a specific model or provider.

## Example Execution State
```json
{
  "id": "interaction-001",
  "base_query": {
    "id": "query-001",
    "content": "I am not able to connect to a server",
    "tags": ["connectivity", "incident"],
    "metadata": {
      "severity": "medium",
      "source": "user"
    },
    "timestamp": "2025-11-04T05:12:00Z"
  },
  "base_context": {
    "id": "ctx-001",
    "content": "Initial troubleshooting context",
    "workspace": {
      "target_host": "10.0.0.15"
    },
    "knowledge": {
      "known_systems": "webapp-prod, db-primary, cache-tier"
    }
  },
  "workflow": {
    "id": "workflow-001",
    "name": "Server Connectivity Troubleshooting",
    "description": "Diagnose network or service issues preventing server connection.",
    "agents": [
      { "id": "agent-planner", "role": "planner" },
      { "id": "agent-diagnostics", "role": "diagnostics" },
      { "id": "agent-remediator", "role": "remediator" }
    ],
    "mode": "sequential",
    "graph": {
      "id": "graph-001",
      "nodes": ["step-1", "step-2", "step-3"],
      "edges": [
        { "from": "step-1", "to": "step-2", "type": "depends_on", "label": "next" },
        { "from": "step-2", "to": "step-3", "type": "depends_on", "label": "next" }
      ]
    },
    "plan_id": "plan-001"
  },
  "steps": [
    {
      "id": "run-001",
      "step_id": "step-1",
      "index": 1,
      "name": "Gather System Logs",
      "status": "success",
      "agent": {
        "id": "agent-diagnostics",
        "name": "Diagnostics Agent",
        "description": "Collects logs and metrics",
        "model": "gpt-4o",
        "role": "diagnostics",
        "system_prompt": "Collect logs from hosts and analyze anomalies.",
        "capabilities": ["log_analysis", "ping_check", "service_status"],
        "parameters": {},
        "last_updated_at": "2025-11-03T22:00:00Z"
      },
      "query": {
        "id": "query-step1",
        "content": "Check recent logs for network errors or service downtime.",
        "timestamp": "2025-11-04T05:12:10Z"
      },
      "answer": {
        "id": "answer-step1",
        "content": "Logs show repeated connection timeouts to 10.0.0.15 from proxy layer.",
        "timestamp": "2025-11-04T05:12:12Z"
      },
      "artifacts": [
        {
          "id": "artifact-001",
          "name": "log_snippet_2025-11-04",
          "path": "/var/log/syslog",
          "type": "log_snippet",
          "content": {
            "snippet": "Nov 04 05:11: timeout connecting to 10.0.0.15 port 443"
          },
          "created_by_step_id": "step-1",
          "created_at": "2025-11-04T05:12:12Z"
        }
      ],
      "started_at": "2025-11-04T05:12:05Z",
      "finished_at": "2025-11-04T05:12:13Z"
    },
    {
      "id": "run-002",
      "step_id": "step-2",
      "index": 2,
      "name": "Run Network Connectivity Checks",
      "status": "error",
      "error": "Ping to 10.0.0.15 failed",
      "agent": {
        "id": "agent-diagnostics",
        "name": "Diagnostics Agent",
        "model": "gpt-4o",
        "role": "diagnostics"
      },
      "curated_tools": [
        {
          "id": "toolinv-001",
          "tool_name": "ping_check",
          "mcp_id": "mcp-network",
          "category": "systems",
          "input": { "target": "10.0.0.15" },
          "output": { "result": "timeout" },
          "status": "error",
          "agent_id": "agent-diagnostics",
          "step_id": "step-2",
          "started_at": "2025-11-04T05:12:15Z",
          "finished_at": "2025-11-04T05:12:16Z",
          "error": "Destination host unreachable"
        }
      ],
      "started_at": "2025-11-04T05:12:14Z",
      "finished_at": "2025-11-04T05:12:17Z"
    },
    {
      "id": "run-003",
      "step_id": "step-3",
      "index": 3,
      "name": "Propose Remediation",
      "status": "success",
      "agent": {
        "id": "agent-remediator",
        "name": "Remediation Agent",
        "model": "gpt-4o",
        "role": "remediator"
      },
      "query": {
        "id": "query-step3",
        "content": "Given the ping failure, what remediation should be suggested?",
        "timestamp": "2025-11-04T05:12:18Z"
      },
      "answer": {
        "id": "answer-step3",
        "content": "The issue is likely a network routing or firewall block. Recommend verifying security group and VPC route configuration for 10.0.0.15.",
        "timestamp": "2025-11-04T05:12:20Z"
      },
      "artifacts": [
        {
          "id": "artifact-002",
          "name": "root_cause_analysis",
          "path": "/diagnostics/root_cause_001.json",
          "type": "root_cause",
          "content": {
            "summary": "Network connectivity failure due to firewall restriction."
          },
          "created_by_step_id": "step-3",
          "created_at": "2025-11-04T05:12:21Z"
        }
      ],
      "started_at": "2025-11-04T05:12:18Z",
      "finished_at": "2025-11-04T05:12:22Z"
    }
  ],
  "summary": "The connectivity issue was caused by a network-level block. Recommended remediation: update firewall rules or verify VPC routes.",
  "created_at": "2025-11-04T05:12:00Z",
  "completed_at": "2025-11-04T05:12:23Z"
}

```