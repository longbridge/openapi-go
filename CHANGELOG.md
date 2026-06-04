# Changelog

## [Unreleased]

### Added

- New public `counter` package for symbol ↔ `counter_id` conversion, backed by
  an embedded ETF / index / warrant directory:
  - `SymbolToCounterID` / `IndexSymbolToCounterID` / `CounterIDToSymbol`.
  - `IsETF` reports whether a symbol resolves to an ETF.
  - `LookupCounterID` resolves locally only (embedded directory + on-disk cache
    + leading-dot index notation).
  - `CacheCounterIDs` persists remotely resolved entries to
    `$LONGBRIDGE_CACHE_DIR/counter-ids.csv` (default
    `~/.longbridge/cache/counter-ids.csv`), one `counter_id` per line in
    lexicographic order.
- `QuoteContext.SymbolToCounterIds` – batch convert symbols to counter IDs via
  `POST /v1/quote/symbol-to-counter-ids`.
- `QuoteContext.ResolveCounterIds` – local-first counter ID resolution with a
  single batched remote fallback and automatic caching of confirmed entries.
- `FundamentalContext.EtfAssetAllocation` – ETF asset allocation (holdings /
  regional / asset class / industry) via `GET /v1/quote/etf-asset-allocation`.

## [v0.24.2] - 2026-06-02

### Fixed

- `calendar.CalendarEventsResponse`: expose `NextDate` cursor so callers can follow pagination (`/v1/quote/finance_calendar` returns results across multiple pages via `next_date`)
- `calendar.CalendarEventInfo.Symbol`: convert raw `counter_id` (e.g. `ST/US/CRM`) to standard symbol format (`CRM.US`)
- `quote.DailyOptionVolume.Symbol`: convert `underlying_counter_id` to symbol format

### Changed

- Add `internal/counter.IDToSymbol` as shared `counter_id` → symbol conversion helper; `sharelist` migrated to use it

## [v4.2.1] - 2026-05-23

### Changed

- `screener` package: endpoints migrated to `/v1/quote/ai/screener/*`; `ScreenerRecommendStrategies` / `ScreenerUserStrategies` now accept a `market` parameter; `ScreenerSearch` accepts typed `ScreenerCondition` objects (Mode B)

### Fixed

- `OperatingFinancial`: renamed `CounterID` → `Symbol` (converts `ST/US/AAPL` → `AAPL.US`)
- `oauth`: fix panic (`sync: unlock of unlocked mutex`) when authorization flow fails

## [v4.2.0] - 2026-05-22

### Added

- 19 new APIs (same as openapi v4.2.0): `FundamentalContext` +9, `QuoteContext` +1, `MarketContext` +3, new `screener` package +5 — see PR [#91](https://github.com/longbridge/openapi-go/pull/91), [#92](https://github.com/longbridge/openapi-go/pull/92)

### Changed

- `ShortPositions`/`ShortTrades`: typed structs, unified US+HK, RFC 3339 timestamps
- `TopMovers`, `RankList`, `ValuationComparison`: typed structs, `counter_id` → symbol, RFC 3339 timestamps

### Breaking changes

- `StockEvents` → `TopMovers`; `StockEventsResponse` → `TopMoversResponse`
- `HkShortPositions` removed; use `ShortPositions(ctx, symbol, count)`
- Response types for `ShortPositions`, `ShortTrades`, `TopMovers`, `RankList`, `ValuationComparison` changed from raw JSON to typed structs

## [v4.1.0] - 2026-05-14

### Added

- New `alert` package with `AlertContext` for price alert management: `List`, `Add`, `Update`, `Delete`.
- New `calendar` package with `CalendarContext` for the finance calendar: `FinanceCalendar` (earnings, dividends, IPOs, macro data, market closures).
- New `dca` package with `DCAContext` for dollar-cost-averaging plan management: `List`, `Create`, `Update`, `Pause`, `Resume`, `Stop`, `History`, `Stats`, `CheckSupport`, `CalcDate`, `SetReminder`.
- New `fundamental` package with `FundamentalContext` covering financial reports, analyst ratings, dividends, EPS forecasts, consensus estimates, valuation (PE/PB/PS), industry valuation, company overview, executives, shareholders, fund holders, corporate actions, investor relations, operating reports, buyback data, and stock ratings (20 methods).
- New `market` package with `MarketContext` for market-level data: `MarketStatus`, `BrokerHolding`, `BrokerHoldingDetail`, `BrokerHoldingDaily`, `AhPremium`, `AhPremiumIntraday`, `TradeStats`, `Anomaly`, `Constituent`.
- New `portfolio` package with `PortfolioContext` for portfolio analysis: `ExchangeRate`, `ProfitAnalysis`, `ProfitAnalysisByMarket`, `ProfitAnalysisDetail`, `ProfitAnalysisFlows`.
- New `sharelist` package with `SharelistContext` for community sharelist management: `List`, `Detail`, `Popular`, `Create`, `Delete`, `AddSecurities`, `RemoveSecurities`, `SortSecurities`.
- `QuoteContext` gains four new methods: `ShortPositions`, `OptionVolume`, `OptionVolumeDaily`, `UpdatePinned`.
- `WatchedSecurity` gains a new `IsPinned bool` field.
- `Config` gains `ExtraHeaders map[string]string` and `WithHeader(key, value string) *Config` for injecting custom HTTP headers into every request.

### Fixed

- `AlertContext.enable` and `AlertContext.disable` (from prior drafts) replaced by a single `AlertContext.Update(item)` method, matching the v4.1.0 breaking change in the Rust SDK.

## [0.23.0] - 2026-03-30

### Added

- New `asset` package with `StatementContext` for accessing statement APIs:
  - `StatementList` – list account statements with date range and pagination.
  - `StatementDownloadURL` – get the download URL for a specific statement file.
- Staging environment support: set `LONGBRIDGE_ENV=staging` to point to `longbridge.xyz` endpoints (HTTP, quote WebSocket, trade WebSocket, OAuth).
- `ContentContext` adds two new methods:
  - `MyTopics(opts *MyTopicsOptions)` — get topics created by the current authenticated user, with optional page/size/topic_type filtering.
  - `CreateTopic(opts *CreateTopicOptions)` — create a new topic; returns the topic ID (`string`) on success.
- New types: `OwnedTopic`, `MyTopicsOptions`, `CreateTopicOptions`, `TopicReply`, `TopicAuthor`, `TopicImage`.

## [0.22.0] - 2026-03-20

### Breaking changes

- **CN endpoint URLs**: Migrated from `longportapp.cn` to `longbridge.cn` (HTTP, quote WebSocket, trade WebSocket).
- **OAuth token storage path**: Changed from `~/.longbridge-openapi/tokens/` to `~/.longbridge/openapi/tokens/`. Existing tokens under the old path will not be read automatically; move them manually or re-authorize.

## [0.21.0] - 2026-03-19

### Added

- New `content` package with `ContentContext` for accessing content APIs:
  - `Topics` – list discussion topics for a symbol.
  - `News` – list news articles for a symbol.
- `QuoteContext.Filings` – list filing documents for a symbol.

## [0.20.0] - 2025-03-10

### Breaking changes

- **Import path**: Update imports from `github.com/longportapp/openapi-go` to `github.com/longbridge/openapi-go`.
- **Config files**: In TOML/YAML, rename the config section from `[longport]` / `longport:` to `[longbridge]` / `longbridge:`.
- **Environment variables**: The recommended prefix is now `LONGBRIDGE_` (e.g. `LONGBRIDGE_APP_KEY`, `LONGBRIDGE_APP_SECRET`, `LONGBRIDGE_ACCESS_TOKEN`). The old `LONGPORT_` prefix is still supported for backward compatibility.
- **Config API**: `WithOAuth` and `FromOAuth` are removed. Use three keys (app key, secret, access token) or `WithOAuthClient` only.
- **Dependencies**: If you depend on them directly, switch from `github.com/longportapp/openapi-protobufs/gen/go` and `github.com/longportapp/openapi-protocol/go` to `github.com/longbridge/openapi-protobufs/gen/go` (v0.7.0) and `github.com/longbridge/openapi-protocol/go` (v0.5.0).

### Added

- OAuth 2.0 authentication support (`WithOAuthClient`, auto-refresh, authorization code flow).

### Changed

- Module path migrated from `github.com/longportapp/openapi-go` to `github.com/longbridge/openapi-go`.
- Dependencies migrated to Longbridge: `openapi-protobufs/gen/go` v0.7.0, `openapi-protocol/go` v0.5.0.
- Config parsing: `longport` renamed to `longbridge` in `parseConfig` (TOML/YAML config block keys updated accordingly).
- Environment variable prefix: recommended prefix is `LONGBRIDGE_`; `LONGPORT_` remains supported for backward compatibility.
- OAuth flow uses `OnOpenURL` callback for opening the authorization page instead of auto-opening the browser.
- Config validation: only three keys or OAuthClient supported.

### Removed

- Config options: `WithOAuth`, `FromOAuth`.
