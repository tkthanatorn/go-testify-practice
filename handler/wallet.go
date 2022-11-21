package handler

import (
	"fmt"
	"gotest/repository"
	"gotest/service"

	"github.com/gofiber/fiber/v2"
)

type walletHandler struct {
	walletServ service.WalletService
}

func NewWalletHandler(walletServ service.WalletService) walletHandler {
	return walletHandler{walletServ: walletServ}
}

func ResponseError(c *fiber.Ctx, err error) error {
	switch v := err.(type) {
	case service.WalletError:
		return c.Status(v.Code).SendString(v.Message)
	default:
		return c.SendStatus(500)
	}
}

func (h walletHandler) OpenAcount(c *fiber.Ctx) error {
	wallet := repository.Wallet{}
	if err := c.BodyParser(&wallet); err != nil {
		e := service.NewErrorUnprocessableEntity()
		return ResponseError(c, e)
	}

	err := h.walletServ.OpenAccount(&wallet)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.SendStatus(201)
}

func (h walletHandler) GetAccount(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := service.NewErrorUnprocessableEntity()
		return ResponseError(c, e)
	}

	wallet, err := h.walletServ.GetAccount(uint64(id))
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(200).JSON(wallet)
}

type TransactionRequest struct {
	Amount float32 `json:"amount" validate:"required,number"`
}

func (h walletHandler) Withdraw(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := service.NewErrorUnprocessableEntity()
		return ResponseError(c, e)
	}

	transaction := TransactionRequest{}
	if err := c.BodyParser(&transaction); err != nil {
		e := service.NewErrorUnprocessableEntity()
		return ResponseError(c, e)
	}

	change, err := h.walletServ.Withdraw(uint64(id), transaction.Amount)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(200).SendString(fmt.Sprintf("Balance: %v", change))
}

func (h walletHandler) Deposit(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := service.NewErrorUnprocessableEntity()
		return ResponseError(c, e)
	}

	transaction := TransactionRequest{}
	if err := c.BodyParser(&transaction); err != nil {
		e := service.NewErrorUnprocessableEntity()
		return ResponseError(c, e)
	}

	change, err := h.walletServ.Deposit(uint64(id), transaction.Amount)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(200).SendString(fmt.Sprintf("Balance: %v", change))
}
