package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/zhenglizhi/policy-fit/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down]")
	}

	action := os.Args[1]

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	switch action {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up completed successfully")
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down completed successfully")
	default:
		log.Fatalf("Unknown action: %s", action)
	}
}

func migrateUp(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS user_account (
			id BIGSERIAL PRIMARY KEY,
			phone VARCHAR(32) UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS analysis_task (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES user_account(id),
			status VARCHAR(32) NOT NULL,
			risk_summary JSONB,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX idx_task_user_id ON analysis_task(user_id)`,
		`CREATE INDEX idx_task_status ON analysis_task(status)`,
		`CREATE TABLE IF NOT EXISTS document (
			id BIGSERIAL PRIMARY KEY,
			task_id BIGINT NOT NULL REFERENCES analysis_task(id) ON DELETE CASCADE,
			doc_type VARCHAR(32) NOT NULL,
			file_name VARCHAR(256) NOT NULL,
			storage_key VARCHAR(512) NOT NULL,
			parse_status VARCHAR(32) NOT NULL,
			parsed_text TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX idx_document_task_id ON document(task_id)`,
		`CREATE TABLE IF NOT EXISTS risk_finding (
			id BIGSERIAL PRIMARY KEY,
			task_id BIGINT NOT NULL REFERENCES analysis_task(id) ON DELETE CASCADE,
			level VARCHAR(16) NOT NULL,
			topic VARCHAR(64) NOT NULL,
			summary TEXT NOT NULL,
			health_evidence JSONB NOT NULL,
			policy_evidence JSONB NOT NULL,
			questions JSONB NOT NULL,
			actions JSONB,
			confidence NUMERIC(4,3),
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX idx_finding_task_id ON risk_finding(task_id)`,
		`CREATE INDEX idx_finding_level ON risk_finding(level)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	tables := []string{
		"risk_finding",
		"document",
		"analysis_task",
		"user_account",
	}

	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}
