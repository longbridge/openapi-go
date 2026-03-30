package content

import (
	"net/url"
	"strconv"
)

func (o *MyTopicsOptions) values() url.Values {
	p := url.Values{}
	if o == nil {
		return p
	}
	if o.Page > 0 {
		p.Set("page", strconv.Itoa(o.Page))
	}
	if o.Size > 0 {
		p.Set("size", strconv.Itoa(o.Size))
	}
	if o.TopicType != "" {
		p.Set("topic_type", o.TopicType)
	}
	return p
}

func (o *ListTopicRepliesOptions) values() url.Values {
	p := url.Values{}
	if o == nil {
		return p
	}
	if o.Page > 0 {
		p.Set("page", strconv.Itoa(o.Page))
	}
	if o.Size > 0 {
		p.Set("size", strconv.Itoa(o.Size))
	}
	return p
}
