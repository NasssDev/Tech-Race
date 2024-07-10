package globals

import (
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("session-store"))
