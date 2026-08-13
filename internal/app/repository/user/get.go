package user

import (
	"context"

	"github.com/TNJKL/bookmark-pkg/pkg/dbutils"
	"github.com/TNJKL/bookmark-user-service/internal/app/model"
)

// GetUserByUsername retrieves a user record from the database by their username
func (r *sqlRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Where("username = ?", username).First(user).Error

	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return user, nil
}

// GetUserByID retrieves a user record from the database by their ID
func (r *sqlRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error

	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return user, nil
}
