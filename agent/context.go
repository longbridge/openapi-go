package agent

import (
	"context"
	nhttp "net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// requestTimeout bounds blocking Conversation/Continue calls. The shared
// http.Client default (see httplib.DefaultTimeout, 15s) is tuned for fast
// REST calls: in blocking mode the server holds the connection silent until
// the whole LLM turn is done, and that can legitimately take longer.
const requestTimeout = 120 * time.Second

// listTimeout bounds the plain REST calls (Workspaces, Agents). Since
// AgentContext's underlying *http.Client has no overall timeout (see
// NewFromCfg), these would otherwise hang forever on a stuck connection;
// they're ordinary fast REST calls, so the shared default is enough.
const listTimeout = httplib.DefaultTimeout

// AgentContext is a client for the Longbridge AI Agent conversation OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	actx, err := agent.NewFromCfg(conf)
//	resp, err := actx.Conversation(context.Background(), "ag_7d3f9b2c", "How has Tesla stock performed recently?", "")
type AgentContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates an AgentContext from a *config.Config.
//
// Both blocking and streamed conversation calls can legitimately run far
// longer than a typical REST call (a full LLM turn), so unless cfg.Client is
// already set, AgentContext uses its own underlying *http.Client with no
// overall timeout: blocking calls (Conversation, Continue) are instead
// bounded by a dedicated per-request timeout (see requestTimeout), and
// streamed calls (ConversationStream, ContinueStream) are left unbounded —
// only ctx cancellation stops them.
func NewFromCfg(cfg *config.Config) (*AgentContext, error) {
	agentCfg := *cfg
	if agentCfg.Client == nil {
		agentCfg.Client = &nhttp.Client{
			Transport: &nhttp.Transport{IdleConnTimeout: 60 * time.Second},
		}
	}
	httpClient, err := httplib.NewFromCfg(&agentCfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &AgentContext{httpClient: httpClient}, nil
}

// NewFromEnv returns an AgentContext configured from environment variables.
func NewFromEnv() (*AgentContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// Workspaces lists the Workspaces the current account belongs to. A
// Workspace is the organizational unit for Agents: find the target
// Workspace here first, then use Agents to list the Agents available in it.
//
// Path: GET /v1/ai/workspaces
func (c *AgentContext) Workspaces(ctx context.Context) (*WorkspacesResponse, error) {
	var resp WorkspacesResponse
	if err := c.httpClient.Get(ctx, "/v1/ai/workspaces", url.Values{}, &resp, httplib.WithRequestTimeout(listTimeout)); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Agents lists the Agents in the specified Workspace. opts is optional; pass
// nil to use the server defaults (page 1, limit 20, no name filter).
//
// The returned Agent.UID is the identifier used by Conversation; only Agents
// with IsPublished set can start conversations.
//
// Path: GET /v1/ai/workspaces/{id}/agents
func (c *AgentContext) Agents(ctx context.Context, workspaceID string, opts *GetAgentsOptions) (*AgentsResponse, error) {
	q := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			q.Set("page", strconv.FormatInt(int64(opts.Page), 10))
		}
		if opts.Limit > 0 {
			q.Set("limit", strconv.FormatInt(int64(opts.Limit), 10))
		}
		if opts.Name != "" {
			q.Set("name", opts.Name)
		}
	}
	path := "/v1/ai/workspaces/" + url.PathEscape(workspaceID) + "/agents"
	var resp AgentsResponse
	if err := c.httpClient.Get(ctx, path, q, &resp, httplib.WithRequestTimeout(listTimeout)); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Conversation asks a question to the specified Agent, blocking until the
// run succeeds, is interrupted, or fails.
//
// chatUID identifies an existing conversation to continue asking in; pass an
// empty string to start a new one.
//
// The Agent generates the answer using capabilities such as market data and
// account access. When the Agent needs more information or confirmation
// from you, the run is interrupted (Status is ConversationStatusInterrupted)
// — send your answers via Continue to resume it.
//
// Path: POST /v1/ai/agents/{id}/conversations
func (c *AgentContext) Conversation(ctx context.Context, agentID, query, chatUID string) (*ConversationResponse, error) {
	path := "/v1/ai/agents/" + url.PathEscape(agentID) + "/conversations"
	body := conversationBody{Query: query, ChatUID: chatUID}
	var resp ConversationResponse
	if err := c.httpClient.Post(ctx, path, body, &resp, httplib.WithRequestTimeout(requestTimeout)); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Continue resumes an interrupted conversation (a Conversation or Continue
// call that returned Status ConversationStatusInterrupted), blocking until
// the run succeeds, is interrupted again, or fails.
//
// chatUID and messageID come from the interrupted response: chatUID is its
// ChatUID, and messageID is its Interrupt.MessageID (as a string). answers
// is keyed by Interrupt.ToolCallID; see AnswersByToolCall's docs for its
// shape. Every question the interrupt asked must be answered.
//
// A single round of conversation may be interrupted multiple times: if the
// run returns ConversationStatusInterrupted again after continuing, call
// Continue again with the new Interrupt.
//
// Path: POST /v1/ai/agents/{id}/conversations/{chat_uid}/messages/{message_id}/continue
func (c *AgentContext) Continue(ctx context.Context, agentID, chatUID, messageID string, answers AnswersByToolCall) (*ConversationResponse, error) {
	path := continuePath(agentID, chatUID, messageID)
	body := continueBody{AnswersByToolCall: answers}
	var resp ConversationResponse
	if err := c.httpClient.Post(ctx, path, body, &resp, httplib.WithRequestTimeout(requestTimeout)); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ConversationStream starts a conversation with the specified Agent,
// returning a *ConversationStream of run-progress events delivered over SSE
// as they arrive, instead of blocking for the final result the way
// Conversation does.
//
// A ConversationStreamEvent of type *WorkflowFinishedEvent carries the run's
// outcome, but it isn't necessarily the last event — the server may still
// emit a few more housekeeping events (e.g. a *ChatTitleUpdatedEvent) before
// closing the connection, so keep calling Next until it returns false
// rather than stopping as soon as you see it. The caller must Close the
// returned stream once done with it.
//
// chatUID identifies an existing conversation to continue asking in; pass an
// empty string to start a new one.
//
// Path: POST /v1/ai/agents/{id}/conversations (Accept: text/event-stream)
func (c *AgentContext) ConversationStream(ctx context.Context, agentID, query, chatUID string) (*ConversationStream, error) {
	path := "/v1/ai/agents/" + url.PathEscape(agentID) + "/conversations"
	body := conversationBody{Query: query, ChatUID: chatUID}
	rc, err := c.httpClient.CallSSE(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}
	return newConversationStream(rc, nil), nil
}

// ContinueStream resumes an interrupted conversation, returning a
// *ConversationStream of run-progress events delivered over SSE as they
// arrive, instead of blocking for the final result the way Continue does.
// See ConversationStream's docs for how to drain it correctly. The caller
// must Close the returned stream once done with it.
//
// Path: POST /v1/ai/agents/{id}/conversations/{chat_uid}/messages/{message_id}/continue (Accept: text/event-stream)
func (c *AgentContext) ContinueStream(ctx context.Context, agentID, chatUID, messageID string, answers AnswersByToolCall) (*ConversationStream, error) {
	path := continuePath(agentID, chatUID, messageID)
	body := continueBody{AnswersByToolCall: answers}
	rc, err := c.httpClient.CallSSE(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}
	// Unlike a brand-new conversation, chat_uid/message_id are already known
	// here — seed them so the final ConversationResponse carries them even
	// if the server doesn't re-emit a ChatStartedEvent for a continued run.
	started := &conversationStarted{ChatUID: chatUID, MessageID: messageID}
	return newConversationStream(rc, started), nil
}

func continuePath(agentID, chatUID, messageID string) string {
	return "/v1/ai/agents/" + url.PathEscape(agentID) +
		"/conversations/" + url.PathEscape(chatUID) +
		"/messages/" + url.PathEscape(messageID) + "/continue"
}
