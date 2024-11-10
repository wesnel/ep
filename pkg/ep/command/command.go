package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"time"
)

type Stringer interface {
	String() string
}

type Runner interface {
	Stringer

	Run(context.Context) (string, error)
}

type Option func(*command)

func WithTimeout(timeout time.Duration) Option {
	return func(cmd *command) {
		cmd.timeout = timeout
	}
}

func WithName(name string) Option {
	return func(cmd *command) {
		cmd.name = name
	}
}

func WithArgs(args ...string) Option {
	return func(cmd *command) {
		cmd.args = append(cmd.args, args...)
	}
}

func WithCleanup(funcs ...func(context.Context) error) Option {
	return func(cmd *command) {
		cmd.cleaners = append(cmd.cleaners, funcs...)
	}
}

func New(
	ctx context.Context,
	opts ...Option,
) (*command, error) {
	cmd := &command{}

	for _, opt := range opts {
		opt(cmd)
	}

	if err := cmd.validate(ctx); err != nil {
		cmd.cleanup(ctx)
		return cmd, fmt.Errorf("invalid command: %w",
			err)
	}

	return cmd, nil
}

type command struct {
	name       string
	args       []string
	validators []func(ctx context.Context, name string, args ...string) error
	cleaners   []func(context.Context) error
	timeout    time.Duration
}

func (cmd *command) Run(ctx context.Context) (string, error) {
	defer cmd.cleanup(ctx)

	if cmd.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cmd.timeout)
		defer cancel()
	}

	slog.LogAttrs(ctx, slog.LevelDebug, "Running command",
		slog.String("command", cmd.String()))
	stdout, err := exec.
		CommandContext(ctx, cmd.name, cmd.args...).
		Output()
	return string(stdout), err
}

func (cmd *command) String() string {
	return exec.Command(cmd.name, cmd.args...).String()
}

func (cmd *command) validate(
	ctx context.Context,
) error {
	errs := make([]error, 0, len(cmd.validators))

	for _, validator := range cmd.validators {
		errs = append(errs,
			validator(ctx, cmd.name, cmd.args...))
	}

	return errors.Join(errs...)
}

func (cmd *command) cleanup(
	ctx context.Context,
) {
	errs := make([]error, 0, len(cmd.cleaners))

	for _, cleanup := range cmd.cleaners {
		errs = append(errs,
			cleanup(ctx))
	}

	if err := errors.Join(errs...); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "Error cleanup up after command",
			slog.String("command", cmd.String()),
			slog.String("error", err.Error()))
	}
}
