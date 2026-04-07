package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"

	// flags
	outputFormat string
	ignoreKeys   []string
	colorOutput  bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// rootCmd is the base command for seatmap-diff.
var rootCmd = &cobra.Command{
	Use:   "seatmap-diff",
	Short: "Diff and audit YAML/JSON infrastructure configs across environments",
	Long: `seatmap-diff compares YAML and JSON configuration files across environments,
highlighting structural and value-level differences to help teams audit
infrastructure changes safely.`,
	Version: version,
}

// diffCmd compares two config files and outputs their differences.
var diffCmd = &cobra.Command{
	Use:   "diff <file-a> <file-b>",
	Short: "Compare two YAML or JSON config files",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

// auditCmd checks a config file against a baseline for unexpected changes.
var auditCmd = &cobra.Command{
	Use:   "audit <baseline> <target>",
	Short: "Audit a config file against a known baseline",
	Args:  cobra.ExactArgs(2),
	RunE:  runAudit,
}

func init() {
	// diff flags
	diffCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json, yaml")
	diffCmd.Flags().StringArrayVarP(&ignoreKeys, "ignore", "i", nil, "Keys to ignore during comparison (can be repeated)")
	diffCmd.Flags().BoolVar(&colorOutput, "color", true, "Enable colored output")

	// audit flags
	auditCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json, yaml")
	auditCmd.Flags().StringArrayVarP(&ignoreKeys, "ignore", "i", nil, "Keys to ignore during audit (can be repeated)")
	auditCmd.Flags().BoolVar(&colorOutput, "color", true, "Enable colored output")

	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(auditCmd)
}

// runDiff is the handler for the diff subcommand.
func runDiff(cmd *cobra.Command, args []string) error {
	fileA := args[0]
	fileB := args[1]

	fmt.Printf("Diffing %s <-> %s\n", fileA, fileB)
	fmt.Printf("Output format: %s | Ignore keys: %v\n", outputFormat, ignoreKeys)

	// TODO: load, parse, and diff the two files
	return nil
}

// runAudit is the handler for the audit subcommand.
func runAudit(cmd *cobra.Command, args []string) error {
	baseline := args[0]
	target := args[1]

	fmt.Printf("Auditing %s against baseline %s\n", target, baseline)
	fmt.Printf("Output format: %s | Ignore keys: %v\n", outputFormat, ignoreKeys)

	// TODO: load, parse, and audit the target against the baseline
	return nil
}
