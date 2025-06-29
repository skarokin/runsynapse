package types

import (
	"encoding/json"
	"fmt"
)

// userID validations
type UserID string

func (u *UserID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if len(s) < 1 || len(s) > 36 {
		return fmt.Errorf("invalid user_id: must be between 1 and 36 characters")
	}
	*u = UserID(s)
	return nil
}

type ThoughtID string

func (t *ThoughtID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if len(s) < 1 || len(s) > 36 {
		return fmt.Errorf("invalid thought_id: must be between 1 and 36 characters")
	}
	*t = ThoughtID(s)
	return nil
}

// query validations
type Query string

func (q *Query) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if len(s) < 1 || len(s) > 256 {
		return fmt.Errorf("invalid query: must be between 1 and 256 characters")
	}
	*q = Query(s)
	return nil
}