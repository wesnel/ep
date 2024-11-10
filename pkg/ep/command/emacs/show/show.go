package show

import (
	"context"
	"ep/pkg/ep/command/emacs/show/elisp"
	"ep/pkg/ep/password"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var format = regexp.MustCompile(`((?P<user>\w+)@)?(?P<host>.+)(:(?P<port>\d+))?`)

func Build(
	ctx context.Context,
	args ...string,
) (*os.File, error) {
	s, err := name(ctx, args...)
	if err != nil {
		return nil, err
	}

	p, err := password.Parse(ctx, format, s)
	if err != nil {
		return nil, err
	}

	var b strings.Builder
	if err := elisp.Show.Execute(&b, p); err != nil {
		return nil, fmt.Errorf("error building elisp: %w",
			err)
	}

	f, err := os.CreateTemp("", "show.*.el")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary elisp file: %w",
			err)
	}

	_, err = f.WriteString(b.String())
	if err != nil {
		return nil, fmt.Errorf("error writing temporary elisp file: %w",
			err)
	}

	return f, nil
}

func name(
	_ context.Context,
	args ...string,
) (string, error) {
	switch {
	case len(args) == 1:
		return args[0], nil
	case len(args) == 2 && args[0] == "show":
		return args[1], nil
	default:
		return "", errors.New("password name incorrectly passed as argument")
	}
}
