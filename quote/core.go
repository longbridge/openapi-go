package quote

import (
	"context"
	"sync"
	"time"

	quotev1 "github.com/longbridge/openapi-protobufs/gen/go/quote"
	protocol "github.com/longbridge/openapi-protocol/go"
	"github.com/longbridge/openapi-protocol/go/client"
	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go"
	"github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/log"
)

type core struct {
	opts            *Options
	transport       Transport
	enableOvernight bool
	language        string
	httpClient      *http.Client

	// WebSocket client; nil until dialed (see ws). Guarded by wsMu.
	client          client.Client
	wsMu            sync.Mutex
	closed          bool
	pendingHandlers []func(client.Client)

	url           string
	mu            sync.Mutex
	subscriptions map[string][]SubType
	store         *store
	profileMu     sync.Mutex
	userProfile   *UserProfile
}

func newCore(opts *Options) (*core, error) {
	if opts.httpClient == nil {
		return nil, errors.New("quote: http client is required (use WithHttpClient)")
	}
	core := &core{
		opts:            opts,
		transport:       opts.transport,
		enableOvernight: opts.enableOvernight,
		language:        string(opts.language),
		httpClient:      opts.httpClient,
		url:             opts.quoteURL,
		subscriptions:   make(map[string][]SubType),
		store:           newStore(),
	}
	if core.transport == TransportHTTP {
		// Pull APIs go over HTTP; the WebSocket is dialed lazily by ws() and
		// the profile fetched by Profile().
		return core, nil
	}
	if _, err := core.ws(context.Background()); err != nil {
		return nil, err
	}
	var err error
	core.userProfile, err = core.queryProfile(context.Background(), core.language)
	if err != nil {
		_ = core.Close()
		return nil, err
	}
	return core, nil
}

// dial connects and authenticates the quote WebSocket and queries the user
// quote profile. Callers hold wsMu.
func (c *core) dial(ctx context.Context) (client.Client, error) {
	opts := c.opts
	getOTP := func() (otp string, err error) {
		otp, err = opts.httpClient.GetOTPV2(context.Background())
		if err != nil {
			return "", errors.Wrap(err, "failed to get otp")
		}
		return otp, nil
	}

	logger := opts.logger
	logger.SetLevel(opts.logLevel)

	clientOpts := []client.ClientOption{client.WithLogger(logger)}
	if opts.enableOvernight {
		clientOpts = append(clientOpts, client.WithConnectMetadata(map[string]string{"need_over_night_quote": "true"}))
	}
	cl := client.New(clientOpts...)

	err := cl.Dial(ctx, opts.quoteURL,
		&protocol.Handshake{
			Version:  1,
			Codec:    protocol.CodecProtobuf,
			Platform: protocol.PlatformOpenapi,
		},
		client.WithAuthTokenGetter(getOTP),
		client.AuthTimeout(opts.lbOpts.AuthTimeout),
		client.DialTimeout(opts.lbOpts.Timeout),
		client.MinGzipSize(opts.lbOpts.MinGzipSize),
		client.ReadBufferSize(opts.lbOpts.ReadBufferSize),
		client.ReadQueueSize(opts.lbOpts.ReadQueueSize),
		client.WriteQueueSize(opts.lbOpts.WriteQueueSize),
	)
	if err != nil {
		// Dial may have connected before authentication failed; release it.
		_ = cl.Close(err)
		return nil, err
	}

	cl.AfterReconnected(func() {
		c.wsMu.Lock()
		closed := c.closed
		c.wsMu.Unlock()
		if closed {
			return
		}
		resubFlag := true

		if err := c.resubscribe(context.Background()); err != nil {
			log.Errorf("faield to do sub, err: %v", err)
			resubFlag = false
		}
		for _, fn := range opts.reconnectCallbacks {
			fn(resubFlag)
		}
	})
	return cl, nil
}

func (c *core) SetQuoteHandler(f func(*PushQuote)) {
	c.withClient(func(cl client.Client) { cl.Subscribe(uint32(quotev1.Command_PushQuoteData), parsePushQuoteFunc(f, c)) })
}

func (c *core) SetTradeHandler(f func(*PushTrade)) {
	c.withClient(func(cl client.Client) { cl.Subscribe(uint32(quotev1.Command_PushTradeData), parsePushTradeFunc(f, c)) })
}

func (c *core) SetDepthHandler(f func(*PushDepth)) {
	c.withClient(func(cl client.Client) { cl.Subscribe(uint32(quotev1.Command_PushDepthData), parsePushDepthFunc(f, c)) })
}

func (c *core) SetBrokersHandler(f func(*PushBrokers)) {
	c.withClient(func(cl client.Client) {
		cl.Subscribe(uint32(quotev1.Command_PushBrokersData), parsePushBrokersFunc(f, c))
	})
}

func (c *core) Subscribe(ctx context.Context, symbols []string, subTypes []SubType, isFirstPush bool) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.doSubscirbe(ctx, symbols, subTypes, isFirstPush)
}

func (c *core) doSubscirbe(ctx context.Context, symbols []string, subTypes []SubType, isFirstPush bool) (err error) {
	req := &quotev1.SubscribeRequest{
		IsFirstPush: isFirstPush,
		Symbol:      symbols,
	}
	err = util.Copy(&req.SubType, subTypes)
	if err != nil {
		return
	}
	err = c.wsDo(ctx, quotev1.Command_Subscribe, req, nil)
	if err != nil {
		return
	}
	for _, symbol := range symbols {
		c.subscriptions[symbol] = subTypes
	}
	return
}

func (c *core) Unsubscribe(ctx context.Context, unSubAll bool, symbols []string, subTypes []SubType) (err error) {
	if !c.connected() {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	req := &quotev1.UnsubscribeRequest{
		Symbol:   symbols,
		UnsubAll: unSubAll,
	}
	err = util.Copy(&req.SubType, subTypes)
	if err != nil {
		return
	}
	err = c.wsDo(ctx, quotev1.Command_Unsubscribe, req, nil)
	if err != nil {
		return
	}
	if unSubAll {
		c.subscriptions = make(map[string][]SubType)
	}
	for _, symbol := range symbols {
		delete(c.subscriptions, symbol)
	}
	return
}

func (c *core) resubscribe(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for symbol, subflags := range c.subscriptions {
		err := c.doSubscirbe(ctx, []string{symbol}, subflags, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *core) Profile() *UserProfile {
	c.profileMu.Lock()
	defer c.profileMu.Unlock()
	if c.userProfile == nil {
		// With TransportHTTP the profile is only fetched on demand; this
		// dials the WebSocket if it is not connected yet.
		profile, err := c.queryProfile(context.Background(), c.language)
		if err != nil {
			log.Errorf("failed to query quote profile: %v", err)
			return nil
		}
		c.userProfile = profile
	}
	return c.userProfile
}

func (c *core) queryProfile(ctx context.Context, lang string) (quoteProfile *UserProfile, err error) {
	req := &quotev1.UserQuoteProfileRequest{
		Language: lang,
	}
	var ret quotev1.UserQuoteProfileResponse
	if err = c.do(ctx, quotev1.Command_QueryUserQuoteProfile, req, &ret); err != nil {
		return
	}
	quoteProfile = &UserProfile{}
	// nolint:govet
	err = util.Copy(&quoteProfile, ret)
	return
}

func (c *core) Subscriptions(ctx context.Context) (subscriptions map[string][]SubType, err error) {
	if !c.connected() {
		// Nothing can be subscribed without a connection; don't dial just to
		// report that.
		return map[string][]SubType{}, nil
	}
	req := &quotev1.SubscriptionRequest{}
	var ret quotev1.SubscriptionResponse
	if err = c.do(ctx, quotev1.Command_Subscription, req, &ret); err != nil {
		return
	}
	subscriptions = make(map[string][]SubType, len(ret.GetSubList()))
	for _, item := range ret.GetSubList() {
		sublist := make([]SubType, 0, len(item.GetSubType()))
		err = util.Copy(&sublist, item.GetSubType())
		if err != nil {
			return
		}
		subscriptions[item.GetSymbol()] = sublist
	}
	return
}

func (c *core) StaticInfo(ctx context.Context, symbols []string) (staticInfos []*StaticInfo, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol: symbols,
	}
	var ret quotev1.SecurityStaticInfoResponse
	if err = c.do(ctx, quotev1.Command_QuerySecurityStaticInfo, req, &ret); err != nil {
		return
	}
	err = util.Copy(&staticInfos, ret.GetSecuStaticInfo())
	return
}

func (c *core) Quote(ctx context.Context, symbols []string) (quotes []*SecurityQuote, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol: symbols,
	}
	var ret quotev1.SecurityQuoteResponse
	if err = c.do(ctx, quotev1.Command_QuerySecurityQuote, req, &ret); err != nil {
		return
	}
	if c.alignsOvernight() {
		// See overnight.go: HTTP always returns overnight quotes.
		for _, q := range ret.GetSecuQuote() {
			q.OverNightQuote = nil
		}
	}
	err = util.Copy(&quotes, ret.GetSecuQuote())
	return
}

func (c *core) OptionQuote(ctx context.Context, symbols []string) (optionQuotes []*OptionQuote, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol: symbols,
	}
	var ret quotev1.OptionQuoteResponse
	if err = c.do(ctx, quotev1.Command_QueryOptionQuote, req, &ret); err != nil {
		return
	}
	err = util.Copy(&optionQuotes, ret.GetSecuQuote())
	return
}

func (c *core) WarrantQuote(ctx context.Context, symbols []string) (warrantQuotes []*WarrantQuote, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol: symbols,
	}
	var ret quotev1.WarrantQuoteResponse
	if err = c.do(ctx, quotev1.Command_QueryWarrantQuote, req, &ret); err != nil {
		return
	}
	err = util.Copy(&warrantQuotes, ret.GetSecuQuote())
	return
}

func (c *core) Depth(ctx context.Context, symbol string) (securityDepth *SecurityDepth, err error) {
	req := &quotev1.SecurityRequest{
		Symbol: symbol,
	}
	var ret quotev1.SecurityDepthResponse
	if err = c.do(ctx, quotev1.Command_QueryDepth, req, &ret); err != nil {
		return
	}
	securityDepth = &SecurityDepth{}
	err = util.Copy(&securityDepth, ret)
	return
}

func (c *core) Brokers(ctx context.Context, symbol string) (securityBorkers *SecurityBrokers, err error) {
	req := &quotev1.SecurityRequest{
		Symbol: symbol,
	}
	var ret quotev1.SecurityBrokersResponse
	if err = c.do(ctx, quotev1.Command_QueryBrokers, req, &ret); err != nil {
		return
	}
	securityBorkers = &SecurityBrokers{}
	err = util.Copy(securityBorkers, ret)
	return
}

func (c *core) Participants(ctx context.Context) (infos []*ParticipantInfo, err error) {
	var ret quotev1.ParticipantBrokerIdsResponse
	if err = c.do(ctx, quotev1.Command_QueryParticipantBrokerIds, nil, &ret); err != nil {
		return
	}
	err = util.Copy(&infos, ret.GetParticipantBrokerNumbers())
	return
}

func (c *core) Trades(ctx context.Context, symbol string, count int32) (trades []*Trade, err error) {
	req := &quotev1.SecurityTradeRequest{
		Symbol: symbol,
		Count:  count,
	}
	var ret quotev1.SecurityTradeResponse
	if err = c.do(ctx, quotev1.Command_QueryTrade, req, &ret); err != nil {
		return
	}
	err = util.Copy(&trades, ret.GetTrades())
	return
}

func (c *core) Intraday(ctx context.Context, symbol string) (lines []*IntradayLine, err error) {
	req := &quotev1.SecurityIntradayRequest{
		Symbol: symbol,
	}
	var ret quotev1.SecurityIntradayResponse
	if err = c.do(ctx, quotev1.Command_QueryIntraday, req, &ret); err != nil {
		return
	}
	raw := ret.GetLines()
	if c.alignsOvernight() && isUSEquitySymbol(symbol) {
		// See overnight.go: HTTP always returns the overnight session.
		kept := make([]*quotev1.Line, 0, len(raw))
		for _, line := range raw {
			if !isUSOvernight(line.GetTimestamp()) {
				kept = append(kept, line)
			}
		}
		raw = kept
	}
	err = util.Copy(&lines, raw)
	return
}

func (c *core) Candlesticks(ctx context.Context, symbol string, period Period, count int32, adjustType AdjustType) (sticks []*Candlestick, err error) {
	req := &quotev1.SecurityCandlestickRequest{
		Symbol:     symbol,
		Period:     quotev1.Period(period),
		Count:      count,
		AdjustType: quotev1.AdjustType(adjustType),
	}
	var ret quotev1.SecurityCandlestickResponse
	if err = c.do(ctx, quotev1.Command_QueryCandlestick, req, &ret); err != nil {
		return
	}
	err = util.Copy(&sticks, ret.GetCandlesticks())
	return
}

func (c *core) HistoryCandlesticksByOffset(ctx context.Context, symbol string, period Period, adjustType AdjustType, isForward bool, dateTime *time.Time, count int32, options ...CandlestickRequestOption) (sticks []*Candlestick, err error) {
	direction := quotev1.Direction_BACKWARD
	if isForward {
		direction = quotev1.Direction_FORWARD
	}
	req := &quotev1.SecurityHistoryCandlestickRequest{
		Symbol:     symbol,
		Period:     quotev1.Period(period),
		AdjustType: quotev1.AdjustType(adjustType),
		QueryType:  quotev1.HistoryCandlestickQueryType_QUERY_BY_OFFSET,
		OffsetRequest: &quotev1.SecurityHistoryCandlestickRequest_OffsetQuery{
			Direction: direction,
			Date:      util.FormatDateSimple(dateTime),
			Minute:    util.FormatMinuteSimple(dateTime),
			Count:     int32(count),
		},
	}

	for _, option := range options {
		option(req)
	}

	return c.historyCandlesticks(ctx, req)
}

func (c *core) HistoryCandlesticksByDate(ctx context.Context, symbol string, period Period, adjustType AdjustType, startDate *time.Time, endDate *time.Time, options ...CandlestickRequestOption) (sticks []*Candlestick, err error) {
	req := &quotev1.SecurityHistoryCandlestickRequest{
		Symbol:     symbol,
		Period:     quotev1.Period(period),
		AdjustType: quotev1.AdjustType(adjustType),
		QueryType:  quotev1.HistoryCandlestickQueryType_QUERY_BY_DATE,
		DateRequest: &quotev1.SecurityHistoryCandlestickRequest_DateQuery{
			StartDate: util.FormatDateSimple(startDate),
			EndDate:   util.FormatDateSimple(endDate),
		},
	}

	for _, option := range options {
		option(req)
	}

	return c.historyCandlesticks(ctx, req)
}

func (c *core) historyCandlesticks(ctx context.Context, req *quotev1.SecurityHistoryCandlestickRequest) (sticks []*Candlestick, err error) {
	var ret quotev1.SecurityCandlestickResponse
	if err = c.do(ctx, quotev1.Command_QueryHistoryCandlestick, req, &ret); err != nil {
		return
	}
	candlesticks, err := c.alignHistoryCandlesticks(ctx, req, ret.GetCandlesticks())
	if err != nil {
		return
	}
	err = util.Copy(&sticks, candlesticks)
	return
}

// alignHistoryCandlesticks applies the WebSocket's US overnight rule to a raw
// HTTP candlestick page (see overnight.go): overnight bars are dropped and a
// count / range-capped window is topped up so it covers the same bars the
// WebSocket would return. The top-up pages with offset queries through
// c.do; a failure there fails the call.
func (c *core) alignHistoryCandlesticks(ctx context.Context, req *quotev1.SecurityHistoryCandlestickRequest, raw []*quotev1.Candlestick) ([]*quotev1.Candlestick, error) {
	if !c.alignsOvernight() {
		return raw, nil
	}
	return alignCandlesticks(ctx, req, raw, func(ctx context.Context, page *quotev1.SecurityHistoryCandlestickRequest) ([]*quotev1.Candlestick, error) {
		var ret quotev1.SecurityCandlestickResponse
		if err := c.do(ctx, quotev1.Command_QueryHistoryCandlestick, page, &ret); err != nil {
			return nil, err
		}
		return ret.GetCandlesticks(), nil
	})
}

func (c *core) OptionChainExpiryDateList(ctx context.Context, symbol string) (times []time.Time, err error) {
	req := &quotev1.SecurityRequest{
		Symbol: symbol,
	}
	var ret quotev1.OptionChainDateListResponse
	if err = c.do(ctx, quotev1.Command_QueryOptionChainDate, req, &ret); err != nil {
		return
	}
	times, err = toTimes(ret.GetExpiryDate())
	if err != nil {
		return
	}
	return
}

func (c *core) WarrantIssuers(ctx context.Context) (infos []*IssuerInfo, err error) {
	var ret quotev1.IssuerInfoResponse
	if err = c.do(ctx, quotev1.Command_QueryWarrantIssuerInfo, nil, &ret); err != nil {
		return
	}
	err = util.Copy(&infos, ret.GetIssuerInfo())
	return
}

func (c *core) WarrantList(ctx context.Context, symbol string, config WarrantFilter, lang WarrantLanguage) (infos []*WarrantInfo, err error) {

	filter := &quotev1.FilterConfig{}
	util.Copy(filter, &config)

	req := &quotev1.WarrantFilterListRequest{
		Symbol:       symbol,
		FilterConfig: filter,
		Language:     int32(lang),
	}

	var ret quotev1.WarrantFilterListResponse
	if err = c.do(ctx, quotev1.Command_QueryWarrantFilterList, req, &ret); err != nil {
		return
	}
	err = util.Copy(&infos, ret.GetWarrantList())
	return
}

func (c *core) TradingSession(ctx context.Context) (sessions []*MarketTradingSession, err error) {
	var ret quotev1.MarketTradePeriodResponse
	if err = c.do(ctx, quotev1.Command_QueryMarketTradePeriod, nil, &ret); err != nil {
		return
	}
	err = util.Copy(&sessions, ret.GetMarketTradeSession())
	return
}

func (c *core) TradingDays(ctx context.Context, market openapi.Market, begin *time.Time, end *time.Time) (days *MarketTradingDay, err error) {
	var (
		tradingDays   []time.Time
		halfTradeDays []time.Time
	)

	req := &quotev1.MarketTradeDayRequest{
		Market: string(market),
		BegDay: util.FormatDateSimple(begin),
		EndDay: util.FormatDateSimple(end),
	}
	var ret quotev1.MarketTradeDayResponse
	if err = c.do(ctx, quotev1.Command_QueryMarketTradeDay, req, &ret); err != nil {
		return
	}
	tradingDays, err = toTimes(ret.GetTradeDay())
	if err != nil {
		return
	}
	halfTradeDays, err = toTimes(ret.GetHalfTradeDay())
	if err != nil {
		return
	}
	days = &MarketTradingDay{
		TradeDay:     tradingDays,
		HalfTradeDay: halfTradeDays,
	}
	return
}

func (c *core) CapitalDistribution(ctx context.Context, symbol string) (capitalDib CapitalDistribution, err error) {
	req := &quotev1.SecurityRequest{
		Symbol: symbol,
	}
	var ret quotev1.CapitalDistributionResponse
	if err = c.do(ctx, quotev1.Command_QueryCapitalFlowDistribution, req, &ret); err != nil {
		return
	}
	err = util.Copy(&capitalDib, ret)
	return
}

func (c *core) CapitalFlow(ctx context.Context, symbol string) (capitalFlowLines []CapitalFlowLine, err error) {
	req := &quotev1.CapitalFlowIntradayRequest{
		Symbol: symbol,
	}
	var ret quotev1.CapitalFlowIntradayResponse
	if err = c.do(ctx, quotev1.Command_QueryCapitalFlowIntraday, req, &ret); err != nil {
		return
	}
	err = util.Copy(&capitalFlowLines, ret.GetCapitalFlowLines())
	return
}

func (c *core) CalcIndex(ctx context.Context, symbols []string, indexes []CalcIndex) (calcIndexes []*SecurityCalcIndex, err error) {
	quoteCalcIndexes := make([]quotev1.CalcIndex, 0, len(indexes))
	util.Copy(&quoteCalcIndexes, indexes)
	req := &quotev1.SecurityCalcQuoteRequest{
		Symbols:   symbols,
		CalcIndex: quoteCalcIndexes,
	}
	var ret quotev1.SecurityCalcQuoteResponse
	if err = c.do(ctx, quotev1.Command_QuerySecurityCalcIndex, req, &ret); err != nil {
		return
	}
	err = util.Copy(&calcIndexes, ret.GetSecurityCalcIndex())
	for _, calcIndex := range calcIndexes {
		doRatio(calcIndex)
	}
	return
}

func (c *core) RealtimeQuote(ctx context.Context, symbols []string) (quotes []*Quote, err error) {
	quotes = make([]*Quote, 0, len(symbols))
	for _, symbol := range symbols {
		if q := c.store.GetQuote(symbol); q != nil {
			quotes = append(quotes, q)
		}
	}
	return
}

func (c *core) RealtimeDepth(ctx context.Context, symbol string) (securityDepth *SecurityDepth, err error) {
	askDepths, bidDepths := c.store.GetDepth(symbol)
	return &SecurityDepth{
		Symbol: symbol,
		Ask:    askDepths,
		Bid:    bidDepths,
	}, nil
}

func (c *core) RealtimeTrades(ctx context.Context, symbol string) (trades []*Trade, err error) {
	return c.store.GetTrades(symbol), nil
}

func (c *core) RealtimeBrokers(ctx context.Context, symbol string) (securityBorkers *SecurityBrokers, err error) {
	askBrokers, bidBorkers := c.store.GetBrokers(symbol)
	return &SecurityBrokers{
		Symbol:     symbol,
		AskBrokers: askBrokers,
		BidBrokers: bidBorkers,
	}, nil
}

func (c *core) Close() error {
	c.wsMu.Lock()
	defer c.wsMu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	c.pendingHandlers = nil
	if c.client == nil {
		return nil
	}
	return c.client.Close(nil)
}

func parsePushQuoteFunc(f func(*PushQuote), core *core) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var err error
		var data quotev1.PushQuote
		if err = packet.Unmarshal(&data); err != nil {
			log.Errorf("quote push event, unmarshal error:%v", err)
			return
		}
		var pb PushQuote
		if err = util.Copy(&pb, data); err != nil {
			log.Errorf("quote push event, copy data error:%v", err)
			return
		}
		core.store.MergeQuote(&pb)
		f(&pb)
	}
}

func parsePushDepthFunc(f func(*PushDepth), core *core) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var err error
		var data quotev1.PushDepth
		if err = packet.Unmarshal(&data); err != nil {
			log.Errorf("quote depth push event, unmarshal error:%v", err)
			return
		}
		var pd PushDepth
		if err = util.Copy(&pd, data); err != nil {
			log.Errorf("quote depth push event, copy data error:%v", err)
			return
		}
		core.store.MergeDepth(&pd)
		f(&pd)
	}
}

func parsePushBrokersFunc(f func(*PushBrokers), core *core) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var err error
		var data quotev1.PushBrokers
		if err = packet.Unmarshal(&data); err != nil {
			log.Errorf("quote brokers push event, unmarshal error:%v", err)
			return
		}
		var pb PushBrokers
		if err = util.Copy(&pb, data); err != nil {
			log.Errorf("quote brokers push event, copy data error:%v", err)
			return
		}
		core.store.MergeBroker(&pb)
		f(&pb)
	}
}

func parsePushTradeFunc(f func(*PushTrade), core *core) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var err error
		var data quotev1.PushTrade
		if err = packet.Unmarshal(&data); err != nil {
			log.Errorf("quote trade push event, unmarshal error:%v", err)
			return
		}
		var pt PushTrade
		if err = util.Copy(&pt, data); err != nil {
			log.Errorf("quote trade push event, copy data error:%v", err)
			return
		}
		core.store.MergeTrade(&pt)
		f(&pt)
	}
}

func toTimes(origin []string) (times []time.Time, err error) {
	times = make([]time.Time, 0, len(origin))
	for _, dateStr := range origin {
		dt, err := util.ParseDateSimple(dateStr)
		if err != nil {
			return nil, err
		}
		times = append(times, dt)
	}
	return
}
