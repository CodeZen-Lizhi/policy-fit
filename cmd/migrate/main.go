package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/zhenglizhi/policy-fit/internal/config"
)

const (
	defaultMigrationDir = "internal/migrations"
)

var migrationNamePattern = regexp.MustCompile(`^(\d+)_([a-zA-Z0-9_]+)\.(up|down)\.sql$`)

type migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

type migrator struct {
	db         *sql.DB
	migrations []migration
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down] [steps|all]")
	}

	action := os.Args[1]
	steps, err := parseDownSteps(os.Args)
	if err != nil {
		log.Fatalf("Invalid down steps: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
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

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	migrationsDir, err := resolveMigrationsDir()
	if err != nil {
		log.Fatalf("Failed to resolve migrations dir: %v", err)
	}

	m, err := newMigrator(db, migrationsDir)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	switch action {
	case "up":
		applied, upErr := m.Up()
		if upErr != nil {
			log.Fatalf("Migration up failed: %v", upErr)
		}
		log.Printf("Migration up completed successfully, applied=%d", applied)
	case "down":
		rolledBack, downErr := m.Down(steps)
		if downErr != nil {
			log.Fatalf("Migration down failed: %v", downErr)
		}
		log.Printf("Migration down completed successfully, rolled_back=%d", rolledBack)
	default:
		log.Fatalf("Unknown action: %s", action)
	}
}

func parseDownSteps(args []string) (int, error) {
	if len(args) < 3 {
		return 1, nil
	}

	raw := args[2]
	if raw == "all" {
		return -1, nil
	}

	steps, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	if steps <= 0 {
		return 0, errors.New("steps must be positive or 'all'")
	}
	return steps, nil
}

func resolveMigrationsDir() (string, error) {
	if value := os.Getenv("MIGRATIONS_DIR"); value != "" {
		path, err := filepath.Abs(value)
		if err != nil {
			return "", err
		}
		return path, nil
	}

	path, err := filepath.Abs(defaultMigrationDir)
	if err != nil {
		return "", err
	}
	return path, nil
}

func newMigrator(db *sql.DB, migrationsDir string) (*migrator, error) {
	migrations, err := loadMigrations(migrationsDir)
	if err != nil {
		return nil, err
	}
	return &migrator{
		db:         db,
		migrations: migrations,
	}, nil
}

func (m *migrator) Up() (int, error) {
	if err := ensureMigrationTable(m.db); err != nil {
		return 0, err
	}

	appliedSet, err := appliedMigrationSet(m.db)
	if err != nil {
		return 0, err
	}

	appliedCount := 0
	for _, migration := range m.migrations {
		if appliedSet[migration.Version] {
			continue
		}
		if err := applyMigration(m.db, migration); err != nil {
			return appliedCount, err
		}
		appliedCount++
	}

	return appliedCount, nil
}

func (m *migrator) Down(steps int) (int, error) {
	if err := ensureMigrationTable(m.db); err != nil {
		return 0, err
	}

	appliedVersions, err := appliedMigrationVersionsDesc(m.db)
	if err != nil {
		return 0, err
	}
	if len(appliedVersions) == 0 {
		return 0, nil
	}

	targetVersions := appliedVersions
	if steps > 0 && steps < len(appliedVersions) {
		targetVersions = appliedVersions[:steps]
	}

	byVersion := make(map[int]migration, len(m.migrations))
	for _, migration := range m.migrations {
		byVersion[migration.Version] = migration
	}

	rolledBack := 0
	for _, version := range targetVersions {
		migration, ok := byVersion[version]
		if !ok {
			return rolledBack, fmt.Errorf("missing down migration file for version=%d", version)
		}
		if err := rollbackMigration(m.db, migration); err != nil {
			return rolledBack, err
		}
		rolledBack++
	}

	return rolledBack, nil
}

func loadMigrations(migrationsDir string) ([]migration, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations dir %s: %w", migrationsDir, err)
	}

	type partial struct {
		version int
		name    string
		upSQL   string
		downSQL string
	}

	store := make(map[int]*partial)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		matches := migrationNamePattern.FindStringSubmatch(fileName)
		if len(matches) != 4 {
			continue
		}

		version, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid migration version %s: %w", matches[1], err)
		}
		name := matches[2]
		direction := matches[3]

		filePath := filepath.Join(migrationsDir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filePath, err)
		}

		if _, ok := store[version]; !ok {
			store[version] = &partial{version: version, name: name}
		}
		p := store[version]
		if p.name != name {
			return nil, fmt.Errorf("migration name mismatch for version %d", version)
		}

		switch direction {
		case "up":
			p.upSQL = string(content)
		case "down":
			p.downSQL = string(content)
		default:
			return nil, fmt.Errorf("unsupported migration direction: %s", direction)
		}
	}

	if len(store) == 0 {
		return nil, fmt.Errorf("no migration files found in %s", migrationsDir)
	}

	var versions []int
	for version := range store {
		versions = append(versions, version)
	}
	sort.Ints(versions)

	migrations := make([]migration, 0, len(versions))
	for _, version := range versions {
		p := store[version]
		if p.upSQL == "" || p.downSQL == "" {
			return nil, fmt.Errorf("migration %d must contain both up and down sql", version)
		}
		migrations = append(migrations, migration{
			Version: p.version,
			Name:    p.name,
			UpSQL:   p.upSQL,
			DownSQL: p.downSQL,
		})
	}

	return migrations, nil
}

func ensureMigrationTable(db *sql.DB) error {
	const statement = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);`
	_, err := db.Exec(statement)
	return err
}

func appliedMigrationSet(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query(`SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]bool)
	for rows.Next() {
		var version int
		if scanErr := rows.Scan(&version); scanErr != nil {
			return nil, scanErr
		}
		result[version] = true
	}
	return result, rows.Err()
}

func appliedMigrationVersionsDesc(db *sql.DB) ([]int, error) {
	rows, err := db.Query(`SELECT version FROM schema_migrations ORDER BY version DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if scanErr := rows.Scan(&version); scanErr != nil {
			return nil, scanErr
		}
		versions = append(versions, version)
	}
	return versions, rows.Err()
}

func applyMigration(db *sql.DB, migration migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer rollbackSilently(tx)

	if _, err = tx.Exec(migration.UpSQL); err != nil {
		return fmt.Errorf("failed to apply up migration %d_%s: %w", migration.Version, migration.Name, err)
	}
	if _, err = tx.Exec(
		`INSERT INTO schema_migrations(version, name) VALUES($1, $2)`,
		migration.Version,
		migration.Name,
	); err != nil {
		return fmt.Errorf("failed to write migration record %d_%s: %w", migration.Version, migration.Name, err)
	}

	return tx.Commit()
}

func rollbackMigration(db *sql.DB, migration migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer rollbackSilently(tx)

	if _, err = tx.Exec(migration.DownSQL); err != nil {
		return fmt.Errorf("failed to apply down migration %d_%s: %w", migration.Version, migration.Name, err)
	}
	if _, err = tx.Exec(`DELETE FROM schema_migrations WHERE version = $1`, migration.Version); err != nil {
		return fmt.Errorf("failed to delete migration record %d_%s: %w", migration.Version, migration.Name, err)
	}

	return tx.Commit()
}

func rollbackSilently(tx *sql.Tx) {
	_ = tx.Rollback()
}
