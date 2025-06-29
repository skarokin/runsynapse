package types

import (
	"github.com/google/uuid"
)

type Thought struct {
	ID        uuid.UUID  `json:"id"`
	Thought   string     `json:"thought"`
	Pinned    bool       `json:"pinned"`
	Created   string     `json:"created_at"`
	Attachments []string `json:"attachments,omitempty"`
}

type LoadFunctionResponse struct {
	Thoughts    []Thought `json:"thoughts,omitempty"`
	PinnedThoughts []Thought `json:"pinned_thoughts,omitempty"`
	HasMoreAbove bool      `json:"more_above"`
}

type LoadThoughtsResponse struct {
	Thoughts    []Thought `json:"thoughts,omitempty"`
	HasMoreAbove bool      `json:"more_above"`
}

// for optimistic UI, returns the new thought w/ its attachments (we do this because thoughtID and attachment url generated in backend)
type NewThoughtResponse struct {
	Thought    Thought    `json:"thought"`
}

type DeleteThoughtResponse struct {
	Success bool   `json:"success"`
}

// pin, unpin, or get pinned
type PinThoughtResponse struct {
	PinnedThoughts []Thought `json:"pinned_thoughts"`
}

// mirrors LoadThoughts by returning thoughts around a pin
type GotoPinResponse struct {
	Thoughts []Thought `json:"thoughts"`
	HasMoreAbove bool    `json:"more_above"`
	HasMoreBelow bool    `json:"more_below"`
}