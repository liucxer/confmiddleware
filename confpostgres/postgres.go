package confpostgres

import (
	"fmt"
	"github.com/liucxer/courier/sqlx/postgresqlconnector"
	"time"

	"github.com/liucxer/courier/envconf"
	"github.com/liucxer/courier/sqlx"
)

type Postgres struct {
	Host            string `env:",upstream"`
	Port            int
	User            string           `env:""`
	Password        envconf.Password `env:""`
	Extra           string
	Extensions      []string
	PoolSize        int
	ConnMaxLifetime envconf.Duration
	Retry
	Database *sqlx.Database `env:"-"`
	*sqlx.DB `env:"-"`
}

func (m *Postgres) SetDefaults() {
	if m.Host == "" {
		m.Host = "127.0.0.1"
	}
	if m.Port == 0 {
		m.Port = 5432
	}

	if m.PoolSize == 0 {
		m.PoolSize = 10
	}

	if m.ConnMaxLifetime == 0 {
		m.ConnMaxLifetime = envconf.Duration(1 * time.Hour)
	}

	if m.Extra == "" {
		m.Extra = "sslmode=disable"
	}
}

func (m *Postgres) URL() string {
	password := m.Password
	if password != "" {
		password = ":" + password
	}
	return fmt.Sprintf("postgres://%s%s@%s:%d", m.User, password, m.Host, m.Port)
}

func (m *Postgres) Conn() error {
	m.SetDefaults()
	db := m.Database.OpenDB(&postgresqlconnector.PostgreSQLConnector{
		Host:       m.URL(),
		Extra:      m.Extra,
		Extensions: m.Extensions,
	})
	db.SetMaxOpenConns(m.PoolSize)
	db.SetMaxIdleConns(m.PoolSize / 2)
	db.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
	m.DB = db
	return nil
}

func (m *Postgres) Init() {
	if m.DB == nil {
		m.Do(m.Conn)
	}
}
