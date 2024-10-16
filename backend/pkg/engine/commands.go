package engine

import (
	"context"
	"time"

	"github.com/mkideal/cli"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

const (
	EngineCommandsOpenBreaker  = "open-breaker"
	EngineCommandsCloseBreaker = "close-breaker"
)

type EngineCommands struct {
	config        config.Config
	observer      *kit.Observer
	engineBreaker *EngineBreaker
}

func NewEngineCommands(observer *kit.Observer, engineBreaker *EngineBreaker, config config.Config) *EngineCommands {
	return &EngineCommands{
		config:        config,
		observer:      observer,
		engineBreaker: engineBreaker,
	}
}

type EngineCommandsOpenBreakerArgs struct {
	cli.Helper
	Timeout int `cli:"*timeout" usage:"how long the circuit is opened in seconds"`
}

func (self *EngineCommands) OpenBreaker(ctx context.Context, command *cli.Context) error {
	args, ok := command.Argv().(*EngineCommandsOpenBreakerArgs)
	if !ok {
		return kit.ErrRunnerGeneric.Raise().With("cannot get command arguments")
	}

	err := self.engineBreaker.Force(ctx, time.Duration(args.Timeout)*time.Second)
	if err != nil {
		return err
	}

	return nil
}

type EngineCommandsCloseBreakerArgs struct {
	cli.Helper
}

func (self *EngineCommands) CloseBreaker(ctx context.Context, command *cli.Context) error {
	_, ok := command.Argv().(*EngineCommandsCloseBreakerArgs)
	if !ok {
		return kit.ErrRunnerGeneric.Raise().With("cannot get command arguments")
	}

	err := self.engineBreaker.Close(ctx)
	if err != nil {
		return err
	}

	return nil
}
