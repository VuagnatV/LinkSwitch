package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	db, err := InitMongo("mongodb://localhost:27017")
	require.NoError(t, err)

	created := URL{
		Long:           "https://alexsoyes.com/solide/",
		Short:          "testeee",
		Clicks:         0,
		ExpirationDate: time.Now().UTC().Truncate(time.Minute),
	}

	url, err := db.InsertOne(created)
	require.NoError(t, err)
	require.Equal(t, created, *url)

	url, err = db.FindOne(url.Short)
	require.NoError(t, err)
	require.Equal(t, created, *url)

}
