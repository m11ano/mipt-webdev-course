package auth

type Client interface {
	ParseJWT(tokenStr string) (*AuthClaims, error)
}
