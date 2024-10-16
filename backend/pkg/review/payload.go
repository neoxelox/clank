package review

import (
	"backend/pkg/feedback"
)

type ReviewPayload struct {
	ID        string                   `json:"id"`
	ProductID string                   `json:"product_id"`
	Feedback  feedback.FeedbackPayload `json:"feedback"`
	Keywords  []string                 `json:"keywords"`
	Sentiment string                   `json:"sentiment"`
	Emotions  []string                 `json:"emotions"`
	Intention string                   `json:"intention"`
	Category  string                   `json:"category"`
	Quality   *int                     `json:"quality"`
}

func NewReviewPayload(review Review) *ReviewPayload {
	feedback := feedback.NewFeedbackPayload(review.Feedback)

	return &ReviewPayload{
		ID:        review.ID,
		ProductID: review.ProductID,
		Feedback:  *feedback,
		Keywords:  review.Keywords,
		Sentiment: review.Sentiment,
		Emotions:  review.Emotions,
		Intention: review.Intention,
		Category:  review.Category,
		Quality:   review.Quality,
	}
}
