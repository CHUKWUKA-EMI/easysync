package deletemessage

import (
	"log"
	"time"

	data "github.com/chukwuka-emi/easysync/Data"
)

// Handler deletes a message from the database
func Handler(input MessageDeleteRequest) error {

	err := data.AstraDBSession.Query(`DELETE FROM easysynk.chats 
	                         WHERE conversation_id=? AND created_at=? IF EXISTS;`,
		input.ConversationID,
		time.UnixMilli(int64(input.ID)),
	).Exec()
	if err != nil {
		log.Println("ERROR DELETING CHAT MESSAGE: ", err)
		return err
	}
	return nil
}
