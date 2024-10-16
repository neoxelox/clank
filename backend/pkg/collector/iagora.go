package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"backend/pkg/config"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/organization"
	"backend/pkg/product"
	"backend/pkg/scraper"
	"backend/pkg/translator"

	"github.com/PuerkitoBio/goquery"
	"github.com/hibiken/asynq"
	"github.com/neoxelox/errors"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"
	"github.com/rs/xid"
)

const (
	IAGORA_COLLECTOR_MAX_REVIEWS_TO_COLLECT   = 167 * 6
	IAGORA_COLLECTOR_MIN_REVIEWS_TO_COLLECT   = 6
	IAGORA_COLLECTOR_DAILY_REVIEWS_TO_COLLECT = IAGORA_COLLECTOR_MIN_REVIEWS_TO_COLLECT * 2
	IAGORA_COLLECTOR_REVIEWS_PER_PAGE         = 6
	IAGORA_COLLECTOR_BASE_URL                 = "https://www.iagora.com/studies/uni/%s::show=reviews:offset=%d?lang=en#main"
)

const (
	IAgoraCollectorCollect  = "collector:collect-iagora-reviews"
	IAgoraCollectorSchedule = "collector:schedule-iagora-reviews"
)

var (
	ErrIAgoraCollectorMalformedReview = errors.New("malformed review")
)

type IAgoraCollectorSettings struct {
	CollectorSettings
	Institution string
}

type IAgoraCollectorJobdata struct {
	CollectorJobdata
	LastCollectedAt *time.Time
}

type IAgoraCollector struct {
	config                 config.Config
	observer               *kit.Observer
	collectorRepository    *CollectorRepository
	productRepository      *product.ProductRepository
	organizationRepository organization.OrganizationRepository
	feedbackRepository     *feedback.FeedbackRepository
	enqueuer               *kit.Enqueuer
	scraper                *scraper.Scraper
}

func NewIAgoraCollector(observer *kit.Observer, collectorRepository *CollectorRepository,
	productRepository *product.ProductRepository, organizationRepository organization.OrganizationRepository,
	feedbackRepository *feedback.FeedbackRepository, enqueuer *kit.Enqueuer, scraper *scraper.Scraper,
	config config.Config) *IAgoraCollector {
	return &IAgoraCollector{
		config:                 config,
		observer:               observer,
		collectorRepository:    collectorRepository,
		productRepository:      productRepository,
		organizationRepository: organizationRepository,
		feedbackRepository:     feedbackRepository,
		enqueuer:               enqueuer,
		scraper:                scraper,
	}
}

func (self *IAgoraCollector) getCollectorProductAndOrganization(ctx context.Context,
	collectorID string) (*Collector, *product.Product, *organization.Organization, error) {
	collector, err := self.collectorRepository.GetByID(ctx, collectorID)
	if err != nil {
		return nil, nil, nil, err
	}

	if collector == nil {
		return nil, nil, nil, nil
	}

	if collector.DeletedAt != nil {
		return nil, nil, nil, nil
	}

	product, err := self.productRepository.GetByID(ctx, collector.ProductID)
	if err != nil {
		return nil, nil, nil, err
	}

	if product == nil {
		return nil, nil, nil, nil
	}

	if product.DeletedAt != nil {
		return nil, nil, nil, nil
	}

	organization, err := self.organizationRepository.GetByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, nil, nil, err
	}

	if organization == nil {
		return nil, nil, nil, nil
	}

	if organization.DeletedAt != nil {
		return nil, nil, nil, nil
	}

	return collector, product, organization, nil
}

func (self *IAgoraCollector) saveAndEnqueue(ctx context.Context, feedbacks []feedback.Feedback) (int, error) {
	newFeedbacks, err := self.feedbackRepository.BulkCreate(ctx, feedbacks)
	if err != nil {
		return 0, err
	}

	for _, feedback := range feedbacks {
		err := self.enqueuer.Enqueue(ctx, translator.FeedbackTranslatorTranslate,
			translator.FeedbackTranslatorTranslateParams{
				FeedbackID: feedback.ID,
			}, asynq.MaxRetry(2), asynq.Unique(12*time.Hour))
		if err != nil {
			self.observer.Error(ctx, err)
		}
	}

	return newFeedbacks, nil
}

type webSchema struct {
	Reviews []struct {
		Author string `json:"author"`
		Date   string `json:"datePublished"` // nolint:tagliatelle
	} `json:"review"` // nolint:tagliatelle
}

type reviewInfo struct {
	CustomerName     string
	CustomerLocation *string
	MetadataRating   *float64
	MetadataVotes    *int
	MetadataLink     string
	PostedAt         time.Time
}

func (self *IAgoraCollector) parseReviewInfo(url string, element *goquery.Selection) (*reviewInfo, error) {
	customerName := element.Find("div.reviewer").Text()
	customerName = strings.TrimSpace(customerName)
	if len(customerName) == 0 {
		return nil, ErrIAgoraCollectorMalformedReview.Raise().With("no customer name").
			Extra(map[string]any{"url": url})
	}

	var customerLocation *string
	rawCustomerLocation := element.Find("div.meta").Text()
	if len(rawCustomerLocation) > 0 {
		rawCustomerLocation = strings.Split(rawCustomerLocation, "\n")[0]
		rawCustomerLocation = strings.TrimSpace(rawCustomerLocation)
		customerLocation = &rawCustomerLocation
	}

	var metadataRating *float64
	rawMetadataRating, _ := element.Find("img.stars-small").Attr("title")
	if len(rawMetadataRating) > 0 {
		rawMetadataRating = strings.Split(rawMetadataRating, "/")[0]
		rawMetadataRating = strings.TrimSpace(rawMetadataRating)
		rating, err := strconv.ParseFloat(rawMetadataRating, 64)
		if err == nil {
			metadataRating = &rating
		}
	}

	var metadataVotes *int
	rawMetadataVotes := element.Find("div.useful > span").Text()
	rawMetadataVotes = strings.TrimSpace(rawMetadataVotes)
	if len(rawMetadataVotes) > 0 {
		votes, err := strconv.ParseInt(rawMetadataVotes, 10, 0)
		if err == nil {
			metadataVotes = kitUtil.Pointer(int(votes))
		}
	}

	metadataLink := url

	root := scraper.GetRoot(element)
	rawSchema := root.Find("head > script[type='application/ld+json']").Text()

	var schema webSchema
	err := json.Unmarshal([]byte(rawSchema), &schema)
	if err != nil {
		return nil, ErrIAgoraCollectorMalformedReview.Raise().With("no web schema").
			Extra(map[string]any{"url": url}).Cause(err)
	}

	var rawPostedAt string
	for _, review := range schema.Reviews {
		if review.Author == customerName {
			rawPostedAt = review.Date
			break
		}
	}
	if len(rawPostedAt) == 0 {
		return nil, ErrIAgoraCollectorMalformedReview.Raise().With("no posted at").
			Extra(map[string]any{"url": url})
	}

	postedAt, err := time.Parse("2006-01-02", rawPostedAt)
	if err != nil {
		return nil, ErrIAgoraCollectorMalformedReview.Raise().With("no posted at").
			Extra(map[string]any{"url": url, "posted_at": rawPostedAt}).Cause(err)
	}

	return &reviewInfo{
		CustomerName:     customerName,
		CustomerLocation: customerLocation,
		MetadataRating:   metadataRating,
		MetadataVotes:    metadataVotes,
		MetadataLink:     metadataLink,
		PostedAt:         postedAt,
	}, nil
}

func (self *IAgoraCollector) parseOverallBasicReviewContent(url string, element *goquery.Selection) string {
	comments := element.Find("div.txt").First().Text()
	comments = strings.TrimSpace(comments)
	comments = strings.TrimSuffix(comments, "Read more\u00a0>")

	careers := element.Find("div.label:contains('Careers') + div.txt").Text()
	careers = strings.TrimSpace(careers)
	careers = strings.TrimSuffix(careers, "Read more\u00a0>")

	pros := element.Find("div.label:contains('Pros') + div.txt").Text()
	pros = strings.TrimSpace(pros)
	pros = strings.TrimSuffix(pros, "Read more\u00a0>")

	cons := element.Find("div.label:contains('Cons') + div.txt").Text()
	cons = strings.TrimSpace(cons)
	cons = strings.TrimSuffix(cons, "Read more\u00a0>")

	return strings.TrimSpace(fmt.Sprintf(`
About careers:
%s
Advantages:
%s
Disadvantages:
%s
Final comments:
%s
`, careers, pros, cons, comments))
}

func (self *IAgoraCollector) parseOverallAdvancedReviewContent(url string, element *goquery.Selection) string {
	wish := element.Find("div.label:contains('I wish I had known...') + div.txt").Text()
	wish = strings.TrimSpace(wish)
	wish = strings.TrimSuffix(wish, "Read more\u00a0>")

	opinion := element.Find("div.label:contains('In my opinion:') + div.txt").Text()
	opinion = strings.TrimSpace(opinion)
	opinion = strings.TrimSuffix(opinion, "Read more\u00a0>")

	recommendation := element.Find("div.label:contains('Personal recommendation') + div.txt").Text()
	recommendation = strings.TrimSpace(recommendation)
	recommendation = strings.TrimSuffix(recommendation, "Read more\u00a0>")

	comments := element.Find("div.label:contains('Final comments') + div.txt").Text()
	comments = strings.TrimSpace(comments)
	comments = strings.TrimSuffix(comments, "Read more\u00a0>")

	return strings.TrimSpace(fmt.Sprintf(`
I wish I had known:
%s
In my opinion:
%s
Personal recommendation:
%s
Final comments:
%s
`, wish, opinion, recommendation, comments))
}

func (self *IAgoraCollector) parseOverallReviewContent(url string, element *goquery.Selection) (string, string) {
	title := element.Find("div.simple-title").Text()
	title = strings.Trim(title, `“”`)
	title = strings.TrimSpace(title)

	if len(title) > 0 {
		return title, self.parseOverallBasicReviewContent(url, element)
	}

	return "", self.parseOverallAdvancedReviewContent(url, element)
}

func (self *IAgoraCollector) parseAcademicReviewContent(url string, element *goquery.Selection) (string, string) {
	recommendations := element.Find("div.label:contains('Course recommendations') + div.txt").Text()
	recommendations = strings.TrimSpace(recommendations)
	recommendations = strings.TrimSuffix(recommendations, "Read more\u00a0>")

	comments := element.Find("div.label:contains('Personal comments') + div.txt").Text()
	comments = strings.TrimSpace(comments)
	comments = strings.TrimSuffix(comments, "Read more\u00a0>")

	return "", strings.TrimSpace(fmt.Sprintf(`
Course recommendations:
%s
Final comments:
%s
`, recommendations, comments))
}

func (self *IAgoraCollector) parseHousingReviewContent(url string, element *goquery.Selection) (string, string) {
	comments := element.Find("div.label:contains('Personal comments') + div.txt").Text()
	comments = strings.TrimSpace(comments)
	comments = strings.TrimSuffix(comments, "Read more\u00a0>")

	return "", comments
}

func (self *IAgoraCollector) parseLanguagesReviewContent(url string, element *goquery.Selection) (string, string) {
	comments := element.Find("div.label:contains('Personal comments') + div.txt").Text()
	comments = strings.TrimSpace(comments)
	comments = strings.TrimSuffix(comments, "Read more\u00a0>")

	return "", comments
}

type IAgoraCollectorCollectParams struct {
	CollectorID string
}

func (self *IAgoraCollector) Collect(ctx context.Context, task *asynq.Task) error {
	params := IAgoraCollectorCollectParams{}

	err := json.Unmarshal(task.Payload(), &params)
	if err != nil {
		self.observer.Error(ctx, kit.ErrWorkerGeneric.Raise().Cause(err))
		return nil
	}

	collector, product, organization, err := self.getCollectorProductAndOrganization(ctx, params.CollectorID)
	if err != nil {
		return err
	} else if collector == nil || product == nil || organization == nil {
		return nil
	}

	settings := collector.Settings.(IAgoraCollectorSettings)
	jobdata := collector.Jobdata.(IAgoraCollectorJobdata)

	reviews := IAGORA_COLLECTOR_MAX_REVIEWS_TO_COLLECT
	if jobdata.LastCollectedAt != nil {
		days := int(time.Since(*jobdata.LastCollectedAt).Hours() / 24)
		reviews = min(IAGORA_COLLECTOR_MAX_REVIEWS_TO_COLLECT, IAGORA_COLLECTOR_DAILY_REVIEWS_TO_COLLECT*days)
	}

	reviews = min(organization.UsageLeft(), reviews)
	if reviews <= 0 {
		return nil
	}
	reviews = int(math.Ceil(float64(reviews)/float64(IAGORA_COLLECTOR_MIN_REVIEWS_TO_COLLECT))) * IAGORA_COLLECTOR_MIN_REVIEWS_TO_COLLECT

	pages := int(math.Ceil(float64(reviews) / float64(IAGORA_COLLECTOR_REVIEWS_PER_PAGE)))
	urls := make([]string, 0, pages)
	for page := range pages {
		urls = append(urls, fmt.Sprintf(IAGORA_COLLECTOR_BASE_URL, url.QueryEscape(settings.Institution), page))
	}

	now := time.Now()
	feedbacks := []feedback.Feedback{}
	mutex := sync.Mutex{}

	err = self.scraper.Scrape(ctx, urls, nil,
		"div.one-review", func(url string, element *goquery.Selection) {
			info, err := self.parseReviewInfo(url, element)
			if err != nil {
				self.observer.Error(ctx, err)
				return
			}

			var title string
			var body string
			_type := element.Find("div.rtype").Text()
			switch strings.ToUpper(_type) {
			case "OVERALL":
				title, body = self.parseOverallReviewContent(url, element)
			case "ACADEMIC":
				title, body = self.parseAcademicReviewContent(url, element)
			case "HOUSING":
				title, body = self.parseHousingReviewContent(url, element)
			case "LANGUAGES":
				title, body = self.parseLanguagesReviewContent(url, element)
			case "STUDENT LIFE":
				return // Ignore student life reviews
			case "EXPENSES":
				return // Ignore expenses reviews
			default:
				self.observer.Error(ctx, ErrIAgoraCollectorMalformedReview.Raise().
					With("unknown type '%s'", _type).Extra(map[string]any{"url": url}))
				return
			}

			content := feedback.CleanContent(title, body)
			if len(content) == 0 {
				return
			}
			hash := feedback.ComputeHash(feedback.FeedbackSourceIAgora, info.CustomerName, content)

			_feedback := feedback.NewFeedback()
			_feedback.ID = xid.New().String()
			_feedback.ProductID = product.ID
			_feedback.Hash = hash
			_feedback.Source = feedback.FeedbackSourceIAgora
			_feedback.Customer.Email = nil
			_feedback.Customer.Name = info.CustomerName
			_feedback.Customer.Picture = feedback.FEEDBACK_CUSTOMER_DEFAULT_PICTURE
			_feedback.Customer.Location = info.CustomerLocation
			_feedback.Customer.Verified = nil
			_feedback.Customer.Reviews = nil
			_feedback.Customer.Link = nil
			_feedback.Content = content
			_feedback.Language = engine.OPTION_UNKNOWN
			_feedback.Translation = ""
			_feedback.Release = engine.OPTION_UNKNOWN
			_feedback.Metadata.Rating = info.MetadataRating
			_feedback.Metadata.Media = nil
			_feedback.Metadata.Verified = nil
			_feedback.Metadata.Votes = info.MetadataVotes
			_feedback.Metadata.Link = &info.MetadataLink
			_feedback.Tokens = 0
			_feedback.PostedAt = info.PostedAt
			_feedback.CollectedAt = now
			_feedback.TranslatedAt = nil
			_feedback.ProcessedAt = nil

			mutex.Lock()
			feedbacks = append(feedbacks, *_feedback)
			mutex.Unlock()
		})
	if err != nil {
		return err
	}

	// This is not as safe and performant as the other collectors as
	// it can potentially create a big feedback array in memory, but
	// the scraper being async makes enhancing this somewhat complicated
	totalFeedbacks := 0
	newFeedbacks := 0
	lastChunk := 0
	for {
		chunk := min(lastChunk+1000, len(feedbacks))

		_newFeedbacks, err := self.saveAndEnqueue(ctx, feedbacks[lastChunk:chunk])
		if err != nil {
			return err
		}

		totalFeedbacks += (chunk - lastChunk)
		newFeedbacks += _newFeedbacks

		lastChunk = chunk
		if lastChunk >= len(feedbacks) {
			break
		}
	}

	jobdata.LastCollectedAt = &now

	// We may rescrape some reviews if the jobdata update isn't in a transaction,
	// but then how to bulk insert in batches in a performant way?
	collector.Jobdata = jobdata
	err = self.collectorRepository.UpdateJobdata(ctx, *collector)
	if err != nil {
		return err
	}

	self.observer.Infof(ctx, "Collected %d IAgora reviews of which %d were duplicated",
		totalFeedbacks, totalFeedbacks-newFeedbacks)

	return nil
}

func (self *IAgoraCollector) Schedule(ctx context.Context, _ *asynq.Task) error {
	ids, err := self.collectorRepository.ListIDsByTypeNotDeleted(ctx, CollectorTypeIAgora)
	if err != nil {
		return err
	}

	for _, id := range ids {
		err := self.enqueuer.Enqueue(ctx, IAgoraCollectorCollect, IAgoraCollectorCollectParams{
			CollectorID: id,
		}, asynq.MaxRetry(2), asynq.Unique(24*time.Hour))
		if err != nil {
			self.observer.Error(ctx, err)
		}
	}

	return nil
}
