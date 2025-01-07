package link

import (
	"log"
	"testing"

	"github.com/frairon/botty"
	"github.com/frairon/linkbot/internal/storage"
	"github.com/stretchr/testify/require"
)

const userId = 5

func TestHome(t *testing.T) {

	state := Home()

	bot, st := setupMockBot(t)
	require.NoError(t, bot.Err())

	session, err := bot.CreateSession(userId)
	require.NoError(t, err)

	require.NoError(t, st.AddLink(userId, "category", "google.com", "headline"))

	// this will actiate home state
	session.ReplaceState(state)

	log.Printf("%s", bot.LastMessage.Text)
	require.Contains(t, bot.LastMessage.Text, "Categories:")

	bot.Stop()
}

func setupMockBot(t *testing.T) (*botty.MockBot[*State], *storage.Storage) {
	st, err := storage.NewStorage(":memory:")
	require.NoError(t, err)

	err = st.AddUser(userId, "test")
	require.NoError(t, err)

	sessionManager, userManager := NewManagers(st)
	botCfg := botty.NewConfig("", sessionManager, userManager, Home)
	b, err := botty.NewMockBot(botCfg)
	require.NoError(t, err)
	return b, st
}
