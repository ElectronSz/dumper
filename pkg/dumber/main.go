package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/OffbyteSecure/dump_util/pkg/dump_util"
	"github.com/OffbyteSecure/dump_util/pkg/dump_util/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	dbType        string
	connStr       string
	outputFile    string
	compress      bool
	batchSize     int
	maxWorkers    int
	excludeTables string // Comma-separated
	version       = "v0.1.1"
	author        = "OffbyteSecure"
	website       = "https://offbytesecure.com"
)

// Colors for output
var (
	infoColor    = color.New(color.FgCyan).SprintFunc()
	successColor = color.New(color.FgGreen).SprintFunc()
	errorColor   = color.New(color.FgRed).SprintFunc()
	warnColor    = color.New(color.FgYellow).SprintFunc()
	bannerColor  = color.New(color.FgMagenta, color.Bold).SprintFunc()
)

// ASCII art for "Dumper" using a fancier "big" font style
var banner = []string{
	"         ____                        ",
	"        |  __|                       ",
	"        | |__ _   _ _ __   ___ _ __  ",
	"        |  __| | | | '_ \\ / __| '_ \\ ",
	"        | |  | |_| | |_) | (__| |_) |",
	"        |_|   \\__,_|_.__/ \\___|_.__/ ",
	"                                     ",
}

// printAnimatedBanner displays the ASCII art with a line-by-line animation
func printAnimatedBanner() {
	fmt.Println()
	for _, line := range banner {
		fmt.Printf("%s%s\n", bannerColor(line), strings.Repeat(" ", 10))
		time.Sleep(100 * time.Millisecond) // Delay for animation effect
	}
	fmt.Printf("%sAuthor: %s\n", infoColor("Info: "), author)
	fmt.Printf("%sWebsite: %s\n", infoColor("Info: "), website)
	fmt.Printf("%sVersion: %s\n", infoColor("Info: "), version)
	fmt.Println()
}

var rootCmd = &cobra.Command{
	Use:   "dumper",
	Short: "A CLI tool for database backups using dump_util",
	Long: `dumper is a command-line tool to create database backups for PostgreSQL, MySQL, or MongoDB.
It supports compression, batch processing, and table exclusion for efficient backups.
Author: ` + author + `
Website: ` + website,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate inputs
		if err := validateInputs(); err != nil {
			log.Fatalf("%s %s", errorColor("Error:"), err)
		}

		fmt.Printf("%s Starting database dump for %s...\n", infoColor("Info:"), dbType)
		opts := &types.BackupOptions{
			Compress:   compress,
			BatchSize:  batchSize,
			MaxWorkers: maxWorkers,
			Exclude:    filterExcludeTables(excludeTables),
		}

		if err := dump_util.DumpDatabase(dbType, connStr, outputFile, opts); err != nil {
			log.Fatalf("%s Failed to dump database: %v", errorColor("Error:"), err)
		}
		fmt.Printf("%s Dump completed successfully: %s\n", successColor("Success:"), outputFile)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of dumper",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s dumper version %s\n", infoColor("Version:"), version)
		fmt.Printf("%s Author: %s\n", infoColor("Author:"), author)
		fmt.Printf("%s Website: %s\n", infoColor("Website:"), website)
	},
}

func init() {
	// Define flags with improved descriptions
	rootCmd.Flags().StringVarP(&dbType, "type", "t", "", "Database type (required: postgres, mysql, mongodb)")
	rootCmd.Flags().StringVarP(&connStr, "conn", "c", "", "Database connection string (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "backup.sql", "Output file path for the dump (default: backup.sql)")
	rootCmd.Flags().BoolVarP(&compress, "compress", "z", false, "Enable gzip compression for the output file")
	rootCmd.Flags().IntVarP(&batchSize, "batch-size", "b", 5000, "Number of rows to process per batch (min: 1)")
	rootCmd.Flags().IntVarP(&maxWorkers, "workers", "w", 5, "Maximum number of concurrent workers (min: 1, max: 50)")
	rootCmd.Flags().StringVar(&excludeTables, "exclude", "", "Comma-separated list of tables/collections to exclude")

	// Mark required flags
	rootCmd.MarkFlagRequired("type")
	rootCmd.MarkFlagRequired("conn")

	// Add version command
	rootCmd.AddCommand(versionCmd)
}

// validateInputs checks flag values for correctness
func validateInputs() error {
	// Validate dbType
	validTypes := map[string]bool{"postgres": true, "mysql": true, "mongodb": true}
	if !validTypes[strings.ToLower(dbType)] {
		return fmt.Errorf("invalid database type: %s (must be postgres, mysql, or mongodb)", dbType)
	}

	// Validate batchSize
	if batchSize < 1 {
		return fmt.Errorf("batch-size must be at least 1, got %d", batchSize)
	}

	// Validate maxWorkers
	if maxWorkers < 1 || maxWorkers > 50 {
		return fmt.Errorf("workers must be between 1 and 50, got %d", maxWorkers)
	}

	// Validate output file
	if outputFile == "" {
		return fmt.Errorf("output file path cannot be empty")
	}

	return nil
}

// filterExcludeTables processes the excludeTables string, removing empty entries
func filterExcludeTables(exclude string) []string {
	if exclude == "" {
		return nil
	}
	var result []string
	for _, table := range strings.Split(exclude, ",") {
		if trimmed := strings.TrimSpace(table); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	// Display animated banner on startup
	printAnimatedBanner()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s Failed to execute command: %v\n", errorColor("Error:"), err)
		os.Exit(1)
	}
}
