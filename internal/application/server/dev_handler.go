package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/jwt"
)

func devTokenHandler(db *pgxpool.Pool, jwtManager *jwt.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			CustomID string `json:"customId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.CustomID == "" {
			http.Error(w, "customId required", http.StatusBadRequest)
			return
		}

		q := sqlc.New(db)
		user, err := q.GetUserByCustomID(context.Background(), body.CustomID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		accessToken, err := jwtManager.GenerateAccessToken(user.ID)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken})
	}
}
