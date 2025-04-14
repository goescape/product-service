package repository

import (
	"context"
	"database/sql"
	"log"
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
	tx, err := s.db.Begin()
	if err != nil {
		log.Default().Println("failed to start db transaction:", err)
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO products (user_id, name, description, price, qty) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var productId string

	err = tx.QueryRowContext(ctx, query, req.UserId, req.Name, req.Description, req.Price, req.Qty).Scan(&productId)
	if err != nil {
		log.Default().Println("failed to insert product:", err)
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Default().Println("failed to commit transaction:", err)
		return nil, err
	}

	return &product.ProductInsertResponse{
		Msg: "Product inserted successfully: " + productId,
	}, nil
}

func (s *store) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	var (
		totalData uint32
		resp      = new(product.ListProductResponse)
	)
	resp.Items = make([]*product.Product, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			user_id,
			name,
			description,
			price,
			qty
		FROM products
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		log.Default().Println("failed to query products:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item product.Product
		if err := rows.Scan(&totalData, &item.Id, &item.UserId, &item.Name, &item.Description, &item.Price, &item.Qty); err != nil {
			log.Default().Println("failed to scan product:", err)
			return nil, err
		}
		resp.Items = append(resp.Items, &item)
	}
	if err := rows.Err(); err != nil {
		log.Default().Println("failed to iterate products:", err)
		return nil, err
	}

	if len(resp.Items) == 0 {
		return &product.ListProductResponse{
			Meta: &product.Meta{
				TotalData:   0,
				TotalPage:   1,
				CurrentPage: req.Page,
				Limit:       req.Limit,
			},
			Items: resp.Items,
		}, nil
	}

	// counting total page
	totalPage := totalData / req.Limit
	if totalData%req.Limit != 0 { // if there is a remainder, add 1 to totalPage
		totalPage++
	}

	resp.Meta = &product.Meta{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
		Limit:       req.Limit,
	}

	return resp, nil
}
