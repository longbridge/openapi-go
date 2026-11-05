package oauth

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/longbridgeapp/assert"
)

func TestCallbackErrorMessage(t *testing.T) {
	assert.Equal(t, "internal_server_error: validate client err:oauth client not found",
		callbackErrorMessage("internal_server_error", "validate client err:oauth client not found"))
	assert.Equal(t, "access_denied", callbackErrorMessage("access_denied", ""))
	assert.Equal(t, "bad [31m: line1 line2", callbackErrorMessage("bad\x1b[31m", "line1\nline2"))
	assert.Equal(t, "a b c d", callbackErrorMessage("a\u202eb\u2028c\u200bd", ""))

	long := callbackErrorMessage(strings.Repeat("你", maxCallbackErrorLen+10), "")
	assert.Equal(t, maxCallbackErrorLen+1, len([]rune(long)))
	assert.True(t, strings.HasSuffix(long, "…"))
}

// callbackURL is redirectURI with params, aimed at 127.0.0.1.
func callbackURL(redirectURI string, params url.Values) string {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return ""
	}
	u.Host = strings.Replace(u.Host, "localhost", "127.0.0.1", 1)
	u.RawQuery = params.Encode()
	return u.String()
}

// callbackRequest sends a request to the local callback server and returns
// the status code and body. It runs in the flow's driver goroutine, so it
// reports failures with t.Errorf (never FailNow) and returns status 0.
func callbackRequest(t *testing.T, method, redirectURI string, params url.Values) (int, string) {
	t.Helper()
	req, err := http.NewRequest(method, callbackURL(redirectURI, params), nil)
	if err != nil {
		t.Errorf("building callback request: %v", err)
		return 0, ""
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("callback request failed: %v", err)
		return 0, ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

// startFlow runs authorizeFlow, calls drive with the redirect URI and the
// expected state once the authorization URL is opened, and waits for drive
// to finish before returning.
func startFlow(t *testing.T, o *OAuth, drive func(redirectURI, state string)) (*oauthToken, error) {
	t.Helper()
	done := make(chan struct{})
	o.WithCallbackPort(0).OnOpenURL(func(authURL string) {
		u, err := url.Parse(authURL)
		if err != nil {
			t.Errorf("parsing authorization URL: %v", err)
			close(done)
			return
		}
		q := u.Query()
		// The server starts serving after openURL returns.
		go func() {
			defer close(done)
			drive(q.Get("redirect_uri"), q.Get("state"))
		}()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	tok, err := o.authorizeFlow(ctx)
	<-done
	return tok, err
}

func TestAuthorizeFlow_IgnoresCallbacksWithAMismatchedState(t *testing.T) {
	type response struct {
		status int
		body   string
	}
	responses := make(chan response, 2)
	o := NewWithBaseURL("my-client", "http://127.0.0.1:1")
	_, err := startFlow(t, o, func(redirectURI, state string) {
		forged := state + "-forged"
		s, b := callbackRequest(t, http.MethodGet, redirectURI, url.Values{"code": {"code"}, "state": {forged}})
		responses <- response{s, b}
		s, b = callbackRequest(t, http.MethodGet, redirectURI, url.Values{
			"error": {"access_denied"}, "error_description": {"visit attacker.example to fix"}, "state": {forged},
		})
		responses <- response{s, b}
		// The real redirect still completes the flow.
		callbackRequest(t, http.MethodGet, redirectURI, url.Values{"error": {"access_denied"}, "state": {state}})
	})

	assert.True(t, err != nil)
	assert.Equal(t, "oauth: authorization failed: access_denied", err.Error())
	for i := 0; i < 2; i++ {
		r := <-responses
		assert.Equal(t, http.StatusBadRequest, r.status)
		assert.Contains(t, r.body, "earlier or different")
		assert.False(t, strings.Contains(r.body, "attacker.example"))
	}
}

func TestAuthorizeFlow_ReportsErrorDescriptionAndEscapesThePage(t *testing.T) {
	bodies := make(chan string, 1)
	o := NewWithBaseURL("my-client", "http://127.0.0.1:1")
	_, err := startFlow(t, o, func(redirectURI, state string) {
		_, body := callbackRequest(t, http.MethodGet, redirectURI, url.Values{
			"error":             {"internal_server_error"},
			"error_description": {"client not found <script>alert(1)</script>"},
			"state":             {state},
		})
		bodies <- body
	})

	assert.True(t, err != nil)
	assert.Equal(t,
		"oauth: authorization failed: internal_server_error: client not found <script>alert(1)</script>",
		err.Error())
	body := <-bodies
	assert.Contains(t, body, "client not found &lt;script&gt;alert(1)&lt;/script&gt;")
	assert.False(t, strings.Contains(body, "<script>"))
}

func TestAuthorizeFlow_ValidCallbackExchangesTheCode(t *testing.T) {
	release := make(chan struct{})
	close(release)
	srv, forms := tokenServer(t, release)
	o := NewWithBaseURL("my-client", srv.URL)
	tok, err := startFlow(t, o, func(redirectURI, state string) {
		callbackRequest(t, http.MethodGet, redirectURI, url.Values{"code": {"the-code"}, "state": {state}})
	})

	assert.NoError(t, err)
	assert.Equal(t, "new-access", tok.AccessToken)
	form := <-forms
	assert.Equal(t, "the-code", form.Get("code"))
	assert.True(t, form.Get("code_verifier") != "")
}

func TestDeliverFirst_KeepsTheFirstAndNeverBlocks(t *testing.T) {
	ch := make(chan int, 1)
	assert.True(t, deliverFirst(ch, 1))
	done := make(chan bool, 1)
	go func() { done <- deliverFirst(ch, 2) }()
	select {
	case sent := <-done:
		assert.False(t, sent)
	case <-time.After(2 * time.Second):
		t.Fatal("deliverFirst blocked on a full channel")
	}
	assert.Equal(t, 1, <-ch)
}

func TestAuthorizeFlow_CallbackAcceptsOnlyGetAndClosesAfterTheFlow(t *testing.T) {
	var callback string
	o := NewWithBaseURL("my-client", "http://127.0.0.1:1")
	_, err := startFlow(t, o, func(redirectURI, state string) {
		callback = callbackURL(redirectURI, url.Values{"error": {"access_denied"}, "state": {state}})
		status, _ := callbackRequest(t, http.MethodPost, redirectURI,
			url.Values{"error": {"access_denied"}, "state": {state}})
		assert.Equal(t, http.StatusMethodNotAllowed, status)
		// A POST does not end the flow; the GET does.
		callbackRequest(t, http.MethodGet, redirectURI, url.Values{"error": {"access_denied"}, "state": {state}})
	})
	assert.Equal(t, "oauth: authorization failed: access_denied", err.Error())

	// The server (keep-alive connections included) is closed with the flow.
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(callback)
	if err == nil {
		resp.Body.Close()
		t.Fatal("callback server still answering after the flow ended")
	}
}
