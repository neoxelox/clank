package review

import (
	"fmt"
	"time"

	"backend/pkg/feedback"
	"backend/pkg/util"

	kitUtil "github.com/neoxelox/kit/util"
)

const (
	REVIEW_MIN_QUALITY = 0
	REVIEW_MAX_QUALITY = 5
)

const (
	ReviewSentimentPositive = "POSITIVE"
	ReviewSentimentNeutral  = "NEUTRAL"
	ReviewSentimentNegative = "NEGATIVE"
)

func IsReviewSentiment(value string) bool {
	return value == ReviewSentimentPositive ||
		value == ReviewSentimentNeutral ||
		value == ReviewSentimentNegative
}

const (
	ReviewEmotionTrust        = "TRUST"
	ReviewEmotionAcceptance   = "ACCEPTANCE"
	ReviewEmotionFear         = "FEAR"
	ReviewEmotionApprehension = "APPREHENSION"
	ReviewEmotionSurprise     = "SURPRISE"
	ReviewEmotionDistraction  = "DISTRACTION"
	ReviewEmotionSadness      = "SADNESS"
	ReviewEmotionPensiveness  = "PENSIVENESS"
	ReviewEmotionDisgust      = "DISGUST"
	ReviewEmotionBoredom      = "BOREDOM"
	ReviewEmotionAnger        = "ANGER"
	ReviewEmotionAnnoyance    = "ANNOYANCE"
	ReviewEmotionAnticipation = "ANTICIPATION"
	ReviewEmotionInterest     = "INTEREST"
	ReviewEmotionJoy          = "JOY"
	ReviewEmotionSerenity     = "SERENITY"
)

func IsReviewEmotion(value string) bool {
	return value == ReviewEmotionTrust ||
		value == ReviewEmotionAcceptance ||
		value == ReviewEmotionFear ||
		value == ReviewEmotionApprehension ||
		value == ReviewEmotionSurprise ||
		value == ReviewEmotionDistraction ||
		value == ReviewEmotionSadness ||
		value == ReviewEmotionPensiveness ||
		value == ReviewEmotionDisgust ||
		value == ReviewEmotionBoredom ||
		value == ReviewEmotionAnger ||
		value == ReviewEmotionAnnoyance ||
		value == ReviewEmotionAnticipation ||
		value == ReviewEmotionInterest ||
		value == ReviewEmotionJoy ||
		value == ReviewEmotionSerenity
}

const (
	ReviewIntentionRetain             = "RETAIN" // BUY_AGAIN / RENEW / RECOMMEND
	ReviewIntentionChurn              = "CHURN"  // RETURN / REFUND / CANCEL / DISCOURAGE
	ReviewIntentionRetainAndRecommend = "RETAIN_AND_RECOMMEND"
	ReviewIntentionChurnAndDiscourage = "CHURN_AND_DISCOURAGE"
)

func IsReviewIntention(value string) bool {
	return value == ReviewIntentionRetain ||
		value == ReviewIntentionChurn ||
		value == ReviewIntentionRetainAndRecommend ||
		value == ReviewIntentionChurnAndDiscourage
}

type Review struct {
	ID         string
	ProductID  string
	Feedback   feedback.Feedback
	Keywords   []string
	Sentiment  string
	Emotions   []string
	Intention  string
	Category   string
	Quality    *int
	CreatedAt  time.Time
	ExportedAt *time.Time
}

func NewReview() *Review {
	return &Review{}
}

func (self Review) String() string {
	return fmt.Sprintf("<Review: %s (%s)>", self.Feedback.Content[:min(len(self.Feedback.Content), 100)], self.ID)
}

func (self Review) Equals(other Review) bool {
	return kitUtil.Equals(self, other)
}

func (self Review) Copy() *Review {
	return kitUtil.Copy(self)
}

type ReviewSearchFilters struct {
	Sources     *[]string
	Releases    *[]string
	Categories  *[]string
	Keywords    *[]string
	Sentiments  *[]string
	Emotions    *[]string
	Intentions  *[]string
	Languages   *[]string
	SeenStartAt *time.Time
	SeenEndAt   *time.Time
}

const (
	ReviewSearchOrdersAscending  = "ASC"
	ReviewSearchOrdersDescending = "DESC"
)

type ReviewSearchOrders struct {
	Recency string
}

type ReviewSearch struct {
	Filters    ReviewSearchFilters
	Orders     ReviewSearchOrders
	Pagination util.Pagination[time.Time]
}
