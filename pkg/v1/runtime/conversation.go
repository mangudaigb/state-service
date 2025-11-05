package runtime

import "time"

//
// ──────────────────────────────── CONVERSATION ────────────────────────────────
//

// Conversation groups multiple interactions (like a chat session).
type Conversation struct {
	ID           string         `json:"id"`
	Interactions []Interaction  `json:"interactions"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

//
// ──────────────────────────────── SNAPSHOT ────────────────────────────────
//

// Snapshot allows saving and replaying execution states.
type Snapshot struct {
	ID        string       `json:"id"`
	StateRef  Conversation `json:"state"`
	Timestamp time.Time    `json:"timestamp"`
}
