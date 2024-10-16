package issue

import "time"

type IssuePayload struct {
	ID          string         `json:"id"`
	ProductID   string         `json:"product_id"`
	Sources     map[string]int `json:"sources"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Steps       []string       `json:"steps"`
	Severities  map[string]int `json:"severities"`
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

func NewIssuePayload(issue Issue) *IssuePayload {
	return &IssuePayload{
		ID:          issue.ID,
		ProductID:   issue.ProductID,
		Sources:     issue.Sources,
		Title:       issue.Title,
		Description: issue.Description,
		Steps:       issue.Steps,
		Severities:  issue.Severities,
		Priority:    issue.Priority,
		Categories:  issue.Categories,
		Releases:    issue.Releases,
		Customers:   issue.Customers,
		AssigneeID:  issue.AssigneeID,
		Quality:     issue.Quality,
		FirstSeenAt: issue.FirstSeenAt,
		LastSeenAt:  issue.LastSeenAt,
		ArchivedAt:  issue.ArchivedAt,
	}
}
