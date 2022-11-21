package service

import "gotest/repository"

type WalletService interface {
	OpenAccount(wallet *repository.Wallet) error
	GetAccount(id uint64) (*repository.Wallet, error)
	Withdraw(id uint64, amount float32) (float32, error)
	Deposit(id uint64, amount float32) (float32, error)
}

type walletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) WalletService {
	return walletService{walletRepo: walletRepo}
}

func (s walletService) OpenAccount(wallet *repository.Wallet) error {
	err := s.walletRepo.Create(wallet)
	if err != nil {
		return NewErrorWalletUnexpected()
	}
	return nil
}

func (s walletService) GetAccount(id uint64) (*repository.Wallet, error) {
	wallet, err := s.walletRepo.Get(id)
	if err != nil {
		return nil, NewErrorWalletUnexpected()
	}

	if wallet.ID == 0 {
		return nil, NewErrorWalletNotFound()
	}

	return wallet, nil
}

func (s walletService) Withdraw(id uint64, amount float32) (float32, error) {
	wallet, err := s.GetAccount(id)
	if err != nil {
		return 0, err
	}

	if wallet.Balance < amount {
		return 0, NewErrorBadRequest("NOT ENOUGH MONEY")
	}

	wallet.Balance = wallet.Balance - amount
	err = s.walletRepo.Update(id, wallet)
	if err != nil {
		return 0, NewErrorWalletUnexpected()
	}

	return wallet.Balance, nil
}

func (s walletService) Deposit(id uint64, amount float32) (float32, error) {
	wallet, err := s.GetAccount(id)
	if err != nil {
		return 0, err
	}

	wallet.Balance = wallet.Balance + amount
	err = s.walletRepo.Update(id, wallet)
	if err != nil {
		return 0, NewErrorWalletUnexpected()
	}

	return wallet.Balance, nil
}
