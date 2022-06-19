package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testjunior/database"
	"testjunior/models"
	"testjunior/tokens"
)

type Foo struct {
	Token string
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//получаем id user-а и base64 refresh token
		guid := r.URL.Query().Get("GUID")
		var user_refresh_token Foo
		err := json.NewDecoder(r.Body).Decode(&user_refresh_token)
		if err != nil {
			log.Fatal(err)
		}

		//проверяем, есть ли юзер с таким id в бд
		session_check_res, err, result := database.Check_session(guid)
		if err != nil {
			result = models.User{Refresh_token: nil}
			log.Fatal(err)
		}

		//проверяем, валидный ли его refresh base64 токен
		//
		token_check_res, _ := database.Check_refresh_token(user_refresh_token.Token, result)
		if (session_check_res) && (token_check_res) {

			//Удаляем старый
			database.DeleteClientSessionByGUID(guid)
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
			//fmt.Fprint(w, "Можно делать рефреш\n")
			//fmt.Fprint(w, result)
		} else if (session_check_res) && (!token_check_res) {
			fmt.Fprint(w, "Invalid refresh session (your refresh session incorrect)\n")
		} else if !session_check_res {
			fmt.Fprint(w, "For this user no accessible session (no user with this GUID)\n")
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
