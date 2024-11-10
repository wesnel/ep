package emacs

import (
	"context"
	"ep/pkg/ep/command"
	"ep/pkg/ep/command/emacs/show"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

func New(
	ctx context.Context,
	args ...string,
) (command.Runner, error) {
	opts, err := options(ctx, args...)
	if err != nil {
		return nil, err
	}

	return command.New(ctx, opts...)
}

type builder func(
	ctx context.Context,
	args ...string,
) (*os.File, error)

func options(
	ctx context.Context,
	args ...string,
) ([]command.Option, error) {
	opts, err := defaults(ctx)
	if err != nil {
		return nil, err
	}

	b := interpret(ctx, args...)

	f, err := b(ctx, args...)
	if err != nil {
		remove(f)
		close(f)
		return nil, err
	}

	opts = append(opts,
		command.WithCleanup(remove(f)),
		command.WithCleanup(close(f)),
		command.WithArgs(
			"--script",
			f.Name()))

	return opts, nil
}

func defaults(
	ctx context.Context,
) ([]command.Option, error) {
	opts := make([]command.Option, 0)

	path, err := exec.LookPath("emacs")
	if err != nil {
		return nil, fmt.Errorf("cannot find emacs: %w",
			err)
	}

	slog.LogAttrs(ctx, slog.LevelDebug, "Found Emacs executable",
		slog.String("path", path))

	opts = append(opts,
		command.WithName(path),
		command.WithArgs(
			"--quick",
			"--batch"))

	return opts, nil
}

func interpret(
	_ context.Context,
	args ...string,
) builder {
	command := "ls"
	if len(args) > 0 {
		command = args[0]
	}

	switch command {
	case "show":
		fallthrough
	default:
		return show.Build
	}
}

func remove(f *os.File) func(context.Context) error {
	return func(_ context.Context) error {
		if f == nil {
			return nil
		}

		return os.Remove(f.Name())
	}
}

func close(f *os.File) func(context.Context) error {
	return func(_ context.Context) error {
		if f == nil {
			return nil
		}

		return f.Close()
	}
}
