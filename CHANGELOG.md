# Changelog

## [Unreleased]

### Breaking changes

- **Removed `GridContext.SubmitStrategyQuestionnaire`** (`POST /v1/record/questionnaire`) and its `jsontypes.SubmitStrategyQuestionnaire` request type (ports longbridge/openapi). The strategy risk-disclosure questionnaire endpoint has been retired and is no longer submitted through the SDK. All other grid APIs are unchanged
- **Removed all symbol↔counter_id conversions from the SDK** (ports longbridge/openapi PR #562). Every endpoint now sends and receives the user-facing symbol (e.g. `AAPL.US`, `HSI.HK`) directly:
  - Request query/body parameters `counter_id` / `counter_ids` / `underlying_counter_id` were renamed to `symbol` / `symbols` and carry the user symbol as-is — no more conversion to the internal `ST/US/AAPL` / `IX/HK/HSI` / `ETF/SH/513050` / `VA/BKKT/BTCUSD` forms. Affects `fundamental` (all endpoints incl. the US series and `comparison_symbols` on `ValuationComparison`), `market`, `quote` (`ShortPositions`, `ShortTrades`, `OptionVolume*`, `USCryptoOverview`), `trade` (`USQueryOrders` body `symbols`), `portfolio` (`ProfitAnalysisDetail`, `ProfitAnalysisFlows`), `dca` (`List`, `Create`, `Stats`, `CheckSupport`, `CalcDate`), `sharelist` (`AddSecurities`, `RemoveSecurities`, `SortSecurities`) and `alert`
  - Response fields previously read from `counter_id` and converted are now read from the server's `symbol` field directly (`ExecutiveGroup`, `ShareholderStock`, `FundHolder`, `OperatingFinancial`, `IndustryRankItem`, `IndustryPeerNode`, asset-allocation items, DCA plans/records, sharelist members, alert groups, US orders/assets — `USStockEntry` reads `full_symbol`, US crypto entries read `symbol`)
  - `DELETE` endpoints no longer send a JSON body: `alert.Delete` sends `ids` and `sharelist.RemoveSecurities` sends `symbols` as query-string parameters; `sharelist.Delete` sends no body at all
  - **Removed** the public `counter` package, the internal converter (embedded US-ETF/IX/WT directories) and `quote.QuoteContext.SymbolToCounterIds` / `ResolveCounterIds` (`POST /v1/quote/symbol-to-counter-ids` is no longer called)
- `fundamental.Ratings` (`GET /v1/quote/ratings`) is temporarily disabled — the endpoint is not yet open; the method is commented out pending release
- `market.TopMoversResponse.NextParams` is now a plain pagination-cursor `string` instead of `json.RawMessage`; an empty string means there are no more pages
- Portfolio flows `ExecutedTimestamp` is normalized to a string (server may send int or string; `null` becomes an empty string)
- **`CalendarContext.FinanceCalendar` gains pagination controls** (ports longbridge/openapi #597) — three new optional trailing parameters `count *int32` (max events per page), `offset *int32` (events to skip), and `next *CalendarPageDirection` (a new `Later` / `Earlier` enum). Previously the method issued a single request with no page-size control, so the server's default page cap (historically ~10 events) made even a one-day query look truncated at 10 results. The response's `NextDate` cursor is unchanged; to retrieve a full window, request a larger `count`, or re-call with the returned `NextDate` as `start` until it comes back empty. Pass `nil` for the new arguments to use the server defaults
- **`QuoteContext.OptionChainInfoByDate` moves off the quote socket onto the HTTP endpoint `GET /v1/gemini/option/option_chain_list`** (ports longbridge/openapi PR #588) and its return type changes shape. The paired `StrikePriceInfo` (`Price` + `CallSymbol` + `PutSymbol` + `Standard`) is **removed** in favour of a flat, one-entry-per-contract `OptionChainContract` (`Symbol`, `ExpiryDate`, `StrikePrice`, `Direction`, `OptionType`, `StandardAttr`, `DaysToExpiry`). Calls and puts are no longer paired by strike, so a strike listed on one side only now yields a single entry; filter on `Direction` where you used to read the two symbol fields. The signature gains an `expiryDate time.Time` (now required, was `*time.Time`) and a `standardOnly bool` that filters out legacy corporate-action contracts server-side (omitted from the query when false). Two new enums come with it: `OptionExpiryCycleType` (`Unknown` / `Monthly` / `Weekly` / `Quarterly` — the server's empty `option_type` means a standard monthly option) and `OptionStandardAttr` (`Unknown` / `Normal` / `Old`), plus `OptionDirection` (`Unknown` / `Put` / `Call`). `QuoteContext.OptionChainExpiryDateList` is unchanged

### Added

- **Grid trading** — new `grid.GridContext` for grid-order management: `Submit` / `Replace` / `Cancel` / `Suspend` / `Restart` grid orders, `List` (paged) and `ListByIds`, `Detail` and `TriggerHistory`, `SubmitStrategyQuestionnaire` (strategy risk-disclosure), and `SymbolInfo` (returns `GridSymbolInfo`: name, last price, lot sizes, price-step rules, channel/authorization) — the security info needed to build a grid order
- **Multi-leg option orders** (ports longbridge/openapi #575, #589, #590) — new `TradeContext.SubmitMultiLeg` (`POST /v1/trade/order/multileg`) submits a multi-leg option combination order (vertical spreads, straddles, strangles, collars, calendar spreads, etc.) whose legs are placed together as a single strategy order. Takes `Side`, `OrderType`, `SubmittedQuantity`, `Strategy` (`MultiLegStrategy`), a list of `SubmitMultiLegOrderLeg` (`Symbol` + `RatioQuantity`), and optional `SubmittedPrice` / `Remark` / `ClientRequestId`; returns the order ID like `SubmitOrder`. `Order`, `OrderDetail`, and the `PushOrderChanged` order-changed push gain an optional `MultiLeg *MultiLegInfo` field — present only for multi-leg orders — carrying `Strategy`, `StrategyName`, `MultilegId`, `Code`, and the combination `Legs` (each with `Symbol`, `Side`, `Position`, `RatioQuantity`, `StrikePrice`, `ExpireDate`, `ContractDirection`). New enums `MultiLegStrategy` (`CoveredCall` … `Strangle`, plus `CalendarCallSpread` / `CalendarPutSpread`), `MultiLegPosition` (`LONG` / `SHORT`), and `ContractDirection` (`C` / `P`)

### Changed

- **`TradeContext.HistoryExecutions` now pages through all results** (ports longbridge/openapi #591). The endpoint caps each response at 1000 records; the method walks the `page` query parameter (1-based) until `has_more` is false, deduping by `trade_id` and stopping early if a page adds nothing new (guarding against the gateway ignoring `page`), bounded to 1000 pages. Previously it issued a single request and silently returned only the first page

## [v0.27.0] - 2026-08-14

### Added

- **AI Agent:** new `agent` package (`AgentContext`) for the AI Agent conversation API:
  - `Workspaces` (`GET /v1/ai/workspaces`) and `Agents` (`GET /v1/ai/workspaces/{id}/agents`)
  - `PublicAgents` (`GET /v1/ai/agents`) — lists all publicly available Agents on the platform (the Explore catalog); unlike `Agents` it is not scoped to a Workspace and returns every published, publicly-shared Agent. Takes the same optional `page` / `limit` / `name` parameters and returns the same `AgentsResponse`
  - `Conversation` / `Continue` — blocking `POST .../conversations` and `.../continue` (dedicated 120s request timeout; the plain `Workspaces` / `Agents` GETs use 15s)
  - `ConversationStream` / `ContinueStream` — SSE variants returning a `*ConversationStream` iterator (`Next` / `Event` / `Err` / `Close`), unbounded except for `ctx` cancellation. Events are modeled as an interface with 7 concrete types (`ChatStartedEvent`, `WorkflowStartedEvent`, `MessageEvent`, `PingEvent`, `ChatFinishedEvent`, `WorkflowFinishedEvent`, `ChatTitleUpdatedEvent`) plus an `OtherEvent` fallback for forward compatibility
  - `Conversation` / `ConversationStream` take a `parentMessageID` to attach a follow-up after a specific earlier message — only valid together with a non-empty `chatUID`
  - `ConversationResponse.FurtherQuestions` — the "you might also ask" follow-up suggestions carried in the `workflow_finished` outputs
  - `Reference` captures the full source payload the server sends — `OriginalIndex`, `RefType` (wire `type`), `ID`, and the nested `Content` (`json.RawMessage`); previously `Title` / `URL` came back empty and `source` / `description` / `published_at` / … were lost entirely
  - `ChatStartedEvent` carries `ChatID`, `Error`, and `ErrorMessage`
  - `Interrupt.Interactions` — a slice of `HumanInteraction` (`ToolCallID`, `InterruptID`, `InteractionType`, `ToolName`, `Questions`, and the raw `ToolArgs`) — and `QuestionOption.Label`; a `null` `questions` / `interactions` list from the server is tolerated (decodes to an empty slice)
  - `message_id` fields accept either a JSON string or a raw number, matching the other SDKs' defensive deserialization
- **http:** new `CallSSE` (streaming call, returns the raw response body instead of buffering it) and `WithRequestTimeout` request option, added to support the streaming AI Agent calls
- **Attached order (take-profit / stop-loss) support** for `SubmitOrder` and `ReplaceOrder`:
  - New types: `AttachedOrderType` (`AttachedOrderTypeProfitTaker` / `AttachedOrderTypeStopLoss` / `AttachedOrderTypeBracket`), `AttachedOrderDetail`, `SubmitAttachedParams`, `ReplaceAttachedParams`
  - `SubmitOrder` / `ReplaceOrder`: new `AttachedParams` field
  - `Order` / `OrderDetail`: new `AttachedOrders []AttachedOrderDetail` field
  - New `TradeContext.OrderDetailAttached(orderId)` and `TradeContext.CancelOrderAttached(orderId)` methods — query/cancel an attached sub-order by its own order ID
  - `GetTodayOrders`: new `OrderId` and `IsAttached` fields — when combined, treats `OrderId` as an attached sub-order ID and returns that sub-order as an `Order` entry (not the parent order)

### Breaking changes

- `OrderDetail.ChargeDetail` is now `*OrderChargeDetail` (previously non-pointer). Attached orders return `nil` for this field; callers must handle the absent case.

## [v0.26.0] - 2026-07-20

### Added

- **US market APIs** — 14 new interfaces for US-region accounts (requires `us_` token):
  - Fundamental: `CompanyOverview`, `ValuationOverview`, `FinancialOverview`, `FinancialStatement`, `KeyFinancialMetrics`, `AnalystConsensus`, `ETFDividendInfo`, `CompanyDividends`, `ETFFiles`
  - Quote: `CryptoOverview` (BTCUSD.BKKT etc.)
  - Trade: `USAssetOverview`, `USRealizedPL`, `QueryUSOrders`, `USOrderDetail`
- **DC-region routing** — `x-dc-region` header derived from token prefix (`us_` → US, others → AP); `NewHTTPFromCfg` for HTTP-only trade context without WebSocket
- **`AllExecutions`** — `GET /v3/trade/execution/all` queries today and historical executions in one call with pagination
- **`OutsideRTH.OptionPreMarket`** — new enum variant for overnight option orders

### Changed

- `OrderTag` enum values updated to match current API definitions; unused variants removed
- `QuoteContext.RealtimeQuote` skips nil cache entries instead of panicking

### Fixed

- `QuoteContext.RealtimeQuote`: nil pointer dereference when cache entry is empty (#104)

## [v0.25.2] - 2026-06-26

### Added

- `market.TradeStatus` models `/v1/quote/market-status` trade status codes,
  including display names, labels, normalization helpers, and status code `2001`.

### Changed

- `market.MarketTimeItem.TradeStatus` and `DelayTradeStatus` now use
  `market.TradeStatus` instead of raw `int32` values.

### Fixed

- `AlertContext.List`: `AlertSymbolGroup.Symbol` was always empty because the JSON tag was `"symbol"` but the API returns `"counter_id"` (e.g. `ST/HK/700`). Fixed by changing the tag to `json:"counter_id"` and converting via `counter.IDToSymbol()` so callers receive the standard symbol format (e.g. `700.HK`).
- `FundamentalContext.Macroeconomic`: `Info.Periodicity` and `Info.Importance` were always zero/empty. Fixed by adding `frequence` and `importance` fields to the v2 wire type `V2MacroeconomicDetail`.

## [v0.25.1] - 2026-06-13

### Added

- **All languages:** `macroeconomic_indicators` gains `keyword` parameter for fuzzy name filtering
- **All languages:** `macroeconomic` switches to `GET /v2/quote/macrodata/{id}`, defaults to `sort=desc`

### Changed

- `MacroeconomicIndicator.name` / `.describe`: `MultiLanguageText` → `string`
- `Macroeconomic.unit` / `.unit_prefix`: `MultiLanguageText` → `string`

## [v0.25.0] - 2026-06-10

### Added

- New public `counter` package for symbol ↔ `counter_id` conversion, backed by
  an embedded ETF (7250) / index (648) / warrant (17693) directory:
  - `SymbolToCounterID` / `IndexSymbolToCounterID` / `CounterIDToSymbol`
  - `IsETF` reports whether a symbol resolves to an ETF
  - `LookupCounterID` resolves locally only (embedded directory + on-disk cache + leading-dot index notation)
  - `CacheCounterIDs` persists remotely resolved entries to `$LONGBRIDGE_CACHE_DIR/counter-ids.csv` (default `~/.longbridge/cache/counter-ids.csv`)
- `QuoteContext.SymbolToCounterIds` — batch convert symbols to counter IDs via `POST /v1/quote/symbol-to-counter-ids`
- `QuoteContext.ResolveCounterIds` — local-first counter ID resolution with batched remote fallback and automatic caching
- `FundamentalContext.EtfAssetAllocation` — ETF asset allocation (holdings / regional / asset class / industry) via `GET /v1/quote/etf-asset-allocation`
- `FundamentalContext.MacroeconomicIndicators(country, offset, limit)` — list macroeconomic indicators via `GET /v1/quote/macrodata`; filter by `MacroeconomicCountry` (HK/CN/US/EU/JP/SG); response includes `Count`
- `FundamentalContext.Macroeconomic(indicatorCode, startDate, endDate, offset, limit)` — historical data for a specific indicator via `GET /v1/quote/macrodata/{indicator_code}`; `startDate` / `endDate` accept `"YYYY-MM-DD"` strings; response includes `Count`
- New types: `MultiLanguageText`, `MacroeconomicCountry`, `MacroeconomicImportance`, `MacroeconomicIndicator`, `MacroeconomicIndicatorListResponse`, `Macroeconomic`, `MacroeconomicResponse`

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
