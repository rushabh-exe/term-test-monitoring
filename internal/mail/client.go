package mailClient

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

var (
	clientID     string
	clientSecret string
	refreshToken string
)

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clientID = os.Getenv("EMAIL_CLIENTID")
	clientSecret = os.Getenv("EMAIL_CLIENTSECRET")
	refreshToken = os.Getenv("EMAIL_REFRESHTOKEN")
}

func getClient() *http.Client {
	loadEnvVariables()
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080",
		Scopes:       []string{gmail.GmailSendScope},
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := config.TokenSource(context.Background(), token)

	return oauth2.NewClient(context.Background(), tokenSource)
}

func SendEmail(service *gmail.Service, recipient, subject, body, attachmentPath string) error {
	var messageStr = []byte("To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: multipart/mixed; boundary=\"boundary\"\r\n" +
		"\r\n--boundary\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" + body + "\r\n" +
		"--boundary\r\n")

	if attachmentPath != "" {
		fileData, err := os.ReadFile(attachmentPath)
		if err != nil {
			return fmt.Errorf("failed to read attachment file: %v", err)
		}

		messageStr = append(messageStr,
			"Content-Disposition: attachment; filename=\""+filepath.Base(attachmentPath)+"\"\r\n"+
				"Content-Type: application/pdf; name=\""+filepath.Base(attachmentPath)+"\"\r\n"+
				"Content-Transfer-Encoding: base64\r\n"+
				"\r\n"+
				base64.StdEncoding.EncodeToString(fileData)+"\r\n"+
				"--boundary--\r\n"...)
	}

	raw := base64.URLEncoding.EncodeToString(messageStr)

	emailMessage := &gmail.Message{
		Raw: raw,
	}

	_, err := service.Users.Messages.Send("me", emailMessage).Do()
	return err
}

func MailClient(recipient, subject, body, attachmentPath string) {
	client := getClient()
	service, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Send the email with attachment
	if err := SendEmail(service, recipient, subject, body, attachmentPath); err != nil {
		log.Fatalf("Unable to send email: %v", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
