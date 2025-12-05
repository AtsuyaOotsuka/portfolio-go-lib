package atylabjwt

import "github.com/stretchr/testify/mock"

type JwtMock struct {
	mock.Mock
}

func (m *JwtMock) CreateJwt(config *JwtConfig) (string, error) {
	args := m.Called(config)
	return args.String(0), args.Error(1)
}

func (m *JwtMock) Validate(jwtSecret string, jwtToken string) error {
	args := m.Called(jwtSecret, jwtToken)
	return args.Error(0)
}

func (m *JwtMock) GetUUID() string {
	args := m.Called()
	return args.String(0)
}

func (m *JwtMock) GetEmail() string {
	args := m.Called()
	return args.String(0)
}
