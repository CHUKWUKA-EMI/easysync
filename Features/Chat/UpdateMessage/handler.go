package updatemessage

import (
	"log"
	"time"

	data "github.com/chukwuka-emi/easysync/Data"
)

// Handler updates a chat message
func Handler(input MessageUpdateRequest) error {

	err := data.AstraDBSession.Query(`UPDATE easysynk.chats 
	                         SET message=?,updated_at=? WHERE conversation_id=? AND created_at=? IF EXISTS;`,
		input.Content,
		input.UpdatedAt,
		input.ConversationID,
		time.UnixMilli(int64(input.ID)),
	).Exec()
	if err != nil {
		log.Println("ERROR UPDATING CHAT MESSAGE: ", err)
		return err
	}
	return nil
}
