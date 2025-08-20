package api

import "net/http"

func NewRouter(h *Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/extrair", h.ExtractImagesHandler)
	router.HandleFunc("/extrair-conteudo-completo", h.ExtractFullContentHandler)
	router.HandleFunc("/extrair-documento-completo", h.ExtractRawDocumentHandler)

	// Servidor de arquivos estáticos
	fs := http.FileServer(http.Dir("imagens_extraidas"))
	router.Handle("/imagens/", http.StripPrefix("/imagens/", fs))

	return router
}
