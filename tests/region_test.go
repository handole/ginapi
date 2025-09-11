package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ginapi/controllers"
	"ginapi/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetRegionsUnauthorized(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	regionController := controllers.RegionController(mt.Client)

	r := gin.Default()
	r.GET("/regions", middlewares.AuthMiddleware(), regionController.GetRegions)

	req, _ := http.NewRequest("GET", "/regions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}
