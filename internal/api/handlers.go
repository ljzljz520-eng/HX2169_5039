package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/graduation-showcase/internal/domain"
)

func NewHandler(service *Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/records", recordsHandler(service))
	mux.HandleFunc("/records/", recordHandler(service))
	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "graduation-showcase"})
}

func recordsHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			var input CreateInput
			if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			record, err := service.CreateRecord(input)
			if err != nil {
				writeError(w, http.StatusUnprocessableEntity, err)
				return
			}
			writeJSON(w, http.StatusCreated, record)
			return
		}
		if request.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, nil)
			return
		}
		records, err := service.QueryRecords(QueryInput{Text: request.URL.Query().Get("q"), Label: request.URL.Query().Get("label"), Status: domain.Status(request.URL.Query().Get("status"))})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, records)
	}
}

func recordHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		id := strings.TrimPrefix(request.URL.Path, "/records/")
		if id == "" {
			writeError(w, http.StatusNotFound, nil)
			return
		}
		switch request.Method {
		case http.MethodGet:
			record, err := service.GetRecord(id)
			if err != nil {
				writeError(w, http.StatusNotFound, err)
				return
			}
			writeJSON(w, http.StatusOK, record)
		case http.MethodPatch:
			var input UpdateInput
			if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			record, err := service.UpdateLabel(id, input)
			if err != nil {
				writeError(w, http.StatusUnprocessableEntity, err)
				return
			}
			writeJSON(w, http.StatusOK, record)
		case http.MethodPost:
			action := request.URL.Query().Get("action")
			if action == "review" {
				handleReview(w, request, service, id)
				return
			}
			if action == "start-review" {
				record, err := service.BeginReview(id, request.URL.Query().Get("actor"))
				if err != nil {
					writeError(w, http.StatusUnprocessableEntity, err)
					return
				}
				writeJSON(w, http.StatusOK, record)
				return
			}
			if action == "archive" {
				record, err := service.ArchiveRecord(id, request.URL.Query().Get("actor"))
				if err != nil {
					writeError(w, http.StatusUnprocessableEntity, err)
					return
				}
				writeJSON(w, http.StatusOK, record)
				return
			}
			writeError(w, http.StatusBadRequest, nil)
		default:
			writeError(w, http.StatusMethodNotAllowed, nil)
		}
	}
}

func handleReview(w http.ResponseWriter, request *http.Request, service *Service, id string) {
	var input ReviewInput
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, err := service.ReviewRecord(id, input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	message := http.StatusText(status)
	if err != nil {
		message = err.Error()
	}
	writeJSON(w, status, map[string]string{"error": message})
}
