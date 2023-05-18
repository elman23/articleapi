package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elman23/articleapi/pkg/models"
	"github.com/gorilla/mux"
)

func (h handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Read the request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var updatedArticle models.Article
	json.Unmarshal(body, &updatedArticle)

	queryStmt := `UPDATE articles SET title = $2, description = $3, content = $4 WHERE id = $1 RETURNING id;`
	err = h.DB.QueryRow(queryStmt, &id, &updatedArticle.Title, &updatedArticle.Desc, &updatedArticle.Content).Scan(&id)
	if err != nil {
		log.Println("Failed to execute query!", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated")
}
