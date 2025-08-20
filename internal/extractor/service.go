package extractor

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcokist/googledocs-extractor/internal/model"
	"google.golang.org/api/docs/v1"
)

type Service struct {
	docsService *docs.Service
	logger      *slog.Logger
}

func NewService(ds *docs.Service, logger *slog.Logger) *Service {
	return &Service{
		docsService: ds,
		logger:      logger,
	}
}

func (s *Service) ExtractImages(documentId string) ([]model.ImageData, error) {
	doc, err := s.docsService.Documents.Get(documentId).Do()
	if err != nil {
		return nil, fmt.Errorf("não foi possível obter o documento: %w", err)
	}
	outputDir := filepath.Join("imagens_extraidas", documentId)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}

	var savedImages []model.ImageData
	var imageId int = 1

	inlineObjectsPtr := make(map[string]*docs.InlineObject)
	for k, v := range doc.InlineObjects {
		inlineObjectsPtr[k] = &v
	}

	contentBlocks, _ := s.processStructuralElements(doc.Body.Content, inlineObjectsPtr, outputDir, documentId, &imageId)
	for _, block := range contentBlocks {
		if block.Type == "image" {
			savedImages = append(savedImages, model.ImageData{
				URL:    block.URL,
				Base64: block.Base64,
			})
		}
	}
	return savedImages, nil
}

func (s *Service) ExtractFullContent(documentId string) ([]model.ContentBlock, error) {
	doc, err := s.docsService.Documents.Get(documentId).Do()
	if err != nil {
		return nil, fmt.Errorf("não foi possível obter o documento: %w", err)
	}
	outputDir := filepath.Join("imagens_extraidas", documentId)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}
	var imageId int = 1
	inlineObjectsPtr := make(map[string]*docs.InlineObject)
	for k, v := range doc.InlineObjects {
		inlineObjectsPtr[k] = &v
	}
	return s.processStructuralElements(doc.Body.Content, inlineObjectsPtr, outputDir, documentId, &imageId)
}

func (s *Service) ExtractRawDocument(documentId string) (*docs.Document, error) {
	return s.docsService.Documents.Get(documentId).Do()
}

func (s *Service) processStructuralElements(elements []*docs.StructuralElement, inlineObjects map[string]*docs.InlineObject, outputDir, documentId string, imageId *int) ([]model.ContentBlock, error) {
	var contentBlocks []model.ContentBlock
	for _, element := range elements {
		if element.Paragraph != nil {
			var currentText strings.Builder
			for _, pElem := range element.Paragraph.Elements {
				if pElem.TextRun != nil {
					currentText.WriteString(pElem.TextRun.Content)
				} else if pElem.InlineObjectElement != nil {
					if currentText.Len() > 0 {
						contentBlocks = append(contentBlocks, model.ContentBlock{Type: "text", Content: strings.TrimSuffix(currentText.String(), "\n")})
						currentText.Reset()
					}
					inlineObjectId := pElem.InlineObjectElement.InlineObjectId
					inlineObject := inlineObjects[inlineObjectId]
					if inlineObject.InlineObjectProperties.EmbeddedObject.ImageProperties != nil {
						imgURL := inlineObject.InlineObjectProperties.EmbeddedObject.ImageProperties.ContentUri
						fileName := fmt.Sprintf("imagem_%d.png", *imageId)
						filePath := filepath.Join(outputDir, fileName)
						_ = downloadFile(imgURL, filePath)
						base64String, _ := encodeFileToBase64(filePath)
						publicURL := filepath.ToSlash(filepath.Join("/imagens", documentId, fileName))
						contentBlocks = append(contentBlocks, model.ContentBlock{Type: "image", URL: publicURL, Base64: base64String})
						*imageId++
					}
				}
			}
			if currentText.Len() > 0 {
				contentBlocks = append(contentBlocks, model.ContentBlock{Type: "text", Content: strings.TrimSuffix(currentText.String(), "\n")})
			}
		} else if element.Table != nil {
			for _, row := range element.Table.TableRows {
				for _, cell := range row.TableCells {
					cellContent, err := s.processStructuralElements(cell.Content, inlineObjects, outputDir, documentId, imageId)
					if err != nil {
						return nil, err
					}
					contentBlocks = append(contentBlocks, cellContent...)
				}
			}
		}
	}
	return contentBlocks, nil
}

func downloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func encodeFileToBase64(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
