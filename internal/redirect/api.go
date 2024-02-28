package redirect

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tombuente/redirect-url/xerrors"
)

const encodingBase = 36

type RedirectAPIHandler struct {
	service RedirectService
}

type URLResponse struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type URLRequest struct {
	URL string `json:"url"`
}

func NewRedirectAPIHandler(service RedirectService) chi.Router {
	h := RedirectAPIHandler{
		service: service,
	}

	r := chi.NewRouter()
	r.Get("/urls/{key}", h.getURL)
	r.Post("/urls/", h.postURL)

	return r
}

func (h RedirectAPIHandler) getURL(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "key"), encodingBase, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.service.GetURL(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, xerrors.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			slog.Error("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	render.JSON(w, r, newURLResponse(url))
}

func (h RedirectAPIHandler) postURL(w http.ResponseWriter, r *http.Request) {
	var req URLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.service.CreateURL(r.Context(), newURLParams(req))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, newURLResponse(url))
}

func newURLResponse(url URL) URLResponse {
	encodedID := strconv.FormatInt(url.ID, 36)
	return URLResponse{
		Key: encodedID,
		URL: url.URL,
	}
}

func newURLParams(req URLRequest) URLParams {
	return URLParams(req)
}
