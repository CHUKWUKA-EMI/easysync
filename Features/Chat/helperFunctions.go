package chat

import (
	"time"

	deletemessage "github.com/chukwuka-emi/easysync/Features/Chat/DeleteMessage"
	updatemessage "github.com/chukwuka-emi/easysync/Features/Chat/UpdateMessage"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type outgoingMessage struct {
	Data       interface{} `json:"data"`
	ActionType actionType  `json:"actionType"`
}

func handleNewChat(
	websocketConn *websocket.Conn,
	message Message,
	conversationID uuid.UUID,
	messageChannel map[*websocket.Conn]bool,
) {

	var chatData Chat

	chatData = Chat{
		ID:             uint64(time.Now().UnixMilli()),
		ConversationID: gocql.UUID(conversationID),
		SenderID:       gocql.UUID(message.SenderID),
		Content:        message.Content,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err := saveMessage(chatData)
	if err != nil {
		websocketConn.WriteMessage(websocket.TextMessage, []byte("We encountered an error. Please try again"))
		return
	}
	for websocketClient := range messageChannel {
		if websocketClient != nil {
			websocketClient.WriteJSON(outgoingMessage{
				ActionType: insertNewMessage,
				Data:       chatData,
			})
		}
	}
}

func updateExistingChat(
	websocketConn *websocket.Conn,
	message Message,
	conversationID uuid.UUID,
	messageChannel map[*websocket.Conn]bool,
) {
	var payload updatemessage.MessageUpdateRequest

	payload = updatemessage.MessageUpdateRequest{
		ID:             message.ChatID,
		ConversationID: gocql.UUID(conversationID),
		Content:        message.Content,
		UpdatedAt:      time.Now(),
	}

	err := updatemessage.Handler(payload)

	if err != nil {
		websocketConn.WriteMessage(websocket.TextMessage, []byte("We encountered an error. Please try again"))
		return
	}

	for websocketClient := range messageChannel {
		if websocketClient != nil {
			websocketClient.WriteJSON(outgoingMessage{
				ActionType: updateMessage,
				Data:       payload,
			})
		}
	}
}

func deleteChat(
	websocketConn *websocket.Conn,
	message Message,
	conversationID uuid.UUID,
	messageChannel map[*websocket.Conn]bool,
) {
	var payload deletemessage.MessageDeleteRequest

	payload = deletemessage.MessageDeleteRequest{
		ID:             message.ChatID,
		ConversationID: gocql.UUID(conversationID),
	}

	err := deletemessage.Handler(payload)

	if err != nil {
		websocketConn.WriteMessage(websocket.TextMessage, []byte("We encountered an error. Please try again"))
		return
	}

	for websocketClient := range messageChannel {
		if websocketClient != nil {
			websocketClient.WriteJSON(outgoingMessage{
				ActionType: deleteMessage,
				Data:       payload,
			})
		}
	}
}
