package middleware

type ctxKey string

const (
	//UserCtx is a context key for user object
	UserCtx ctxKey = "UserEntry"

	//BoardCtx is a context key for board object
	BoardCtx ctxKey = "BoardEntry"
)
