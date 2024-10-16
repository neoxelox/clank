package engine

const (
	OPTION_UNKNOWN = "UNKNOWN"
)

type Usage struct {
	Input  int
	Output int
}

type Feedback struct {
	Content string
}

type Issue struct {
	Title       string
	Description string
	Steps       []string
	Severity    string
	Category    string
}

type Suggestion struct {
	Title       string
	Description string
	Reason      string
	Importance  string
	Category    string
}

type Review struct {
	Content   string
	Keywords  []string
	Sentiment string
	Emotions  []string
	Intention string
	Category  string
}
