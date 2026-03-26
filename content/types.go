package content

import "time"

// --- Legacy types (topics/news by symbol) ---

// TopicItem is a discussion topic for a security.
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

// NewsItem is a news article for a security.
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

// --- Community topic/reply types ---

// TopicAuthor is the author of a topic or reply.
type TopicAuthor struct {
	MemberID string
	Name     string
	Avatar   string
}

// TopicImage is an image attached to a topic or reply.
type TopicImage struct {
	URL string
	Sm  string
	Lg  string
}

// OwnedTopic is a full topic record returned by the topic-detail and list-mine endpoints.
type OwnedTopic struct {
	ID            string
	Title         string
	Description   string
	Body          string
	Author        TopicAuthor
	Tickers       []string
	Hashtags      []string
	Images        []TopicImage
	LikesCount    int64
	CommentsCount int64
	ViewsCount    int64
	SharesCount   int64
	TopicType     string
	DetailURL     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TopicReply is a reply on a topic.
type TopicReply struct {
	ID            string
	TopicID       string
	Body          string
	ReplyToID     string
	Author        TopicAuthor
	Images        []TopicImage
	LikesCount    int64
	CommentsCount int64
	CreatedAt     time.Time
}
