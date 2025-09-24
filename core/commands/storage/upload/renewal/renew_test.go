package renewal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenewRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request *RenewRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &RenewRequest{
				CID:         "QmTestHash",
				Duration:    30,
				Token:       common.HexToAddress("0x1234567890123456789012345678901234567890"),
				Price:       1000,
				RenterID:    "12D3KooWTest",
				OriginalEnd: time.Now().Add(24 * time.Hour),
				NewEnd:      time.Now().Add(30*24*time.Hour + 24*time.Hour),
				TotalCost:   30000,
			},
			wantErr: false,
		},
		{
			name: "invalid duration",
			request: &RenewRequest{
				CID:      "QmTestHash",
				Duration: 0, // Invalid: zero duration
				Token:    common.HexToAddress("0x1234567890123456789012345678901234567890"),
				Price:    1000,
			},
			wantErr: true,
		},
		{
			name: "empty file hash",
			request: &RenewRequest{
				CID:      "", // Invalid: empty file hash
				Duration: 30,
				Token:    common.HexToAddress("0x1234567890123456789012345678901234567890"),
				Price:    1000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRenewRequest(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAutoRenewalConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *RenewalInfo
		valid  bool
	}{
		{
			name: "valid config",
			config: &RenewalInfo{
				CID:             "QmTestHash",
				RenewalDuration: 30,
				Token:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
				Price:           1000,
				Enabled:         true,
				CreatedAt:       time.Now(),
				NextRenewalAt:   time.Now().Add(30 * 24 * time.Hour),
			},
			valid: true,
		},
		{
			name: "invalid duration",
			config: &RenewalInfo{
				CID:             "QmTestHash",
				RenewalDuration: -1, // Invalid: negative duration
				Token:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
				Price:           1000,
				Enabled:         true,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAutoRenewalConfig(tt.config)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAutoRenewalService_StartStop(t *testing.T) {
	// Mock context params
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a mock auto-renewal service
	service := &AutoRenewalService{
		ctx:           ctx,
		cancel:        cancel,
		checkInterval: 100 * time.Millisecond, // Short interval for testing
		running:       false,
	}

	// Test starting the service
	err := service.Start()
	require.NoError(t, err)
	assert.True(t, service.IsRunning())

	// Test starting already running service
	err = service.Start()
	assert.Error(t, err)

	// Test stopping the service
	err = service.Stop()
	require.NoError(t, err)
	assert.False(t, service.IsRunning())

	// Test stopping already stopped service
	err = service.Stop()
	assert.Error(t, err)
}

func TestCalculateRenewalCost(t *testing.T) {
	tests := []struct {
		name         string
		shardSize    int64
		price        int64
		duration     int
		rate         int64
		shardCount   int
		expectedCost int64
	}{
		{
			name:         "basic calculation",
			shardSize:    1024 * 1024, // 1MB
			price:        1000,        // 1000 µBTT per GiB per day
			duration:     30,          // 30 days
			rate:         1000000,     // 1:1 rate
			shardCount:   10,          // 10 shards
			expectedCost: 290,         // Expected cost calculation (corrected)
		},
		{
			name:         "zero duration",
			shardSize:    1024 * 1024,
			price:        1000,
			duration:     0,
			rate:         1000000,
			shardCount:   10,
			expectedCost: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := calculateRenewalCost(tt.shardSize, tt.price, tt.duration, tt.rate, tt.shardCount)
			assert.Equal(t, tt.expectedCost, cost)
		})
	}
}

// Helper functions for testing

func validateRenewRequest(req *RenewRequest) error {
	if req.CID == "" {
		return fmt.Errorf("file hash cannot be empty")
	}
	if req.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}

func validateAutoRenewalConfig(config *RenewalInfo) error {
	if config.RenewalDuration <= 0 {
		return fmt.Errorf("renewal duration must be positive")
	}
	return nil
}

func calculateRenewalCost(shardSize, price int64, duration int, rate int64, shardCount int) int64 {
	if duration <= 0 {
		return 0
	}

	// Convert shard size to GiB
	gib := float64(shardSize) / (1024 * 1024 * 1024)

	// Calculate cost per shard
	costPerShard := int64(gib * float64(price) * float64(duration))

	// Apply rate conversion
	costPerShard = costPerShard * rate / 1000000

	// Multiply by shard count
	return costPerShard * int64(shardCount)
}

// Benchmark tests

func BenchmarkRenewalCostCalculation(b *testing.B) {
	shardSize := int64(1024 * 1024) // 1MB
	price := int64(1000)
	duration := 30
	rate := int64(1000000)
	shardCount := 10

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateRenewalCost(shardSize, price, duration, rate, shardCount)
	}
}

func BenchmarkAutoRenewalConfigValidation(b *testing.B) {
	config := &RenewalInfo{
		CID:             "QmTestHash",
		RenewalDuration: 30,
		Token:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Price:           1000,
		Enabled:         true,
		CreatedAt:       time.Now(),
		NextRenewalAt:   time.Now().Add(30 * 24 * time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateAutoRenewalConfig(config)
	}
}
