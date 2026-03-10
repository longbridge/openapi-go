# Changelog

## [0.20.0] - 2025-03-10

### Breaking changes

- **Import path**: Update imports from `github.com/longportapp/openapi-go` to `github.com/longbridge/openapi-go`.
- **Config files**: In TOML/YAML, rename the config section from `[longport]` / `longport:` to `[longbridge]` / `longbridge:`.
- **Config API**: `WithOAuth` and `FromOAuth` are removed. Use three keys (app key, secret, access token) or `WithOAuthClient` only.
- **Dependencies**: If you depend on them directly, switch from `github.com/longportapp/openapi-protobufs/gen/go` and `github.com/longportapp/openapi-protocol/go` to `github.com/longbridge/openapi-protobufs/gen/go` (v0.7.0) and `github.com/longbridge/openapi-protocol/go` (v0.5.0).

### Added

- OAuth 2.0 authentication support (`WithOAuthClient`, auto-refresh, authorization code flow).

### Changed

- Module path migrated from `github.com/longportapp/openapi-go` to `github.com/longbridge/openapi-go`.
- Dependencies migrated to Longbridge: `openapi-protobufs/gen/go` v0.7.0, `openapi-protocol/go` v0.5.0.
- Config parsing: `longport` renamed to `longbridge` in `parseConfig` (TOML/YAML config block keys updated accordingly).
- OAuth flow uses `OnOpenURL` callback for opening the authorization page instead of auto-opening the browser.
- Config validation: only three keys or OAuthClient supported.

### Removed

- Config options: `WithOAuth`, `FromOAuth`.
