package sharelist

type SharelistList struct {
	Sharelists           []*SharelistInfo
	SubscribedSharelists []*SharelistInfo
	TailMark             string
}

type SharelistDetail struct {
	Sharelist *SharelistInfo
	Scopes    *SharelistScopes
}

type SharelistInfo struct {
	Id               int64
	Name             string
	Description      string
	Cover            string
	SubscribersCount int64
	CreatedAt        string
	EditedAt         string
	ThisYearChg      string
	Subscribed       bool
	Chg              string
	SharelistType    int32
	IndustryCode     string
	Stocks           []*SharelistStock
}

type SharelistStock struct {
	Symbol                  string
	Name                    string
	Market                  string
	Code                    string
	Intro                   string
	UnreadChangeLogCategory string
	Change                  string
	LastDone                string
	TradeStatus             int32
	Latency                 bool
}

type SharelistScopes struct {
	Subscription bool
	IsSelf       bool
}
