package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// Colors
var (
	bold      = color.New(color.Bold)
	dim       = color.New(color.FgHiBlack)
	green     = color.New(color.FgGreen)
	red       = color.New(color.FgRed)
	yellow    = color.New(color.FgYellow)
	cyan      = color.New(color.FgCyan)
	boldGreen = color.New(color.Bold, color.FgGreen)
	boldRed   = color.New(color.Bold, color.FgRed)
	boldCyan  = color.New(color.Bold, color.FgCyan)
)

func printLogo() {
	boldCyan.Println(`
  ╔═╗╔╗╔╦  ╦╔═╗╦ ╦╦ ╔╦╗
  ║╣ ║║║╚╗╔╝╠═╣║ ║║  ║
  ╚═╝╝╚╝ ╚╝ ╩ ╩╚═╝╩═╝╩ `)
	dim.Println("  Secure secrets management for teams")
	fmt.Println()
}

func success(msg string, args ...interface{}) {
	green.Printf("  ✓ ")
	fmt.Printf(msg+"\n", args...)
}

func info(msg string, args ...interface{}) {
	cyan.Printf("  ℹ ")
	fmt.Printf(msg+"\n", args...)
}

func warn(msg string, args ...interface{}) {
	yellow.Printf("  ⚠ ")
	fmt.Printf(msg+"\n", args...)
}

func fail(msg string, args ...interface{}) {
	red.Printf("  ✗ ")
	fmt.Printf(msg+"\n", args...)
}

func fatal(msg string, args ...interface{}) {
	fail(msg, args...)
	os.Exit(1)
}

func prompt(label string) string {
	cyan.Printf("  ? ")
	bold.Printf("%s: ", label)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptPassword(label string) string {
	cyan.Printf("  ? ")
	bold.Printf("%s: ", label)
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return ""
	}
	return string(password)
}

func promptConfirm(label string) bool {
	answer := prompt(label + " (y/N)")
	return strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes"
}

// Spinner
type spinner struct {
	msg    string
	done   chan bool
}

func startSpinner(msg string) *spinner {
	s := &spinner{msg: msg, done: make(chan bool)}
	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				fmt.Printf("\r  %s %s", cyan.Sprint(frames[i%len(frames)]), s.msg)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	return s
}

func (s *spinner) stop(successMsg string) {
	s.done <- true
	fmt.Printf("\r  %s %s\n", green.Sprint("✓"), successMsg)
}

func (s *spinner) fail(errMsg string) {
	s.done <- true
	fmt.Printf("\r  %s %s\n", red.Sprint("✗"), errMsg)
}

// Table rendering
func printTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		dim.Println("  No data to display.")
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	fmt.Print("  ")
	for i, h := range headers {
		dim.Printf("%-*s", widths[i]+3, strings.ToUpper(h))
	}
	fmt.Println()

	// Print separator
	fmt.Print("  ")
	for i := range headers {
		dim.Print(strings.Repeat("─", widths[i]+2) + " ")
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		fmt.Print("  ")
		for i, cell := range row {
			if i < len(widths) {
				fmt.Printf("%-*s", widths[i]+3, cell)
			}
		}
		fmt.Println()
	}
}

func printKeyValue(key, value string) {
	dim.Printf("  %s: ", key)
	fmt.Println(value)
}
