## Go Google Docs Extractor API

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Uma API robusta em Go para extrair conteúdo de documentos do Google Docs, incluindo imagens (com conversão para Base64), texto estruturado e o objeto de dados completo da API. A aplicação é estruturada, concorrente e produz logs em formato JSON.

## ✨ Funcionalidades

- **API RESTful**: Interface HTTP para interagir com a lógica de extração.
- **Múltiplos Endpoints**: Oferece diferentes "visões" dos dados do documento:
  - Extração de apenas imagens (com URL e Base64).
  - Extração de conteúdo completo (texto e imagens em ordem sequencial).
  - Extração do objeto de documento "bruto" (raw) da API do Google.
- **Processamento Concorrente**: Capaz de processar múltiplos documentos em uma única requisição de forma paralela, utilizando goroutines.
- **Autenticação Segura**: Utiliza o fluxo OAuth 2.0 para autorizar o acesso aos documentos de forma segura.
- **Log Estruturado**: Emite logs em formato JSON utilizando o pacote `slog`, ideal para monitoramento e análise.
- **Estrutura de Projeto Profissional**: O código é organizado seguindo as convenções da comunidade Go (`/cmd`, `/internal`), promovendo manutenibilidade e escalabilidade.

## 📂 Estrutura do Projeto

```
meu-projeto/
├── cmd/
│   └── doc-extractor-api/
│       └── main.go           # Ponto de entrada: configuração e inicialização
├── internal/
│   ├── api/
│   │   ├── handler.go        # Handlers HTTP (a "cola" entre a web e a lógica)
│   │   └── router.go         # Definição das rotas da API
│   ├── extractor/
│   │   └── service.go        # Lógica de negócio principal (como extrair dados)
│   ├── gdocs/
│   │   └── client.go         # Lógica de autenticação e comunicação com a API Google
│   └── model/
│       └── types.go          # Definição das estruturas de dados (structs)
├── Dockerfile                # Receita para construir a imagem Docker
├── .dockerignore             # Arquivos a serem ignorados pelo Docker
├── go.mod
├── go.sum
└── credentials.json          # Arquivo de credenciais da API Google (NÃO versionar em Git público)
└── token.json                # Gerado após a primeira autenticação (NÃO versionar)
```

## 🚀 Começando

Siga os passos abaixo para configurar e executar o projeto localmente.

### Pré-requisitos

- [Go](https://go.dev/doc/install) (versão 1.21 ou superior)
- [Docker](https://docs.docker.com/get-docker/)
- Uma Conta Google
- Um projeto no [Google Cloud Platform](https://console.cloud.google.com/)

### Configuração

1.  **Clone o Repositório** (ou copie os arquivos para uma nova pasta):

    ```bash
    # Exemplo com git
    git clone [https://sua-url-de-repositorio.git](https://sua-url-de-repositorio.git)
    cd seu-projeto
    ```

2.  **Configure as Credenciais do Google Cloud**:
    - No seu projeto do Google Cloud, vá para "APIs e Serviços" > "Biblioteca" e ative a **Google Docs API**.
    - Vá para "APIs e Serviços" > "Credenciais".
    - Clique em "Criar Credenciais" > "ID do cliente OAuth".
    - Selecione o tipo de aplicativo como **"Aplicativo para computador"**.
    - Após a criação, faça o download do arquivo JSON.
    - **Renomeie o arquivo baixado para `credentials.json`** e coloque-o na pasta raiz do projeto.

3.  **Inicialize o Módulo Go**:
    (Se você ainda não o fez, substitua `meu-projeto` pelo nome do seu módulo)

    ```bash
    go mod init meu-projeto
    ```

4.  **Instale as Dependências**:
    ```bash
    go mod tidy
    ```

## 🐳 Executando a Aplicação com Docker

O fluxo de trabalho com Docker é dividido em dois passos principais: uma autenticação única feita localmente, seguida pela execução normal via contêiner.

### Passo 1: Autenticação Inicial (Execução Única)

Devido ao fluxo de autenticação interativo do Google, o primeiro passo precisa ser feito fora do Docker para gerar o arquivo de token (`token.json`). **Você só precisa fazer isso uma vez.**

1.  Execute a aplicação localmente com o Go:
    ```bash
    go run ./cmd/doc-extractor-api
    ```
2.  Siga as instruções no terminal: copie a URL para o navegador, autorize o aplicativo e cole o código de autorização de volta no terminal.
3.  Um arquivo `token.json` será criado na raiz do projeto. Após a sua criação, você pode parar o servidor local (`Ctrl + C`).

### Passo 2: Construindo e Executando com Docker

Esta é a forma padrão de executar a aplicação no dia a dia.

1.  **Construa a Imagem Docker:**
    Este comando lê o `Dockerfile` e empacota sua aplicação em uma imagem chamada `doc-extractor-api`.

    ```bash
    docker build -t doc-extractor-api .
    ```

2.  **Execute o Contêiner Docker:**
    Este comando inicia um contêiner a partir da imagem que acabamos de construir.

    ```bash
    docker run -p 8080:8080 --rm --name my-doc-extractor -v "$(pwd)/token.json:/app/token.json" doc-extractor-api
    ```

    - `-p 8080:8080`: Mapeia a porta 8080 do seu computador para a porta 8080 do contêiner.
    - `--rm`: Remove o contêiner automaticamente quando ele for parado.
    - `--name`: Dá um nome fácil de lembrar para o contêiner em execução.
    - `-v "$(pwd)/token.json:/app/token.json"`: **(A parte mais importante)** Monta o `token.json` da sua máquina local para dentro do contêiner. Isso permite que a aplicação pule a etapa de autenticação interativa. (No Windows CMD, use `%cd%` no lugar de `$(pwd)`).

O servidor estará rodando em `http://localhost:8080` e pronto para receber requisições, sem pedir autorização no terminal.

## 📖 Uso da API

A API aceita múltiplos parâmetros `doc_id` em todos os endpoints para processamento em lote.

---

### 1. Extrair Somente Imagens

Retorna uma lista de todas as imagens encontradas, com suas URLs de acesso e dados em Base64.

- **Método**: `GET`
- **Endpoint**: `/extrair`
- **Exemplo de Requisição**:
  ```http
  http://localhost:8080/extrair?doc_id=ID_DOCUMENTO_1&doc_id=ID_DOCUMENTO_2
  ```
- **Exemplo de Resposta JSON**:
  ```json
  {
    "ID_DOCUMENTO_1": {
      "status": "success",
      "images": [
        {
          "url": "/imagens/ID_DOCUMENTO_1/imagem_1.png",
          "base64": "iVBORw0KGgoAAA..."
        }
      ]
    }
  }
  ```

---

### 2. Extrair Conteúdo Completo (Texto e Imagens)

Retorna uma lista ordenada de blocos de conteúdo (texto e imagem), preservando a estrutura do documento, incluindo conteúdo dentro de tabelas.

- **Método**: `GET`
- **Endpoint**: `/extrair-conteudo-completo`
- **Exemplo de Requisição**:
  ```http
  http://localhost:8080/extrair-conteudo-completo?doc_id=ID_DOCUMENTO_1
  ```
- **Exemplo de Resposta JSON**:
  ```json
  {
    "ID_DOCUMENTO_1": {
      "status": "success",
      "content": [
        {
          "type": "text",
          "content": "Este é o primeiro parágrafo."
        },
        {
          "type": "image",
          "url": "/imagens/ID_DOCUMENTO_1/imagem_1.png",
          "base64": "iVBORw0KGgoAAA..."
        },
        {
          "type": "text",
          "content": "Este texto vem depois da imagem."
        }
      ]
    }
  }
  ```

---

### 3. Extrair Documento Completo (Raw)

Retorna o objeto de documento completo e não processado da API do Google, oferecendo máxima flexibilidade para análise detalhada.

- **Método**: `GET`
- **Endpoint**: `/extrair-documento-completo`
- **Exemplo de Requisição**:
  ```http
  http://localhost:8080/extrair-documento-completo?doc_id=ID_DOCUMENTO_1
  ```
- **Exemplo de Resposta JSON**:
  ```json
  {
    "ID_DOCUMENTO_1": {
      "documentId": "ID_DOCUMENTO_1",
      "title": "Título do Documento",
      "body": {
        "content": [
          // ... estrutura completa e detalhada de parágrafos, tabelas, etc.
        ]
      }
      // ... muitos outros campos da API do Google
    }
  }
  ```

## 📝 Logging

A aplicação utiliza o pacote `log/slog` do Go para gerar logs estruturados em formato JSON no terminal. Isso facilita a análise e integração com sistemas de monitoramento.

## ⚖️ Licença

Este projeto é distribuído sob a licença MIT.
