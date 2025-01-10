package database

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/may20xx/booking/config"
	"github.com/may20xx/booking/pkg/log"
)

var (
	db     *sqlx.DB
	dbOnce sync.Once
	dbErr  error
)

const (
	maxRetries        = 5
	retryDelay        = 2 * time.Second
	connectionTimeout = 10 * time.Second
	maxOpenConns      = 25
	maxIdleConns      = 25
	connMaxLifetime   = 5 * time.Minute
)

func InitDatabase() error {
	dbOnce.Do(func() {
		setting := config.GetConfig()

		port, err := strconv.Atoi(setting.DBPort)
		if err != nil {
			log.Msg.Panic(err)
		}
		setting.DBPort = strconv.Itoa(port)

		connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			setting.DBHost, port, setting.DBUser, setting.DBPassword, setting.DBName)

		db, dbErr = connectDBWithRetry(connectionString, maxRetries)
		if dbErr == nil {
			configureConnectionPool(db)
		}
	})
	return dbErr
}

func connectDBWithRetry(conn string, maxRetries int) (*sqlx.DB, error) {
	for i := 0; i < maxRetries; i++ {
		db, err := sqlx.Connect("postgres", conn)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
			err = db.PingContext(ctx)
			cancel()
			if err == nil {
				log.Msg.Info("Connected to PostgreSQL successfully! ✅")
				return db, nil
			}
		}
		log.Msg.Warnf("Failed to connect to database, retrying (%d/%d)...", i+1, maxRetries)
		time.Sleep(retryDelay)
	}
	return nil, fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

func configureConnectionPool(db *sqlx.DB) {
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
}

func GetDatabase(ctx context.Context) (*sqlx.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database connection lost: %w", err)
	}
	return db, nil
}

func CloseDatabase() error {
	if db != nil {
		log.Msg.Info("Closing database connection...")
		return db.Close()
	}
	return nil
}

func GracefulShutdown(ctx context.Context) {
	if err := CloseDatabase(); err != nil {
		log.Msg.Error("Error closing database connection: ", err)
	}
}
