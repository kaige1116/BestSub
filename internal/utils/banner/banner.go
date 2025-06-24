package banner

import (
	"fmt"
	"strings"
	"time"

	"github.com/bestruirui/bestsub/internal/utils/color"
)

func Print(version, commit, date, builtBy string) {
	logo := `
  ██████╗ ███████╗███████╗████████╗███████╗██╗   ██╗██████╗ 
  ██╔══██╗██╔════╝██╔════╝╚══██╔══╝██╔════╝██║   ██║██╔══██╗
  ██████╔╝█████╗  ███████╗   ██║   ███████╗██║   ██║██████╔╝
  ██╔══██╗██╔══╝  ╚════██║   ██║   ╚════██║██║   ██║██╔══██╗
  ██████╔╝███████╗███████║   ██║   ███████║╚██████╔╝██████╔╝
  ╚═════╝ ╚══════╝╚══════╝   ╚═╝   ╚══════╝ ╚═════╝ ╚═════╝ 
`

	fmt.Print(color.Cyan + color.Bold)
	fmt.Println(logo)
	fmt.Print(color.Reset)

	fmt.Print(color.Blue + color.Bold)
	fmt.Println("          🚀 BestSub - Best Subscription Manager")
	fmt.Print(color.Reset)

	fmt.Print(color.Dim)
	fmt.Println("  " + strings.Repeat("─", 60))
	fmt.Print(color.Reset)

	printInfo("Version", version, color.Green)
	printInfo("Commit", commit[:min(8, len(commit))], color.Yellow)
	printInfo("Build Time", formatDate(date), color.Blue)
	printInfo("Built By", builtBy, color.Purple)

	fmt.Print(color.Dim)
	fmt.Println("  " + strings.Repeat("═", 60))
	fmt.Print(color.Reset)
}

func printInfo(label, value, print_color string) {
	fmt.Printf("  %s%-12s%s %s%s%s\n",
		color.Dim, label+":", color.Reset,
		print_color, value, color.Reset)
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
