package feedback

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/neoxelox/kit/util"
)

const (
	FEEDBACK_CUSTOMER_DEFAULT_PICTURE = "https://clank.so/images/pictures/feedback-customer.png"
)

const (
	FeedbackSourceTrustpilot = "TRUSTPILOT"
	FeedbackSourcePlayStore  = "PLAY_STORE"
	FeedbackSourceAppStore   = "APP_STORE"
	FeedbackSourceAmazon     = "AMAZON"
	FeedbackSourceIAgora     = "IAGORA"
	FeedbackSourceWebhook    = "WEBHOOK"
	FeedbackSourceWidget     = "WIDGET"
)

func IsFeedbackSource(value string) bool {
	return value == FeedbackSourceTrustpilot ||
		value == FeedbackSourcePlayStore ||
		value == FeedbackSourceAppStore ||
		value == FeedbackSourceAmazon ||
		value == FeedbackSourceIAgora ||
		value == FeedbackSourceWebhook ||
		value == FeedbackSourceWidget
}

type FeedbackCustomer struct {
	Email    *string
	Name     string
	Picture  string
	Location *string
	Verified *bool
	Reviews  *int
	Link     *string
}

type FeedbackMetadata struct {
	Rating   *float64
	Media    *[]string
	Verified *bool
	Votes    *int
	Link     *string
}

type Feedback struct {
	ID           string
	ProductID    string
	Hash         string
	Source       string
	Customer     FeedbackCustomer
	Content      string
	Language     string
	Translation  string
	Release      string
	Metadata     FeedbackMetadata
	Tokens       int
	PostedAt     time.Time
	CollectedAt  time.Time
	TranslatedAt *time.Time
	ProcessedAt  *time.Time
}

func NewFeedback() *Feedback {
	return &Feedback{}
}

func (self Feedback) String() string {
	return fmt.Sprintf("<Feedback: %s (%s)>", self.Customer.Name, self.ID)
}

func (self Feedback) Equals(other Feedback) bool {
	return util.Equals(self, other)
}

func (self Feedback) Copy() *Feedback {
	return util.Copy(self)
}

const (
	// nolint: lll, revive
	PUNCTUATION_MARKS = "!\"#%&'()*,./:;?@[\\]^_`{|}~\xA0¡¦§¨©ª«¬\xAD®¯²³´µ¶·¸¹º»¿‐‑‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞™"
)

func CleanContent(title string, body string) string {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)

	if len(title) == 0 {
		return body
	}

	if len(body) == 0 {
		return title
	}

	if strings.HasPrefix(body, strings.TrimRight(title, PUNCTUATION_MARKS)) { // nolint:staticcheck
		return body
	}

	return title + "\n" + body
}

func ComputeHash(source string, customer string, content string) string {
	hash := sha1.Sum([]byte(source + customer + content)) // nolint:gosec
	return hex.EncodeToString(hash[:])
}
