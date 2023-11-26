package data

import "log"

// CreateConversationTable creates conversations table
func CreateConversationTable() {
	err := AstraDBSession.Query(`CREATE TABLE easysynk.conversations 
	       (id uuid, 
			members set<uuid>, 
			is_open boolean,
			created_at timestamp,
			PRIMARY KEY (id, created_at)
			);`).Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = AstraDBSession.Query(`CREATE INDEX member_idx ON easysynk.conversations 
	       (members);`).Exec()
	if err != nil {
		log.Fatal(err)
	}
}

// CreateChatsTable creates chats table
func CreateChatsTable() {
	err := AstraDBSession.Query(`CREATE TABLE easysynk.chats 
	       (id bigint, 
			conversation_id uuid,
			sender_id uuid, 
			message text, 
			created_at timestamp, 
			updated_at timestamp,
			PRIMARY KEY ((conversation_id,id),created_at)
			);`).Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = AstraDBSession.Query(`CREATE INDEX message_idx ON easysynk.chats 
	       (message);`).Exec()
	if err != nil {
		log.Fatal(err)
	}
}
