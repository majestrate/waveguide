package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"waveguide/lib/log"
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

func (db *pgDB) Close() error {
	return db.conn.Close()
}

func (db *pgDB) CreateTables() (err error) {
	tables := map[string]string{
		"video_users": "user_id SERIAL PRIMARY KEY, user_name VARCHAR(255) NOT NULL, user_email VARCHAR(255) NOT NULL, user_logincred VARCHAR(255) NOT NULL",
		"videos":      "video_id SERIAL PRIMARY KEY, video_name VARCHAR(255) NOT NULL, video_uploader INTEGER REFERENCES video_users(user_id) ON DELETE CASCADE, video_upload_date INTEGER NOT NULL, video_description TEXT NOT NULL, video_torrent_url VARCHAR(255) NOT NULL",
		"webseeds":    "webseed_url TEXT NOT NULL, video_id INTEGER REFERENCES videos(video_id) ON DELETE CASCADE",
	}

	tables_order := []string{
		"video_users",
		"videos",
		"webseeds",
	}

	for _, name := range tables_order {
		q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", name, tables[name])
		log.Debugf("sql init: %s", q)
		_, err = db.conn.Exec(q)
		if err != nil {
			return
		}
	}
	return
}

func (db *pgDB) NextVideoID() (id int64, err error) {
	err = db.conn.QueryRow("SELECT coalesce(max(video_id),0) FROM videos").Scan(&id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (db *pgDB) GetFrontpageVideos() (list model.VideoList, err error) {
	var rows *sql.Rows
	rows, err = db.conn.Query("SELECT video_id, video_name, video_upload_date FROM videos ORDER BY video_upload_date DESC LIMIT 10")
	if err == sql.ErrNoRows {
		err = nil
		return
	} else if err == nil {
		for rows.Next() {
			var info model.VideoInfo
			rows.Scan(&info.VideoID, &info.Title, &info.UploadedAt)
			list = append(list, info)
		}
		rows.Close()
	}
	return
}

func (db *pgDB) RegisterVideo(video *model.VideoInfo) error {
	result, err := db.conn.Exec("INSERT INTO videos (video_name, video_description, video_upload_date, video_torrent_url) VALUES ( $1, $2, $3, $4)", video.Title, video.Description, video.UploadedAt, video.TorrentURL)
	if err == nil {
		video.VideoID, err = result.LastInsertId()
		if err == nil {
			for idx := range video.WebSeeds {
				_, err = db.conn.Exec("INSERT INTO webseeds(video_id, webseed_url) VALUES ($1, $2)", video.VideoID, video.WebSeeds[idx])
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

func (db *pgDB) GetVideoInfo(id int64) (info *model.VideoInfo, err error) {
	info = new(model.VideoInfo)
	err = db.conn.QueryRow("SELECT video_id, video_name, video_description, video_upload_date, video_torrent_url FROM videos WHERE video_id = $1", id).Scan(&info.VideoID, &info.Title, &info.Description, &info.UploadedAt, &info.TorrentURL)
	if err == nil {
		var rows *sql.Rows
		rows, err = db.conn.Query("SELECT webseed_url FROM webseeds WHERE video_id = $1", id)
		if err == nil {
			for rows.Next() {
				var url string
				rows.Scan(&url)
				info.WebSeeds = append(info.WebSeeds, url)
			}
			rows.Close()
		}
	} else {
		info = nil
	}
	return
}
