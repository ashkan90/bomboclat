package generator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var workingDir string

func init() {
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func TestGenerator(t *testing.T) {
	testDir := filepath.Join(workingDir, "testdata")
	interfacesDir := filepath.Join(testDir, "interfaces")
	generatedDir := filepath.Join(testDir, "generated")

	// Test setup - testdata dizinlerini oluştur
	err := os.MkdirAll(interfacesDir, 0755)
	require.NoError(t, err)
	//defer os.RemoveAll(testDir)

	// Test interface dosyasını oluştur
	interfaceContent := `
package interfaces

import "context"

type UserService interface {
   GetUser(ctx context.Context, id int) (*User, error)
   GetProfile(ctx context.Context, id int) (*Profile, error)
   UpdateUser(ctx context.Context, user *User) error
}

type User struct {
   ID   int
   Name string
}

type Profile struct {
   ID   int
   Name string
   Surname string
   Address string
}
`
	err = os.WriteFile("./testdata/interfaces/service.go", []byte(interfaceContent), 0644)
	require.NoError(t, err)
	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		errMessage string
		validate   func(t *testing.T, outputPath string)
	}{
		{
			name: "Generate UserService",
			args: []string{
				"-input", interfacesDir,
				"-output", filepath.Join(generatedDir, "user_service.go"),
				"-package", "generated",
				"-interfaces", "UserService",
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				// Read generated file
				content, err := os.ReadFile(outputPath)
				require.NoError(t, err)

				// Check package name
				assert.Contains(t, string(content), "package generated")

				// Check if method is generated
				assert.Contains(t, string(content), "GetUserSpawn")

				// Check spawn option type
				assert.Contains(t, string(content), "type GetUserSpawnOption struct")

				// Check interface implementation
				assert.Contains(t, string(content), "func (o GetUserSpawnOption) opt()")

				// Check handler implementation
				assert.Contains(t, string(content), "type GetUserHandler struct")
			},
		},
		{
			name: "Generate UserService",
			args: []string{
				"-input", interfacesDir,
				"-output", filepath.Join(generatedDir, "user_service.go"),
				"-package", "generated",
				"-interfaces", "UserService",
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				// Read generated file
				content, err := os.ReadFile(outputPath)
				require.NoError(t, err)

				// Check package name
				assert.Contains(t, string(content), "package generated")

				// Check if all methods are generated
				assert.Contains(t, string(content), "GetUserSpawn")
				assert.Contains(t, string(content), "GetProfileSpawn")
				assert.Contains(t, string(content), "UpdateUserSpawn")

				// Check spawn option types
				assert.Contains(t, string(content), "type GetUserSpawnOption struct")
				assert.Contains(t, string(content), "type GetProfileSpawnOption struct")
				assert.Contains(t, string(content), "type UpdateUserSpawnOption struct")

				// Check interface implementation
				assert.Contains(t, string(content), "func (o GetUserSpawnOption) opt()")

				// Check handler implementation
				assert.Contains(t, string(content), "type GetUserHandler struct")
				assert.Contains(t, string(content), "func (h *GetUserHandler) Execute(ctx context.Context, opt Option)")
			},
		},
		{
			name: "Generate multiple interfaces",
			args: []string{
				"-input", interfacesDir,
				"-output", filepath.Join(generatedDir, "user_service.go"),
				"-package", "generated",
				"-interfaces", "UserService,OrderService",
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				content, err := os.ReadFile(outputPath)
				require.NoError(t, err)

				// Check both services
				assert.Contains(t, string(content), "GetUserSpawn")
				assert.Contains(t, string(content), "CreateOrderSpawn")
			},
		},
		{
			name: "Missing input path",
			args: []string{
				"-output", filepath.Join(generatedDir, "user_service.go"),
				"-package", "generated",
				"-interfaces", "UserService",
			},
			wantErr:    true,
			errMessage: "input path is required",
		},
		{
			name: "Missing output path",
			args: []string{
				"-input", interfacesDir,
				"-package", "generated",
				"-interfaces", "UserService",
			},
			wantErr:    true,
			errMessage: "output path is required",
		},
		{
			name: "Invalid interface name",
			args: []string{
				"-input", interfacesDir,
				"-output", filepath.Join(generatedDir, "user_service.go"),
				"-package", "generated",
				"-interfaces", "NonExistentService",
			},
			wantErr:    true,
			errMessage: "interface NonExistentService not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create output directory if needed
			if !tt.wantErr {
				outDir := filepath.Dir(getOutputPath(tt.args))
				err := os.MkdirAll(outDir, 0755)
				require.NoError(t, err)
			}

			// Save original args and restore them after test
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = append([]string{"cmd"}, tt.args...)

			// Run generator
			err = run()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
				return
			}

			assert.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, getOutputPath(tt.args))
			}
		})
	}
}

func TestGeneratorWithFeatures(t *testing.T) {
	testDir := filepath.Join(workingDir, "testdata")
	interfacesDir := filepath.Join(testDir, "interfaces")
	generatedDir := filepath.Join(testDir, "generated")

	// Test setup
	err := os.MkdirAll(interfacesDir, 0755)
	require.NoError(t, err)
	//defer os.RemoveAll(testDir)

	// Test interface dosyasını oluştur
	interfaceContent := `
package interfaces

import "context"

type UserService interface {
   GetUser(ctx context.Context, id int) (*User, error)
   GetUserWithCache(ctx context.Context, id int) (*User, error)
   GetUsersBulk(ctx context.Context, ids []int) ([]*User, error)
	UpdateProfile(ctx context.Context, user *User, profile *Profile) error
	NotifyUser(ctx context.Context, user *User, message string) error
}

type ProfileService interface {
    GetProfile(ctx context.Context, userId int) (*Profile, error)
    ValidateProfile(ctx context.Context, profile *Profile) (bool, error)
}

type User struct {
   ID   int
   Name string
	Status string
}

type Profile struct {
	ID int
	UserID int
	Address string
}
`
	err = os.WriteFile(filepath.Join(interfacesDir, "service.go"), []byte(interfaceContent), 0644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		args     []string
		features []string
		validate func(t *testing.T, content string)
	}{
		{
			name:     "Generate with Cache Feature",
			features: []string{"worker_pool", "cache", "retry", "rate_limit", "circuit_breaker", "fallback", "timeout", "metrics", "dag"},
			validate: func(t *testing.T, content string) {
				// Cache yapılarını kontrol et
				assert.Contains(t, content, "type Cache interface")
				assert.Contains(t, content, "type InMemoryCache struct")
				assert.Contains(t, content, "func NewInMemoryCache()")

				// Cache middleware kontrolü
				assert.Contains(t, content, "type CacheMiddleware struct")
				assert.Contains(t, content, "func NewCacheMiddleware")
				assert.Contains(t, content, "type cacheHandler struct")
				assert.Contains(t, content, "func (h *cacheHandler) Execute")

				// Base functionality de olmalı
				assert.Contains(t, content, "type GetUserSpawnOption struct")
			},
		},
		{
			name:     "Generate with DAG",
			features: []string{"dag", "timeout", "worker_pool"},
			validate: func(t *testing.T, content string) {
				// Cache yapılarını kontrol et
				assert.Contains(t, content, "type Cache interface")
				assert.Contains(t, content, "type InMemoryCache struct")
				assert.Contains(t, content, "func NewInMemoryCache()")

				// Cache middleware kontrolü
				assert.Contains(t, content, "type CacheMiddleware struct")
				assert.Contains(t, content, "func NewCacheMiddleware")
				assert.Contains(t, content, "type cacheHandler struct")
				assert.Contains(t, content, "func (h *cacheHandler) Execute")

				// Base functionality de olmalı
				assert.Contains(t, content, "type GetUserSpawnOption struct")
			},
		},
		{
			name:     "Generate with Bulkhead Feature",
			features: []string{"base", "bulkhead"},
			validate: func(t *testing.T, content string) {
				// Bulkhead yapılarını kontrol et
				assert.Contains(t, content, "type BulkheadConfig struct")
				assert.Contains(t, content, "type Bulkhead struct")
				assert.Contains(t, content, "func NewBulkhead(")

				// Bulkhead middleware kontrolü
				assert.Contains(t, content, "WithBulkhead(maxCalls, queueCapacity int)")

				// Base functionality de olmalı
				assert.Contains(t, content, "type GetUsersBulkSpawnOption struct")
			},
		},
		{
			name:     "Generate with Multiple Features",
			features: []string{"base", "cache", "bulkhead", "ratelimit"},
			validate: func(t *testing.T, content string) {
				// Cache kontrolü
				assert.Contains(t, content, "type Cache interface")

				// Bulkhead kontrolü
				assert.Contains(t, content, "type BulkheadConfig struct")

				// Rate limiter kontrolü
				assert.Contains(t, content, "type RateLimiter struct")
				assert.Contains(t, content, "func NewRateLimiter(")

				// Options kontrolü
				assert.Contains(t, content, "WithCache(")
				assert.Contains(t, content, "WithBulkhead(")
				assert.Contains(t, content, "WithRateLimit(")
			},
		},
		{
			name:     "Generate with Events Feature",
			features: []string{"base", "events"},
			validate: func(t *testing.T, content string) {
				// Event yapılarını kontrol et
				assert.Contains(t, content, "type EventType string")
				assert.Contains(t, content, "type Event struct")
				assert.Contains(t, content, "type Observer interface")

				// Event handling kontrolü
				assert.Contains(t, content, "func (e *EventEmitter) Emit(")
				assert.Contains(t, content, "func (e *EventEmitter) Register(")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := filepath.Join(generatedDir, fmt.Sprintf("%s.go", strings.ReplaceAll(strings.ToLower(tt.name), " ", "_")))

			// Create generator config
			config := Config{
				InputPath:  interfacesDir,
				OutputPath: outputPath,
				Package:    "generated",
				Interfaces: []string{"UserService", "ProfileService"},
				Features:   tt.features,
			}

			// Generate code
			genConfig := GeneratorConfig{
				Package:    config.Package,
				InputPath:  config.InputPath,
				OutputPath: config.OutputPath,
				Interfaces: config.Interfaces,
				Features:   config.Features,
			}

			gen, err := NewGenerator(genConfig)
			require.NoError(t, err)

			err = gen.Generate()
			require.NoError(t, err)

			// Read generated content
			content, err := os.ReadFile(outputPath)
			require.NoError(t, err)

			// Run validation
			tt.validate(t, string(content))

			// Try to compile generated code
			//cmd := exec.Command("go", "build", outputPath)
			//cmd.Dir = generatedDir
			//output, err := cmd.CombinedOutput()
			//if err != nil {
			//	t.Fatalf("Generated code does not compile: %v\nOutput: %s", err, output)
			//}
		})
	}
}

func getOutputPath(args []string) string {
	for i, arg := range args {
		if arg == "-output" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

// Helper function to run the generator
func run() error {
	main()
	return nil
}

// Test directory setup and cleanup
func TestMain(m *testing.M) {
	// Setup
	err := os.MkdirAll("./testdata/generated", 0755)
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	// os.RemoveAll("./testdata/generated")

	os.Exit(code)
}
