package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

func PrintBanner(version, commit, date, builtBy string) {
	logo := `
  ██████╗ ███████╗███████╗████████╗███████╗██╗   ██╗██████╗ 
  ██╔══██╗██╔════╝██╔════╝╚══██╔══╝██╔════╝██║   ██║██╔══██╗
  ██████╔╝█████╗  ███████╗   ██║   ███████╗██║   ██║██████╔╝
  ██╔══██╗██╔══╝  ╚════██║   ██║   ╚════██║██║   ██║██╔══██╗
  ██████╔╝███████╗███████║   ██║   ███████║╚██████╔╝██████╔╝
  ╚═════╝ ╚══════╝╚══════╝   ╚═╝   ╚══════╝ ╚═════╝ ╚═════╝ 
`

	fmt.Print(ColorCyan + ColorBold)
	fmt.Println(logo)
	fmt.Print(ColorReset)

	fmt.Print(ColorBlue + ColorBold)
	fmt.Println("          🚀 BestSub - Best Subscription Manager")
	fmt.Print(ColorReset)

	fmt.Print(ColorDim)
	fmt.Println("  " + strings.Repeat("─", 60))
	fmt.Print(ColorReset)

	printInfo("Version", version, ColorGreen)
	printInfo("Commit", commit[:min(8, len(commit))], ColorYellow)
	printInfo("Build Time", formatDate(date), ColorBlue)
	printInfo("Built By", builtBy, ColorPurple)

	fmt.Print(ColorDim)
	fmt.Println("  " + strings.Repeat("═", 60))
	fmt.Print(ColorReset)
}

func printInfo(label, value, color string) {
	fmt.Printf("  %s%-12s%s %s%s%s\n",
		ColorDim, label+":", ColorReset,
		color, value, ColorReset)
}

func formatDate(date string) string {
	if date == "unknown" || date == "" {
		return "unknown"
	}

	layouts := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, date); err == nil {
			return t.Format("2006-01-02 15:04")
		}
	}

	return date
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
