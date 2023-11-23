package models

import "time"

type PromoModel struct {
	Id            string `db:"id" json:"id"`
	Promo_name  string `db:"promo_name" form:"promo_name" json:"promo_name"`
	Description   string `db:"description" form:"description" json:"description"`
	Discount_type string `db:"discount_type" form:"discount_type" json:"discount_type"`
	Flat_amount  int `db:"flat_amount" form:"flat_amount" json:"flat_amount"`
	Percent_amount float64 `db:"percent_amount" form:"percent_amount" json:"percent_amount"`
	Created_at    *time.Time `db:"created_at" json:"created_at"`
}