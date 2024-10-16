package issue

import (
	"fmt"
	"time"

	"backend/pkg/engine"
	"backend/pkg/util"

	kitUtil "github.com/neoxelox/kit/util"
)

const (
	ISSUE_MIN_QUALITY       = 0
	ISSUE_MAX_QUALITY       = 5
	ISSUE_SIMILAR_THRESHOLD = 0.6
)

const (
	IssueSeverityCritical = "CRITICAL"
	IssueSeverityHigh     = "HIGH"
	IssueSeverityMedium   = "MEDIUM"
	IssueSeverityLow      = "LOW"
)

func IsIssueSeverity(value string) bool {
	return value == IssueSeverityCritical ||
		value == IssueSeverityHigh ||
		value == IssueSeverityMedium ||
		value == IssueSeverityLow
}

var IssueSeverityWeight = map[string]int{
	IssueSeverityCritical: 4,
	IssueSeverityHigh:     3,
	IssueSeverityMedium:   2,
	IssueSeverityLow:      1,
}

type Issue struct {
	ID               string
	ProductID        string
	Embedding        []float32
	Sources          map[string]int
	Title            string
	Description      string
	Steps            []string
	Severities       map[string]int
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

func NewIssue() *Issue {
	return &Issue{}
}

func (self Issue) Severity() string {
	return computeSeverity(self.Severities)
}

func (self Issue) Category() string {
	return computeCategory(self.Categories)
}

func (self Issue) String() string {
	return fmt.Sprintf("<Issue: %s (%s)>", self.Title, self.ID)
}

func (self Issue) Equals(other Issue) bool {
	return kitUtil.Equals(self, other)
}

func (self Issue) Copy() *Issue {
	return kitUtil.Copy(self)
}

func computeSeverity(severities map[string]int) string {
	majorSeverity := IssueSeverityLow
	maxCount := 0
	for severity, count := range severities {
		if count > maxCount || (count == maxCount && severity > majorSeverity) {
			majorSeverity = severity
			maxCount = count
		}
	}

	return majorSeverity
}

func ComputePriority(severities map[string]int, customers int) int {
	return IssueSeverityWeight[computeSeverity(severities)] * customers
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
	IssueSearchFiltersAssigneesUnassigned = "UNASSIGNED"
)

const (
	IssueSearchFiltersStatusActive     = "ACTIVE"
	IssueSearchFiltersStatusRegressed  = "REGRESSED"
	IssueSearchFiltersStatusArchived   = "ARCHIVED"
	IssueSearchFiltersStatusUnarchived = "UNARCHIVED"
)

func IsIssueSearchFiltersStatus(value string) bool {
	return value == IssueSearchFiltersStatusActive ||
		value == IssueSearchFiltersStatusRegressed ||
		value == IssueSearchFiltersStatusArchived ||
		value == IssueSearchFiltersStatusUnarchived
}

type IssueSearchFilters struct {
	Embedding        *[]float32
	Sources          *[]string
	Severities       *[]string
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
	IssueSearchOrdersAscending  = "ASC"
	IssueSearchOrdersDescending = "DESC"
)

type IssueSearchOrders struct {
	Relevance string
}

type IssueSearchCursor struct {
	Priority int
	SeenAt   time.Time
}

type IssueSearch struct {
	Filters    IssueSearchFilters
	Orders     IssueSearchOrders
	Pagination util.Pagination[IssueSearchCursor]
}

type PartialIssue struct {
	ID          string
	FeedbackID  string
	Title       string
	Description string
	Steps       []string
	Severity    string
	Category    string
	CreatedAt   time.Time
}

func NewPartialIssue() *PartialIssue {
	return &PartialIssue{}
}

func (self PartialIssue) String() string {
	return fmt.Sprintf("<PartialIssue: %s (%s)>", self.Title, self.ID)
}

func (self PartialIssue) Equals(other PartialIssue) bool {
	return kitUtil.Equals(self, other)
}

func (self PartialIssue) Copy() *PartialIssue {
	return kitUtil.Copy(self)
}
