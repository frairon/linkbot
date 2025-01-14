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
	log.Printf("buttons: %v", bot.LastMessageButtons())

	require.Contains(t, bot.LastMessage.Text, "Categories:")

	bot.Send(userId, "back")

	bot.Stop()
}

func TestAddLink(t *testing.T) {
	bot, st := setupMockBot(t)
	require.NoError(t, bot.Err())
	defer bot.Stop()

	session, err := bot.CreateSession(userId)
	require.NoError(t, err)

	// this will actiate home state
	session.ReplaceState(Home())
	bot.Send(userId, "add")

	require.Contains(t, bot.LastMessage.Text, "Paste the link")
	bot.Send(userId, "google.com")

	require.Contains(t, bot.LastMessage.Text, "try again")

	bot.Send(userId, "https://en.wikipedia.org/wiki/Go_(programming_language)")
	log.Printf("%s", bot.LastMessage.Text)

	bot.Send(userId, "Add")
	cats, err := st.ListCategories(userId)
	require.NoError(t, err)
	require.Len(t, cats, 1)
	require.Equal(t, "somecategory", cats[0].Category)
	require.Equal(t, 1, cats[0].Count)
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
