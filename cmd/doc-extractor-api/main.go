package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/marcokist/googledocs-extractor/internal/api"
	"github.com/marcokist/googledocs-extractor/internal/extractor"
	"github.com/marcokist/googledocs-extractor/internal/gdocs"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	ctx := context.Background()

	logger.Info("Iniciando a aplicação...")

	gdocsService, err := gdocs.NewService(ctx)
	if err != nil {
		logger.Error("Erro fatal ao inicializar o serviço do Google", "erro", err)
		os.Exit(1)
	}

	extractorService := extractor.NewService(gdocsService, logger) // <-- Passa o logger

	apiHandler := api.NewHandler(extractorService, logger) // <-- Passa o logger

	router := api.NewRouter(apiHandler)

	port := "8080"
	addr := ":" + port

	logger.Info("Servidor iniciado e escutando", "endereço", fmt.Sprintf("http://localhost:%s", port))
	logger.Info("Endpoints disponíveis:",
		"extrair_imagens", "/extrair?doc_id=ID...",
		"extrair_conteudo_completo", "/extrair-conteudo-completo?doc_id=ID...",
		"extrair_documento_completo", "/extrair-documento-completo?doc_id=ID...",
	)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("Não foi possível iniciar o servidor", "erro", err)
		os.Exit(1)
	}
}
