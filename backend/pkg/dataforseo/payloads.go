package dataforseo

type trustpilotResponseTaskResult struct {
	Domain   string `json:"domain"`
	Type     string `json:"type"`
	SEDomain string `json:"se_domain"`
	Location string `json:"location"`
	CheckURL string `json:"check_url"`
	Datetime string `json:"datetime"`
	Title    string `json:"title"`
	Rating   struct {
		RatingType string  `json:"rating_type"`
		Value      float64 `json:"value"`
		VotesCount int     `json:"votes_count"`
		RatingMax  int     `json:"rating_max"`
	} `json:"rating"`
	ReviewsCount int `json:"reviews_count"`
	ItemsCount   int `json:"items_count"`
	Items        []struct {
		Type         string  `json:"type"`
		RankGroup    int     `json:"rank_group"`
		RankAbsolute int     `json:"rank_absolute"`
		Position     string  `json:"position"`
		URL          *string `json:"url"`
		Rating       struct {
			RatingType string  `json:"rating_type"`
			Value      float64 `json:"value"`
			VotesCount int     `json:"votes_count"`
			RatingMax  int     `json:"rating_max"`
		} `json:"rating"`
		Verified     bool     `json:"verified"`
		Language     string   `json:"language"`
		Timestamp    string   `json:"timestamp"`
		Title        string   `json:"title"`
		ReviewText   string   `json:"review_text"`
		ReviewImages []string `json:"review_images"`
		UserProfile  struct {
			Name         string  `json:"name"`
			URL          string  `json:"url"`
			ImageURL     *string `json:"image_url"`
			Location     string  `json:"location"`
			ReviewsCount int     `json:"reviews_count"`
		} `json:"user_profile"`
		Responses []struct {
			Title     string `json:"title"`
			Text      string `json:"text"`
			Timestamp string `json:"timestamp"`
		} `json:"responses"`
	} `json:"items"`
}

type playStoreResponseTaskResult struct {
	AppID        string `json:"app_id"`
	Type         string `json:"type"`
	SEDomain     string `json:"se_domain"`
	LocationCode int    `json:"location_code"`
	LanguageCode string `json:"language_code"`
	CheckURL     string `json:"check_url"`
	Datetime     string `json:"datetime"`
	Title        string `json:"title"`
	Rating       struct {
		RatingType string  `json:"rating_type"`
		Value      float64 `json:"value"`
		VotesCount int     `json:"votes_count"`
		RatingMax  int     `json:"rating_max"`
	} `json:"rating"`
	ReviewsCount int `json:"reviews_count"`
	ItemsCount   int `json:"items_count"`
	Items        []struct {
		Type         string  `json:"type"`
		RankGroup    int     `json:"rank_group"`
		RankAbsolute int     `json:"rank_absolute"`
		Position     string  `json:"position"`
		Version      *string `json:"version"`
		Rating       struct {
			RatingType string  `json:"rating_type"`
			Value      float64 `json:"value"`
			VotesCount int     `json:"votes_count"`
			RatingMax  int     `json:"rating_max"`
		} `json:"rating"`
		Timestamp    string `json:"timestamp"`
		ID           string `json:"id"`
		HelpfulCount int    `json:"helpful_count"`
		Title        string `json:"title"`
		ReviewText   string `json:"review_text"`
		UserProfile  struct {
			ProfileName     string `json:"profile_name"`
			ProfileImageURL string `json:"profile_image_url"`
		} `json:"user_profile"`
		Responses []struct {
			Author    string `json:"author"`
			Title     string `json:"title"`
			Text      string `json:"text"`
			Timestamp string `json:"timestamp"`
		} `json:"responses"`
	} `json:"items"`
}

type appStoreResponseTaskResult struct {
	AppID        string `json:"app_id"`
	Type         string `json:"type"`
	SEDomain     string `json:"se_domain"`
	LocationCode int    `json:"location_code"`
	LanguageCode string `json:"language_code"`
	CheckURL     string `json:"check_url"`
	Datetime     string `json:"datetime"`
	Title        string `json:"title"`
	Rating       struct {
		RatingType string  `json:"rating_type"`
		Value      float64 `json:"value"`
		VotesCount int     `json:"votes_count"`
		RatingMax  int     `json:"rating_max"`
	} `json:"rating"`
	ReviewsCount int `json:"reviews_count"`
	ItemsCount   int `json:"items_count"`
	Items        []struct {
		Type         string  `json:"type"`
		RankGroup    int     `json:"rank_group"`
		RankAbsolute int     `json:"rank_absolute"`
		Position     string  `json:"position"`
		Version      *string `json:"version"`
		Rating       struct {
			RatingType string  `json:"rating_type"`
			Value      float64 `json:"value"`
			VotesCount int     `json:"votes_count"`
			RatingMax  int     `json:"rating_max"`
		} `json:"rating"`
		Timestamp   string `json:"timestamp"`
		ID          string `json:"id"`
		Title       string `json:"title"`
		ReviewText  string `json:"review_text"`
		UserProfile struct {
			ProfileName     string  `json:"profile_name"`
			ProfileImageURL *string `json:"profile_image_url"`
		} `json:"user_profile"`
	} `json:"items"`
}

type amazonResponseTaskResult struct {
	ASIN         string         `json:"asin"`
	Type         string         `json:"type"`
	SEDomain     string         `json:"se_domain"`
	LocationCode int            `json:"location_code"`
	LanguageCode string         `json:"language_code"`
	CheckURL     string         `json:"check_url"`
	Datetime     string         `json:"datetime"`
	Spell        map[string]any `json:"spell"`
	Title        string         `json:"title"`
	Image        struct {
		Type     string `json:"type"`
		Alt      string `json:"alt"`
		URL      string `json:"url"`
		ImageURL string `json:"image_url"`
	} `json:"image"`
	Rating struct {
		Type       string  `json:"type"`
		Position   string  `json:"position"`
		RatingType string  `json:"rating_type"`
		Value      float64 `json:"value"`
		VotesCount int     `json:"votes_count"`
		RatingMax  int     `json:"rating_max"`
	} `json:"rating"`
	ReviewsCount int      `json:"reviews_count"`
	ItemTypes    []string `json:"item_types"`
	ItemsCount   int      `json:"items_count"`
	Items        []struct {
		Type         string `json:"type"`
		RankGroup    int    `json:"rank_group"`
		RankAbsolute int    `json:"rank_absolute"`
		Position     string `json:"position"`
		XPath        string `json:"xpath"` // nolint:tagliatelle
		Verified     bool   `json:"verified"`
		Subtitle     string `json:"subtitle"`
		HelpfulVotes int    `json:"helpful_votes"`
		Images       []struct {
			Type     string `json:"type"`
			Alt      string `json:"alt"`
			URL      string `json:"url"`
			ImageURL string `json:"image_url"`
		} `json:"images"`
		Videos []struct {
			Type    string `json:"type"`
			Source  string `json:"source"`
			Preview string `json:"preview"`
		} `json:"videos"`
		UserProfile struct {
			Name         string  `json:"name"`
			Avatar       string  `json:"avatar"`
			URL          string  `json:"url"`
			ReviewsCount *int    `json:"reviews_count"`
			Locations    *string `json:"locations"`
		} `json:"user_profile"`
		Title           string  `json:"title"`
		URL             *string `json:"url"`
		ReviewText      string  `json:"review_text"`
		PublicationDate string  `json:"publication_date"`
		Rating          struct {
			RatingType string  `json:"rating_type"`
			Value      float64 `json:"value"`
			VotesCount int     `json:"votes_count"`
			RatingMax  int     `json:"rating_max"`
		} `json:"rating"`
	} `json:"items"`
}

type responseTaskResult interface {
	trustpilotResponseTaskResult | playStoreResponseTaskResult |
		appStoreResponseTaskResult | amazonResponseTaskResult
}

type responseTaskData struct {
	Tag string `json:"tag"`
}

type responseTask[R responseTaskResult] struct {
	ID            string           `json:"id"`
	StatusCode    int              `json:"status_code"`
	StatusMessage string           `json:"status_message"`
	Time          string           `json:"time"`
	Cost          float64          `json:"cost"`
	ResultCount   int              `json:"result_count"`
	Path          []string         `json:"path"`
	Data          responseTaskData `json:"data"`
	Result        []R              `json:"result"`
}

func isResponseStatusCodeOk(statusCode int) bool {
	return statusCode >= 20000 && statusCode <= 29999
}

func isResponseStatusCodeInvalidRequest(statusCode int) bool {
	return (statusCode >= 40000 && statusCode <= 49999) && statusCode != 40101
}

func isResponseStatusCodeInternalError(statusCode int) bool {
	return (statusCode >= 50000 && statusCode <= 59999) || statusCode == 40101
}

type response[R responseTaskResult] struct {
	Version       string            `json:"version"`
	StatusCode    int               `json:"status_code"`
	StatusMessage string            `json:"status_message"`
	Time          string            `json:"time"`
	Cost          float64           `json:"cost"`
	TasksCount    int               `json:"tasks_count"`
	TasksError    int               `json:"tasks_error"`
	Tasks         []responseTask[R] `json:"tasks"`
}
