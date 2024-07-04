package messages

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type FlashMessage struct {
	Message     string `json:"message"`
	MessageType string `json:"type"`
}


func SetFlash(w http.ResponseWriter, r *http.Request, message string, messageType string) {
	flash := FlashMessage{Message: message, MessageType: messageType}
	encoded, err := encode(flash)
	if err != nil {
		// Handle error
		return
	}
	
	// Create the new flash cookie
	c := &http.Cookie{
		Name:  "flash",
		Value: encoded,
		Path:  "/",
	}

	// Check if a flash cookie already exists
	if existingCookie, err := r.Cookie("flash"); err == nil {
		// If the cookie exists, replace its value
		existingCookie.Value = encoded
		http.SetCookie(w, existingCookie)
	} else {
		// If the cookie does not exist, set the new cookie
		http.SetCookie(w, c)
	}
}

//-------------------------------------------------

func encode(flash FlashMessage) (string, error) {
	jsonData, err := json.Marshal(flash)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(jsonData), nil
}

func decode(src string, flash *FlashMessage) error {
	data, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, flash)
}



func GetFlash(w http.ResponseWriter, r *http.Request) (string, string) {
	c, err := r.Cookie("flash")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return "", ""
		default:
			return "", ""
		}
	}

	var flash FlashMessage
	err = decode(c.Value, &flash)
	if err != nil {
		return "", ""
	}

	// Delete the cookie after retrieving the value
	dc := &http.Cookie{Name: "flash", MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)

	return flash.Message, flash.MessageType
}