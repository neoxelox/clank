package suggestion

import (
	"encoding/json"
	"time"

	"github.com/pgvector/pgvector-go"
)

const (
	SUGGESTION_FEEDBACK_MODEL_TABLE = "\"suggestion_feedback\""
	SUGGESTION_MODEL_TABLE          = "\"suggestion\""
)

type SuggestionFeedbackModel struct {
	SuggestionID string `db:"suggestion_id"`
	FeedbackID   string `db:"feedback_id"`
}

type SuggestionModel struct {
	ID               string          `db:"id"`
	ProductID        string          `db:"product_id"`
	Embedding        pgvector.Vector `db:"embedding"`
	Sources          []byte          `db:"sources"`
	Title            string          `db:"title"`
	Description      string          `db:"description"`
	Reason           string          `db:"reason"`
	Importances      []byte          `db:"importances"`
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

func NewSuggestionModel(suggestion Suggestion) *SuggestionModel {
	embedding := pgvector.NewVector(suggestion.Embedding)

	sources, err := json.Marshal(suggestion.Sources)
	if err != nil {
		panic(err)
	}

	importances, err := json.Marshal(suggestion.Importances)
	if err != nil {
		panic(err)
	}

	categories, err := json.Marshal(suggestion.Categories)
	if err != nil {
		panic(err)
	}

	releases, err := json.Marshal(suggestion.Releases)
	if err != nil {
		panic(err)
	}

	return &SuggestionModel{
		ID:               suggestion.ID,
		ProductID:        suggestion.ProductID,
		Embedding:        embedding,
		Sources:          sources,
		Title:            suggestion.Title,
		Description:      suggestion.Description,
		Reason:           suggestion.Reason,
		Importances:      importances,
		Priority:         suggestion.Priority,
		Categories:       categories,
		Releases:         releases,
		Customers:        suggestion.Customers,
		AssigneeID:       suggestion.AssigneeID,
		Quality:          suggestion.Quality,
		FirstSeenAt:      suggestion.FirstSeenAt,
		LastSeenAt:       suggestion.LastSeenAt,
		CreatedAt:        suggestion.CreatedAt,
		ArchivedAt:       suggestion.ArchivedAt,
		LastAggregatedAt: suggestion.LastAggregatedAt,
		ExportedAt:       suggestion.ExportedAt,
	}
}

func (self *SuggestionModel) ToEntity() *Suggestion {
	embedding := self.Embedding.Slice()

	var sources map[string]int
	err := json.Unmarshal(self.Sources, &sources)
	if err != nil {
		panic(err)
	}

	var importances map[string]int
	err = json.Unmarshal(self.Importances, &importances)
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

	return &Suggestion{
		ID:               self.ID,
		ProductID:        self.ProductID,
		Embedding:        embedding,
		Sources:          sources,
		Title:            self.Title,
		Description:      self.Description,
		Reason:           self.Reason,
		Importances:      importances,
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
	PARTIAL_SUGGESTION_MODEL_TABLE = "\"partial_suggestion\""
)

type PartialSuggestionModel struct {
	ID          string    `db:"id"`
	FeedbackID  string    `db:"feedback_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Reason      string    `db:"reason"`
	Importance  string    `db:"importance"`
	Category    string    `db:"category"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewPartialSuggestionModel(partial PartialSuggestion) *PartialSuggestionModel {
	return &PartialSuggestionModel{
		ID:          partial.ID,
		FeedbackID:  partial.FeedbackID,
		Title:       partial.Title,
		Description: partial.Description,
		Reason:      partial.Reason,
		Importance:  partial.Importance,
		Category:    partial.Category,
		CreatedAt:   partial.CreatedAt,
	}
}

func (self *PartialSuggestionModel) ToEntity() *PartialSuggestion {
	return &PartialSuggestion{
		ID:          self.ID,
		FeedbackID:  self.FeedbackID,
		Title:       self.Title,
		Description: self.Description,
		Reason:      self.Reason,
		Importance:  self.Importance,
		Category:    self.Category,
		CreatedAt:   self.CreatedAt,
	}
}
