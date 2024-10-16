package suggestion

import "time"

type SuggestionPayload struct {
	ID          string         `json:"id"`
	ProductID   string         `json:"product_id"`
	Sources     map[string]int `json:"sources"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Reason      string         `json:"reason"`
	Importances map[string]int `json:"importances"`
	Priority    int            `json:"priority"`
	Categories  map[string]int `json:"categories"`
	Releases    map[string]int `json:"releases"`
	Customers   int            `json:"customers"`
	AssigneeID  *string        `json:"assignee_id"`
	Quality     *int           `json:"quality"`
	FirstSeenAt time.Time      `json:"first_seen_at"`
	LastSeenAt  time.Time      `json:"last_seen_at"`
	ArchivedAt  *time.Time     `json:"archived_at"`
}

func NewSuggestionPayload(suggestion Suggestion) *SuggestionPayload {
	return &SuggestionPayload{
		ID:          suggestion.ID,
		ProductID:   suggestion.ProductID,
		Sources:     suggestion.Sources,
		Title:       suggestion.Title,
		Description: suggestion.Description,
		Reason:      suggestion.Reason,
		Importances: suggestion.Importances,
		Priority:    suggestion.Priority,
		Categories:  suggestion.Categories,
		Releases:    suggestion.Releases,
		Customers:   suggestion.Customers,
		AssigneeID:  suggestion.AssigneeID,
		Quality:     suggestion.Quality,
		FirstSeenAt: suggestion.FirstSeenAt,
		LastSeenAt:  suggestion.LastSeenAt,
		ArchivedAt:  suggestion.ArchivedAt,
	}
}
