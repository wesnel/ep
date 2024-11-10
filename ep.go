package main

import (
	"context"
	"ep/pkg/ep/command/emacs"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var verbose bool

func init() {
	flags()
	logging()
}

func flags() {
	flag.Usage = usage
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.Parse()
}

func usage() {
	var b strings.Builder
	flag.CommandLine.SetOutput(&b)

	fmt.Fprintf(&b, "ep: Emacs passwords\n\n")
	fmt.Fprintf(&b, "The following flags are available in all commands:\n")
	flag.PrintDefaults()
	fmt.Fprintf(&b, "Usage:\n")
	fmt.Fprintf(&b, "  %s [show] pass-name\n", os.Args[0])
	fmt.Fprintf(&b, "    Show existing password.\n")

	fmt.Fprintf(os.Stderr, b.String())
}

func logging() {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})))
}

func main() {
	ctx := context.Background()
	args := flag.Args()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		attrs := []slog.Attr{
			slog.Any("positional", args),
		}

		flag.VisitAll(func(f *flag.Flag) {
			attrs = append(attrs,
				slog.String(f.Name, f.Value.String()))
		})

		slog.LogAttrs(ctx, slog.LevelDebug, "Args", attrs...)
	}

	cmd, err := emacs.New(ctx, args...)
	if err != nil {
		panic(err.Error())
	}

	stdout, err := cmd.Run(ctx)
	if err != nil {
		panic(fmt.Sprintf(`Error while executing "%s": "%s"`,
			cmd.String(),
			err.Error()))
	}

	fmt.Fprintln(os.Stdout, stdout)
}
