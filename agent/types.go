// Package agent provides a client for the Longbridge AI Agent conversation
// OpenAPI. It covers listing Workspaces and Agents, starting and continuing
// conversations, and streaming run progress over SSE.
//
// Reference: https://open.longbridge.com/en/docs/ai/chat/conversation
package agent

import (
	"bytes"
	"encoding/json"
)

// ─── Workspaces ───────────────────────────────────────────────────────────

// Workspace is a Workspace the current account belongs to.
type Workspace struct {
	// ID is the Workspace ID.
	ID string `json:"id"`
	// Name is the Workspace name.
	Name string `json:"name"`
	// CreatedAt is the creation time, Unix timestamp in seconds.
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt is the last updated time, Unix timestamp in seconds.
	UpdatedAt int64 `json:"updated_at"`
}

// WorkspacesResponse is the response for AgentContext.Workspaces.
type WorkspacesResponse struct {
	// Workspaces the current account belongs to.
	Workspaces []Workspace `json:"workspaces"`
}

// ─── Agents ───────────────────────────────────────────────────────────────

// Agent is an Agent in a Workspace.
type Agent struct {
	// UID is the Agent UID, used as the path parameter of
	// AgentContext.Conversation.
	UID string `json:"uid"`
	// Name is the Agent name.
	Name string `json:"name"`
	// Description is the Agent description.
	Description string `json:"description"`
	// Mode is the Agent mode, e.g. "chat".
	Mode string `json:"mode"`
	// Icon is the icon URL.
	Icon string `json:"icon"`
	// IsPublished reports whether the Agent is published; only published
	// Agents can start conversations.
	IsPublished bool `json:"is_published"`
	// PublishedAt is the publish time, Unix timestamp in seconds; 0 if
	// unpublished.
	PublishedAt int64 `json:"published_at"`
	// CreatedAt is the creation time, Unix timestamp in seconds.
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt is the last updated time, Unix timestamp in seconds.
	UpdatedAt int64 `json:"updated_at"`
}

// AgentsResponse is the response for AgentContext.Agents.
type AgentsResponse struct {
	// Agents is the Agent list.
	Agents []Agent `json:"agents"`
	// Total is the total number of matching Agents.
	Total int32 `json:"total"`
}

// GetAgentsOptions holds the optional query parameters for
// AgentContext.Agents. The zero value uses the server defaults (page 1,
// limit 20, no name filter).
type GetAgentsOptions struct {
	// Page is the page number, starting at 1. 0 uses the server default.
	Page int32
	// Limit is the page size. 0 uses the server default.
	Limit int32
	// Name fuzzy-searches by Agent name. Empty omits the filter.
	Name string
}

// ─── Conversation ─────────────────────────────────────────────────────────

// ConversationStatus is the final run status of a conversation.
type ConversationStatus string

const (
	// ConversationStatusSucceeded means the run completed successfully.
	ConversationStatusSucceeded ConversationStatus = "succeeded"
	// ConversationStatusInterrupted means the run is paused, waiting for
	// AgentContext.Continue.
	ConversationStatusInterrupted ConversationStatus = "interrupted"
	// ConversationStatusFailed means the run failed.
	ConversationStatusFailed ConversationStatus = "failed"
	// ConversationStatusStopped means the run was stopped.
	ConversationStatusStopped ConversationStatus = "stopped"
)

// Reference is a source referenced by the answer.
type Reference struct {
	// Index is the reference index.
	Index int32 `json:"index"`
	// Title is the reference title.
	Title string `json:"title"`
	// URL is the reference URL.
	URL string `json:"url"`
}

// QuestionOption is one option of a Question.
type QuestionOption struct {
	// Description is the option text.
	Description string `json:"description"`
}

// Question is one question the Agent needs you to answer to resume an
// interrupted conversation.
type Question struct {
	// Question is the question text.
	Question string `json:"question"`
	// Options are the answer options; empty means free-form answer.
	Options []QuestionOption `json:"options"`
	// MultiSelect reports whether multiple options may be selected.
	MultiSelect bool `json:"multi_select"`
}

// Interrupt is present when a conversation run is interrupted, waiting for
// AgentContext.Continue.
type Interrupt struct {
	// NodeID is the ID of the node that triggered the interrupt.
	NodeID string `json:"node_id"`
	// ToolCallID is the tool call ID of this inquiry; used as the outer key
	// of AnswersByToolCall when continuing.
	ToolCallID string `json:"tool_call_id"`
	// Questions you need to answer.
	Questions []Question `json:"questions"`
	// MessageID is the ID of the paused message.
	MessageID int64 `json:"message_id"`
	// ChatID is the ID of the owning conversation.
	ChatID int64 `json:"chat_id"`
}

// ConversationError is present when a conversation run failed.
type ConversationError struct {
	// Code is the error code.
	Code int32 `json:"code"`
	// Message is the error message.
	Message string `json:"message"`
}

// AnswersByToolCall is answers keyed by Interrupt.ToolCallID, each value
// being a map of question text to answer, used as the request body of
// AgentContext.Continue and AgentContext.ContinueStream.
//
// The outer key corresponds to Interrupt.ToolCallID. The inner key is the
// Question.Question text, and the value is the answer selected or typed for
// it. Every question the interrupt asked must be answered.
type AnswersByToolCall map[string]map[string]string

// ConversationResponse is the response for AgentContext.Conversation and
// AgentContext.Continue, and the outcome carried by WorkflowFinishedEvent
// when streaming.
type ConversationResponse struct {
	// ChatUID identifies the conversation, used for follow-up questions and
	// troubleshooting.
	ChatUID string `json:"chat_uid"`
	// MessageID is the message ID of this round.
	MessageID string `json:"message_id"`
	// Status is the final run status.
	Status ConversationStatus `json:"status"`
	// Answer is the final answer text; valid when Status is
	// ConversationStatusSucceeded.
	Answer string `json:"answer"`
	// References are sources referenced by the answer; nil if none.
	References []Reference `json:"references"`
	// ElapsedTime is the run duration in seconds.
	ElapsedTime float64 `json:"elapsed_time"`
	// Interrupt is present only when Status is ConversationStatusInterrupted.
	Interrupt *Interrupt `json:"interrupt"`
	// Error is present only when the run failed.
	Error *ConversationError `json:"error"`
}

// UnmarshalJSON accepts message_id as either a JSON string or a JSON
// number: the blocking response's top-level message_id is a quoted string,
// but callers building a ConversationResponse from streamed events may feed
// it the numeric form seen in some SSE payloads (see ChatStartedEvent).
func (r *ConversationResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		ChatUID     string             `json:"chat_uid"`
		MessageID   json.RawMessage    `json:"message_id"`
		Status      ConversationStatus `json:"status"`
		Answer      string             `json:"answer"`
		References  []Reference        `json:"references"`
		ElapsedTime float64            `json:"elapsed_time"`
		Interrupt   *Interrupt         `json:"interrupt"`
		Error       *ConversationError `json:"error"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	messageID, err := unmarshalStringOrInt(raw.MessageID)
	if err != nil {
		return err
	}
	r.ChatUID = raw.ChatUID
	r.MessageID = messageID
	r.Status = raw.Status
	r.Answer = raw.Answer
	r.References = raw.References
	r.ElapsedTime = raw.ElapsedTime
	r.Interrupt = raw.Interrupt
	r.Error = raw.Error
	return nil
}

// conversationBody is the request body of AgentContext.Conversation and
// AgentContext.ConversationStream.
type conversationBody struct {
	Query   string `json:"query"`
	ChatUID string `json:"chat_uid,omitempty"`
}

// continueBody is the request body of AgentContext.Continue and
// AgentContext.ContinueStream.
type continueBody struct {
	AnswersByToolCall AnswersByToolCall `json:"answers_by_tool_call"`
}

// ─── Streaming events ─────────────────────────────────────────────────────

// ConversationStreamEvent is implemented by every event observed while
// draining a *ConversationStream. Use a type switch to handle it:
//
//	switch e := stream.Event().(type) {
//	case *agent.MessageEvent:
//		fmt.Print(e.Text)
//	case *agent.WorkflowFinishedEvent:
//		result = e.ConversationResponse
//	case *agent.HumanInteractionRequiredEvent:
//		interrupt = e.Interrupt
//	}
type ConversationStreamEvent interface {
	conversationStreamEvent()
}

// ChatStartedEvent reports that the run has started.
type ChatStartedEvent struct {
	// ChatUID identifies the conversation.
	ChatUID string `json:"chat_uid"`
	// MessageID is the message ID of this round.
	MessageID string `json:"message_id"`
}

func (*ChatStartedEvent) conversationStreamEvent() {}

// UnmarshalJSON accepts message_id as either a JSON string or a JSON
// number: the docs' SSE example encodes it as a raw number here, unlike the
// blocking response's top-level message_id, which is a quoted string.
func (e *ChatStartedEvent) UnmarshalJSON(data []byte) error {
	var raw struct {
		ChatUID   string          `json:"chat_uid"`
		MessageID json.RawMessage `json:"message_id"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	messageID, err := unmarshalStringOrInt(raw.MessageID)
	if err != nil {
		return err
	}
	e.ChatUID = raw.ChatUID
	e.MessageID = messageID
	return nil
}

// WorkflowStartedInputs echoes the run's inputs, carried by
// WorkflowStartedEvent.
type WorkflowStartedInputs struct {
	// ChatID is the ID of the owning conversation.
	ChatID int64 `json:"chat_id"`
	// ChatUID identifies the conversation.
	ChatUID string `json:"chat_uid"`
	// MessageID is the message ID of this round.
	MessageID string `json:"message_id"`
	// Query is the question that was asked.
	Query string `json:"query"`
}

// UnmarshalJSON accepts message_id as either a JSON string or a JSON
// number; see ChatStartedEvent.UnmarshalJSON.
func (i *WorkflowStartedInputs) UnmarshalJSON(data []byte) error {
	var raw struct {
		ChatID    int64           `json:"chat_id"`
		ChatUID   string          `json:"chat_uid"`
		MessageID json.RawMessage `json:"message_id"`
		Query     string          `json:"query"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	messageID, err := unmarshalStringOrInt(raw.MessageID)
	if err != nil {
		return err
	}
	i.ChatID = raw.ChatID
	i.ChatUID = raw.ChatUID
	i.MessageID = messageID
	i.Query = raw.Query
	return nil
}

// WorkflowStartedEvent is observed right after ChatStartedEvent on every run
// seen so far.
type WorkflowStartedEvent struct {
	// HitCache reports whether this run's answer was served from a cache.
	HitCache bool `json:"hit_cache"`
	// Inputs echoes the run's inputs.
	Inputs WorkflowStartedInputs `json:"inputs"`
	// StartedAt is a Unix timestamp in seconds.
	StartedAt int64 `json:"started_at"`
	// WorkflowID is an internal workflow run ID.
	WorkflowID int64 `json:"workflow_id"`
}

func (*WorkflowStartedEvent) conversationStreamEvent() {}

// MessageEvent is an incremental piece of the answer.
type MessageEvent struct {
	// Text is the incremental answer text.
	Text string `json:"text"`
}

func (*MessageEvent) conversationStreamEvent() {}

// PingEvent is a heartbeat with no payload, observed at arbitrary points in
// the stream (including in between MessageEvent chunks).
type PingEvent struct{}

func (*PingEvent) conversationStreamEvent() {}

// ChatFinishedEvent is observed once all MessageEvents for this round have
// been sent, shortly before WorkflowFinishedEvent.
type ChatFinishedEvent struct {
	// ChatID is the ID of the owning conversation.
	ChatID int64 `json:"chat_id"`
	// ChatUID identifies the conversation.
	ChatUID string `json:"chat_uid"`
	// MessageID is the message ID of this round.
	MessageID string `json:"message_id"`
	// Error has been empty in every run observed so far.
	Error string `json:"error"`
	// ErrorMessage has been empty in every run observed so far.
	ErrorMessage string `json:"error_message"`
}

func (*ChatFinishedEvent) conversationStreamEvent() {}

// UnmarshalJSON accepts message_id as either a JSON string or a JSON
// number; see ChatStartedEvent.UnmarshalJSON.
func (e *ChatFinishedEvent) UnmarshalJSON(data []byte) error {
	var raw struct {
		ChatID       int64           `json:"chat_id"`
		ChatUID      string          `json:"chat_uid"`
		MessageID    json.RawMessage `json:"message_id"`
		Error        string          `json:"error"`
		ErrorMessage string          `json:"error_message"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	messageID, err := unmarshalStringOrInt(raw.MessageID)
	if err != nil {
		return err
	}
	e.ChatID = raw.ChatID
	e.ChatUID = raw.ChatUID
	e.MessageID = messageID
	e.Error = raw.Error
	e.ErrorMessage = raw.ErrorMessage
	return nil
}

// WorkflowFinishedEvent reports that the run finished (succeeded, failed, or
// stopped), carrying the run's outcome. Not emitted when the run is
// interrupted instead — see HumanInteractionRequiredEvent. Not necessarily
// the last event of the stream — the server may still emit a few more
// housekeeping events (e.g. ChatTitleUpdatedEvent) before actually closing
// the connection, so keep draining the stream until Next returns false
// rather than stopping as soon as this is seen.
type WorkflowFinishedEvent struct {
	*ConversationResponse
}

func (*WorkflowFinishedEvent) conversationStreamEvent() {}

// HumanInteractionRequiredEvent reports that the run is paused, waiting for
// more information from the caller — send answers via AgentContext.Continue
// or AgentContext.ContinueStream to resume it. An interrupted run does not
// emit WorkflowFinishedEvent; this event carries the interrupt details
// instead, and the stream ends with ChatFinishedEvent right after it.
type HumanInteractionRequiredEvent struct {
	*Interrupt
}

func (*HumanInteractionRequiredEvent) conversationStreamEvent() {}

// ChatTitleUpdatedEvent reports the server auto-generating a short title for
// the conversation, as a UI convenience. Can arrive before or after
// WorkflowFinishedEvent; not tied to the run's outcome.
type ChatTitleUpdatedEvent struct {
	// ChatID is the ID of the owning conversation.
	ChatID int64 `json:"chat_id"`
	// ChatUID identifies the conversation.
	ChatUID string `json:"chat_uid"`
	// Source describes where the title came from, e.g. "ai_generated".
	Source string `json:"source"`
	// Title is the new (possibly truncated) title.
	Title string `json:"title"`
	// UpdatedAt is a Unix timestamp in seconds.
	UpdatedAt int64 `json:"updated_at"`
}

func (*ChatTitleUpdatedEvent) conversationStreamEvent() {}

// OtherEvent is an event type not recognized by this SDK version, carried as
// raw JSON so callers aren't broken by future additions to the API. Event
// is the SSE envelope's discriminator string, so callers can at least tell
// these apart instead of getting an opaque blob.
type OtherEvent struct {
	// Event is the SSE envelope's "event" field (the event type name).
	Event string
	// Data is the SSE envelope's "data" field.
	Data json.RawMessage
}

func (*OtherEvent) conversationStreamEvent() {}

// workflowOutputs are the outputs of a workflow_finished SSE event.
type workflowOutputs struct {
	// Answer is the final answer text; present when the run succeeded.
	Answer *string `json:"answer"`
	// References are sources referenced by the answer.
	References []Reference `json:"references"`
	// Interrupt is present only when the status is "interrupted".
	Interrupt *Interrupt `json:"interrupt"`
	// Error is present only when the run failed.
	Error *ConversationError `json:"error"`
}

// workflowFinishedPayload is the payload of a workflow_finished SSE event.
type workflowFinishedPayload struct {
	Status      ConversationStatus `json:"status"`
	ElapsedTime float64            `json:"elapsed_time"`
	Outputs     workflowOutputs    `json:"outputs"`
}

// unmarshalStringOrInt decodes data as either a JSON string or a JSON
// number, returning its value as a string either way. Some SSE event
// payloads encode message_id as a raw JSON number, unlike the blocking
// response's top-level message_id, which is a quoted string.
func unmarshalStringOrInt(data []byte) (string, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return "", nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return "", err
		}
		return s, nil
	}
	return string(data), nil
}
