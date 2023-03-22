package v2ex

const APIBaseURL string = "https://www.v2ex.com/api/v2/"

const (
	TokenScopeRegular    TokenScope = "regular"
	TokenScopeEverything TokenScope = "everything"
)

const (
	TokenExpiration30Days  TokenExpiration = 2592000
	TokenExpiration60Days  TokenExpiration = 5184000
	TokenExpiration90Days  TokenExpiration = 7776000
	TokenExpiration180Days TokenExpiration = 15552000
)
