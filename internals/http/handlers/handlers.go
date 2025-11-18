package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Handler struct {
	log     *slog.Logger
	client  pizzalndv1.PizzaLandClient
	timeout time.Duration
}

func NewHandler(log *slog.Logger, conn *grpc.ClientConn, timeout time.Duration) *Handler {
	return &Handler{
		log:     log,
		client:  pizzalndv1.NewPizzaLandClient(conn),
		timeout: timeout,
	}
}

func (h *Handler) ListPizzas(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	// Parse query parameters
	page, _ := strconv.ParseUint(r.URL.Query().Get("page"), 10, 32)
	offsetParam, _ := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 32)
	limit, _ := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 32)
	categoryID, _ := strconv.ParseUint(r.URL.Query().Get("category"), 10, 32)
	categoryName := r.URL.Query().Get("category_name")
	search := r.URL.Query().Get("search")

	// Calculate offset - support both page and offset parameters
	var offset uint64
	if offsetParam > 0 {
		offset = offsetParam
	} else if page > 0 {
		// Default limit if not specified with page
		if limit == 0 {
			limit = 12
		}
		offset = (page - 1) * limit
	}

	// Default values
	if limit == 0 {
		limit = 12
	}
	// Validate limit - must be one of [12, 24, 36, 48] for gRPC, but allow any for flexibility
	// We'll use 12 as default if invalid
	if limit != 12 && limit != 24 && limit != 36 && limit != 48 {
		// Round to nearest valid value or use 12
		if limit < 12 {
			limit = 12
		} else if limit <= 18 {
			limit = 12
		} else if limit <= 30 {
			limit = 24
		} else if limit <= 42 {
			limit = 36
		} else {
			limit = 48
		}
	}

	var req *pizzalndv1.ListRequest
	if categoryID > 0 {
		req = &pizzalndv1.ListRequest{
			CategoryId: wrapperspb.UInt32(uint32(categoryID)),
			Offset:     uint32(offset),
			Limit:      uint32(limit),
		}
	} else if categoryName != "" {
		req = &pizzalndv1.ListRequest{
			CategoryName: wrapperspb.String(categoryName),
			Offset:       uint32(offset),
			Limit:        uint32(limit),
		}
	} else {
		req = &pizzalndv1.ListRequest{
			Offset: uint32(offset),
			Limit:  uint32(limit),
		}
	}

	resp, err := h.client.List(ctx, req)
	if err != nil {
		h.handleError(w, err, "failed to list pizzas")
		return
	}

	// Filter by search if provided
	pizzas := resp.GetPizza()
	if search != "" {
		filtered := make([]*pizzalndv1.PizzaProperties, 0)
		searchLower := strings.ToLower(search)
		for _, pizza := range pizzas {
			if strings.Contains(strings.ToLower(pizza.GetTitle()), searchLower) {
				filtered = append(filtered, pizza)
			}
		}
		pizzas = filtered
	}

	// Convert to frontend format
	result := make([]map[string]interface{}, 0, len(pizzas))
	for _, pizza := range pizzas {
		result = append(result, h.pizzaToJSON(pizza))
	}

	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) GetPizza(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/api/pizzas/")
	
	var req *pizzalndv1.GetRequest
	if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		req = &pizzalndv1.GetRequest{
			Identifier: &pizzalndv1.GetRequest_PizzaId{PizzaId: id},
		}
	} else {
		req = &pizzalndv1.GetRequest{
			Identifier: &pizzalndv1.GetRequest_PizzaName{PizzaName: idStr},
		}
	}

	resp, err := h.client.Get(ctx, req)
	if err != nil {
		h.handleError(w, err, "failed to get pizza")
		return
	}

	h.writeJSON(w, http.StatusOK, h.pizzaToJSON(resp.GetPizza()))
}

func (h *Handler) SavePizza(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	var pizza pizzalndv1.PizzaProperties
	if err := json.NewDecoder(r.Body).Decode(&pizza); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req := &pizzalndv1.SaveRequest{
		Pizza: &pizza,
	}

	resp, err := h.client.Save(ctx, req)
	if err != nil {
		h.handleError(w, err, "failed to save pizza")
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"pizza_id": resp.GetPizzaId(),
	})
}

func (h *Handler) UpdatePizza(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/api/pizzas/")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid pizza id")
		return
	}

	var pizza pizzalndv1.PizzaProperties
	if err := json.NewDecoder(r.Body).Decode(&pizza); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req := &pizzalndv1.UpdateRequest{
		Identifier: &pizzalndv1.UpdateRequest_Id{Id: id},
	}

	if pizza.Category > 0 {
		req.Category = wrapperspb.UInt32(pizza.Category)
	}
	if pizza.Title != "" {
		req.Title = wrapperspb.String(pizza.Title)
	}
	if pizza.Description != nil {
		req.Description = wrapperspb.String(pizza.Description.Value)
	}
	if len(pizza.Types) > 0 {
		req.Types = pizza.Types
	}
	if pizza.Price > 0 {
		req.Price = wrapperspb.Float(pizza.Price)
	}
	if len(pizza.Sizes) > 0 {
		req.Sizes = pizza.Sizes
	}
	if pizza.Rating > 0 {
		req.Rating = wrapperspb.UInt32(pizza.Rating)
	}
	if pizza.ImageUrl != "" {
		req.ImageUrl = wrapperspb.String(pizza.ImageUrl)
	}

	resp, err := h.client.Update(ctx, req)
	if err != nil {
		h.handleError(w, err, "failed to update pizza")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": resp.GetSuccess(),
	})
}

func (h *Handler) RemovePizza(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/api/pizzas/")
	
	var req *pizzalndv1.RemoveRequest
	if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		req = &pizzalndv1.RemoveRequest{
			Identifier: &pizzalndv1.RemoveRequest_PizzaId{PizzaId: id},
		}
	} else {
		req = &pizzalndv1.RemoveRequest{
			Identifier: &pizzalndv1.RemoveRequest_PizzaName{PizzaName: idStr},
		}
	}

	resp, err := h.client.Remove(ctx, req)
	if err != nil {
		h.handleError(w, err, "failed to remove pizza")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": resp.GetSuccess(),
	})
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	// For now, return a simple list. In a real app, you'd want a GetCategoryList endpoint
	// that returns categories, not pizzas. For now, we'll return a mock list.
	categories := []map[string]interface{}{
		{"id": 0, "name": "Все"},
		{"id": 1, "name": "Мясные"},
		{"id": 2, "name": "Вегетарианская"},
		{"id": 3, "name": "Гриль"},
		{"id": 4, "name": "Острые"},
		{"id": 5, "name": "Закрытые"},
	}

	h.writeJSON(w, http.StatusOK, categories)
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	idStr = strings.TrimSuffix(idStr, "/pizzas")
	
	var req *pizzalndv1.GetCategoryRequest
	if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
		req = &pizzalndv1.GetCategoryRequest{
			Identifier: &pizzalndv1.GetCategoryRequest_CategoryId{CategoryId: uint32(id)},
		}
	} else {
		req = &pizzalndv1.GetCategoryRequest{
			Identifier: &pizzalndv1.GetCategoryRequest_CategoryName{CategoryName: idStr},
		}
	}

	resp, err := h.client.GetCategory(ctx, req)
	if err != nil {
		h.handleError(w, err, "failed to get category")
		return
	}

	h.writeJSON(w, http.StatusOK, resp.GetCategory())
}

func (h *Handler) GetCategoryPizzas(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	idStr = strings.TrimSuffix(idStr, "/pizzas")
	
	offset, _ := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 32)
	limit, _ := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 32)
	if limit == 0 {
		limit = 12
	}

	var req *pizzalndv1.GetCategoryListRequest
	if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
		req = &pizzalndv1.GetCategoryListRequest{
			Identifier: &pizzalndv1.GetCategoryListRequest_CategoryId{CategoryId: uint32(id)},
			Offset:     uint32(offset),
			Limit:      uint32(limit),
		}
	} else {
		req = &pizzalndv1.GetCategoryListRequest{
			Identifier: &pizzalndv1.GetCategoryListRequest_CategoryName{CategoryName: idStr},
			Offset:     uint32(offset),
			Limit:      uint32(limit),
		}
	}

	resp, err := h.client.GetCategoryList(ctx, req)
	if err != nil {
		h.handleError(w, err, "failed to get category pizzas")
		return
	}

	// Convert to frontend format
	pizzas := resp.GetPizza().GetPizza()
	result := make([]map[string]interface{}, 0, len(pizzas))
	for _, pizza := range pizzas {
		result = append(result, h.pizzaToJSON(pizza))
	}

	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to encode response", slog.Any("error", err))
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func (h *Handler) pizzaToJSON(pizza *pizzalndv1.PizzaProperties) map[string]interface{} {
	types := pizza.GetTypes()
	if types == nil {
		types = []int32{}
	}
	sizes := pizza.GetSizes()
	if sizes == nil {
		sizes = []int32{}
	}
	
	result := map[string]interface{}{
		"title":    pizza.GetTitle(),
		"price":    pizza.GetPrice(),
		"imageUrl": pizza.GetImageUrl(),
		"types":    types,
		"sizes":    sizes,
		"category": pizza.GetCategory(),
		"rating":   pizza.GetRating(),
	}
	if pizza.GetId() != nil {
		result["id"] = pizza.GetId().GetValue()
	}
	if pizza.GetDescription() != nil {
		result["description"] = pizza.GetDescription().GetValue()
	}
	return result
}

func (h *Handler) handleError(w http.ResponseWriter, err error, defaultMsg string) {
	st, ok := status.FromError(err)
	if !ok {
		h.writeError(w, http.StatusInternalServerError, defaultMsg)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		h.writeError(w, http.StatusNotFound, st.Message())
	case codes.InvalidArgument:
		h.writeError(w, http.StatusBadRequest, st.Message())
	case codes.AlreadyExists:
		h.writeError(w, http.StatusConflict, st.Message())
	default:
		h.writeError(w, http.StatusInternalServerError, defaultMsg)
	}
}

