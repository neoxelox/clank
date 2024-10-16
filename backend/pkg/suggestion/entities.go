package suggestion

import (
	"fmt"
	"time"

	"backend/pkg/engine"
	"backend/pkg/util"

	kitUtil "github.com/neoxelox/kit/util"
)

const (
	SUGGESTION_MIN_QUALITY       = 0
	SUGGESTION_MAX_QUALITY       = 5
	SUGGESTION_SIMILAR_THRESHOLD = 0.6
)

const (
	SuggestionImportanceCritical = "CRITICAL"
	SuggestionImportanceHigh     = "HIGH"
	SuggestionImportanceMedium   = "MEDIUM"
	SuggestionImportanceLow      = "LOW"
)

func IsSuggestionImportance(value string) bool {
	return value == SuggestionImportanceCritical ||
		value == SuggestionImportanceHigh ||
		value == SuggestionImportanceMedium ||
		value == SuggestionImportanceLow
}

var SuggestionImportanceWeight = map[string]int{
	SuggestionImportanceCritical: 4,
	SuggestionImportanceHigh:     3,
	SuggestionImportanceMedium:   2,
	SuggestionImportanceLow:      1,
}

type Suggestion struct {
	ID               string
	ProductID        string
	Embedding        []float32
	Sources          map[string]int
	Title            string
	Description      string
	Reason           string
	Importances      map[string]int
	Priority         int
	Categories       map[string]int
	Releases         map[string]int
	Customers        int
	AssigneeID       *string
	Quality          *int
	FirstSeenAt      time.Time
	LastSeenAt       time.Time
	CreatedAt        time.Time
	ArchivedAt       *time.Time
	LastAggregatedAt *time.Time
	ExportedAt       *time.Time
}

func NewSuggestion() *Suggestion {
	return &Suggestion{}
}

func (self Suggestion) Importance() string {
	return computeImportance(self.Importances)
}

func (self Suggestion) Category() string {
	return computeCategory(self.Categories)
}

func (self Suggestion) String() string {
	return fmt.Sprintf("<Suggestion: %s (%s)>", self.Title, self.ID)
}

func (self Suggestion) Equals(other Suggestion) bool {
	return kitUtil.Equals(self, other)
}

func (self Suggestion) Copy() *Suggestion {
	return kitUtil.Copy(self)
}

func computeImportance(importances map[string]int) string {
	majorImportance := SuggestionImportanceLow
	maxCount := 0
	for importance, count := range importances {
		if count > maxCount || (count == maxCount && importance > majorImportance) {
			majorImportance = importance
			maxCount = count
		}
	}

	return majorImportance
}

func ComputePriority(importances map[string]int, customers int) int {
	return SuggestionImportanceWeight[computeImportance(importances)] * customers
}
func computeCategory(categories map[string]int) string {
	majorCategory := engine.OPTION_UNKNOWN
	maxCount := 0
	for category, count := range categories {
		if count > maxCount || (count == maxCount && category > majorCategory) {
			majorCategory = category
			maxCount = count
		}
	}

	return majorCategory
}

const (
	SuggestionSearchFiltersAssigneesUnassigned = "UNASSIGNED"
)

const (
	SuggestionSearchFiltersStatusActive     = "ACTIVE"
	SuggestionSearchFiltersStatusRegressed  = "REGRESSED"
	SuggestionSearchFiltersStatusArchived   = "ARCHIVED"
	SuggestionSearchFiltersStatusUnarchived = "UNARCHIVED"
)

func IsSuggestionSearchFiltersStatus(value string) bool {
	return value == SuggestionSearchFiltersStatusActive ||
		value == SuggestionSearchFiltersStatusRegressed ||
		value == SuggestionSearchFiltersStatusArchived ||
		value == SuggestionSearchFiltersStatusUnarchived
}

type SuggestionSearchFilters struct {
	Embedding        *[]float32
	Sources          *[]string
	Importances      *[]string
	Releases         *[]string
	Categories       *[]string
	Assignees        *[]string
	Status           *string
	FirstSeenStartAt *time.Time
	FirstSeenEndAt   *time.Time
	LastSeenStartAt  *time.Time
	LastSeenEndAt    *time.Time
}

const (
	SuggestionSearchOrdersAscending  = "ASC"
	SuggestionSearchOrdersDescending = "DESC"
)

type SuggestionSearchOrders struct {
	Relevance string
}

type SuggestionSearchCursor struct {
	Priority int
	SeenAt   time.Time
}

type SuggestionSearch struct {
	Filters    SuggestionSearchFilters
	Orders     SuggestionSearchOrders
	Pagination util.Pagination[SuggestionSearchCursor]
}

type PartialSuggestion struct {
	ID          string
	FeedbackID  string
	Title       string
	Description string
	Reason      string
	Importance  string
	Category    string
	CreatedAt   time.Time
}

func NewPartialSuggestion() *PartialSuggestion {
	return &PartialSuggestion{}
}

func (self PartialSuggestion) String() string {
	return fmt.Sprintf("<PartialSuggestion: %s (%s)>", self.Title, self.ID)
}

func (self PartialSuggestion) Equals(other PartialSuggestion) bool {
	return kitUtil.Equals(self, other)
}

func (self PartialSuggestion) Copy() *PartialSuggestion {
	return kitUtil.Copy(self)
}
