package handler

import (
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockContextEntity est le mock de ContextEntityInterface
type MockContextEntity struct {
	mock.Mock
}

func (m *MockContextEntity) UploadImage(imageURL, publicID string) (*uploader.UploadResult, error) {
	args := m.Called(imageURL, publicID)
	return args.Get(0).(*uploader.UploadResult), args.Error(1)
}

func (m *MockContextEntity) UploadVideo(videoURL, publicID string) (*uploader.UploadResult, error) {
	args := m.Called(videoURL, publicID)
	return args.Get(0).(*uploader.UploadResult), args.Error(1)
}

func TestUploadVideoHandlerGin(t *testing.T) {
	// Crée une instance du mock
	mockEntity := new(MockContextEntity)

	// Définit ce que le mock doit retourner
	mockResult := &uploader.UploadResult{
		SecureURL: "http://example.com/video.mp4",
		PublicID:  "sample-video-id",
	}
	mockEntity.On("UploadVideo", "http://example.com/video.mp4", "sample-video-id").Return(mockResult, nil)

	// Crée un routeur Gin et attache le handler
	r := gin.Default()
	r.GET("/upload_video", UploadVideoHandlerGin(mockEntity))

	// Crée une requête HTTP simulée
	req, _ := http.NewRequest("GET", "/upload_video?url=http://example.com/video.mp4&id=sample-video-id", nil)
	w := httptest.NewRecorder()

	// Effectue la requête
	r.ServeHTTP(w, req)

	// Vérifie le code de statut HTTP
	assert.Equal(t, http.StatusOK, w.Code)

	// Vérifie le corps de la réponse
	expectedResponse := `{"StatusCode":200,"Message":"success","Data":{"data":{"SecureURL":"http://example.com/video.mp4","PublicID":"sample-video-id"}}}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	// Vérifie que le mock a été appelé avec les bons arguments
	mockEntity.AssertCalled(t, "UploadVideo", "http://example.com/video.mp4", "sample-video-id")
}

func TestUploadImageHandler(t *testing.T) {
	// Crée une instance du mock
	mockEntity := new(MockContextEntity)

	// Définit ce que le mock doit retourner
	mockResult := &uploader.UploadResult{
		SecureURL: "http://example.com/image.jpg",
	}
	mockEntity.On("UploadImage", "http://example.com/image.jpg", "sample-image-id").Return(mockResult, nil)

	// Crée un routeur HTTP et attache le handler
	handler := UploadImageHandler(mockEntity)

	// Crée une requête HTTP simulée avec les paramètres nécessaires
	req, err := http.NewRequest("GET", "/upload_image?url=http://example.com/image.jpg&id=sample-image-id", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rec := httptest.NewRecorder()

	// Effectue la requête
	handler.ServeHTTP(rec, req)

	// Vérifie le code de statut HTTP
	assert.Equal(t, http.StatusOK, rec.Code)

	// Vérifie le corps de la réponse
	expectedResponse := "Image uploaded successfully! Delivery URL: http://example.com/image.jpg\n"
	assert.Equal(t, expectedResponse, rec.Body.String())

	// Vérifie que le mock a été appelé avec les bons arguments
	mockEntity.AssertCalled(t, "UploadImage", "http://example.com/image.jpg", "sample-image-id")
}

//func TestUploadImageHandler_BadRequest(t *testing.T) {
//	// Crée une instance du mock
//	mockEntity := new(MockContextEntity)
//
//	// Crée un routeur HTTP et attache le handler
//	handler := UploadImageHandler(mockEntity)
//
//	// Crée une requête HTTP simulée sans le paramètre URL
//	req, err := http.NewRequest("GET", "/upload_image?id=sample-image-id", nil)
//	if err != nil {
//		t.Fatalf("Failed to create request: %v", err)
//	}
//	rec := httptest.NewRecorder()
//
//	// Effectue la requête
//	handler.ServeHTTP(rec, req)
//
//	// Vérifie le code de statut HTTP
//	assert.Equal(t, http.StatusBadRequest, rec.Code)
//
//	// Vérifie le corps de la réponse
//	expectedResponse := "URL parameter is required"
//	assert.Equal(t, expectedResponse, rec.Body.String())
//}
