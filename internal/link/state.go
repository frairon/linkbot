package link

import (
	"github.com/frairon/botty"
	"github.com/frairon/linkbot/internal/storage"
)

// some aliases so we don't have to keep repeating it
type (
	LinkState   = botty.State[*State]
	LinkSession = botty.Session[*State]
)

type State struct {
	st *storage.Storage
}
