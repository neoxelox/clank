package issue

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

type IssueRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewIssueRepository(observer *kit.Observer, database *kit.Database, config config.Config) *IssueRepository {
	return &IssueRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *IssueRepository) Create(ctx context.Context, issue Issue, feedback feedback.Feedback) (*Issue, error) {
	i := NewIssueModel(issue)

	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		stmt := sqlf.
			InsertInto(ISSUE_MODEL_TABLE).
			Set("id", i.ID).
			Set("product_id", i.ProductID).
			Set("embedding", i.Embedding).
			Set("sources", i.Sources).
			Set("title", i.Title).
			Set("description", i.Description).
			Set("steps", i.Steps).
			Set("severities", i.Severities).
			Set("priority", i.Priority).
			Set("categories", i.Categories).
			Set("releases", i.Releases).
			Set("customers", i.Customers).
			Set("assignee_id", i.AssigneeID).
			Set("quality", i.Quality).
			Set("first_seen_at", i.FirstSeenAt).
			Set("last_seen_at", i.LastSeenAt).
			Set("created_at", i.CreatedAt).
			Set("archived_at", i.ArchivedAt).
			Set("last_aggregated_at", i.LastAggregatedAt).
			Set("exported_at", i.ExportedAt).
			Returning("*").To(&i)

		err := self.database.Query(ctx, stmt)
		if err != nil {
			return err
		}

		stmt = sqlf.
			InsertInto(ISSUE_FEEDBACK_MODEL_TABLE).
			Set("issue_id", issue.ID).
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

	return i.ToEntity(), nil
}

func (self *IssueRepository) GetByID(ctx context.Context, id string) (*Issue, error) {
	var i IssueModel

	stmt := sqlf.
		Select("*").To(&i).
		From(ISSUE_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return i.ToEntity(), nil
}

func (self *IssueRepository) GetByIDForUpdate(ctx context.Context, id string) (*Issue, error) {
	var i IssueModel

	stmt := sqlf.
		Select("*").To(&i).
		From(ISSUE_MODEL_TABLE).
		Where("id = ?", id).
		Clause("FOR NO KEY UPDATE NOWAIT")

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return i.ToEntity(), nil
}

func (self *IssueRepository) ListByProductID(ctx context.Context,
	productID string, search IssueSearch) (*util.Page[Issue, IssueSearchCursor], error) {
	var is []IssueModel

	stmt := sqlf.
		Select("*").To(&is).
		From(ISSUE_MODEL_TABLE).
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

	if search.Filters.Severities != nil && len(*search.Filters.Severities) > 0 {
		stmt.
			Where("").
			SubQuery("(", ")", sqlf.
				Select("severity.key").
				From("jsonb_each_text(severities) as severity").
				OrderBy("severity.value::int DESC", "severity.key DESC").
				Limit(1)).
			In(util.Spread(*search.Filters.Severities)...)
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
			return value == IssueSearchFiltersAssigneesUnassigned
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
		case IssueSearchFiltersStatusActive:
			stmt.
				Where("(archived_at IS NULL OR (archived_at IS NOT NULL AND last_seen_at > archived_at))")
		case IssueSearchFiltersStatusRegressed:
			stmt.
				Where("archived_at IS NOT NULL").
				Where("last_seen_at > archived_at")
		case IssueSearchFiltersStatusArchived:
			stmt.
				Where("archived_at IS NOT NULL")
		case IssueSearchFiltersStatusUnarchived:
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
		if search.Orders.Relevance == IssueSearchOrdersAscending {
			stmt.
				Where("(priority, last_seen_at, id) > (?, ?, ?)", search.Pagination.From.Value.Priority,
					search.Pagination.From.Value.SeenAt, search.Pagination.From.ID)
		} else {
			stmt.
				Where("(priority, last_seen_at, id) < (?, ?, ?)", search.Pagination.From.Value.Priority,
					search.Pagination.From.Value.SeenAt, search.Pagination.From.ID)
		}
	}

	if search.Orders.Relevance == IssueSearchOrdersAscending {
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
			return &util.Page[Issue, IssueSearchCursor]{}, nil
		}

		return nil, err
	}

	items := make([]Issue, 0, len(is))
	for _, i := range is {
		items = append(items, *i.ToEntity())
	}

	var cursor *util.Cursor[IssueSearchCursor]
	if len(items) == search.Pagination.Limit {
		cursor = &util.Cursor[IssueSearchCursor]{
			Value: IssueSearchCursor{
				Priority: items[search.Pagination.Limit-1].Priority,
				SeenAt:   items[search.Pagination.Limit-1].LastSeenAt,
			},
			ID: items[search.Pagination.Limit-1].ID,
		}
	}

	return &util.Page[Issue, IssueSearchCursor]{
		Items: items,
		Next:  cursor,
	}, nil
}

func (self *IssueRepository) ListByEmbeddingAndProductID(ctx context.Context,
	embedding []float32, threshold float64, limit int, productID string) ([]Issue, error) {
	var result []struct {
		IssueModel
		Score float64 `db:"score"`
	}

	stmt := sqlf.
		Select("*").To(&result).
		From("").
		SubQuery("(", ")", sqlf.
			Select("*").
			Select("1 - (embedding <=> ?::vector) AS score", pgvector.NewVector(embedding)).
			From(ISSUE_MODEL_TABLE).
			Where("product_id = ?", productID).
			OrderBy("score DESC").
			Limit(limit)).
		Where("score >= ?", threshold)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []Issue{}, nil
		}

		return nil, err
	}

	entities := make([]Issue, 0, len(result))
	for _, res := range result {
		entities = append(entities, *res.IssueModel.ToEntity())
	}

	return entities, nil
}

func (self *IssueRepository) ListFeedbacks(ctx context.Context,
	id string, pagination util.Pagination[time.Time]) (*util.Page[feedback.Feedback, time.Time], error) {
	var fs []feedback.FeedbackModel

	stmt := sqlf.
		Select(feedback.FEEDBACK_MODEL_TABLE+".*").To(&fs).
		From(feedback.FEEDBACK_MODEL_TABLE).
		Join(ISSUE_FEEDBACK_MODEL_TABLE,
			ISSUE_FEEDBACK_MODEL_TABLE+".feedback_id = "+feedback.FEEDBACK_MODEL_TABLE+".id").
		Where(ISSUE_FEEDBACK_MODEL_TABLE+".issue_id = ?", id)

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

func (self *IssueRepository) UpdateAssignee(ctx context.Context, id string, assigneeID *string) error {
	stmt := sqlf.
		Update(ISSUE_MODEL_TABLE).
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

func (self *IssueRepository) UpdateQuality(ctx context.Context, id string, quality int) error {
	stmt := sqlf.
		Update(ISSUE_MODEL_TABLE).
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

func (self *IssueRepository) UpdateArchivedAt(ctx context.Context, id string, archivedAt *time.Time) error {
	stmt := sqlf.
		Update(ISSUE_MODEL_TABLE).
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

func (self *IssueRepository) UpdateAggregated(ctx context.Context, issue Issue,
	feedback feedback.Feedback) error {
	i := NewIssueModel(issue)

	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		stmt := sqlf.
			Update(ISSUE_MODEL_TABLE).
			Set("embedding", i.Embedding).
			Set("sources", i.Sources).
			Set("title", i.Title).
			Set("description", i.Description).
			Set("steps", i.Steps).
			Set("severities", i.Severities).
			Set("priority", i.Priority).
			Set("categories", i.Categories).
			Set("releases", i.Releases).
			Set("customers", i.Customers).
			Set("first_seen_at", i.FirstSeenAt).
			Set("last_seen_at", i.LastSeenAt).
			Set("last_aggregated_at", i.LastAggregatedAt).
			Where("id = ?", i.ID)

		affected, err := self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		if affected != 1 {
			return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
		}

		stmt = sqlf.
			InsertInto(ISSUE_FEEDBACK_MODEL_TABLE).
			Set("issue_id", issue.ID).
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

type PartialIssueRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewPartialIssueRepository(observer *kit.Observer, database *kit.Database,
	config config.Config) *PartialIssueRepository {
	return &PartialIssueRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *PartialIssueRepository) Create(ctx context.Context, partial PartialIssue) (*PartialIssue, error) {
	p := NewPartialIssueModel(partial)

	stmt := sqlf.
		InsertInto(PARTIAL_ISSUE_MODEL_TABLE).
		Set("id", p.ID).
		Set("feedback_id", p.FeedbackID).
		Set("title", p.Title).
		Set("description", p.Description).
		Set("steps", p.Steps).
		Set("severity", p.Severity).
		Set("category", p.Category).
		Set("created_at", p.CreatedAt).
		Returning("*").To(&p)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return p.ToEntity(), nil
}

func (self *PartialIssueRepository) BulkCreate(ctx context.Context, partials []PartialIssue) error {
	if len(partials) == 0 {
		return nil
	}

	stmt := sqlf.
		InsertInto(PARTIAL_ISSUE_MODEL_TABLE)

	for _, partial := range partials {
		p := NewPartialIssueModel(partial)

		stmt.
			NewRow().
			Set("id", p.ID).
			Set("feedback_id", p.FeedbackID).
			Set("title", p.Title).
			Set("description", p.Description).
			Set("steps", p.Steps).
			Set("severity", p.Severity).
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

func (self *PartialIssueRepository) GetByID(ctx context.Context, id string) (*PartialIssue, error) {
	var p PartialIssueModel

	stmt := sqlf.
		Select("*").To(&p).
		From(PARTIAL_ISSUE_MODEL_TABLE).
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

func (self *PartialIssueRepository) ListIDsByCreatedAt(ctx context.Context,
	pagination util.Pagination[time.Time]) (*util.Page[string, time.Time], error) {
	var result []struct {
		ID        string    `db:"id"`
		CreatedAt time.Time `db:"created_at"`
	}

	stmt := sqlf.
		Select("id, created_at").To(&result).
		From(PARTIAL_ISSUE_MODEL_TABLE)

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

func (self *PartialIssueRepository) Delete(ctx context.Context, id string) error {
	stmt := sqlf.
		DeleteFrom(PARTIAL_ISSUE_MODEL_TABLE).
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
