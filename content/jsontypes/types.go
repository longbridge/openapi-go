package jsontypes

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

type TopicAuthor struct {
	MemberId string `json:"member_id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}

type TopicImage struct {
	Url string `json:"url"`
	Sm  string `json:"sm"`
	Lg  string `json:"lg"`
}

type OwnedTopic struct {
	Id            string       `json:"id"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Body          string       `json:"body"`
	Author        *TopicAuthor `json:"author"`
	Tickers       []string     `json:"tickers"`
	Hashtags      []string     `json:"hashtags"`
	Images        []*TopicImage `json:"images"`
	LikesCount    int32        `json:"likes_count"`
	CommentsCount int32        `json:"comments_count"`
	ViewsCount    int32        `json:"views_count"`
	SharesCount   int32        `json:"shares_count"`
	TopicType     string       `json:"topic_type"`
	License       int32        `json:"license"`
	DetailUrl     string       `json:"detail_url"`
	CreatedAt     int64        `json:"created_at,string"`
	UpdatedAt     int64        `json:"updated_at,string"`
}

type OwnedTopicList struct {
	Items []*OwnedTopic `json:"items"`
}

type createTopicItem struct {
	Id string `json:"id"`
}

type CreateTopicResponse struct {
	Item createTopicItem `json:"item"`
}

type CreateTopicRequest struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	TopicType string   `json:"topic_type,omitempty"`
	Tickers   []string `json:"tickers,omitempty"`
	Hashtags  []string `json:"hashtags,omitempty"`
	License   int32    `json:"license,omitempty"`
}

