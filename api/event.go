package api

import (
	"context"
	"net/url"
	"strconv"

	"github.com/celerfi/cedra-go-kit/client"
	"github.com/celerfi/cedra-go-kit/types"
)

type EventAPI struct {
	c *client.Client
}

func newEventAPI(c *client.Client) *EventAPI {
	return &EventAPI{c: c}
}

func (e *EventAPI) GetEventsByEventType(ctx context.Context, eventType string, limit, start *uint64) ([]types.Event, error) {
	params := paginationParams(limit, start)
	var events []types.Event
	if err := e.c.Get(ctx, "/events/"+url.PathEscape(eventType), params, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (e *EventAPI) GetAccountEventsByCreationNumber(ctx context.Context, address string, creationNumber uint64, limit, start *uint64) ([]types.Event, error) {
	params := paginationParams(limit, start)
	var events []types.Event
	path := "/accounts/" + address + "/events/" + strconv.FormatUint(creationNumber, 10)
	if err := e.c.Get(ctx, path, params, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (e *EventAPI) GetAccountEventsByEventType(ctx context.Context, address, eventType string, limit, start *uint64) ([]types.Event, error) {
	params := paginationParams(limit, start)
	var events []types.Event
	path := "/accounts/" + address + "/events/" + url.PathEscape(eventType)
	if err := e.c.Get(ctx, path, params, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func paginationParams(limit, start *uint64) url.Values {
	params := url.Values{}
	if limit != nil {
		params.Set("limit", strconv.FormatUint(*limit, 10))
	}
	if start != nil {
		params.Set("start", strconv.FormatUint(*start, 10))
	}
	return params
}
