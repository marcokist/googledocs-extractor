package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/marcokist/googledocs-extractor/internal/extractor"
	"github.com/marcokist/googledocs-extractor/internal/model"
)

type Handler struct {
	extractor *extractor.Service
	logger    *slog.Logger
}

func NewHandler(e *extractor.Service, logger *slog.Logger) *Handler { // <-- Recebe o logger
	return &Handler{
		extractor: e,
		logger:    logger,
	}
}

func (h *Handler) ExtractImagesHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Requisição recebida",
		"metodo", r.Method,
		"endpoint", r.URL.Path,
		"query", r.URL.RawQuery,
	)
	docIDs := r.URL.Query()["doc_id"]
	if len(docIDs) == 0 {
		http.Error(w, "O parâmetro 'doc_id' é obrigatório", http.StatusBadRequest)
		return
	}
	results := make(map[string]model.ResponseData)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, docID := range docIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			imageData, err := h.extractor.ExtractImages(id)
			mu.Lock()
			if err != nil {
				results[id] = model.ResponseData{Status: "error", Message: err.Error()}
			} else {
				results[id] = model.ResponseData{Status: "success", Images: imageData}
			}
			mu.Unlock()
		}(docID)
	}
	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) ExtractFullContentHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Requisição recebida",
		"metodo", r.Method,
		"endpoint", r.URL.Path,
		"query", r.URL.RawQuery,
	)
	docIDs := r.URL.Query()["doc_id"]
	if len(docIDs) == 0 {
		http.Error(w, "O parâmetro 'doc_id' é obrigatório", http.StatusBadRequest)
		return
	}
	results := make(map[string]model.FullContentResponseData)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, docID := range docIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			content, err := h.extractor.ExtractFullContent(id)
			mu.Lock()
			if err != nil {
				results[id] = model.FullContentResponseData{Status: "error", Message: err.Error()}
			} else {
				results[id] = model.FullContentResponseData{Status: "success", Content: content}
			}
			mu.Unlock()
		}(docID)
	}
	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) ExtractRawDocumentHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Requisição recebida",
		"metodo", r.Method,
		"endpoint", r.URL.Path,
		"query", r.URL.RawQuery,
	)
	docIDs := r.URL.Query()["doc_id"]
	if len(docIDs) == 0 {
		http.Error(w, "O parâmetro 'doc_id' é obrigatório", http.StatusBadRequest)
		return
	}
	results := make(map[string]any)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, docID := range docIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			doc, err := h.extractor.ExtractRawDocument(id)
			mu.Lock()
			if err != nil {
				results[id] = model.ErrorResponse{Status: "error", Message: err.Error()}
			} else {
				results[id] = doc
			}
			mu.Unlock()
		}(docID)
	}
	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
