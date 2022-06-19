package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"testjunior/database"
	"testjunior/tokens"
)

type Tokens struct {
	Access_token  string
	Refresh_token string
}

func RecieveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		guid := r.URL.Query().Get("GUID")

		access_token, refresh_token64 := tokens.GenerateAllTokens(guid)
		refresh_hash := tokens.Bcrypt(refresh_token64)
		database.Save_tokens(refresh_hash, guid)

		t := Tokens{
			Access_token:  access_token,
			Refresh_token: refresh_token64,
		}

		mess, err := json.Marshal(t)
		if err != nil {
			log.Panic(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mess)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
