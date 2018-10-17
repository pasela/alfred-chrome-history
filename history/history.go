package history

import (
	"database/sql"
	"net/url"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pasela/alfred-chrome-history/profile"
	"github.com/pasela/alfred-chrome-history/utils"
)

const (
	dbFileName = "History"
	escapeChar = '\\'
)

type Entry struct {
	URL           string
	Title         string
	VisitCount    int
	TypedCount    int
	LastVisitTime time.Time
	Hidden        int
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
			url, title, visit_count, typed_count, last_visit_time, hidden
		FROM
			urls
		WHERE
			(title LIKE ? ESCAPE ? OR url LIKE ? ESCAPE ?)
			AND hidden = 0
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
		var visit int64
		err := rows.Scan(
			&entry.URL,
			&entry.Title,
			&entry.VisitCount,
			&entry.TypedCount,
			&visit,
			&entry.Hidden,
		)
		if err != nil {
			return nil, err
		}
		entry.LastVisitTime = convertChromeTime(visit)
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// https://code.google.com/p/chromium/codesearch#chromium/src/base/time/time.h
func convertChromeTime(msec int64) time.Time {
	sec := msec / 1000000
	nsec := (msec % 1000000) * 1000

	t := time.Unix(sec, nsec)
	t = t.AddDate(-369, 0, 0) // 369 = 1970 - 1601
	return t
}
