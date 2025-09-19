package models

import (
	"auth-backend/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

type User struct {
	ID           int
	Name         string
	PhoneNumber  *string // nullable
	Email        *string // nullable
	PasswordHash string
	Role         string // "customer", "staff", "admin"
	IsVerified   bool
	GoogleID     *string // nullable
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RoleID       int
}

// --------------------- Fetch Users ---------------------

func GetUserByPhone(phone string) (*User, error) {
	u := &User{}
	query := `
		SELECT u.id, u.name, u.phone, u.email, u.password_hash, u.google_id, 
		       u.is_verified, u.role_id, r.name, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.phone=$1
	`
	err := DB.QueryRow(context.Background(), query, phone).Scan(
		&u.ID, &u.Name, &u.PhoneNumber, &u.Email, &u.PasswordHash, &u.GoogleID,
		&u.IsVerified, &u.RoleID, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("❌ GetUserByPhone error: %v", err)
		return nil, err
	}
	return u, nil
}

func GetUserByEmail(email string) (*User, error) {
	u := &User{}
	err := DB.QueryRow(context.Background(),
		`SELECT u.id, u.name, u.phone, u.email, u.password_hash, u.google_id, 
		        u.is_verified, u.role_id, r.name, u.created_at, u.updated_at
		 FROM users u
		 JOIN roles r ON u.role_id = r.id
		 WHERE u.email=$1`, email).
		Scan(&u.ID, &u.Name, &u.PhoneNumber, &u.Email, &u.PasswordHash, &u.GoogleID,
			&u.IsVerified, &u.RoleID, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("❌ GetUserByEmail error: %v", err)
		return nil, err
	}
	return u, nil
}

func GetUserByID(id int) (*User, error) {
	u := &User{}
	err := DB.QueryRow(context.Background(),
		`SELECT u.id, u.name, u.phone, u.email, u.password_hash, u.google_id, 
		        u.is_verified, u.role_id, r.name, u.created_at, u.updated_at
		 FROM users u
		 JOIN roles r ON u.role_id = r.id
		 WHERE u.id=$1`, id).
		Scan(&u.ID, &u.Name, &u.PhoneNumber, &u.Email, &u.PasswordHash, &u.GoogleID,
			&u.IsVerified, &u.RoleID, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("❌ GetUserByID error: %v", err)
		return nil, err
	}
	return u, nil
}

// --------------------- Create / Fetch Users ---------------------
func CreateOrFetchUser(user *User, extra map[string]interface{}) (*User, bool, error) {
	log.Printf("👉 CreateOrFetchUser called with Name=%s, Phone=%v, Email=%v, GoogleID=%v, Role=%s",
		user.Name, user.PhoneNumber, user.Email, user.GoogleID, user.Role)

	// Ensure admin/staff is always verified
	if user.Role == "admin" || user.Role == "staff" {
		user.IsVerified = true
	}

	// Ensure password is hashed
	if user.PasswordHash != "" {
		if !strings.HasPrefix(user.PasswordHash, "$2a$") { // bcrypt hash prefix
			hash, err := utils.HashPassword(user.PasswordHash)
			if err != nil {
				return nil, false, fmt.Errorf("failed to hash password: %w", err)
			}
			user.PasswordHash = hash
		}
	}

	// Resolve role_id
	switch user.Role {
	case "admin":
		user.RoleID = 1
	case "staff":
		user.RoleID = 2
	default:
		user.RoleID = 3
	}

	var existing *User
	var err error

	// Lookup by phone
	if user.PhoneNumber != nil && *user.PhoneNumber != "" {
		existing, err = GetUserByPhone(*user.PhoneNumber)
		if err != nil {
			return nil, false, err
		}
		if existing != nil {
			mergeUserFields(existing, user)
			if err := UpdateUser(existing, extra); err != nil {
				log.Printf("❌ UpdateUser error: %v", err)
			}
			return existing, false, nil
		}
	}

	// Lookup by email
	if user.Email != nil && *user.Email != "" {
		existing, err = GetUserByEmail(*user.Email)
		if err != nil {
			return nil, false, err
		}
		if existing != nil {
			mergeUserFields(existing, user)
			if err := UpdateUser(existing, extra); err != nil {
				log.Printf("❌ UpdateUser error: %v", err)
			}
			return existing, false, nil
		}
	}

	// Lookup by GoogleID
	if user.GoogleID != nil && *user.GoogleID != "" {
		var id int
		row := DB.QueryRow(context.Background(), `SELECT id FROM users WHERE google_id=$1`, *user.GoogleID)
		if err := row.Scan(&id); err == nil {
			existing, _ := GetUserByID(id)
			if existing != nil {
				mergeUserFields(existing, user)
				if err := UpdateUser(existing, extra); err != nil {
					log.Printf("❌ UpdateUser error: %v", err)
				}
				return existing, false, nil
			}
		}
	}

	// Insert new user
	err = DB.QueryRow(context.Background(),
		`INSERT INTO users (name, phone, email, password_hash, google_id, role_id, is_verified, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
		 RETURNING id`,
		user.Name, user.PhoneNumber, user.Email, user.PasswordHash, user.GoogleID, user.RoleID, user.IsVerified,
	).Scan(&user.ID)

	if err != nil {
		log.Printf("❌ Insert user failed: %v", err)
		return nil, false, fmt.Errorf("insert user failed: %w", err)
	}

	// Insert role extension data
	if err := UpdateRoleDetails(user.ID, user.Role, extra); err != nil {
		log.Printf("❌ UpdateRoleDetails error: %v", err)
	}

	log.Printf("✅ User created successfully: ID=%d", user.ID)
	created, _ := GetUserByID(user.ID)
	return created, true, nil
}

// --------------------- Merge Helper ---------------------
func mergeUserFields(existing, incoming *User) {
	// Merge name
	if incoming.Name != "" {
		existing.Name = incoming.Name
	}

	// Merge email if existing is nil/empty
	if incoming.Email != nil && (existing.Email == nil || *existing.Email == "") {
		existing.Email = incoming.Email
	}

	// Merge phone if existing is nil/empty
	if incoming.PhoneNumber != nil && (existing.PhoneNumber == nil || *existing.PhoneNumber == "") {
		existing.PhoneNumber = incoming.PhoneNumber
	}

	// Merge GoogleID if existing is nil/empty
	if incoming.GoogleID != nil && (existing.GoogleID == nil || *existing.GoogleID == "") {
		existing.GoogleID = incoming.GoogleID
	}

	// Merge password hash safely
	if incoming.PasswordHash != "" {
		if existing.PasswordHash == "" || strings.HasPrefix(incoming.PasswordHash, "$2a$") {
			existing.PasswordHash = incoming.PasswordHash
		}
	}

	// Merge verification status: only upgrade to true
	if incoming.IsVerified {
		existing.IsVerified = true
	}
}

// --------------------- Update User ---------------------

func UpdateUser(user *User, extra map[string]interface{}) error {
	// Update main users table
	_, err := DB.Exec(context.Background(),
		`UPDATE users SET name=$1, phone=$2, email=$3, password_hash=$4, google_id=$5, role_id=$6, is_verified=$7, updated_at=NOW()
		 WHERE id=$8`,
		user.Name, user.PhoneNumber, user.Email, user.PasswordHash, user.GoogleID, user.RoleID, user.IsVerified, user.ID)
	if err != nil {
		return err
	}

	// Update role extension table
	if err := UpdateRoleDetails(user.ID, user.Role, extra); err != nil {
		log.Printf("❌ UpdateRoleDetails error: %v", err)
		return err
	}

	return nil
}

// --------------------- Role Extension Table ---------------------

func UpdateRoleDetails(userID int, role string, extra map[string]interface{}) error {
	switch role {
	case "admin":
		level := "admin"
		if extra != nil {
			if l, ok := extra["level"].(string); ok && l != "" {
				level = l
			}
		}
		_, err := DB.Exec(context.Background(),
			`INSERT INTO admin_roles (user_id, level)
			 VALUES ($1, $2)
			 ON CONFLICT (user_id) DO UPDATE SET level=$2`, userID, level)
		return err

	case "staff":
		dept := "general"
		if extra != nil {
			if d, ok := extra["dept"].(string); ok && d != "" {
				dept = d
			}
		}
		_, err := DB.Exec(context.Background(),
			`INSERT INTO staff_roles (user_id, dept)
			 VALUES ($1, $2)
			 ON CONFLICT (user_id) DO UPDATE SET dept=$2`, userID, dept)
		return err

	case "customer":
		points := 0
		if extra != nil {
			if p, ok := extra["loyalty_points"].(int); ok {
				points = p
			}
		}
		_, err := DB.Exec(context.Background(),
			`INSERT INTO customer_roles (user_id, loyalty_points)
			 VALUES ($1, $2)
			 ON CONFLICT (user_id) DO UPDATE SET loyalty_points=$2`, userID, points)
		return err
	}
	return fmt.Errorf("invalid role: %s", role)
}

// --------------------- Other Helpers ---------------------

func UpdateIsVerified(userID int, verified bool) error {
	_, err := DB.Exec(context.Background(),
		"UPDATE users SET is_verified=$1, updated_at=NOW() WHERE id=$2",
		verified, userID)
	return err
}

func GetAllUsers() ([]*User, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT u.id, u.name, u.phone, u.email, u.password_hash, u.google_id, 
		        u.is_verified, u.role_id, r.name, u.created_at, u.updated_at
		 FROM users u
		 JOIN roles r ON u.role_id = r.id
		 ORDER BY u.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.PhoneNumber, &u.Email, &u.PasswordHash, &u.GoogleID,
			&u.IsVerified, &u.RoleID, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			log.Printf("❌ Scan user error: %v", err)
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetAdminLevel(userID int) (string, error) {
	var level string
	err := DB.QueryRow(context.Background(),
		`SELECT level FROM admin_roles WHERE user_id=$1`, userID).Scan(&level)
	if err != nil {
		return "", err
	}
	return level, nil
}

func GetStaffDept(userID int) (string, error) {
	var dept string
	err := DB.QueryRow(context.Background(),
		`SELECT dept FROM staff_roles WHERE user_id=$1`, userID).Scan(&dept)
	if err != nil {
		return "", err
	}
	return dept, nil
}

func GetCustomerPoints(userID int) (int, error) {
	var points int
	err := DB.QueryRow(context.Background(),
		`SELECT loyalty_points FROM customer_roles WHERE user_id=$1`, userID).Scan(&points)
	if err != nil {
		return 0, err
	}
	return points, nil
}
