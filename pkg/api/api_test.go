package api_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dakolli/opencode-go-kit/pkg/api"
	"github.com/dakolli/opencode-go-kit/pkg/client"
)

func TestWrapperCoverage(t *testing.T) {
	// 1. Get all methods defined on client.Invoker (the OpenAPI spec)
	invokerType := reflect.TypeOf((*client.Invoker)(nil)).Elem()
	invokerMethods := make(map[string]bool)
	for i := 0; i < invokerType.NumMethod(); i++ {
		invokerMethods[invokerType.Method(i).Name] = true
	}

	// 2. Get all methods defined on your wrapper *api.API
	apiType := reflect.TypeOf(&api.API{})
	apiMethods := make(map[string]bool)
	for i := 0; i < apiType.NumMethod(); i++ {
		name := apiType.Method(i).Name
		// Normalize name if your wrapper uses a prefix like "Get" or "Post"
		name = strings.TrimPrefix(name, "Get")
		name = strings.TrimPrefix(name, "Post")
		apiMethods[name] = true
	}

	// 3. Compare and output organized lists
	var implemented []string
	var missing []string

	for method := range invokerMethods {
		if apiMethods[method] {
			implemented = append(implemented, method)
		} else {
			missing = append(missing, method)
		}
	}

	total := len(invokerMethods)
	pct := (float64(len(implemented)) / float64(total)) * 100

	t.Logf("=== API Wrapper Implementation Report ===")
	t.Logf("Coverage: %d/%d operations wrapped (%.2f%%)", len(implemented), total, pct)

	t.Logf("\n[UNIMPLEMENTED ENDPOINTS]:")
	for _, m := range missing {
		t.Logf("  - %s", m)
	}

	t.Logf("\n[IMPLEMENTED ENDPOINTS]:")
	for _, m := range implemented {
		t.Logf("  - %s", m)
	}
}
