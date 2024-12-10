package generator

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	InputPath  string
	OutputPath string
	Package    string
	Interfaces []string
	Features   []string
}

func main() {
	var config Config

	// Parse command line flags
	flag.StringVar(&config.InputPath, "input", "", "Input directory path containing interfaces")
	flag.StringVar(&config.OutputPath, "output", "", "Output file path for generated code")
	flag.StringVar(&config.Package, "package", "generated", "Package name for generated code")
	interfaceList := flag.String("interfaces", "", "Comma-separated list of interface names to process")
	featureList := flag.String("features", "", "Comma-separated list of feature names to process")

	flag.Parse()

	// Validate flags
	if config.InputPath == "" {
		log.Fatal("input path is required")
	}
	if config.OutputPath == "" {
		log.Fatal("output path is required")
	}
	if *interfaceList == "" {
		log.Fatal("at least one interface name is required")
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(config.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Check if input directory exists
	if _, err := os.Stat(config.InputPath); os.IsNotExist(err) {
		log.Fatalf("Input directory does not exist: %s", config.InputPath)
	}

	// Parse interface list
	config.Interfaces = filepath.SplitList(*interfaceList)
	config.Features = strings.Split(*featureList, ",")

	// Create generator config
	genConfig := GeneratorConfig{
		Package:    config.Package,
		InputPath:  config.InputPath,
		OutputPath: config.OutputPath,
		Interfaces: config.Interfaces,
		Features:   config.Features,
	}

	// Initialize generator
	gen, err := NewGenerator(genConfig)
	if err != nil {
		log.Fatalf("Failed to initialize generator: %v", err)
	}

	// Generate code
	if err := gen.Generate(); err != nil {
		log.Fatalf("Code generation failed: %v", err)
	}

	fmt.Printf("Successfully generated code for interfaces: %v\n", config.Interfaces)
	fmt.Printf("Output written to: %s\n", config.OutputPath)
}
