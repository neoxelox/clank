package organization

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/util"
)

const (
	OrganizationTasksDowngradeEndedTrials = "organization:downgrade-ended-trials"
	OrganizationTasksComputeUsage         = "organization:compute-usage"
	OrganizationTasksScheduleComputeUsage = "organization:schedule-compute-usage"
)

type OrganizationTasks struct {
	config                 config.Config
	observer               *kit.Observer
	database               *kit.Database
	organizationRepository OrganizationRepository
	enqueuer               *kit.Enqueuer
}

func NewOrganizationTasks(observer *kit.Observer, database *kit.Database,
	organizationRepository OrganizationRepository, enqueuer *kit.Enqueuer,
	config config.Config) *OrganizationTasks {
	return &OrganizationTasks{
		config:                 config,
		observer:               observer,
		database:               database,
		organizationRepository: organizationRepository,
		enqueuer:               enqueuer,
	}
}

func (self *OrganizationTasks) DowngradeEndedTrials(ctx context.Context, _ *asynq.Task) error {
	err := self.organizationRepository.DowngradeTrialsByEndsAt(ctx, time.Now())
	if err != nil {
		return err
	}

	return nil
}

type OrganizationTasksComputeUsageParams struct {
	OrganizationID string
}

func (self *OrganizationTasks) ComputeUsage(ctx context.Context, task *asynq.Task) error {
	params := OrganizationTasksComputeUsageParams{}

	err := json.Unmarshal(task.Payload(), &params)
	if err != nil {
		self.observer.Error(ctx, kit.ErrWorkerGeneric.Raise().Cause(err))
		return nil
	}

	usagePerProduct, err := self.organizationRepository.CountCollectedFeedbacksThisMonthPerProduct(ctx, params.OrganizationID)
	if err != nil {
		return err
	}

	usage := 0
	for _, count := range usagePerProduct {
		usage += count
	}

	err = self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		organization, err := self.organizationRepository.GetByIDForUpdate(ctx, params.OrganizationID)
		if err != nil {
			return err
		}

		if organization == nil {
			return nil
		}

		if organization.DeletedAt != nil {
			return nil
		}

		now := time.Now()
		if organization.Usage.LastComputedAt.Before(util.StartOfDay(util.StartOfMonth(now))) {
			organization.Capacity.Extra -= min(max(0, organization.Usage.Value-organization.Capacity.Included), organization.Capacity.Extra)
		}

		organization.Usage.Value = usage
		organization.Usage.LastComputedAt = now

		err = self.organizationRepository.UpdateOrganizationUsage(ctx, *organization)
		if err != nil {
			return err
		}

		if len(usagePerProduct) > 0 {
			err = self.organizationRepository.UpdateProductsUsage(ctx, usagePerProduct)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (self *OrganizationTasks) ScheduleComputeUsage(ctx context.Context, _ *asynq.Task) error {
	ids, err := self.organizationRepository.ListIDsByNotDeleted(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {
		err := self.enqueuer.Enqueue(ctx, OrganizationTasksComputeUsage, OrganizationTasksComputeUsageParams{
			OrganizationID: id,
		}, asynq.Queue("critical"))
		if err != nil {
			self.observer.Error(ctx, err)
		}
	}

	return nil
}
