package bot

import (
	"strings"
)

type CommandHandler interface {
	Handle(bs *botSession, command string, args ...string) bool
}

type FuncCommandHandler func(bs *botSession, command string, args ...string) bool

func (f FuncCommandHandler) Handle(bs *botSession, command string, args ...string) bool {
	return f(bs, command, args...)
}

type HandlerMap map[string]CommandHandler

func (hm HandlerMap) Handle(bs *botSession, command string, args ...string) bool {
	cmd, ok := hm[command]
	if !ok {
		return false
	}
	return cmd.Handle(bs, command, args...)
}

func (hm HandlerMap) Set(command string, sc CommandHandler) {
	hm[strings.TrimPrefix(command, "/")] = sc
}
