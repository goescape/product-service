package repository

import (
	"context"
	"database/sql"
	"product-svc/proto/product"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

var _ ProductRepository = &store{}

type ProductRepository interface {
	InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error)
	ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error)
}

func (s *store) InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error) {
	// Implement the logic to insert a product into the database
	// and return the response.
	return &product.ProductInsertResponse{
		Msg: "Product inserted successfully",
	}, nil
}

func (s *store) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	// Implement the logic to list products from the database
	// and return the response.
	items := make([]*product.Product, 0)
	meta := &product.Meta{
		TotalData:   0,
		TotalPage:   1,
		CurrentPage: 1,
		Limit:       10,
	}
	return &product.ListProductResponse{
		Items: items,
		Meta:  meta,
	}, nil
}
