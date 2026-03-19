package content

import "time"

// TopicItem is a discussion topic for a security
type TopicItem struct {
	// Topic ID
	Id string
	// Title
	Title string
	// Description
	Description string
	// URL
	Url string
	// Published time
	PublishedAt time.Time
	// Comments count
	CommentsCount int32
	// Likes count
	LikesCount int32
	// Shares count
	SharesCount int32
}

// NewsItem is a news article for a security
type NewsItem struct {
	// News ID
	Id string
	// Title
	Title string
	// Description
	Description string
	// URL
	Url string
	// Published time
	PublishedAt time.Time
	// Comments count
	CommentsCount int32
	// Likes count
	LikesCount int32
	// Shares count
	SharesCount int32
}
