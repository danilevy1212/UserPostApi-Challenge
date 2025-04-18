package inmemory

import (
	"context"
	"database/sql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/models"
)

type PingFunc func(context.Context) error

type UserCreateFunc func(context.Context, models.User) (*models.User, error)
type UserGetAllFunc func(context.Context) ([]*models.User, error)
type UserGetByIDFunc func(context.Context, uint64) (*models.User, error)
type UserDeleteByIDFunc func(context.Context, uint64) error
type UserUpdateFunc func(context.Context, models.UserUpdate) (*models.User, error)

var InMemoryDBPingFn PingFunc = func(c context.Context) error {
	return nil
}
var InMemoryUserCreateFn UserCreateFunc = func(ctx context.Context, user models.User) (*models.User, error) {
	return &models.User{
		ID:    1,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
var InMemoryUserGetAllFn UserGetAllFunc = func(ctx context.Context) ([]*models.User, error) {
	return []*models.User{
		{
			ID:    1,
			Name:  "John Doe",
			Email: "johnnydoe@gmail.com",
		},
		{
			ID:    2,
			Name:  "Daniel Levy Moreno",
			Email: "danielmorenolevy@gmail.com",
		},
	}, nil
}
var InMemoryUserGetByIDFn UserGetByIDFunc = func(ctx context.Context, id uint64) (*models.User, error) {
	return &models.User{
		ID:    id,
		Name:  "Daniel Levy Moreno",
		Email: "danielmorenolevy@gmail.com",
	}, nil
}
var InMemoryUserDeleteByIDFn UserDeleteByIDFunc = func(ctx context.Context, i uint64) error {
	return nil
}
var InMemoryUserUpdateFn UserUpdateFunc = func(ctx context.Context, u models.UserUpdate) (*models.User, error) {
	return &models.User{
		ID:    1,
		Name:  "Daniel Levy Moreno",
		Email: "danielmorenolevy@gmail.com",
	}, nil
}

type PostCreateFunc func(ctx context.Context, post models.Post) (*models.Post, error)
type PostGetAllFunc func(ctx context.Context) ([]*models.Post, error)
type PostGetByIDFunc func(ctx context.Context, id uint64) (*models.Post, error)
type PostDeleteByIDFunc func(ctx context.Context, id uint64) error
type PostUpdateFunc func(ctx context.Context, post models.PostUpdate) (*models.Post, error)

var InMemoryPostCreateFn PostCreateFunc = func(ctx context.Context, post models.Post) (*models.Post, error) {
	return &models.Post{
		ID:      1,
		Title:   "coolio",
		Content: "coolest content",
		UserID:  1,
	}, nil
}

var InMemoryPostGetAllFn PostGetAllFunc = func(ctx context.Context) ([]*models.Post, error) {
	return []*models.Post{
		{
			ID:      1,
			Title:   "coolio",
			Content: "coolest content",
			UserID:  1,
		},
		{
			ID:      2,
			Title:   "another coolio",
			Content: "another coolest content",
			UserID:  1,
		},
		{
			ID:      3,
			Title:   "more coolio",
			Content: "coolest content?",
			UserID:  2,
		},
	}, nil
}

var InMemoryPostGetByIDFn PostGetByIDFunc = func(ctx context.Context, id uint64) (*models.Post, error) {
	return &models.Post{
		ID:      id,
		Title:   "coolio",
		Content: "coolest content",
		UserID:  1,
	}, nil
}

var InMemoryPostDeleteByIDFn PostDeleteByIDFunc = func(ctx context.Context, id uint64) error {
	return nil
}

var InMemoryPostUpdateFn PostUpdateFunc = func(ctx context.Context, post models.PostUpdate) (*models.Post, error) {
	return &models.Post{
		ID:      *post.ID,
		Title:   "coolio",
		Content: "coolest content",
		UserID:  1,
	}, nil
}

type InMemoryDB struct{}

func (im *InMemoryDB) Connection() *sql.DB {
	return nil
}

func (im *InMemoryDB) Ping(ctx context.Context) error {
	return InMemoryDBPingFn(ctx)
}

func (im *InMemoryDB) UserCreate(ctx context.Context, user models.User) (*models.User, error) {
	return InMemoryUserCreateFn(ctx, user)
}

func (im *InMemoryDB) UserGetAll(ctx context.Context) ([]*models.User, error) {
	return InMemoryUserGetAllFn(ctx)
}

func (im *InMemoryDB) UserGetByID(ctx context.Context, id uint64) (*models.User, error) {
	return InMemoryUserGetByIDFn(ctx, id)
}

func (im *InMemoryDB) UserDeleteByID(ctx context.Context, id uint64) error {
	return InMemoryUserDeleteByIDFn(ctx, id)
}

func (im *InMemoryDB) UserUpdate(ctx context.Context, user models.UserUpdate) (*models.User, error) {
	return InMemoryUserUpdateFn(ctx, user)
}


func (im *InMemoryDB) PostCreate(ctx context.Context, post models.Post) (*models.Post, error) {
	return InMemoryPostCreateFn(ctx, post)
}

func (im *InMemoryDB) PostGetAll(ctx context.Context) ([]*models.Post, error) {
	return InMemoryPostGetAllFn(ctx)
}

func (im *InMemoryDB) PostGetByID(ctx context.Context, id uint64) (*models.Post, error) {
	return InMemoryPostGetByIDFn(ctx, id)
}

func (im *InMemoryDB) PostDeleteByID(ctx context.Context, id uint64) error {
	return InMemoryPostDeleteByIDFn(ctx, id)
}

func (im *InMemoryDB) PostUpdate(ctx context.Context, post models.PostUpdate) (*models.Post, error) {
	return InMemoryPostUpdateFn(ctx, post)
}
