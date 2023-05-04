// file: http_handler_test.go
package interfaces

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"evolve/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"evolve/services/transaction/domain"
)

type MockTransactionApplication struct {
	CreateTransferFunc func(ctx context.Context, transfer domain.Transfer) (*services.DefaultResponse, error)
}

func (m *MockTransactionApplication) CreateTransfer(ctx context.Context, transfer domain.Transfer) (*services.DefaultResponse, error) {
	return m.CreateTransferFunc(ctx, transfer)
}

func TestHTTPHandler_CreateTransfer(t *testing.T) {
	tt := []struct {
		name               string
		requestBody        domain.Transfer
		createTransferFunc func(ctx context.Context, transfer domain.Transfer) (*services.DefaultResponse, error)
		expectedResponse   *services.DefaultResponse
		expectedStatusCode int
	}{
		{
			name:        "successful transfer",
			requestBody: domain.Transfer{TransactionPin: "evolve"},
			createTransferFunc: func(ctx context.Context, transfer domain.Transfer) (*services.DefaultResponse, error) {
				return &services.DefaultResponse{Success: true, Message: "successful transfer"}, nil
			},
			expectedResponse:   &services.DefaultResponse{Success: true, Message: "successful transfer"},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "failed to decode request body",
			requestBody: domain.Transfer{},
			createTransferFunc: func(ctx context.Context, transfer domain.Transfer) (*services.DefaultResponse, error) {
				return &services.DefaultResponse{}, errors.New("error")
			},
			expectedResponse:   &services.DefaultResponse{Message: "unable to decode request body", Success: false},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "failed to create transfer",
			requestBody: domain.Transfer{},
			createTransferFunc: func(ctx context.Context, transfer domain.Transfer) (*services.DefaultResponse, error) {
				return nil, errors.New("error")
			},
			expectedResponse:   &services.DefaultResponse{Message: "error", Success: false},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := HTTPHandler{
				TransactionApplication: &MockTransactionApplication{
					CreateTransferFunc: tc.createTransferFunc,
				},
			}
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			handler.CreateTransfer(rr, req)

			type data struct {
				Data services.DefaultResponse `json:"data"`
			}
			var response data
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				return
			}
			if response.Data.Message != tc.expectedResponse.Message {
				t.Errorf("expected message %v, but got %v", tc.expectedResponse.Message, response.Data.Message)
			}

			if response.Data.Success != tc.expectedResponse.Success {
				t.Errorf("expected success %v, but got %v", tc.expectedResponse.Success, response.Data.Success)
			}

			if rr.Code != tc.expectedStatusCode {
				t.Errorf("expected status code %v, but got %v", tc.expectedStatusCode, rr.Code)
			}
		})
	}
}
