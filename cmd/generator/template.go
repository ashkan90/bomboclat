// generator/templates/template.go
package generator

const baseTemplate = `
package {{.Package}}

import (
   "context"
   "fmt"
   "sync"
   "time"
    {{range .Imports}}
    "{{.}}"
    {{- end}}
)

// Referencing custom types from source package
{{range .Types}}
type {{.Name}} {{.Definition}}
{{end}}

// Option is the interface all spawn options must implement
type Option interface {
   opt()
}

// SpawnConfig represents configuration for spawn operations
type SpawnConfig struct {
   Timeout time.Duration
   Retry   RetryConfig
}

// RetryConfig represents retry configuration
type RetryConfig struct {
   MaxAttempts int
   Delay       time.Duration
   BackoffFunc func(attempt int) time.Duration
}

// Default configurations
var (
   DefaultSpawnConfig = SpawnConfig{
   	Timeout: 30 * time.Second,
   	Retry: RetryConfig{
   		MaxAttempts: 3,
   		Delay:      100 * time.Millisecond,
   	},
   }
)

{{range .Interfaces}}
{{$interface := .}}
{{range .Methods}}
// {{.Name}}SpawnOption represents spawn configuration for {{$interface.Name}}.{{.Name}}
type {{.Name}}SpawnOption struct {
   fn     {{.FuncType}}
   args   []interface{}
   config SpawnConfig
}

func (o {{.Name}}SpawnOption) opt() {}

// {{.Name}}Spawn creates a spawn option for {{$interface.Name}}.{{.Name}}
func {{.Name}}Spawn(fn {{.FuncType}}, args ...interface{}) Option {
   return &{{.Name}}SpawnOption{
   	fn:     fn,
   	args:   args,
   	config: DefaultSpawnConfig,
   }
}

// {{.Name}}Handler handles execution of {{$interface.Name}}.{{.Name}}
type {{.Name}}Handler struct{}

func (h *{{.Name}}Handler) Execute(ctx context.Context, opt Option) (interface{}, error) {
    o, ok := opt.(*{{.Name}}SpawnOption)
    if !ok {
        return nil, fmt.Errorf("invalid option type: %T", opt)
    }

    if len(o.args) != {{len .Params}} {
        return nil, fmt.Errorf("expected {{len .Params}} arguments, got %d", len(o.args))
    }

    {{range $i, $p := .Params}}
    arg{{$i}}, ok := o.args[{{$i}}].({{$p.Type}})
    if !ok {
        return nil, fmt.Errorf("argument {{$i}} must be {{$p.Type}}")
    }
    {{end}}

    // Execute with timeout
    resultCh := make(chan struct {
        {{if len .Returns}}
        {{range $i, $r := .Returns}}
        r{{$i}} {{$r.Type}}
        {{end}}
        {{end}}
    }, 1)

	go func() {
		{{$method := .}}
		{{if len .Returns}}
		{{range $i, $r := .Returns}}r{{$i}}{{if not (last $i $method.Returns)}}, {{end}}{{end}} := {{end}}o.fn({{range $i, $_ := .Params}}arg{{$i}}, {{end}})
		resultCh <- struct {
			{{if len $method.Returns}}
			{{range $i, $r := $method.Returns}}
			r{{$i}} {{$r.Type}}
			{{end}}
			{{end}}
		}{ {{if len $method.Returns}}{{range $i, $r := $method.Returns}}r{{$i}}, {{end}}{{end}} }
	}()

    select {
    case result := <-resultCh:
        {{if hasError $method.Returns}}
		if result.r{{lastIndex $method.Returns}} != nil {
			return nil, result.r{{lastIndex $method.Returns}}
		}
		{{end}}
		{{if eq (len $method.Returns) 0}}
		return nil, nil
		{{else if eq (len $method.Returns) 1}}
			{{if isError (index $method.Returns 0)}}
			return nil, result.r0
			{{else}}
			return result.r0, nil
			{{end}}
		{{else}}
			{{$nonErrorReturns := nonErrorReturns $method.Returns}}
			{{if eq (len $nonErrorReturns) 1}}
			return result.r{{index $nonErrorReturns 0}}, nil
			{{else}}
			// Multiple return values, wrap them in a struct
			return struct{ {{range $i, $r := $nonErrorReturns}}
				R{{$i}} {{(index $method.Returns $r).Type}}
			{{end}} }{ {{range $i, $r := $nonErrorReturns}}
				R{{$i}}: result.r{{$r}},
			{{end}} }, nil
			{{end}}
		{{end}}
    case <-time.After(o.config.Timeout):
        return nil, fmt.Errorf("execution timeout after %v", o.config.Timeout)
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}
{{end}}
{{end}}

// ResultHandler handles execution results
type ResultHandler struct {
   values []interface{}
   err    error
   mu     sync.RWMutex
}

func (h *ResultHandler) AddResult(value interface{}, err error) {
   h.mu.Lock()
   defer h.mu.Unlock()

   if err != nil {
   	if h.err == nil {
   		h.err = err
   	} else {
   		h.err = fmt.Errorf("%v; %v", h.err, err)
   	}
   	return
   }

   h.values = append(h.values, value)
}

func (h *ResultHandler) Results() ([]interface{}, error) {
   h.mu.RLock()
   defer h.mu.RUnlock()
   return h.values, h.err
}

// Execute executes spawn options concurrently
func Execute(ctx context.Context, opts ...Option) ([]interface{}, error) {
   handler := &ResultHandler{}
   var wg sync.WaitGroup

   for _, opt := range opts {
   	wg.Add(1)
   	go func(opt Option) {
   		defer wg.Done()

   		switch o := opt.(type) {
   		{{range .Interfaces}}
   		{{range .Methods}}
   		case *{{.Name}}SpawnOption:
   			h := &{{.Name}}Handler{}
   			value, err := h.Execute(ctx, o)
   			handler.AddResult(value, err)
   		{{end}}
   		{{end}}
   		}
   	}(opt)
   }

   wg.Wait()
   return handler.Results()
}

func Param(arg interface{}) interface{} { return arg }
`

type TemplateData struct {
	Package    string
	Interfaces []Interface
	Imports    []string
	Types      []CustomType
	Features   []Feature
}

func NewTemplateData(config GeneratorConfig, interfaces []Interface) TemplateData {
	return TemplateData{
		Package:    config.Package,
		Interfaces: interfaces,
	}
}
