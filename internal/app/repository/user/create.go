package user

import (
	"context"

	"github.com/TNJKL/bookmark-pkg/pkg/dbutils"
	"github.com/TNJKL/bookmark-user-service/internal/app/model"
)

// CreateUser inserts a new user record into the database and returns the created user
func (r *sqlRepository) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return newUser, nil
}
