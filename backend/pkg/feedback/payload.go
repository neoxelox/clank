package feedback

import "time"

type FeedbackPayloadCustomer struct {
	Email    *string `json:"email"`
	Name     string  `json:"name"`
	Picture  string  `json:"picture"`
	Location *string `json:"location"`
	Verified *bool   `json:"verified"`
	Reviews  *int    `json:"reviews"`
	Link     *string `json:"link"`
}

type FeedbackPayloadMetadata struct {
	Rating   *float64  `json:"rating"`
	Media    *[]string `json:"media"`
	Verified *bool     `json:"verified"`
	Votes    *int      `json:"votes"`
	Link     *string   `json:"link"`
}

type FeedbackPayload struct {
	ID          string                  `json:"id"`
	ProductID   string                  `json:"product_id"`
	Source      string                  `json:"source"`
	Customer    FeedbackPayloadCustomer `json:"customer"`
	Content     string                  `json:"content"`
	Language    string                  `json:"language"`
	Translation string                  `json:"translation"`
	Release     string                  `json:"release"`
	Metadata    FeedbackPayloadMetadata `json:"metadata"`
	PostedAt    time.Time               `json:"posted_at"`
}

func NewFeedbackPayload(feedback Feedback) *FeedbackPayload {
	return &FeedbackPayload{
		ID:        feedback.ID,
		ProductID: feedback.ProductID,
		Source:    feedback.Source,
		Customer: FeedbackPayloadCustomer{
			Email:    feedback.Customer.Email,
			Name:     feedback.Customer.Name,
			Picture:  feedback.Customer.Picture,
			Location: feedback.Customer.Location,
			Verified: feedback.Customer.Verified,
			Reviews:  feedback.Customer.Reviews,
			Link:     feedback.Customer.Link,
		},
		Content:     feedback.Content,
		Language:    feedback.Language,
		Translation: feedback.Translation,
		Release:     feedback.Release,
		Metadata: FeedbackPayloadMetadata{
			Rating:   feedback.Metadata.Rating,
			Media:    feedback.Metadata.Media,
			Verified: feedback.Metadata.Verified,
			Votes:    feedback.Metadata.Votes,
			Link:     feedback.Metadata.Link,
		},
		PostedAt: feedback.PostedAt,
	}
}
