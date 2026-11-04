package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/longbridgeapp/assert"
	"golang.org/x/oauth2"
)

func TestBuildAuthorizeURL_CarriesPKCEChallenge(t *testing.T) {
	o := NewWithBaseURL("my-client", oauthBaseURL)
	authURL, state, verifier := buildAuthorizeURL(o.oauth2Config("http://localhost:60355/callback"))

	u, err := url.Parse(authURL)
	assert.NoError(t, err)
	q := u.Query()
	assert.Equal(t, "S256", q.Get("code_challenge_method"))
	assert.Equal(t, oauth2.S256ChallengeFromVerifier(verifier), q.Get("code_challenge"))
	assert.Equal(t, state, q.Get("state"))
	assert.Equal(t, "my-client", q.Get("client_id"))

	_, _, verifier2 := buildAuthorizeURL(o.oauth2Config("http://localhost:60355/callback"))
	assert.NotEqual(t, verifier, verifier2)
}

func revokeServer(t *testing.T, status int, body string) (*httptest.Server, <-chan url.Values) {
	t.Helper()
	forms := make(chan url.Values, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/oauth2/revoke", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		forms <- r.PostForm
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv, forms
}

func TestRevoke_PrefersRefreshToken(t *testing.T) {
	srv, forms := revokeServer(t, http.StatusOK, "")
	o := NewWithBaseURL("my-client", srv.URL)
	o.token = &oauthToken{AccessToken: "access", RefreshToken: "refresh", ExpiresAt: time.Now().Unix() + 7200}

	assert.NoError(t, o.Revoke(context.Background()))
	form := <-forms
	assert.Equal(t, "refresh", form.Get("token"))
	assert.Equal(t, "refresh_token", form.Get("token_type_hint"))
	assert.Equal(t, "my-client", form.Get("client_id"))
	// Revoke leaves the local token untouched.
	assert.Equal(t, "access", o.token.AccessToken)
}

func TestRevoke_FallsBackToAccessToken(t *testing.T) {
	srv, forms := revokeServer(t, http.StatusOK, "")
	o := NewWithBaseURL("my-client", srv.URL)
	o.token = &oauthToken{AccessToken: "access", ExpiresAt: time.Now().Unix() + 7200}

	assert.NoError(t, o.Revoke(context.Background()))
	form := <-forms
	assert.Equal(t, "access", form.Get("token"))
	assert.Equal(t, "access_token", form.Get("token_type_hint"))
}

func TestRevoke_ServerError(t *testing.T) {
	srv, _ := revokeServer(t, http.StatusBadRequest,
		`{"error":"unsupported_token_type","error_description":"token type not supported"}`)
	o := NewWithBaseURL("my-client", srv.URL)
	o.token = &oauthToken{AccessToken: "access", ExpiresAt: time.Now().Unix() + 7200}

	err := o.Revoke(context.Background())
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "unsupported_token_type: token type not supported")
}

func TestRevoke_WithoutToken(t *testing.T) {
	o := NewWithBaseURL("my-client", "http://127.0.0.1:1")
	err := o.Revoke(context.Background())
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "no token available")
}

func TestRevokeErrorMessage(t *testing.T) {
	assert.Equal(t, "invalid_request", revokeErrorMessage(400, []byte(`{"error":"invalid_request"}`)))
	assert.Equal(t, "server returned HTTP 502", revokeErrorMessage(502, []byte("bad gateway")))
}

func TestClear_DropsMemoryAndFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	o := New("my-client")
	path, err := o.tokenPath()
	assert.NoError(t, err)
	tok := &oauthToken{AccessToken: "access", ExpiresAt: time.Now().Unix() + 7200}
	assert.NoError(t, saveTokenToPath(path, tok))
	o.token = tok

	assert.NoError(t, o.Clear())
	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))
	_, err = o.AccessToken(context.Background())
	assert.True(t, err != nil)

	// Clearing again (missing file) is not an error.
	assert.NoError(t, o.Clear())
}

func TestClear_KeepsMemoryTokenWhenDeleteFails(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	o := New("my-client")
	path, err := o.tokenPath()
	assert.NoError(t, err)
	// A non-empty directory at the token path cannot be removed with os.Remove.
	assert.NoError(t, os.MkdirAll(filepath.Join(path, "x"), 0o700))
	o.token = &oauthToken{AccessToken: "access", ExpiresAt: time.Now().Unix() + 7200}

	assert.True(t, o.Clear() != nil)
	got, err := o.AccessToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "access", got)
}

// tokenServer serves /oauth2/token, blocking each request until release is
// closed, and records each request's form.
func tokenServer(t *testing.T, release <-chan struct{}) (*httptest.Server, <-chan url.Values) {
	t.Helper()
	forms := make(chan url.Values, 4)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		forms <- r.PostForm
		<-release
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"new-access","refresh_token":"new-refresh","token_type":"Bearer","expires_in":7200}`))
	}))
	t.Cleanup(srv.Close)
	return srv, forms
}

func TestAccessToken_RefreshDiscardedAfterClear(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	release := make(chan struct{})
	srv, forms := tokenServer(t, release)
	o := NewWithBaseURL("my-client", srv.URL)
	path, err := o.tokenPath()
	assert.NoError(t, err)
	o.token = &oauthToken{AccessToken: "old", RefreshToken: "refresh", ExpiresAt: time.Now().Unix() - 1}

	errs := make(chan error, 1)
	go func() {
		_, err := o.AccessToken(context.Background())
		errs <- err
	}()
	<-forms // refresh is in flight
	assert.NoError(t, o.Clear())
	close(release)

	err = <-errs
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "cleared")
	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr))
	o.mu.Lock()
	assert.True(t, o.token == nil)
	o.mu.Unlock()
}

func TestBuild_RefreshDiscardedAfterClear(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	release := make(chan struct{})
	srv, forms := tokenServer(t, release)
	o := NewWithBaseURL("my-client", srv.URL)
	path, err := o.tokenPath()
	assert.NoError(t, err)
	assert.NoError(t, saveTokenToPath(path, &oauthToken{
		AccessToken: "old", RefreshToken: "refresh", ExpiresAt: time.Now().Unix() - 1,
	}))

	errs := make(chan error, 1)
	go func() { errs <- o.Build(context.Background()) }()
	<-forms // refresh is in flight
	assert.NoError(t, o.Clear())
	close(release)

	err = <-errs
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "cleared")
	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr))
}

func TestExchange_SendsPKCEVerifier(t *testing.T) {
	release := make(chan struct{})
	close(release)
	srv, forms := tokenServer(t, release)
	cfg := NewWithBaseURL("my-client", srv.URL).oauth2Config("http://localhost:60355/callback")
	_, _, verifier := buildAuthorizeURL(cfg)

	_, err := cfg.Exchange(context.Background(), "code", oauth2.VerifierOption(verifier))
	assert.NoError(t, err)
	form := <-forms
	assert.Equal(t, verifier, form.Get("code_verifier"))
	assert.Equal(t, "code", form.Get("code"))
}

func TestRevoke_WaitsForRefreshInFlight(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	release := make(chan struct{})
	tokenForms := make(chan url.Values, 1)
	revokeForms := make(chan url.Values, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		switch r.URL.Path {
		case "/oauth2/token":
			tokenForms <- r.PostForm
			<-release
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"new-access","refresh_token":"new-refresh","token_type":"Bearer","expires_in":7200}`))
		case "/oauth2/revoke":
			revokeForms <- r.PostForm
		}
	}))
	t.Cleanup(srv.Close)

	o := NewWithBaseURL("my-client", srv.URL)
	o.token = &oauthToken{AccessToken: "old", RefreshToken: "old-refresh", ExpiresAt: time.Now().Unix() - 1}

	accessErrs := make(chan error, 1)
	go func() {
		_, err := o.AccessToken(context.Background())
		accessErrs <- err
	}()
	<-tokenForms // refresh is in flight

	revokeErrs := make(chan error, 1)
	go func() { revokeErrs <- o.Revoke(context.Background()) }()
	select {
	case <-revokeForms:
		t.Fatal("Revoke must wait for the refresh in flight")
	case <-time.After(100 * time.Millisecond):
	}

	close(release)
	assert.NoError(t, <-accessErrs)
	assert.NoError(t, <-revokeErrs)
	form := <-revokeForms
	assert.Equal(t, "new-refresh", form.Get("token"))
}

func TestAccessToken_ConcurrentCallsRefreshOnce(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	release := make(chan struct{})
	srv, forms := tokenServer(t, release)
	o := NewWithBaseURL("my-client", srv.URL)
	o.token = &oauthToken{AccessToken: "old", RefreshToken: "refresh", ExpiresAt: time.Now().Unix() - 1}

	const n = 3
	results := make(chan string, n)
	for i := 0; i < n; i++ {
		go func() {
			tok, err := o.AccessToken(context.Background())
			assert.NoError(t, err)
			results <- tok
		}()
	}
	<-forms
	close(release)
	for i := 0; i < n; i++ {
		assert.Equal(t, "new-access", <-results)
	}
	assert.Equal(t, 0, len(forms))
}
