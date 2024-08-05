package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ahmedkhaeld/goblend/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	PostgreSQL *pgxpool.Pool
	MariaDB    *sql.DB
	MongoDB    *mongo.Client
	Driver     string
}

func NewDatabase() *Database {
	return &Database{
		Driver: config.Env.DatabaseDriver,
	}
}

func (d *Database) Connect() error {
	var err error

	switch d.Driver {
	case "postgres":
		d.PostgreSQL, err = connectPostgreSQL()
	case "mariadb":
		d.MariaDB, err = connectMariaDB()
	case "mongodb":
		d.MongoDB, err = connectMongoDB()
	default:
		return fmt.Errorf("unsupported database driver: %s", d.Driver)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", d.Driver, err)
	}

	return nil
}

func connectPostgreSQL() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Env.DatabaseUser, config.Env.DatabasePass,
		config.Env.DatabaseHost, config.Env.DatabasePort,
		config.Env.DatabaseName, config.Env.DatabaseSSLMode)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL connection string: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return pool, nil
}

func connectMariaDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Env.DatabaseUser, config.Env.DatabasePass,
		config.Env.DatabaseHost, config.Env.DatabasePort, config.Env.DatabaseName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MariaDB connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping MariaDB: %w", err)
	}

	return db, nil
}

func connectMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.Env.MongoDBURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

func (d *Database) Close() {
	if d.PostgreSQL != nil {
		d.PostgreSQL.Close()
	}
	if d.MariaDB != nil {
		err := d.MariaDB.Close()
		if err != nil {
			fmt.Printf("error closing MariaDB connection: %v\n", err)
		}
	}
	if d.MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := d.MongoDB.Disconnect(ctx)
		if err != nil {
			fmt.Printf("error disconnecting MongoDB: %v\n", err)
		}
	}
}

// Additional methods for database operations can be added here
