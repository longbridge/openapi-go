package jsontypes

// --- Legacy types (topics/news by symbol) ---

type TopicItem struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Url           string `json:"url"`
	PublishedAt   int64  `json:"published_at,string"`
	CommentsCount int32  `json:"comments_count"`
	LikesCount    int32  `json:"likes_count"`
	SharesCount   int32  `json:"shares_count"`
}

type TopicList struct {
	Items []*TopicItem `json:"items"`
}

type NewsItem struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Url           string `json:"url"`
	PublishedAt   int64  `json:"published_at,string"`
	CommentsCount int32  `json:"comments_count"`
	LikesCount    int32  `json:"likes_count"`
	SharesCount   int32  `json:"shares_count"`
}

type NewsList struct {
	Items []*NewsItem `json:"items"`
}

// --- Community topic/reply types ---

// TopicAuthor is the author of a topic or reply.
type TopicAuthor struct {
	MemberID string `json:"member_id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}

// TopicImage is an image attached to a topic or reply.
type TopicImage struct {
	URL string `json:"url"`
	Sm  string `json:"sm"`
	Lg  string `json:"lg"`
}

// OwnedTopic is a topic created by the authenticated user, returned by topic detail or list-mine.
type OwnedTopic struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Body          string       `json:"body"`
	Author        TopicAuthor  `json:"author"`
	Tickers       []string     `json:"tickers"`
	Hashtags      []string     `json:"hashtags"`
	Images        []TopicImage `json:"images"`
	LikesCount    string       `json:"likes_count"`
	CommentsCount string       `json:"comments_count"`
	ViewsCount    string       `json:"views_count"`
	SharesCount   string       `json:"shares_count"`
	TopicType     string       `json:"topic_type"`
	DetailURL     string       `json:"detail_url"`
	CreatedAt     string       `json:"created_at"`
	UpdatedAt     string       `json:"updated_at"`
}

// TopicDetailResponse wraps a single OwnedTopic from the detail endpoint.
type TopicDetailResponse struct {
	Item OwnedTopic `json:"item"`
}

// TopicReply is a reply on a topic.
type TopicReply struct {
	ID            string       `json:"id"`
	TopicID       string       `json:"topic_id"`
	Body          string       `json:"body"`
	ReplyToID     string       `json:"reply_to_id"`
	Author        TopicAuthor  `json:"author"`
	Images        []TopicImage `json:"images"`
	LikesCount    string       `json:"likes_count"`
	CommentsCount string       `json:"comments_count"`
	CreatedAt     string       `json:"created_at"`
}

// TopicRepliesResponse wraps a list of replies.
type TopicRepliesResponse struct {
	Items []*TopicReply `json:"items"`
}

// TopicReplyResponse wraps a single reply (create response).
type TopicReplyResponse struct {
	Item TopicReply `json:"item"`
}
