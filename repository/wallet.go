package repository

type Wallet struct {
	ID      uint64  `json:"id" `
	Name    string  `json:"name" validate:"required" `
	Balance float32 `json:"balance" validate:"required,number"`
}

type WalletRepository interface {
	Get(id uint64) (*Wallet, error)
	Create(wallet *Wallet) error
	Update(id uint64, wallet *Wallet) error
}
