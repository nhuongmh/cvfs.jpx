package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/pkg/errors"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/schema.sql
var schema string

type DB struct {
	SqlDB        *sql.DB
	QueryBuilder *squirrel.StatementBuilderType
	url          string
}

func ConnectDB(ctx context.Context, dbFileUrl string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbFileUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Failed open database file")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "Failed ping Database")
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &DB{
		db,
		&psql,
		dbFileUrl,
	}, nil

}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *DB) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.SqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		logger.Log.Fatal().Err(err).Msgf("db down") // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.SqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *DB) Close() error {
	logger.Log.Info().Msgf("Disconnected from database: %s", s.url)
	return s.SqlDB.Close()
}

func (db *DB) Migrate() error {
	_, err := db.SqlDB.Exec(schema)
	if err != nil {
		return errors.Wrap(err, "Failed migrate database")
	}

	return nil
}
