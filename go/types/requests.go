package types

// NewThoughtRequest works entirely with form data, so no struct needed

type RequestsOnlyRequiringUserID struct {
	UserID UserID `json:"user_id"`
}

type LoadThoughtsRequest struct {
	UserID UserID `json:"user_id"`
	Cursor ThoughtID `json:"cursor"`
}

type TogglePinRequest struct {
	UserID   UserID `json:"user_id"`
	ThoughtID string `json:"thought_id"`
}

type DeleteThoughtRequest struct {
	UserID   UserID `json:"user_id"`
	ThoughtID ThoughtID `json:"thought_id"`
}

type SearchThoughtsRequest struct {
	UserID UserID `json:"user_id"`
	Query  Query `json:"query"`
}