package generator

import (
	"strings"
)

type FuncSignature struct {
	Name       string
	Args       []Parameter
	ReturnType []Parameter
}

type Method struct {
	Name    string
	Doc     string
	Params  []Parameter
	Returns []Parameter
}

type Parameter struct {
	Name       string
	Type       string
	IsVariadic bool
	Doc        string // parametreyi açıklayan yorum
}

type GeneratorConfig struct {
	Package    string
	InputPath  string
	OutputPath string
	Features   []string
	Interfaces []string
	Imports    []string
	Signatures []FuncSignature
	Options    []OptionType
}

type OptionType struct {
	Name   string
	Fields []Parameter
}

type Interface struct {
	Name    string
	Methods []Method
	Doc     string
}

type CustomType struct {
	Name       string
	Definition string
}

// Method için yardımcı fonksiyonlar
func (m Method) FuncType() string {
	var buf strings.Builder
	buf.WriteString("func(")

	// Parametreleri yaz
	for i, p := range m.Params {
		if i > 0 {
			buf.WriteString(", ")
		}
		if p.IsVariadic {
			buf.WriteString("...")
		}
		buf.WriteString(p.Type)
	}
	buf.WriteString(")")

	// Return tiplerini yaz
	if len(m.Returns) > 0 {
		if len(m.Returns) == 1 && m.Returns[0].Name == "" {
			buf.WriteString(" " + m.Returns[0].Type)
		} else {
			buf.WriteString(" (")
			for i, r := range m.Returns {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(r.Type)
			}
			buf.WriteString(")")
		}
	}

	return buf.String()
}

func (m Method) String() string {
	var buf strings.Builder
	buf.WriteString(m.Name + "(")

	for i, p := range m.Params {
		if i > 0 {
			buf.WriteString(", ")
		}
		if p.Name != "" {
			buf.WriteString(p.Name + " ")
		}
		if p.IsVariadic {
			buf.WriteString("...")
		}
		buf.WriteString(p.Type)
	}
	buf.WriteString(")")

	if len(m.Returns) > 0 {
		buf.WriteString(" ")
		if len(m.Returns) == 1 && m.Returns[0].Name == "" {
			buf.WriteString(m.Returns[0].Type)
		} else {
			buf.WriteString("(")
			for i, r := range m.Returns {
				if i > 0 {
					buf.WriteString(", ")
				}
				if r.Name != "" {
					buf.WriteString(r.Name + " ")
				}
				buf.WriteString(r.Type)
			}
			buf.WriteString(")")
		}
	}

	return buf.String()
}

type Feature string

const (
	FeatureBase           Feature = "base"
	FeatureWorkerPool     Feature = "worker_pool"
	FeatureCache          Feature = "cache"
	FeatureTimeout        Feature = "timeout"
	FeatureRetry          Feature = "retry"
	FeatureRateLimit      Feature = "rate_limit"
	FeatureCircuitBreaker Feature = "circuit_breaker"
	FeatureFallback       Feature = "fallback"
	FeatureMetrics        Feature = "metrics"
	FeatureErrorHandling  Feature = "error_handling"
	FeatureDAG            Feature = "dag"
)

type FeatureConfig struct {
	RequiredImports []string
	RequiredTypes   []string
	Template        string
}

var featureConfigs = map[Feature]FeatureConfig{
	FeatureBase: {
		RequiredImports: []string{
			"context",
			"fmt",
			"sync",
			"time",
		},
		RequiredTypes: nil,
		Template: `
package generated

{{if len .Imports}}
import(
	{{range .Imports}}
	"{{.}}"
	{{- end}}
)
{{end}}


// Referencing custom types from source package
{{range .Types}}
type {{.Name}} {{.Definition}}
{{end}}

// Option is the interface all spawn options must implement
type Option interface {
    Opt()
}

{{range .Interfaces}}
{{$interface := .}}
{{range .Methods}}
// {{.Name}}SpawnOption represents spawn configuration for {{$interface.Name}}.{{.Name}}
type {{.Name}}SpawnOption struct {
    fn     {{.FuncType}}
    args   []interface{}
}

func (o {{.Name}}SpawnOption) Opt() {}

// {{.Name}}Spawn creates a spawn option for {{$interface.Name}}.{{.Name}}
func {{.Name}}Spawn(fn {{.FuncType}}, args ...interface{}) Option {
    return &{{.Name}}SpawnOption{
        fn:     fn,
        args:   args,
    }
}

// {{.Name}}Handler handles execution of {{$interface.Name}}.{{.Name}}
type {{.Name}}Handler struct{}

func (h *{{.Name}}Handler) Execute(ctx context.Context, opt Option) (interface{}, error) {
	{{if hasFeature $.Features "error_handling"}}
    defer func() {
        if r := recover(); r != nil {
            if execConfig.errorHandler != nil {
                execConfig.errorHandler.HandlePanic(r, debug.Stack())
            }
            return nil, fmt.Errorf("panic recovered in {{.Name}}: %v", r)
        }
    }()
    {{end}}

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

type ExecuteOption interface {
    applyExecute(*executeOptions)
}

{{if hasAnyFeature .Features "rate_limit" "circuit_breaker" "retry" "dag"}}
type executeOptions struct {
    {{if hasFeature .Features "rate_limit"}}
    rateLimiter *RateLimiter
    {{end}}
    {{if hasFeature .Features "circuit_breaker"}}
    circuitBreaker *CircuitBreaker
    {{end}}
    {{if hasFeature .Features "retry"}}
    retry *RetryConfig
    {{end}}
	{{if hasFeature .Features "fallback"}}
	fallbackRegistry *FallbackRegistry
	{{end}}
	{{if hasFeature .Features "timeout"}}
	timeout TimeoutConfig
	{{end}}
	{{if hasFeature .Features "metrics"}}
	metrics MetricsCollector
	{{end}}
	{{if hasFeature .Features "error_handling"}}
	errorHandler ErrorHandler
	{{end}}
}
{{end}}

// Execute executes spawn options concurrently
func Execute(ctx context.Context, opts []Option, execOpts ...ExecuteOption) ([]interface{}, error) {
    {{if hasAnyFeature .Features "rate_limit" "circuit_breaker" "retry" "fallback" "timeout" "dag"}}
    execConfig := &executeOptions{}
    for _, opt := range execOpts {
        opt.applyExecute(execConfig)
    }
    {{end}}
	
	handler := &ResultHandler{}
    var wg sync.WaitGroup

    for _, opt := range opts {
        wg.Add(1)
        go func(opt Option) {
			defer func() {
				if r := recover(); r != nil {
					{{if hasFeature $.Features "error_handling"}}
					if execConfig.errorHandler != nil {
						execConfig.errorHandler.HandlePanic(r, debug.Stack())
					}
					{{end}}
					handler.AddResult(nil, fmt.Errorf("panic recovered: %v", r))
				}
				wg.Done()
			}()

			{{if hasFeature .Features "circuit_breaker"}}
            // Circuit breaker check
            if execConfig.circuitBreaker != nil {
                if !execConfig.circuitBreaker.shouldAllow() {
                    handler.AddResult(nil, fmt.Errorf("circuit breaker is open"))
                    return
                }
            }
            {{end}}

            {{if hasFeature .Features "rate_limit"}}
            if execConfig.rateLimiter != nil {
                if !execConfig.rateLimiter.Allow() {
                    {{if hasFeature .Features "circuit_breaker"}}
                    if execConfig.circuitBreaker != nil {
                        execConfig.circuitBreaker.recordFailure()
                    }
                    {{end}}
                    handler.AddResult(nil, fmt.Errorf("rate limit exceeded"))
                    return
                }
            }
            {{end}}

			{{if hasFeature .Features "timeout"}}
			if execConfig.timeout.Duration > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, execConfig.timeout.Duration)
				defer func() {
					cancel()
					if execConfig.timeout.OnTimeout != nil && ctx.Err() == context.DeadlineExceeded {
						execConfig.timeout.OnTimeout()
					}
				}()
			}
			{{end}}

            switch o := opt.(type) {
            {{range .Interfaces}}
            {{range .Methods}}
            case *{{.Name}}SpawnOption:
                h := &{{.Name}}Handler{}

				{{if hasFeature $.Features "metrics"}}
				start := time.Now()
				methodName := "{{.Name}}"
				{{end}}

                value, err := h.Execute(ctx, o)

				{{if hasFeature $.Features "circuit_breaker"}}
				if execConfig.circuitBreaker != nil {
				   if err != nil {
					   execConfig.circuitBreaker.recordFailure()
				   } else {
					   execConfig.circuitBreaker.recordSuccess()
				   }
				}
				{{end}}
				
				{{if hasFeature $.Features "fallback"}}
				if err != nil && execConfig.fallbackRegistry != nil {
				   returnType := reflect.TypeOf(o).Elem().Field(0).Type.Out(0)
				   if fallbackFn, ok := execConfig.fallbackRegistry.Get(returnType); ok {
					   fallbackValue, fallbackErr := fallbackFn(err)
					   if fallbackErr == nil {
						   value = fallbackValue
						   err = nil
					   }
				   }
				}
				{{end}}			

				{{if hasFeature $.Features "metrics"}}
				if execConfig.metrics != nil {
					duration := time.Since(start).Seconds()
					tags := map[string]string{
						"method": methodName,
						"status": fmt.Sprintf("%t", err == nil),
					}
					
					execConfig.metrics.ObserveHistogram("execution_duration_seconds", duration, tags)
					execConfig.metrics.IncCounter("execution_total", tags)
					
					if err != nil {
						execConfig.metrics.IncCounter("execution_errors_total", tags)
					}
					
					{{if hasFeature $.Features "circuit_breaker"}}
					if execConfig.circuitBreaker != nil {
						execConfig.metrics.RecordGauge("circuit_breaker_state", float64(execConfig.circuitBreaker.GetState()), tags)
					}
					{{end}}
					
					{{if hasFeature $.Features "rate_limit"}}
					if execConfig.rateLimiter != nil {
						execConfig.metrics.IncCounter("rate_limit_exceeded_total", tags)
					}
					{{end}}
				}
				{{end}}

				{{if hasFeature $.Features "error_handling"}}
				if err != nil && execConfig.errorHandler != nil {
					execConfig.errorHandler.HandleError(err, debug.Stack())
				}
				{{end}}

                handler.AddResult(value, err)
            {{end}}
            {{end}}
            }
        }(opt)
    }

    wg.Wait()
    return handler.Results()
}
`,
	},
	FeatureWorkerPool: {
		RequiredImports: []string{
			"sync",
			"context",
		},
		Template: `
type WorkerPool struct {
   size    int
   jobs    chan func()
   results chan interface{}
}

func NewWorkerPool(size int) *WorkerPool {
   pool := &WorkerPool{
       size:    size,
       jobs:    make(chan func()),
       results: make(chan interface{}),
   }
   pool.start()
   return pool
}

func (p *WorkerPool) start() {
   for i := 0; i < p.size; i++ {
       go func() {
           for job := range p.jobs {
               job()
           }
       }()
   }
}

func (p *WorkerPool) Submit(job func()) {
   p.jobs <- job
}`,
	},
	FeatureCache: {
		RequiredImports: []string{
			"time",
			"sync",
		},
		Template: `
type Cache interface {
   Get(key string) (interface{}, bool)
   Set(key string, value interface{}, ttl time.Duration)
   Delete(key string)
}

type inMemoryCache struct {
   data  map[string]cacheItem
   mu    sync.RWMutex
}

type cacheItem struct {
   value      interface{}
   expiration time.Time
}

func NewInMemoryCache() Cache {
   return &inMemoryCache{
       data: make(map[string]cacheItem),
   }
}

func (c *inMemoryCache) Get(key string) (interface{}, bool) {
   c.mu.RLock()
   defer c.mu.RUnlock()

   item, exists := c.data[key]
   if !exists || time.Now().After(item.expiration) {
       return nil, false
   }
   return item.value, true
}

func (c *inMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
   c.mu.Lock()
   defer c.mu.Unlock()

   c.data[key] = cacheItem{
       value:      value,
       expiration: time.Now().Add(ttl),
   }
}

func (c *inMemoryCache) Delete(key string) {
   c.mu.Lock()
   defer c.mu.Unlock()
   delete(c.data, key)
}`,
	},
	FeatureTimeout: {
		RequiredImports: []string{
			"context",
			"time",
		},
		Template: `
type TimeoutConfig struct {
    Duration time.Duration
    // Timeout gerçekleşince yapılacak cleanup işlemleri için
    OnTimeout func()
}

type WithTimeoutOption struct {
    config TimeoutConfig
}

func (o WithTimeoutOption) applyExecute(opts *executeOptions) {
    opts.timeout = o.config
}

func WithTimeout(duration time.Duration) ExecuteOption {
    return WithTimeoutOption{
        config: TimeoutConfig{
            Duration: duration,
        },
    }
}

func WithTimeoutAndCleanup(duration time.Duration, onTimeout func()) ExecuteOption {
    return WithTimeoutOption{
        config: TimeoutConfig{
            Duration:  duration,
            OnTimeout: onTimeout,
        },
    }
}
`,
	},
	FeatureRateLimit: {
		RequiredImports: []string{
			"context",
			"time",
			"sync",
		},
		Template: `
type RateLimiter struct {
    tokens  chan struct{}
    ticker  *time.Ticker
    mu      sync.RWMutex
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, rate),
        ticker: time.NewTicker(interval / time.Duration(rate)),
    }
    
    // Initial fill
    for i := 0; i < rate; i++ {
        rl.tokens <- struct{}{}
    }
    
    // Refill tokens
    go func() {
        for range rl.ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
            }
        }
    }()
    
    return rl
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}

type WithRateLimitOption struct {
    limiter *RateLimiter
}

func (o WithRateLimitOption) applyExecute(opts *executeOptions) {
    opts.rateLimiter = o.limiter
}

func WithRateLimit(rate int, interval time.Duration) ExecuteOption {
    return WithRateLimitOption{
        limiter: NewRateLimiter(rate, interval),
    }
}`,
	},
	FeatureRetry: {
		RequiredImports: []string{
			"context",
			"time",
			"math",
		},
		Template: `
// RetryConfig defines retry behavior
type RetryConfig struct {
    MaxAttempts      int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffFactor   float64
}

// RetryOption is the option for retry configuration
type RetryOption struct {
    Config RetryConfig
}

// WithRetry configures retry behavior
func WithRetry(maxAttempts int, initialDelay time.Duration) RetryOption {
    return RetryOption{
        Config: RetryConfig{
            MaxAttempts:    maxAttempts,
            InitialDelay:   initialDelay,
            MaxDelay:       initialDelay * 10,
            BackoffFactor:  2.0,
        },
    }
}

// WithCustomRetry allows full customization of retry behavior
func WithCustomRetry(config RetryConfig) RetryOption {
    return RetryOption{Config: config}
}

func DoWithRetry(ctx context.Context, config RetryConfig, op func() (interface{}, error)) (interface{}, error) {
    var lastErr error
    
    for attempt := 0; attempt < config.MaxAttempts; attempt++ {
        result, err := op()
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        if attempt == config.MaxAttempts-1 {
            break
        }
        
        delay := config.InitialDelay * time.Duration(math.Pow(config.BackoffFactor, float64(attempt)))
        if delay > config.MaxDelay {
            delay = config.MaxDelay
        }
        
        select {
        case <-time.After(delay):
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
    
    return nil, fmt.Errorf("all retry attempts failed: %w", lastErr)
}
`,
	},
	FeatureCircuitBreaker: {
		RequiredImports: []string{
			"time",
			"sync",
			"sync/atomic",
		},
		Template: `
type CircuitState int32

const (
    StateClosed CircuitState = iota  // Normal operation
    StateOpen                        // Failing, fast fail
    StateHalfOpen                    // Testing if service recovered
)

type CircuitBreaker struct {
    state           atomic.Int32
    failures        atomic.Int32
    failureThreshold int32
    resetTimeout    time.Duration
    lastFailure     atomic.Int64
}

func NewCircuitBreaker(failureThreshold int32, resetTimeout time.Duration) *CircuitBreaker {
    cb := &CircuitBreaker{
        failureThreshold: failureThreshold,
        resetTimeout:     resetTimeout,
    }
    cb.state.Store(int32(StateClosed))
    return cb
}

func (cb *CircuitBreaker) GetState() CircuitState {
    return CircuitState(cb.state.Load())
}

func (cb *CircuitBreaker) recordSuccess() {
    cb.failures.Store(0)
    cb.state.Store(int32(StateClosed))
}

func (cb *CircuitBreaker) recordFailure() {
    failures := cb.failures.Add(1)
    cb.lastFailure.Store(time.Now().UnixNano())
    
    if failures >= cb.failureThreshold {
        cb.state.Store(int32(StateOpen))
    }
}

func (cb *CircuitBreaker) shouldAllow() bool {
    switch CircuitState(cb.state.Load()) {
    case StateClosed:
        return true
    case StateOpen:
        lastFailure := time.Unix(0, cb.lastFailure.Load())
        if time.Since(lastFailure) > cb.resetTimeout {
            cb.state.Store(int32(StateHalfOpen))
            return true
        }
        return false
    case StateHalfOpen:
        return true
    default:
        return false
    }
}

type WithCircuitBreakerOption struct {
    breaker *CircuitBreaker
}

func (o WithCircuitBreakerOption) applyExecute(opts *executeOptions) {
    opts.circuitBreaker = o.breaker
}

func WithCircuitBreaker(failureThreshold int32, resetTimeout time.Duration) ExecuteOption {
    return WithCircuitBreakerOption{
        breaker: NewCircuitBreaker(failureThreshold, resetTimeout),
    }
}
`,
	},
	FeatureFallback: {
		RequiredImports: []string{
			"context",
			"fmt",
			"reflect",
		},
		Template: `
// FallbackFunc represents a function that provides fallback value
type FallbackFunc func(error) (interface{}, error)

// FallbackRegistry holds fallback functions for different types
type FallbackRegistry struct {
    fallbacks map[reflect.Type]FallbackFunc
    mu        sync.RWMutex
}

func NewFallbackRegistry() *FallbackRegistry {
    return &FallbackRegistry{
        fallbacks: make(map[reflect.Type]FallbackFunc),
    }
}

func (r *FallbackRegistry) Register(typ reflect.Type, fn FallbackFunc) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.fallbacks[typ] = fn
}

func (r *FallbackRegistry) Get(typ reflect.Type) (FallbackFunc, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    fn, ok := r.fallbacks[typ]
    return fn, ok
}

type WithFallbackOption struct {
    registry *FallbackRegistry
}

func (o WithFallbackOption) applyExecute(opts *executeOptions) {
    opts.fallbackRegistry = o.registry
}

func WithFallback(fallbacks map[reflect.Type]FallbackFunc) ExecuteOption {
    registry := NewFallbackRegistry()
    for typ, fn := range fallbacks {
        registry.Register(typ, fn)
    }
    return WithFallbackOption{registry: registry}
}

// Helper functions for common types
func WithTypedFallback[T any](fallback T) ExecuteOption {
    registry := NewFallbackRegistry()
    registry.Register(reflect.TypeOf((*T)(nil)).Elem(), func(err error) (interface{}, error) {
        return fallback, nil
    })
    return WithFallbackOption{registry: registry}
}
`,
	},
	FeatureMetrics: {
		RequiredImports: []string{
			"time",
			"sync",
			"sync/atomic",
			"sort",
			"strings",
		},
		Template: `

type histogramValues struct {
    values []float64
    mu     sync.RWMutex
}

func newHistogramValues() *histogramValues {
    return &histogramValues{
        values: make([]float64, 0),
    }
}

func (h *histogramValues) Append(value float64) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.values = append(h.values, value)
}

func (h *histogramValues) Values() []float64 {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return append([]float64{}, h.values...)  // Return a copy
}

// Metrics interface for different metric collectors
type MetricsCollector interface {
    IncCounter(name string, tags map[string]string)
    ObserveHistogram(name string, value float64, tags map[string]string)
    RecordGauge(name string, value float64, tags map[string]string)
	Range(f func(name string, metricType string, value interface{}) bool)
}

// Simple in-memory metrics implementation
type InMemoryMetrics struct {
    counters   sync.Map // string -> *atomic.Int64
    histograms sync.Map // string -> *histogramValues
    gauges     sync.Map // string -> float64
}

func NewInMemoryMetrics() *InMemoryMetrics {
    return &InMemoryMetrics{}
}

func (m *InMemoryMetrics) IncCounter(name string, tags map[string]string) {
    key := formatMetricKey(name, tags)
    value, _ := m.counters.LoadOrStore(key, &atomic.Int64{})
    value.(*atomic.Int64).Add(1)
}

func (m *InMemoryMetrics) ObserveHistogram(name string, value float64, tags map[string]string) {
    key := formatMetricKey(name, tags)
    hist, _ := m.histograms.LoadOrStore(key, newHistogramValues())
    hist.(*histogramValues).Append(value)
}

func (m *InMemoryMetrics) RecordGauge(name string, value float64, tags map[string]string) {
    key := formatMetricKey(name, tags)
    m.gauges.Store(key, value)
}

func (m *InMemoryMetrics) Range(f func(name string, metricType string, value interface{}) bool) {
    m.counters.Range(func(key, value interface{}) bool {
        return f(key.(string), "counter", value.(*atomic.Int64).Load())
    })
    
    m.histograms.Range(func(key, value interface{}) bool {
        return f(key.(string), "histogram", value.(*histogramValues).Values())
    })
    
    m.gauges.Range(func(key, value interface{}) bool {
        return f(key.(string), "gauge", value)
    })
}

func formatMetricKey(name string, tags map[string]string) string {
    if len(tags) == 0 {
        return name
    }
    // Sort tags for consistent key generation
    var pairs []string
    for k, v := range tags {
        pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
    }
    sort.Strings(pairs)
    return fmt.Sprintf("%s{%s}", name, strings.Join(pairs, ","))
}

type WithMetricsOption struct {
    collector MetricsCollector
}

func (o WithMetricsOption) applyExecute(opts *executeOptions) {
    opts.metrics = o.collector
}

func WithMetrics(collector MetricsCollector) ExecuteOption {
    return WithMetricsOption{collector: collector}
}
`,
	},
	FeatureErrorHandling: {
		RequiredImports: []string{
			"runtime/debug",
			"fmt",
		},
		Template: `
// ErrorHandler defines how errors should be processed
type ErrorHandler interface {
    HandleError(err error, stackTrace []byte)
    HandlePanic(value interface{}, stackTrace []byte)
}

// DefaultErrorHandler provides basic error handling
type DefaultErrorHandler struct{}

func (h *DefaultErrorHandler) HandleError(err error, stackTrace []byte) {
    // Just pass through the error, let the caller handle it
}

func (h *DefaultErrorHandler) HandlePanic(value interface{}, stackTrace []byte) {
    // Convert panic to error
    err := fmt.Errorf("panic recovered: %v\n%s", value, stackTrace)
    h.HandleError(err, stackTrace)
}

type WithErrorHandlerOption struct {
    handler ErrorHandler
}

func (o WithErrorHandlerOption) applyExecute(opts *executeOptions) {
    opts.errorHandler = o.handler
}

func WithErrorHandler(handler ErrorHandler) ExecuteOption {
    return WithErrorHandlerOption{handler: handler}
}`,
	},
	FeatureDAG: {
		RequiredImports: []string{
			"context",
			"fmt",
			"sync",
			"time",
		},
		Template: `
type NodeRef struct {
    id string
}

type StepBuilderBase interface {
    // GetFirstOption returns the first option in the chain
    GetFirstOption() Option
    // GetDependentOperations returns all subsequent operations
    GetDependentOperations() []DependentOperation
    // GetCondition returns the condition function if any
    GetCondition() func(interface{}) bool
}

func (s *StepBuilder[T]) GetFirstOption() Option {
    return s.first
}

func (s *StepBuilder[T]) GetDependentOperations() []DependentOperation {
    return s.then
}

func (s *StepBuilder[T]) GetCondition() func(interface{}) bool {
    if s.condition == nil {
        return nil
    }
    return func(i interface{}) bool {
        if t, ok := i.(T); ok {
            return s.condition(t)
        }
        return true
    }
}

type DependentOperation struct {
    operation DependentOption
    name      string
    condition func(interface{}) bool 
}

type DependentOption interface {
    Option
    WithDependency(ctx context.Context, dependencyID string, result interface{}) (Option, error)
}

type Flow struct {
    dag       *DAG
    lastNodeID int
    cache     Cache              
}

func NewFlow(opts ...FlowOption) *Flow {
    f := &Flow{
        dag:       NewDAG(),
        lastNodeID: 0,
    }
    for _, opt := range opts {
        opt(f)
    }
    return f
}

type FlowOption func(*Flow)

func WithCache(cache Cache) FlowOption {
    return func(f *Flow) {
        f.cache = cache
    }
}

func Chain[T any](first Option) *StepBuilder[T] {
    return &StepBuilder[T]{
        first: first,
    }
}

type StepBuilder[T any] struct {
    first     Option
    then      []DependentOperation
    condition func(T) bool
}

func (s *StepBuilder[T]) Then(operation DependentOption, name string) *StepBuilder[T] {
    s.then = append(s.then, DependentOperation{
        operation: operation,
        name:      name,
    })
    return s
}

func (s *StepBuilder[T]) When(condition func(T) bool) *StepBuilder[T] {
    s.condition = condition
    return s
}

func (s *StepBuilder[T]) ThenWhen(operation DependentOption, name string, condition func(interface{}) bool) *StepBuilder[T] {
    s.then = append(s.then, DependentOperation{
        operation: operation,
        name:      name,
        condition: condition,
    })
    return s
}

type DAG struct {
    nodes    map[string]*DAGNode
    mu       sync.RWMutex
}

func NewDAG() *DAG {
   return &DAG{
       nodes: make(map[string]*DAGNode),
   }
}

type DAGNode struct {
    ID           string
    Option       Option
    Dependencies []string
    Result       interface{}
    Error        error
    condition    func(interface{}) bool
    executed     bool
    skipped      bool
	skipReason   string
}

func AddToFlow[T any](f *Flow, step *StepBuilder[T]) []NodeRef {
    // First operation
    firstRef := f.Spawn(step.first)
    if step.condition != nil {
        f.dag.nodes[firstRef.id].condition = func(i interface{}) bool {
            if t, ok := i.(T); ok {
                return step.condition(t)
            }
            return true
        }
    }
    
    // All subsequent operations
    var refs []NodeRef
    for _, op := range step.then {
        ref := f.Then(op.operation, firstRef)
        if op.condition != nil {
            f.dag.nodes[ref.id].condition = op.condition
        }
        refs = append(refs, ref)
    }
    
    return refs
}

func (f *Flow) Add(step StepBuilderBase) []NodeRef {
    // First operation
    firstRef := f.Spawn(step.GetFirstOption())
    
    // Set condition if exists
    if condition := step.GetCondition(); condition != nil {
        f.dag.nodes[firstRef.id].condition = condition
    }
    
    // All subsequent operations
    var refs []NodeRef
    for _, op := range step.GetDependentOperations() {
        ref := f.Then(op.operation, firstRef)
        if op.condition != nil {
            f.dag.nodes[ref.id].condition = op.condition
        }
        refs = append(refs, ref)
    }
    
    return refs
}

func (f *Flow) Then(opt DependentOption, dependency NodeRef) NodeRef {
    id := fmt.Sprintf("node_%d", f.lastNodeID)
    f.lastNodeID++
    
    f.dag.mu.Lock()
    f.dag.nodes[id] = &DAGNode{
        ID:           id,
        Option:       opt,
        Dependencies: []string{dependency.id},
    }
    f.dag.mu.Unlock()
    
    return NodeRef{id: id}
}

func (f *Flow) Spawn(opt Option) NodeRef {
    id := fmt.Sprintf("node_%d", f.lastNodeID)
    f.lastNodeID++
    
    f.dag.mu.Lock()
    f.dag.nodes[id] = &DAGNode{
        ID:     id,
        Option: opt,
    }
    f.dag.mu.Unlock()
    
    return NodeRef{id: id}
}

// Get retrieves the result for a specific node reference
func (f *Flow) Get(ref NodeRef) (interface{}, error) {
    f.dag.mu.RLock()
    defer f.dag.mu.RUnlock()

    node, exists := f.dag.nodes[ref.id]
    if !exists {
        return nil, fmt.Errorf("node %s not found", ref.id)
    }

    // Node skip edilmiş mi kontrol et
    if node.skipped {
        return nil, fmt.Errorf("node was skipped: %s", node.skipReason)
    }

    // Hata varsa döndür
    if node.Error != nil {
        return nil, node.Error
    }

    return node.Result, nil
}

// Generic type-safe getter for specific types
func GetTyped[T any](f *Flow, ref NodeRef) (T, error) {
    result, err := f.Get(ref)
    if err != nil {
        return *new(T), err
    }

    if result == nil {
        return *new(T), nil
    }

    typed, ok := result.(T)
    if !ok {
        return *new(T), fmt.Errorf("invalid result type: expected %T, got %T", *new(T), result)
    }

    return typed, nil
}

func (f *Flow) Execute(ctx context.Context, execOpts ...ExecuteOption) error {
    var executed = make(map[string]bool)
    
    for len(executed) < len(f.dag.nodes) {
        readyNodes := make([]*DAGNode, 0)
        for id, node := range f.dag.nodes {
            // İlk node için debug log ekleyelim
            fmt.Printf("Checking node %s: executed=%v, skipped=%v\n", id, executed[id], node.skipped)
            
            if executed[id] {
                continue
            }

            if node.skipped {
                executed[id] = true
                continue
            }

            // İlk node için dependency olmamalı
            canExecute := true
            if len(node.Dependencies) > 0 {
                fmt.Printf("Node %s has dependencies: %v\n", id, node.Dependencies)
                for _, depID := range node.Dependencies {
                    depNode := f.dag.nodes[depID]
                    if !executed[depID] {
                        canExecute = false
                        fmt.Printf("Dependency %s not yet executed\n", depID)
                        break
                    }
                    
                    if depNode.skipped {
                        node.skipped = true
                        node.skipReason = fmt.Sprintf("dependency %s was skipped", depID)
                        canExecute = false
                        fmt.Printf("Dependency %s was skipped\n", depID)
                        break
                    }
                }
            }

            if canExecute {
                fmt.Printf("Adding node %s to ready nodes\n", id)
                readyNodes = append(readyNodes, node)
            }
        }

        // Execute ready nodes with more detailed logging
        for _, node := range readyNodes {
            fmt.Printf("Executing node %s\n", node.ID)
			// Check cache first
			if f.cache != nil {
				if cachedResult, ok := f.cache.Get(node.ID); ok {
					node.Result = cachedResult
					executed[node.ID] = true
					continue
				}
			}
            
            opt := node.Option
			if len(node.Dependencies) > 0 && opt != nil {
				if depOpt, ok := opt.(DependentOption); ok {
					for _, depID := range node.Dependencies {
						depNode := f.dag.nodes[depID]
						var err error
						opt, err = depOpt.WithDependency(ctx, depID, depNode.Result)
						if err != nil {
							return fmt.Errorf("failed to process dependency %s: %w", depID, err)
						}
					}
				}
			}

            results, err := Execute(ctx, []Option{opt}, execOpts...)
            if err != nil {
                node.Error = err
                return fmt.Errorf("node %s failed: %w", node.ID, err)
            }

            if len(results) > 0 {
                node.Result = results[0]
                fmt.Printf("Node %s executed successfully with result: %+v\n", node.ID, node.Result)
                
                // Condition kontrolünü doğru sırayla yapalım
                if node.Result != nil && node.condition != nil {
                    conditionResult := node.condition(node.Result)
                    fmt.Printf("Node %s condition result: %v\n", node.ID, conditionResult)
                    if !conditionResult {
                        node.skipped = true
                        node.skipReason = "condition returned false after execution"
                    }
                }
            } else {
                fmt.Printf("Node %s executed but returned no results\n", node.ID)
            }
            
            executed[node.ID] = true
        }

        // Her iterasyonda durum özeti
        fmt.Printf("\nExecution state after iteration:\n")
        for id, node := range f.dag.nodes {
            fmt.Printf("Node %s: executed=%v, skipped=%v, reason=%s\n", 
                id, executed[id], node.skipped, node.skipReason)
        }
        fmt.Println()
    }

    return nil
}

`,
	},
}

var featureDependencies = map[Feature][]Feature{
	FeatureBase:           {},
	FeatureWorkerPool:     {},
	FeatureCache:          {},
	FeatureRateLimit:      {FeatureBase},
	FeatureCircuitBreaker: {FeatureBase},
	FeatureRetry:          {FeatureBase},
	FeatureFallback:       {FeatureBase},
	FeatureMetrics:        {FeatureBase},
	FeatureErrorHandling:  {FeatureBase},
	FeatureDAG:            {FeatureBase},
}
