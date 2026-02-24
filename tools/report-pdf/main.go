package main

import (
	"fmt"
	"os"

	"github.com/eenemeene/kitamanager-go/tools/report-pdf/internal/auth"
	"github.com/eenemeene/kitamanager-go/tools/report-pdf/internal/config"
	"github.com/eenemeene/kitamanager-go/tools/report-pdf/internal/pdf"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Logging in as %s...\n", cfg.Email)
	cookies, err := auth.Login(cfg.APIURL, cfg.Email, cfg.Password, cfg.BaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Login failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Login successful (%d cookies)\n", len(cookies))

	fmt.Println("Initializing PDF generator...")
	gen, err := pdf.NewGenerator(cookies, cfg.BaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize PDF generator: %v\n", err)
		os.Exit(1)
	}
	defer gen.Close()

	var failed []string
	for _, report := range cfg.Reports {
		fmt.Printf("Generating %s report...\n", report)
		if err := gen.GenerateReport(report, cfg.OrgID, cfg.Year, cfg.OutputDir); err != nil {
			fmt.Fprintf(os.Stderr, "  Failed: %v\n", err)
			failed = append(failed, report)
			continue
		}
	}

	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "\nFailed reports: %v\n", failed)
		os.Exit(1)
	}

	fmt.Println("\nAll reports generated successfully!")
}
