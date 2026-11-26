package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/auth"
	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/inventory"
	"github.com/golang-jwt/jwt/v5"
)

// fakeRepository reemplaza a PostgreSQL en los tests de la capa HTTP.
type fakeRepository struct {
	products map[string]inventory.Product
	created  inventory.ProductInput
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{products: make(map[string]inventory.Product)}
}

func (f *fakeRepository) ListProducts(ctx context.Context) ([]inventory.Product, error) {
	products := make([]inventory.Product, 0, len(f.products))
	for _, product := range f.products {
		products = append(products, product)
	}
	return products, nil
}

func (f *fakeRepository) GetProduct(ctx context.Context, id string) (*inventory.Product, error) {
	product, ok := f.products[id]
	if !ok {
		return nil, inventory.ErrNotFound
	}
	return &product, nil
}

func (f *fakeRepository) CreateProduct(ctx context.Context, input inventory.ProductInput) (*inventory.Product, error) {
	f.created = input
	product := inventory.Product{
		ID:            "new-id",
		SKU:           input.SKU,
		Name:          input.Name,
		Unit:          input.Unit,
		PurchasePrice: input.PurchasePrice,
		SalePrice:     input.SalePrice,
		Active:        input.Active,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	f.products[product.ID] = product
	return &product, nil
}

func (f *fakeRepository) UpdateProduct(ctx context.Context, id string, input inventory.ProductInput) (*inventory.Product, error) {
	product, ok := f.products[id]
	if !ok {
		return nil, inventory.ErrNotFound
	}
	product.SKU = input.SKU
	product.Name = input.Name
	product.Unit = input.Unit
	product.PurchasePrice = input.PurchasePrice
	product.SalePrice = input.SalePrice
	product.Active = input.Active
	product.UpdatedAt = time.Now()
	f.products[id] = product
	return &product, nil
}

func (f *fakeRepository) DeleteProduct(ctx context.Context, id string) error {
	if _, ok := f.products[id]; !ok {
		return inventory.ErrNotFound
	}
	delete(f.products, id)
	return nil
}

func newTestHandler(repo inventory.Repository) http.Handler {
	return NewServer(
		inventory.NewService(repo),
		auth.NewValidator("test_secret"),
		"inventory-service-go",
	).Handler()
}

func signedToken(t *testing.T) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub":   "user-1",
		"email": "admin@inventory.local",
		"name":  "Admin",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("test_secret"))
	if err != nil {
		t.Fatalf("no se pudo firmar el token de prueba: %v", err)
	}
	return signed
}

func perform(handler http.Handler, method, path, token, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

const validProductBody = `{"sku":"LAP-001","name":"Laptop","unit":"UND","purchasePrice":900,"salePrice":1200,"active":true}`

func TestHealthIsPublic(t *testing.T) {
	recorder := perform(newTestHandler(newFakeRepository()), http.MethodGet, "/health", "", "")

	if recorder.Code != http.StatusOK {
		t.Fatalf("esperaba 200, obtuve %d", recorder.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("respuesta no JSON: %v", err)
	}
	if body["status"] != "ok" || body["service"] != "inventory-service-go" {
		t.Errorf("health incorrecto: %v", body)
	}
}

func TestProductsRequireValidJWT(t *testing.T) {
	handler := newTestHandler(newFakeRepository())

	if recorder := perform(handler, http.MethodGet, "/api/products", "", ""); recorder.Code != http.StatusUnauthorized {
		t.Errorf("sin token: esperaba 401, obtuve %d", recorder.Code)
	}
	if recorder := perform(handler, http.MethodGet, "/api/products", "token-invalido", ""); recorder.Code != http.StatusUnauthorized {
		t.Errorf("token inválido: esperaba 401, obtuve %d", recorder.Code)
	}
}

func TestListProductsWithValidJWT(t *testing.T) {
	repo := newFakeRepository()
	repo.products["p-1"] = inventory.Product{ID: "p-1", SKU: "LAP-001", Name: "Laptop"}

	recorder := perform(newTestHandler(repo), http.MethodGet, "/api/products", signedToken(t), "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("esperaba 200, obtuve %d", recorder.Code)
	}

	var products []inventory.Product
	if err := json.NewDecoder(recorder.Body).Decode(&products); err != nil {
		t.Fatalf("respuesta no JSON: %v", err)
	}
	if len(products) != 1 || products[0].SKU != "LAP-001" {
		t.Errorf("lista incorrecta: %+v", products)
	}
}

func TestGetProductNotFound(t *testing.T) {
	recorder := perform(newTestHandler(newFakeRepository()), http.MethodGet, "/api/products/no-existe", signedToken(t), "")

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("esperaba 404, obtuve %d", recorder.Code)
	}
}

func TestCreateProduct(t *testing.T) {
	repo := newFakeRepository()
	recorder := perform(newTestHandler(repo), http.MethodPost, "/api/products", signedToken(t), validProductBody)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d: %s", recorder.Code, recorder.Body.String())
	}
	if repo.created.SKU != "LAP-001" || repo.created.Unit != "UND" {
		t.Errorf("Repository recibió datos incorrectos: %+v", repo.created)
	}

	var product inventory.Product
	if err := json.NewDecoder(recorder.Body).Decode(&product); err != nil {
		t.Fatalf("respuesta no JSON: %v", err)
	}
	if product.ID != "new-id" || product.Name != "Laptop" {
		t.Errorf("producto creado incorrecto: %+v", product)
	}
}

func TestCreateProductRejectsExtraFields(t *testing.T) {
	repo := newFakeRepository()
	body := `{"sku":"LAP-001","name":"Laptop","unit":"UND","roles":["ADMIN"]}`

	recorder := perform(newTestHandler(repo), http.MethodPost, "/api/products", signedToken(t), body)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d", recorder.Code)
	}
	if repo.created.SKU != "" {
		t.Error("no se debió llamar a CreateProduct con un payload inválido")
	}
}

func TestCreateProductRejectsInvalidPayload(t *testing.T) {
	handler := newTestHandler(newFakeRepository())
	recorder := perform(handler, http.MethodPost, "/api/products", signedToken(t), `{"name":"Laptop","unit":"UND"}`)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d", recorder.Code)
	}
}

func TestUpdateAndDeleteProduct(t *testing.T) {
	repo := newFakeRepository()
	repo.products["p-1"] = inventory.Product{ID: "p-1", SKU: "VIEJO", Name: "Viejo", Unit: "UND"}
	handler := newTestHandler(repo)

	updated := perform(handler, http.MethodPut, "/api/products/p-1", signedToken(t), validProductBody)
	if updated.Code != http.StatusOK {
		t.Fatalf("PUT: esperaba 200, obtuve %d", updated.Code)
	}
	if repo.products["p-1"].SKU != "LAP-001" {
		t.Errorf("producto no actualizado: %+v", repo.products["p-1"])
	}

	deleted := perform(handler, http.MethodDelete, "/api/products/p-1", signedToken(t), "")
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("DELETE: esperaba 204, obtuve %d", deleted.Code)
	}
	if recorder := perform(handler, http.MethodDelete, "/api/products/p-1", signedToken(t), ""); recorder.Code != http.StatusNotFound {
		t.Errorf("DELETE repetido: esperaba 404, obtuve %d", recorder.Code)
	}
}
