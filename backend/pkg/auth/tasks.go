package auth

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

const (
	AuthTasksDeleteExpiredSignInCodes = "auth:delete-expired-sign-in-codes"
)

type AuthTasks struct {
	config               config.Config
	observer             *kit.Observer
	signInCodeRepository SignInCodeRepository
}

func NewAuthTasks(observer *kit.Observer, signInCodeRepository SignInCodeRepository,
	config config.Config) *AuthTasks {
	return &AuthTasks{
		config:               config,
		observer:             observer,
		signInCodeRepository: signInCodeRepository,
	}
}

func (self *AuthTasks) DeleteExpiredSignInCodes(ctx context.Context, _ *asynq.Task) error {
	err := self.signInCodeRepository.DeleteByExpiresAt(ctx, time.Now())
	if err != nil {
		return err
	}

	return nil
}
