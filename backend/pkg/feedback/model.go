package feedback

import (
	"encoding/json"
	"time"
)

const (
	FEEDBACK_MODEL_TABLE = "\"feedback\""
)

type FeedbackModel struct {
	ID           string     `db:"id"`
	ProductID    string     `db:"product_id"`
	Hash         string     `db:"hash"`
	Source       string     `db:"source"`
	Customer     []byte     `db:"customer"`
	Content      string     `db:"content"`
	Language     string     `db:"language"`
	Translation  string     `db:"translation"`
	Release      string     `db:"release"`
	Metadata     []byte     `db:"metadata"`
	Tokens       int        `db:"tokens"`
	PostedAt     time.Time  `db:"posted_at"`
	CollectedAt  time.Time  `db:"collected_at"`
	TranslatedAt *time.Time `db:"translated_at"`
	ProcessedAt  *time.Time `db:"processed_at"`
}

func NewFeedbackModel(feedback Feedback) *FeedbackModel {
	customer, err := json.Marshal(feedback.Customer)
	if err != nil {
		panic(err)
	}

	metadata, err := json.Marshal(feedback.Metadata)
	if err != nil {
		panic(err)
	}

	return &FeedbackModel{
		ID:           feedback.ID,
		ProductID:    feedback.ProductID,
		Hash:         feedback.Hash,
		Source:       feedback.Source,
		Customer:     customer,
		Content:      feedback.Content,
		Language:     feedback.Language,
		Translation:  feedback.Translation,
		Release:      feedback.Release,
		Metadata:     metadata,
		Tokens:       feedback.Tokens,
		PostedAt:     feedback.PostedAt,
		CollectedAt:  feedback.CollectedAt,
		TranslatedAt: feedback.TranslatedAt,
		ProcessedAt:  feedback.ProcessedAt,
	}
}

func (self *FeedbackModel) ToEntity() *Feedback {
	var customer FeedbackCustomer
	err := json.Unmarshal(self.Customer, &customer)
	if err != nil {
		panic(err)
	}

	var metadata FeedbackMetadata
	err = json.Unmarshal(self.Metadata, &metadata)
	if err != nil {
		panic(err)
	}

	return &Feedback{
		ID:           self.ID,
		ProductID:    self.ProductID,
		Hash:         self.Hash,
		Source:       self.Source,
		Customer:     customer,
		Content:      self.Content,
		Language:     self.Language,
		Translation:  self.Translation,
		Release:      self.Release,
		Metadata:     metadata,
		Tokens:       self.Tokens,
		PostedAt:     self.PostedAt,
		CollectedAt:  self.CollectedAt,
		TranslatedAt: self.TranslatedAt,
		ProcessedAt:  self.ProcessedAt,
	}
}
