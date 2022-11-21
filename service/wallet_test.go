// go:build unit
package service_test

import (
	"errors"
	"gotest/repository"
	"gotest/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAccount(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		// arrange
		successInput := repository.Wallet{
			Name:    "John Doe",
			Balance: 1000,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Create", &successInput).Return(nil)
		serv := service.NewWalletService(repo)

		// act
		err := serv.OpenAccount(&successInput)
		// assert
		assert.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// arrange
		errorInput := repository.Wallet{
			Name:    "John Doe",
			Balance: 1000,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.
			On("Create", &errorInput).
			Return(service.NewErrorWalletUnexpected())
		serv := service.NewWalletService(repo)

		// act
		err := serv.OpenAccount(&errorInput)
		// assert
		assert.ErrorIs(t, err, service.NewErrorWalletUnexpected())
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("Succesful", func(t *testing.T) {
		// arrange
		var input uint64 = 1
		wallet := repository.Wallet{
			ID:      1,
			Name:    "John Doe",
			Balance: 1000,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Get", input).Return(&wallet, nil)
		serv := service.NewWalletService(repo)
		// act
		result, _ := serv.GetAccount(input)
		// assert
		assert.Equal(t, wallet.ID, result.ID)
	})

	t.Run("ErrorNotFound", func(t *testing.T) {
		// arrange
		var input uint64 = 2
		wallet := repository.Wallet{
			ID:      0,
			Name:    "",
			Balance: 0,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Get", input).Return(&wallet, nil)
		serv := service.NewWalletService(repo)
		// act
		result, err := serv.GetAccount(input)
		_ = result
		// assert
		assert.ErrorIs(t, err, service.NewErrorWalletNotFound())
	})

	t.Run("ErrorUnexpected", func(t *testing.T) {
		// arrange
		var input uint64 = 3
		wallet := repository.Wallet{
			ID:      0,
			Name:    "",
			Balance: 0,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Get", input).Return(&wallet, errors.New("database down"))
		serv := service.NewWalletService(repo)
		// act
		result, err := serv.GetAccount(input)
		_ = result
		// assert
		assert.ErrorIs(t, err, service.NewErrorWalletUnexpected())
	})
}

func TestWithdrawSuccessful(t *testing.T) {
	wallet := repository.Wallet{
		ID:      1,
		Name:    "John Doe",
		Balance: 1000,
	}

	repo := repository.NewWalletRepositoryMock()
	repo.On("Get", uint64(1)).Return(&wallet, nil)
	repo.On("Update", uint64(1)).Return(nil)
	serv := service.NewWalletService(repo)

	type testCase struct {
		Name     string
		ID       uint64
		Amount   float32
		Expected float32
	}

	tests := []testCase{
		{Name: "Withdraw 200", ID: 1, Amount: 200, Expected: 800},
		{Name: "Withdraw All", ID: 1, Amount: 1000, Expected: 0},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, _ := serv.Withdraw(test.ID, test.Amount)
			assert.Equal(t, result, test.Expected)
		})
	}
}

func TestWithdrawError(t *testing.T) {
	t.Run("Error From GetAccount", func(t *testing.T) {
		// arrange
		var id uint64 = 1
		var amount float32 = 200
		wallet := repository.Wallet{
			ID:      0,
			Name:    "",
			Balance: 0,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Get", uint64(1)).Return(&wallet, nil)
		serv := service.NewWalletService(repo)
		// act
		result, err := serv.Withdraw(id, amount)
		_ = result
		// assert
		assert.Error(t, err)
	})

	t.Run("Error Not Enough Money", func(t *testing.T) {
		// arrange
		var id uint64 = 1
		var amount float32 = 2000
		wallet := repository.Wallet{
			ID:      1,
			Name:    "John Doe",
			Balance: 1000,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Get", id).Return(&wallet, nil)
		repo.On("Update", id).Return(nil)
		serv := service.NewWalletService(repo)
		// act
		result, err := serv.Withdraw(id, amount)
		_ = result
		// assert
		assert.ErrorIs(t, err, service.NewErrorBadRequest("NOT ENOUGH MONEY"))
	})

	t.Run("Error Unexpected", func(t *testing.T) {
		// arrange
		var id uint64 = 1
		var amount float32 = 1000
		wallet := repository.Wallet{
			ID:      1,
			Name:    "John Doe",
			Balance: 1000,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Get", id).Return(&wallet, nil)
		repo.On("Update", id).Return(errors.New("database down"))
		serv := service.NewWalletService(repo)
		// act
		result, err := serv.Withdraw(id, amount)
		_ = result
		// assert
		assert.ErrorIs(t, err, service.NewErrorWalletUnexpected())
	})
}

func TestDepositSuccessful(t *testing.T) {
	wallet := repository.Wallet{
		ID:      1,
		Name:    "John Doe",
		Balance: 1000,
	}

	repo := repository.NewWalletRepositoryMock()
	repo.On("Get", uint64(1)).Return(&wallet, nil)
	repo.On("Update", uint64(1)).Return(nil)
	serv := service.NewWalletService(repo)

	type testCase struct {
		Name     string
		ID       uint64
		Amount   float32
		Expected float32
	}

	tests := []testCase{
		{Name: "Deposit 200", ID: 1, Amount: 200, Expected: 1200},
		{Name: "Deposit 1000", ID: 1, Amount: 1000, Expected: 2200},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, _ := serv.Deposit(test.ID, test.Amount)
			assert.Equal(t, result, test.Expected)
		})
	}
}

func TestDepositError(t *testing.T) {
	t.Run("Error From GetAccount", func(t *testing.T) {
		// arrange
		var id uint64 = 1
		var amount float32 = 200
		wallet := repository.Wallet{
			ID:      0,
			Name:    "",
			Balance: 0,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Get", uint64(1)).Return(&wallet, nil)
		serv := service.NewWalletService(repo)
		// act
		result, err := serv.Deposit(id, amount)
		_ = result
		// assert
		assert.Error(t, err)
	})

	t.Run("Error Unexpected", func(t *testing.T) {
		// arrange
		var id uint64 = 1
		var amount float32 = 1000
		wallet := repository.Wallet{
			ID:      1,
			Name:    "John Doe",
			Balance: 1000,
		}

		repo := repository.NewWalletRepositoryMock()
		repo.On("Get", id).Return(&wallet, nil)
		repo.On("Update", id).Return(errors.New("database down"))
		serv := service.NewWalletService(repo)
		// act
		result, err := serv.Deposit(id, amount)
		_ = result
		// assert
		assert.ErrorIs(t, err, service.NewErrorWalletUnexpected())
	})
}
