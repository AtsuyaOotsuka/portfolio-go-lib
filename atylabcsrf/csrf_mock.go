package atylabcsrf

type CsrfPkgMockStruct struct{}

func (c *CsrfPkgMockStruct) GenerateNonceString() string {
	return "mocked_nonce"
}

func (c *CsrfPkgMockStruct) GenerateCSRFCookieToken(secret string, timestamp int64, nonceStr string) string {
	return "mocked_csrf_token"
}

func (c *CsrfPkgMockStruct) ValidateCSRFCookieToken(token string, secret string, timestamp int64) error {
	return nil
}
