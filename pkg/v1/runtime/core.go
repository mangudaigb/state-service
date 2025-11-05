package runtime

import "time"

//
// ──────────────────────────────── CORE ENTITIES ────────────────────────────────
//

// Context holds the active working state of the system at a given step.
type Context struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Cognitive map[string]string `json:"cognitive,omitempty"` // ephemeral reasoning or variables
	Workspace map[string]string `json:"workspace,omitempty"` // environment states like paths or files or rest endpoints
	Knowledge map[string]string `json:"knowledge,omitempty"` // persistent facts or goals

	Logs      map[string]string `json:"logs,omitempty"`
	Metrics   map[string]string `json:"metrics,omitempty"`
	Systems   map[string]string `json:"systems,omitempty"`
	Incidents map[string]string `json:"incidents,omitempty"`
}

// Query represents the user's or planner's request/question.
type Query struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Tags      []string          `json:"tags,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// Answer captures the system’s or agent’s answer/output.
type Answer struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

//
// ──────────────────────────────── MCP & TOOLING ────────────────────────────────
//

type ToolCategory string

const (
	ToolCategoryLogs      ToolCategory = "logs"
	ToolCategoryDatabases ToolCategory = "databases"
	ToolCategoryDocs      ToolCategory = "docs"
	ToolCategoryMetrics   ToolCategory = "metrics"
	ToolCategorySystems   ToolCategory = "systems"
	ToolCategoryIncidents ToolCategory = "incidents"
	ToolCategoryPlanner   ToolCategory = "planner"
)

// Tool is an executable capability provided by an MCP.
type Tool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    ToolCategory      `json:"category"`
	Inputs      map[string]string `json:"inputs,omitempty"`
	Outputs     map[string]string `json:"outputs,omitempty"`
}

// MCP represents a model capability provider (e.g., GitHub, FileSystem, Terminal).
type MCP struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tools       []Tool            `json:"tools"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type MCPToolRef struct {
	McpID    string       `json:"mcpId"`
	ToolName string       `json:"toolName"`
	Category ToolCategory `json:"category"`
}

//
// ──────────────────────────────── AGENTS ────────────────────────────────
//

// Agent represents a worker agent (planner, coder, executor, etc.)
type Agent struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Model         string            `json:"model"`
	Role          string            `json:"role"` // planner, coder, tester, etc.
	SystemPrompt  string            `json:"systemPrompt"`
	UserPrompt    string            `json:"userPrompt"`
	Capabilities  []string          `json:"capabilities"`
	Parameters    map[string]string `json:"parameters"`
	LastUpdatedAt time.Time         `json:"lastUpdatedAt"`
}

type AgentRef struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}
