package labsd

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
)

type Client struct{ Socket string }

func (c Client) Call(ctx context.Context, operation, app string) error {
	_, err := c.StatusOrCall(ctx, operation, app)
	return err
}

func (c Client) Install(ctx context.Context, app, compose string) error {
	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", c.Socket)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := json.NewEncoder(conn).Encode(Request{Operation: "InstallApp", App: app, Compose: compose}); err != nil {
		return err
	}
	var response Response
	if err := json.NewDecoder(conn).Decode(&response); err != nil {
		return err
	}
	if !response.OK {
		return errors.New(response.Message)
	}
	return nil
}

func (c Client) List(ctx context.Context) ([]string, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", c.Socket)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := json.NewEncoder(conn).Encode(Request{Operation: "ListApps"}); err != nil {
		return nil, err
	}
	var response Response
	if err := json.NewDecoder(conn).Decode(&response); err != nil {
		return nil, err
	}
	if !response.OK {
		return nil, errors.New(response.Message)
	}
	return response.Apps, nil
}

func hasApp(names []string, id string) bool {
	for _, name := range names {
		if strings.EqualFold(name, id) {
			return true
		}
	}
	return false
}

func (c Client) Status(ctx context.Context, app string) (string, error) {
	return c.StatusOrCall(ctx, "StatusApp", app)
}

func (c Client) Logs(ctx context.Context, app string, lines int) (string, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", c.Socket)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if err := json.NewEncoder(conn).Encode(Request{Operation: "LogsApp", App: app, Lines: lines}); err != nil {
		return "", err
	}
	var response Response
	if err := json.NewDecoder(conn).Decode(&response); err != nil {
		return "", err
	}
	if !response.OK {
		return "", errors.New(response.Message)
	}
	return response.Logs, nil
}

func (c Client) StatusOrCall(ctx context.Context, operation, app string) (string, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", c.Socket)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if err := json.NewEncoder(conn).Encode(Request{Operation: operation, App: app}); err != nil {
		return "", err
	}
	var response Response
	if err := json.NewDecoder(conn).Decode(&response); err != nil {
		return "", err
	}
	if !response.OK {
		return "", errors.New(response.Message)
	}
	return response.Status, nil
}
