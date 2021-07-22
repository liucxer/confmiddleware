package confmiddleware

import (
	"context"
	"github.com/go-courier/sqlx/v2"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	pg := &Postgres{
		Host: "10.0.9.212",
		Password: "123456",
		User: "postgres",
		Database: &sqlx.Database{
			Name: "osm",
		},
	}
	pg.SetDefaults()
	pg.Init()

	_, err := pg.QueryContext(context.Background(), "SELECT 1")
	require.NoError(t, err)
}
