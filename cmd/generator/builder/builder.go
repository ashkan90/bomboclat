package builder

import (
	"bomboclat/cmd/generator"
	"sync"
)

type Builder struct {
	pkg        string
	imports    []string
	templates  []Template
	generators []Generator
}

type Template interface {
	Name() string
	Content() string
	Dependencies() []string
}

type Generator interface {
	Generate(pkg string, signatures []generator.FuncSignature) ([]byte, error)
}

type TemplateRegistry struct {
	templates map[string]Template
	mu        sync.RWMutex
}

func (r *TemplateRegistry) Register(t Template) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[t.Name()] = t
}

type TestGenerator struct {
	templateRegistry *TemplateRegistry
}

const testTemplate = `
{{range .Signatures}}
func Test{{.Name}}Spawn(t *testing.T) {
    tests := []struct {
        name    string
        fn      {{.FuncType}}
        args    []interface{}
        want    {{.ReturnType}}
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            spawn := {{.Name}}Spawn(tt.fn, tt.args...)
            ctx := context.Background()
            
            executor := NewExecutor()
            result := executor.Execute(ctx, spawn)
            
            if tt.wantErr {
                assert.Error(t, result.Err)
                return
            }
            
            assert.NoError(t, result.Err)
            got := result.Get{{.Name}}(0)
            assert.Equal(t, tt.want, got)
        })
    }
}

func TestConcurrent{{.Name}}Spawn(t *testing.T) {
    var wg sync.WaitGroup
    workers := 10
    iterations := 100
    
    executor := NewExecutor()
    
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            for j := 0; j < iterations; j++ {
                ctx := context.Background()
                
                // Create mock function and args
                fn := func(ctx context.Context, id int) ({{.ReturnType}}, error) {
                    return &{{.Name}}Model{ID: id}, nil
                }
                
                spawn := {{.Name}}Spawn(fn, ctx, j)
                result := executor.Execute(ctx, spawn)
                
                assert.NoError(t, result.Err)
                got := result.Get{{.Name}}(0)
                assert.Equal(t, j, got.ID)
            }
        })
    }
    
    wg.Wait()
}
{{end}}

func TestDAGExecution(t *testing.T) {
    executor := NewExecutor()
    dag := NewDAG(executor)
    
    // Add test nodes
    err := dag.AddNode("root", mockSpawn())
    assert.NoError(t, err)
    
    err = dag.AddNode("child1", mockSpawn(), "root")
    assert.NoError(t, err)
    
    err = dag.AddNode("child2", mockSpawn(), "root")
    assert.NoError(t, err)
    
    // Execute
    ctx := context.Background()
    result := dag.Execute(ctx)
    
    assert.NoError(t, result.Err)
    assert.Len(t, result.Values, 3)
}

func TestRateLimiting(t *testing.T) {
    metrics := newTestMetrics(t)
    limiter := NewTokenBucketLimiter(2, time.Millisecond*100, metrics)
    
    ctx := context.Background()
    
    // Should allow 2 immediate acquisitions
    assert.NoError(t, limiter.Acquire(ctx))
    assert.NoError(t, limiter.Acquire(ctx))
    
    // Third should block
    ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Millisecond*50)
    defer cancel()
    
    err := limiter.Acquire(ctxWithTimeout)
    assert.Error(t, err)
    assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestBackpressure(t *testing.T) {
    metrics := newTestMetrics(t)
    limiter := NewTokenBucketLimiter(1, time.Millisecond*100, metrics)
    
    backoff := &ExponentialBackoff{
        initial: time.Millisecond * 10,
        max:     time.Millisecond * 100,
        factor:  2,
    }
    
    handler := NewBackpressureHandler(
        mockHandler(t),
        limiter,
        3,
        backoff,
    )
    
    ctx := context.Background()
    
    // Simulate load
    for i := 0; i < 10; i++ {
        go func() {
            _, err := handler.Execute(ctx, mockSpawn())
            assert.NoError(t, err)
        }()
    }
}
`
