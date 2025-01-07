package storage

import (
	"fmt"
	"time"

	"github.com/frairon/linkbot/internal/storage/models"
	"github.com/gofrs/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type CategoryCounts struct {
	Category string
	Count    int
}

func (s *Storage) ListCategories(userID int64) ([]CategoryCounts, error) {
	rows, err := s.db.Query("select category, count(*) from `user_links` where user_id=$1 and hidden != true group by 1 order by added DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("error reading categories: %w", err)
	}
	var categories []CategoryCounts
	for rows.Next() {
		var cc CategoryCounts
		if err := rows.Scan(&cc.Category, &cc.Count); err != nil {
			return nil, fmt.Errorf("error reading category: %w", err)
		}
		categories = append(categories, cc)
	}
	return categories, nil
}

func (s *Storage) AddLink(userId int64, category string, link string, headline string) error {
	linkId, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("error generating link-id: %w", err)
	}
	ul := &models.UserLink{
		LinkID:   linkId.String(),
		UserID:   userId,
		Category: category,
		Link:     link,
		Headline: headline,
		Added:    time.Now(),
	}

	return ul.Insert(s.db, boil.Infer())
}
