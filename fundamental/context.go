package fundamental

import (
	"context"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/fundamental/jsontypes"
)

// FundamentalContext is a client for fundamental financial data.
type FundamentalContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a FundamentalContext from a Config.
func NewFromCfg(cfg *config.Config) (*FundamentalContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &FundamentalContext{httpClient: httpClient}, nil
}

// NewFromEnv creates a FundamentalContext from environment variables.
func NewFromEnv() (*FundamentalContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

func parseDecimalOpt(s string) *decimal.Decimal {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return &d
}

func convertValuationMetric(m *jsontypes.ValuationMetricData) *ValuationMetricData {
	if m == nil {
		return nil
	}
	result := &ValuationMetricData{
		Desc:   m.Desc,
		High:   parseDecimalOpt(m.High),
		Low:    parseDecimalOpt(m.Low),
		Median: parseDecimalOpt(m.Median),
	}
	for _, p := range m.List {
		result.List = append(result.List, &ValuationPoint{
			Timestamp: p.Timestamp,
			Value:     parseDecimalOpt(p.Value),
		})
	}
	return result
}

func convertValuationHistoryMetric(m *jsontypes.ValuationHistoryMetric) *ValuationHistoryMetric {
	if m == nil {
		return nil
	}
	result := &ValuationHistoryMetric{
		Desc:   m.Desc,
		High:   parseDecimalOpt(m.High),
		Low:    parseDecimalOpt(m.Low),
		Median: parseDecimalOpt(m.Median),
	}
	for _, p := range m.List {
		result.List = append(result.List, &ValuationPoint{
			Timestamp: p.Timestamp,
			Value:     parseDecimalOpt(p.Value),
		})
	}
	return result
}

// FinancialReport returns financial statement data for a security.
// kind selects the statement type; period is optional (pass empty string to omit).
func (c *FundamentalContext) FinancialReport(ctx context.Context, symbol string, kind FinancialReportKind, period FinancialReportPeriod) (*FinancialReports, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("kind", string(kind))
	if period != "" {
		values.Add("report", string(period))
	}
	var resp jsontypes.FinancialReports
	if err := c.httpClient.Get(ctx, "/v1/quote/financial-reports", values, &resp); err != nil {
		return nil, err
	}
	return &FinancialReports{List: resp.List}, nil
}

// InstitutionRating returns analyst ratings for a security (latest snapshot + consensus summary).
func (c *FundamentalContext) InstitutionRating(ctx context.Context, symbol string) (*InstitutionRating, error) {
	cid := util.SymbolToCounterID(symbol)
	values := url.Values{}
	values.Add("counter_id", cid)

	var (
		latestResp  jsontypes.InstitutionRatingLatest
		summaryResp jsontypes.InstitutionRatingSummary
		latestErr   error
		summaryErr  error
		wg          sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		v := url.Values{}
		v.Add("counter_id", cid)
		latestErr = c.httpClient.Get(ctx, "/v1/quote/institution-rating-latest", v, &latestResp)
	}()
	go func() {
		defer wg.Done()
		v := url.Values{}
		v.Add("counter_id", cid)
		summaryErr = c.httpClient.Get(ctx, "/v1/quote/institution-ratings", v, &summaryResp)
	}()
	wg.Wait()
	if latestErr != nil {
		return nil, latestErr
	}
	if summaryErr != nil {
		return nil, summaryErr
	}

	result := &InstitutionRating{
		Latest: &InstitutionRatingLatest{
			Evaluate: RatingEvaluate{
				Buy: latestResp.Evaluate.Buy, Over: latestResp.Evaluate.Over,
				Hold: latestResp.Evaluate.Hold, Under: latestResp.Evaluate.Under,
				Sell: latestResp.Evaluate.Sell, NoOpinion: latestResp.Evaluate.NoOpinion,
				Total: latestResp.Evaluate.Total,
				StartDate: latestResp.Evaluate.StartDate, EndDate: latestResp.Evaluate.EndDate,
			},
			Target: RatingTarget{
				HighestPrice: parseDecimalOpt(latestResp.Target.HighestPrice),
				LowestPrice:  parseDecimalOpt(latestResp.Target.LowestPrice),
				PrevClose:    parseDecimalOpt(latestResp.Target.PrevClose),
				StartDate:    latestResp.Target.StartDate,
				EndDate:      latestResp.Target.EndDate,
			},
			IndustryId: latestResp.IndustryId, IndustryName: latestResp.IndustryName,
			IndustryRank: latestResp.IndustryRank, IndustryTotal: latestResp.IndustryTotal,
			IndustryMean: latestResp.IndustryMean, IndustryMedian: latestResp.IndustryMedian,
		},
		Summary: &InstitutionRatingSummary{
			CcySymbol: summaryResp.CcySymbol,
			Change:    parseDecimalOpt(summaryResp.Change),
			Evaluate: RatingSummaryEvaluate{
				Buy: summaryResp.Evaluate.Buy, Date: summaryResp.Evaluate.Date,
				Hold: summaryResp.Evaluate.Hold, Sell: summaryResp.Evaluate.Sell,
				StrongBuy: summaryResp.Evaluate.StrongBuy, Under: summaryResp.Evaluate.Under,
			},
			Recommend: summaryResp.Recommend,
			Target:    parseDecimalOpt(summaryResp.Target),
			UpdatedAt: summaryResp.UpdatedAt,
		},
	}
	return result, nil
}

// InstitutionRatingDetail returns historical analyst rating time-series for a security.
func (c *FundamentalContext) InstitutionRatingDetail(ctx context.Context, symbol string) (*InstitutionRatingDetail, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.InstitutionRatingDetail
	if err := c.httpClient.Get(ctx, "/v1/quote/institution-ratings/detail", values, &resp); err != nil {
		return nil, err
	}
	result := &InstitutionRatingDetail{CcySymbol: resp.CcySymbol}
	for _, item := range resp.Evaluate.List {
		result.Evaluate.List = append(result.Evaluate.List, &InstitutionRatingDetailEvaluateItem{
			Buy: item.Buy, Date: item.Date, Hold: item.Hold, Sell: item.Sell,
			StrongBuy: item.StrongBuy, NoOpinion: item.NoOpinion, Under: item.Under,
		})
	}
	result.Target = InstitutionRatingDetailTarget{
		DataPercent:        parseDecimalOpt(resp.Target.DataPercent),
		PredictionAccuracy: parseDecimalOpt(resp.Target.PredictionAccuracy),
		UpdatedAt:          resp.Target.UpdatedAt,
	}
	for _, item := range resp.Target.List {
		result.Target.List = append(result.Target.List, &InstitutionRatingDetailTargetItem{
			AvgTarget: parseDecimalOpt(item.AvgTarget), Date: item.Date,
			MaxTarget: parseDecimalOpt(item.MaxTarget), MinTarget: parseDecimalOpt(item.MinTarget),
			Meet: item.Meet, Price: parseDecimalOpt(item.Price), Timestamp: item.Timestamp,
		})
	}
	return result, nil
}

// Dividend returns dividend history for a security.
func (c *FundamentalContext) Dividend(ctx context.Context, symbol string) (*DividendList, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.DividendList
	if err := c.httpClient.Get(ctx, "/v1/quote/dividends", values, &resp); err != nil {
		return nil, err
	}
	return convertDividendList(&resp), nil
}

// DividendDetail returns detailed dividend information for a security.
func (c *FundamentalContext) DividendDetail(ctx context.Context, symbol string) (*DividendList, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.DividendList
	if err := c.httpClient.Get(ctx, "/v1/quote/dividends/details", values, &resp); err != nil {
		return nil, err
	}
	return convertDividendList(&resp), nil
}

func convertDividendList(resp *jsontypes.DividendList) *DividendList {
	result := &DividendList{}
	for _, item := range resp.List {
		result.List = append(result.List, &DividendItem{
			Symbol:      util.CounterIDToSymbol(item.CounterId),
			Id:          item.Id,
			Desc:        item.Desc,
			RecordDate:  item.RecordDate,
			ExDate:      item.ExDate,
			PaymentDate: item.PaymentDate,
		})
	}
	return result
}

// ForecastEps returns EPS forecasts for a security.
func (c *FundamentalContext) ForecastEps(ctx context.Context, symbol string) (*ForecastEps, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.ForecastEps
	if err := c.httpClient.Get(ctx, "/v1/quote/forecast-eps", values, &resp); err != nil {
		return nil, err
	}
	result := &ForecastEps{}
	for _, item := range resp.Items {
		result.Items = append(result.Items, &ForecastEpsItem{
			ForecastEpsMedian:  parseDecimalOpt(item.ForecastEpsMedian),
			ForecastEpsMean:    parseDecimalOpt(item.ForecastEpsMean),
			ForecastEpsLowest:  parseDecimalOpt(item.ForecastEpsLowest),
			ForecastEpsHighest: parseDecimalOpt(item.ForecastEpsHighest),
			InstitutionTotal:   item.InstitutionTotal,
			InstitutionUp:      item.InstitutionUp,
			InstitutionDown:    item.InstitutionDown,
			ForecastStartDate:  item.ForecastStartDate,
			ForecastEndDate:    item.ForecastEndDate,
		})
	}
	return result, nil
}

// Consensus returns financial consensus estimates for a security.
func (c *FundamentalContext) Consensus(ctx context.Context, symbol string) (*FinancialConsensus, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.FinancialConsensus
	if err := c.httpClient.Get(ctx, "/v1/quote/financial-consensus-detail", values, &resp); err != nil {
		return nil, err
	}
	result := &FinancialConsensus{
		CurrentIndex:  resp.CurrentIndex,
		Currency:      resp.Currency,
		OptPeriods:    resp.OptPeriods,
		CurrentPeriod: resp.CurrentPeriod,
	}
	for _, r := range resp.List {
		report := &ConsensusReport{
			FiscalYear:   r.FiscalYear,
			FiscalPeriod: r.FiscalPeriod,
			PeriodText:   r.PeriodText,
		}
		for _, d := range r.Details {
			report.Details = append(report.Details, &ConsensusDetail{
				Key: d.Key, Name: d.Name, Desc: d.Desc,
				Actual: parseDecimalOpt(d.Actual), Estimate: parseDecimalOpt(d.Estimate),
				CompValue: parseDecimalOpt(d.CompValue), CompDesc: d.CompDesc,
				Comp: d.Comp, IsReleased: d.IsReleased,
			})
		}
		result.List = append(result.List, report)
	}
	return result, nil
}

// Valuation returns valuation metrics (PE/PB/PS/dividend yield) for a security.
func (c *FundamentalContext) Valuation(ctx context.Context, symbol string) (*ValuationData, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("indicator", "pe")
	values.Add("range", "1")
	var resp jsontypes.ValuationData
	if err := c.httpClient.Get(ctx, "/v1/quote/valuation", values, &resp); err != nil {
		return nil, err
	}
	return &ValuationData{Metrics: ValuationMetricsData{
		Pe:     convertValuationMetric(resp.Metrics.Pe),
		Pb:     convertValuationMetric(resp.Metrics.Pb),
		Ps:     convertValuationMetric(resp.Metrics.Ps),
		DvdYld: convertValuationMetric(resp.Metrics.DvdYld),
	}}, nil
}

// ValuationHistory returns historical valuation data for a security.
func (c *FundamentalContext) ValuationHistory(ctx context.Context, symbol string) (*ValuationHistoryResponse, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.ValuationHistoryResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/valuation/detail", values, &resp); err != nil {
		return nil, err
	}
	return &ValuationHistoryResponse{History: ValuationHistoryData{
		Metrics: ValuationHistoryMetrics{
			Pe: convertValuationHistoryMetric(resp.History.Metrics.Pe),
			Pb: convertValuationHistoryMetric(resp.History.Metrics.Pb),
			Ps: convertValuationHistoryMetric(resp.History.Metrics.Ps),
		},
	}}, nil
}

// IndustryValuation returns valuation comparison against industry peers.
func (c *FundamentalContext) IndustryValuation(ctx context.Context, symbol string) (*IndustryValuationList, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.IndustryValuationList
	if err := c.httpClient.Get(ctx, "/v1/quote/industry-valuation-comparison", values, &resp); err != nil {
		return nil, err
	}
	result := &IndustryValuationList{}
	for _, item := range resp.List {
		iv := &IndustryValuationItem{
			Symbol:         util.CounterIDToSymbol(item.CounterId),
			Name:           item.Name,
			Currency:       item.Currency,
			Assets:         parseDecimalOpt(item.Assets),
			Bps:            parseDecimalOpt(item.Bps),
			Eps:            parseDecimalOpt(item.Eps),
			Dps:            parseDecimalOpt(item.Dps),
			DivYld:         parseDecimalOpt(item.DivYld),
			DivPayoutRatio: parseDecimalOpt(item.DivPayoutRatio),
			FiveYAvgDps:    parseDecimalOpt(item.FiveYAvgDps),
			Pe:             parseDecimalOpt(item.Pe),
		}
		for _, h := range item.History {
			iv.History = append(iv.History, &IndustryValuationHistory{
				Date: h.Date,
				Pe:   parseDecimalOpt(h.Pe),
				Pb:   parseDecimalOpt(h.Pb),
				Ps:   parseDecimalOpt(h.Ps),
			})
		}
		result.List = append(result.List, iv)
	}
	return result, nil
}

func convertValuationDist(d *jsontypes.ValuationDist) *ValuationDist {
	if d == nil {
		return nil
	}
	return &ValuationDist{
		Low:       parseDecimalOpt(d.Low),
		High:      parseDecimalOpt(d.High),
		Median:    parseDecimalOpt(d.Median),
		Value:     parseDecimalOpt(d.Value),
		Ranking:   parseDecimalOpt(d.Ranking),
		RankIndex: d.RankIndex,
		RankTotal: d.RankTotal,
	}
}

// IndustryValuationDist returns valuation distribution within the industry.
func (c *FundamentalContext) IndustryValuationDist(ctx context.Context, symbol string) (*IndustryValuationDist, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.IndustryValuationDist
	if err := c.httpClient.Get(ctx, "/v1/quote/industry-valuation-distribution", values, &resp); err != nil {
		return nil, err
	}
	return &IndustryValuationDist{
		Pe: convertValuationDist(resp.Pe),
		Pb: convertValuationDist(resp.Pb),
		Ps: convertValuationDist(resp.Ps),
	}, nil
}

// Company returns company overview information.
func (c *FundamentalContext) Company(ctx context.Context, symbol string) (*CompanyOverview, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.CompanyOverview
	if err := c.httpClient.Get(ctx, "/v1/quote/comp-overview", values, &resp); err != nil {
		return nil, err
	}
	return &CompanyOverview{
		Name: resp.Name, CompanyName: resp.CompanyName, Founded: resp.Founded,
		ListingDate: resp.ListingDate, Market: resp.Market, Region: resp.Region,
		Address: resp.Address, OfficeAddress: resp.OfficeAddress, Website: resp.Website,
		IssuePrice: parseDecimalOpt(resp.IssuePrice), SharesOffered: resp.SharesOffered,
		Chairman: resp.Chairman, Secretary: resp.Secretary, AuditInst: resp.AuditInst,
		Category: resp.Category, YearEnd: resp.YearEnd, Employees: resp.Employees,
		Phone: resp.Phone, Fax: resp.Fax, Email: resp.Email, LegalRepr: resp.LegalRepr,
		Manager: resp.Manager, BusLicense: resp.BusLicense, AccountingFirm: resp.AccountingFirm,
		SecuritiesRep: resp.SecuritiesRep, LegalCounsel: resp.LegalCounsel,
		ZipCode: resp.ZipCode, Ticker: resp.Ticker, Icon: resp.Icon,
		Profile: resp.Profile, AdsRatio: resp.AdsRatio, Sector: resp.Sector,
	}, nil
}

// Executive returns executive and board member information.
func (c *FundamentalContext) Executive(ctx context.Context, symbol string) (*ExecutiveList, error) {
	values := url.Values{}
	values.Add("counter_ids", util.SymbolToCounterID(symbol))
	var resp jsontypes.ExecutiveList
	if err := c.httpClient.Get(ctx, "/v1/quote/company-professionals", values, &resp); err != nil {
		return nil, err
	}
	result := &ExecutiveList{}
	for _, g := range resp.ProfessionalList {
		group := &ExecutiveGroup{
			Symbol:     util.CounterIDToSymbol(g.CounterId),
			ForwardUrl: g.ForwardUrl,
			Total:      g.Total,
		}
		for _, p := range g.Professionals {
			group.Professionals = append(group.Professionals, &Professional{
				Id: p.Id, Name: p.Name, NameZhcn: p.NameZhcn, NameEn: p.NameEn,
				Title: p.Title, Biography: p.Biography, Photo: p.Photo, WikiUrl: p.WikiUrl,
			})
		}
		result.ProfessionalList = append(result.ProfessionalList, group)
	}
	return result, nil
}

// Shareholder returns major shareholders for a security.
func (c *FundamentalContext) Shareholder(ctx context.Context, symbol string) (*ShareholderList, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.ShareholderList
	if err := c.httpClient.Get(ctx, "/v1/quote/shareholders", values, &resp); err != nil {
		return nil, err
	}
	result := &ShareholderList{ForwardUrl: resp.ForwardUrl, Total: resp.Total}
	for _, s := range resp.Shareholders {
		sh := &Shareholder{
			ShareholderId:   s.ShareholderId,
			ShareholderName: s.ShareholderName,
			InstitutionType: s.InstitutionType,
			PercentOfShares: parseDecimalOpt(s.PercentOfShares),
			SharesChanged:   parseDecimalOpt(s.SharesChanged),
			ReportDate:      s.ReportDate,
		}
		for _, stock := range s.Stocks {
			sh.Stocks = append(sh.Stocks, &ShareholderStock{
				Symbol: util.CounterIDToSymbol(stock.CounterId),
				Code:   stock.Code,
				Market: stock.Market,
				Chg:    stock.Chg,
			})
		}
		result.Shareholders = append(result.Shareholders, sh)
	}
	return result, nil
}

// FundHolder returns funds and ETFs that hold a security.
func (c *FundamentalContext) FundHolder(ctx context.Context, symbol string) (*FundHolders, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.FundHolders
	if err := c.httpClient.Get(ctx, "/v1/quote/fund-holders", values, &resp); err != nil {
		return nil, err
	}
	result := &FundHolders{}
	for _, f := range resp.Lists {
		result.Lists = append(result.Lists, &FundHolder{
			Code:          f.Code,
			Symbol:        util.CounterIDToSymbol(f.CounterId),
			Currency:      f.Currency,
			Name:          f.Name,
			PositionRatio: parseDecimalOpt(f.PositionRatio),
			ReportDate:    f.ReportDate,
		})
	}
	return result, nil
}

// CorpAction returns corporate actions (dividends, splits, buybacks, etc.).
func (c *FundamentalContext) CorpAction(ctx context.Context, symbol string) (*CorpActions, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("req_type", "1")
	values.Add("version", "3")
	var resp jsontypes.CorpActions
	if err := c.httpClient.Get(ctx, "/v1/quote/company-act", values, &resp); err != nil {
		return nil, err
	}
	result := &CorpActions{}
	for _, item := range resp.Items {
		ca := &CorpActionItem{
			Id: item.Id, Date: item.Date, DateStr: item.DateStr,
			DateType: item.DateType, DateZone: item.DateZone,
			ActType: item.ActType, ActDesc: item.ActDesc, Action: item.Action,
			Recent: item.Recent, IsDelay: item.IsDelay, DelayContent: item.DelayContent,
			Security: item.Security,
		}
		if item.Live != nil {
			ca.Live = &CorpActionLive{
				Id:        item.Live.Id,
				Status:    item.Live.Status,
				StartedAt: item.Live.StartedAt,
				Name:      item.Live.Name,
				Icon:      item.Live.Icon,
			}
		}
		result.Items = append(result.Items, ca)
	}
	return result, nil
}

// InvestRelation returns investor relations / investment holdings.
func (c *FundamentalContext) InvestRelation(ctx context.Context, symbol string) (*InvestRelations, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("count", "0")
	var resp jsontypes.InvestRelations
	if err := c.httpClient.Get(ctx, "/v1/quote/invest-relations", values, &resp); err != nil {
		return nil, err
	}
	result := &InvestRelations{ForwardUrl: resp.ForwardUrl}
	for _, s := range resp.InvestSecurities {
		result.InvestSecurities = append(result.InvestSecurities, &InvestSecurity{
			CompanyId:       s.CompanyId,
			CompanyName:     s.CompanyName,
			CompanyNameEn:   s.CompanyNameEn,
			CompanyNameZhcn: s.CompanyNameZhcn,
			Symbol:          util.CounterIDToSymbol(s.CounterId),
			Currency:        s.Currency,
			PercentOfShares: parseDecimalOpt(s.PercentOfShares),
			SharesRank:      s.SharesRank,
			SharesValue:     parseDecimalOpt(s.SharesValue),
		})
	}
	return result, nil
}

// Operating returns operating metrics and financial report summaries.
func (c *FundamentalContext) Operating(ctx context.Context, symbol string) (*OperatingList, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.OperatingList
	if err := c.httpClient.Get(ctx, "/v1/quote/operatings", values, &resp); err != nil {
		return nil, err
	}
	result := &OperatingList{}
	for _, item := range resp.List {
		op := &OperatingItem{
			Id: item.Id, Report: item.Report, Title: item.Title,
			Txt: item.Txt, Latest: item.Latest, Keywords: item.Keywords, WebUrl: item.WebUrl,
			Financial: OperatingFinancial{
				Code: item.Financial.Code, CounterId: item.Financial.CounterId,
				Currency: item.Financial.Currency, Name: item.Financial.Name,
				Region: item.Financial.Region, Report: item.Financial.Report,
				ReportTxt: item.Financial.ReportTxt,
			},
		}
		for _, ind := range item.Financial.Indicators {
			op.Financial.Indicators = append(op.Financial.Indicators, &OperatingIndicator{
				FieldName: ind.FieldName, IndicatorName: ind.IndicatorName,
				IndicatorValue: ind.IndicatorValue, Yoy: parseDecimalOpt(ind.Yoy),
			})
		}
		result.List = append(result.List, op)
	}
	return result, nil
}

// Buyback returns buyback data for a security.
func (c *FundamentalContext) Buyback(ctx context.Context, symbol string) (*BuybackData, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.BuybackData
	if err := c.httpClient.Get(ctx, "/v1/quote/buy-backs", values, &resp); err != nil {
		return nil, err
	}
	result := &BuybackData{}
	if resp.RecentBuybacks != nil {
		result.RecentBuybacks = &RecentBuybacks{
			Currency:           resp.RecentBuybacks.Currency,
			NetBuybackTtm:      parseDecimalOpt(resp.RecentBuybacks.NetBuybackTtm),
			NetBuybackYieldTtm: parseDecimalOpt(resp.RecentBuybacks.NetBuybackYieldTtm),
		}
	}
	for _, h := range resp.BuybackHistory {
		result.BuybackHistory = append(result.BuybackHistory, &BuybackHistoryItem{
			FiscalYear:           h.FiscalYear,
			FiscalYearRange:      h.FiscalYearRange,
			NetBuyback:           parseDecimalOpt(h.NetBuyback),
			NetBuybackYield:      parseDecimalOpt(h.NetBuybackYield),
			NetBuybackGrowthRate: parseDecimalOpt(h.NetBuybackGrowthRate),
			Currency:             h.Currency,
		})
	}
	for _, r := range resp.BuybackRatios {
		result.BuybackRatios = append(result.BuybackRatios, &BuybackRatios{
			NetBuybackPayoutRatio:     parseDecimalOpt(r.NetBuybackPayoutRatio),
			NetBuybackToCashflowRatio: parseDecimalOpt(r.NetBuybackToCashflowRatio),
		})
	}
	return result, nil
}

// Ratings returns stock ratings for a security.
func (c *FundamentalContext) Ratings(ctx context.Context, symbol string) (*StockRatings, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.StockRatings
	if err := c.httpClient.Get(ctx, "/v1/quote/ratings", values, &resp); err != nil {
		return nil, err
	}
	result := &StockRatings{
		StyleTxtName:        resp.StyleTxtName,
		ScaleTxtName:        resp.ScaleTxtName,
		ReportPeriodTxt:     resp.ReportPeriodTxt,
		MultiScore:          resp.MultiScore,
		MultiLetter:         resp.MultiLetter,
		MultiScoreChange:    resp.MultiScoreChange,
		IndustryName:        resp.IndustryName,
		IndustryRank:        resp.IndustryRank,
		IndustryTotal:       resp.IndustryTotal,
		IndustryMeanScore:   resp.IndustryMeanScore,
		IndustryMedianScore: resp.IndustryMedianScore,
	}
	for _, cat := range resp.Ratings {
		rc := &RatingCategory{Kind: cat.Kind}
		for _, g := range cat.SubIndicators {
			group := &RatingSubIndicatorGroup{
				Indicator: RatingIndicator{
					Name: g.Indicator.Name, Score: g.Indicator.Score, Letter: g.Indicator.Letter,
				},
			}
			for _, leaf := range g.SubIndicators {
				group.SubIndicators = append(group.SubIndicators, &RatingLeafIndicator{
					Name: leaf.Name, Value: leaf.Value, ValueType: leaf.ValueType,
					Score: leaf.Score, Letter: leaf.Letter,
				})
			}
			rc.SubIndicators = append(rc.SubIndicators, group)
		}
		result.Ratings = append(result.Ratings, rc)
	}
	return result, nil
}
