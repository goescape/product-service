package usecases

import (
	"context"
	"log"
	"product-svc/proto/product"
	repository "product-svc/repository/product"
)

type usecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) *usecase {
	return &usecase{
		repo: repo,
	}
}

var _ ProductUsecase = &usecase{}

type ProductUsecase interface {
	InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error)
	ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error)
}

func (u *usecase) InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error) {
	resp, err := u.repo.InsertProduct(ctx, req)
	if err != nil {
		log.Default().Println("failed to insert product:", err)
		return nil, err
	}

	return resp, nil
}

func (u *usecase) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	resp, err := u.repo.ListProduct(ctx, req)
	if err != nil {
		log.Default().Println("failed to list product:", err)
		return nil, err
	}

	return resp, nil
}
