package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slash10k/internal/models"
	"strconv"
)

const (
	BaseUrl    = "https://true.torfstack.com/"
	DebtsUrl   = BaseUrl + "api/debt"
	JournalUrl = BaseUrl + "api/journal"
)

type DebtClient interface {
	AddDebt(ctx context.Context, name string, amount int64, reason string) error
	GetAllDebts(ctx context.Context) (*models.AllDebtsResponse, error)
	GetJournalEntries(ctx context.Context, name string) (*models.JournalEntries, error)
}

type AdminClient interface {
	DebtClient

	AddPlayer(ctx context.Context, name string) error
	DeletePlayer(ctx context.Context, name string) error
}

type client struct{}

var _ DebtClient = &client{}
var _ AdminClient = &client{}

func NewClient() *client {
	return &client{}
}

func (c client) AddDebt(ctx context.Context, name string, amount int64, reason string) error {
	var jsonData []byte
	if reason != "" {
		jsonData = []byte(fmt.Sprintf(`{"description": "%s"}`, reason))
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/%s/%s",
			DebtsUrl,
			name,
			strconv.FormatInt(amount, 10),
		),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("could not create debt post request: %w", err)
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("could not send debt post request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("debt post request was NOT successful (200): received %s", res.Status)
	}
	return nil
}

func (c client) GetAllDebts(ctx context.Context) (*models.AllDebtsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, DebtsUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create all debts get request: %w", err)
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("could not send all debts get request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("all debts get request was NOT successful (200): received %s", res.Status)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read all debts response: %w", err)
	}
	var debts models.AllDebtsResponse
	if err = json.Unmarshal(b, &debts); err != nil {
		return nil, fmt.Errorf("could not unmarshal all debts response: %w", err)
	}
	return &debts, nil
}

func (c client) GetJournalEntries(ctx context.Context, name string) (*models.JournalEntries, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s",
			JournalUrl,
			name,
		),
		nil,
	)
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("could not send journal entries get request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("journal entries get request was NOT successful (200): received %s", res.Status)
	}
	var entries models.JournalEntries
	if err = json.NewDecoder(res.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("could not decode journal entries response: %w", err)
	}
	return &entries, nil
}

func (c client) AddPlayer(ctx context.Context, name string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%sapi/admin/player/%s",
			BaseUrl,
			name,
		),
		nil,
	)
	req.Header.Set("Authorization", basicAuth("admin", os.Getenv("ADMIN_PASSWORD")))
	if err != nil {
		return fmt.Errorf("could not create add player request: %w", err)
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("could not send player add request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("player add request was NOT successful (204): received %s", res.Status)
	}
	return nil
}

func (c client) DeletePlayer(ctx context.Context, name string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf(
			"%sapi/admin/player/%s",
			BaseUrl,
			name,
		),
		nil,
	)
	req.Header.Set("Authorization", basicAuth("admin", os.Getenv("ADMIN_PASSWORD")))
	if err != nil {
		return fmt.Errorf("could not create player delete request: %w", err)
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("could not send player delete request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("player delete request was NOT successful (204): received %s", res.Status)
	}
	return nil
}
