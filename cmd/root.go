package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"visuche/internal/csv"
	"visuche/internal/git"
	"visuche/internal/github"
	"visuche/internal/i18n"
	"visuche/internal/stats"

	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var repo string
var since string
var until string
var author string
var label string
var csvOutput bool
var lang string
var langJP bool

var rootCmd = &cobra.Command{
	Use:   "visuche",
	Short: "A visualization tool for GitHub repository metrics and CI/CD analytics.",
	Long:  `visuche (visualization check) analyzes GitHub repositories to provide insights on PR metrics, lead times, and CI/CD performance.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no arguments provided, use interactive mode
		if repo == "" && since == "" && until == "" {
			runInteractiveMode()
			return
		}

		// Traditional argument-based mode
		runAnalysis()
	},
}

func getTargetRepo() (string, error) {
	if repo != "" {
		return repo, nil
	}

	detectedRepo, err := git.GetRepoFromGitRemote()
	if err == nil {
		prompt := promptui.Select{
			Label: fmt.Sprintf("? Found repository '%s'. Use this one?", detectedRepo),
			Items: []string{"Yes", "No"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return "", fmt.Errorf("prompt failed %w", err)
		}

		if result == "Yes" {
			return detectedRepo, nil
		}
	}

	// Ask for manual input
	prompt := promptui.Prompt{
		Label: "Please enter the repository in 'owner/repo' format",
		Validate: func(input string) error {
			if len(strings.Split(input, "/")) != 2 || strings.TrimSpace(input) == "" {
				return fmt.Errorf("invalid format, please use 'owner/repo'")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed %w", err)
	}
	return result, nil
}

func init() {
	cobra.OnInitialize(applyLanguageSetting)

	rootCmd.PersistentFlags().StringVar(&repo, "repo", "", "Specify the GitHub repository in 'owner/repo' format")
	rootCmd.PersistentFlags().StringVar(&since, "since", "", "Fetch PRs created after this date (YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVar(&until, "until", "", "Fetch PRs created before this date (YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVar(&author, "author", "", "Filter PRs by author username")
	rootCmd.PersistentFlags().StringVar(&label, "label", "", "Filter PRs by label name")
	rootCmd.PersistentFlags().BoolVar(&csvOutput, "csv", false, "Export results to CSV file")
	rootCmd.PersistentFlags().StringVar(&lang, "lang", "en", "Output language (en/jp)")
	rootCmd.PersistentFlags().BoolVar(&langJP, "jp", false, "Use Japanese output (shortcut for --lang=jp)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func applyLanguageSetting() {
	selected := strings.ToLower(lang)
	if langJP {
		selected = "jp"
	}
	i18n.SetLanguage(selected)
}

// CalculateLeadTimes calculates the lead time for each pull request.
// It returns a new slice containing only closed or merged PRs with their lead time calculated.
func CalculateLeadTimes(prs []github.PullRequest) []github.PullRequest {
	processedPRs := make([]github.PullRequest, 0, len(prs))
	for _, pr := range prs {
		var endAt time.Time
		if pr.Merged && !pr.MergedAt.IsZero() {
			endAt = pr.MergedAt
		} else if !pr.ClosedAt.IsZero() {
			endAt = pr.ClosedAt
		}

		if !endAt.IsZero() {
			pr.LeadTime = endAt.Sub(pr.CreatedAt)
		}

		// Keep open PRs as well so metrics like TotalPRs/WIP are accurate.
		processedPRs = append(processedPRs, pr)
	}
	return processedPRs
}

// displayStatsTable displays PR statistics in a formatted table
func displayStatsTable(statistics stats.Stats) {
	fmt.Println("\n" + i18n.T("📊 Pull Request Statistics"))
	fmt.Println("=" + strings.Repeat("=", 50))

	// Basic Statistics Table
	fmt.Println("\n" + i18n.T("🔢 Basic Metrics:"))
	basicTable := tablewriter.NewWriter(os.Stdout)
	basicTable.SetHeader([]string{i18n.T("Metric"), i18n.T("Value")})
	basicTable.SetBorder(true)
	basicTable.Append([]string{i18n.T("Total PRs"), fmt.Sprintf("%d", statistics.TotalPRs)})
	basicTable.Append([]string{i18n.T("Merged PRs"), fmt.Sprintf("%d", statistics.MergedPRs)})
	basicTable.Append([]string{i18n.T("WIP PRs"), fmt.Sprintf("%d", statistics.WIPPRCount)})
	basicTable.Append([]string{i18n.T("Releases (main/master merges)"), fmt.Sprintf("%d", statistics.ReleaseCount)})
	if statistics.TotalPRs > 0 {
		basicTable.Append([]string{i18n.T("Merge Rate"), fmt.Sprintf("%.1f%%", float64(statistics.MergedPRs)/float64(statistics.TotalPRs)*100)})
	}
	basicTable.Render()

	// Timing Statistics Table
	fmt.Println("\n" + i18n.T("⏱️ Timing Metrics:"))
	timingTable := tablewriter.NewWriter(os.Stdout)
	timingTable.SetHeader([]string{i18n.T("Metric"), i18n.T("Average"), i18n.T("Median")})
	timingTable.SetBorder(true)
	timingTable.Append([]string{
		i18n.T("Lead Time"),
		formatDuration(statistics.AverageLeadTime),
		formatDuration(statistics.MedianLeadTime),
	})
	timingTable.Append([]string{
		i18n.T("Review Time"),
		formatDuration(statistics.AverageReviewTime),
		"-",
	})
	timingTable.Append([]string{
		i18n.T("Merge Wait Time"),
		formatDuration(statistics.AverageMergeWaitTime),
		formatDuration(statistics.MedianMergeWaitTime),
	})
	timingTable.Append([]string{
		i18n.T("Approval→Merge Time"),
		formatDuration(statistics.AverageApprovalToMerge),
		formatDuration(statistics.MedianApprovalToMerge),
	})
	timingTable.Append([]string{
		i18n.T("Commit→PR Time"),
		formatDuration(statistics.AverageCommitToPRTime),
		"-",
	})
	timingTable.Render()

	// Code Change Statistics Table
	fmt.Println("\n" + i18n.T("💻 Code Change Metrics:"))
	codeTable := tablewriter.NewWriter(os.Stdout)
	codeTable.SetHeader([]string{i18n.T("Metric"), i18n.T("Average")})
	codeTable.SetBorder(true)
	codeTable.Append([]string{i18n.T("Files Changed"), fmt.Sprintf("%.1f", statistics.AverageFilesChanged)})
	codeTable.Append([]string{i18n.T("Lines Added"), fmt.Sprintf("%.1f", statistics.AverageAdditions)})
	codeTable.Append([]string{i18n.T("Lines Deleted"), fmt.Sprintf("%.1f", statistics.AverageDeletions)})
	codeTable.Append([]string{i18n.T("Commits per PR"), fmt.Sprintf("%.1f", statistics.AverageCommitsPerPR)})
	codeTable.Append([]string{i18n.T("Commit Frequency/Week"), fmt.Sprintf("%.1f", statistics.CommitFrequencyPerWeek)})
	codeTable.Render()

	// Collaboration Statistics Table
	fmt.Println("\n" + i18n.T("👥 Collaboration Metrics:"))
	collabTable := tablewriter.NewWriter(os.Stdout)
	collabTable.SetHeader([]string{i18n.T("Metric"), i18n.T("Value")})
	collabTable.SetBorder(true)
	collabTable.Append([]string{i18n.T("Avg Reviewers per PR"), fmt.Sprintf("%.1f", statistics.AverageReviewersPerPR)})
	collabTable.Append([]string{i18n.T("Self-Merge Rate"), fmt.Sprintf("%.1f%%", statistics.SelfMergeRate)})
	collabTable.Render()

	// Review Comment Analysis (focus on code review comments only)
	if statistics.PRsWithReviewComments > 0 {
		fmt.Println("\n" + i18n.T("💬 Code Review Analysis:"))
		reviewTable := tablewriter.NewWriter(os.Stdout)
		reviewTable.SetHeader([]string{i18n.T("Metric"), i18n.T("Average"), i18n.T("Median"), i18n.T("Max")})
		reviewTable.SetBorder(true)

		reviewTable.Append([]string{
			i18n.T("Review Comments per PR"),
			fmt.Sprintf("%.1f", statistics.AverageReviewCommentsPerPR),
			fmt.Sprintf("%.1f", statistics.MedianReviewCommentsPerPR),
			fmt.Sprintf("%d", statistics.MaxReviewCommentsInPR),
		})
		reviewTable.Render()

		// Review Coverage Statistics
		fmt.Println("\n" + i18n.T("📈 Review Coverage:"))
		coverageTable := tablewriter.NewWriter(os.Stdout)
		coverageTable.SetHeader([]string{i18n.T("Metric"), i18n.T("Count"), i18n.T("Percentage")})
		coverageTable.SetBorder(true)

		if statistics.TotalPRs > 0 {
			reviewCommentCoverage := float64(statistics.PRsWithReviewComments) / float64(statistics.TotalPRs) * 100.0

			coverageTable.Append([]string{i18n.T("PRs with Review Comments"), fmt.Sprintf("%d", statistics.PRsWithReviewComments), fmt.Sprintf("%.1f%%", reviewCommentCoverage)})
			coverageTable.Append([]string{i18n.T("PRs without Review Comments"), fmt.Sprintf("%d", statistics.PRsWithoutReviewComments), fmt.Sprintf("%.1f%%", 100.0-reviewCommentCoverage)})
		}

		coverageTable.Render()

		// Review Density Analysis
		fmt.Println("\n" + i18n.T("🔍 Review Quality:"))
		densityTable := tablewriter.NewWriter(os.Stdout)
		densityTable.SetHeader([]string{i18n.T("Metric"), i18n.T("Value")})
		densityTable.SetBorder(true)

		// Calculate review density based on review comments only
		reviewDensity := 0.0
		totalReviewComments := int(statistics.AverageReviewCommentsPerPR * float64(statistics.TotalPRs))
		if statistics.AverageAdditions+statistics.AverageDeletions > 0 {
			totalLines := int((statistics.AverageAdditions + statistics.AverageDeletions) * float64(statistics.TotalPRs))
			reviewDensity = float64(totalReviewComments) / float64(totalLines) * 100.0
		}

		densityTable.Append([]string{i18n.T("Review Comment Density"), i18n.Sprintf("%.2f comments/100 lines", reviewDensity)})
		densityTable.Render()
	} else {
		// Show a message when no review comments are found
		fmt.Println("\n" + i18n.T("💬 Code Review Analysis:"))
		fmt.Printf(i18n.Sprintf("📝 No code review comments found in this period (%d PRs analyzed)", statistics.TotalPRs) + "\n")
		fmt.Printf(i18n.T("💡 This could indicate:") + "\n")
		fmt.Printf(i18n.T("   • Code quality is consistently high") + "\n")
		fmt.Printf(i18n.T("   • Team does reviews via other channels") + "\n")
		fmt.Printf(i18n.T("   • PRs are small and self-explanatory") + "\n")
	}

	// Merge Type Statistics Table
	if len(statistics.MergeTypeTrend) > 0 {
		fmt.Println("\n" + i18n.T("🔀 Merge Type Distribution:"))
		mergeTable := tablewriter.NewWriter(os.Stdout)
		mergeTable.SetHeader([]string{i18n.T("Merge Type"), i18n.T("Percentage")})
		mergeTable.SetBorder(true)
		for mergeType, percentage := range statistics.MergeTypeTrend {
			mergeTable.Append([]string{mergeType, fmt.Sprintf("%.1f%%", percentage)})
		}
		mergeTable.Render()
	}

	fmt.Println()
}

// formatDuration formats a time.Duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	} else {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
}

// runInteractiveMode runs the interactive mode for repository and date selection
func runInteractiveMode() {
	fmt.Println("🎯 Welcome to visuche - Interactive GitHub Analytics")
	fmt.Println("=" + strings.Repeat("=", 50))

	// Step 1: Repository selection
	targetRepo, err := getInteractiveRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	repo = targetRepo

	// Step 2: Analysis type selection
	analysisType, err := selectAnalysisType()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Date range selection
	startDate, endDate, err := selectDateRange()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	since = startDate
	until = endDate

	// Step 4: Optional filters
	if analysisType == "PR Analysis" {
		author, label = selectOptionalFilters()
	}

	// Step 5: Run analysis
	fmt.Printf("\n✅ Configuration:\n")
	fmt.Printf("  Repository: %s\n", repo)
	fmt.Printf("  Analysis: %s\n", analysisType)
	fmt.Printf("  Period: %s to %s\n", since, until)
	if author != "" {
		fmt.Printf("  Author: %s\n", author)
	}
	if label != "" {
		fmt.Printf("  Label: %s\n", label)
	}

	confirm := promptui.Select{
		Label: "Proceed with analysis?",
		Items: []string{"Yes", "No"},
	}
	_, result, err := confirm.Run()
	if err != nil || result != "Yes" {
		fmt.Println("❌ Analysis cancelled")
		return
	}

	// Run the appropriate analysis based on type
	if analysisType == "Actions Analysis" {
		runActionsAnalysis()
	} else {
		runAnalysis()
	}
}

// runAnalysis performs the actual analysis with current settings
func runAnalysis() {
	// Determine the target repository
	targetRepo, err := getTargetRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	repo = targetRepo

	fmt.Printf(i18n.Sprintf("✅ Using repository: %s\n", repo))

	// Fetch pull requests
	fmt.Println(i18n.T("📥 Fetching pull requests..."))
	prs, err := github.FetchPullRequests(repo, since, until, author, label, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching pull requests: %v\n", err)
		os.Exit(1)
	}

	// Calculate lead times
	processedPRs := CalculateLeadTimes(prs)

	// Fetch comment timing data
	processedPRs = github.FetchPRCommentTiming(repo, processedPRs)

	// Calculate stats
	statistics := stats.CalculateStats(processedPRs)

	// Display stats
	displayStatsTable(statistics)

	// Output to CSV if requested
	if csvOutput {
		repoNameForFile := strings.ReplaceAll(repo, "/", "-")
		csvFilename := fmt.Sprintf("visuche_%s.csv", repoNameForFile)
		if err := csv.WritePullRequestsToCSV(csvFilename, processedPRs); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("📁 CSV output: %s\n", csvFilename)
	}
}

// getInteractiveRepo gets repository interactively
func getInteractiveRepo() (string, error) {
	detectedRepo, err := git.GetRepoFromGitRemote()
	if err == nil {
		prompt := promptui.Select{
			Label: fmt.Sprintf("Found repository '%s' in current directory. Use this?", detectedRepo),
			Items: []string{"Yes, use detected repo", "No, enter manually"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return "", fmt.Errorf("prompt failed %w", err)
		}

		if result == "Yes, use detected repo" {
			return detectedRepo, nil
		}
	}

	// Manual entry
	prompt := promptui.Prompt{
		Label: "Enter GitHub repository (owner/repo format)",
		Validate: func(input string) error {
			if len(strings.Split(input, "/")) != 2 || strings.TrimSpace(input) == "" {
				return fmt.Errorf("invalid format, please use 'owner/repo'")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed %w", err)
	}
	return result, nil
}

// selectAnalysisType allows user to select analysis type
func selectAnalysisType() (string, error) {
	prompt := promptui.Select{
		Label: "Select analysis type",
		Items: []string{
			"PR Analysis - Pull Request metrics and lead times",
			"Actions Analysis - CI/CD performance and workflow insights",
		},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(result, "Actions Analysis") {
		return "Actions Analysis", nil
	}
	return "PR Analysis", nil
}

// selectDateRange allows user to select date range with simplified options
func selectDateRange() (string, string, error) {
	now := time.Now()

	prompt := promptui.Select{
		Label: "Select time period",
		Items: []string{
			"Last 1 week",
			"Last 2 weeks",
			"Last 1 month",
			"Custom range (flexible input)",
		},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", "", err
	}

	switch result {
	case "Last 1 week":
		return now.AddDate(0, 0, -7).Format("2006-01-02"), now.Format("2006-01-02"), nil
	case "Last 2 weeks":
		return now.AddDate(0, 0, -14).Format("2006-01-02"), now.Format("2006-01-02"), nil
	case "Last 1 month":
		return now.AddDate(0, -1, 0).Format("2006-01-02"), now.Format("2006-01-02"), nil
	case "Custom range (flexible input)":
		return getEnhancedCustomDateRange()
	default:
		return "", "", fmt.Errorf("unknown selection")
	}
}

// getCustomDateRange gets custom date range from user (legacy function)
func getCustomDateRange() (string, string, error) {
	return getEnhancedCustomDateRange()
}

// getEnhancedCustomDateRange provides flexible custom date input with smart parsing
func getEnhancedCustomDateRange() (string, string, error) {
	now := time.Now()

	fmt.Println("\n🗓️  Custom Date Range Input")
	fmt.Println("=========================")
	fmt.Println("Supported formats:")
	fmt.Println("  • YYYY-MM-DD (e.g., 2024-01-15)")
	fmt.Println("  • Relative: '30 days ago', '2 weeks ago', '3 months ago'")
	fmt.Println("  • Keywords: 'today', 'yesterday', 'last monday'")
	fmt.Println("  • Shortcuts: '2024-01' (whole month), '2024-Q1' (quarter)")
	fmt.Println()

	// Start date input with enhanced parsing
	startPrompt := promptui.Prompt{
		Label: "Enter start date",
		Validate: func(input string) error {
			_, err := ParseFlexibleDate(input, now)
			if err != nil {
				return fmt.Errorf("invalid date format: %v", err)
			}
			return nil
		},
	}
	startInput, err := startPrompt.Run()
	if err != nil {
		return "", "", err
	}

	startDate, _ := ParseFlexibleDate(startInput, now)

	// End date input with smart defaults
	endPrompt := promptui.Prompt{
		Label: fmt.Sprintf("Enter end date (default: today - %s)", now.Format("2006-01-02")),
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return nil // Allow empty for default
			}
			_, err := ParseFlexibleDate(input, now)
			if err != nil {
				return fmt.Errorf("invalid date format: %v", err)
			}
			return nil
		},
	}
	endInput, err := endPrompt.Run()
	if err != nil {
		return "", "", err
	}

	var endDate time.Time
	if strings.TrimSpace(endInput) == "" {
		endDate = now // Default to today
	} else {
		endDate, _ = ParseFlexibleDate(endInput, now)
	}

	// Validate date range
	if endDate.Before(startDate) {
		return "", "", fmt.Errorf("end date cannot be before start date")
	}

	fmt.Printf("✅ Selected period: %s to %s\n",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	return startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), nil
}

// selectOptionalFilters allows user to set optional filters
func selectOptionalFilters() (string, string) {
	var selectedAuthor, selectedLabel string

	// Author filter
	authorPrompt := promptui.Select{
		Label: "Filter by author?",
		Items: []string{"No filter", "Specify author"},
	}
	_, authorResult, err := authorPrompt.Run()
	if err == nil && authorResult == "Specify author" {
		prompt := promptui.Prompt{
			Label: "Enter GitHub username",
		}
		selectedAuthor, _ = prompt.Run()
	}

	// Label filter
	labelPrompt := promptui.Select{
		Label: "Filter by label?",
		Items: []string{"No filter", "Specify label"},
	}
	_, labelResult, err := labelPrompt.Run()
	if err == nil && labelResult == "Specify label" {
		prompt := promptui.Prompt{
			Label: "Enter label name",
		}
		selectedLabel, _ = prompt.Run()
	}

	return selectedAuthor, selectedLabel
}

// ParseFlexibleDate parses various date input formats
func ParseFlexibleDate(input string, baseDate time.Time) (time.Time, error) {
	input = strings.TrimSpace(strings.ToLower(input))

	// Standard YYYY-MM-DD format
	if date, err := time.Parse("2006-01-02", input); err == nil {
		return date, nil
	}

	// Keywords
	switch input {
	case "today":
		return baseDate, nil
	case "yesterday":
		return baseDate.AddDate(0, 0, -1), nil
	}

	// Relative dates with regex patterns
	relativePatterns := map[string]func(int) time.Time{
		`(\d+)\s*days?\s*ago`:   func(n int) time.Time { return baseDate.AddDate(0, 0, -n) },
		`(\d+)\s*weeks?\s*ago`:  func(n int) time.Time { return baseDate.AddDate(0, 0, -n*7) },
		`(\d+)\s*months?\s*ago`: func(n int) time.Time { return baseDate.AddDate(0, -n, 0) },
		`(\d+)\s*years?\s*ago`:  func(n int) time.Time { return baseDate.AddDate(-n, 0, 0) },
	}

	for pattern, calculator := range relativePatterns {
		if matches := regexp.MustCompile(pattern).FindStringSubmatch(input); matches != nil {
			if num, err := strconv.Atoi(matches[1]); err == nil {
				return calculator(num), nil
			}
		}
	}

	// Month shortcuts (e.g., "2024-01" -> first day of January 2024)
	if matches := regexp.MustCompile(`^(\d{4})-(\d{1,2})$`).FindStringSubmatch(input); matches != nil {
		year, _ := strconv.Atoi(matches[1])
		month, _ := strconv.Atoi(matches[2])
		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, baseDate.Location()), nil
	}

	// Quarter shortcuts (e.g., "2024-Q1" -> first day of Q1 2024)
	if matches := regexp.MustCompile(`^(\d{4})-q(\d)$`).FindStringSubmatch(input); matches != nil {
		year, _ := strconv.Atoi(matches[1])
		quarter, _ := strconv.Atoi(matches[2])
		if quarter < 1 || quarter > 4 {
			return time.Time{}, fmt.Errorf("invalid quarter: %d", quarter)
		}
		month := (quarter-1)*3 + 1
		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, baseDate.Location()), nil
	}

	// Day names (basic implementation for "last monday", etc.)
	dayNames := map[string]time.Weekday{
		"sunday": time.Sunday, "monday": time.Monday, "tuesday": time.Tuesday,
		"wednesday": time.Wednesday, "thursday": time.Thursday, "friday": time.Friday, "saturday": time.Saturday,
	}

	if matches := regexp.MustCompile(`^last\s+(\w+)$`).FindStringSubmatch(input); matches != nil {
		if targetDay, exists := dayNames[matches[1]]; exists {
			// Find the most recent occurrence of that day
			days := int(baseDate.Weekday() - targetDay)
			if days <= 0 {
				days += 7
			}
			return baseDate.AddDate(0, 0, -days), nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized date format: %s", input)
}
