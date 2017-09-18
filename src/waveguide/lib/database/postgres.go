package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"waveguide/lib/model"
)

type pgDB struct {
	url  string
	conn *sql.DB
}

func (db *pgDB) Init() (err error) {
	db.conn, err = sql.Open("postgres", db.url)
	if err == nil {
		err = db.CreateTables()
	}
	return
}

func (db *pgDB) CreateTables() (err error) {
	tables := map[string]string{
		"videos": "video_id SERIAL PRIMARY KEY, video_name TEXT NOT NULL, video_uploader INTEGER",
	}
	for name := range tables {
		q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", name, tables[name])
		_, err = db.conn.Exec(q)
		if err != nil {
			return
		}
	}
	return
}

func (db *pgDB) GetFrontpageVideos() (list model.VideoList, err error) {
	return
}
