package generated

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Referencing custom types from source package

type User struct {
	ID     int
	Name   string
	Status string
}

type Profile struct {
	ID      int
	UserID  int
	Address string
}

// Option is the interface all spawn options must implement
type Option interface {
	Opt()
}

// GetUserSpawnOption represents spawn configuration for UserService.GetUser
type GetUserSpawnOption struct {
	fn   func(context.Context, int) (*User, error)
	args []interface{}
}

func (o GetUserSpawnOption) Opt() {}

// GetUserSpawn creates a spawn option for UserService.GetUser
func GetUserSpawn(fn func(context.Context, int) (*User, error), args ...interface{}) Option {
	return &GetUserSpawnOption{
		fn:   fn,
		args: args,
	}
}

// GetUserHandler handles execution of UserService.GetUser
type GetUserHandler struct{}

func (h *GetUserHandler) Execute(ctx context.Context, opt Option) (interface{}, error) {

	o, ok := opt.(*GetUserSpawnOption)
	if !ok {
		return nil, fmt.Errorf("invalid option type: %T", opt)
	}

	if len(o.args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", len(o.args))
	}

	arg0, ok := o.args[0].(context.Context)
	if !ok {
		return nil, fmt.Errorf("argument 0 must be context.Context")
	}

	arg1, ok := o.args[1].(int)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be int")
	}

	// Execute with timeout
	resultCh := make(chan struct {
		r0 *User

		r1 error
	}, 1)

	go func() {

		r0, r1 := o.fn(arg0, arg1)
		resultCh <- struct {
			r0 *User

			r1 error
		}{r0, r1}
	}()

	select {
	case result := <-resultCh:

		if result.r1 != nil {
			return nil, result.r1
		}

		return result.r0, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetUserWithCacheSpawnOption represents spawn configuration for UserService.GetUserWithCache
type GetUserWithCacheSpawnOption struct {
	fn   func(context.Context, int) (*User, error)
	args []interface{}
}

func (o GetUserWithCacheSpawnOption) Opt() {}

// GetUserWithCacheSpawn creates a spawn option for UserService.GetUserWithCache
func GetUserWithCacheSpawn(fn func(context.Context, int) (*User, error), args ...interface{}) Option {
	return &GetUserWithCacheSpawnOption{
		fn:   fn,
		args: args,
	}
}

// GetUserWithCacheHandler handles execution of UserService.GetUserWithCache
type GetUserWithCacheHandler struct{}

func (h *GetUserWithCacheHandler) Execute(ctx context.Context, opt Option) (interface{}, error) {

	o, ok := opt.(*GetUserWithCacheSpawnOption)
	if !ok {
		return nil, fmt.Errorf("invalid option type: %T", opt)
	}

	if len(o.args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", len(o.args))
	}

	arg0, ok := o.args[0].(context.Context)
	if !ok {
		return nil, fmt.Errorf("argument 0 must be context.Context")
	}

	arg1, ok := o.args[1].(int)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be int")
	}

	// Execute with timeout
	resultCh := make(chan struct {
		r0 *User

		r1 error
	}, 1)

	go func() {

		r0, r1 := o.fn(arg0, arg1)
		resultCh <- struct {
			r0 *User

			r1 error
		}{r0, r1}
	}()

	select {
	case result := <-resultCh:

		if result.r1 != nil {
			return nil, result.r1
		}

		return result.r0, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetUsersBulkSpawnOption represents spawn configuration for UserService.GetUsersBulk
type GetUsersBulkSpawnOption struct {
	fn   func(context.Context, []int) ([]*User, error)
	args []interface{}
}

func (o GetUsersBulkSpawnOption) Opt() {}

// GetUsersBulkSpawn creates a spawn option for UserService.GetUsersBulk
func GetUsersBulkSpawn(fn func(context.Context, []int) ([]*User, error), args ...interface{}) Option {
	return &GetUsersBulkSpawnOption{
		fn:   fn,
		args: args,
	}
}

// GetUsersBulkHandler handles execution of UserService.GetUsersBulk
type GetUsersBulkHandler struct{}

func (h *GetUsersBulkHandler) Execute(ctx context.Context, opt Option) (interface{}, error) {

	o, ok := opt.(*GetUsersBulkSpawnOption)
	if !ok {
		return nil, fmt.Errorf("invalid option type: %T", opt)
	}

	if len(o.args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", len(o.args))
	}

	arg0, ok := o.args[0].(context.Context)
	if !ok {
		return nil, fmt.Errorf("argument 0 must be context.Context")
	}

	arg1, ok := o.args[1].([]int)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be []int")
	}

	// Execute with timeout
	resultCh := make(chan struct {
		r0 []*User

		r1 error
	}, 1)

	go func() {

		r0, r1 := o.fn(arg0, arg1)
		resultCh <- struct {
			r0 []*User

			r1 error
		}{r0, r1}
	}()

	select {
	case result := <-resultCh:

		if result.r1 != nil {
			return nil, result.r1
		}

		return result.r0, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// UpdateProfileSpawnOption represents spawn configuration for UserService.UpdateProfile
type UpdateProfileSpawnOption struct {
	fn   func(context.Context, *User, *Profile) error
	args []interface{}
}

func (o UpdateProfileSpawnOption) Opt() {}

// UpdateProfileSpawn creates a spawn option for UserService.UpdateProfile
func UpdateProfileSpawn(fn func(context.Context, *User, *Profile) error, args ...interface{}) Option {
	return &UpdateProfileSpawnOption{
		fn:   fn,
		args: args,
	}
}

// UpdateProfileHandler handles execution of UserService.UpdateProfile
type UpdateProfileHandler struct{}

func (h *UpdateProfileHandler) Execute(ctx context.Context, opt Option) (interface{}, error) {

	o, ok := opt.(*UpdateProfileSpawnOption)
	if !ok {
		return nil, fmt.Errorf("invalid option type: %T", opt)
	}

	if len(o.args) != 3 {
		return nil, fmt.Errorf("expected 3 arguments, got %d", len(o.args))
	}

	arg0, ok := o.args[0].(context.Context)
	if !ok {
		return nil, fmt.Errorf("argument 0 must be context.Context")
	}

	arg1, ok := o.args[1].(*User)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be *User")
	}

	arg2, ok := o.args[2].(*Profile)
	if !ok {
		return nil, fmt.Errorf("argument 2 must be *Profile")
	}

	// Execute with timeout
	resultCh := make(chan struct {
		r0 error
	}, 1)

	go func() {

		r0 := o.fn(arg0, arg1, arg2)
		resultCh <- struct {
			r0 error
		}{r0}
	}()

	select {
	case result := <-resultCh:

		if result.r0 != nil {
			return nil, result.r0
		}

		return nil, result.r0

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// NotifyUserSpawnOption represents spawn configuration for UserService.NotifyUser
type NotifyUserSpawnOption struct {
	fn   func(context.Context, *User, string) error
	args []interface{}
}

func (o NotifyUserSpawnOption) Opt() {}

// NotifyUserSpawn creates a spawn option for UserService.NotifyUser
func NotifyUserSpawn(fn func(context.Context, *User, string) error, args ...interface{}) Option {
	return &NotifyUserSpawnOption{
		fn:   fn,
		args: args,
	}
}

// NotifyUserHandler handles execution of UserService.NotifyUser
type NotifyUserHandler struct{}

func (h *NotifyUserHandler) Execute(ctx context.Context, opt Option) (interface{}, error) {

	o, ok := opt.(*NotifyUserSpawnOption)
	if !ok {
		return nil, fmt.Errorf("invalid option type: %T", opt)
	}

	if len(o.args) != 3 {
		return nil, fmt.Errorf("expected 3 arguments, got %d", len(o.args))
	}

	arg0, ok := o.args[0].(context.Context)
	if !ok {
		return nil, fmt.Errorf("argument 0 must be context.Context")
	}

	arg1, ok := o.args[1].(*User)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be *User")
	}

	arg2, ok := o.args[2].(string)
	if !ok {
		return nil, fmt.Errorf("argument 2 must be string")
	}

	// Execute with timeout
	resultCh := make(chan struct {
		r0 error
	}, 1)

	go func() {

		r0 := o.fn(arg0, arg1, arg2)
		resultCh <- struct {
			r0 error
		}{r0}
	}()

	select {
	case result := <-resultCh:

		if result.r0 != nil {
			return nil, result.r0
		}

		return nil, result.r0

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetProfileSpawnOption represents spawn configuration for ProfileService.GetProfile
type GetProfileSpawnOption struct {
	fn   func(context.Context, int) (*Profile, error)
	args []interface{}
}

func (o GetProfileSpawnOption) Opt() {}

// GetProfileSpawn creates a spawn option for ProfileService.GetProfile
func GetProfileSpawn(fn func(context.Context, int) (*Profile, error), args ...interface{}) Option {
	return &GetProfileSpawnOption{
		fn:   fn,
		args: args,
	}
}

// GetProfileHandler handles execution of ProfileService.GetProfile
type GetProfileHandler struct{}

func (h *GetProfileHandler) Execute(ctx context.Context, opt Option) (interface{}, error) {

	o, ok := opt.(*GetProfileSpawnOption)
	if !ok {
		return nil, fmt.Errorf("invalid option type: %T", opt)
	}

	if len(o.args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", len(o.args))
	}

	arg0, ok := o.args[0].(context.Context)
	if !ok {
		return nil, fmt.Errorf("argument 0 must be context.Context")
	}

	arg1, ok := o.args[1].(int)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be int")
	}

	// Execute with timeout
	resultCh := make(chan struct {
		r0 *Profile

		r1 error
	}, 1)

	go func() {

		r0, r1 := o.fn(arg0, arg1)
		resultCh <- struct {
			r0 *Profile

			r1 error
		}{r0, r1}
	}()

	select {
	case result := <-resultCh:

		if result.r1 != nil {
			return nil, result.r1
		}

		return result.r0, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ValidateProfileSpawnOption represents spawn configuration for ProfileService.ValidateProfile
type ValidateProfileSpawnOption struct {
	fn   func(context.Context, *Profile) (bool, error)
	args []interface{}
}

func (o ValidateProfileSpawnOption) Opt() {}

// ValidateProfileSpawn creates a spawn option for ProfileService.ValidateProfile
func ValidateProfileSpawn(fn func(context.Context, *Profile) (bool, error), args ...interface{}) Option {
	return &ValidateProfileSpawnOption{
		fn:   fn,
		args: args,
	}
}

// ValidateProfileHandler handles execution of ProfileService.ValidateProfile
type ValidateProfileHandler struct{}

func (h *ValidateProfileHandler) Execute(ctx context.Context, opt Option) (interface{}, error) {

	o, ok := opt.(*ValidateProfileSpawnOption)
	if !ok {
		return nil, fmt.Errorf("invalid option type: %T", opt)
	}

	if len(o.args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", len(o.args))
	}

	arg0, ok := o.args[0].(context.Context)
	if !ok {
		return nil, fmt.Errorf("argument 0 must be context.Context")
	}

	arg1, ok := o.args[1].(*Profile)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be *Profile")
	}

	// Execute with timeout
	resultCh := make(chan struct {
		r0 bool

		r1 error
	}, 1)

	go func() {

		r0, r1 := o.fn(arg0, arg1)
		resultCh <- struct {
			r0 bool

			r1 error
		}{r0, r1}
	}()

	select {
	case result := <-resultCh:

		if result.r1 != nil {
			return nil, result.r1
		}

		return result.r0, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

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

type executeOptions struct {
	rateLimiter *RateLimiter

	circuitBreaker *CircuitBreaker

	retry *RetryConfig

	fallbackRegistry *FallbackRegistry

	timeout TimeoutConfig

	metrics MetricsCollector
}

// Execute executes spawn options concurrently
func Execute(ctx context.Context, opts []Option, execOpts ...ExecuteOption) ([]interface{}, error) {

	execConfig := &executeOptions{}
	for _, opt := range execOpts {
		opt.applyExecute(execConfig)
	}

	handler := &ResultHandler{}
	var wg sync.WaitGroup

	for _, opt := range opts {
		wg.Add(1)
		go func(opt Option) {
			defer func() {
				if r := recover(); r != nil {

					handler.AddResult(nil, fmt.Errorf("panic recovered: %v", r))
				}
				wg.Done()
			}()

			// Circuit breaker check
			if execConfig.circuitBreaker != nil {
				if !execConfig.circuitBreaker.shouldAllow() {
					handler.AddResult(nil, fmt.Errorf("circuit breaker is open"))
					return
				}
			}

			if execConfig.rateLimiter != nil {
				if !execConfig.rateLimiter.Allow() {

					if execConfig.circuitBreaker != nil {
						execConfig.circuitBreaker.recordFailure()
					}

					handler.AddResult(nil, fmt.Errorf("rate limit exceeded"))
					return
				}
			}

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

			switch o := opt.(type) {

			case *GetUserSpawnOption:
				h := &GetUserHandler{}

				start := time.Now()
				methodName := "GetUser"

				value, err := h.Execute(ctx, o)

				if execConfig.circuitBreaker != nil {
					if err != nil {
						execConfig.circuitBreaker.recordFailure()
					} else {
						execConfig.circuitBreaker.recordSuccess()
					}
				}

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

					if execConfig.circuitBreaker != nil {
						execConfig.metrics.RecordGauge("circuit_breaker_state", float64(execConfig.circuitBreaker.GetState()), tags)
					}

					if execConfig.rateLimiter != nil {
						execConfig.metrics.IncCounter("rate_limit_exceeded_total", tags)
					}

				}

				handler.AddResult(value, err)

			case *GetUserWithCacheSpawnOption:
				h := &GetUserWithCacheHandler{}

				start := time.Now()
				methodName := "GetUserWithCache"

				value, err := h.Execute(ctx, o)

				if execConfig.circuitBreaker != nil {
					if err != nil {
						execConfig.circuitBreaker.recordFailure()
					} else {
						execConfig.circuitBreaker.recordSuccess()
					}
				}

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

					if execConfig.circuitBreaker != nil {
						execConfig.metrics.RecordGauge("circuit_breaker_state", float64(execConfig.circuitBreaker.GetState()), tags)
					}

					if execConfig.rateLimiter != nil {
						execConfig.metrics.IncCounter("rate_limit_exceeded_total", tags)
					}

				}

				handler.AddResult(value, err)

			case *GetUsersBulkSpawnOption:
				h := &GetUsersBulkHandler{}

				start := time.Now()
				methodName := "GetUsersBulk"

				value, err := h.Execute(ctx, o)

				if execConfig.circuitBreaker != nil {
					if err != nil {
						execConfig.circuitBreaker.recordFailure()
					} else {
						execConfig.circuitBreaker.recordSuccess()
					}
				}

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

					if execConfig.circuitBreaker != nil {
						execConfig.metrics.RecordGauge("circuit_breaker_state", float64(execConfig.circuitBreaker.GetState()), tags)
					}

					if execConfig.rateLimiter != nil {
						execConfig.metrics.IncCounter("rate_limit_exceeded_total", tags)
					}

				}

				handler.AddResult(value, err)

			case *UpdateProfileSpawnOption:
				h := &UpdateProfileHandler{}

				start := time.Now()
				methodName := "UpdateProfile"

				value, err := h.Execute(ctx, o)

				if execConfig.circuitBreaker != nil {
					if err != nil {
						execConfig.circuitBreaker.recordFailure()
					} else {
						execConfig.circuitBreaker.recordSuccess()
					}
				}

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

					if execConfig.circuitBreaker != nil {
						execConfig.metrics.RecordGauge("circuit_breaker_state", float64(execConfig.circuitBreaker.GetState()), tags)
					}

					if execConfig.rateLimiter != nil {
						execConfig.metrics.IncCounter("rate_limit_exceeded_total", tags)
					}

				}

				handler.AddResult(value, err)

			case *NotifyUserSpawnOption:
				h := &NotifyUserHandler{}

				start := time.Now()
				methodName := "NotifyUser"

				value, err := h.Execute(ctx, o)

				if execConfig.circuitBreaker != nil {
					if err != nil {
						execConfig.circuitBreaker.recordFailure()
					} else {
						execConfig.circuitBreaker.recordSuccess()
					}
				}

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

					if execConfig.circuitBreaker != nil {
						execConfig.metrics.RecordGauge("circuit_breaker_state", float64(execConfig.circuitBreaker.GetState()), tags)
					}

					if execConfig.rateLimiter != nil {
						execConfig.metrics.IncCounter("rate_limit_exceeded_total", tags)
					}

				}

				handler.AddResult(value, err)

			case *GetProfileSpawnOption:
				h := &GetProfileHandler{}

				start := time.Now()
				methodName := "GetProfile"

				value, err := h.Execute(ctx, o)

				if execConfig.circuitBreaker != nil {
					if err != nil {
						execConfig.circuitBreaker.recordFailure()
					} else {
						execConfig.circuitBreaker.recordSuccess()
					}
				}

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

					if execConfig.circuitBreaker != nil {
						execConfig.metrics.RecordGauge("circuit_breaker_state", float64(execConfig.circuitBreaker.GetState()), tags)
					}

					if execConfig.rateLimiter != nil {
						execConfig.metrics.IncCounter("rate_limit_exceeded_total", tags)
					}

				}

				handler.AddResult(value, err)

			case *ValidateProfileSpawnOption:
				h := &ValidateProfileHandler{}

				start := time.Now()
				methodName := "ValidateProfile"

				value, err := h.Execute(ctx, o)

				if execConfig.circuitBreaker != nil {
					if err != nil {
						execConfig.circuitBreaker.recordFailure()
					} else {
						execConfig.circuitBreaker.recordSuccess()
					}
				}

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

					if execConfig.circuitBreaker != nil {
						execConfig.metrics.RecordGauge("circuit_breaker_state", float64(execConfig.circuitBreaker.GetState()), tags)
					}

					if execConfig.rateLimiter != nil {
						execConfig.metrics.IncCounter("rate_limit_exceeded_total", tags)
					}

				}

				handler.AddResult(value, err)

			}
		}(opt)
	}

	wg.Wait()
	return handler.Results()
}

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
}

type inMemoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
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
}

type CircuitState int32

const (
	StateClosed   CircuitState = iota // Normal operation
	StateOpen                         // Failing, fast fail
	StateHalfOpen                     // Testing if service recovered
)

type CircuitBreaker struct {
	state            atomic.Int32
	failures         atomic.Int32
	failureThreshold int32
	resetTimeout     time.Duration
	lastFailure      atomic.Int64
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
	dag        *DAG
	lastNodeID int
	cache      Cache
}

func NewFlow(opts ...FlowOption) *Flow {
	f := &Flow{
		dag:        NewDAG(),
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
	nodes map[string]*DAGNode
	mu    sync.RWMutex
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
	return append([]float64{}, h.values...) // Return a copy
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

type RateLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	mu     sync.RWMutex
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
}

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxAttempts   int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

// RetryOption is the option for retry configuration
type RetryOption struct {
	Config RetryConfig
}

// WithRetry configures retry behavior
func WithRetry(maxAttempts int, initialDelay time.Duration) RetryOption {
	return RetryOption{
		Config: RetryConfig{
			MaxAttempts:   maxAttempts,
			InitialDelay:  initialDelay,
			MaxDelay:      initialDelay * 10,
			BackoffFactor: 2.0,
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
}
