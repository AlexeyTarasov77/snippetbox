package constants

type CtxKey string

const (
	FlashCtxKey = CtxKey("flash")
	UserCtxKey = CtxKey("user")
	UserIDCtxKey = CtxKey("userID")
)