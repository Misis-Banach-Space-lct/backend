package model

import (
	"context"
	"time"
)

var VideosTableName = "videos"

type Video struct {
	Id          int       `json:"id"`    // serial
	Title       string    `json:"title"` // unique
	ContentType string    `json:"contentType"`
	Url         string    `json:"-"`
	CreatedAt   time.Time `json:"createdAt"` // default = current timestamp
	UpdatedAt   time.Time `json:"updatedAt"` // default = current timestamp
}

type VideoCreate struct {
	Title       string `json:"title" validate:"required"`
	ContentType string `json:"contentType" validate:"required,oneof=camera file"`
	Url         string `json:"url" validate:"required"`
}

type VideoGroupUpdate struct {
	Action  string `json:"action" validate:"required,oneof=add remove"`
	VideoId int    `json:"videoId" validate:"required,gt=0"`
	GroupId int    `json:"groupId" validate:"required,gte=0"`
}

type VideoRepository interface {
	InsertOne(c context.Context, videoData VideoCreate) error
	InsertMany(c context.Context, videosData []VideoCreate) error
	FindOne(c context.Context, filter string, value any) (Video, error)
	FindMany(c context.Context, contentType string) ([]Video, error)
	DeleteOne(c context.Context, videoId int) error // deletes video source and all related frames
	AddToGroup(c context.Context, videoId, groupId int) error
	RemoveFromGroup(c context.Context, videoId, groupId int) error
}
