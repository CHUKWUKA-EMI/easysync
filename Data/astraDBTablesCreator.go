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
			PRIMARY KEY (conversation_id,created_at)
			) WITH CLUSTERING ORDER BY (created_at DESC);`).Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = AstraDBSession.Query(`CREATE CUSTOM INDEX chat_id_idx ON easysynk.chats 
	(id) USING 'StorageAttachedIndex';`).Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = AstraDBSession.Query(`CREATE CUSTOM INDEX chat_message_idx ON easysynk.chats 
	(message) USING 'StorageAttachedIndex' 
	WITH OPTIONS = {'case_sensitive': 'false', 'normalize': 'true', 'ascii': 'true'}; `).Exec()
	if err != nil {
		log.Fatal(err)
	}
}
