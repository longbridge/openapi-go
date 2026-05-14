package jsontypes

type WatchedSecurity struct {
	Symbol    string `json:"symbol"`
	Market    string `json:"market"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	WatchedAt int64  `json:"watched_at,string"`
	IsPinned  bool   `json:"is_pinned"`
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

type ShortPosition struct {
	Symbol             string `json:"symbol"`
	Date               string `json:"date"`
	ShortSellQty       int64  `json:"short_sell_qty,string"`
	MarketTurnover     string `json:"market_turnover"`
	ShortSellTurnover  string `json:"short_sell_turnover"`
	ShortSellRatio     string `json:"short_sell_ratio"`
}

type ShortPositionsResponse struct {
	List []*ShortPosition `json:"list"`
}

type OptionVolumeStats struct {
	Symbol      string `json:"symbol"`
	CallVolume  int64  `json:"call_volume,string"`
	PutVolume   int64  `json:"put_volume,string"`
	CallPutRatio string `json:"call_put_ratio"`
}

type OptionVolumeResponse struct {
	List []*OptionVolumeStats `json:"list"`
}

type DailyOptionVolume struct {
	Date         string `json:"date"`
	CallVolume   int64  `json:"call_volume,string"`
	PutVolume    int64  `json:"put_volume,string"`
	CallPutRatio string `json:"call_put_ratio"`
}

type OptionVolumeDailyResponse struct {
	Symbol string               `json:"symbol"`
	Items  []*DailyOptionVolume `json:"items"`
}
