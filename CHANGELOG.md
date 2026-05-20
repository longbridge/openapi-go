# Changelog

## [Unreleased]

### Added

- **Go:** Six new `FundamentalContext` methods (merged from PR #91):
  - `BusinessSegments` — GET `/v1/quote/fundamentals/business-segments`
  - `BusinessSegmentsHistory` — GET `/v1/quote/fundamentals/business-segments/history`
  - `InstitutionRatingViews` — GET `/v1/quote/ratings/institutional`
  - `IndustryRank` — GET `/v1/quote/industry/rank`
  - `IndustryPeers` — GET `/v1/quote/industries/peers`
  - `FinancialReportSnapshot` — GET `/v1/quote/financials/earnings-snapshot`

- **Go:** New `screener` package with `ScreenerContext` (5 methods):
  - `ScreenerRecommendStrategies` — GET `/v1/quote/screener/strategies/recommend`
  - `ScreenerUserStrategies` — GET `/v1/quote/screener/strategies/mine`
  - `ScreenerStrategy` — GET `/v1/quote/screener/strategy`
  - `ScreenerSearch` — POST `/v1/quote/screener/search`
  - `ScreenerIndicators` — GET `/v1/quote/screener/indicators`

- **Go:** Three new `FundamentalContext` methods:
  - `ShareholderTop` — GET `/v1/quote/shareholders/top`
  - `ShareholderDetail` — GET `/v1/quote/shareholders/holding`
  - `ValuationComparison` — GET `/v1/quote/compare/valuation`

- **Go:** Two new `QuoteContext` methods:
  - `ShortPositions(ctx, symbol, count)` — GET `/v1/quote/short-positions/hk` or `/us` (auto-detected from symbol suffix)
  - `ShortTrades` — GET `/v1/quote/short-trades/hk` or `/us` (auto-detected)

- **Go:** Three new `MarketContext` methods:
  - `TopMovers` — POST `/v1/quote/market/stock-events`
  - `RankCategories` — GET `/v1/quote/market/rank/categories`
  - `RankList` — GET `/v1/quote/market/rank/list`

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
