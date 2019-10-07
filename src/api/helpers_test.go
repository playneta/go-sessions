package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/playneta/go-sessions/src/models"
	mock_repositories "github.com/playneta/go-sessions/src/repositories/mocks"
	mock_services "github.com/playneta/go-sessions/src/services/mocks"
	"go.uber.org/zap"
)

type suite struct {
	gmock          *gomock.Controller
	userRepo       *mock_repositories.MockUser
	accountService *mock_services.MockAccount
	user           *models.User
	context        echo.Context
	request        *http.Request
	recorder       *httptest.ResponseRecorder
	api            *API
}

func newTestSuite(t *testing.T, method string, body io.Reader, headers map[string]string) *suite {
	// Creating mocks
	ctrl := gomock.NewController(t)

	userRepo := mock_repositories.NewMockUser(ctrl)
	accountService := mock_services.NewMockAccount(ctrl)

	// Basic setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Creating api instance
	logger := zap.NewNop().Sugar()

	// Creating instance of API
	a := &API{
		logger:         logger,
		accountService: accountService,
		userRepo:       userRepo,
	}

	return &suite{
		gmock:          ctrl,
		userRepo:       userRepo,
		accountService: accountService,
		request:        req,
		recorder:       rec,
		context:        c,
		api:            a,
	}
}

func (s *suite) authorize() {
	ts, _ := time.Parse(time.RFC1123, time.RFC1123)

	// Authorizing user
	user := &models.User{
		ID:        1,
		Email:     "user@example.com",
		Password:  "my_tokenized_password",
		Token:     "",
		CreatedAt: ts,
		UpdatedAt: ts,
	}
	s.user = user
	s.context.Set("user", user)
}

func (s *suite) close() {
	s.gmock.Finish()
}
