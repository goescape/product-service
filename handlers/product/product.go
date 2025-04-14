package handlers

import (
	"context"
	"product-svc/proto/product"
	usecases "product-svc/usecases/product"
)

type Handler struct {
	usecase usecases.ProductUsecase
	product.UnimplementedProductServiceServer
}

func NewHandler(uc usecases.ProductUsecase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

func (h *Handler) InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error) {
	return h.usecase.InsertProduct(ctx, req)
}

func (h *Handler) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	return h.usecase.ListProduct(ctx, req)
}
