package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	dbConfig "github.com/chukwuka-emi/easysync/Data"
	chat "github.com/chukwuka-emi/easysync/Features/Chat"
)

func main() {

	if os.Getenv("GO_ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file", err.Error())
		}
	}

	dbConfig.InitDatabaseConnection()

	// dbConfig.CreateConversationTable()
	// dbConfig.CreateChatsTable()

	registerServices()
	chat.InitializeChatHub()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(cors.Default())

	handleRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5002"
	}
	router.Run(fmt.Sprintf(":%s", port))
}
