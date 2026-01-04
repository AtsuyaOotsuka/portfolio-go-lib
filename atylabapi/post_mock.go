package atylabapi

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type ApiPostStructMock struct {
	client  *http.Client
	apiKey  string
	baseUrl string
	mock.Mock
}

func (s *ApiPostStructMock) Post(
	path string,
	body map[string][]string,
	ctx *ApiCtxSvc,
) ([]byte, error) {
	args := s.Called(path, body, ctx)
	return args.Get(0).([]byte), args.Error(1)
}
