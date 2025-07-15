package info

import (
	"fmt"
	"strings"
	"time"

	"github.com/bestruirui/bestsub/internal/utils/color"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
	Author    = "bestrui"
	Repo      = "https://github.com/bestruirui/bestsub"
)

func Banner() {
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

	printInfo("Version", Version, color.Green)
	printInfo("Commit", Commit[:min(8, len(Commit))], color.Yellow)
	printInfo("Build Time", formatDate(BuildTime), color.Blue)
	printInfo("Built By", Author, color.Purple)
	printInfo("Repo", Repo, color.Cyan)

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
