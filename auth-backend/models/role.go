package models

import (
	"context"
	"log"
)

// --------------------- Admin Role ---------------------
type AdminRole struct {
	UserID int
	Level  string // "admin", "manager"
}

func GetAdminRole(userID int) (*AdminRole, error) {
	role := &AdminRole{}
	err := DB.QueryRow(context.Background(),
		`SELECT user_id, level FROM admin_roles WHERE user_id=$1`, userID).
		Scan(&role.UserID, &role.Level)
	if err != nil {
		log.Printf("❌ GetAdminRole error: %v", err)
		return nil, err
	}
	return role, nil
}

func CreateOrUpdateAdminRole(role *AdminRole) error {
	_, err := DB.Exec(context.Background(),
		`INSERT INTO admin_roles (user_id, level) VALUES ($1, $2)
		 ON CONFLICT (user_id) DO UPDATE SET level=$2`,
		role.UserID, role.Level)
	if err != nil {
		log.Printf("❌ CreateOrUpdateAdminRole error: %v", err)
	}
	return err
}

// --------------------- Staff Role ---------------------
type StaffRole struct {
	UserID int
	Dept   string
}

func GetStaffRole(userID int) (*StaffRole, error) {
	role := &StaffRole{}
	err := DB.QueryRow(context.Background(),
		`SELECT user_id, dept FROM staff_roles WHERE user_id=$1`, userID).
		Scan(&role.UserID, &role.Dept)
	if err != nil {
		log.Printf("❌ GetStaffRole error: %v", err)
		return nil, err
	}
	return role, nil
}

func CreateOrUpdateStaffRole(role *StaffRole) error {
	_, err := DB.Exec(context.Background(),
		`INSERT INTO staff_roles (user_id, dept) VALUES ($1, $2)
		 ON CONFLICT (user_id) DO UPDATE SET dept=$2`,
		role.UserID, role.Dept)
	if err != nil {
		log.Printf("❌ CreateOrUpdateStaffRole error: %v", err)
	}
	return err
}

// --------------------- Customer Role ---------------------
type CustomerRole struct {
	UserID        int
	LoyaltyPoints int
}

func GetCustomerRole(userID int) (*CustomerRole, error) {
	role := &CustomerRole{}
	err := DB.QueryRow(context.Background(),
		`SELECT user_id, loyalty_points FROM customer_roles WHERE user_id=$1`, userID).
		Scan(&role.UserID, &role.LoyaltyPoints)
	if err != nil {
		log.Printf("❌ GetCustomerRole error: %v", err)
		return nil, err
	}
	return role, nil
}

func CreateOrUpdateCustomerRole(role *CustomerRole) error {
	_, err := DB.Exec(context.Background(),
		`INSERT INTO customer_roles (user_id, loyalty_points) VALUES ($1, $2)
		 ON CONFLICT (user_id) DO UPDATE SET loyalty_points=$2`,
		role.UserID, role.LoyaltyPoints)
	if err != nil {
		log.Printf("❌ CreateOrUpdateCustomerRole error: %v", err)
	}
	return err
}
