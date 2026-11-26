package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/auth"
	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/inventory"
)

// Server agrupa los handlers HTTP. No contiene lógica de negocio.
type Server struct {
	service     *inventory.Service
	validator   *auth.Validator
	serviceName string
}

func NewServer(service *inventory.Service, validator *auth.Validator, serviceName string) *Server {
	return &Server{service: service, validator: validator, serviceName: serviceName}
}

// Handler devuelve el router: /health es público, /api/* exige JWT.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)

	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/products", s.listProducts)
	protected.HandleFunc("GET /api/products/{id}", s.getProduct)
	protected.HandleFunc("POST /api/products", s.createProduct)
	protected.HandleFunc("PUT /api/products/{id}", s.updateProduct)
	protected.HandleFunc("DELETE /api/products/{id}", s.deleteProduct)

	mux.Handle("/api/", s.requireAuth(protected))
	return mux
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": s.serviceName,
	})
}

func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.service.ListProducts(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	product, err := s.service.GetProduct(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	input, err := decodeProductInput(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "payload inválido")
		return
	}

	product, err := s.service.CreateProduct(r.Context(), input)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, product)
}

func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request) {
	input, err := decodeProductInput(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "payload inválido")
		return
	}

	product, err := s.service.UpdateProduct(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request) {
	if err := s.service.DeleteProduct(r.Context(), r.PathValue("id")); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeProductInput(r *http.Request) (inventory.ProductInput, error) {
	var input inventory.ProductInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return inventory.ProductInput{}, err
	}
	return input, nil
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, inventory.ErrNotFound):
		writeError(w, http.StatusNotFound, "No encontrado")
	case errors.Is(err, inventory.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
