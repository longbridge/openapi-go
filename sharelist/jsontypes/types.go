package jsontypes

type SharelistList struct {
	Sharelists           []*SharelistInfo `json:"sharelists"`
	SubscribedSharelists []*SharelistInfo `json:"subscribed_sharelists"`
	TailMark             string           `json:"tail_mark"`
}

type SharelistDetail struct {
	Sharelist *SharelistInfo  `json:"sharelist"`
	Scopes    *SharelistScopes `json:"scopes"`
}

type SharelistInfo struct {
	Id               int64             `json:"id,string"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Cover            string            `json:"cover"`
	SubscribersCount int64             `json:"subscribers_count"`
	CreatedAt        string            `json:"created_at"`
	EditedAt         string            `json:"edited_at"`
	ThisYearChg      string            `json:"this_year_chg"`
	Subscribed       bool              `json:"subscribed"`
	Chg              string            `json:"chg"`
	SharelistType    int32             `json:"sharelist_type"`
	IndustryCode     string            `json:"industry_code"`
	Stocks           []*SharelistStock `json:"stocks"`
}

type SharelistStock struct {
	CounterId                string `json:"counter_id"`
	Name                     string `json:"name"`
	Market                   string `json:"market"`
	Code                     string `json:"code"`
	Intro                    string `json:"intro"`
	UnreadChangeLogCategory  string `json:"unread_change_log_category"`
	Change                   string `json:"change"`
	LastDone                 string `json:"last_done"`
	TradeStatus              int32  `json:"trade_status"`
	Latency                  bool   `json:"latency"`
}

type SharelistScopes struct {
	Subscription bool `json:"subscription"`
	IsSelf       bool `json:"self"`
}
