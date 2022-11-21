package service

type WalletError struct {
	Code    int
	Message string
}

func (e WalletError) Error() string {
	return e.Message
}

func NewErrorWalletNotFound() WalletError {
	return WalletError{
		Code:    404,
		Message: "WALLET NOT FOUND",
	}
}

func NewErrorWalletUnexpected() WalletError {
	return WalletError{
		Code:    500,
		Message: "UNEXPECTED ERROR",
	}
}

func NewErrorBadRequest(msg string) WalletError {
	return WalletError{
		Code:    400,
		Message: msg,
	}
}

func NewErrorUnprocessableEntity() WalletError {
	return WalletError{
		Code:    422,
		Message: "INVALID PARAMETER",
	}
}
