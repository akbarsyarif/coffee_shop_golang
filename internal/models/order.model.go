package models

import "time"

type OrderModel struct {
	Id         string      `db:"id" json:"id" valid:"-"`
	User_id         string `db:"user_id" form:"user_id" json:"user_id,omitempty" valid:"-"`
	Full_name  string      `db:"full_name" form:"full_name" json:"full_name" valid:"-"`
	Status     string      `db:"status_name" form:"status_name" json:"status_name" valid:"-"`
	Shipping   string      `db:"shipping_name" form:"shipping_name" json:"shipping_name" valid:"-"`
	Total      int         `db:"total" form:"total" json:"total" valid:"-"`
	Created_at *time.Time  `db:"created_at" json:"created_at" valid:"-"`
	Product []OrderProductModel `form:"product" json:"product,omitempty" valid:"-"`
	// Deleted_at    *time.Time `db:"deleted_at" json:"deleted_at"`
}

type OrderProductModel struct {
	// Id         string      `db:"id" json:"id"`
	Order_id  string      `db:"order_id" form:"order_id" json:"order_id,omitempty" valid:"-"`
	Full_name  string      `db:"full_name" form:"full_name" json:"full_name,omitempty" valid:"-"`
	Product_image interface{} `db:"product_image" form:"product_image" json:"product_image,omitempty" valid:"-"`
	Product_name  string `db:"product_name" form:"product_name" json:"product_name,omitempty" valid:"-"`
	Quantity int `db:"quantity" form:"quantity" json:"quantity,omitempty" valid:"-"`
	Size     string      `db:"size_name" form:"size_name" json:"size_name,omitempty" valid:"-"`
	WithIce     bool     `db:"with_ice" form:"with_ice" json:"with_ice" valid:"-"`
	Status     string      `db:"status_name" form:"status_name" json:"status_name,omitempty" valid:"-"`
	Sub_total      int      `db:"sub_total" form:"sub_total" json:"sub_total,omitempty" valid:"-"`
	Shipping   string      `db:"shipping_name" form:"shipping_name" json:"shipping_name,omitempty" valid:"-"`
	Created_at *time.Time  `db:"created_at" json:"created_at,omitempty" valid:"-"`
	// Deleted_at    *time.Time `db:"deleted_at" json:"deleted_at"`
}

