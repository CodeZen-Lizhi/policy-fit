package main

import (
	"fmt"
	"os"

	"github.com/zhenglizhi/policy-fit/internal/config"
)

func main() {
	if len(os.Args) > 1 {
		path := os.Args[1]
		if err := config.ValidateEnvFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "env validation failed (%s): %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("env validation passed: %s\n", path)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "env validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("env validation passed: env=%s file=%s\n", cfg.AppEnv, cfg.ConfigFile)
}
