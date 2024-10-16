package issue

import (
	"encoding/json"
	"time"

	"github.com/pgvector/pgvector-go"
)

const (
	ISSUE_FEEDBACK_MODEL_TABLE = "\"issue_feedback\""
	ISSUE_MODEL_TABLE          = "\"issue\""
)

type IssueFeedbackModel struct {
	IssueID    string `db:"issue_id"`
	FeedbackID string `db:"feedback_id"`
}

type IssueModel struct {
	ID               string          `db:"id"`
	ProductID        string          `db:"product_id"`
	Embedding        pgvector.Vector `db:"embedding"`
	Sources          []byte          `db:"sources"`
	Title            string          `db:"title"`
	Description      string          `db:"description"`
	Steps            []string        `db:"steps"`
	Severities       []byte          `db:"severities"`
	Priority         int             `db:"priority"`
	Categories       []byte          `db:"categories"`
	Releases         []byte          `db:"releases"`
	Customers        int             `db:"customers"`
	AssigneeID       *string         `db:"assignee_id"`
	Quality          *int            `db:"quality"`
	FirstSeenAt      time.Time       `db:"first_seen_at"`
	LastSeenAt       time.Time       `db:"last_seen_at"`
	CreatedAt        time.Time       `db:"created_at"`
	ArchivedAt       *time.Time      `db:"archived_at"`
	LastAggregatedAt *time.Time      `db:"last_aggregated_at"`
	ExportedAt       *time.Time      `db:"exported_at"`
}

func NewIssueModel(issue Issue) *IssueModel {
	embedding := pgvector.NewVector(issue.Embedding)

	sources, err := json.Marshal(issue.Sources)
	if err != nil {
		panic(err)
	}

	severities, err := json.Marshal(issue.Severities)
	if err != nil {
		panic(err)
	}

	categories, err := json.Marshal(issue.Categories)
	if err != nil {
		panic(err)
	}

	releases, err := json.Marshal(issue.Releases)
	if err != nil {
		panic(err)
	}

	return &IssueModel{
		ID:               issue.ID,
		ProductID:        issue.ProductID,
		Embedding:        embedding,
		Sources:          sources,
		Title:            issue.Title,
		Description:      issue.Description,
		Steps:            issue.Steps,
		Severities:       severities,
		Priority:         issue.Priority,
		Categories:       categories,
		Releases:         releases,
		Customers:        issue.Customers,
		AssigneeID:       issue.AssigneeID,
		Quality:          issue.Quality,
		FirstSeenAt:      issue.FirstSeenAt,
		LastSeenAt:       issue.LastSeenAt,
		CreatedAt:        issue.CreatedAt,
		ArchivedAt:       issue.ArchivedAt,
		LastAggregatedAt: issue.LastAggregatedAt,
		ExportedAt:       issue.ExportedAt,
	}
}

func (self *IssueModel) ToEntity() *Issue {
	embedding := self.Embedding.Slice()

	var sources map[string]int
	err := json.Unmarshal(self.Sources, &sources)
	if err != nil {
		panic(err)
	}

	var severities map[string]int
	err = json.Unmarshal(self.Severities, &severities)
	if err != nil {
		panic(err)
	}

	var categories map[string]int
	err = json.Unmarshal(self.Categories, &categories)
	if err != nil {
		panic(err)
	}

	var releases map[string]int
	err = json.Unmarshal(self.Releases, &releases)
	if err != nil {
		panic(err)
	}

	return &Issue{
		ID:               self.ID,
		ProductID:        self.ProductID,
		Embedding:        embedding,
		Sources:          sources,
		Title:            self.Title,
		Description:      self.Description,
		Steps:            self.Steps,
		Severities:       severities,
		Priority:         self.Priority,
		Categories:       categories,
		Releases:         releases,
		Customers:        self.Customers,
		AssigneeID:       self.AssigneeID,
		Quality:          self.Quality,
		FirstSeenAt:      self.FirstSeenAt,
		LastSeenAt:       self.LastSeenAt,
		CreatedAt:        self.CreatedAt,
		ArchivedAt:       self.ArchivedAt,
		LastAggregatedAt: self.LastAggregatedAt,
		ExportedAt:       self.ExportedAt,
	}
}

const (
	PARTIAL_ISSUE_MODEL_TABLE = "\"partial_issue\""
)

type PartialIssueModel struct {
	ID          string    `db:"id"`
	FeedbackID  string    `db:"feedback_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Steps       []string  `db:"steps"`
	Severity    string    `db:"severity"`
	Category    string    `db:"category"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewPartialIssueModel(partial PartialIssue) *PartialIssueModel {
	return &PartialIssueModel{
		ID:          partial.ID,
		FeedbackID:  partial.FeedbackID,
		Title:       partial.Title,
		Description: partial.Description,
		Steps:       partial.Steps,
		Severity:    partial.Severity,
		Category:    partial.Category,
		CreatedAt:   partial.CreatedAt,
	}
}

func (self *PartialIssueModel) ToEntity() *PartialIssue {
	return &PartialIssue{
		ID:          self.ID,
		FeedbackID:  self.FeedbackID,
		Title:       self.Title,
		Description: self.Description,
		Steps:       self.Steps,
		Severity:    self.Severity,
		Category:    self.Category,
		CreatedAt:   self.CreatedAt,
	}
}
