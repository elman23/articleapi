package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/elman23/articleapi/pkg/models"
	"github.com/gorilla/mux"
)

func (h handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint: [GetArtcle].")

	vars := mux.Vars(r)
	id := vars["id"]

	queryStmt := `SELECT * FROM articles WHERE id = $1 ;`
	results, err := h.DB.Query(queryStmt, id)
	if err != nil {
		log.Println("Failed to execute query!", err)
		w.WriteHeader(500)
		return
	}

	var article models.Article
	for results.Next() {
		err = results.Scan(&article.Id, &article.Title, &article.Desc, &article.Content)
		if err != nil {
			log.Println("Failed to scan item!", err)
			w.WriteHeader(500)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}
