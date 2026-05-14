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

// ShortPosition is a single short interest data point
type ShortPosition struct {
	Timestamp           string `json:"timestamp"`
	Rate                string `json:"rate"`
	AvgDailyShareVolume string `json:"avg_daily_share_volume"`
	CurrentSharesShort  string `json:"current_shares_short"`
	DaysToCover         string `json:"days_to_cover"`
	Close               string `json:"close"`
}

// ShortPositionStats contains short interest data for a security
type ShortPositionStats struct {
	Symbol  string           `json:"counter_id"`
	Data    []*ShortPosition `json:"data"`
	Sources int32            `json:"sources"`
}

// OptionVolumeStats contains total call/put volume for a security
type OptionVolumeStats struct {
	CallVolume string `json:"c"`
	PutVolume  string `json:"p"`
}

// OptionVolumeDailyStat is a single daily option volume data point
type OptionVolumeDailyStat struct {
	Symbol                   string `json:"underlying_counter_id"`
	Timestamp                string `json:"timestamp"`
	TotalVolume              string `json:"total_volume"`
	TotalPutVolume           string `json:"total_put_volume"`
	TotalCallVolume          string `json:"total_call_volume"`
	PutCallVolumeRatio       string `json:"put_call_volume_ratio"`
	TotalOpenInterest        string `json:"total_open_interest"`
	TotalPutOpenInterest     string `json:"total_put_open_interest"`
	TotalCallOpenInterest    string `json:"total_call_open_interest"`
	PutCallOpenInterestRatio string `json:"put_call_open_interest_ratio"`
}

// OptionVolumeDaily contains a list of daily option volume stats
type OptionVolumeDaily struct {
	Stats []*OptionVolumeDailyStat `json:"stats"`
}
