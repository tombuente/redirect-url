package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tombuente/redirect-url/xerrors"
)

type RedirectHandler struct {
	service RedirectService
}

func NewRedirectHandler(service RedirectService) chi.Router {
	h := RedirectHandler{
		service: service,
	}

	r := chi.NewRouter()
	r.Get("/{key}", h.getRedirect)

	return r
}

func (h RedirectHandler) getRedirect(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "key"), encodingBase, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
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

	http.Redirect(w, r, url.URL, http.StatusTemporaryRedirect)
}
