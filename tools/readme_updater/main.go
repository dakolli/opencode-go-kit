package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/dakolli/opencode-go-kit/pkg/api"
	"github.com/dakolli/opencode-go-kit/pkg/client"
)

func main() {
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
		name = strings.TrimPrefix(name, "Get")
		name = strings.TrimPrefix(name, "Post")
		apiMethods[name] = true
	}

	// 3. Compare and organize lists
	var implemented []string
	var missing []string

	for method := range invokerMethods {
		if apiMethods[method] {
			implemented = append(implemented, method)
		} else {
			missing = append(missing, method)
		}
	}

	sort.Strings(implemented)
	sort.Strings(missing)

	total := len(invokerMethods)
	pct := (float64(len(implemented)) / float64(total)) * 100

	// 4. Generate the Markdown segment
	var builder strings.Builder
	builder.WriteString("<!-- COVERAGE_START -->\n")
	builder.WriteString(fmt.Sprintf("[![API Coverage](https://img.shields.io/badge/Coverage-%.2f%%25-brightgreen)](#)\n\n", pct))
	builder.WriteString(fmt.Sprintf("We have wrapped **%d out of %d** (%0.2f%%) OpenAPI client methods in our clean, typed API wrapper layer.\n\n", len(implemented), total, pct))

	builder.WriteString("### Covered Endpoints\n\n")
	for _, m := range implemented {
		builder.WriteString(fmt.Sprintf("- [x] `%s`\n", m))
	}

	if len(missing) > 0 {
		builder.WriteString("\n### Uncovered Endpoints\n\n")
		for _, m := range missing {
			builder.WriteString(fmt.Sprintf("- [ ] `%s`\n", m))
		}
	}
	builder.WriteString("\n<!-- COVERAGE_END -->")

	newSegment := builder.String()

	// 5. Read README.md and replace the block
	readmePath := "README.md"
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		readmePath = "../../README.md" // Fallback if run from a different subfolder
	}

	contentBytes, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Printf("Error reading README: %v\n", err)
		os.Exit(1)
	}

	content := string(contentBytes)

	startTag := "<!-- COVERAGE_START -->"
	endTag := "<!-- COVERAGE_END -->"

	startIndex := strings.Index(content, startTag)
	endIndex := strings.Index(content, endTag)

	if startIndex == -1 || endIndex == -1 {
		fmt.Printf("Error: could not find coverage placeholder tags in %s\n", readmePath)
		os.Exit(1)
	}

	// Reconstruct the file contents
	updatedContent := content[:startIndex] + newSegment + content[endIndex+len(endTag):]

	err = os.WriteFile(readmePath, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Printf("Error writing updated README: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated %s with API coverage: %.2f%%\n", readmePath, pct)
}
