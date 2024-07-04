package messages

import (
	"log"
	// "encoding/base64"
	// "encoding/json"
	"net/http"
	// "time"
	"github.com/gorilla/sessions"
    // "encoding/json"
)

type FlashMessage struct {
	Message     string `json:"message"`
	MessageType string `json:"type"`
}


var store = sessions.NewCookieStore([]byte("your-secret-key"))

func SetFlash(w http.ResponseWriter, r *http.Request, message string, messageType string) {
    session, _ := store.Get(r, "flash-session")
    session.AddFlash(FlashMessage{Message: message, MessageType: messageType})
    session.Save(r, w)
    log.Printf("Flash message set in session: %s (%s)", message, messageType)
}

func GetFlash(w http.ResponseWriter, r *http.Request) (string, string) {
    session, _ := store.Get(r, "flash-session")
    flashes := session.Flashes()
    if len(flashes) == 0 {
        return "", ""
    }
    session.Save(r, w)
    flash := flashes[0].(FlashMessage)
    log.Printf("Retrieved flash from session: %s (%s)", flash.Message, flash.MessageType)
    return flash.Message, flash.MessageType
}