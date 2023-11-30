package repositories

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"database/sql"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	*sqlx.DB
}

func InitializeUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) RepositoryCountUser() ([]int, error) {
	var total_data = []int{}
	query := `
		select
			count(*) as "total_user"
		from
			users
		`
	err := r.Select(&total_data, query)
	if err != nil {
		return nil, err
	}
	return total_data, nil
}

func (r *UserRepository) RepositoryGetAllUser(page string) ([]models.UserModel, error) {
	result := []models.UserModel{}
	offset, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	offset = (offset - 1) * 5
	query := `
			select 
				u.id, u.full_name, u.email, up.profile_pic, up.phone_number, up.address, ur.role_name, u.isverified, u.created_at
			from
				users u
			join
				users_profile up on u.id = up.users_id
			join
				user_role ur on u.user_role_id = ur.id
				`
	query += ` limit 5 offset $1`
	
	err = r.Select(&result, query, offset);
	// err := r.Select(&result, query, "%Test%");
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserRepository) RepositoryGetUserDetail(userId string) ([]models.UserModel, error) {
	result := []models.UserModel{}
	query := `
			select 
				u.id, u.full_name, u.email, up.profile_pic, up.phone_number, up.address, ur.role_name, u.isverified, u.created_at
			from
				users u
			join
				users_profile up on u.id = up.users_id
			join
				user_role ur on u.user_role_id = ur.id
			where u.id = $1
				`
	
	err := r.Select(&result, query, userId);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserRepository) RepositoryCreateUser(body *models.UserModel, client *sqlx.Tx) (*sqlx.Rows, error) {
	query := `
			insert into users (full_name, email, pwd)
			values (:full_name, :email, :pwd)
			returning id
			`
	result, err := client.NamedQuery(query, body);
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserRepository) RepositoryCreateUserProfile(userId any, client *sqlx.Tx) (sql.Result, error) {
	query := `
	insert into users_profile (users_id)
	values ($1)
	`
	result, err:= client.Exec(query, userId);
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (r *UserRepository) RepositoryUpdateUser(body *models.UserModel, userId, imageUrl string) (sql.Result, error) {
	query := `update users_profile set `
	
	var filteredBody []string
	filterBody := make(map[string]interface{})
	filterBody["id"] = userId
	if imageUrl != "" {
		filteredBody = append(filteredBody, "profile_pic = :profile_pic")
		filterBody["profile_pic"] = imageUrl
	}
	if body.Phone_number != "" {
		filteredBody = append(filteredBody, "phone_number = :phone_number")
		filterBody["phone_number"] = body.Phone_number
	}
	if body.Address != "" {
		filteredBody = append(filteredBody, "address = :address")
		filterBody["address"] = body.Address
	}
	// if body.Full_name != "" {
		// 	filteredBody = append(filteredBody, "full_name = :full_name")
		// 	filterBody["full_name"] = body.Full_name
		// }
		if len(filteredBody) > 0 {
			query += strings.Join(filteredBody, ", ")
		}
		query += `, updated_at = now() where users_id = :id`
		
		result, err := r.NamedExec(query, filterBody);
		if err != nil {
			return nil, err
		}
		return result, nil
}

func (r *UserRepository) RepositoryDeleteUser(userId string) (sql.Result, error) {
	query := `
	delete from users
			where id = $1
			returning full_name
			`
			
			result, err := r.Exec(query, userId);
			if err != nil {
		return nil, err
	}
	return result, nil
}

// func (r *UserRepository) RepositorySelectId(email string, client *sqlx.Tx) (any, error) {
// 	var result []string
// 	querySelect := `select id from users where email = $1`
// 	err := client.Select(&result, querySelect, email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }