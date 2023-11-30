package repositories

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PromoRepository struct {
	*sqlx.DB
}

func InitializePromoRepository(db *sqlx.DB) *PromoRepository {
	return &PromoRepository{db}
}

func (r *PromoRepository) RepositoryCountPromo() ([]int, error) {
	var total_data = []int{}
	query := `
		select
			count(*) as "total_promo"
		from
			promo
		`
	err := r.Select(&total_data, query)
	if err != nil {
		return nil, err
	}
	return total_data, nil
}

func (r *PromoRepository) RepositoryGetPromo(page string) ([]models.PromoModel, error) {
	result := []models.PromoModel{}
	offset, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	offset = (offset - 1) * 5
	query := `select id, promo_name, description, discount_type, flat_amount, percent_amount, created_at from promo p order by id asc limit 5 offset $1`

	err = r.Select(&result, query, offset);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PromoRepository) RepositoryCreatePromo(body *models.PromoModel) (sql.Result, error) {
	query := `
	insert into promo (promo_name, description, discount_type, flat_amount, percent_amount)
			values (:promo_name, :description, :discount_type, :flat_amount, :percent_amount)
			returning promo_name
			`
		
	result, err := r.NamedExec(query, body);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PromoRepository) RepositoryUpdatePromo(body *models.PromoModel, promoId string) (sql.Result, error) {
	query := `update promo set `

	var filteredBody []string
	filterBody := make(map[string]interface{})
	filterBody["id"] = promoId
	if body.Promo_name != "" {
		filteredBody = append(filteredBody, "promo_name = :promo_name")
		filterBody["promo_name"] = body.Promo_name
	}
	if body.Description != "" {
		filteredBody = append(filteredBody, "description = :description")
		filterBody["description"] = body.Description
	}
	if body.Discount_type != "" {
		filteredBody = append(filteredBody, "discount_type = :discount_type")
		filterBody["discount_type"] = body.Discount_type
	}
	if body.Flat_amount != 0 {
		filteredBody = append(filteredBody, "flat_amount = :flat_amount")
		filterBody["flat_amount"] = body.Flat_amount
	}
	if body.Percent_amount != 0 {
		filteredBody = append(filteredBody, "percent_amount = :percent_amount")
		filterBody["percent_amount"] = body.Percent_amount
	}
	if len(filteredBody) > 0 {
		query += strings.Join(filteredBody, ", ")
	}
	query += `, updated_at = now() where id = :id`
	log.Println(len(query))

	result, err := r.NamedExec(query, filterBody);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PromoRepository) RepositoryDeletePromo(promoId string) (sql.Result, error) {
	query := `
			delete from promo
			where id = $1
			returning id
			`
		
	result, err := r.Exec(query, promoId);
	if err != nil {
		return nil, err
	}
	return result, nil
}

