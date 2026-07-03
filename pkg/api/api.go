package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/dakolli/opencode-go-kit/pkg/client"
)

type API struct {
	Client *client.Client
}

type ClientCfg struct {
	URL      string
	Username string
	Password string
}

// use env vars
func NewClientCFG(url, username, password string) (ClientCfg, error) {
	if url == "" {
		url = os.Getenv("OPENCODE_URL")
	}
	if username == "" {
		username = os.Getenv("OPENCODE_USERNAME")
	}
	if password == "" {
		password = os.Getenv("OPENCODE_PASSWORD")
	}
	if url == "" || username == "" || password == "" {
		return ClientCfg{}, fmt.Errorf("OPENCODE_URL, OPENCODE_USERNAME, and OPENCODE_PASSWORD must be set")
	}
	return ClientCfg{
		URL:      url,
		Username: username,
		Password: password,
	}, nil
}

func NewAPI(config ClientCfg) (*API, error) {
	client, err := client.NewClient(config.URL,
		client.WithRequestEditor(func(ctx context.Context, req *http.Request) error {
			req.SetBasicAuth(config.Username, config.Password)
			return nil
		}))
	if err != nil {
		return nil, err
	}
	return &API{Client: client}, nil
}

// Extract asserts that the result of an API call was successful and converts it
// to the expected concrete type T. If the API call failed, it returns an error.
func Extract[T any, R any](res R, err error) (T, error) {
	var zero T
	if err != nil {
		return zero, fmt.Errorf("API transport error: %w", err)
	}

	val, ok := any(res).(T)
	if !ok {
		// If the response is a BadRequestError or any other error struct
		// that supports a Validate() string method, let's print that.
		if validationErr, ok := any(res).(interface{ Validate() string }); ok {
			return zero, fmt.Errorf("API bad request: %s", validationErr.Validate())
		}
		return zero, fmt.Errorf("unexpected API response type: %T (expected %T)", res, zero)
	}

	return val, nil
}

func GetRoutes(c *client.Client) []string {
	var routes []string
	if c == nil {
		return routes
	}

	// We get the reflect.Type of the client.Invoker interface.
	// Since client.Invoker is an interface, passing a nil pointer of its pointer-type
	// (*client.Invoker)(nil) and calling Elem() returns the reflect.Type representing
	// the interface itself.
	invokerType := reflect.TypeOf((*client.Invoker)(nil)).Elem()

	for i := 0; i < invokerType.NumMethod(); i++ {
		method := invokerType.Method(i)
		routes = append(routes, method.Name)
	}

	return routes
}
