package gdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

func NewService(ctx context.Context) (*docs.Service, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("não foi possível ler o arquivo de credenciais: %w", err)
	}

	config, err := google.ConfigFromJSON(b, docs.DocumentsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("não foi possível analisar o arquivo de segredos do cliente: %w", err)
	}

	client := getClient(config)

	srv, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar o serviço do Docs: %w", err)
	}
	return srv, nil
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Acesse esta URL no seu navegador e autorize o aplicativo: \n%v\n", authURL)
	fmt.Print("Digite o código de autorização do navegador: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Não foi possível ler o código de autorização: %v", err)
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Não foi possível obter o token da web: %v", err)
	}
	return tok
}
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Salvando o arquivo de token em: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Não foi possível salvar o token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
