package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// AccountChannelCtx returns the account channel (e.g. "lb") of the current
// session, read from the access token's JWT payload without verifying it. It
// returns "" when the token is not a JWT or carries no channel; an error is
// returned only when an OAuth session cannot obtain its token.
func (c *Client) AccountChannelCtx(ctx context.Context) (string, error) {
	token := c.opts.AccessToken
	if c.opts.OAuthClient != nil {
		t, err := c.opts.OAuthClient.AccessToken(ctx)
		if err != nil {
			return "", err
		}
		token = t
	}
	return accountChannelFromToken(token), nil
}

// accountChannelFromToken reads the account channel from a JWT access token:
// the "ac" claim, or "account_channel" inside the "sub" claim (a nested JSON
// string), or a top-level "account_channel" claim. Region prefixes such as
// "us_m_" sit before the first '.', so they do not affect the payload.
func accountChannelFromToken(token string) string {
	parts := strings.Split(strings.TrimSpace(stripBearerPrefix(token)), ".")
	if len(parts) != 3 {
		return ""
	}
	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(parts[1], "="))
	if err != nil {
		return ""
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}
	if ch := claimString(claims, "ac"); ch != "" {
		return ch
	}
	if sub := claimString(claims, "sub"); sub != "" {
		var subject map[string]interface{}
		if err := json.Unmarshal([]byte(sub), &subject); err == nil {
			if ch := claimString(subject, "account_channel"); ch != "" {
				return ch
			}
		}
	}
	return claimString(claims, "account_channel")
}

// claimString returns claims[key] when it is a non-empty string.
func claimString(claims map[string]interface{}, key string) string {
	s, _ := claims[key].(string)
	return strings.TrimSpace(s)
}
