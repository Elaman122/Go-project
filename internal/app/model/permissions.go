package model

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

// Permissions holds the permission codes for a single user.
type Permissions []string

// Include checks whether the Permissions slice contains a specific permission code.
func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}

	return false
}

type PermissionModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// GetAllForUser returns all permission codes for a specific user in a Permissions slice.
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
		SELECT permissions.code
		FROM permissions
			INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
			INNER JOIN users ON users_permissions.user_id = users.id
		WHERE users.id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		m.ErrorLog.Printf("Error querying permissions for user: %v", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			m.ErrorLog.Printf("Error scanning row: %v", err)
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		m.ErrorLog.Printf("Error iterating through rows: %v", err)
		return nil, err
	}

	return permissions, nil
}

// AddForUser adds the provided codes for a specific user.
func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	query := `
		INSERT INTO users_permissions
		SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	if err != nil {
		m.ErrorLog.Printf("Error adding permissions for user: %v", err)
		return err
	}

	return nil
}

func (m PermissionModel) CheckPermission(userID int64, requiredPermissionID int) (bool, error) {
    var permissionID int
    query := `
		SELECT permission_id
		FROM users_permissions
		WHERE user_id = $1
    `

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, query, userID).Scan(&permissionID)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil // No permission found for the user
        }
        m.ErrorLog.Printf("Error checking permission: %v", err)
        return false, err
    }

    // Check if the user has the required permission
    if permissionID == requiredPermissionID {
        return true, nil // User has the required permission
    }

    return false, nil // User does not have the required permission
}