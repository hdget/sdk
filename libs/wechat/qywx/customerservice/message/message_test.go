package message

import (
	"testing"
)

func TestNewTextMessage(t *testing.T) {
	// 测试文本消息
	data := []byte(`{
		"msgtype": "text",
		"text": {
			"content": "hello world"
		}
	}`)

	msg, err := New(data)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if msg.GetKind() != MessageKindText {
		t.Errorf("GetKind() = %v, want %v", msg.GetKind(), MessageKindText)
	}

	textMsg, ok := msg.(*TextNormalMessage)
	if !ok {
		t.Fatal("msg is not TextNormalMessage")
	}

	if textMsg.Text == nil || textMsg.Text.Content != "hello world" {
		t.Errorf("Text content = %v, want 'hello world'", textMsg.Text)
	}
}

func TestNewImageMessage(t *testing.T) {
	// 测试图片消息
	data := []byte(`{
		"msgtype": "image",
		"image": {
			"media_id": "test_media_id"
		}
	}`)

	msg, err := New(data)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if msg.GetKind() != MessageKindImage {
		t.Errorf("GetKind() = %v, want %v", msg.GetKind(), MessageKindImage)
	}

	imageMsg, ok := msg.(*ImageNormalMessage)
	if !ok {
		t.Fatal("msg is not ImageNormalMessage")
	}

	if imageMsg.GetMediaID() != "test_media_id" {
		t.Errorf("MediaID = %v, want 'test_media_id'", imageMsg.GetMediaID())
	}
}

func TestNewEnterSessionEvent(t *testing.T) {
	// 测试进入会话事件
	data := []byte(`{
		"msgtype": "event",
		"event": {
			"event_type": "enter_session",
			"open_kfid": "wkAJ2GCAAASSm4_FhToWMFea0xAFfd3Q",
			"external_userid": "wmAJ2GCAAAme1XQRC-NI-q0_ZM9ukoAw",
			"scene": "123",
			"welcome_code": "test_code"
		}
	}`)

	msg, err := New(data)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if msg.GetKind() != MessageKindEventEnterSession {
		t.Errorf("GetKind() = %v, want %v", msg.GetKind(), MessageKindEventEnterSession)
	}

	eventMsg, ok := msg.(*EnterSessionEventMessage)
	if !ok {
		t.Fatal("msg is not EnterSessionEventMessage")
	}

	if eventMsg.Event == nil || eventMsg.Event.WelcomeCode != "test_code" {
		t.Errorf("WelcomeCode = %v, want 'test_code'", eventMsg.Event)
	}
}

func TestReplyText(t *testing.T) {
	// 测试回复文本消息
	msg := &Message{
		ExternalUserID: "test_user",
		OpenKfID:       "test_kf",
	}

	replyData, err := msg.ReplyText("test reply")
	if err != nil {
		t.Fatalf("ReplyText() error = %v", err)
	}

	t.Logf("Reply data: %s", string(replyData))
}
