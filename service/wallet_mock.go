package service

import (
	"gotest/repository"

	"github.com/stretchr/testify/mock"
)

type walletServiceMock struct {
	mock.Mock
}

func NewWalletServiceMock() *walletServiceMock {
	return &walletServiceMock{}
}

func (m *walletServiceMock) OpenAccount(wallet *repository.Wallet) error {
	c := m.Called(wallet)
	return c.Error(0)
}

func (m *walletServiceMock) GetAccount(id uint64) (*repository.Wallet, error) {
	c := m.Called(id)
	return c.Get(0).(*repository.Wallet), c.Error(1)
}

func (m *walletServiceMock) Withdraw(id uint64, amount float32) (float32, error) {
	c := m.Called(id, amount)
	return c.Get(0).(float32), c.Error(1)
}

func (m *walletServiceMock) Deposit(id uint64, amount float32) (float32, error) {
	c := m.Called(id, amount)
	return c.Get(0).(float32), c.Error(1)
}
