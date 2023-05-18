package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/elman23/articleapi/pkg/mocks"
	"github.com/gorilla/mux"
)

func (h handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, article := range mocks.Articles {
		if article.Id == id {
			mocks.Articles = append(mocks.Articles[:index], mocks.Articles[index+1:]...)

			w.Header().Add("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Deleted")
			break
		}
	}
}
