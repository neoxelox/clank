package scraper

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/neoxelox/errors"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

const (
	SCRAPER_TIMEOUT       = 5 * time.Second
	SCRAPER_MAX_URL_SIZE  = 256     // 256 B
	SCRAPER_MAX_BODY_SIZE = 2 << 20 // 2 MB
)

var (
	ErrScraperGeneric  = errors.New("scraper failed")
	ErrScraperTimedOut = errors.New("scraper timed out")
)

type Scraper struct {
	config   config.Config
	observer *kit.Observer
	scraper  *colly.Collector
}

func NewScraper(observer *kit.Observer, config config.Config) *Scraper {
	scraper := colly.NewCollector(
		colly.UserAgent("Random"),
		colly.MaxDepth(0),
		colly.ParseHTTPErrorResponse(),
		colly.AllowURLRevisit(),
		colly.MaxBodySize(SCRAPER_MAX_BODY_SIZE),
		colly.IgnoreRobotsTxt(),
		colly.Async(true),
		colly.DetectCharset(),
	)

	scraper.SetRequestTimeout(SCRAPER_TIMEOUT)
	extensions.RandomUserAgent(scraper)
	extensions.URLLengthFilter(scraper, SCRAPER_MAX_URL_SIZE)

	if config.Service.Environment == kit.EnvDevelopment {
		scraper.SetDebugger(&debug.LogDebugger{})
	}

	return &Scraper{
		config:   config,
		observer: observer,
		scraper:  scraper,
	}
}

func (self *Scraper) Scrape(ctx context.Context, urls []string, rules []*colly.LimitRule,
	selector string, callback func(url string, element *goquery.Selection)) error {
	scraper := self.scraper.Clone()

	domains := make([]string, 0, len(urls))
	for _, _url := range urls {
		url, err := url.Parse(_url)
		if err != nil {
			return ErrScraperGeneric.Raise().Cause(err)
		}

		domains = append(domains, url.Hostname())
	}
	scraper.AllowedDomains = domains

	err := scraper.Limits(rules)
	if err != nil {
		return ErrScraperGeneric.Raise().Cause(err)
	}

	var errors int
	scraper.OnError(func(response *colly.Response, err error) {
		if err != nil {
			_url := response.Request.URL.String()

			if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
				err = ErrScraperTimedOut.Raise().
					Extra(map[string]any{"url": _url, "timeout": SCRAPER_TIMEOUT}).
					Cause(err)
			} else {
				err = ErrScraperGeneric.Raise().
					Extra(map[string]any{"url": _url, "status": response.StatusCode}).
					Cause(err)
			}

			self.observer.Error(ctx, err)
			errors++
		}
	})

	scraper.OnHTML(selector, func(element *colly.HTMLElement) {
		url := element.Request.URL.String()

		html, err := element.DOM.Html()
		if err != nil {
			self.observer.Error(ctx, ErrScraperGeneric.Raise().Extra(map[string]any{"url": url}).Cause(err))
			errors++
			return
		}

		html = strings.ReplaceAll(html, "<br>", "\n")
		html = strings.ReplaceAll(html, "<br/>", "\n")
		html = strings.ReplaceAll(html, "<noscript>", "")
		html = strings.ReplaceAll(html, "</noscript>", "")

		element.DOM.SetHtml(html)

		callback(url, element.DOM)
	})

	for _, url := range urls {
		err := scraper.Visit(url)
		if err != nil {
			self.observer.Error(ctx, ErrScraperGeneric.Raise().Extra(map[string]any{"url": url}).Cause(err))
			errors++
			break
		}
	}

	scraper.Wait()

	if errors > 0 {
		return ErrScraperGeneric.Raise().With("%d errors encountered while scraping", errors)
	}

	return nil
}

func GetRoot(element *goquery.Selection) *goquery.Selection {
	root := element

	for {
		if len(root.Parent().Nodes) == 0 {
			break
		}
		root = root.Parent()
	}

	return root
}
