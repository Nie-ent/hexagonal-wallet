package repository

import (
	"context"
	"hexagonal/ccgo01/internal/core/model"
	"time"

	"gorm.io/gorm"
)

type WalletGORM struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `gorm:"column:user_id;not null;uniqueIndex"`
	Balance   float64   `gorm:"column:balance;type:decimal(18,2);not null;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}

func (WalletGORM) TableName() string {
	return "wallets"
}

func (w WalletGORM) toDomain() model.Wallet {
	return model.Wallet{
		ID:        w.ID,
		UserID:    w.UserID,
		Balance:   w.Balance,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

func fromDomain(w model.Wallet) WalletGORM {
	return WalletGORM{
		ID:        w.ID,
		UserID:    w.UserID,
		Balance:   w.Balance,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

type PostgresWalletRepository struct {
	db *gorm.DB
}

func NewPostgresWalletRepository(db *gorm.DB) *PostgresWalletRepository {
	return &PostgresWalletRepository{db: db}
}

func (r *PostgresWalletRepository) Create(ctx context.Context, wallet *model.Wallet) error {
	gormWallet := fromDomain(*wallet)
	result := r.db.WithContext(ctx).Create(&gormWallet)
	if result.Error != nil {
		return result.Error
	}
	*wallet = gormWallet.toDomain()
	return nil
}

func (r *PostgresWalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	var gormWallet WalletGORM
	result := r.db.WithContext(ctx).First(&gormWallet, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	domainWallet := gormWallet.toDomain()
	return &domainWallet, nil
}

func (r *PostgresWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	var gormWallet WalletGORM
	result := r.db.WithContext(ctx).First(&gormWallet, "user_id = ?", userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	domainWallet := gormWallet.toDomain()
	return &domainWallet, nil
}

func (r *PostgresWalletRepository) UpdateBalance(ctx context.Context, id string, balance float64) error {
	result := r.db.WithContext(ctx).
		Model(&WalletGORM{}).
		Where("id = ?", id).
		Update("balance", balance)
	return result.Error
}

func (r *PostgresWalletRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&WalletGORM{}, "id = ?", id)
	return result.Error
}
