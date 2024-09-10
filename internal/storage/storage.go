package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/frairon/linkbot/internal/storage/models"
	"github.com/frairon/linkbot/internal/storage/schema"
	migrate_sqlite "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Storage struct {
	db *sql.DB
}

var ErrDuplicate = errors.New("duplicate database entry")

func NewStorage(dbConn string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbConn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	boil.DebugMode = true

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping database: %v", err)
	}
	driver, err := migrate_sqlite.WithInstance(db, &migrate_sqlite.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver")
	}
	err = schema.ApplySchema(false, driver)
	if err != nil {
		return nil, fmt.Errorf("error applying schema. %v", err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) ctx() context.Context {
	return context.Background()
}

func (s *Storage) UserExists(userId int64) bool {
	exists, err := models.UserExists(s.ctx(), s.db, userId)
	if err != nil {
		log.Printf("error finding user: %v", err)
	}

	return exists
}

func (s *Storage) AddUser(userId int64, name string) error {
	exists, err := models.UserExists(s.ctx(), s.db, userId)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}
	log.Printf("userId: %d", userId)
	u := models.User{ID: userId, Name: null.StringFrom(name)}

	if exists {
		log.Printf("user exists.")
		_, err = u.Update(s.ctx(), s.db, boil.Infer())
		return err
	}
	log.Printf("user does not exist, will create: %#v", u)
	// insert user
	return u.Insert(s.ctx(), s.db, boil.Infer())
}

func (s *Storage) ListUsers() ([]*models.User, error) {
	return models.Users().All(s.ctx(), s.db)
}

func (s *Storage) GetUser(userId int64) (*models.User, error) {
	return models.FindUser(s.ctx(), s.db, userId)
}

func (s *Storage) DeleteUser(user *models.User) error {
	_, err := user.Delete(s.ctx(), s.db)
	return err
}

func (s *Storage) UserSessions() ([]*models.UserSession, error) {
	return models.UserSessions().All(s.ctx(), s.db)
}

func (s *Storage) StoreSession(userId int64, chatId int64, lastUserAction time.Time, sessionData string) error {
	sess, err := models.FindUserSession(s.ctx(), s.db, userId)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error finding existing session: %v", err)
	}

	if sess == nil {
		sess = &models.UserSession{
			UserID:         userId,
			ChatID:         chatId,
			LastUserAction: null.TimeFrom(lastUserAction),
			Data:           null.StringFrom(sessionData),
		}
		return sess.Insert(s.ctx(), s.db, boil.Infer())
	}

	sess.ChatID = chatId
	sess.LastUserAction = null.TimeFrom(lastUserAction)
	sess.Data = null.StringFrom(sessionData)
	_, err = sess.Update(s.ctx(), s.db, boil.Infer())
	return err
}

const residentMacsKey = "resident-macs"

func (s *Storage) GetResidentMacs() (*[]string, error) {
	return GetSettingByKey[[]string](s, residentMacsKey)
}

func (s *Storage) SetResidentMacs(macs []string) error {
	return SetSetting(s, residentMacsKey, macs)
}

func GetSettingByKey[T any](s *Storage, key string) (*T, error) {
	var target T
	setting, err := models.FindSetting(s.ctx(), s.db, null.StringFrom(key))
	if err == sql.ErrNoRows {
		return &target, nil
	}

	return &target, json.Unmarshal([]byte(setting.Value.String), &target)
}

func SetSetting[T any](s *Storage, key string, value T) error {
	marshalled, err := json.Marshal(value)
	if err != nil {
		return err
	}
	setting := &models.Setting{
		Key:   null.StringFrom(key),
		Value: null.StringFrom(string(marshalled)),
	}

	exists, err := models.Settings(qm.Where("key=?", key)).Exists(s.ctx(), s.db)
	if err != nil {
		return err
	}
	if !exists {
		return setting.Insert(s.ctx(), s.db, boil.Infer())
	}
	_, err = setting.Update(s.ctx(), s.db, boil.Infer())
	return err
}
