package user

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

const (
	UserTasksDeleteExpiredInvitations = "user:delete-expired-invitations"
)

type UserTasks struct {
	config               config.Config
	observer             *kit.Observer
	invitationRepository InvitationRepository
}

func NewUserTasks(observer *kit.Observer, invitationRepository InvitationRepository,
	config config.Config) *UserTasks {
	return &UserTasks{
		config:               config,
		observer:             observer,
		invitationRepository: invitationRepository,
	}
}

func (self *UserTasks) DeleteExpiredInvitations(ctx context.Context, _ *asynq.Task) error {
	err := self.invitationRepository.DeleteByExpiresAt(ctx, time.Now())
	if err != nil {
		return err
	}

	return nil
}
