package post

import (
	"boilerplate-elastic-search/internal/pkg/domain"
	"boilerplate-elastic-search/internal/pkg/storage"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service service
}

func New(storage storage.PostStorer) Handler {
	return Handler{
		service: service{storage: storage},
	}
}

// POST /api/v1/posts
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	res, err := h.service.create(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrConflict:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(res)
	_, _ = w.Write(body)
}

// PATCH /api/v1/posts:id
func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	req.ID = httprouter.ParamsFromContext(r.Context()).ByName("id")

	if err := h.service.update(r.Context(), req); err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/posts/:id

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	if err := h.service.delete(r.Context(), DeleteRequest{ID: id}); err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)

		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) Find(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	res, err := h.service.find(r.Context(), FindRequest{ID: id})
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(res)
	_, _ = w.Write(body)
}

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	keyword := httprouter.ParamsFromContext(r.Context()).ByName("keyword")

	res, err := h.service.search(r.Context(), SearchRequest{Keyword: keyword})
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(res)
	_, _ = w.Write(body)
}
