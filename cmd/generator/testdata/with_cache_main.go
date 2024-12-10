package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"sync"
	"time"

	"bomboclat/cmd/generator/testdata/generated"
)

type CustomErrorHandler struct {
	errors []string
	mu     sync.Mutex
}

func (h *CustomErrorHandler) HandleError(err error, stackTrace []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.errors = append(h.errors, fmt.Sprintf("Error: %v\nStack: %s", err, stackTrace))
}

func (h *CustomErrorHandler) HandlePanic(value interface{}, stackTrace []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.errors = append(h.errors, fmt.Sprintf("Panic: %v\nStack: %s", value, stackTrace))
}

type profileService struct {
}

type userService struct {
	cache generated.Cache
	pool  *generated.WorkerPool
}

func NewUserService() *userService {
	return &userService{
		cache: generated.NewInMemoryCache(),
		pool:  generated.NewWorkerPool(runtime.NumCPU()),
	}
}

func (s *userService) GetUser(ctx context.Context, id int) (*generated.User, error) {
	time.Sleep(100 * time.Millisecond)
	return &generated.User{
		ID:     id,
		Name:   fmt.Sprintf("User %d", id),
		Status: "active",
	}, nil
}

func (s *userService) GetUserWithCache(ctx context.Context, id int) (*generated.User, error) {
	if cached, ok := s.cache.Get(fmt.Sprintf("user:%d", id)); ok {
		if user, ok := cached.(*generated.User); ok {
			return user, nil
		}
	}

	user, err := s.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	s.cache.Set(fmt.Sprintf("user:%d", id), user, 5*time.Minute)
	return user, nil
}

func (s *userService) GetUsersBulk(ctx context.Context, ids []int) ([]*generated.User, error) {
	// Worker pool ile birden fazla kullanıcıyı parallel getir
	results := make(chan *generated.User, len(ids))
	errors := make(chan error, len(ids))

	for _, id := range ids {
		id := id // copy id
		s.pool.Submit(func() {
			user, err := s.GetUser(ctx, id)
			if err != nil {
				errors <- err
				return
			}
			results <- user
		})
	}

	// gather result
	var users []*generated.User
	for i := 0; i < len(ids); i++ {
		select {
		case user := <-results:
			users = append(users, user)
		case err := <-errors:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return users, nil
}

func (s *userService) UpdateProfile(ctx context.Context, user *generated.User, profile *generated.Profile) error {
	log.Printf("updateProfile has been called with parameters User:%+v Profile:%+v", user, profile)
	return nil
}

func (s *userService) NotifyUser(ctx context.Context, user *generated.User, message string) error {
	log.Printf("notifyUser has been called with parameters User:%+v Message:%s", user, message)
	return nil
}

func (p *profileService) GetProfile(ctx context.Context, userId int) (*generated.Profile, error) {
	return &generated.Profile{
		ID:      19,
		UserID:  userId,
		Address: "my-address-info",
	}, nil
}

func (p *profileService) ValidateProfile(ctx context.Context, profile *generated.Profile) (bool, error) {
	log.Printf("validateProfile has been called with parameters Profile:%+v", profile)
	return true, nil
}

func main() {
	ctx := context.Background()
	service := NewUserService()

	fmt.Println("Test 1: Single operations in parallel")
	results, err := generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserSpawn(service.GetUser, ctx, 1),
			generated.GetUserSpawn(service.GetUser, ctx, 2),
			generated.GetUserSpawn(service.GetUser, ctx, 3),
		},
	)
	if err != nil {
		log.Fatalf("Execute failed: %v", err)
	}
	for i, result := range results {
		if user, ok := result.(*generated.User); ok {
			fmt.Printf("User %d: %+v\n", i+1, user)
		}
	}

	fmt.Println("\nTest 2: Bulk operation with worker pool")
	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUsersBulkSpawn(service.GetUsersBulk, ctx, []int{4, 5, 6, 7, 8}),
		},
	)
	if err != nil {
		log.Fatalf("Execute failed: %v", err)
	}
	if users, ok := results[0].([]*generated.User); ok {
		fmt.Printf("Bulk users: %+v\n", users)
	}

	fmt.Println("\nTest 3: Mixed operations with cache")
	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserWithCacheSpawn(service.GetUserWithCache, ctx, 9), // First call - DB hit
			generated.GetUserWithCacheSpawn(service.GetUserWithCache, ctx, 9), // Second call - Cache hit
			generated.GetUsersBulkSpawn(service.GetUsersBulk, ctx, []int{10, 11, 12}),
		},
	)
	if err != nil {
		log.Fatalf("Execute failed: %v", err)
	}
	fmt.Println("First GetUserWithCache (DB):", results[0])
	fmt.Println("Second GetUserWithCache (Cache):", results[1])
	fmt.Println("Bulk result:", results[2])

	fmt.Println("\nTest 4: Timeout scenario")
	ctxTimeout, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	_, err = generated.Execute(ctxTimeout,
		[]generated.Option{
			generated.GetUsersBulkSpawn(service.GetUsersBulk, ctxTimeout, []int{13, 14, 15, 16, 17}),
		},
	)
	if err != nil {
		fmt.Printf("Expected timeout error: %v\n", err)
	}

	fmt.Println("\nTest 5: Rate Limited Operations")
	start := time.Now()

	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserSpawn(service.GetUser, ctx, 1),
			generated.GetUserSpawn(service.GetUser, ctx, 2),
			//generated.GetUserSpawn(service.GetUser, ctx, 3),
		},
		generated.WithRateLimit(2, time.Second),
	)
	if err != nil {
		log.Fatalf("Execute failed: %v", err)
	}
	fmt.Printf("Rate limited operations completed in %v\n", time.Since(start))

	fmt.Println("\nTest 6: Retry with Backoff")
	// Simüle edilmiş hata durumu için geçici bir fonksiyon
	failingService := func(ctx context.Context, id int) (*generated.User, error) {
		if time.Now().UnixNano()%2 == 0 { // Rastgele hata
			return nil, fmt.Errorf("temporary error")
		}
		return service.GetUser(ctx, id)
	}

	retryConfig := generated.RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      1 * time.Second,
		BackoffFactor: 2.0,
	}

	result, err := generated.DoWithRetry(ctx, retryConfig, func() (interface{}, error) {
		return failingService(ctx, 1)
	})
	if err != nil {
		fmt.Printf("Retry failed after all attempts: %v\n", err)
	} else {
		fmt.Printf("Retry succeeded: %+v\n", result)
	}

	fmt.Println("\nTest 7: Circuit Breaker Pattern")

	// Create a service that fails frequently
	failingService = func(ctx context.Context, id int) (*generated.User, error) {
		if time.Now().Unix()%3 == 0 { // her 3 çağrıdan 1'i başarısız
			return nil, fmt.Errorf("service error")
		}
		return service.GetUser(ctx, id)
	}

	// Circuit breaker: 3 hata sonrası aç, 5 saniye sonra tekrar dene
	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserSpawn(failingService, ctx, 1),
			generated.GetUserSpawn(failingService, ctx, 2),
			generated.GetUserSpawn(failingService, ctx, 3),
			generated.GetUserSpawn(failingService, ctx, 4),
			generated.GetUserSpawn(failingService, ctx, 5),
		},
		generated.WithCircuitBreaker(3, 5*time.Second),
	)

	if err != nil {
		fmt.Printf("Circuit breaker test results: %v\n", err)

		// Wait for reset timeout and try again
		time.Sleep(6 * time.Second)

		results, err = generated.Execute(ctx,
			[]generated.Option{
				generated.GetUserSpawn(service.GetUser, ctx, 1), // normal service ile dene
			},
			generated.WithCircuitBreaker(3, 5*time.Second),
		)

		if err != nil {
			fmt.Printf("After reset attempt failed: %v\n", err)
		} else {
			fmt.Println("Circuit closed, service recovered")
		}
	}

	fmt.Println("\nTest 8: Fallback Pattern")

	// Hata veren bir servis oluştur
	failingUserService := func(ctx context.Context, id int) (*generated.User, error) {
		return nil, fmt.Errorf("service temporarily unavailable")
	}

	// Fallback değerleri tanımla
	fallbacks := map[reflect.Type]generated.FallbackFunc{
		reflect.TypeOf(&generated.User{}): func(err error) (interface{}, error) {
			// Cache'den al veya default değer dön
			return &generated.User{
				ID:   -1,
				Name: "Default User",
			}, nil
		},
	}

	// Execute with fallback
	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserSpawn(failingUserService, ctx, 1),
		},
		generated.WithFallback(fallbacks),
	)

	if err != nil {
		fmt.Printf("Fallback test failed: %v\n", err)
	} else {
		fmt.Printf("Fallback test succeeded with default value: %+v\n", results[0])
	}

	// Generic helper kullanımı
	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserSpawn(failingUserService, ctx, 1),
		},
		generated.WithTypedFallback(&generated.User{ID: -1, Name: "Static Fallback User"}),
	)

	if err != nil {
		fmt.Printf("Typed fallback test failed: %v\n", err)
	} else {
		fmt.Printf("Typed fallback test succeeded with static value: %+v\n", results[0])
	}

	fmt.Println("\nTest 9: Custom Timeout Pattern")

	// Uzun süren bir işlem simülasyonu
	longRunningOperation := func(ctx context.Context, id int) (*generated.User, error) {
		select {
		case <-time.After(2 * time.Second):
			return &generated.User{ID: id, Name: "Late User"}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Cleanup işlemi
	cleanupFn := func() {
		fmt.Println("Performing cleanup after timeout")
		// Cleanup logic: Connections kapatma, resource temizleme vs.
	}

	// Execute with custom timeout
	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserSpawn(longRunningOperation, ctx, 1),
			generated.GetUserSpawn(service.GetUser, ctx, 2), // Normal operation
		},
		generated.WithTimeoutAndCleanup(1*time.Second, cleanupFn),
	)

	if err != nil {
		fmt.Printf("Timeout test result: %v\n", err)
	} else {
		fmt.Printf("Unexpected success in timeout test: %+v\n", results)
	}

	// Timeout ve diğer feature'ların birlikte kullanımı
	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserSpawn(longRunningOperation, ctx, 1),
		},
		generated.WithTimeout(1*time.Second),
		generated.WithTypedFallback(&generated.User{ID: -1, Name: "Timeout Fallback User"}),
	)

	if err != nil {
		fmt.Printf("Timeout with fallback failed: %v\n", err)
	} else {
		fmt.Printf("Timeout with fallback succeeded: %+v\n", results[0])
	}
	fmt.Println("\nTest 10: Metrics Collection")

	metrics := generated.NewInMemoryMetrics()

	// Execute with metrics
	results, err = generated.Execute(ctx,
		[]generated.Option{
			generated.GetUserSpawn(service.GetUser, ctx, 1),
			generated.GetUserSpawn(service.GetUserWithCache, ctx, 2),
			// Fail case
			generated.GetUserSpawn(func(ctx context.Context, id int) (*generated.User, error) {
				return nil, fmt.Errorf("simulated error")
			}, ctx, 3),
		},
		generated.WithMetrics(metrics),
		generated.WithRateLimit(10, time.Second),
		generated.WithCircuitBreaker(3, 5*time.Second),
	)

	// Print metrics summary
	fmt.Println("\nMetrics Summary:")
	metrics.Range(func(name string, metricType string, value interface{}) bool {
		fmt.Printf("Name: %s | Type: %s Value: %v\n", name, metricType, value)
		return true
	})

	fmt.Println("\nTest 11: Error Handling - Recovery")

	//errorHandler := &CustomErrorHandler{}
	//
	//// Panic yapan bir fonksiyon
	//panicFunc := func(ctx context.Context, id int) (*generated.User, error) {
	//	panic("something went wrong")
	//}
	//
	//results, err = generated.Execute(ctx,
	//	[]generated.Option{
	//		generated.GetUserSpawn(service.GetUser, ctx, 1),
	//		generated.GetUserSpawn(panicFunc, ctx, 2),
	//	},
	//	generated.WithErrorHandler(errorHandler),
	//)
	//
	//fmt.Printf("Execute completed with error: %v\n", err)
	//fmt.Println("\nCaptured Errors:")
	//for _, e := range errorHandler.errors {
	//	fmt.Println(e)
	//}
	//
	//// Normal işlem sonuçlarını kontrol et
	//for i, result := range results {
	//	if result != nil {
	//		fmt.Printf("Result %d: %+v\n", i, result)
	//	} else {
	//		fmt.Printf("Result %d: <error/panic>\n", i)
	//	}
	//}

	fmt.Println("\nTest 12: DAG Execution")

	userSvc := &userService{}
	profileSvc := &profileService{}

	flow := generated.NewFlow(
		generated.WithCache(generated.NewInMemoryCache()),
	)

	// Build a chain that:
	// 1. Gets a user
	// 2. If user is active:
	//    - Gets their profile
	//    - Validates profile
	//    - If valid, updates user with profile
	//    - Notifies user about update
	userChain := generated.Chain[*generated.User](
		generated.GetUserSpawn(userSvc.GetUser, ctx, 1),
	).When(func(user *generated.User) bool {
		return user.Status == "active"
	}).Then(
		&GetProfileOption{
			fn: profileSvc.GetProfile,
		},
		"get_profile",
	).ThenWhen(
		&ValidateProfileOption{
			fn: profileSvc.ValidateProfile,
		},
		"validate_profile",
		func(result interface{}) bool {
			profile, ok := result.(*generated.Profile)
			return ok && profile != nil
		},
	).ThenWhen(
		&UpdateProfileOption{
			fn: userSvc.UpdateProfile,
		},
		"update_profile",
		func(result interface{}) bool {
			valid, ok := result.(bool)
			return ok && valid
		},
	).Then(
		&NotifyUserOption{
			fn:      userSvc.NotifyUser,
			message: "Your profile has been updated!",
		},
		"notify_user",
	)

	// Add chain to flow
	refs := flow.Add(userChain)

	// Execute the flow
	if err := flow.Execute(ctx); err != nil {
		fmt.Printf("Flow execution failed: %v\n", err)
		return
	}

	// Process results
	for i, ref := range refs {
		result, err := flow.Get(ref)
		if err != nil {
			fmt.Printf("Operation %d failed: %v\n", i, err)
			continue
		}

		if result == nil {
			fmt.Printf("Operation %d was skipped\n", i)
			continue
		}

		fmt.Printf("Operation %d completed successfully: %+v\n", i, result)
	}
}

// Option implementations for dependent operations
type GetProfileOption struct {
	fn func(context.Context, int) (*generated.Profile, error)
}

func (o *GetProfileOption) Opt() {}

func (o *GetProfileOption) WithDependency(ctx context.Context, depID string, result interface{}) (generated.Option, error) {
	user, ok := result.(*generated.User)
	if !ok {
		return nil, fmt.Errorf("expected User, got %T", result)
	}
	return generated.GetProfileSpawn(o.fn, ctx, user.ID), nil
}

type ValidateProfileOption struct {
	fn func(context.Context, *generated.Profile) (bool, error)
}

func (o *ValidateProfileOption) Opt() {}

func (o *ValidateProfileOption) WithDependency(ctx context.Context, depID string, result interface{}) (generated.Option, error) {
	profile, ok := result.(*generated.Profile)
	if !ok {
		return nil, fmt.Errorf("expected Profile, got %T", result)
	}
	return generated.ValidateProfileSpawn(o.fn, ctx, profile), nil
}

type UpdateProfileOption struct {
	fn   func(context.Context, *generated.User, *generated.Profile) error
	user *generated.User
}

func (o *UpdateProfileOption) Opt() {}

func (o *UpdateProfileOption) WithDependency(ctx context.Context, depID string, result interface{}) (generated.Option, error) {
	// This one needs both User and Profile results
	switch depID {
	case "get_user":
		user, ok := result.(*generated.User)
		if !ok {
			return nil, fmt.Errorf("expected User, got %T", result)
		}
		// Store user and wait for profile
		return &UpdateProfileOption{fn: o.fn, user: user}, nil
	case "get_profile":
		profile, ok := result.(*generated.Profile)
		if !ok {
			return nil, fmt.Errorf("expected Profile, got %T", result)
		}
		// Now we can create the final spawn with both user and profile
		return generated.UpdateProfileSpawn(o.fn, ctx, o.user, profile), nil
	default:
		return nil, fmt.Errorf("unknown dependency ID: %s", depID)
	}
}

type NotifyUserOption struct {
	fn      func(context.Context, *generated.User, string) error
	message string
}

func (o *NotifyUserOption) Opt() {}

func (o *NotifyUserOption) WithDependency(ctx context.Context, depID string, result interface{}) (generated.Option, error) {
	user, ok := result.(*generated.User)
	if !ok {
		return nil, fmt.Errorf("expected User, got %T", result)
	}
	return generated.NotifyUserSpawn(o.fn, ctx, user, o.message), nil
}

var _ generated.DependentOption = (*ValidateProfileOption)(nil)
