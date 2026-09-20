package service

import (
	"context"
	"fmt"
	"hexagonal/ccgo01/internal/core/model"
	"hexagonal/ccgo01/internal/core/port"
	"sync"
)

type WalletService struct {
	repo port.WalletRepositry
	mu   sync.RWMutex
}

func NewWalletService(repo port.WalletRepositry) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) getWallet(ctx context.Context, id string) (*model.Wallet, error) {
	wallet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet %s not found", id)
	}
	return wallet, nil
}

func (s *WalletService) CreateWallet(ctx context.Context, userID string) (*model.Wallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("check existing wallet %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user %s already has a wallet", userID)
	}

	wallet := &model.Wallet{
		UserID:  userID,
		Balance: 0,
	}
	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("create wallet: %w", err)
	}
	return wallet, nil
}

func (s *WalletService) GetWallet(ctx context.Context, id string) (*model.Wallet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getWallet(ctx, id)
}

func (s *WalletService) Deposit(ctx context.Context, id string, amount float64) (*model.Wallet, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("deposit amount must be positive")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, err := s.getWallet(ctx, id)
	if err != nil {
		return nil, err
	}
	newBalance := wallet.Balance + amount
	if err := s.repo.UpdateBalance(ctx, id, newBalance); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}
	wallet.Balance = newBalance
	return wallet, nil
}

func (s *WalletService) Withdraw(ctx context.Context, id string, amount float64) (*model.Wallet, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("withdraw amount must be positive")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, err := s.getWallet(ctx, id)
	if err != nil {
		return nil, err
	}
	if wallet.Balance < amount {
		return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", wallet.Balance, amount)
	}
	newBalance := wallet.Balance - amount
	if err := s.repo.UpdateBalance(ctx, id, newBalance); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	wallet.Balance = newBalance
	return wallet, nil
}

func (s *WalletService) DeleteWalelt(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete wallet: %w", err)
	}
	return nil
}
