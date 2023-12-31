package repository

import (
	"context"
	"lct/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type videoPgRepository struct {
	db *pgxpool.Pool
}

func NewVideoPgRepository(db *pgxpool.Pool) (model.VideoRepository, error) {
	ctx := context.Background()

	_, err := db.Exec(ctx, `
		create table if not exists `+model.VideosTableName+`(
			id serial primary key,
			title text not null unique,
			source text not null unique,
			processedSource text default '',
			status text default 'processing',
			createdAt timestamp default current_timestamp,
			updatedAt timestamp default current_timestamp
		);
	`)
	if err != nil {
		return nil, err
	}

	return &videoPgRepository{
		db: db,
	}, nil
}

func (vr *videoPgRepository) InsertOne(c context.Context, videoData model.VideoCreate) error {
	tx, err := vr.db.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	var videoId int
	err = tx.QueryRow(c, `
		insert into `+model.VideosTableName+`(title, source)
		values($1, $2)
		returning id
	`, videoData.Title, videoData.Source).Scan(&videoId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(c, `
		insert into `+model.VideosTableName+"_"+model.GroupsTableName+`
		values($1, $2)
	`, videoId, videoData.GroupId)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func (vr *videoPgRepository) InsertMany(c context.Context, videoData []model.VideoCreate) error {
	tx, err := vr.db.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	for _, video := range videoData {
		innerTx, err := tx.Begin(c)
		if err != nil {
			return err
		}
		defer innerTx.Rollback(c)

		var videoId int
		err = innerTx.QueryRow(c, `
			insert into `+model.VideosTableName+`(title, source)
			values($1, $2)
			returning id
		`, video.Title, video.Source).Scan(&videoId)
		if err != nil {
			return err
		}

		_, err = innerTx.Exec(c, `
			insert into `+model.VideosTableName+"_"+model.GroupsTableName+`
			values($1, $2)
		`, videoId, video.GroupId)
		if err != nil {
			return err
		}

		err = innerTx.Commit(c)
		if err != nil {
			return err
		}
	}

	return tx.Commit(c)
}

func (vr *videoPgRepository) FindOne(c context.Context, filter string, value any, userGroupIds []int) (model.Video, error) {
	var video model.Video

	sql := `select * from ` + model.VideosTableName
	if filter == "status" {
		sql += ` where ` + filter + ` = '` + value.(string) + `'`
	} else if filter == "groupId" {
		sql += ` where id in (select videoId from ` + model.VideosTableName + "_" + model.GroupsTableName + ` where groupId = ` + value.(string) + `)`
	}

	sql += ` and id in (select videoId from ` + model.VideosTableName + "_" + model.GroupsTableName + ` where groupId = any($1))`

	err := vr.db.QueryRow(c, sql, userGroupIds).Scan(&video.Id, &video.Title, &video.Source, &video.ProcessedSource, &video.Status, &video.CreatedAt, &video.UpdatedAt)
	if err != nil {
		return video, err
	}

	return video, nil
}

func (vr *videoPgRepository) FindMany(c context.Context, filter string, value any, offset, limit int, userGroupIds []int) ([]model.Video, error) {
	var videos []model.Video

	sql := `select * from ` + model.VideosTableName
	if filter == "" {
		sql += ` where id in (select videoId from ` + model.VideosTableName + "_" + model.GroupsTableName + ` where groupId = any($1))`
	} else {
		if filter == "status" {
			sql += ` where ` + filter + ` = '` + value.(string) + `'`
		} else if filter == "groupId" {
			sql += ` where id in (select videoId from ` + model.VideosTableName + "_" + model.GroupsTableName + ` where groupId = ` + value.(string) + `)`
		}
		sql += ` 
			and id in (select videoId from ` + model.VideosTableName + "_" + model.GroupsTableName + ` where groupId = any($1))
		`
	}
	sql += `order by id desc offset $2 limit $3`

	rows, err := vr.db.Query(c, sql, userGroupIds, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var video model.Video
		if err := rows.Scan(&video.Id, &video.Title, &video.Source, &video.ProcessedSource, &video.Status, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (vr *videoPgRepository) DeleteOne(c context.Context, videoId int) error {
	tx, err := vr.db.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
		delete from `+model.VideosTableName+"_"+model.GroupsTableName+`
		where videoId = $1
	`, videoId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(c, `
		delete from `+model.VideosTableName+`
		where id = $1
	`, videoId)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func (vr *videoPgRepository) AddToGroup(c context.Context, videoId, groupId int) error {
	_, err := vr.db.Exec(c, `
		insert into `+model.VideosTableName+"_"+model.GroupsTableName+`
		values ($1, $2)
	`, videoId, groupId)
	return err
}

func (vr *videoPgRepository) RemoveFromGroup(c context.Context, videoId, groupId int) error {
	_, err := vr.db.Exec(c, `
		delete from `+model.VideosTableName+"_"+model.GroupsTableName+`
		where videoId = $1 and groupId = $2
	`, videoId, groupId)
	return err
}
