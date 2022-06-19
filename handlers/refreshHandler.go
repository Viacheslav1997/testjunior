package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"testjunior/database"
	"testjunior/tokens"
)

type Foo struct {
	RefreshToken string
	AccessToken  string
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//получаем id user-а и tokens
		guid := r.URL.Query().Get("GUID")
		var user_tokens Foo
		err := json.NewDecoder(r.Body).Decode(&user_tokens)
		if err != nil {
			log.Fatal(err)
		}

		//проверяем, есть ли юзер с таким id и рефреш токеном в бд
		session_check_res, err, res := database.Check_session(guid, user_tokens.RefreshToken)

		//проверяем, был ли выдан данный access токен вместе с этим refresh
		tokens_coonect_res := database.ValidateSession(user_tokens.RefreshToken, user_tokens.AccessToken)

		if (session_check_res) && (tokens_coonect_res) {

			//Удаляем старый
			database.DeleteClientSessionByGUID(guid, res.Refresh_token)
			//Создаем два новых токена
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

		} else {
			http.Error(w, "Error", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
