package review

import (
	"context"
	"fmt"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/feedback"
	"backend/pkg/util"
)

type ReviewRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewReviewRepository(observer *kit.Observer, database *kit.Database, config config.Config) *ReviewRepository {
	return &ReviewRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *ReviewRepository) Create(ctx context.Context, review Review) (*Review, error) {
	m := ReviewFeedbackModel{
		Review:   *NewReviewModel(review),
		Feedback: *feedback.NewFeedbackModel(review.Feedback),
	}

	stmt := sqlf.
		InsertInto(REVIEW_MODEL_TABLE).
		Set("id", m.Review.ID).
		Set("product_id", m.Review.ProductID).
		Set("feedback_id", m.Review.FeedbackID).
		Set("keywords", m.Review.Keywords).
		Set("sentiment", m.Review.Sentiment).
		Set("emotions", m.Review.Emotions).
		Set("intention", m.Review.Intention).
		Set("category", m.Review.Category).
		Set("quality", m.Review.Quality).
		Set("created_at", m.Review.CreatedAt).
		Set("exported_at", m.Review.ExportedAt).
		Returning("*").To(&m.Review)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return m.Review.ToEntity(m.Feedback), nil
}

func (self *ReviewRepository) GetByID(ctx context.Context, id string) (*Review, error) {
	var m ReviewFeedbackModel

	stmt := sqlf.
		Select(REVIEW_MODEL_TABLE+".*").
		Select(`0 as "notate:feedback"`).
		Select(feedback.FEEDBACK_MODEL_TABLE+".*").To(&m).
		From(REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+REVIEW_MODEL_TABLE+".feedback_id").
		Where(REVIEW_MODEL_TABLE+".id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return m.Review.ToEntity(m.Feedback), nil
}

func (self *ReviewRepository) ListByProductID(ctx context.Context,
	productID string, search ReviewSearch) (*util.Page[Review, time.Time], error) {
	var ms []ReviewFeedbackModel

	stmt := sqlf.
		Select(REVIEW_MODEL_TABLE+".*").
		Select(`0 as "notate:feedback"`).
		Select(feedback.FEEDBACK_MODEL_TABLE+".*").To(&ms).
		From(REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+REVIEW_MODEL_TABLE+".feedback_id").
		Where(REVIEW_MODEL_TABLE+".product_id = ?", productID)

	if search.Filters.Sources != nil && len(*search.Filters.Sources) > 0 {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE + ".source").In(util.Spread(*search.Filters.Sources)...)
	}

	if search.Filters.Releases != nil && len(*search.Filters.Releases) > 0 {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE + ".release").In(util.Spread(*search.Filters.Releases)...)
	}

	if search.Filters.Categories != nil && len(*search.Filters.Categories) > 0 {
		stmt.
			Where(REVIEW_MODEL_TABLE + ".category").In(util.Spread(*search.Filters.Categories)...)
	}

	if search.Filters.Keywords != nil && len(*search.Filters.Keywords) > 0 {
		stmt.
			Where(REVIEW_MODEL_TABLE+".keywords && ?", *search.Filters.Keywords)
	}

	if search.Filters.Sentiments != nil && len(*search.Filters.Sentiments) > 0 {
		stmt.
			Where(REVIEW_MODEL_TABLE + ".sentiment").In(util.Spread(*search.Filters.Sentiments)...)
	}

	if search.Filters.Emotions != nil && len(*search.Filters.Emotions) > 0 {
		stmt.
			Where(REVIEW_MODEL_TABLE+".emotions && ?", *search.Filters.Emotions)
	}

	if search.Filters.Intentions != nil && len(*search.Filters.Intentions) > 0 {
		stmt.
			Where(REVIEW_MODEL_TABLE + ".intention").In(util.Spread(*search.Filters.Intentions)...)
	}

	if search.Filters.Languages != nil && len(*search.Filters.Languages) > 0 {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE + ".language").In(util.Spread(*search.Filters.Languages)...)
	}

	if search.Filters.SeenStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *search.Filters.SeenStartAt)
	}

	if search.Filters.SeenEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *search.Filters.SeenEndAt)
	}

	if search.Pagination.From != nil {
		if search.Orders.Recency == ReviewSearchOrdersAscending {
			stmt.
				Where(fmt.Sprintf("(%s.posted_at, %s.id) > (?, ?)",
					feedback.FEEDBACK_MODEL_TABLE, feedback.FEEDBACK_MODEL_TABLE),
					search.Pagination.From.Value, search.Pagination.From.ID)
		} else {
			stmt.
				Where(fmt.Sprintf("(%s.posted_at, %s.id) < (?, ?)",
					feedback.FEEDBACK_MODEL_TABLE, feedback.FEEDBACK_MODEL_TABLE),
					search.Pagination.From.Value, search.Pagination.From.ID)
		}
	}

	if search.Orders.Recency == ReviewSearchOrdersAscending {
		stmt.
			OrderBy(feedback.FEEDBACK_MODEL_TABLE+".posted_at ASC", feedback.FEEDBACK_MODEL_TABLE+".id ASC")
	} else {
		stmt.
			OrderBy(feedback.FEEDBACK_MODEL_TABLE+".posted_at DESC", feedback.FEEDBACK_MODEL_TABLE+".id DESC")
	}

	stmt.
		Limit(search.Pagination.Limit)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &util.Page[Review, time.Time]{}, nil
		}

		return nil, err
	}

	items := make([]Review, 0, len(ms))
	for _, m := range ms {
		items = append(items, *m.Review.ToEntity(m.Feedback))
	}

	var cursor *util.Cursor[time.Time]
	if len(items) == search.Pagination.Limit {
		cursor = &util.Cursor[time.Time]{
			Value: items[search.Pagination.Limit-1].Feedback.PostedAt,
			ID:    items[search.Pagination.Limit-1].Feedback.ID,
		}
	}

	return &util.Page[Review, time.Time]{
		Items: items,
		Next:  cursor,
	}, nil
}

func (self *ReviewRepository) UpdateQuality(ctx context.Context, id string, quality int) error {
	stmt := sqlf.
		Update(REVIEW_MODEL_TABLE).
		Set("quality", quality).
		Where("id = ?", id)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}
