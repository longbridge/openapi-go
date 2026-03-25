# Changelog

## [0.23.0] - 2026-03-24

### Added

- `ContentContext.TopicsMine` – list discussion topics created by the current authenticated user, with pagination and type filtering.
- `ContentContext.CreateTopic` – create a new discussion topic for the current authenticated user.
- New types: `OwnedTopic`, `Author`, `Image`, `TopicsMineOptions`, `CreateTopicOptions`.

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
