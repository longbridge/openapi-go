package content

// ListMyTopicsOptions holds optional parameters for listing topics created by the authenticated user.
type ListMyTopicsOptions struct {
	Page      int    // optional, default 1
	Size      int    // optional, default 50, range 1–500
	TopicType string // optional: "article" | "post"; empty returns all
}

// CreateTopicOptions holds the parameters for creating a new topic.
//
// Stock symbols mentioned in Body (e.g. "700.HK", "TSLA.US") are automatically
// recognized and linked by the platform. Use Tickers to associate additional
// symbols not mentioned in the body.
type CreateTopicOptions struct {
	Title     string   // required for topic_type "article"; optional for "post"
	Body      string   // required; plain text for "post", Markdown for "article"
	TopicType string   // optional: "post" (default) | "article"
	Tickers   []string // optional, max 10; format: "{symbol}.{market}"
	Hashtags  []string // optional, max 5
}

// ListTopicRepliesOptions holds optional pagination parameters for listing replies on a topic.
type ListTopicRepliesOptions struct {
	Page int // optional, default 1
	Size int // optional, default 20, range 1–50
}

// CreateReplyOptions holds the parameters for posting a reply to a topic.
//
// Body is plain text only — Markdown is not rendered.
// Stock symbols mentioned in Body (e.g. "700.HK", "TSLA.US") are automatically
// recognized and linked by the platform.
//
// Permission: the authenticated user must hold a funded Longbridge account.
// Rate limit: first 3 replies per topic have no interval requirement; subsequent
// replies require incrementally longer waits (3 s → 5 s → 8 s → 13 s → 21 s → 34 s → 55 s cap).
type CreateReplyOptions struct {
	Body      string // required; plain text only
	ReplyToID string // optional; ID of the reply to nest under; empty or "0" = top-level
}
