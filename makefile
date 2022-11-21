unit-service:
	go test gotest/service -v -cover -tags=unit

unit-handler:
	go test gotest/handler -v -cover -tags=unit