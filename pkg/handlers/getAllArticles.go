package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/elman23/articleapi/pkg/models"
)

func (h handler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint: [GetAllArticles].")

	results, err := h.DB.Query("SELECT * FROM articles;")
	if err != nil {
		log.Println("Failed to execute query!", err)
		w.WriteHeader(500)
		return
	}

	var articles = make([]models.Article, 0)
	for results.Next() {
		var article models.Article
		err = results.Scan(&article.Id, &article.Title, &article.Desc, &article.Content)
		if err != nil {
			log.Println("Failed to scan item!", err)
			w.WriteHeader(500)
			return
		}

		articles = append(articles, article)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}
