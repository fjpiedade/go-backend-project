package main

import (
	"context"
	"social/internal/store"
)

// Mock que implementa UserStorageInterface
type mockUserStore struct {
	createFn  func(context.Context, *store.User) error
	getByIDFn func(context.Context, int64) (*store.User, error)
	getAllFn  func(context.Context) ([]store.User, error)
	updateFn  func(context.Context, *store.User) error
	deleteFn  func(context.Context, int64) error
}

func (m *mockUserStore) Create(ctx context.Context, user *store.User) error {
	return m.createFn(ctx, user)
}

func (m *mockUserStore) GetByID(ctx context.Context, id int64) (*store.User, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockUserStore) GetAll(ctx context.Context) ([]store.User, error) {
	return m.getAllFn(ctx)
}

func (m *mockUserStore) Update(ctx context.Context, user *store.User) error {
	return m.updateFn(ctx, user)
}

func (m *mockUserStore) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

// Mock do PostStore (necessário para instanciar Storage)
type mockPostStore struct {
	createFn  func(context.Context, *store.Post) error
	getByIDFn func(context.Context, int64) (*store.Post, error)
	deleteFn  func(context.Context, int64) error
	updateFn  func(context.Context, *store.Post) error
}

func (m *mockPostStore) Create(ctx context.Context, post *store.Post) error {
	if m.createFn != nil {
		return m.createFn(ctx, post)
	}
	return nil
}

func (m *mockPostStore) GetByID(ctx context.Context, id int64) (*store.Post, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockPostStore) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockPostStore) Update(ctx context.Context, post *store.Post) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, post)
	}
	return nil
}
