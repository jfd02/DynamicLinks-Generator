package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"dynamic-links-generator/api/apperrors"
	"dynamic-links-generator/api/models"
	"dynamic-links-generator/api/service"

	"github.com/rs/zerolog/log"
)

type Handler interface {
	CreateLink(w http.ResponseWriter, r *http.Request)
	ExchangeShortLink(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	linkService service.LinkService
}

func NewHandler(linkService service.LinkService) Handler {
	return &handler{
		linkService: linkService,
	}
}

func (h *handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var rawReq map[string]any
	if err := json.NewDecoder(r.Body).Decode(&rawReq); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", "INVALID_ARGUMENT")
		return
	}

	createReq, err := h.linkService.PrepareDynamicLinkRequest(rawReq)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidURLFormat):
			WriteErrorResponse(w, http.StatusBadRequest, "longDynamicLink is not parsable", "INVALID_ARGUMENT")
		case errors.Is(err, apperrors.ErrHostInvalid):
			WriteErrorResponse(w, http.StatusBadRequest, "Host is invalid", "INVALID_ARGUMENT")
		case errors.Is(err, apperrors.ErrInvalidFormat),
			errors.Is(err, apperrors.ErrMissingHost),
			errors.Is(err, apperrors.ErrMissingLink):
			WriteErrorResponse(w, http.StatusBadRequest, err.Error(), "INVALID_ARGUMENT")
		default:
			WriteErrorResponse(w, http.StatusBadRequest, "Invalid request format", "INVALID_ARGUMENT")
		}
		return
	}

	link, err := h.linkService.CreateDynamicLink(r.Context(), createReq)
	if errors.Is(err, apperrors.ErrDomainLinkNotAllowed) {
		WriteErrorResponse(w, http.StatusBadRequest, "'link' parameter contains a host that is not in the allow list", "INVALID_ARGUMENT")
		return
	} else if errors.Is(err, apperrors.ErrInvalidAppStoreID) {
		WriteErrorResponse(w, http.StatusBadRequest, "'isbn' parameter contains a non-numeric value", "INVALID_ARGUMENT")
		return
	} else if err != nil {
		log.Error().Err(err).Msg("Failed to create dynamic link")
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create link", "INTERNAL")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func (h *handler) ExchangeShortLink(w http.ResponseWriter, r *http.Request) {
	var req models.ExchangeShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RequestedLink == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid or missing requestedLink", "INVALID_ARGUMENT")
		return
	}

	link, err := h.linkService.ResolveShortPath(r.Context(), req.RequestedLink)
	switch {
	case errors.Is(err, apperrors.ErrLinkNotFound):
		WriteErrorResponse(w, http.StatusNotFound, "Link not found", "NOT_FOUND")
	case errors.Is(err, apperrors.ErrInvalidRequestedLink):
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid requested link", "INVALID_ARGUMENT")
	case err != nil:
		log.Error().Err(err).Msg("Failed to resolve short link")
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to resolve link", "INTERNAL")
	default:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(link)
	}
}

func WriteErrorResponse(w http.ResponseWriter, code int, message string, status string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: models.ErrorDetails{
			Code:    code,
			Message: message,
			Status:  status,
		},
	})
}
