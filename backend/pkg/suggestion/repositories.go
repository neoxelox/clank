package suggestion

import (
	"context"
	"fmt"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/feedback"
	"backend/pkg/util"

	"github.com/pgvector/pgvector-go"
)

type SuggestionRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewSuggestionRepository(observer *kit.Observer, database *kit.Database,
	config config.Config) *SuggestionRepository {
	return &SuggestionRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *SuggestionRepository) Create(ctx context.Context,
	suggestion Suggestion, feedback feedback.Feedback) (*Suggestion, error) {
	s := NewSuggestionModel(suggestion)

	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		stmt := sqlf.
			InsertInto(SUGGESTION_MODEL_TABLE).
			Set("id", s.ID).
			Set("product_id", s.ProductID).
			Set("embedding", s.Embedding).
			Set("sources", s.Sources).
			Set("title", s.Title).
			Set("description", s.Description).
			Set("reason", s.Reason).
			Set("importances", s.Importances).
			Set("priority", s.Priority).
			Set("categories", s.Categories).
			Set("releases", s.Releases).
			Set("customers", s.Customers).
			Set("assignee_id", s.AssigneeID).
			Set("quality", s.Quality).
			Set("first_seen_at", s.FirstSeenAt).
			Set("last_seen_at", s.LastSeenAt).
			Set("created_at", s.CreatedAt).
			Set("archived_at", s.ArchivedAt).
			Set("last_aggregated_at", s.LastAggregatedAt).
			Set("exported_at", s.ExportedAt).
			Returning("*").To(&s)

		err := self.database.Query(ctx, stmt)
		if err != nil {
			return err
		}

		stmt = sqlf.
			InsertInto(SUGGESTION_FEEDBACK_MODEL_TABLE).
			Set("suggestion_id", suggestion.ID).
			Set("feedback_id", feedback.ID)

		affected, err := self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		if affected != 1 {
			return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.ToEntity(), nil
}

func (self *SuggestionRepository) GetByID(ctx context.Context, id string) (*Suggestion, error) {
	var s SuggestionModel

	stmt := sqlf.
		Select("*").To(&s).
		From(SUGGESTION_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return s.ToEntity(), nil
}

func (self *SuggestionRepository) GetByIDForUpdate(ctx context.Context, id string) (*Suggestion, error) {
	var s SuggestionModel

	stmt := sqlf.
		Select("*").To(&s).
		From(SUGGESTION_MODEL_TABLE).
		Where("id = ?", id).
		Clause("FOR NO KEY UPDATE NOWAIT")

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return s.ToEntity(), nil
}

func (self *SuggestionRepository) ListByProductID(ctx context.Context,
	productID string, search SuggestionSearch) (*util.Page[Suggestion, SuggestionSearchCursor], error) {
	var ss []SuggestionModel

	stmt := sqlf.
		Select("*").To(&ss).
		From(SUGGESTION_MODEL_TABLE).
		Where("product_id = ?", productID)

	if search.Filters.Embedding != nil && len(*search.Filters.Embedding) > 0 {
		stmt.
			Where("(1 - (embedding <=> ?::vector)) >= ?", pgvector.NewVector(*search.Filters.Embedding), 0.3)
	}

	if search.Filters.Sources != nil && len(*search.Filters.Sources) > 0 {
		stmt.
			Where("").
			SubQuery("EXISTS (", ")", sqlf.
				Select("TRUE").
				From("jsonb_object_keys(sources) AS source").
				Where("source").In(util.Spread(*search.Filters.Sources)...))
	}

	if search.Filters.Importances != nil && len(*search.Filters.Importances) > 0 {
		stmt.
			Where("").
			SubQuery("(", ")", sqlf.
				Select("importance.key").
				From("jsonb_each_text(importances) as importance").
				OrderBy("importance.value::int DESC", "importance.key DESC").
				Limit(1)).
			In(util.Spread(*search.Filters.Importances)...)
	}

	if search.Filters.Releases != nil && len(*search.Filters.Releases) > 0 {
		stmt.
			Where("").
			SubQuery("EXISTS (", ")", sqlf.
				Select("TRUE").
				From("jsonb_object_keys(releases) AS release").
				Where("release").In(util.Spread(*search.Filters.Releases)...))
	}

	if search.Filters.Categories != nil && len(*search.Filters.Categories) > 0 {
		stmt.
			Where("").
			SubQuery("(", ")", sqlf.
				Select("category.key").
				From("jsonb_each_text(categories) as category").
				OrderBy("category.value::int DESC", "category.key DESC").
				Limit(1)).
			In(util.Spread(*search.Filters.Categories)...)
	}

	if search.Filters.Assignees != nil {
		unAssignees, assignees := util.Extract(*search.Filters.Assignees, func(value string) bool {
			return value == SuggestionSearchFiltersAssigneesUnassigned
		})

		if len(unAssignees) > 0 && len(assignees) > 0 {
			stmt.
				Where("(assignee_id IS NULL OR assignee_id").In(util.Spread(assignees)...).Clause(")")
		} else if len(unAssignees) > 0 {
			stmt.
				Where("assignee_id IS NULL")
		} else if len(assignees) > 0 {
			stmt.
				Where("assignee_id").In(util.Spread(assignees)...)
		}
	}

	if search.Filters.Status != nil {
		switch *search.Filters.Status {
		case SuggestionSearchFiltersStatusActive:
			stmt.
				Where("(archived_at IS NULL OR (archived_at IS NOT NULL AND last_seen_at > archived_at))")
		case SuggestionSearchFiltersStatusRegressed:
			stmt.
				Where("archived_at IS NOT NULL").
				Where("last_seen_at > archived_at")
		case SuggestionSearchFiltersStatusArchived:
			stmt.
				Where("archived_at IS NOT NULL")
		case SuggestionSearchFiltersStatusUnarchived:
			stmt.
				Where("archived_at IS NULL")
		}
	}

	if search.Filters.FirstSeenStartAt != nil {
		stmt.
			Where("first_seen_at >= ?", *search.Filters.FirstSeenStartAt)
	}

	if search.Filters.FirstSeenEndAt != nil {
		stmt.
			Where("first_seen_at <= ?", *search.Filters.FirstSeenEndAt)
	}

	if search.Filters.LastSeenStartAt != nil {
		stmt.
			Where("last_seen_at >= ?", *search.Filters.LastSeenStartAt)
	}

	if search.Filters.LastSeenEndAt != nil {
		stmt.
			Where("last_seen_at <= ?", *search.Filters.LastSeenEndAt)
	}

	if search.Pagination.From != nil {
		if search.Orders.Relevance == SuggestionSearchOrdersAscending {
			stmt.
				Where("(priority, last_seen_at, id) > (?, ?, ?)", search.Pagination.From.Value.Priority,
					search.Pagination.From.Value.SeenAt, search.Pagination.From.ID)
		} else {
			stmt.
				Where("(priority, last_seen_at, id) < (?, ?, ?)", search.Pagination.From.Value.Priority,
					search.Pagination.From.Value.SeenAt, search.Pagination.From.ID)
		}
	}

	if search.Orders.Relevance == SuggestionSearchOrdersAscending {
		stmt.
			OrderBy("priority ASC", "last_seen_at ASC", "id ASC")
	} else {
		stmt.
			OrderBy("priority DESC", "last_seen_at DESC", "id DESC")
	}

	stmt.
		Limit(search.Pagination.Limit)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &util.Page[Suggestion, SuggestionSearchCursor]{}, nil
		}

		return nil, err
	}

	items := make([]Suggestion, 0, len(ss))
	for _, s := range ss {
		items = append(items, *s.ToEntity())
	}

	var cursor *util.Cursor[SuggestionSearchCursor]
	if len(items) == search.Pagination.Limit {
		cursor = &util.Cursor[SuggestionSearchCursor]{
			Value: SuggestionSearchCursor{
				Priority: items[search.Pagination.Limit-1].Priority,
				SeenAt:   items[search.Pagination.Limit-1].LastSeenAt,
			},
			ID: items[search.Pagination.Limit-1].ID,
		}
	}

	return &util.Page[Suggestion, SuggestionSearchCursor]{
		Items: items,
		Next:  cursor,
	}, nil
}

func (self *SuggestionRepository) ListByEmbeddingAndProductID(ctx context.Context,
	embedding []float32, threshold float64, limit int, productID string) ([]Suggestion, error) {
	var result []struct {
		SuggestionModel
		Score float64 `db:"score"`
	}

	stmt := sqlf.
		Select("*").To(&result).
		From("").
		SubQuery("(", ")", sqlf.
			Select("*").
			Select("1 - (embedding <=> ?::vector) AS score", pgvector.NewVector(embedding)).
			From(SUGGESTION_MODEL_TABLE).
			Where("product_id = ?", productID).
			OrderBy("score DESC").
			Limit(limit)).
		Where("score >= ?", threshold)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []Suggestion{}, nil
		}

		return nil, err
	}

	entities := make([]Suggestion, 0, len(result))
	for _, res := range result {
		entities = append(entities, *res.SuggestionModel.ToEntity())
	}

	return entities, nil
}

func (self *SuggestionRepository) ListFeedbacks(ctx context.Context,
	id string, pagination util.Pagination[time.Time]) (*util.Page[feedback.Feedback, time.Time], error) {
	var fs []feedback.FeedbackModel

	stmt := sqlf.
		Select(feedback.FEEDBACK_MODEL_TABLE+".*").To(&fs).
		From(feedback.FEEDBACK_MODEL_TABLE).
		Join(SUGGESTION_FEEDBACK_MODEL_TABLE,
			SUGGESTION_FEEDBACK_MODEL_TABLE+".feedback_id = "+feedback.FEEDBACK_MODEL_TABLE+".id").
		Where(SUGGESTION_FEEDBACK_MODEL_TABLE+".suggestion_id = ?", id)

	if pagination.From != nil {
		stmt.
			Where(fmt.Sprintf("(%s.posted_at, %s.id) < (?, ?)",
				feedback.FEEDBACK_MODEL_TABLE, feedback.FEEDBACK_MODEL_TABLE),
				pagination.From.Value, pagination.From.ID)
	}

	stmt.
		OrderBy(feedback.FEEDBACK_MODEL_TABLE+".posted_at DESC", feedback.FEEDBACK_MODEL_TABLE+".id DESC").
		Limit(pagination.Limit)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &util.Page[feedback.Feedback, time.Time]{}, nil
		}

		return nil, err
	}

	items := make([]feedback.Feedback, 0, len(fs))
	for _, f := range fs {
		items = append(items, *f.ToEntity())
	}

	var cursor *util.Cursor[time.Time]
	if len(items) == pagination.Limit {
		cursor = &util.Cursor[time.Time]{
			Value: items[pagination.Limit-1].PostedAt,
			ID:    items[pagination.Limit-1].ID,
		}
	}

	return &util.Page[feedback.Feedback, time.Time]{
		Items: items,
		Next:  cursor,
	}, nil
}

func (self *SuggestionRepository) UpdateAssignee(ctx context.Context, id string, assigneeID *string) error {
	stmt := sqlf.
		Update(SUGGESTION_MODEL_TABLE).
		Set("assignee_id", assigneeID).
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

func (self *SuggestionRepository) UpdateQuality(ctx context.Context, id string, quality int) error {
	stmt := sqlf.
		Update(SUGGESTION_MODEL_TABLE).
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

func (self *SuggestionRepository) UpdateArchivedAt(ctx context.Context, id string, archivedAt *time.Time) error {
	stmt := sqlf.
		Update(SUGGESTION_MODEL_TABLE).
		Set("archived_at", archivedAt).
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

func (self *SuggestionRepository) UpdateAggregated(ctx context.Context, suggestion Suggestion,
	feedback feedback.Feedback) error {
	s := NewSuggestionModel(suggestion)

	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		stmt := sqlf.
			Update(SUGGESTION_MODEL_TABLE).
			Set("embedding", s.Embedding).
			Set("sources", s.Sources).
			Set("title", s.Title).
			Set("description", s.Description).
			Set("reason", s.Reason).
			Set("importances", s.Importances).
			Set("priority", s.Priority).
			Set("categories", s.Categories).
			Set("releases", s.Releases).
			Set("customers", s.Customers).
			Set("first_seen_at", s.FirstSeenAt).
			Set("last_seen_at", s.LastSeenAt).
			Set("last_aggregated_at", s.LastAggregatedAt).
			Where("id = ?", s.ID)

		affected, err := self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		if affected != 1 {
			return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
		}

		stmt = sqlf.
			InsertInto(SUGGESTION_FEEDBACK_MODEL_TABLE).
			Set("suggestion_id", suggestion.ID).
			Set("feedback_id", feedback.ID)

		affected, err = self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		if affected != 1 {
			return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

type PartialSuggestionRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewPartialSuggestionRepository(observer *kit.Observer, database *kit.Database,
	config config.Config) *PartialSuggestionRepository {
	return &PartialSuggestionRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *PartialSuggestionRepository) Create(ctx context.Context,
	partial PartialSuggestion) (*PartialSuggestion, error) {
	p := NewPartialSuggestionModel(partial)

	stmt := sqlf.
		InsertInto(PARTIAL_SUGGESTION_MODEL_TABLE).
		Set("id", p.ID).
		Set("feedback_id", p.FeedbackID).
		Set("title", p.Title).
		Set("description", p.Description).
		Set("reason", p.Reason).
		Set("importance", p.Importance).
		Set("category", p.Category).
		Set("created_at", p.CreatedAt).
		Returning("*").To(&p)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return p.ToEntity(), nil
}

func (self *PartialSuggestionRepository) BulkCreate(ctx context.Context, partials []PartialSuggestion) error {
	if len(partials) == 0 {
		return nil
	}

	stmt := sqlf.
		InsertInto(PARTIAL_SUGGESTION_MODEL_TABLE)

	for _, partial := range partials {
		p := NewPartialSuggestionModel(partial)

		stmt.
			NewRow().
			Set("id", p.ID).
			Set("feedback_id", p.FeedbackID).
			Set("title", p.Title).
			Set("description", p.Description).
			Set("reason", p.Reason).
			Set("importance", p.Importance).
			Set("category", p.Category).
			Set("created_at", p.CreatedAt)
	}

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != len(partials) {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, len(partials))
	}

	return nil
}

func (self *PartialSuggestionRepository) GetByID(ctx context.Context, id string) (*PartialSuggestion, error) {
	var p PartialSuggestionModel

	stmt := sqlf.
		Select("*").To(&p).
		From(PARTIAL_SUGGESTION_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return p.ToEntity(), nil
}

func (self *PartialSuggestionRepository) ListIDsByCreatedAt(ctx context.Context,
	pagination util.Pagination[time.Time]) (*util.Page[string, time.Time], error) {
	var result []struct {
		ID        string    `db:"id"`
		CreatedAt time.Time `db:"created_at"`
	}

	stmt := sqlf.
		Select("id, created_at").To(&result).
		From(PARTIAL_SUGGESTION_MODEL_TABLE)

	if pagination.From != nil {
		stmt.
			Where("(created_at, id) > (?, ?)", pagination.From.Value, pagination.From.ID)
	}

	stmt.
		OrderBy("created_at ASC", "id ASC").
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
			Value: result[pagination.Limit-1].CreatedAt,
			ID:    result[pagination.Limit-1].ID,
		}
	}

	return &util.Page[string, time.Time]{
		Items: items,
		Next:  cursor,
	}, nil
}

func (self *PartialSuggestionRepository) Delete(ctx context.Context, id string) error {
	stmt := sqlf.
		DeleteFrom(PARTIAL_SUGGESTION_MODEL_TABLE).
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
