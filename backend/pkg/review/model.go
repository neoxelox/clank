package review

import (
	"time"

	"backend/pkg/feedback"
)

const (
	REVIEW_MODEL_TABLE = "\"review\""
)

type ReviewFeedbackModel struct {
	Review   ReviewModel            `db:"review" scan:"follow"`
	Feedback feedback.FeedbackModel `db:"feedback" scan:"notate"`
}

type ReviewModel struct {
	ID         string     `db:"id"`
	ProductID  string     `db:"product_id"`
	FeedbackID string     `db:"feedback_id"`
	Keywords   []string   `db:"keywords"`
	Sentiment  string     `db:"sentiment"`
	Emotions   []string   `db:"emotions"`
	Intention  string     `db:"intention"`
	Category   string     `db:"category"`
	Quality    *int       `db:"quality"`
	CreatedAt  time.Time  `db:"created_at"`
	ExportedAt *time.Time `db:"exported_at"`
}

func NewReviewModel(review Review) *ReviewModel {
	return &ReviewModel{
		ID:         review.ID,
		ProductID:  review.ProductID,
		FeedbackID: review.Feedback.ID,
		Keywords:   review.Keywords,
		Sentiment:  review.Sentiment,
		Emotions:   review.Emotions,
		Intention:  review.Intention,
		Category:   review.Category,
		Quality:    review.Quality,
		CreatedAt:  review.CreatedAt,
		ExportedAt: review.ExportedAt,
	}
}

func (self *ReviewModel) ToEntity(_feedback feedback.FeedbackModel) *Review {
	feedback := _feedback.ToEntity()

	return &Review{
		ID:         self.ID,
		ProductID:  self.ProductID,
		Feedback:   *feedback,
		Keywords:   self.Keywords,
		Sentiment:  self.Sentiment,
		Emotions:   self.Emotions,
		Intention:  self.Intention,
		Category:   self.Category,
		Quality:    self.Quality,
		CreatedAt:  self.CreatedAt,
		ExportedAt: self.ExportedAt,
	}
}
