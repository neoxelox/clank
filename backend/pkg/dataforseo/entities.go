package dataforseo

import "time"

type TrustpilotCustomer struct {
	Link     string
	Name     string
	Picture  *string
	Location string
	Reviews  int
}

type TrustpilotReview struct {
	Page      string
	Link      *string
	Title     string
	Content   string
	Images    []string
	Customer  TrustpilotCustomer
	Rating    float64
	Votes     int
	Verified  bool
	Timestamp time.Time
}

type PlayStoreCustomer struct {
	Name    string
	Picture string
}

type PlayStoreReview struct {
	ID        string
	Page      string
	Title     string
	Content   string
	Customer  PlayStoreCustomer
	Rating    float64
	Votes     int
	Release   *string
	Timestamp time.Time
}

type AppStoreCustomer struct {
	Name    string
	Picture *string
}

type AppStoreReview struct {
	ID        string
	Page      string
	Title     string
	Content   string
	Customer  AppStoreCustomer
	Rating    float64
	Release   *string
	Timestamp time.Time
}

type AmazonCustomer struct {
	Link     string
	Name     string
	Picture  string
	Location *string
	Reviews  *int
}

type AmazonReview struct {
	Page      string
	Link      *string
	Title     string
	Subtitle  string
	Content   string
	Images    []string
	Videos    []string
	Customer  AmazonCustomer
	Rating    float64
	Votes     int
	Verified  bool
	Timestamp time.Time
}

type Customer interface {
	TrustpilotCustomer | PlayStoreCustomer |
		AppStoreCustomer | AmazonCustomer
}

type Review interface {
	TrustpilotReview | PlayStoreReview |
		AppStoreReview | AmazonReview
}

type Task[R Review] struct {
	ID         string
	Status     int
	Message    string
	Cost       float64
	Identifier string
	Reviews    []R
}

type Perspective struct {
	Location string
	Language string
}

var PlayStorePerspectives []Perspective = []Perspective{
	{
		Location: "Australia",
		Language: "English",
	},
	{
		Location: "Canada",
		Language: "English",
	},
	{
		Location: "Ireland",
		Language: "English",
	},
	{
		Location: "New Zealand",
		Language: "English",
	},
	{
		Location: "United Kingdom",
		Language: "English",
	},
	{
		Location: "United States",
		Language: "English",
	},
}

var AppStorePerspectives []Perspective = []Perspective{
	{
		Location: "Australia",
		Language: "English (Australia)",
	},
	{
		Location: "Canada",
		Language: "English",
	},
	{
		Location: "Ireland",
		Language: "English (United Kingdom)",
	},
	{
		Location: "New Zealand",
		Language: "English (Australia)",
	},
	{
		Location: "United Kingdom",
		Language: "English (United Kingdom)",
	},
	{
		Location: "United States",
		Language: "English",
	},
}

var AmazonPerspectives []Perspective = []Perspective{
	{
		Location: "Australia",
		Language: "English (Australia)",
	},
	{
		Location: "Canada",
		Language: "English (Canada)",
	},
	{
		Location: "United Kingdom",
		Language: "English (United Kingdom)",
	},
	{
		Location: "United States",
		Language: "English (United States)",
	},
}
