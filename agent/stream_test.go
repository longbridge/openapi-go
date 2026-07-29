package agent

import (
	"io"
	"strings"
	"testing"
)

// sseFrame formats a raw SSE frame the way the server sends it: an
// "event: message" line (the SSE transport's own event name, always
// "message") followed by a "data:" line carrying the real event envelope.
func sseFrame(data string) string {
	return "event: message\ndata: " + data + "\n\n"
}

// The `data:` payloads of the three example SSE frames from
// https://open.longbridge.com/en/docs/ai/chat/conversation, plus the four
// event types the docs don't mention (workflow_started, ping,
// chat_finished, chat_title_updated).
const (
	chatStartedFrame      = `{"event":"chat_started","workflow_run_id":"wr_1","data":{"chat_uid":"ct_9f2c1a5b","message_id":42}}`
	workflowStartedFrame  = `{"event":"workflow_started","workflow_run_id":"wr_1","data":{"hit_cache":false,"inputs":{"chat_id":834552,"chat_uid":"ct_9f2c1a5b","message_id":42,"query":"How has Tesla stock performed recently?"},"started_at":1784545150,"workflow_id":176476}}`
	messageFrame          = `{"event":"message","workflow_run_id":"wr_1","data":{"text":"Tesla"}}`
	pingFrame             = `{"event":"ping","workflow_run_id":"wr_1","data":null}`
	chatFinishedFrame     = `{"event":"chat_finished","workflow_run_id":"wr_1","data":{"chat_id":834552,"chat_uid":"ct_9f2c1a5b","error":"","error_message":"","message_id":42}}`
	workflowFinishedFrame = `{"event":"workflow_finished","workflow_run_id":"wr_1","data":{"status":"succeeded","elapsed_time":3.21,"outputs":{"answer":"Tesla (TSLA.US) recently..."}}}`
	chatTitleUpdatedFrame = `{"event":"chat_title_updated","workflow_run_id":"wr_1","data":{"chat_id":834552,"chat_uid":"ct_9f2c1a5b","source":"ai_generated","title":"Tesla stock performance","updated_at":1784546957}}`

	// From https://open.longbridge.com/en/docs/ai/chat/events: the interrupt
	// event fired mid-run instead of workflow_finished.
	humanInteractionRequiredFrame = `{"event":"human_interaction_required","workflow_run_id":"wr_1","data":{"node_id":"n_ask_human","tool_call_id":"call_abc123","questions":[{"question":"Which time range would you like to check?","options":[{"description":"Past week"},{"description":"Past month"}],"multi_select":false}],"message_id":43,"chat_id":1001}}`
)

func newTestStream(frames ...string) *ConversationStream {
	body := io.NopCloser(strings.NewReader(strings.Join(frames, "")))
	return newConversationStream(body, nil)
}

func TestConversationStreamFullSequence(t *testing.T) {
	// Exercises the fuller, real-world sequence beyond the docs' three-event
	// example, including chat_title_updated arriving *after*
	// workflow_finished (observed live against the real API).
	stream := newTestStream(
		sseFrame(chatStartedFrame),
		sseFrame(workflowStartedFrame),
		sseFrame(messageFrame),
		sseFrame(pingFrame),
		sseFrame(chatFinishedFrame),
		sseFrame(workflowFinishedFrame),
		sseFrame(chatTitleUpdatedFrame),
	)
	defer stream.Close()

	if !stream.Next() {
		t.Fatalf("expected chat_started event, got err: %v", stream.Err())
	}
	chatStarted, ok := stream.Event().(*ChatStartedEvent)
	if !ok {
		t.Fatalf("expected *ChatStartedEvent, got %T", stream.Event())
	}
	if chatStarted.ChatUID != "ct_9f2c1a5b" || chatStarted.MessageID != "42" {
		t.Errorf("unexpected ChatStartedEvent: %+v", chatStarted)
	}

	if !stream.Next() {
		t.Fatalf("expected workflow_started event, got err: %v", stream.Err())
	}
	workflowStarted, ok := stream.Event().(*WorkflowStartedEvent)
	if !ok {
		t.Fatalf("expected *WorkflowStartedEvent, got %T", stream.Event())
	}
	if workflowStarted.HitCache {
		t.Error("expected HitCache = false")
	}
	if workflowStarted.Inputs.ChatUID != "ct_9f2c1a5b" || workflowStarted.Inputs.MessageID != "42" {
		t.Errorf("unexpected WorkflowStartedEvent.Inputs: %+v", workflowStarted.Inputs)
	}
	if workflowStarted.WorkflowID != 176476 {
		t.Errorf("WorkflowID = %d, want 176476", workflowStarted.WorkflowID)
	}

	if !stream.Next() {
		t.Fatalf("expected message event, got err: %v", stream.Err())
	}
	message, ok := stream.Event().(*MessageEvent)
	if !ok {
		t.Fatalf("expected *MessageEvent, got %T", stream.Event())
	}
	if message.Text != "Tesla" {
		t.Errorf("Text = %q, want %q", message.Text, "Tesla")
	}

	if !stream.Next() {
		t.Fatalf("expected ping event, got err: %v", stream.Err())
	}
	if _, ok := stream.Event().(*PingEvent); !ok {
		t.Fatalf("expected *PingEvent, got %T", stream.Event())
	}

	if !stream.Next() {
		t.Fatalf("expected chat_finished event, got err: %v", stream.Err())
	}
	chatFinished, ok := stream.Event().(*ChatFinishedEvent)
	if !ok {
		t.Fatalf("expected *ChatFinishedEvent, got %T", stream.Event())
	}
	if chatFinished.ChatUID != "ct_9f2c1a5b" || chatFinished.MessageID != "42" {
		t.Errorf("unexpected ChatFinishedEvent: %+v", chatFinished)
	}
	if chatFinished.Error != "" || chatFinished.ErrorMessage != "" {
		t.Errorf("expected empty Error/ErrorMessage, got %+v", chatFinished)
	}

	if !stream.Next() {
		t.Fatalf("expected workflow_finished event, got err: %v", stream.Err())
	}
	workflowFinished, ok := stream.Event().(*WorkflowFinishedEvent)
	if !ok {
		t.Fatalf("expected *WorkflowFinishedEvent, got %T", stream.Event())
	}
	if workflowFinished.ChatUID != "ct_9f2c1a5b" || workflowFinished.MessageID != "42" {
		t.Errorf("expected chat_uid/message_id threaded from chat_started, got %+v", workflowFinished.ConversationResponse)
	}
	if workflowFinished.Status != ConversationStatusSucceeded {
		t.Errorf("Status = %v, want %v", workflowFinished.Status, ConversationStatusSucceeded)
	}
	if workflowFinished.Answer != "Tesla (TSLA.US) recently..." {
		t.Errorf("Answer = %q", workflowFinished.Answer)
	}

	// Arrives *after* workflow_finished in this (real, observed) ordering.
	if !stream.Next() {
		t.Fatalf("expected chat_title_updated event, got err: %v", stream.Err())
	}
	chatTitleUpdated, ok := stream.Event().(*ChatTitleUpdatedEvent)
	if !ok {
		t.Fatalf("expected *ChatTitleUpdatedEvent, got %T", stream.Event())
	}
	if chatTitleUpdated.ChatUID != "ct_9f2c1a5b" || chatTitleUpdated.Source != "ai_generated" || chatTitleUpdated.Title != "Tesla stock performance" {
		t.Errorf("unexpected ChatTitleUpdatedEvent: %+v", chatTitleUpdated)
	}

	if stream.Next() {
		t.Fatalf("expected stream to end, got event %T", stream.Event())
	}
	if err := stream.Err(); err != nil {
		t.Errorf("unexpected error at end of stream: %v", err)
	}
}

func TestConversationStreamInterruptedSequence(t *testing.T) {
	// Per the docs, an interrupted run emits human_interaction_required and
	// then chat_finished directly — no workflow_finished in between.
	stream := newTestStream(
		sseFrame(chatStartedFrame),
		sseFrame(workflowStartedFrame),
		sseFrame(humanInteractionRequiredFrame),
		sseFrame(chatFinishedFrame),
	)
	defer stream.Close()

	if !stream.Next() {
		t.Fatalf("expected chat_started event, got err: %v", stream.Err())
	}
	if _, ok := stream.Event().(*ChatStartedEvent); !ok {
		t.Fatalf("expected *ChatStartedEvent, got %T", stream.Event())
	}

	if !stream.Next() {
		t.Fatalf("expected workflow_started event, got err: %v", stream.Err())
	}
	if _, ok := stream.Event().(*WorkflowStartedEvent); !ok {
		t.Fatalf("expected *WorkflowStartedEvent, got %T", stream.Event())
	}

	if !stream.Next() {
		t.Fatalf("expected human_interaction_required event, got err: %v", stream.Err())
	}
	interrupted, ok := stream.Event().(*HumanInteractionRequiredEvent)
	if !ok {
		t.Fatalf("expected *HumanInteractionRequiredEvent, got %T", stream.Event())
	}
	if interrupted.NodeID != "n_ask_human" || interrupted.ToolCallID != "call_abc123" {
		t.Errorf("unexpected HumanInteractionRequiredEvent: %+v", interrupted.Interrupt)
	}
	if interrupted.MessageID != 43 || interrupted.ChatID != 1001 {
		t.Errorf("unexpected MessageID/ChatID: %+v", interrupted.Interrupt)
	}
	if len(interrupted.Questions) != 1 || interrupted.Questions[0].Question != "Which time range would you like to check?" {
		t.Fatalf("unexpected Questions: %+v", interrupted.Questions)
	}
	if len(interrupted.Questions[0].Options) != 2 || interrupted.Questions[0].MultiSelect {
		t.Errorf("unexpected Questions[0]: %+v", interrupted.Questions[0])
	}

	if !stream.Next() {
		t.Fatalf("expected chat_finished event, got err: %v", stream.Err())
	}
	if _, ok := stream.Event().(*ChatFinishedEvent); !ok {
		t.Fatalf("expected *ChatFinishedEvent, got %T", stream.Event())
	}

	if stream.Next() {
		t.Fatalf("expected stream to end, got event %T", stream.Event())
	}
	if err := stream.Err(); err != nil {
		t.Errorf("unexpected error at end of stream: %v", err)
	}
}

func TestConversationStreamUnknownEventFallsBackToOther(t *testing.T) {
	stream := newTestStream(sseFrame(`{"event":"some_future_event","data":{"foo":"bar"}}`))
	defer stream.Close()

	if !stream.Next() {
		t.Fatalf("expected an event, got err: %v", stream.Err())
	}
	other, ok := stream.Event().(*OtherEvent)
	if !ok {
		t.Fatalf("expected *OtherEvent, got %T", stream.Event())
	}
	if other.Event != "some_future_event" {
		t.Errorf("Event = %q, want %q", other.Event, "some_future_event")
	}
	if !strings.Contains(string(other.Data), `"foo":"bar"`) {
		t.Errorf("Data = %s", other.Data)
	}
}

func TestConversationStreamWithoutTrailingBlankLine(t *testing.T) {
	// The final frame of a real stream may not be followed by a blank line
	// before the connection closes; the reader must still surface it.
	body := io.NopCloser(strings.NewReader("event: message\ndata: " + messageFrame))
	stream := newConversationStream(body, nil)
	defer stream.Close()

	if !stream.Next() {
		t.Fatalf("expected an event, got err: %v", stream.Err())
	}
	if _, ok := stream.Event().(*MessageEvent); !ok {
		t.Fatalf("expected *MessageEvent, got %T", stream.Event())
	}
	if stream.Next() {
		t.Fatalf("expected stream to end, got event %T", stream.Event())
	}
}

func TestContinueStreamSeedsStartedInfo(t *testing.T) {
	// ContinueStream already knows chat_uid/message_id from the caller
	// (unlike a brand-new conversation), so WorkflowFinishedEvent should
	// carry them even without an intervening ChatStartedEvent.
	body := io.NopCloser(strings.NewReader(sseFrame(workflowFinishedFrame)))
	stream := newConversationStream(body, &conversationStarted{ChatUID: "ct_9f2c1a5b", MessageID: "43"})
	defer stream.Close()

	if !stream.Next() {
		t.Fatalf("expected an event, got err: %v", stream.Err())
	}
	workflowFinished, ok := stream.Event().(*WorkflowFinishedEvent)
	if !ok {
		t.Fatalf("expected *WorkflowFinishedEvent, got %T", stream.Event())
	}
	if workflowFinished.ChatUID != "ct_9f2c1a5b" || workflowFinished.MessageID != "43" {
		t.Errorf("unexpected ConversationResponse: %+v", workflowFinished.ConversationResponse)
	}
}
