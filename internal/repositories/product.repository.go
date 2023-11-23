package repositories

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	*sqlx.DB
}

func InitializeRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) RepositoryGetAllProduct() ([]models.ProductModel, error) {
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
	err := r.Select(&result, query);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProductRepository) RepositoryCountProduct() (*sql.Rows, error) {
	query := `
		SELECT
			COUNT(*) AS "Total_product"
		FROM
			products p `
	result, err := r.Query(query)
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
	// rows, err := r.Query(query, productId);
	// if err != nil {
		// 	return nil, err
		// }
		// defer rows.Close()
		
		// var result []models.ProductModel
		
		// for rows.Next() {
			// 	var res models.ProductModel
			// 	err := rows.Scan(&res.Id, &res.Product_image, &res.Product_name)
			// 	if err != nil {
				// 		return result, err
				// 	}
				// 	result = append(result, res)
	// }
	// if err = rows.Err(); err != nil {
		// 	return result, err
		// }
		// return result, nil
	}

func (r *ProductRepository) RepositoryCreateProduct(body *models.ProductModel) (sql.Result, error) {
	query := `
			insert into products (product_name, description, price, category_id, promo_id)
			values (:product_name, :description, :price, (select id from category where category_name = :category), (select id from promo where promo_name = :promo))
			returning product_name, description, price
			`
		
	result, err := r.NamedExec(query, body);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProductRepository) RepositoryUpdateProduct(body *models.ProductModel, productId string) (sql.Result, error) {
	query := `update products set `

	var filteredBody []string
	filterBody := make(map[string]interface{})
	filterBody["id"] = productId
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
	// log.Println(len(query))
	if len(query) < 56 {
		err := errors.New("Please Input at Least One Change")
		return nil, err
	}

	result, err := r.NamedExec(query, filterBody);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProductRepository) RepositoryDeleteProduct(productId string) (sql.Result, error) {
	query := `
			delete from products
			where id = $1
			returning id
			`
		
	result, err := r.Exec(query, productId);
	if err != nil {
		return nil, err
	}
	return result, nil
}