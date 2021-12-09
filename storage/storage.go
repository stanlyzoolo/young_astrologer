package storage

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type Role struct {
	Name     string
	Password string
	Database string
}

// New returns sql string to create new role with password.
func (r *Role) New() string {
	return fmt.Sprintf("CREATE USER %s WITH ENCRYPTED PASSWORD '%s'",
		r.Name, r.Password)
}

// Alter retyrns sql string to alter user`s permissions.
func (r *Role) Alter() string {
	return fmt.Sprintf("ALTER USER %s WITH SUPERUSER", r.Name)
}

// CreateDB returns sql string to create a database.
func (r *Role) CreateDB() string {
	return fmt.Sprintf("CREATE DATABASE %s", r.Database)
}

// Grant returns sql string to grant privileges on database to new role.
func (r *Role) Grant() string {
	return fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", r.Database, r.Name)
}

// CreateTable returns sql string to create table with required fields.
func (r *Role) CreateTable() string {
	query := fmt.Sprintf(`CREATE TABLE %s (
		unique_id serial,
		copyright varchar,
		date varchar,
		explanation varchar,
		hdurl varchar,
		media_type varchar,
		service_version varchar,
		title varchar,
		url varchar,
		image bytea,
		PRIMARY KEY (unique_id)
	)`, os.Getenv("PSQLNASATABLE"))

	return query
}

// IsExist checks existence of username with name from .env.
func (r *Role) IsExist(ctx context.Context, conn *pgx.Conn, logger *zap.Logger) bool {
	query := fmt.Sprintf("SELECT 1 FROM pg_roles WHERE rolname='%s'", r.Name)
	result, err := conn.Exec(ctx, query)
	
	if err != nil {
		logger.Error("sql execution failed -> get roles",
			zap.String("package", ""),
			zap.String("func", "IsExist"),
			zap.Error(err))
	}

	info := strings.Split(string(result), " ")
	if info[1] != "0" {
		return false
	}

	return true
}
