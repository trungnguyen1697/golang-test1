package repository

import (
	"database/sql"
	"encoding/json"
	"golang-test1/app/model"
	"golang-test1/platform/database"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type productRepository struct {
	db *database.DB
}

func (p productRepository) Create(product *model.Product, categoryIDs []uuid.UUID) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Convert attributes to JSON
	attributesJSON, err := json.Marshal(product.Attributes)
	if err != nil {
		return err
	}

	// Insert product
	query := `
		INSERT INTO products (id, sku, name, description, price, sale_price, cost_price, 
		                     stock_quantity, status, attributes, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err = tx.Exec(
		query,
		product.ID,
		product.SKU,
		product.Name,
		product.Description,
		product.Price,
		product.SalePrice,
		product.CostPrice,
		product.StockQuantity,
		product.Status,
		attributesJSON,
		product.CreatedAt,
		product.UpdatedAt,
		false,
	)

	if err != nil {
		return err
	}

	// Insert product categories
	for _, categoryID := range categoryIDs {
		_, err = tx.Exec(
			"INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)",
			product.ID,
			categoryID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
func (r *productRepository) GetByID(id uuid.UUID) (*model.Product, error) {
	// Define a temporary struct with nullable fields
	type productScan struct {
		ID            uuid.UUID       `db:"id"`
		SKU           string          `db:"sku"`
		Name          string          `db:"name"`
		Description   sql.NullString  `db:"description"`
		Price         float64         `db:"price"`
		SalePrice     sql.NullFloat64 `db:"sale_price"`
		CostPrice     sql.NullFloat64 `db:"cost_price"`
		StockQuantity int             `db:"stock_quantity"`
		Status        string          `db:"status"`
		Attributes    string          `db:"attributes"`
		CreatedAt     time.Time       `db:"created_at"`
		UpdatedAt     time.Time       `db:"updated_at"`
	}

	var scanProduct productScan

	query := `
        SELECT p.id, p.sku, p.name, p.description, p.price, p.sale_price, p.cost_price, 
               p.stock_quantity, p.status, to_json(p.attributes) as attributes, 
               p.created_at, p.updated_at
        FROM products p
        WHERE p.id = $1
    `

	err := r.db.Get(&scanProduct, query, id)
	if err != nil {
		return nil, err
	}

	// Convert to model.Product
	product := model.Product{
		ID:            scanProduct.ID,
		SKU:           scanProduct.SKU,
		Name:          scanProduct.Name,
		Description:   scanProduct.Description.String,
		Price:         scanProduct.Price,
		StockQuantity: scanProduct.StockQuantity,
		Status:        scanProduct.Status,
		CreatedAt:     scanProduct.CreatedAt,
		UpdatedAt:     scanProduct.UpdatedAt,
	}

	// Handle nullable fields
	if scanProduct.SalePrice.Valid {
		product.SalePrice = scanProduct.SalePrice.Float64
	}

	if scanProduct.CostPrice.Valid {
		product.CostPrice = scanProduct.CostPrice.Float64
	}

	// Parse attributes JSON
	if scanProduct.Attributes != "" {
		if err := json.Unmarshal([]byte(scanProduct.Attributes), &product.Attributes); err != nil {
			return nil, err
		}
	}

	// Get categories
	categories, err := r.GetCategories(id)
	if err != nil {
		return nil, err
	}
	product.Categories = categories

	return &product, nil
}

func (r *productRepository) Update(product *model.Product, categoryIDs []uuid.UUID) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Convert attributes to JSON
	attributesJSON, err := json.Marshal(product.Attributes)
	if err != nil {
		return err
	}

	product.UpdatedAt = time.Now()

	// Update product
	query := `
		UPDATE products
		SET sku = $1, name = $2, description = $3, price = $4, sale_price = $5, 
		    cost_price = $6, stock_quantity = $7, status = $8, attributes = $9, updated_at = $10
		WHERE id = $11
	`

	_, err = tx.Exec(
		query,
		product.SKU,
		product.Name,
		product.Description,
		product.Price,
		product.SalePrice,
		product.CostPrice,
		product.StockQuantity,
		product.Status,
		attributesJSON,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return err
	}

	// Delete existing product categories
	_, err = tx.Exec("DELETE FROM product_categories WHERE product_id = $1", product.ID)
	if err != nil {
		return err
	}

	// Insert new product categories
	for _, categoryID := range categoryIDs {
		_, err = tx.Exec(
			"INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)",
			product.ID,
			categoryID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.QueryRow("DELETE FROM products WHERE id = $1 RETURNING id", id).Scan(&id)
}

func (r *productRepository) List(offset, limit int, search string, categoryID *uuid.UUID, status string) ([]model.Product, int, error) {
	var total int

	// Define a scan struct with nullable fields
	type productScan struct {
		ID            uuid.UUID       `db:"id"`
		SKU           string          `db:"sku"`
		Name          string          `db:"name"`
		Description   sql.NullString  `db:"description"`
		Price         float64         `db:"price"`
		SalePrice     sql.NullFloat64 `db:"sale_price"`
		CostPrice     sql.NullFloat64 `db:"cost_price"`
		StockQuantity int             `db:"stock_quantity"`
		Status        string          `db:"status"`
		Attributes    string          `db:"attributes"`
		CreatedAt     time.Time       `db:"created_at"`
		UpdatedAt     time.Time       `db:"updated_at"`
	}

	var scanProducts []productScan

	// Base queries
	countQuery := `SELECT COUNT(*) FROM products p`
	listQuery := `
		SELECT p.id, p.sku, p.name, p.description, p.price, p.sale_price, p.cost_price, 
		       p.stock_quantity, p.status, to_json(p.attributes) as attributes, 
		       p.created_at, p.updated_at
		FROM products p
	`

	// Build WHERE clause
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause += " WHERE (p.name ILIKE $" + strconv.Itoa(argIndex) + " OR p.sku ILIKE $" + strconv.Itoa(argIndex) + " OR p.description ILIKE $" + strconv.Itoa(argIndex) + ")"
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if categoryID != nil {
		if whereClause == "" {
			whereClause += " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "p.id IN (SELECT product_id FROM product_categories WHERE category_id = $" + strconv.Itoa(argIndex) + ")"
		args = append(args, *categoryID)
		argIndex++
	}

	if status != "" {
		if whereClause == "" {
			whereClause += " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "p.status = $" + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	// Complete queries
	countQuery += whereClause
	listQuery += whereClause + " ORDER BY p.created_at DESC LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	// Get total count
	err := r.db.Get(&total, countQuery, args[:argIndex-1]...)
	if err != nil {
		return nil, 0, err
	}

	// Get products
	err = r.db.Select(&scanProducts, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Convert scanProducts to model.Product
	products := make([]model.Product, len(scanProducts))
	for i, scanProduct := range scanProducts {
		products[i] = model.Product{
			ID:            scanProduct.ID,
			SKU:           scanProduct.SKU,
			Name:          scanProduct.Name,
			Description:   scanProduct.Description.String,
			Price:         scanProduct.Price,
			StockQuantity: scanProduct.StockQuantity,
			Status:        scanProduct.Status,
			CreatedAt:     scanProduct.CreatedAt,
			UpdatedAt:     scanProduct.UpdatedAt,
		}

		// Handle nullable fields
		if scanProduct.SalePrice.Valid {
			products[i].SalePrice = scanProduct.SalePrice.Float64
		}

		if scanProduct.CostPrice.Valid {
			products[i].CostPrice = scanProduct.CostPrice.Float64
		}

		// Parse attributes JSON
		if scanProduct.Attributes != "" {
			if err := json.Unmarshal([]byte(scanProduct.Attributes), &products[i].Attributes); err != nil {
				return nil, 0, err
			}
		}
	}

	// Get categories for each product
	for i := range products {
		categories, err := r.GetCategories(products[i].ID)
		if err != nil {
			return nil, 0, err
		}
		products[i].Categories = categories
	}

	return products, total, nil
}

// ListWithFilters gets a list of products with comprehensive filtering options
func (r *productRepository) ListWithFilters(
	offset, limit int,
	search string,
	categoryID *uuid.UUID,
	status string,
	minPrice, maxPrice *float64,
	minStock, maxStock *int,
	createdAfter, createdBefore *time.Time,
	sortBy, sortOrder string,
) ([]model.Product, int, error) {
	var total int

	// Define a scan struct with nullable fields
	type productScan struct {
		ID            uuid.UUID       `db:"id"`
		SKU           string          `db:"sku"`
		Name          string          `db:"name"`
		Description   sql.NullString  `db:"description"`
		Price         float64         `db:"price"`
		SalePrice     sql.NullFloat64 `db:"sale_price"`
		CostPrice     sql.NullFloat64 `db:"cost_price"`
		StockQuantity int             `db:"stock_quantity"`
		Status        string          `db:"status"`
		Attributes    string          `db:"attributes"`
		CreatedAt     time.Time       `db:"created_at"`
		UpdatedAt     time.Time       `db:"updated_at"`
	}

	var scanProducts []productScan

	// Base queries
	countQuery := `SELECT COUNT(*) FROM products p`
	listQuery := `
		SELECT p.id, p.sku, p.name, p.description, p.price, p.sale_price, p.cost_price, 
		       p.stock_quantity, p.status, to_json(p.attributes) as attributes, 
		       p.created_at, p.updated_at
		FROM products p
	`

	// Build WHERE clause
	whereClause := " WHERE 1=1" // Start with a true condition to simplify adding AND clauses
	args := []interface{}{}
	argIndex := 1

	// Add search filter
	if search != "" {
		whereClause += " AND (p.name ILIKE $" + strconv.Itoa(argIndex) + " OR p.sku ILIKE $" + strconv.Itoa(argIndex) + " OR p.description ILIKE $" + strconv.Itoa(argIndex) + ")"
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Add category filter
	if categoryID != nil {
		whereClause += " AND p.id IN (SELECT product_id FROM product_categories WHERE category_id = $" + strconv.Itoa(argIndex) + ")"
		args = append(args, *categoryID)
		argIndex++
	}

	// Add status filter
	if status != "" {
		whereClause += " AND p.status = $" + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	// Add price range filters
	if minPrice != nil {
		whereClause += " AND p.price >= $" + strconv.Itoa(argIndex)
		args = append(args, *minPrice)
		argIndex++
	}
	if maxPrice != nil {
		whereClause += " AND p.price <= $" + strconv.Itoa(argIndex)
		args = append(args, *maxPrice)
		argIndex++
	}

	// Add stock quantity filters
	if minStock != nil {
		whereClause += " AND p.stock_quantity >= $" + strconv.Itoa(argIndex)
		args = append(args, *minStock)
		argIndex++
	}
	if maxStock != nil {
		whereClause += " AND p.stock_quantity <= $" + strconv.Itoa(argIndex)
		args = append(args, *maxStock)
		argIndex++
	}

	// Add date range filters
	if createdAfter != nil {
		whereClause += " AND p.created_at >= $" + strconv.Itoa(argIndex)
		args = append(args, *createdAfter)
		argIndex++
	}
	if createdBefore != nil {
		whereClause += " AND p.created_at <= $" + strconv.Itoa(argIndex)
		args = append(args, *createdBefore)
		argIndex++
	}

	// Add sort order
	orderByClause := " ORDER BY p." + sortBy + " " + sortOrder

	// Complete queries
	countQuery += whereClause
	listQuery += whereClause + orderByClause + " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	// Get total count
	err := r.db.Get(&total, countQuery, args[:argIndex-1]...)
	if err != nil {
		return nil, 0, err
	}

	// Get products
	err = r.db.Select(&scanProducts, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Convert scanProducts to model.Product
	products := make([]model.Product, len(scanProducts))
	for i, scanProduct := range scanProducts {
		products[i] = model.Product{
			ID:            scanProduct.ID,
			SKU:           scanProduct.SKU,
			Name:          scanProduct.Name,
			Description:   scanProduct.Description.String,
			Price:         scanProduct.Price,
			StockQuantity: scanProduct.StockQuantity,
			Status:        scanProduct.Status,
			CreatedAt:     scanProduct.CreatedAt,
			UpdatedAt:     scanProduct.UpdatedAt,
		}

		// Handle nullable fields
		if scanProduct.SalePrice.Valid {
			products[i].SalePrice = scanProduct.SalePrice.Float64
		}

		if scanProduct.CostPrice.Valid {
			products[i].CostPrice = scanProduct.CostPrice.Float64
		}

		// Parse attributes JSON
		if scanProduct.Attributes != "" {
			if err := json.Unmarshal([]byte(scanProduct.Attributes), &products[i].Attributes); err != nil {
				return nil, 0, err
			}
		}
	}

	// Get categories for each product
	for i := range products {
		categories, err := r.GetCategories(products[i].ID)
		if err != nil {
			return nil, 0, err
		}
		products[i].Categories = categories
	}

	return products, total, nil
}

func (r *productRepository) GetCategories(productID uuid.UUID) ([]model.Category, error) {
	var categories []model.Category

	query := `
		SELECT c.*
		FROM categories c
		JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = $1
	`

	err := r.db.Select(&categories, query, productID)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}
