package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// Set from package.json by the build script.
var Version = "dev"

const banner = `
██████╗ ███████╗██╗   ██╗███████╗██████╗  █████╗  ██████╗███████╗
██╔══██╗██╔════╝██║   ██║██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██║  ██║█████╗  ██║   ██║███████╗██████╔╝███████║██║     █████╗
██║  ██║██╔══╝  ╚██╗ ██╔╝╚════██║██╔═══╝ ██╔══██║██║     ██╔══╝
██████╔╝███████╗ ╚████╔╝ ███████║██║     ██║  ██║╚██████╗███████╗
╚═════╝ ╚══════╝  ╚═══╝  ╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚══════╝
`

const (
	colorYellow = "\x1b[93m"
	colorWhite  = "\x1b[97m"
	colorCyan   = "\x1b[1;96m"
	colorMuted  = "\x1b[90m"
	colorReset  = "\x1b[0m"
)

var RootCmd = &cobra.Command{
	Use:   "devspace",
	Short: "Scaffold applications, Dockerfiles and Kubernetes manifests",
	Long: `Build your next application from an interactive terminal wizard.

WHAT DEVSPACE DOES
  + Creates frontend and backend project structures.
  + Generates Dockerfiles and Kubernetes manifests.
  + Lets you choose the project folder and service name.
  + Configures replicas, CPU, memory and service exposure.

SUPPORTED STACKS
  Backend   Go | Django | Express | Actix-web
  Frontend  React | Next.js | Vue | Svelte | Backend only

REQUIREMENTS
  The npm launcher requires Node.js 22 or newer.
  Install the tools needed by your chosen framework.
  Docker and kubectl are needed when building or deploying.
  Local scaffolding does not require a GitHub token.`,
	Example: `  devspace                     Show the animated welcome screen
  devspace create              Start the project wizard
  devspace create --help       Explain every wizard setting
  devspace --version           Show the installed version
  devspace completion --help   Set up shell completions`,
	Version:       Version,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		printWelcome(cmd.OutOrStdout())
		return nil
	},
}

func init() {
	defaultHelp := RootCmd.HelpFunc()

	RootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		out := cmd.OutOrStdout()

		var buffer strings.Builder
		cmd.SetOut(&buffer)

		func() {
			defer cmd.SetOut(out)
			defaultHelp(cmd, args)
		}()

		printBanner(out)
		printStyledHelp(out, buffer.String())
	})

	RootCmd.SetVersionTemplate("Devspace {{.Version}}\n")

	// The installer prints a welcome screen without waiting for input.
	RootCmd.AddCommand(&cobra.Command{
		Use:    "welcome",
		Hidden: true,
		Args:   cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			printStaticWelcome(cmd.OutOrStdout())
		},
	})

	createCmd.Short = "Create a project using the interactive wizard"
	createCmd.Args = cobra.NoArgs

	createCmd.Long = `Create a project by answering the following questions.

PROJECT
  Location      Select an existing folder.
                Enter or . uses the current terminal folder.
                Quoted paths are accepted.
                Use ~ for your home folder.

  Service name  Use 1-46 characters.
                Start with a letter and end with a letter or number.
                Hyphens are allowed. Example: trial-1.
                Spaces become hyphens; capitals become lowercase.
                Windows reserved names are rejected.

FRAMEWORKS
  Backend       Go, Python/Django, Node.js/Express or Rust/Actix-web.
  Frontend      React, Next.js, Vue, Svelte or backend only.

KUBERNETES
  Replicas      Number of application instances; minimum 1.
  CPU           CPU request per pod, such as 100m or 250m.
  Memory        Memory request per pod, such as 128Mi or 256Mi.
  Service       ClusterIP, NodePort or LoadBalancer.

OUTPUT
  Files are generated inside <location>/<service-name>.
  Framework setup may download dependencies from the internet.
  Review generated manifests before deployment.

CONTROLS
  Arrow keys    Move through choices.
  Enter         Confirm your answer.
  Ctrl+C        Cancel the current prompt.`

	createCmd.Example = `  devspace create

  Example answers:
    Location:      .
    Service name:  trial-1
    Backend:       Node.js (Express)
    Frontend:      React (Vite)`
}

// Enable colors in supported terminals, including Windows.
func terminalOutput(w io.Writer) (io.Writer, bool) {
	file, ok := w.(*os.File)
	if !ok {
		return w, false
	}

	_, noColor := os.LookupEnv("NO_COLOR")
	if noColor || os.Getenv("TERM") == "dumb" {
		return w, false
	}

	if !isatty.IsTerminal(file.Fd()) &&
		!isatty.IsCygwinTerminal(file.Fd()) {
		return w, false
	}

	return colorable.NewColorable(file), true
}

// Combine the SD card and title into aligned rows.
func bannerLines() []string {
	card := []string{
		"  ╭─────┐   ",
		"  │     └─╮ ",
		"  │ ━━━━━ │ ",
		"  │ ━━━   │ ",
		"  │       │ ",
		"  ╰───────╯ ",
	}

	var title []string

	for _, line := range strings.Split(banner, "\n") {
		if strings.TrimSpace(line) != "" {
			title = append(title, line)
		}
	}

	rows := len(card)
	if len(title) > rows {
		rows = len(title)
	}

	lines := make([]string, rows)

	for i := 0; i < rows; i++ {
		iconLine, titleLine := "", ""

		if i < len(card) {
			iconLine = card[i]
		}

		if i < len(title) {
			titleLine = title[i]
		}

		lines[i] = fmt.Sprintf("%-13s%s", iconLine, titleLine)
	}

	return lines
}

func printBrandFooter(w io.Writer, useColor bool) {
	muted, reset := "", ""

	if useColor {
		muted = colorMuted
		reset = colorReset
	}

	fmt.Fprintf(
		w,
		"\n%s  Full-stack scaffolding | Devspace %s%s\n",
		muted, Version, reset,
	)

	fmt.Fprintf(
		w,
		"%s  --------------------------------------------------%s\n\n",
		muted, reset,
	)
}

func printBanner(w io.Writer) {
	out, useColor := terminalOutput(w)
	yellow, reset := "", ""

	if useColor {
		yellow = colorYellow
		reset = colorReset
	}

	fmt.Fprintln(out)

	for _, line := range bannerLines() {
		fmt.Fprintln(out, yellow+line+reset)
	}

	printBrandFooter(out, useColor)
}

func printCommands(w io.Writer, useColor bool) {
	cyan, reset := "", ""

	if useColor {
		cyan = colorCyan
		reset = colorReset
	}

	fmt.Fprintf(w, "  %sGET STARTED%s\n\n", cyan, reset)
	fmt.Fprintln(w, "    devspace create       Create a project")
	fmt.Fprintln(w, "    devspace --help       View commands and information")
	fmt.Fprintln(w, "    devspace --version    Show installed version")
	fmt.Fprintln(w)
}

func printStaticWelcome(w io.Writer) {
	printBanner(w)

	out, useColor := terminalOutput(w)
	printCommands(out, useColor)
}

// Animate until the user presses Enter or Ctrl+C.
func printWelcome(w io.Writer) {
	out, useColor := terminalOutput(w)

	inputIsTerminal := isatty.IsTerminal(os.Stdin.Fd()) ||
		isatty.IsCygwinTerminal(os.Stdin.Fd())

	if !useColor || !inputIsTerminal {
		printStaticWelcome(w)
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	enterPressed := make(chan struct{})

	go func() {
		var input [1]byte

		for {
			n, err := os.Stdin.Read(input[:])

			if err != nil ||
				(n > 0 && (input[0] == '\n' || input[0] == '\r')) {
				close(enterPressed)
				return
			}
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	frame := 0
	screen := sparkleScreen(frame)
	lineCount := strings.Count(screen, "\n")

	fmt.Fprint(out, screen)
	defer fmt.Fprint(out, colorReset, "\n")

	for {
		select {
		case <-ticker.C:
			frame++

			fmt.Fprintf(out, "\x1b[%dA\r", lineCount)
			fmt.Fprint(out, sparkleScreen(frame))

		case <-enterPressed:
			return

		case <-interrupt:
			return
		}
	}
}

func sparkleScreen(frame int) string {
	lines := bannerLines()

	width := 0
	for _, line := range lines {
		if length := len([]rune(line)); length > width {
			width = length
		}
	}

	// Leave a small pause between highlight sweeps.
	shimmerColumn := (frame * 2) % (width + 12)

	var screen strings.Builder
	screen.WriteString("\n")

	for _, line := range lines {
		// Rune indexes keep Unicode block characters aligned.
		for column, character := range []rune(line) {
			if character == ' ' {
				screen.WriteByte(' ')
				continue
			}

			inShimmer := column >= shimmerColumn &&
				column < shimmerColumn+3

			if inShimmer || rand.Intn(100) < 7 {
				screen.WriteString(colorWhite)
			} else {
				screen.WriteString(colorYellow)
			}

			screen.WriteRune(character)
		}

		screen.WriteString(colorReset)
		screen.WriteString("\x1b[K\n")
	}

	printBrandFooter(&screen, true)
	printCommands(&screen, true)

	fmt.Fprintf(
		&screen,
		"%s  Press Enter or Ctrl+C to return to your terminal.%s\n",
		colorMuted, colorReset,
	)

	return screen.String()
}

func printStyledHelp(w io.Writer, helpText string) {
	out, useColor := terminalOutput(w)

	if !useColor {
		fmt.Fprint(out, helpText)
		return
	}

	headings := map[string]bool{
		"WHAT DEVSPACE DOES":      true,
		"SUPPORTED STACKS":        true,
		"REQUIREMENTS":            true,
		"Usage:":                  true,
		"Aliases:":                true,
		"Examples:":               true,
		"Available Commands:":     true,
		"Flags:":                  true,
		"Global Flags:":           true,
		"Additional help topics:": true,
		"PROJECT":                 true,
		"FRAMEWORKS":              true,
		"KUBERNETES":              true,
		"OUTPUT":                  true,
		"CONTROLS":                true,
	}

	lines := strings.Split(helpText, "\n")

	for i, line := range lines {
		if headings[strings.TrimSpace(line)] {
			lines[i] = colorCyan + line + colorReset
		}
	}

	fmt.Fprint(out, strings.Join(lines, "\n"))
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(RootCmd.ErrOrStderr(), "Error:", err)
		fmt.Fprintln(
			RootCmd.ErrOrStderr(),
			"Run 'devspace --help' for usage.",
		)
		os.Exit(1)
	}
}
