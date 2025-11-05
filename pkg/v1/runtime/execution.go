package runtime

import "time"

//
// ──────────────────────────────── EXECUTION ────────────────────────────────
//

type ArtifactType string

const (
	ArtifactLogSnippet  ArtifactType = "log_snippet"
	ArtifactQueryResult ArtifactType = "query_result"
	ArtifactDocSummary  ArtifactType = "doc_summary"
	ArtifactRootCause   ArtifactType = "root_cause"
	ArtifactRemediation ArtifactType = "remediation"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusStop    Status = "stop"
	StatusError   Status = "error"
	StatusSuccess Status = "success"
)

// Artifact represents tangible outputs (code files, logs, documents, etc.)
type Artifact struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Path            string            `json:"path"`
	Type            ArtifactType      `json:"type"`
	Content         map[string]string `json:"content"`
	CreatedByStepID string            `json:"createdByStepId"`
	CreatedAt       time.Time         `json:"createdAt"`
}

// Message generalizes both Query and Answer for message trace.
type Message struct {
	ID        string            `json:"id"`
	Role      string            `json:"role"` // user, planner, agent, etc.
	Content   string            `json:"content"`
	StepId    string            `json:"stepId"`
	AgentId   string            `json:"agentId"`
	Sequence  int               `json:"sequence"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// McpToolInvocation captures an invocation of a specific MCP Tool by an agent during a step.
type McpToolInvocation struct {
	ID         string            `json:"id"`
	McpId      string            `json:"mcpId"`
	ToolName   string            `json:"toolName"`
	Category   ToolCategory      `json:"category"`
	Input      map[string]string `json:"input,omitempty"`
	Output     map[string]string `json:"output,omitempty"`
	Status     Status            `json:"status"`
	Error      string            `json:"error,omitempty"`
	AgentId    string            `json:"agentId"`
	StepId     string            `json:"stepId"`
	StartedAt  time.Time         `json:"startedAt"`
	FinishedAt time.Time         `json:"finishedAt"`
}

// Step is a single executable step in the plan, executed by an agent.
type Step struct {
	ID                string              `json:"id"`
	Sequence          int                 `json:"sequence"`
	Name              string              `json:"name"`
	Status            Status              `json:"status"`
	Error             string              `json:"error,omitempty"`
	Agent             Agent               `json:"agent"`
	AvailableToolRefs []MCPToolRef        `json:"availableToolRefs"`
	InputContext      Context             `json:"inputContext"`
	OutputContext     Context             `json:"outputContext"`
	Query             Query               `json:"query"`
	Answer            Answer              `json:"answer"`
	Result            Context             `json:"result"` // I might not need this but for an intermediate step
	Artifacts         []Artifact          `json:"artifacts,omitempty"`
	CuratedTools      []McpToolInvocation `json:"actions,omitempty"`
	StartedAt         time.Time           `json:"startedAt"`
	FinishedAt        time.Time           `json:"finishedAt"`
	InputStepID       string              `json:"inputStepId,omitempty"`
	Version           int                 `json:"version"`
}

type EdgeType string

const (
	DependsOn  EdgeType = "depends_on"
	Triggers   EdgeType = "triggers"
	DependsAll EdgeType = "depends_on_all"
	DependsAny EdgeType = "depends_on_any"
)

// Edge describes the dependency between two steps. And how they are to be run
type Edge struct {
	FromStepId string   `json:"fromStepId"`
	ToStepId   string   `json:"toStepId"`
	Type       EdgeType `json:"type"`
}

type ExecutionNode struct {
	StepId string `json:"stepId"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

// ExecutionGraph describes the execution graph of the workflow. This holds the runtime details.
type ExecutionGraph struct {
	ID      string          `json:"id"`
	Nodes   []ExecutionNode `json:"nodes"`
	Edges   []Edge          `json:"edges"`
	Version int             `json:"version"`
}

// Workflow describes the full execution structure and participating agents.
type Workflow struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	AgentRefs        []AgentRef        `json:"agentRefs"`
	AvailableMcpRefs []string          `json:"availableMcpRefs"`
	Variables        map[string]string `json:"variables,omitempty"`
	Mode             string            `json:"mode"`
	ExecutionGraph   *ExecutionGraph   `json:"executionGraph"`
}

// Interaction represents a single user ↔ AI exchange (atomic Q&A cycle).
// Causal relationship: BaseQuery, BaseContext -> Plan -> Workflow -> Steps -> Messages, Artifacts -> Answer
type Interaction struct {
	ID          string     `json:"id"`
	BaseQuery   *Query     `json:"baseQuery"`
	BaseContext *Context   `json:"baseContext"`
	Plan        *Plan      `json:"plan"`
	Workflow    *Workflow  `json:"workflow"`
	CurrentStep string     `json:"currentStep"`
	StepIds     []string   `json:"stepIds"`
	Messages    []Message  `json:"messages"`
	Summary     string     `json:"summary,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}
