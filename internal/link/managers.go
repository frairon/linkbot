package link

import (
	"github.com/frairon/botty"
	"github.com/frairon/linkbot/internal/storage"
)

type sessionContextManager struct {
	st *storage.Storage
}

func NewManagers(st *storage.Storage) (botty.SessionContextManager[*State], botty.UserManager) {
	return &sessionContextManager{
			st: st,
		}, &userManager{
			st: st,
		}
}

func (scm *sessionContextManager) CreateSessionContext(userId, chatId int64) *State {
	return &State{}
}
func (scm *sessionContextManager) StoreSession(session botty.StoredSession[*State]) error {
	// TODO store them
	return nil
}
func (scm *sessionContextManager) LoadSessions() ([]botty.StoredSession[*State], error) {
	// TODO load them
	return nil, nil
}

type userManager struct {
	st *storage.Storage
}

func (um *userManager) ListUsers() ([]botty.User, error) {
	users, err := um.st.ListUsers()
	if err != nil {
		return nil, err
	}
	var botUsers []botty.User
	for _, user := range users {
		botUsers = append(botUsers, botty.User{
			ID:   user.ID,
			Name: user.Name,
		})
	}
	return botUsers, nil
}

func (um *userManager) AddUser(userID int64, userName string) error {
	return um.st.AddUser(userID, userName)
}
func (um *userManager) UserExists(userID int64) bool {
	return um.st.UserExists(userID)
}
func (um *userManager) DeleteUser(userID int64) error {
	usr, err := um.st.GetUser(userID)
	if err != nil {
		return err
	}
	return um.st.DeleteUser(usr)
}
