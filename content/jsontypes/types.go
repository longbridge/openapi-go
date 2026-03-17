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
