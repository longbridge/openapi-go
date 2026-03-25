package content

import "time"

// Author is the author of a topic
type Author struct {
	// Member ID
	MemberId string
	// Display name
	Name string
	// Avatar URL
	Avatar string
}

// Image is an image attached to a topic
type Image struct {
	// Original image URL
	Url string
	// Small thumbnail URL
	Sm string
	// Large image URL
	Lg string
}

// OwnedTopic is a topic created by the current user
type OwnedTopic struct {
	// Topic ID
	Id string
	// Title
	Title string
	// Plain text excerpt
	Description string
	// Markdown body
	Body string
	// Author
	Author *Author
	// Related stock tickers, format: {symbol}.{market}
	Tickers []string
	// Hashtag names
	Hashtags []string
	// Images
	Images []*Image
	// Likes count
	LikesCount int32
	// Comments count
	CommentsCount int32
	// Views count
	ViewsCount int32
	// Shares count
	SharesCount int32
	// Content type: "article" or "post"
	TopicType string
	// License: 0=none, 1=original, 2=non-original
	License int32
	// URL to the full topic page
	DetailUrl string
	// Created time
	CreatedAt time.Time
	// Updated time
	UpdatedAt time.Time
}

// TopicsMineOptions are the options for TopicsMine
type TopicsMineOptions struct {
	// Page number (default 1)
	Page int32
	// Number of records per page, range 1~500 (default 50)
	Size int32
	// Filter by content type: "article" or "post"; empty returns all
	TopicType string
}

// CreateTopicOptions are the options for TopicsCreate
type CreateTopicOptions struct {
	// Topic title (required)
	Title string
	// Topic body in Markdown format (required)
	Body string
	// Content type: "article" or "post" (default "post")
	TopicType string
	// Related stock tickers, format: {symbol}.{market}, max 10
	Tickers []string
	// Hashtag names, max 5
	Hashtags []string
	// License: 0=none (default), 1=original, 2=non-original
	License int32
}

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
