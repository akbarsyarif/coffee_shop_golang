package repositories

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	*sqlx.DB
	// Tx sqlx.Tx
}

func InitializeOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) RepositoryCountOrder(userId, status string) ([]int, error) {
	var total_data = []int{}
	query := `
		select
			count(*) as "total_order"
		from
			"order" o
		join
			status s on o.status_id = s.id
		where o.users_id = $1 and s."name" like $2
		`
	err := r.Select(&total_data, query, userId, status)
	if err != nil {
		return nil, err
	}
	return total_data, nil
}

func (r *OrderRepository) RepositoryGetAllOrder() ([]models.OrderModel, error) {
	result := []models.OrderModel{}
	query := `
			select
				o.id, u.full_name, sh."name" as "shipping_name", s."name" as "status_name", total, o.created_at
			from
				"order" o
			join
				users u on o.users_id  = u.id
			join
				status s on o.status_id = s.id
			join
				shipping sh on o.shipping_id = sh.id
			order by id asc
			`
	
	err := r.Select(&result, query);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepository) RepositoryGetOrderPerUser(userId, status, page string) ([]models.OrderModel, error) {
	result := []models.OrderModel{}
	offset, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	offset = (offset - 1) * 5

	query := `
			select
				o.id, u.full_name, sh."name" as "shipping_name", s."name" as "status_name", total, o.created_at
			from
				"order" o
			join
				users u on o.users_id  = u.id
			join
				status s on o.status_id = s.id
			join
				shipping sh on o.shipping_id = sh.id
			where o.users_id = $1 and s."name" like $2 and deleted_at is null
			`
	query += ` limit 5 offset $3`
	// log.Println(query)
	err = r.Select(&result, query, userId, status, offset);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepository) RepositoryGetOrderDetail(orderId, userId string) ([]models.OrderProductModel, error) {
	result := []models.OrderProductModel{}

	query := `
			select
				p.product_image, p.product_name, o.id as "order_id",
				u.full_name, op.quantity, s."size" as "size_name",
				op.ice as "with_ice", sh."name" as "shipping_name",
				st."name" as "status_name", op.sub_total, o.created_at
			from
				order_products op
			join
				"order" o  on op.order_id  = o.id
			join
				products p on op.products_id = p.id
			join
				users u on o.users_id  = u.id
			join
				sizes s on op.sizes_id = s.id
			join
				status st on o.status_id = st.id
			join
				shipping sh on o.shipping_id = sh.id
			where o.id = $1 and o.users_id = $2 and op.deleted_at is null
			`
	// log.Println(query)
	err := r.Select(&result, query, orderId, userId);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepository) RepositoryCreateOrder(body *models.OrderModel, client *sqlx.Tx) (*sqlx.Rows, error) {
	query := `
			insert into "order" (users_id, total, shipping_id)
			values (:user_id, :total, (select id from shipping where "name" = :shipping_name))
			returning id
			`
	result, err := client.NamedQuery(query, body);
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (r *OrderRepository) RepositoryCreateOrderProduct(body *models.OrderModel, client *sqlx.Tx, orderId string) (sql.Result, error) {
	query := `
			insert into order_products (order_id, products_id, quantity, sizes_id, ice, sub_total)
			values 
			`

	var filteredBody []string
	filterBody := make(map[string]interface{})
	filterBody["order_id"] = orderId

	j := 2
	for i := 0; i < len(body.Product); i++ {
		filteredBody = append(filteredBody, "(:order_id")
		filteredBody = append(filteredBody, fmt.Sprintf(`(select id from products where product_name = :product_name%d)`, j))
		filterBody[fmt.Sprintf("product_name%d", j)] = body.Product[i].Product_name
		filteredBody = append(filteredBody, fmt.Sprintf(`:quantity%d`, j))
		filterBody[fmt.Sprintf("quantity%d", j)] = body.Product[i].Quantity
		filteredBody = append(filteredBody, fmt.Sprintf(`(select id from sizes where size = :size_name%d)`, j))
		filterBody[fmt.Sprintf("size_name%d", j)] = body.Product[i].Size
		filteredBody = append(filteredBody, fmt.Sprintf(`:with_ice%d`, j))
		filterBody[fmt.Sprintf("with_ice%d", j)] = body.Product[i].WithIce
		filteredBody = append(filteredBody, fmt.Sprintf(`:sub_total%d)`, j))
		filterBody[fmt.Sprintf("sub_total%d", j)] = body.Product[i].Sub_total
		j++
	}
	if len(filteredBody) > 0 {
		query += strings.Join(filteredBody, ", ")
	}

	result, err := client.NamedExec(query, filterBody);
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (r *OrderRepository) RepositoryUpdateOrder(body *models.OrderModel, OrderId string) (sql.Result, error) {
	query := `update "order"
				set status_id = (select id from status where "name" = $2), updated_at = now()
				where id = $1`

	result, err := r.Exec(query, OrderId, body.Status);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepository) RepositoryDeleteOrder(orderId string, client *sqlx.Tx) (sql.Result, error) {
	query := `update "order" set deleted_at = now() where id = $1`
		
	result, err := client.Exec(query, orderId);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepository) RepositoryDeleteOrderProduct(orderId string, client *sqlx.Tx) (sql.Result, error) {
	query := `update order_products set deleted_at = now() where order_id = $1`
		
	result, err := client.Exec(query, orderId);
	if err != nil {
		return nil, err
	}
	return result, nil
}