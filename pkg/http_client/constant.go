package httpclient

type AuthMethod string

const (
	AuthMethodJWT          AuthMethod = "JWT"
	AuthMethodAPIKey       AuthMethod = "APIKey"
	AuthMethodBearerAPIKey AuthMethod = "BearerAPIKey"
)

type AuthHeader struct {
	Method AuthMethod
	Token  string
}
