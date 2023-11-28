package chat

import (
	"log"

	data "github.com/chukwuka-emi/easysync/Data"
)

func saveMessage(chatData Chat) error {
	err := data.AstraDBSession.Query(`INSERT into easysynk.chats 
	                         (id,conversation_id,sender_id,message,created_at,updated_at) 
							 VALUES (?,?,?,?,?,?);`,
		chatData.ID,
		chatData.ConversationID,
		chatData.SenderID,
		chatData.Content,
		chatData.CreatedAt,
		chatData.UpdatedAt,
	).Exec()
	if err != nil {
		log.Println("ERROR SAVING CHAT MESSAGE: ", err)
		return err
	}
	return nil
}
