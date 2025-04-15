package repository

import (
	"context"
	"database/sql"
	"log"
	"product-svc/proto/product"
	"strings"

	"github.com/lib/pq"
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
	ReduceProduct(ctx context.Context, req *product.ReduceProductsRequest) (*product.ReduceProductsResponse, error)
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
		args      = make([]interface{}, 0)
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
	`
	if req.ProductIds == "" {
		query += `LIMIT $1 OFFSET $2`
		args = append(args, req.Limit, (req.Page-1)*req.Limit)
	} else {
		// example: xx1, xx2, xx3
		query += `WHERE id = ANY($1)`
		args = append(args, pq.Array(strings.Split(req.ProductIds, ",")))
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
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

func (s *store) ReduceProduct(ctx context.Context, req *product.ReduceProductsRequest) (*product.ReduceProductsResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Default().Println("failed to start db transaction:", err)
		return nil, err
	}
	defer tx.Rollback()

	for _, item := range req.Items {
		query := `UPDATE products SET qty = qty - $1 WHERE id = $2 AND qty >= $1`
		res, err := tx.ExecContext(ctx, query, item.Qty, item.ProductId)
		if err != nil {
			log.Default().Println("failed to update product:", err)
			return nil, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Default().Println("failed to get rows affected:", err)
			return nil, err
		}

		if rowsAffected == 0 {
			return &product.ReduceProductsResponse{
				Msg: "Product quantity not sufficient",
			}, nil
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Default().Println("failed to commit transaction:", err)
		return nil, err
	}

	return &product.ReduceProductsResponse{
		Msg: "Product quantity reduced successfully",
	}, nil
}
