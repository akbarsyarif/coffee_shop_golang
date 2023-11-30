package repositories

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	*sqlx.DB
}

func InitializeRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) RepositoryCountProduct(params *models.ProductParams) ([]int, error) {
	var total_data = []int{}
	query := `
		SELECT
			COUNT(*) AS "total_product"
		from
			products p 
		join
			category c on p.category_id = c.id	
			`

	var filteredParams []string
	filterParams := []any{}
	if params.Product_name != "" {
		filteredParams = append(filteredParams, "p.product_name ilike $"+fmt.Sprint(len(filterParams)+1))
		filterParams = append(filterParams, "%"+params.Product_name+"%")
	}
	if params.Max_price != "" {
		filteredParams = append(filteredParams, "p.price <  $"+fmt.Sprint(len(filterParams)+1))
		filterParams = append(filterParams, params.Max_price)
	}
	if params.Min_price != "" {
		filteredParams = append(filteredParams, "p.price >  $"+fmt.Sprint(len(filterParams)+1))
		filterParams = append(filterParams, params.Min_price)	
	}
	if params.Category != "" {
		filteredParams = append(filteredParams, "c.category_name = $"+fmt.Sprint(len(filterParams)+1))
		filterParams = append(filterParams, params.Category)
	}
	if len(filteredParams) > 0 {
		query += " where " + strings.Join(filteredParams, " and ")
	}

	err := r.Select(&total_data, query, filterParams...)
	if err != nil {
		return nil, err
	}
	return total_data, nil
}

func (r *ProductRepository) RepositoryGetAllProduct(params *models.ProductParams) ([]models.ProductModel, error) {
	result := []models.ProductModel{}
	
	query := `
			select 
				p.id, p.product_image, p.product_name, c.category_name as "category", p.description, p.rating, p.price, pr.promo_name as "promo", p.created_at
			from
				products p
			join
				category c on p.category_id = c.id
			join
				promo pr on promo_id = pr.id
				`
	// test := []string {"hainan", "ayam"}
	var filteredParams []string
	filterParams := []any{}
	if params.Product_name != "" {
		filteredParams = append(filteredParams, "p.product_name ilike $"+fmt.Sprint(len(filterParams)+1))
		filterParams = append(filterParams, "%"+params.Product_name+"%")
	}
	if params.Max_price != "" {
		filteredParams = append(filteredParams, "p.price <  $"+fmt.Sprint(len(filterParams)+1))
		filterParams = append(filterParams, params.Max_price)
	}
	if params.Min_price != "" {
		filteredParams = append(filteredParams, "p.price >  $"+fmt.Sprint(len(filterParams)+1))
		filterParams = append(filterParams, params.Min_price)	
	}
	if params.Category != "" {
		filteredParams = append(filteredParams, "c.category_name = $"+fmt.Sprint(len(filterParams)+1))
		filterParams = append(filterParams, params.Category)
	}
	if len(filteredParams) > 0 {
		query += " where " + strings.Join(filteredParams, " and ")
	}
	
	i := 0
	var filteredSort []string
	filterSort := []string{}
	if params.Name != "" {
		filterSort = append(filterSort, params.Name)
		filteredSort = append(filteredSort, fmt.Sprintf("p.product_name %v", filterSort[i]))
		i++
	}
	if params.Price != "" {
		filterSort = append(filterSort, params.Price)
		filteredSort = append(filteredSort, fmt.Sprintf("p.price %v", filterSort[i]))
		i++
	}
	if params.Created_at != "" {
		filterSort = append(filterSort, params.Created_at)
		filteredSort = append(filteredSort, fmt.Sprintf("p.price %v", filterSort[i]))
	}
	if len(filteredSort) > 0 {
		query += " order by " + strings.Join(filteredSort, ", ")
	}
	page, err := strconv.Atoi(params.Page)
	if err != nil {
		return nil, err
	}
	pagination := fmt.Sprintf(" limit 6 offset %d", (page - 1) * 6)
	query += pagination

	err = r.Select(&result, query, filterParams...);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProductRepository) RepositoryGetProductDetail(productId string) ([]models.ProductModel, error) {
	result := []models.ProductModel{}
	
	query := `
			select 
				p.id, p.product_image, p.product_name, c.category_name as "category", p.description, p.rating, p.price, pr.promo_name as "promo", p.created_at
			from
				products p
			join
				category c on p.category_id = c.id
			join
				promo pr on promo_id = pr.id
			where
				p.id = $1
				`

	err := r.Select(&result, query, productId);
	if err != nil {
		return nil, err
	}
	
	return result, nil
	}

func (r *ProductRepository) RepositoryCreateProduct(body *models.ProductModel) (*sqlx.Rows, error) {
	query := `
			insert into products (product_name, description, price, category_id, promo_id)
			values (:product_name, :description, :price, (select id from category where category_name = :category), (select id from promo where promo_name = :promo))
			returning id
			`
		
	result, err := r.NamedQuery(query, body);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProductRepository) RepositoryInsertimage(imageUrl, productId string) (sql.Result, error) {
	query := `update products set product_image = $1 where id = $2`

	result, err := r.Exec(query, imageUrl, productId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProductRepository) RepositoryUpdateProduct(body *models.ProductModel, productId, imageUrl string) (sql.Result, error) {
	query := `update products set `

	var filteredBody []string
	filterBody := make(map[string]interface{})
	filterBody["id"] = productId
	if imageUrl != "" {
		filteredBody = append(filteredBody, "product_image = :product_image")
		filterBody["product_image"] = imageUrl
	}
	if body.Product_name != "" {
		filteredBody = append(filteredBody, "product_name = :product_name")
		filterBody["product_name"] = body.Product_name
	}
	if body.Description != "" {
		filteredBody = append(filteredBody, "description = :description")
		filterBody["description"] = body.Description
	}
	if body.Price != 0 {
		filteredBody = append(filteredBody, "price = :price")
		filterBody["price"] = body.Price
	}
	if body.Category != "" {
		filteredBody = append(filteredBody, "category_id = (select id from category where category_name = :category)")
		filterBody["category"] = body.Category
	}
	if body.Promo != "" {
		filteredBody = append(filteredBody, "promo_id = (select id from promo where promo_name = :promo)")
		filterBody["promo"] = body.Promo
	}
	if len(filteredBody) > 0 {
		query += strings.Join(filteredBody, ", ")
	}
	query += `, updated_at = now() where id = :id`

	result, err := r.NamedExec(query, filterBody);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProductRepository) RepositoryDeleteProduct(productId string) (sql.Result, error) {
	query := `
			update products set deleted_at = now()
			where id = $1
			`
		
	result, err := r.Exec(query, productId);
	if err != nil {
		return nil, err
	}
	return result, nil
}