package constant

import "time"

const (
	LIMIT_LOGIN      = "limit:login:"
	BLACK_LIST       = "black_list:"
	BLACK_LIST_COUNT = "black_list:count:"
	HEARTBEAT        = "heartbeat:"
)

const (
	TTL_LIMIT_lOGIN = 5 * time.Minute
	HEARTBEAT_TTL   = 1 * time.Minute
)
