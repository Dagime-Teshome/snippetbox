package mock

import (
	"time"

	"github.com/Dagime-Teshome/snippetbox/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type snippetModel struct{}

func (s *snippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (s *snippetModel) GetById(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (s *snippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
