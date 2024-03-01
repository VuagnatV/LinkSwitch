package services_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"linkSwitch/database"
	"linkSwitch/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	db, err := database.InitMongo("mongodb://localhost:27017")
	require.NoError(t, err)

	contentType := "application/json"

	server := httptest.NewServer(services.RunServer(db))
	defer server.Close()

	created := database.URL{
		Long: "https://alexsoyes.com/solid/",
	}

	payload, err := serializeURL(created)
	require.NoError(t, err)

	response, err := http.Post(server.URL, contentType, payload)
	require.NoError(t, err)
	require.NotNil(t, response)

	defer response.Body.Close()

	var url database.URL

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &url)
	if err != nil {
		return
	}

	fmt.Println("Response body:", string(body))

	response, err = http.Get(server.URL + "/" + "url.Short")
	require.NoError(t, err)
	require.NotNil(t, response)

	response, err = http.Get(server.URL + "/stats/")
	require.NoError(t, err)
	require.NotNil(t, response)

	response, err = http.Get(server.URL + "/stats/" + url.Short)
	require.NoError(t, err)
	require.NotNil(t, response)

}

func serializeURL(url database.URL) (io.Reader, error) {
	data, err := json.Marshal(url)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
