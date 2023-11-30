package models

import "time"

type PromoModel struct {
	Id            string `db:"id" json:"id,omitempty" valid:"-"`
	Promo_name  string `db:"promo_name" form:"promo_name" json:"promo_name,omitempty" valid:"optional"`
	Description   string `db:"description" form:"description" json:"description,omitempty" valid:"optional"`
	Discount_type string `db:"discount_type" form:"discount_type" json:"discount_type,omitempty" valid:"optional"`
	Flat_amount  int `db:"flat_amount" form:"flat_amount" json:"flat_amount,omitempty" valid:"optional"`
	Percent_amount float64 `db:"percent_amount" form:"percent_amount" json:"percent_amount,omitempty" valid:"optional"`
	Created_at    *time.Time `db:"created_at" json:"created_at,omitempty" valid:"-"`
}