package logger

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/Ladicle/tabwriter"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"

	"github.com/vikbert/taskr/v3/errors"
	"github.com/vikbert/taskr/v3/experiments"
	"github.com/vikbert/taskr/v3/internal/term"
	"github.com/vikbert/taskr/v3/internal/version"
)

var (
	ErrPromptCancelled = errors.New("prompt cancelled")
	ErrNoTerminal      = errors.New("no terminal")
)

// createPtermPrintFunc creates a PrintFunc that uses pterm for coloring
func createPtermPrintFunc(style *pterm.Style) PrintFunc {
	return func(w io.Writer, format string, args ...any) {
		text := fmt.Sprintf(format, args...)
		coloredText := style.Sprint(text)
		fmt.Fprint(w, coloredText)
	}
}

type (
	Color     func() PrintFunc
	PrintFunc func(io.Writer, string, ...any)
)

func Default() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle())
}

func Blue() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgBlue))
}

func Green() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgGreen))
}

func Cyan() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgCyan))
}

func Yellow() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgYellow))
}

func Magenta() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgMagenta))
}

func Red() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgRed))
}

func BrightBlue() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgLightBlue))
}

func BrightGreen() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgLightGreen))
}

func BrightCyan() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgLightCyan))
}

func BrightYellow() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgLightYellow))
}

func BrightMagenta() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgLightMagenta))
}

func BrightRed() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgLightRed))
}

func BoldYellow() PrintFunc {
	return createPtermPrintFunc(pterm.NewStyle(pterm.FgYellow, pterm.Bold))
}

// Logger is just a wrapper that prints stuff to STDOUT or STDERR,
// with optional color.
type Logger struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	Verbose    bool
	Color      bool
	AssumeYes  bool
	AssumeTerm bool // Used for testing
}

// Outf prints stuff to STDOUT.
func (l *Logger) Outf(color Color, s string, args ...any) {
	l.FOutf(l.Stdout, color, s, args...)
}

// FOutf prints stuff to the given writer.
func (l *Logger) FOutf(w io.Writer, color Color, s string, args ...any) {
	if len(args) == 0 {
		s, args = "%s", []any{s}
	}
	if !l.Color {
		color = Default
	}
	print := color()
	print(w, s, args...)
}

// VerboseOutf prints stuff to STDOUT if verbose mode is enabled.
func (l *Logger) VerboseOutf(color Color, s string, args ...any) {
	if l.Verbose {
		l.Outf(color, s, args...)
	}
}

// Errf prints stuff to STDERR.
func (l *Logger) Errf(color Color, s string, args ...any) {
	if len(args) == 0 {
		s, args = "%s", []any{s}
	}
	if !l.Color {
		color = Default
	}
	print := color()
	print(l.Stderr, s, args...)
}

// VerboseErrf prints stuff to STDERR if verbose mode is enabled.
func (l *Logger) VerboseErrf(color Color, s string, args ...any) {
	if l.Verbose {
		l.Errf(color, s, args...)
	}
}

func (l *Logger) Warnf(message string, args ...any) {
	l.Errf(Yellow, message, args...)
}

func (l *Logger) Prompt(color Color, prompt string, defaultValue string, continueValues ...string) error {
	if l.AssumeYes {
		l.Outf(color, "%s [assuming yes]\n", prompt)
		return nil
	}

	if !l.AssumeTerm && !term.IsTerminal() {
		return ErrNoTerminal
	}

	if len(continueValues) == 0 {
		return errors.New("no continue values provided")
	}

	l.Outf(color, "%s [%s/%s]: ", prompt, strings.ToLower(continueValues[0]), strings.ToUpper(defaultValue))

	reader := bufio.NewReader(l.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if !slices.Contains(continueValues, input) {
		return ErrPromptCancelled
	}

	return nil
}

func (l *Logger) PrintExperiments() error {
	w := tabwriter.NewWriter(l.Stdout, 0, 8, 0, ' ', 0)
	for _, x := range experiments.List() {
		if !x.Active() {
			continue
		}
		l.FOutf(w, Yellow, "* ")
		l.FOutf(w, Green, x.Name)
		l.FOutf(w, Default, ": \t%s\n", x.String())
	}
	return w.Flush()
}

func (l *Logger) PrintBanner() {
	l.PrintBannerWithProject("")
}

func (l *Logger) PrintBannerWithProject(project string) {
	// Print empty line before banner
	fmt.Println()

	if project != "" {
		_ = pterm.DefaultBigText.WithLetters(
			putils.LettersFromStringWithStyle(project, pterm.FgCyan.ToStyle()),
		).Render()
	} else {
		_ = pterm.DefaultBigText.WithLetters(
			putils.LettersFromStringWithStyle("TASK", pterm.FgCyan.ToStyle()),
			putils.LettersFromStringWithStyle("R", pterm.FgLightYellow.ToStyle()),
		).Render()
	}

	// Print version information
	l.Outf(BrightCyan, "Taskr: v%s\n", version.GetVersionWithBuildInfo())

	// Print empty line after banner
	fmt.Println()
}
