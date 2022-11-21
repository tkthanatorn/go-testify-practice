// go:build unit
package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotest/handler"
	"gotest/repository"
	"gotest/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestOpenAccount(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		// arrange
		wallet := repository.Wallet{
			Name:    "John Doe",
			Balance: 1000,
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(wallet); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("OpenAccount", &wallet).Return(nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank", handler.OpenAcount)
		req := httptest.NewRequest(http.MethodPost, "/bank", &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, resp.StatusCode, 201)
	})

	t.Run("Unprocessable Entity Error", func(t *testing.T) {
		// arrange
		walletMap := map[string]interface{}{
			"balance": "Hi",
		}
		wallet := repository.Wallet{}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(walletMap); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("OpenAccount", &wallet).Return(nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank", handler.OpenAcount)
		req := httptest.NewRequest(http.MethodPost, "/bank", &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, resp.StatusCode, 422)
		serv.AssertNotCalled(t, "OpenAccount", &wallet)
	})

	t.Run("Process Error", func(t *testing.T) {
		// arrange
		wallet := repository.Wallet{
			Name:    "John Doe",
			Balance: 1000,
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(wallet); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("OpenAccount", &wallet).Return(service.NewErrorWalletUnexpected())
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank", handler.OpenAcount)
		req := httptest.NewRequest(http.MethodPost, "/bank", &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, resp.StatusCode, 500)
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		var id uint64 = 1
		wallet := repository.Wallet{
			ID:      1,
			Name:    "John Doe",
			Balance: 1000,
		}

		serv := service.NewWalletServiceMock()
		serv.On("GetAccount", id).Return(&wallet, nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Get("/bank/:id", handler.GetAccount)
		url := fmt.Sprintf("/bank/%v", id)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Unprocessable Entity", func(t *testing.T) {
		var id uint64 = 1
		wallet := repository.Wallet{
			ID:      1,
			Name:    "John Doe",
			Balance: 1000,
		}

		serv := service.NewWalletServiceMock()
		serv.On("GetAccount", id).Return(&wallet, nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Get("/bank/:id", handler.GetAccount)
		url := fmt.Sprintf("/bank/%v", "error")
		req := httptest.NewRequest(http.MethodGet, url, nil)
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 422, resp.StatusCode)
	})

	t.Run("Process Error", func(t *testing.T) {
		var id uint64 = 1
		wallet := repository.Wallet{
			ID:      0,
			Name:    "",
			Balance: 0,
		}

		serv := service.NewWalletServiceMock()
		serv.On("GetAccount", id).Return(&wallet, service.NewErrorWalletNotFound())
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Get("/bank/:id", handler.GetAccount)
		url := fmt.Sprintf("/bank/%v", id)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestWithdraw(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		var id uint64 = 1
		transaction := handler.TransactionRequest{
			Amount: 200,
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(transaction); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("Withdraw", id, transaction.Amount).Return(float32(800), nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank/withdraw/:id", handler.Withdraw)
		url := fmt.Sprintf("/bank/withdraw/%v", id)
		req := httptest.NewRequest(http.MethodPost, url, &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Pass Param Error", func(t *testing.T) {
		var id uint64 = 1
		transaction := handler.TransactionRequest{
			Amount: 200,
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(transaction); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("Withdraw", id, transaction.Amount).Return(float32(800), nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank/withdraw/:id", handler.Withdraw)
		url := fmt.Sprintf("/bank/withdraw/%v", "error")
		req := httptest.NewRequest(http.MethodPost, url, &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 422, resp.StatusCode)
		serv.AssertNotCalled(t, "Withdraw")
	})

	t.Run("Pass Body Error", func(t *testing.T) {
		var id uint64 = 1
		transaction := map[string]interface{}{
			"amount": "Hello World",
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(transaction); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("Withdraw", id, float32(200)).Return(float32(800), nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank/withdraw/:id", handler.Withdraw)
		url := fmt.Sprintf("/bank/withdraw/%v", id)
		req := httptest.NewRequest(http.MethodPost, url, &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 422, resp.StatusCode)
		serv.AssertNotCalled(t, "Withdraw")
	})

	t.Run("Withdraw Error", func(t *testing.T) {
		var id uint64 = 1
		transaction := handler.TransactionRequest{
			Amount: 200,
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(transaction); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("Withdraw", id, transaction.Amount).Return(float32(0), service.NewErrorWalletUnexpected())
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank/withdraw/:id", handler.Withdraw)
		url := fmt.Sprintf("/bank/withdraw/%v", id)
		req := httptest.NewRequest(http.MethodPost, url, &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestDeposit(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		var id uint64 = 1
		transaction := handler.TransactionRequest{
			Amount: 200,
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(transaction); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("Deposit", id, transaction.Amount).Return(float32(1200), nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Get("/bank/deposit/:id", handler.Deposit)
		url := fmt.Sprintf("/bank/deposit/%v", id)
		req := httptest.NewRequest(http.MethodGet, url, &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Pass Param Error", func(t *testing.T) {
		var id uint64 = 1
		transaction := handler.TransactionRequest{
			Amount: 200,
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(transaction); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("Deposit", id, transaction.Amount).Return(float32(1200), nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank/deposit/:id", handler.Deposit)
		url := fmt.Sprintf("/bank/deposit/%v", "error")
		req := httptest.NewRequest(http.MethodPost, url, &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 422, resp.StatusCode)
		serv.AssertNotCalled(t, "Deposit")
	})

	t.Run("Pass Body Error", func(t *testing.T) {
		var id uint64 = 1
		transaction := map[string]interface{}{
			"amount": "Hello World",
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(transaction); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("Deposit", id, float32(200)).Return(float32(1200), nil)
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank/deposit/:id", handler.Deposit)
		url := fmt.Sprintf("/bank/deposit/%v", id)
		req := httptest.NewRequest(http.MethodPost, url, &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 422, resp.StatusCode)
		serv.AssertNotCalled(t, "Deposit")
	})

	t.Run("Deposit Error", func(t *testing.T) {
		var id uint64 = 1
		transaction := handler.TransactionRequest{
			Amount: 200,
		}

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(transaction); err != nil {
			t.Log(err)
		}

		serv := service.NewWalletServiceMock()
		serv.On("Deposit", id, transaction.Amount).Return(float32(0), service.NewErrorWalletUnexpected())
		handler := handler.NewWalletHandler(serv)

		app := fiber.New()
		app.Post("/bank/deposit/:id", handler.Deposit)
		url := fmt.Sprintf("/bank/deposit/%v", id)
		req := httptest.NewRequest(http.MethodPost, url, &buff)
		req.Header.Set("Content-Type", "application/json")
		// act
		resp, _ := app.Test(req)
		// assert
		assert.Equal(t, 500, resp.StatusCode)
	})
}
