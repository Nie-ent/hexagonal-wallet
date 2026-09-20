package service_test

import (
	"context"
	"hexagonal/ccgo01/internal/core/model"
	"hexagonal/ccgo01/internal/core/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	wallets map[string]*model.Wallet
}

func (m *mockRepo) Create(_ context.Context, wallet *model.Wallet) error {
	wallet.ID = "mock-id-001"
	m.wallets[wallet.ID] = wallet
	return nil
}

func (m *mockRepo) GetByID(_ context.Context, id string) (*model.Wallet, error) {
	w, ok := m.wallets[id]
	if !ok {
		return nil, nil
	}
	return w, nil
}

func (m *mockRepo) GetByUserID(_ context.Context, userID string) (*model.Wallet, error) {
	for _, w := range m.wallets {
		if w.UserID == userID {
			return w, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) UpdateBalance(_ context.Context, id string, balance float64) error {
	if w, ok := m.wallets[id]; ok {
		w.Balance = balance
	}
	return nil
}

func (m *mockRepo) Delete(_ context.Context, id string) error {
	delete(m.wallets, id)
	return nil
}

func newMockRepo() *mockRepo {
	return &mockRepo{wallets: make(map[string]*model.Wallet)}
}

func TestCreateWallet_Success(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, err := svc.CreateWallet(context.Background(), "user-1")
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, "user-1", wallet.UserID)
	assert.Equal(t, 0.0, wallet.Balance)
}

func TestCreateWallet_DuplicateUser(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	_, err := svc.CreateWallet(context.Background(), "user-1")
	assert.NoError(t, err)

	_, err = svc.CreateWallet(context.Background(), "user-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already has a wallet")
}

func TestDeposit_Success(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	updated, err := svc.Deposit(context.Background(), wallet.ID, 100.0)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, updated.Balance)
}

func TestDeposit_NegativeAmount(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	_, err := svc.Deposit(context.Background(), wallet.ID, -50.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")
}

func TestWithdraw_Success(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	svc.Deposit(context.Background(), wallet.ID, 200.0)
	updated, err := svc.Withdraw(context.Background(), wallet.ID, 50.0)
	assert.NoError(t, err)
	assert.Equal(t, 150.0, updated.Balance)
}

func TestWithdraw_InsufficientBalance(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	_, err := svc.Withdraw(context.Background(), wallet.ID, 10.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient balance")
}

func TestGetWallet_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	_, err := svc.GetWallet(context.Background(), "non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDeleteWallet_Success(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)
	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	err := svc.DeleteWalelt(context.Background(), wallet.ID)
	assert.NoError(t, err)

	_, err = svc.GetWallet(context.Background(), wallet.ID)
	assert.Error(t, err)
}
