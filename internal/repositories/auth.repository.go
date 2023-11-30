package repositories

import (
	"akbarsyarif/coffeeshopgolang/internal/models"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	*sqlx.DB
	// Tx sqlx.Tx
}

func InitializeAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) RepositoryRegister(body *models.UserModel, hashedpassword string, client *sqlx.Tx)  (*sqlx.Rows, error) {
	query := `insert into users (full_name, email, pwd) values ($1, $2, $3) returning id`
	values := []any{body.Full_name, body.Email, hashedpassword}
	result, err := client.Queryx(query, values...)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (r *AuthRepository) RepositoryCreateUserProfile(userId string, client *sqlx.Tx) error {
	query := `
			insert into users_profile (users_id)
			values ($1)
			`
	_, err:= client.Exec(query, userId);
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) RepositoryGetPassword(body *models.GetUserInfoModel) ([]models.GetUserInfoModel, error) {
	result := []models.GetUserInfoModel{}
	query := `
			select 
				u.id, u.full_name, u.pwd, ur.role_name, u.isverified
			from users u
			join user_role ur on u.user_role_id = ur.id
			where email = $1
			`
	values := []any{body.Email}
	if err := r.Select(&result, query, values...); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *AuthRepository) RepositoryLogout(token string) error {
	query := "insert into blacklist (blacklist_token) values ($1)"
	
	if _, err := r.Exec(query, token); err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) RepositoryCheckToken(token string) ([]string, error) {
	result := []string{}
	query := "select blacklist_token from blacklist where blacklist_token like $1"

	if err := r.Select(&result, query, token); err != nil {
		return nil, err
	}
	return result, nil
}