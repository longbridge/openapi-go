package agent

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

// conversationStarted holds the chat_uid/message_id captured from an
// earlier ChatStartedEvent (nil if one hasn't been observed yet), threaded
// through to WorkflowFinishedEvent, which doesn't repeat them.
type conversationStarted struct {
	ChatUID   string
	MessageID string
}

// ConversationStream is an open server-sent-events stream of
// ConversationStreamEvent values from AgentContext.ConversationStream or
// AgentContext.ContinueStream.
//
// Call Next repeatedly to advance through the stream; it returns false once
// the stream ends or an error occurs — call Err to tell the two apart. Call
// Close once done, whether or not the stream was drained to the end.
//
// Usage:
//
//	stream, err := actx.ConversationStream(ctx, agentID, query, "")
//	if err != nil {
//		// handle err
//	}
//	defer stream.Close()
//	for stream.Next() {
//		switch e := stream.Event().(type) {
//		case *agent.MessageEvent:
//			fmt.Print(e.Text)
//		case *agent.WorkflowFinishedEvent:
//			result = e.ConversationResponse
//		}
//	}
//	if err := stream.Err(); err != nil {
//		// handle err
//	}
type ConversationStream struct {
	body    io.ReadCloser
	scanner *bufio.Scanner
	started *conversationStarted
	event   ConversationStreamEvent
	err     error
}

func newConversationStream(body io.ReadCloser, started *conversationStarted) *ConversationStream {
	scanner := bufio.NewScanner(body)
	// A single "message" event can carry a sizeable chunk of the answer;
	// grow well past bufio.Scanner's 64KB default max line length.
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	return &ConversationStream{body: body, scanner: scanner, started: started}
}

// Next advances the stream to the next event. It returns false once the
// stream ends or an error occurs; call Err to distinguish the two.
func (s *ConversationStream) Next() bool {
	if s.err != nil {
		return false
	}
	payload, err := s.readEventData()
	if err != nil {
		if err != io.EOF {
			s.err = err
		}
		return false
	}
	event, err := s.decodeEvent(payload)
	if err != nil {
		s.err = err
		return false
	}
	s.event = event
	return true
}

// Event returns the event Next just advanced to. Only valid after a call to
// Next that returned true.
func (s *ConversationStream) Event() ConversationStreamEvent {
	return s.event
}

// Err returns the first error encountered while reading the stream, if any.
func (s *ConversationStream) Err() error {
	return s.err
}

// Close closes the underlying connection. Safe to call more than once.
func (s *ConversationStream) Close() error {
	return s.body.Close()
}

// readEventData reads raw SSE lines up to the next blank line, returning the
// joined "data:" payload. It returns io.EOF once the stream ends cleanly
// with no more data.
func (s *ConversationStream) readEventData() (string, error) {
	var dataLines []string
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if line == "" {
			if len(dataLines) > 0 {
				return strings.Join(dataLines, "\n"), nil
			}
			continue
		}
		// Other SSE fields (event:, id:, retry:, or ":" comments) carry no
		// information the API uses — the real event type is nested inside
		// the "data:" JSON payload's own "event" field (see sseEnvelope).
		if rest, ok := cutPrefix(line, "data:"); ok {
			dataLines = append(dataLines, strings.TrimPrefix(rest, " "))
		}
	}
	if err := s.scanner.Err(); err != nil {
		return "", err
	}
	if len(dataLines) > 0 {
		return strings.Join(dataLines, "\n"), nil
	}
	return "", io.EOF
}

func cutPrefix(s, prefix string) (string, bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

// sseEnvelope is the outer JSON object carried by every "data:" line. The
// SSE transport's own event name is always "message" — the real event type
// is this envelope's "event" field instead.
type sseEnvelope struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// decodeEvent parses one SSE frame's data payload into a
// ConversationStreamEvent, threading chat_uid/message_id captured from an
// earlier ChatStartedEvent into WorkflowFinishedEvent (which doesn't repeat
// them).
func (s *ConversationStream) decodeEvent(payload string) (ConversationStreamEvent, error) {
	var env sseEnvelope
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		return nil, err
	}

	switch env.Event {
	case "chat_started":
		var e ChatStartedEvent
		if err := json.Unmarshal(env.Data, &e); err != nil {
			return nil, err
		}
		s.started = &conversationStarted{ChatUID: e.ChatUID, MessageID: e.MessageID}
		return &e, nil

	case "workflow_started":
		var e WorkflowStartedEvent
		if err := json.Unmarshal(env.Data, &e); err != nil {
			return nil, err
		}
		return &e, nil

	case "message":
		var e MessageEvent
		if err := json.Unmarshal(env.Data, &e); err != nil {
			return nil, err
		}
		return &e, nil

	case "ping":
		return &PingEvent{}, nil

	case "chat_finished":
		var e ChatFinishedEvent
		if err := json.Unmarshal(env.Data, &e); err != nil {
			return nil, err
		}
		return &e, nil

	case "chat_title_updated":
		var e ChatTitleUpdatedEvent
		if err := json.Unmarshal(env.Data, &e); err != nil {
			return nil, err
		}
		return &e, nil

	case "workflow_finished":
		var payload workflowFinishedPayload
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			return nil, err
		}
		answer := ""
		if payload.Outputs.Answer != nil {
			answer = *payload.Outputs.Answer
		}
		resp := &ConversationResponse{
			Status:      payload.Status,
			Answer:      answer,
			References:  payload.Outputs.References,
			ElapsedTime: payload.ElapsedTime,
			Interrupt:   payload.Outputs.Interrupt,
			Error:       payload.Outputs.Error,
		}
		if s.started != nil {
			resp.ChatUID = s.started.ChatUID
			resp.MessageID = s.started.MessageID
		}
		return &WorkflowFinishedEvent{ConversationResponse: resp}, nil

	default:
		return &OtherEvent{Event: env.Event, Data: env.Data}, nil
	}
}
