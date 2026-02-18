package mysql

import (
	"database/sql"

	"github.com/Dagime-Teshome/snippetbox/pkg/models"
)

// make functions to create , get , delete snippets from database

type SnippetModel struct {
	Db *sql.DB
}

func (sm *SnippetModel) Insert(title, content, expires string) (int, error) {
	sql := `INSERT INTO snippets (title,content,created,expires) VALUES (?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`
	result, err := sm.Db.Exec(sql, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (sm *SnippetModel) GetById(id int) (*models.Snippet, error) {
	query := `SELECT * FROM snippets where id = ?`
	snippet := models.Snippet{}
	row := sm.Db.QueryRow(query, id)
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return &snippet, nil
}

func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	query := `	SELECT id, title, content, created, expires FROM snippets
				WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := sm.Db.Query(query)
	if err != nil {
		return nil, err
	}
	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
