package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"task/internal/models"
	"task/internal/service"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service service.Redirecter
	cache   service.Cache
}

func NewHandler(service service.Redirecter, cache service.Cache) *Handler {
	return &Handler{
		service: service,
		cache:   cache,
	}
}

func (h *Handler) InitRoutes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/admin/redirects", h.adminRedirects)
	router.GET("/admin/redirects/:id", h.adminRedirectsID)
	router.POST("/admin/redirects", h.createRedirect)
	router.PATCH("/admin/redirects/:id", h.updateRedirect)
	router.DELETE("/admin/redirects/:id", h.deleteRedirect)

	router.GET("/redirects", h.redirects)

	return router
}

// handler для пользовательского метода
// для передачи запрашеваемой ссылки берет query parameter как v
func (h *Handler) redirects(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	link := r.URL.Query().Get("v")

	correctLink, ok := h.cache.Get(link)

	// если запрашеваемой ссылки нет в кэше, ищем ссылку в базе
	if !ok {
		activeLink, err := h.service.ActiveLinkByHistory(link)
		if err != nil && !errors.Is(err, service.ErrNoEntry) {
			h.httpError(w, http.StatusInternalServerError, err)
			return
		}

		// если ActiveLinkByHistory(link) вернула пустую строку, то значит активная ссылка это запрашеваемая ссылка
		if activeLink == "" {
			activeLink = link
		} else {
			w.WriteHeader(301)
			w.Write([]byte(activeLink))
		}

		h.cache.Add(link, activeLink)
		return
	}

	if correctLink != link {
		w.WriteHeader(301)
		w.Write([]byte(correctLink))
	}

	log.Printf("FROM Cashe: key - %s, value - %s", link, correctLink)
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
		if errors.Is(err, service.ErrNoEntry) {
			h.httpError(w, http.StatusBadRequest, nil)
			return
		}
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
