package main

import (
	"net/url"
	"time"
)

type PublicationId int

type PublicationType int

const (
	PublicationTypeStory PublicationType = iota
	PublicationTypeComment
)

type PublicationCommunityTag struct {
	Username ProfileId
	Tag      Tag
}

type PublicationVote struct {
	Time          time.Time // in sql, "voteid" as UNIX seconds
	PublicationId PublicationId
	ProfileId     ProfileId `sql:"username"`
	Up            bool
	Score         int `sql:"num"`
}

type Publication struct {
	PublicationId   PublicationId
	TypeId          PublicationType
	Username        ProfileId
	Time            time.Time
	Url             url.URL
	Title           string
	Mail            bool
	Tag             Tag
	Tag2            Tag
	Text            string
	Md              string
	Domain          string
	Score           int // TODO(calculate)
	Deleted         bool
	Draft           bool
	Parent          PublicationId
	Locked          bool
	NoKill          bool
	CC              map[ProfileId]Void
	CommunityTags   map[Tag]Void
	CommunityTagses PublicationCommunityTag
	Votes           map[PublicationVote]Void
	SavedBy         map[ProfileId]Void
	SharedBy        map[ProfileId]Void
	BadgedBy        map[ProfileId]Void
	BadgedKids      map[PublicationId]Void
	CubbedBy        map[ProfileId]Void
	Kids            map[PublicationId]Void
}
