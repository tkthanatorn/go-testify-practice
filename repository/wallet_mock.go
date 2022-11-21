package repository

import "github.com/stretchr/testify/mock"

type walletRepositoryMock struct {
	mock.Mock
}

func NewWalletRepositoryMock() *walletRepositoryMock {
	return &walletRepositoryMock{}
}

func (m *walletRepositoryMock) Get(id uint64) (*Wallet, error) {
	c := m.Called(id)
	return c.Get(0).(*Wallet), c.Error(1)
}

func (m *walletRepositoryMock) Create(wallet *Wallet) error {
	c := m.Called(wallet)
	return c.Error(0)
}

func (m *walletRepositoryMock) Update(id uint64, wallet *Wallet) error {
	c := m.Called(id)
	return c.Error(0)
}
