package tengoutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

const (
	maxAllocs = 5000
	maxConsts = 500
)

func NewSecureScript(input []byte, globals map[string]interface{}) (*tengo.Script, error) {
	s := tengo.NewScript(input)

	modules := stdlib.GetModuleMap(
		// `os` is excluded, should *not* be importable from script.
		"math", "text", "times", "rand", "fmt", "json", "base64", "hex", "enum",
	)
	modules.AddBuiltinModule("http", createHTTPModule())
	s.SetImports(modules)
	s.SetMaxAllocs(maxAllocs)
	s.SetMaxConstObjects(maxConsts)

	for name, v := range globals {
		if err := s.Add(name, v); err != nil {
			return nil, fmt.Errorf("new secure script: declare globals: %w", err)
		}
	}

	return s, nil
}

func createHTTPModule() map[string]tengo.Object {
	return map[string]tengo.Object{
		"get": &tengo.UserFunction{
			Name: "get",
			Value: func(args ...tengo.Object) (tengo.Object, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("expected 1 argument, got %d", len(args))
				}

				url, ok := tengo.ToString(args[0])
				if !ok {
					return nil, fmt.Errorf("expected argument 1 (URL) to be a string")
				}

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				if err != nil {
					return nil, err
				}
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}

				return &tengo.Map{
					Value: map[string]tengo.Object{
						"body": &tengo.String{Value: string(body)},
						"code": &tengo.Int{Value: int64(resp.StatusCode)},
					},
				}, nil
			},
		},
	}
}
