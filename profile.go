package main

import (
	"net/http"
	"net/mail"
	"time"
)

type ProfileId string

type VoteId int

type Tag string

type Domain string

type Void struct{}

type ProfileVote struct {
	Time              time.Time // in sql, "voteid" as UNIX seconds
	PublicationId     PublicationId
	PublicationPoster ProfileId
	PublicationDomain Domain
	Up                bool
}

type ProfileSuggestedTag struct {
	PublicationId PublicationId
	Tag           Tag
}

type Donation struct {
	Cents int
	Time  time.Time
}

type Profile struct {
	Id                    ProfileId
	Joined                time.Time
	Inactive              bool
	Clout                 float64
	WordCount             bool
	AverageComNumerator   int
	AverageComDenominator int
	IgnoreNewbies         bool
	GlobalIgnored         bool
	NewTabs               bool
	PublicationTabs       bool
	ReplyAlerts           bool
	PmAlerts              bool
	FollowerAlerts        bool
	Shoutouts             bool
	BadgeAlerts           bool
	SavedNotifications    bool
	FeedTimes             bool
	ShareCounts           bool
	ShowGlobalFiltered    bool
	FollowsBadges         bool
	EmbedVideos           bool
	Bio                   string
	Email                 bool
	HubskiStyle           string // TODO(enum)
	Homepage              string // TODO(enum)
	FollowerCount         int
	PostsCount            int
	SharedsCount          int
	UnreadNotifications   bool
	LastComTime           time.Time
	ComCloutDate          time.Time
	Zen                   bool
	Spammer               bool
	Submitted             map[PublicationId]Void
	Saved                 map[PublicationId]Void
	Sticky                map[PublicationId]Void
	Hidden                map[PublicationId]Void
	Mail                  map[PublicationId]Void
	Drafts                map[PublicationId]Void
	Shareds               map[PublicationId]Void
	Cubbeds               map[PublicationId]Void // "Cubbeds" are shared comments
	Votes                 map[ProfileVote]Void
	SuggestedTags         map[ProfileSuggestedTag]Void
	Badged                map[PublicationId]Void
	Badges                map[PublicationId]Void
	Ignoring              map[ProfileId]Void
	Muting                map[ProfileId]Void
	Hushing               map[ProfileId]Void
	Blocking              map[ProfileId]Void
	IgnoringTag           map[Tag]Void
	IgnoringDomain        map[Domain]Void
	IgnoredBy             map[ProfileId]Void
	MutedBy               map[ProfileId]Void
	HushedBy              map[ProfileId]Void
	BlockedBy             map[ProfileId]Void
	Followed              map[ProfileId]Void
	Follower              map[ProfileId]Void
	PersonalTags          map[Tag]Void
	FollowedTags          map[Tag]Void
	FollowedDomain        map[Domain]Void
	Notified              map[PublicationId]Void
	Donations             map[Donation]Void
}

type ProfilePasswordHash struct {
	ProfileId ProfileId
	HashType  string // TODO(enum)
	Hash1     string
	Hash2     string
}

type ProfileCookies struct {
	ProfileId ProfileId
	Cookie    http.Cookie
}

type ProfileEmails struct {
	ProfileId ProfileId
	Email     mail.Address
}
