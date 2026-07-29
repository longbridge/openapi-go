package agent

import (
	"encoding/json"
	"testing"
)

// The `data` payload of the "Run succeeded" example from
// https://open.longbridge.com/en/docs/ai/chat/conversation
const succeededJSON = `{
	"chat_uid": "ct_9f2c1a5b",
	"message_id": "42",
	"status": "succeeded",
	"answer": "Tesla (TSLA.US) recently...",
	"references": [
		{ "index": 1, "title": "...", "url": "..." }
	],
	"elapsed_time": 3.21
}`

// The `data` payload of the "Run interrupted" example from the same page.
const interruptedJSON = `{
	"chat_uid": "ct_9f2c1a5b",
	"message_id": "43",
	"status": "interrupted",
	"answer": "",
	"references": null,
	"elapsed_time": 1.05,
	"interrupt": {
		"node_id": "n_ask_human",
		"tool_call_id": "call_abc123",
		"questions": [
			{
				"question": "Which time range would you like to check?",
				"options": [
					{ "description": "Past week" },
					{ "description": "Past month" }
				],
				"multi_select": false
			}
		],
		"message_id": 43,
		"chat_id": 1001
	}
}`

func TestUnmarshalSucceededConversationResponse(t *testing.T) {
	var resp ConversationResponse
	if err := json.Unmarshal([]byte(succeededJSON), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.ChatUID != "ct_9f2c1a5b" {
		t.Errorf("ChatUID = %q", resp.ChatUID)
	}
	if resp.MessageID != "42" {
		t.Errorf("MessageID = %q", resp.MessageID)
	}
	if resp.Status != ConversationStatusSucceeded {
		t.Errorf("Status = %v", resp.Status)
	}
	if resp.Answer != "Tesla (TSLA.US) recently..." {
		t.Errorf("Answer = %q", resp.Answer)
	}
	if len(resp.References) != 1 || resp.References[0].Index != 1 {
		t.Errorf("References = %+v", resp.References)
	}
	if resp.ElapsedTime != 3.21 {
		t.Errorf("ElapsedTime = %v", resp.ElapsedTime)
	}
	if resp.Interrupt != nil {
		t.Errorf("Interrupt = %+v, want nil", resp.Interrupt)
	}
	if resp.Error != nil {
		t.Errorf("Error = %+v, want nil", resp.Error)
	}
}

func TestUnmarshalInterruptedConversationResponse(t *testing.T) {
	var resp ConversationResponse
	if err := json.Unmarshal([]byte(interruptedJSON), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Status != ConversationStatusInterrupted {
		t.Errorf("Status = %v", resp.Status)
	}
	if resp.References != nil {
		t.Errorf("References = %+v, want nil", resp.References)
	}
	interrupt := resp.Interrupt
	if interrupt == nil {
		t.Fatal("Interrupt = nil, want non-nil")
	}
	if interrupt.NodeID != "n_ask_human" {
		t.Errorf("NodeID = %q", interrupt.NodeID)
	}
	if interrupt.ToolCallID != "call_abc123" {
		t.Errorf("ToolCallID = %q", interrupt.ToolCallID)
	}
	if interrupt.MessageID != 43 {
		t.Errorf("MessageID = %d", interrupt.MessageID)
	}
	if interrupt.ChatID != 1001 {
		t.Errorf("ChatID = %d", interrupt.ChatID)
	}
	if len(interrupt.Questions) != 1 {
		t.Fatalf("Questions = %+v", interrupt.Questions)
	}
	q := interrupt.Questions[0]
	if len(q.Options) != 2 {
		t.Errorf("Options = %+v", q.Options)
	}
	if q.MultiSelect {
		t.Error("MultiSelect = true, want false")
	}
}

func TestUnmarshalWorkspacesResponse(t *testing.T) {
	const j = `{
		"workspaces": [
			{ "id": "1001", "name": "My Workspace", "created_at": 1742000000, "updated_at": 1742001000 }
		]
	}`
	var resp WorkspacesResponse
	if err := json.Unmarshal([]byte(j), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(resp.Workspaces) != 1 || resp.Workspaces[0].ID != "1001" {
		t.Errorf("Workspaces = %+v", resp.Workspaces)
	}
}

func TestUnmarshalAgentsResponse(t *testing.T) {
	const j = `{
		"agents": [
			{
				"uid": "ag_7d3f9b2c",
				"name": "US Stock Analyst",
				"description": "Answers US stock questions with market and fundamental data",
				"mode": "chat",
				"icon": "https://cdn.longbridge.com/icons/agent.png",
				"is_published": true,
				"published_at": 1742000000,
				"created_at": 1741000000,
				"updated_at": 1742001000
			}
		],
		"total": 12
	}`
	var resp AgentsResponse
	if err := json.Unmarshal([]byte(j), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Total != 12 {
		t.Errorf("Total = %d", resp.Total)
	}
	if len(resp.Agents) != 1 || resp.Agents[0].UID != "ag_7d3f9b2c" {
		t.Errorf("Agents = %+v", resp.Agents)
	}
	if !resp.Agents[0].IsPublished {
		t.Error("IsPublished = false, want true")
	}
}

func TestContinueBodyMarshal(t *testing.T) {
	body := continueBody{AnswersByToolCall: AnswersByToolCall{
		"call_abc123": {"Which time range would you like to check?": "Past month"},
	}}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var round continueBody
	if err := json.Unmarshal(b, &round); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if round.AnswersByToolCall["call_abc123"]["Which time range would you like to check?"] != "Past month" {
		t.Errorf("round-tripped body = %+v", round)
	}
}
