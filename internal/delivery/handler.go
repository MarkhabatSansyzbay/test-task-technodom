package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"task/internal/models"
	"task/internal/service"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service service.Redirecter
}

func NewHandler(service service.Redirecter) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/admin/redirects", h.adminRedirects)
	router.GET("/admin/redirects/:id", h.adminRedirectsID)
	router.POST("/admin/redirects", h.createRedirect)
	router.PATCH("/admin/redirects/:id", h.updateRedirect)
	router.DELETE("/admin/redirects/:id", h.deleteRedirect)

	return router
}

func (h *Handler) deleteRedirect(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idParam := p.ByName("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.httpError(w, http.StatusNotFound, nil)
		return
	}

	if err := h.service.DeleteRedirect(id); err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) updateRedirect(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idParam := p.ByName("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.httpError(w, http.StatusNotFound, nil)
		return
	}

	var newLink models.Link
	if err := json.NewDecoder(r.Body).Decode(&newLink); err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.service.UpdateRedirect(id, newLink.ActiveLink); err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) createRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var redirect models.Link
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&redirect); err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.service.CreateRedirect(redirect); err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) adminRedirectsID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idParam := p.ByName("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.httpError(w, http.StatusNotFound, nil)
		return
	}

	res, err := h.service.RedirectByID(id)
	if err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) adminRedirects(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirects, err := h.service.Redirects()
	if err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(redirects); err != nil {
		h.httpError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) httpError(w http.ResponseWriter, status int, err error) {
	if status == http.StatusInternalServerError {
		log.Println(err)
	}

	http.Error(w, http.StatusText(status), status)
}
