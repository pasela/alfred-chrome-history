package history

import (
	"database/sql"
	"net/url"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pasela/alfred-chrome-history/profile"
	"github.com/pasela/alfred-chrome-history/utils"
)

const (
	dbFileName = "History"
	escapeChar = '\\'
)

type Entry struct {
	URL   string
	Title string
}

type History struct {
	db *sql.DB
}

func GetHistoryPath(profileName string) (string, error) {
	path, err := profile.GetProfilePath(profileName)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, dbFileName), nil
}

func Open(file string) (*History, error) {
	dsn, err := url.Parse(file)
	if err != nil {
		return nil, err
	}
	dsn.Scheme = "file"
	q := dsn.Query()
	q.Set("mode", "ro")
	q.Set("immutable", "1")
	q.Set("_query_only", "1")
	dsn.RawQuery = q.Encode()

	db, err := sql.Open("sqlite3", dsn.String())
	if err != nil {
		return nil, err
	}

	return &History{db}, nil
}

func (h *History) Close() error {
	return h.db.Close()
}

func (h *History) Query(url, title string) ([]Entry, error) {
	u := "%" + utils.EscapeLike(url, escapeChar) + "%"
	t := "%" + utils.EscapeLike(title, escapeChar) + "%"

	rows, err := h.db.Query(`
		SELECT
			url, title
		FROM
			urls
		WHERE
			title LIKE ? ESCAPE ? OR url LIKE ? ESCAPE ?
		ORDER BY
			visit_count DESC, typed_count DESC, last_visit_time DESC
	`, u, string(escapeChar), t, string(escapeChar))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]Entry, 0)
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.URL, &entry.Title); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
