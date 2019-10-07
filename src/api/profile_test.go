package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/playneta/go-sessions/src/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile(t *testing.T) {
	suite := newTestSuite(t, http.MethodGet, nil, nil)
	suite.authorize()
	defer suite.close()

	// Calling instance
	err := suite.api.Profile(suite.context)
	require.NoError(t, err)

	// Checking response body
	{
		require.Equal(t, http.StatusOK, suite.recorder.Code)
		checkUser := new(models.User)
		err := json.NewDecoder(suite.recorder.Body).Decode(checkUser)
		require.NoError(t, err)
		assert.Equal(t, suite.user, checkUser)
	}
}
