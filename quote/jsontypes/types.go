package jsontypes

type WatchedSecurity struct {
	Symbol    string `json:"symbol"`
	Market    string `json:"market"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	WatchedAt int64  `json:"watched_at,string"`
}

type WatchedGroup struct {
	Id        int64              `json:"id,string"`
	Name      string             `json:"name"`
	Securites []*WatchedSecurity `json:"securities"`
}

type WatchedGroupList struct {
	Groups []*WatchedGroup `json:"groups"`
}

type Security struct {
	Symbol string `json:"symbol"`
	NameCN string `json:"name_cn"`
	NameEN string `json:"name_en"`
	NameHK string `json:"name_hk"`
}

type SecurityList struct {
	List []*Security
}

type FilingItem struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	FileName    string   `json:"file_name"`
	FileUrls    []string `json:"file_urls"`
	PublishAt   int64    `json:"publish_at,string"`
}

type FilingList struct {
	Items []*FilingItem `json:"items"`
}
