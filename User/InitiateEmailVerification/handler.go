package initiateemailverification

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	services "github.com/chukwuka-emi/easysync/Services"
	"github.com/chukwuka-emi/easysync/utils"
	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

// Handler sends a verification code to the provided email
func Handler(c *gin.Context) {
	var input request
	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := utils.Generate8DigitsNumber()

	email := &services.EmailService{
		Sender:    services.EmailUser{Name: "EasySync", Email: os.Getenv("SENDER_EMAIL")},
		Recipient: services.EmailUser{Name: "", Email: input.Email},
		Subject:   fmt.Sprintf("EasySync Confirmation Code: %d", code),
		Content: fmt.Sprintf(`<h1>Confirm your email address</h1>
		          <p>Use the code below to confirm your email address</p>
				  <p style="letter-spacing:4px;text-align:center;font-size:30px;font-weight:bold;">%d</p>
		         `, code),
	}
	codeExpiration := 180 * time.Second
	_, err = services.RedisClient.SetEx(ctx, input.Email, code, codeExpiration).Result()
	if err != nil {
		log.Fatal("Error saving verification code to cache", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	email.SendEmail()
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("A confirmation code has been sent to the email. It expires in %s", codeExpiration)})

}
