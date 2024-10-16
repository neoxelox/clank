package feedback

import (
	"context"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/util"
)

type FeedbackRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewFeedbackRepository(observer *kit.Observer, database *kit.Database, config config.Config) *FeedbackRepository {
	return &FeedbackRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *FeedbackRepository) Create(ctx context.Context, feedback Feedback) (*Feedback, error) {
	f := NewFeedbackModel(feedback)

	stmt := sqlf.
		InsertInto(FEEDBACK_MODEL_TABLE).
		Set("id", f.ID).
		Set("product_id", f.ProductID).
		Set("hash", f.Hash).
		Set("source", f.Source).
		Set("customer", f.Customer).
		Set("content", f.Content).
		Set("language", f.Language).
		Set("translation", f.Translation).
		Set("release", f.Release).
		Set("metadata", f.Metadata).
		Set("tokens", f.Tokens).
		Set("posted_at", f.PostedAt).
		Set("collected_at", f.CollectedAt).
		Set("translated_at", f.TranslatedAt).
		Set("processed_at", f.ProcessedAt).
		Returning("*").To(&f)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return f.ToEntity(), nil
}

func (self *FeedbackRepository) BulkCreate(ctx context.Context, feedbacks []Feedback) (int, error) {
	if len(feedbacks) == 0 {
		return 0, nil
	}

	stmt := sqlf.
		InsertInto(FEEDBACK_MODEL_TABLE)

	for _, feedback := range feedbacks {
		f := NewFeedbackModel(feedback)

		stmt.
			NewRow().
			Set("id", f.ID).
			Set("product_id", f.ProductID).
			Set("hash", f.Hash).
			Set("source", f.Source).
			Set("customer", f.Customer).
			Set("content", f.Content).
			Set("language", f.Language).
			Set("translation", f.Translation).
			Set("release", f.Release).
			Set("metadata", f.Metadata).
			Set("tokens", f.Tokens).
			Set("posted_at", f.PostedAt).
			Set("collected_at", f.CollectedAt).
			Set("translated_at", f.TranslatedAt).
			Set("processed_at", f.ProcessedAt)
	}

	stmt.
		Clause("ON CONFLICT DO NOTHING")

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (self *FeedbackRepository) GetByID(ctx context.Context, id string) (*Feedback, error) {
	var f FeedbackModel

	stmt := sqlf.
		Select("*").To(&f).
		From(FEEDBACK_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return f.ToEntity(), nil
}

func (self *FeedbackRepository) GetByIDForUpdate(ctx context.Context, id string) (*Feedback, error) {
	var f FeedbackModel

	stmt := sqlf.
		Select("*").To(&f).
		From(FEEDBACK_MODEL_TABLE).
		Where("id = ?", id).
		Clause("FOR NO KEY UPDATE")

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return f.ToEntity(), nil
}

func (self *FeedbackRepository) ListIDsByNotTranslated(ctx context.Context,
	pagination util.Pagination[time.Time]) (*util.Page[string, time.Time], error) {
	var result []struct {
		ID          string    `db:"id"`
		CollectedAt time.Time `db:"collected_at"`
	}

	stmt := sqlf.
		Select("id, collected_at").To(&result).
		From(FEEDBACK_MODEL_TABLE).
		Where("translated_at IS NULL")

	if pagination.From != nil {
		stmt.
			Where("(collected_at, id) > (?, ?)", pagination.From.Value, pagination.From.ID)
	}

	stmt.
		OrderBy("collected_at ASC", "id ASC").
		Limit(pagination.Limit)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &util.Page[string, time.Time]{}, nil
		}

		return nil, err
	}

	items := make([]string, 0, len(result))
	for _, res := range result {
		items = append(items, res.ID)
	}

	var cursor *util.Cursor[time.Time]
	if len(result) == pagination.Limit {
		cursor = &util.Cursor[time.Time]{
			Value: result[pagination.Limit-1].CollectedAt,
			ID:    result[pagination.Limit-1].ID,
		}
	}

	return &util.Page[string, time.Time]{
		Items: items,
		Next:  cursor,
	}, nil
}

func (self *FeedbackRepository) ListIDsByNotProcessed(ctx context.Context,
	pagination util.Pagination[time.Time]) (*util.Page[string, time.Time], error) {
	var result []struct {
		ID           string    `db:"id"`
		TranslatedAt time.Time `db:"translated_at"`
	}

	stmt := sqlf.
		Select("id, translated_at").To(&result).
		From(FEEDBACK_MODEL_TABLE).
		Where("processed_at IS NULL").
		Where("translated_at IS NOT NULL")

	if pagination.From != nil {
		stmt.
			Where("(translated_at, id) > (?, ?)", pagination.From.Value, pagination.From.ID)
	}

	stmt.
		OrderBy("translated_at ASC", "id ASC").
		Limit(pagination.Limit)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &util.Page[string, time.Time]{}, nil
		}

		return nil, err
	}

	items := make([]string, 0, len(result))
	for _, res := range result {
		items = append(items, res.ID)
	}

	var cursor *util.Cursor[time.Time]
	if len(result) == pagination.Limit {
		cursor = &util.Cursor[time.Time]{
			Value: result[pagination.Limit-1].TranslatedAt,
			ID:    result[pagination.Limit-1].ID,
		}
	}

	return &util.Page[string, time.Time]{
		Items: items,
		Next:  cursor,
	}, nil
}

func (self *FeedbackRepository) UpdateTranslated(ctx context.Context, feedback Feedback) error {
	f := NewFeedbackModel(feedback)

	stmt := sqlf.
		Update(FEEDBACK_MODEL_TABLE).
		Set("language", f.Language).
		Set("translation", f.Translation).
		Set("tokens", f.Tokens).
		Set("translated_at", f.TranslatedAt).
		Where("id = ?", f.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *FeedbackRepository) UpdateProcessed(ctx context.Context, feedback Feedback) error {
	f := NewFeedbackModel(feedback)

	stmt := sqlf.
		Update(FEEDBACK_MODEL_TABLE).
		Set("tokens", f.Tokens).
		Set("processed_at", f.ProcessedAt).
		Where("id = ?", f.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *FeedbackRepository) UpdateTokens(ctx context.Context, id string, tokens int) error {
	stmt := sqlf.
		Update(FEEDBACK_MODEL_TABLE).
		Set("tokens", tokens).
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
