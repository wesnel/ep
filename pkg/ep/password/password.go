package password

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
)

var expected = map[string]setter{
	"host": func(
		ctx context.Context,
		p *password,
		s string,
	) error {
		slog.LogAttrs(ctx, slog.LevelDebug, "Received value for host",
			slog.String("value", s))

		p.Host = s

		return nil
	},
	"user": func(
		ctx context.Context,
		p *password,
		s string,
	) error {
		slog.LogAttrs(ctx, slog.LevelDebug, "Received value for user",
			slog.String("value", s))

		p.User = s

		return nil
	},
	"port": func(
		ctx context.Context,
		p *password,
		s string,
	) error {
		slog.LogAttrs(ctx, slog.LevelDebug, "Received value for port",
			slog.String("value", s))

		if s == "" {
			return nil
		}

		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		if i < 0 {
			return fmt.Errorf("port %d is not a non-negative number", i)
		}

		p.Port = uint64(i)

		return nil
	},
}

type setter func(
	ctx context.Context,
	p *password,
	s string,
) error

type password struct {
	Host string
	User string
	Port uint64
}

func Parse(
	ctx context.Context,
	pattern *regexp.Regexp,
	s string,
) (*password, error) {
	p := &password{}
	names := pattern.SubexpNames()
	matches := pattern.FindStringSubmatch(s)

	for name, f := range expected {
		if !slices.Contains(names, name) {
			return nil, fmt.Errorf("missing regexp subexp with name %s", name)
		}

		i := pattern.SubexpIndex(name)

		if i < 0 {
			return nil, fmt.Errorf("missing regexp subexp with name %s", name)
		}

		if err := f(ctx, p, matches[i]); err != nil {
			return nil, err
		}
	}

	return p, nil
}
