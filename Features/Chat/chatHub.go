package chat

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Hub houses the core of the chat service
type Hub struct {
	Upgrader      websocket.Upgrader
	Channels      map[uuid.UUID]map[*websocket.Conn]bool
	ChannelsMutex sync.Mutex
}

// Message is the strcuture of a chat payload
type Message struct {
	SenderID uuid.UUID `json:"senderId"`
	Content  string    `json:"content"`
}

// ChatHub points to the chat hub
var ChatHub *Hub

// InitializeChatHub initializes the Hub
func InitializeChatHub() {
	ChatHub = &Hub{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		Channels:      make(map[uuid.UUID]map[*websocket.Conn]bool),
		ChannelsMutex: sync.Mutex{},
	}
}

// HandleUpgrade handles websocket upgrades/connections
func (h *Hub) HandleUpgrade(ctx *gin.Context) {
	websocketConn, err := h.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	var conversationID uuid.UUID
	conversationIDStr := ctx.Request.URL.Query().Get("conversationId")

	if conversationIDStr == "" {
		log.Println("Channel not specified.")
		websocketConn.Close()
		return
	}

	conversationID, err = uuid.Parse(conversationIDStr)

	if err != nil {
		log.Println("Error parsing conversationID", err.Error())
		websocketConn.Close()
		return
	}

	h.ChannelsMutex.Lock()
	if _, ok := h.Channels[conversationID]; !ok {
		h.Channels[conversationID] = make(map[*websocket.Conn]bool)
	}

	h.Channels[conversationID][websocketConn] = true
	h.ChannelsMutex.Unlock()

	go handleMessaging(websocketConn, conversationID)
}

func handleMessaging(websocketConn *websocket.Conn, conversationID uuid.UUID) {
	defer closeConnection(websocketConn, conversationID)

	for {
		var incomingMessage Message
		err := websocketConn.ReadJSON(&incomingMessage)
		if err != nil {
			break
		}
		chatData := Chat{
			ID:             uint64(time.Now().UnixMilli()),
			ConversationID: gocql.UUID(conversationID),
			SenderID:       gocql.UUID(incomingMessage.SenderID),
			Message:        incomingMessage.Content,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		go sendMessage(websocketConn, chatData)
	}
}

func sendMessage(websocketConn *websocket.Conn, chatData Chat) {
	ChatHub.ChannelsMutex.Lock()
	defer ChatHub.ChannelsMutex.Unlock()
	messageChannel, ok := ChatHub.Channels[uuid.UUID(chatData.ConversationID)]
	if ok {
		err := saveMessage(chatData)
		if err != nil {
			websocketConn.WriteMessage(websocket.TextMessage, []byte("We encountered an error. Please try again"))
			return
		}
		for websocketClient := range messageChannel {
			if websocketClient != nil {
				websocketClient.WriteJSON(Chat{
					ID:             chatData.ID,
					SenderID:       chatData.SenderID,
					ConversationID: chatData.ConversationID,
					Message:        chatData.Message,
					CreatedAt:      chatData.CreatedAt,
					UpdatedAt:      chatData.UpdatedAt,
				})
			}
		}
	}
}

func closeConnection(websocketConn *websocket.Conn, conversationID uuid.UUID) {
	ChatHub.ChannelsMutex.Lock()
	delete(ChatHub.Channels[conversationID], websocketConn)
	websocketConn.Close()
	ChatHub.ChannelsMutex.Unlock()
}
