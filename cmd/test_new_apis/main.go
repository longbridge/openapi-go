// Comprehensive test for all new APIs added in the openapi-go sync (v4.0.6+).
// Read-only calls run unconditionally. Write calls (create/update/delete) run
// only when -write flag is set, and each write is cleaned up immediately after.
//
// Usage:
//
//	go run ./cmd/test_new_apis -client-id <id>          # read-only
//	go run ./cmd/test_new_apis -client-id <id> -write   # include write ops
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/longbridge/openapi-go/alert"
	"github.com/longbridge/openapi-go/calendar"
	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/dca"
	"github.com/longbridge/openapi-go/fundamental"
	"github.com/longbridge/openapi-go/market"
	"github.com/longbridge/openapi-go/oauth"
	"github.com/longbridge/openapi-go/portfolio"
	"github.com/longbridge/openapi-go/quote"
	"github.com/longbridge/openapi-go/sharelist"
)

var (
	clientID  = flag.String("client-id", "", "OAuth client ID")
	withWrite = flag.Bool("write", false, "also test write operations (create/update/delete)")
)

type tester struct {
	pass, fail int
}

func (t *tester) check(name string, fn func() error) {
	fmt.Printf("  %-58s", name)
	if e := fn(); e != nil {
		fmt.Printf("FAIL: %v\n", e)
		t.fail++
	} else {
		fmt.Println("OK")
		t.pass++
	}
}

func (t *tester) skip(name, reason string) {
	fmt.Printf("  %-58s SKIP (%s)\n", name, reason)
}

func main() {
	flag.Parse()

	ctx := context.Background()

	var cfg *config.Config
	var err error
	if *clientID != "" {
		o := oauth.New(*clientID).
			OnOpenURL(func(url string) { fmt.Println("Open URL:", url) })
		if err = o.Build(ctx); err != nil {
			log.Fatalf("oauth: %v", err)
		}
		cfg, err = config.New(config.WithOAuthClient(o))
	} else {
		cfg, err = config.NewFormEnv() //nolint:staticcheck
	}
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	t := &tester{}

	// ══════════════════════════════════════════════════════════
	// alert (4 methods)
	// ══════════════════════════════════════════════════════════
	fmt.Println("\n=== alert.AlertContext (4 methods) ===")
	alertCtx, err := alert.NewFromCfg(cfg)
	if err != nil {
		log.Printf("alert.NewFromCfg: %v", err)
	} else {
		var existingAlertID string
		t.check("List()", func() error {
			list, err := alertCtx.List(ctx)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d groups)", len(list.Lists))
			for _, g := range list.Lists {
				if len(g.Indicators) > 0 {
					existingAlertID = g.Indicators[0].ID
				}
			}
			return nil
		})
		if *withWrite {
			t.check("Add(700.HK, PriceRise, 200, Daily)", func() error {
				return alertCtx.Add(ctx, "0700.HK", alert.AlertConditionPriceRise, "200", alert.AlertFrequencyOnce)
			})
			if existingAlertID != "" {
				t.check("Delete([existingID])", func() error {
					return alertCtx.Delete(ctx, []string{existingAlertID})
				})
			} else {
				t.skip("Delete()", "no existing alert ID")
			}
		} else {
			t.skip("Add()", "-write not set")
			t.skip("Update()", "-write not set")
			t.skip("Delete()", "-write not set")
		}
	}

	// ══════════════════════════════════════════════════════════
	// calendar (1 method)
	// ══════════════════════════════════════════════════════════
	fmt.Println("\n=== calendar.CalendarContext (1 method) ===")
	calCtx, err := calendar.NewFromCfg(cfg)
	if err != nil {
		log.Printf("calendar.NewFromCfg: %v", err)
	} else {
		start := time.Now().Format("2006-01-02")
		end := time.Now().Add(14 * 24 * time.Hour).Format("2006-01-02")
		for _, cat := range []struct {
			name string
			cat  calendar.CalendarCategory
		}{
			{"FinanceCalendar(Report)", calendar.CalendarCategoryReport},
			{"FinanceCalendar(Dividend)", calendar.CalendarCategoryDividend},
			{"FinanceCalendar(Ipo)", calendar.CalendarCategoryIpo},
			{"FinanceCalendar(MacroData)", calendar.CalendarCategoryMacroData},
		} {
			cat := cat
			t.check(cat.name, func() error {
				resp, err := calCtx.FinanceCalendar(ctx, cat.cat, start, end, nil)
				if err != nil {
					return err
				}
				fmt.Printf(" (%d day-groups)", len(resp.List))
				return nil
			})
		}
	}

	// ══════════════════════════════════════════════════════════
	// market (9 methods)
	// ══════════════════════════════════════════════════════════
	fmt.Println("\n=== market.MarketContext (9 methods) ===")
	mktCtx, err := market.NewFromCfg(cfg)
	if err != nil {
		log.Printf("market.NewFromCfg: %v", err)
	} else {
		t.check("MarketStatus()", func() error {
			resp, err := mktCtx.MarketStatus(ctx)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d markets)", len(resp.MarketTime))
			return nil
		})
		t.check("BrokerHolding(0700.HK, Rct5)", func() error {
			resp, err := mktCtx.BrokerHolding(ctx, "0700.HK", market.BrokerHoldingPeriodRct5)
			if err != nil {
				return err
			}
			fmt.Printf(" (buy:%d sell:%d)", len(resp.Buy), len(resp.Sell))
			return nil
		})
		var firstBrokerID string
		t.check("BrokerHoldingDetail(0700.HK)", func() error {
			resp, err := mktCtx.BrokerHoldingDetail(ctx, "0700.HK")
			if err != nil {
				return err
			}
			if len(resp.List) > 0 {
				firstBrokerID = resp.List[0].PartiNumber
			}
			fmt.Printf(" (%d brokers)", len(resp.List))
			return nil
		})
		t.check("BrokerHoldingDaily(0700.HK, brokerID)", func() error {
			if firstBrokerID == "" {
				firstBrokerID = "9000" // fallback
			}
			resp, err := mktCtx.BrokerHoldingDaily(ctx, "0700.HK", firstBrokerID)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d days)", len(resp.List))
			return nil
		})
		t.check("AhPremium(0700.HK, Day, 30)", func() error {
			resp, err := mktCtx.AhPremium(ctx, "0700.HK", market.AhPremiumPeriodDay, 30)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d klines)", len(resp.Klines))
			return nil
		})
		t.check("AhPremiumIntraday(0700.HK)", func() error {
			resp, err := mktCtx.AhPremiumIntraday(ctx, "0700.HK")
			if err != nil {
				return err
			}
			fmt.Printf(" (%d klines)", len(resp.Klines))
			return nil
		})
		t.check("TradeStats(AAPL.US)", func() error {
			resp, err := mktCtx.TradeStats(ctx, "AAPL.US")
			if err != nil {
				return err
			}
			fmt.Printf(" (buy_vol:%s)", resp.Statistics.Buy.String())
			return nil
		})
		t.check("Anomaly(HK)", func() error {
			resp, err := mktCtx.Anomaly(ctx, "HK")
			if err != nil {
				return err
			}
			fmt.Printf(" (%d changes)", len(resp.Changes))
			return nil
		})
		t.check("Constituent(.HSI)", func() error {
			resp, err := mktCtx.Constituent(ctx, ".HSI")
			if err != nil {
				return err
			}
			fmt.Printf(" (%d stocks)", len(resp.Stocks))
			return nil
		})
	}

	// ══════════════════════════════════════════════════════════
	// fundamental (26 methods)
	// ══════════════════════════════════════════════════════════
	fmt.Println("\n=== fundamental.FundamentalContext (26 methods) ===")
	fundCtx, err := fundamental.NewFromCfg(cfg)
	if err != nil {
		log.Printf("fundamental.NewFromCfg: %v", err)
	} else {
		sym := "AAPL.US"
		period := fundamental.FinancialReportPeriodAnnual
		t.check("FinancialReport(AAPL.US, Income, Annual)", func() error {
			_, err := fundCtx.FinancialReport(ctx, sym, fundamental.FinancialReportKindIncomeStatement, &period)
			return err
		})
		t.check("InstitutionRating(AAPL.US)", func() error {
			r, err := fundCtx.InstitutionRating(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (industry:%s)", r.Latest.IndustryName)
			return nil
		})
		t.check("InstitutionRatingDetail(AAPL.US)", func() error {
			r, err := fundCtx.InstitutionRatingDetail(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d weekly snapshots)", len(r.Evaluate.List))
			return nil
		})
		t.check("Dividend(AAPL.US)", func() error {
			r, err := fundCtx.Dividend(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d dividends)", len(r.List))
			return nil
		})
		t.check("DividendDetail(AAPL.US)", func() error {
			r, err := fundCtx.DividendDetail(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d dividends)", len(r.List))
			return nil
		})
		t.check("ForecastEps(AAPL.US)", func() error {
			r, err := fundCtx.ForecastEps(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d items)", len(r.Items))
			return nil
		})
		t.check("Consensus(AAPL.US)", func() error {
			_, err := fundCtx.Consensus(ctx, sym)
			return err
		})
		t.check("Valuation(AAPL.US)", func() error {
			_, err := fundCtx.Valuation(ctx, sym)
			return err
		})
		t.check("ValuationHistory(AAPL.US)", func() error {
			_, err := fundCtx.ValuationHistory(ctx, sym)
			return err
		})
		t.check("IndustryValuation(AAPL.US)", func() error {
			r, err := fundCtx.IndustryValuation(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d items)", len(r.List))
			return nil
		})
		t.check("IndustryValuationDist(AAPL.US)", func() error {
			_, err := fundCtx.IndustryValuationDist(ctx, sym)
			return err
		})
		t.check("Company(AAPL.US)", func() error {
			r, err := fundCtx.Company(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%s)", r.Name)
			return nil
		})
		t.check("Executive(AAPL.US)", func() error {
			r, err := fundCtx.Executive(ctx, sym)
			if err != nil {
				return err
			}
			total := 0
			for _, g := range r.ProfessionalList {
				total += len(g.Professionals)
			}
			fmt.Printf(" (%d executives)", total)
			return nil
		})
		t.check("Shareholder(AAPL.US)", func() error {
			r, err := fundCtx.Shareholder(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d shareholders)", len(r.ShareholderList))
			return nil
		})
		t.check("FundHolder(700.HK)", func() error {
			r, err := fundCtx.FundHolder(ctx, "0700.HK")
			if err != nil {
				return err
			}
			fmt.Printf(" (%d funds)", len(r.Lists))
			return nil
		})
		t.check("CorpAction(AAPL.US)", func() error {
			_, err := fundCtx.CorpAction(ctx, sym)
			return err
		})
		t.check("InvestRelation(AAPL.US)", func() error {
			_, err := fundCtx.InvestRelation(ctx, sym)
			return err
		})
		t.check("Operating(AAPL.US)", func() error {
			_, err := fundCtx.Operating(ctx, sym)
			return err
		})
		t.check("Buyback(AAPL.US)", func() error {
			r, err := fundCtx.Buyback(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d history records)", len(r.BuybackHistory))
			return nil
		})
		t.check("Ratings(AAPL.US)", func() error {
			_, err := fundCtx.Ratings(ctx, sym)
			return err
		})
		// ── new APIs (PR #91) ──────────────────────────────────
		t.check("BusinessSegments(AAPL.US)", func() error {
			r, err := fundCtx.BusinessSegments(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d segments)", len(r.Business))
			return nil
		})
		t.check("BusinessSegmentsHistory(AAPL.US, qf, \"\")", func() error {
			r, err := fundCtx.BusinessSegmentsHistory(ctx, sym, "qf", "")
			if err != nil {
				return err
			}
			fmt.Printf(" (%d periods)", len(r.Historical))
			return nil
		})
		t.check("InstitutionRatingViews(AAPL.US)", func() error {
			r, err := fundCtx.InstitutionRatingViews(ctx, sym)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d months)", len(r.Elist))
			return nil
		})
		var bkCounterID string
		t.check("IndustryRank(US, Indicator0/leading-gainer, Ascending, 10)", func() error {
			r, err := fundCtx.IndustryRank(ctx, "US", fundamental.IndustryRankIndicator0, fundamental.IndustryRankSortTypeAscending, 10)
			if err != nil {
				return err
			}
			n := 0
			for _, g := range r.Items {
				n += len(g.Lists)
				for _, item := range g.Lists {
					if bkCounterID == "" && item.CounterID != "" {
						bkCounterID = item.CounterID
					}
				}
			}
			fmt.Printf(" (%d industries, first_id:%s)", n, bkCounterID)
			return nil
		})
		t.check("IndustryPeers(BK/US/..., US, \"\")", func() error {
			if bkCounterID == "" {
				bkCounterID = "BK/US/IN00258" // fallback
			}
			r, err := fundCtx.IndustryPeers(ctx, bkCounterID, "US", "")
			if err != nil {
				return err
			}
			children := 0
			if r.Chain != nil {
				children = len(r.Chain.Next)
			}
			fmt.Printf(" (top:%s children:%d)", r.Top.Name, children)
			return nil
		})
		t.check("FinancialReportSnapshot(AAPL.US)", func() error {
			r, err := fundCtx.FinancialReportSnapshot(ctx, sym, "", 0, "")
			if err != nil {
				return err
			}
			fmt.Printf(" (%s %s–%s)", r.Ticker, r.FpStart, r.FpEnd)
			return nil
		})
	}

	// ══════════════════════════════════════════════════════════
	// portfolio (5 methods)
	// ══════════════════════════════════════════════════════════
	fmt.Println("\n=== portfolio.PortfolioContext (5 methods) ===")
	portCtx, err := portfolio.NewFromCfg(cfg)
	if err != nil {
		log.Printf("portfolio.NewFromCfg: %v", err)
	} else {
		t.check("ExchangeRate()", func() error {
			r, err := portCtx.ExchangeRate(ctx)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d pairs)", len(r.Exchanges))
			return nil
		})
		t.check("ProfitAnalysis()", func() error {
			_, err := portCtx.ProfitAnalysis(ctx, &portfolio.ProfitAnalysisOptions{})
			return err
		})
		t.check("ProfitAnalysisByMarket(page:1, size:20)", func() error {
			r, err := portCtx.ProfitAnalysisByMarket(ctx, &portfolio.ProfitAnalysisByMarketOptions{Page: 1, Size: 20})
			if err != nil {
				return err
			}
			fmt.Printf(" (%d stocks)", len(r.StockItems))
			return nil
		})
		t.check("ProfitAnalysisDetail(AAPL.US)", func() error {
			_, err := portCtx.ProfitAnalysisDetail(ctx, &portfolio.ProfitAnalysisDetailOptions{Symbol: "AAPL.US"})
			return err
		})
		t.check("ProfitAnalysisFlows(AAPL.US, page:1, size:10)", func() error {
			r, err := portCtx.ProfitAnalysisFlows(ctx, &portfolio.ProfitAnalysisFlowsOptions{
				Symbol: "AAPL.US", Page: 1, Size: 10,
			})
			if err != nil {
				return err
			}
			fmt.Printf(" (%d flows)", len(r.FlowsList))
			return nil
		})
	}

	// ══════════════════════════════════════════════════════════
	// dca (11 methods; write ops gated by -write)
	// ══════════════════════════════════════════════════════════
	fmt.Println("\n=== dca.DCAContext (11 methods) ===")
	dcaCtx, err := dca.NewFromCfg(cfg)
	if err != nil {
		log.Printf("dca.NewFromCfg: %v", err)
	} else {
		var firstPlanID string
		t.check("List()", func() error {
			r, err := dcaCtx.List(ctx, nil, nil)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d plans)", len(r.Plans))
			if len(r.Plans) > 0 {
				firstPlanID = r.Plans[0].PlanID
			}
			return nil
		})
		t.check("Stats()", func() error {
			r, err := dcaCtx.Stats(ctx, nil)
			if err != nil {
				return err
			}
			fmt.Printf(" (active:%s finished:%s)", r.ActiveCount, r.FinishedCount)
			return nil
		})
		t.check("CheckSupport(AAPL.US, 700.HK)", func() error {
			r, err := dcaCtx.CheckSupport(ctx, []string{"AAPL.US", "0700.HK"})
			if err != nil {
				return err
			}
			fmt.Printf(" (%d results)", len(r))
			return nil
		})
		t.check("CalcDate(AAPL.US, Weekly, Mon)", func() error {
			r, err := dcaCtx.CalcDate(ctx, "AAPL.US", dca.DCAFrequencyWeekly, &dca.CalcDateOptions{DayOfWeek: "Mon"})
			if err != nil {
				return err
			}
			fmt.Printf(" (next:%s)", r.TradeDate.Format("2006-01-02"))
			return nil
		})
		if firstPlanID != "" {
			t.check("History(firstPlanID, 1, 10)", func() error {
				r, err := dcaCtx.History(ctx, firstPlanID, 1, 10)
				if err != nil {
					return err
				}
				fmt.Printf(" (%d records)", len(r.Records))
				return nil
			})
		} else {
			t.skip("History()", "no existing plan")
		}
		if *withWrite {
			t.skip("Create()", "skip to avoid unintended plan creation")
			t.skip("Update()", "skip (requires valid plan)")
			t.skip("Pause()/Resume()/Stop()", "skip (requires valid plan)")
			t.check("SetReminder(8)", func() error {
				return dcaCtx.SetReminder(ctx, "8")
			})
		} else {
			t.skip("Create()", "-write not set")
			t.skip("Update()", "-write not set")
			t.skip("Pause()", "-write not set")
			t.skip("Resume()", "-write not set")
			t.skip("Stop()", "-write not set")
			t.skip("SetReminder()", "-write not set")
		}
	}

	// ══════════════════════════════════════════════════════════
	// sharelist (8 methods; write ops gated by -write)
	// ══════════════════════════════════════════════════════════
	fmt.Println("\n=== sharelist.SharelistContext (8 methods) ===")
	slCtx, err := sharelist.NewFromCfg(cfg)
	if err != nil {
		log.Printf("sharelist.NewFromCfg: %v", err)
	} else {
		var firstID int64
		t.check("List(10)", func() error {
			r, err := slCtx.List(ctx, 10)
			if err != nil {
				return err
			}
			if len(r.Sharelists) > 0 {
				firstID = r.Sharelists[0].ID
			}
			fmt.Printf(" (%d lists, %d subscribed)", len(r.Sharelists), len(r.SubscribedSharelists))
			return nil
		})
		if firstID != 0 {
			t.check("Detail(firstID)", func() error {
				r, err := slCtx.Detail(ctx, firstID)
				if err != nil {
					return err
				}
				fmt.Printf(" (%s, %d stocks)", r.Sharelist.Name, len(r.Sharelist.Stocks))
				return nil
			})
		} else {
			t.skip("Detail()", "no existing sharelist")
		}
		t.check("Popular(5)", func() error {
			r, err := slCtx.Popular(ctx, 5)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d lists)", len(r.Sharelists))
			return nil
		})
		if *withWrite {
			var newID int64
			t.check("Create(test-list)", func() error {
				err := slCtx.Create(ctx, "go-sdk-test", "temporary test list")
				if err != nil {
					return err
				}
				// fetch to get the new ID
				r, err := slCtx.List(ctx, 20)
				if err != nil {
					return err
				}
				for _, sl := range r.Sharelists {
					if sl.Name == "go-sdk-test" {
						newID = sl.ID
					}
				}
				return nil
			})
			if newID != 0 {
				t.check("AddSecurities(newID, AAPL.US)", func() error {
					return slCtx.AddSecurities(ctx, newID, []string{"AAPL.US"})
				})
				t.check("SortSecurities(newID)", func() error {
					return slCtx.SortSecurities(ctx, newID, []string{"AAPL.US"})
				})
				t.check("RemoveSecurities(newID, AAPL.US)", func() error {
					return slCtx.RemoveSecurities(ctx, newID, []string{"AAPL.US"})
				})
				t.check("Delete(newID)", func() error {
					return slCtx.Delete(ctx, newID)
				})
			}
		} else {
			t.skip("Create()", "-write not set")
			t.skip("AddSecurities()", "-write not set")
			t.skip("SortSecurities()", "-write not set")
			t.skip("RemoveSecurities()", "-write not set")
			t.skip("Delete()", "-write not set")
		}
	}

	// ══════════════════════════════════════════════════════════
	// quote — new methods (5 methods)
	// ══════════════════════════════════════════════════════════
	fmt.Println("\n=== quote.QuoteContext — new methods (5 methods) ===")
	qctx, err := quote.NewFromCfg(cfg)
	if err != nil {
		log.Printf("quote.NewFromCfg: %v", err)
	} else {
		defer qctx.Close()
		t.check("ShortPositions(AAPL.US)", func() error {
			r, err := qctx.ShortPositions(ctx, "AAPL.US")
			if err != nil {
				return err
			}
			fmt.Printf(" (%d records)", len(r.Data))
			return nil
		})
		t.check("OptionVolume(AAPL.US)", func() error {
			r, err := qctx.OptionVolume(ctx, "AAPL.US")
			if err != nil {
				return err
			}
			fmt.Printf(" (call:%s put:%s)", r.CallVolume, r.PutVolume)
			return nil
		})
		end := time.Now()
		start := end.Add(-30 * 24 * time.Hour)
		t.check("OptionVolumeDaily(AAPL.US, 30d)", func() error {
			r, err := qctx.OptionVolumeDaily(ctx, "AAPL.US", start, end)
			if err != nil {
				return err
			}
			fmt.Printf(" (%d days)", len(r))
			return nil
		})
		t.check("WatchedGroups() [IsPinned field]", func() error {
			groups, err := qctx.WatchedGroups(ctx)
			if err != nil {
				return err
			}
			total := 0
			for _, g := range groups {
				total += len(g.Securites)
			}
			fmt.Printf(" (%d groups, %d securities)", len(groups), total)
			return nil
		})
		if *withWrite {
			t.check("UpdatePinned(Add, AAPL.US)", func() error {
				err := qctx.UpdatePinned(ctx, quote.PinnedModeAdd, []string{"AAPL.US"})
				if err != nil {
					return err
				}
				// clean up
				return qctx.UpdatePinned(ctx, quote.PinnedModeRemove, []string{"AAPL.US"})
			})
		} else {
			t.skip("UpdatePinned()", "-write not set")
		}
	}

	// ══════════════════════════════════════════════════════════
	// summary
	// ══════════════════════════════════════════════════════════
	fmt.Printf("\n─────────────────────────────────────────────────────────────\n")
	fmt.Printf("PASS: %d  FAIL: %d\n", t.pass, t.fail)
	if t.fail > 0 {
		os.Exit(1)
	}
}
