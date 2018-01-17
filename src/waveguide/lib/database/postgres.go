package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"time"
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
		"videos":      "video_id SERIAL PRIMARY KEY, video_name VARCHAR(255) NOT NULL, video_uploader INTEGER REFERENCES video_users(user_id) ON DELETE CASCADE, video_upload_date INTEGER NOT NULL, video_description TEXT NOT NULL, video_metainfo_url TEXT NOT NULL",
		"webseeds":    "webseed_url TEXT NOT NULL, video_id INTEGER REFERENCES videos(video_id) ON DELETE CASCADE",
		"oauth_users": "video_user_id INTEGER REFERENCES video_users, oauth_user_id TEXT NOT NULL",
	}

	tables_order := []string{
		"video_users",
		"videos",
		"webseeds",
		"oauth_users",
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

func (db *pgDB) GetFrontpageVideos() (list model.VideoList, err error) {
	var rows *sql.Rows
	rows, err = db.conn.Query("SELECT video_id, video_name, video_upload_date, video_metainfo_url FROM videos ORDER BY video_upload_date DESC LIMIT 10")
	if err == sql.ErrNoRows {
		err = nil
		return
	} else if err == nil {
		for rows.Next() {
			var info model.VideoInfo
			var id int64
			rows.Scan(&id, &info.Title, &info.UploadedAt, &info.TorrentURL)
			info.VideoID = fmt.Sprintf("%d", id)
			list.Videos = append(list.Videos, info)
			if list.LastUpdated.Unix() < info.UploadedAt {
				list.LastUpdated = time.Unix(info.UploadedAt, 0)
			}
		}
		rows.Close()
	}
	return
}

func (db *pgDB) RegisterVideo(video *model.VideoInfo) error {
	return db.conn.QueryRow("INSERT INTO videos (video_name, video_description, video_upload_date, video_metainfo_url) VALUES ($1, $2, $3, $4) RETURNING video_id", video.Title, video.Description, video.UploadedAt, "").Scan(&video.VideoID)
}

func (db *pgDB) SetVideoMetaInfo(idstr string, url string) (err error) {
	var id int64
	id, err = strconv.ParseInt(idstr, 10, 64)
	if err == nil {
		_, err = db.conn.Exec("UPDATE videos SET video_metainfo_url = $1 WHERE video_id = $2", url, id)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

func (db *pgDB) AddWebseed(idstr string, url string) (err error) {
	var id int64
	id, err = strconv.ParseInt(idstr, 10, 64)
	if err == nil {
		_, err = db.conn.Exec("INSERT INTO webseeds(video_id, webseed_url) VALUES ($1, $2)", id, url)
	}
	return
}

func (db *pgDB) GetVideoInfo(idstr string) (info *model.VideoInfo, err error) {

	var id int64
	id, err = strconv.ParseInt(idstr, 10, 64)
	if err == nil {
		info = &model.VideoInfo{
			VideoID: idstr,
		}
		err = db.conn.QueryRow("SELECT video_name, video_description, video_upload_date, video_metainfo_url FROM videos WHERE video_id = $1", id).Scan(&info.Title, &info.Description, &info.UploadedAt, &info.TorrentURL)
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
	}
	return
}

func (db *pgDB) GetUserByName(name string) (user *model.UserInfo, err error) {
	user = new(model.UserInfo)
	err = db.conn.QueryRow("SELECT user_id, user_name FROM video_users WHERE user_name = $1", name).Scan(&user.UserID, &user.Name)
	if err == sql.ErrNoRows {
		user = nil
		err = nil
	} else if err != nil {
		user = nil
	}
	return
}

func (db *pgDB) GetVideosForUserByName(name string) (list *model.VideoFeed, err error) {
	var user *model.UserInfo
	user, err = db.GetUserByName(name)
	if user != nil {
		var rows *sql.Rows
		rows, err = db.conn.Query("SELECT video_id, video_name, video_description, video_upload_date, video_metainfo_url FROM videos WHERE video_uploader = $1", user.UserID)
		if err == sql.ErrNoRows {
			err = nil
		} else if err == nil {
			list = new(model.VideoFeed)
			list.Owner = user
			for rows.Next() {
				var info model.VideoInfo
				info.UserID = user.UserID
				var id int64
				rows.Scan(&id, &info.Title, &info.Description, &info.UploadedAt, &info.TorrentURL)
				info.VideoID = fmt.Sprintf("%d", id)
				list.List.Videos = append(list.List.Videos, info)
			}
			rows.Close()
		}
	}
	return
}
