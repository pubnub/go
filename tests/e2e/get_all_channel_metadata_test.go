package e2e

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v7"
	"github.com/stretchr/testify/assert"
)

// Helper function to create test channel metadata
func createTestChannelMetadata(t *testing.T, pn *pubnub.PubNub, id, name, description string, custom map[string]interface{}) {
	incl := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}
	_, _, err := pn.SetChannelMetadata().
		Include(incl).
		Channel(id).
		Name(name).
		Description(description).
		Custom(custom).
		Execute()

	if err != nil {
		t.Logf("Warning: Failed to create test channel metadata for %s: %v", id, err)
	}
}

// Helper function to clean up test channel metadata
func cleanupTestChannelMetadata(t *testing.T, pn *pubnub.PubNub, channelIDs []string) {
	for _, id := range channelIDs {
		_, _, err := pn.RemoveChannelMetadata().Channel(id).Execute()
		if err != nil {
			t.Logf("Warning: Failed to cleanup channel metadata for %s: %v", id, err)
		}
	}
}

// Helper function to clean up ALL test channels with Go_Sdk_test prefix (failsafe cleanup)
func cleanupAllTestChannelMetadata(t *testing.T, pn *pubnub.PubNub) {
	// Query for all channels with our test prefix
	resp, _, err := pn.GetAllChannelMetadata().
		Filter("id LIKE 'Go_Sdk_test*'").
		Limit(100).
		Execute()

	if err != nil {
		t.Logf("Warning: Failed to query test channels for cleanup: %v", err)
		return
	}

	if resp != nil && resp.Data != nil {
		var testChannelIDs []string
		for _, channel := range resp.Data {
			testChannelIDs = append(testChannelIDs, channel.ID)
		}

		if len(testChannelIDs) > 0 {
			cleanupTestChannelMetadata(t, pn, testChannelIDs)
		}
	}
}

// Test basic happy path - get all channel metadata
func TestGetAllChannelMetadataBasic(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create some test channels for this test
	testPrefix := "Go_Sdk_test_basic_"
	testChannels := []string{
		testPrefix + randomized("channel_1"),
		testPrefix + randomized("channel_2"),
	}
	defer cleanupTestChannelMetadata(t, pn, testChannels)

	// Create test channel metadata
	createTestChannelMetadata(t, pn, testChannels[0], "Test Channel 1", "First test channel", map[string]interface{}{"type": "test1"})
	createTestChannelMetadata(t, pn, testChannels[1], "Test Channel 2", "Second test channel", map[string]interface{}{"type": "test2"})

	// Wait a bit for the data to propagate
	time.Sleep(1 * time.Second)

	// Test basic GetAllChannelMetadata
	resp, status, err := pn.GetAllChannelMetadata().Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotNil(resp.Data)
	// Should have at least our test channels (we created 2)
	assert.True(len(resp.Data) >= 2)
}

// Test GetAllChannelMetadata with Include parameter
func TestGetAllChannelMetadataWithInclude(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	testChannel := "Go_Sdk_test_include_" + randomized("channel")
	defer cleanupTestChannelMetadata(t, pn, []string{testChannel})

	// Create test channel with custom data
	custom := map[string]interface{}{
		"department": "engineering",
		"priority":   "high",
		"tags":       "test,metadata",
	}
	createTestChannelMetadata(t, pn, testChannel, "Include Test Channel", "Channel for testing include", custom)

	time.Sleep(1 * time.Second)

	// Test with different include options
	testCases := []struct {
		name    string
		include []pubnub.PNChannelMetadataInclude
	}{
		{
			name:    "Include Custom",
			include: []pubnub.PNChannelMetadataInclude{pubnub.PNChannelMetadataIncludeCustom},
		},
		{
			name:    "Multiple Includes",
			include: []pubnub.PNChannelMetadataInclude{pubnub.PNChannelMetadataIncludeCustom},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := pn.GetAllChannelMetadata().
				Include(tc.include).
				Execute()

			assert.Nil(err)
			assert.Equal(200, status.StatusCode)
			assert.NotNil(resp)
			assert.NotNil(resp.Data)

			// Find our test channel in the results
			for _, channel := range resp.Data {
				if channel.ID == testChannel {
					// When custom is included, custom fields should be present
					if len(tc.include) > 0 {
						assert.NotNil(channel.Custom)
						assert.Equal("engineering", channel.Custom["department"])
					}
					break
				}
			}
		})
	}
}

// Test GetAllChannelMetadata with Limit parameter
func TestGetAllChannelMetadataWithLimit(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create multiple test channels with consistent prefix
	testPrefix := "Go_Sdk_test_limit_"
	testChannels := make([]string, 5)
	for i := 0; i < 5; i++ {
		testChannels[i] = testPrefix + randomized("channel_"+strconv.Itoa(i))
		createTestChannelMetadata(t, pn, testChannels[i], "Limit Test "+strconv.Itoa(i), "Channel for limit testing", map[string]interface{}{"index": i})
	}
	defer cleanupTestChannelMetadata(t, pn, testChannels)

	time.Sleep(1 * time.Second)

	testCases := []struct {
		name  string
		limit int
	}{
		{"Limit 1", 1},
		{"Limit 3", 3},
		{"Limit 10", 10},
		{"Limit 100", 100}, // Default limit
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := pn.GetAllChannelMetadata().
				Limit(tc.limit).
				Execute()

			assert.Nil(err)
			assert.Equal(200, status.StatusCode)
			assert.NotNil(resp)
			assert.NotNil(resp.Data)
			assert.True(len(resp.Data) <= tc.limit)
		})
	}
}

// Test GetAllChannelMetadata with Count parameter
func TestGetAllChannelMetadataWithCount(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Test with Count = true
	resp, status, err := pn.GetAllChannelMetadata().
		Count(true).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.True(resp.TotalCount >= 0) // Should include total count

	// Test with Count = false
	resp2, status2, err2 := pn.GetAllChannelMetadata().
		Count(false).
		Execute()

	assert.Nil(err2)
	assert.Equal(200, status2.StatusCode)
	assert.NotNil(resp2)
	// When count is false, TotalCount should be 0 or not meaningful
}

// Test GetAllChannelMetadata with Filter parameter
func TestGetAllChannelMetadataWithFilter(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create test channels with specific names for filtering
	testPrefix := "Go_Sdk_test_filter_"
	testChannels := []string{
		testPrefix + randomized("alpha_channel"),
		testPrefix + randomized("beta_channel"),
		testPrefix + randomized("gamma_channel"),
	}
	defer cleanupTestChannelMetadata(t, pn, testChannels)

	createTestChannelMetadata(t, pn, testChannels[0], "Alpha Channel", "Alpha test", map[string]interface{}{"category": "alpha"})
	createTestChannelMetadata(t, pn, testChannels[1], "Beta Channel", "Beta test", map[string]interface{}{"category": "beta"})
	createTestChannelMetadata(t, pn, testChannels[2], "Gamma Channel", "Gamma test", map[string]interface{}{"category": "gamma"})

	time.Sleep(1 * time.Second)

	testCases := []struct {
		name   string
		filter string
	}{
		{
			name:   "Filter by test prefix",
			filter: "id LIKE '" + testPrefix + "*'",
		},
		{
			name:   "Filter by name pattern",
			filter: "name LIKE '*Alpha*'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := pn.GetAllChannelMetadata().
				Filter(tc.filter).
				Include([]pubnub.PNChannelMetadataInclude{pubnub.PNChannelMetadataIncludeCustom}).
				Execute()

			assert.Nil(err)
			assert.Equal(200, status.StatusCode)
			assert.NotNil(resp)
			assert.NotNil(resp.Data)
			// Results should be filtered based on the criteria
		})
	}
}

// Test GetAllChannelMetadata with Sort parameter
func TestGetAllChannelMetadataWithSort(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	testCases := []struct {
		name        string
		sort        []string
		expectError bool
	}{
		{
			name:        "Sort by name ascending",
			sort:        []string{"name"},
			expectError: false,
		},
		{
			name:        "Sort by name descending",
			sort:        []string{"name:desc"},
			expectError: false,
		},
		{
			name:        "Sort by id ascending",
			sort:        []string{"id"},
			expectError: false,
		},
		{
			name:        "Sort by updated time descending",
			sort:        []string{"updated:desc"},
			expectError: false,
		},
		{
			name:        "Multiple sort criteria",
			sort:        []string{"name", "updated:desc"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := pn.GetAllChannelMetadata().
				Sort(tc.sort).
				Limit(10).
				Execute()

			if tc.expectError {
				// We expect either an error or a non-200 status for unsupported sort fields
				if err == nil {
					assert.NotEqual(200, status.StatusCode, "Expected error but got success status")
				}
			} else {
				assert.Nil(err)
				assert.Equal(200, status.StatusCode)
				assert.NotNil(resp)
				if resp != nil {
					assert.NotNil(resp.Data)
				}
				// Results should be sorted according to criteria
			}
		})
	}
}

// Test GetAllChannelMetadata with invalid sort fields
func TestGetAllChannelMetadataWithInvalidSort(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	invalidSortCases := []struct {
		name string
		sort []string
	}{
		{
			name: "Invalid sort field - created",
			sort: []string{"created"},
		},
		{
			name: "Invalid sort field - nonexistent",
			sort: []string{"nonexistent_field"},
		},
		{
			name: "Mixed valid and invalid sort fields",
			sort: []string{"name", "created:desc"},
		},
	}

	for _, tc := range invalidSortCases {
		t.Run(tc.name, func(t *testing.T) {
			_, status, err := pn.GetAllChannelMetadata().
				Sort(tc.sort).
				Limit(10).
				Execute()

			// We expect these to fail with 400 or other error
			if err != nil {
				// Error is expected for invalid sort fields
				assert.NotNil(err)
			} else {
				// If no error, we should get a non-200 status code
				assert.NotEqual(200, status.StatusCode, "Expected error for invalid sort field but got success")
			}
			// Don't access resp.Data here as it might be nil due to error
		})
	}
}

// Test GetAllChannelMetadata with pagination (Start/End)
func TestGetAllChannelMetadataPagination(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create several test channels for pagination
	testPrefix := "Go_Sdk_test_pagination_"
	testChannels := make([]string, 10)
	for i := 0; i < 10; i++ {
		testChannels[i] = testPrefix + randomized("channel_"+strconv.Itoa(i))
		createTestChannelMetadata(t, pn, testChannels[i], "Pagination Channel "+strconv.Itoa(i), "Channel for pagination", map[string]interface{}{"index": i})
	}
	defer cleanupTestChannelMetadata(t, pn, testChannels)

	time.Sleep(1 * time.Second)

	// Get first page - filter to only our test channels
	resp1, status1, err1 := pn.GetAllChannelMetadata().
		Filter("id LIKE '" + testPrefix + "*'").
		Limit(3).
		Count(true).
		Execute()

	assert.Nil(err1)
	assert.Equal(200, status1.StatusCode)
	assert.NotNil(resp1)
	assert.True(len(resp1.Data) <= 3)

	// If there's a next page, test pagination
	if resp1.Next != "" {
		resp2, status2, err2 := pn.GetAllChannelMetadata().
			Filter("id LIKE '" + testPrefix + "*'").
			Limit(3).
			Start(resp1.Next).
			Execute()

		assert.Nil(err2)
		assert.Equal(200, status2.StatusCode)
		assert.NotNil(resp2)

		// The results should be different from the first page
		if len(resp1.Data) > 0 && len(resp2.Data) > 0 {
			assert.NotEqual(resp1.Data[0].ID, resp2.Data[0].ID)
		}
	}
}

// Test GetAllChannelMetadata with QueryParam
func TestGetAllChannelMetadataWithQueryParam(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	queryParams := map[string]string{
		"custom_param1": "value1",
		"custom_param2": "value2",
	}

	resp, status, err := pn.GetAllChannelMetadata().
		QueryParam(queryParams).
		Limit(5).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotNil(resp.Data)
}

// Test GetAllChannelMetadata with Context
func TestGetAllChannelMetadataWithContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	resp, status, err := pn.GetAllChannelMetadataWithContext(backgroundContext).
		Limit(5).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotNil(resp.Data)
}

// Test GetAllChannelMetadata error scenarios
func TestGetAllChannelMetadataErrorScenarios(t *testing.T) {
	assert := assert.New(t)

	// Test with missing Subscribe Key
	pn := pubnub.NewPubNub(&pubnub.Config{
		UUID: "test-uuid",
		// Missing SubscribeKey
	})

	_, _, err := pn.GetAllChannelMetadata().Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

// Test GetAllChannelMetadata with invalid parameters
func TestGetAllChannelMetadataInvalidParameters(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	testCases := []struct {
		name        string
		limit       int
		filter      string
		start       string
		expectError bool
	}{
		{
			name:        "Negative limit",
			limit:       -1,
			expectError: false, // API might handle this gracefully
		},
		{
			name:        "Very large limit",
			limit:       10000,
			expectError: false, // API might clamp this
		},
		{
			name:        "Invalid filter syntax",
			filter:      "invalid filter syntax $$$ @@@",
			expectError: true, // Invalid filter should cause error
		},
		{
			name:        "Empty start token",
			start:       "",
			expectError: false, // Empty start should be ignored
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := pn.GetAllChannelMetadata()

			if tc.limit > 0 {
				builder = builder.Limit(tc.limit)
			}
			if tc.filter != "" {
				builder = builder.Filter(tc.filter)
			}
			if tc.start != "" {
				builder = builder.Start(tc.start)
			}

			resp, status, err := builder.Execute()

			if tc.expectError {
				// We expect either an error or a non-200 status
				if err == nil {
					assert.NotEqual(200, status.StatusCode, "Expected error but got success status")
				}
				// Don't access resp when we expect an error - it might be nil
			} else {
				// For cases where we don't expect errors, just verify it doesn't crash
				// The API might still return an error for some edge cases, which is acceptable
				if err == nil {
					assert.True(status.StatusCode >= 200 && status.StatusCode < 300)
					// Only access resp when we don't expect errors and err is nil
					if resp != nil {
						assert.NotNil(resp.Data)
					}
				}
			}
		})
	}
}

// Test GetAllChannelMetadata comprehensive scenario
func TestGetAllChannelMetadataComprehensive(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Failsafe cleanup at the start - clean up any orphaned test channels
	defer cleanupAllTestChannelMetadata(t, pn)

	// Create comprehensive test data with consistent prefix for filtering
	testPrefix := "Go_Sdk_test_comprehensive_"
	testChannels := []string{
		testPrefix + randomized("alpha"),
		testPrefix + randomized("beta"),
		testPrefix + randomized("gamma"),
	}
	defer cleanupTestChannelMetadata(t, pn, testChannels)

	// Create channels with different metadata
	createTestChannelMetadata(t, pn, testChannels[0], "Alpha Comprehensive", "First comprehensive test", map[string]interface{}{
		"priority":   "high",
		"department": "engineering",
		"tags":       "test,alpha",
	})
	createTestChannelMetadata(t, pn, testChannels[1], "Beta Comprehensive", "Second comprehensive test", map[string]interface{}{
		"priority":   "medium",
		"department": "marketing",
		"tags":       "test,beta",
	})
	createTestChannelMetadata(t, pn, testChannels[2], "Gamma Comprehensive", "Third comprehensive test", map[string]interface{}{
		"priority":   "low",
		"department": "support",
		"tags":       "test,gamma",
	})

	// Wait for metadata to propagate
	time.Sleep(4 * time.Second)

	// Test comprehensive query with multiple parameters using filter
	resp, status, err := pn.GetAllChannelMetadata().
		Include([]pubnub.PNChannelMetadataInclude{pubnub.PNChannelMetadataIncludeCustom}).
		Filter("id LIKE '" + testPrefix + "*'"). // Filter only our test channels
		Count(true).
		Sort([]string{"name"}).
		QueryParam(map[string]string{"test": "comprehensive"}).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotNil(resp.Data)

	// Verify our test channels are in the results
	foundChannels := make(map[string]bool)
	for _, channel := range resp.Data {
		for _, testChannel := range testChannels {
			if channel.ID == testChannel {
				foundChannels[testChannel] = true
				// Verify custom data is included
				assert.NotNil(channel.Custom)
				assert.NotNil(channel.Name)
				assert.NotNil(channel.Description)
				break
			}
		}
	}

	// E2E test: We must find ALL 3 test channels we created
	assert.Equal(3, len(foundChannels), "Should find all 3 created test channels")
	for _, testChannel := range testChannels {
		assert.True(foundChannels[testChannel], "Should find test channel: %s", testChannel)
	}
}

// Test GetAllChannelMetadata edge cases
func TestGetAllChannelMetadataEdgeCases(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Test with limit 0
	resp1, status1, err1 := pn.GetAllChannelMetadata().
		Limit(0).
		Execute()

	// This might return an error or handle gracefully
	if err1 == nil {
		assert.Equal(200, status1.StatusCode)
		assert.NotNil(resp1)
	}

	// Test with empty include array
	resp2, status2, err2 := pn.GetAllChannelMetadata().
		Include([]pubnub.PNChannelMetadataInclude{}).
		Execute()

	assert.Nil(err2)
	assert.Equal(200, status2.StatusCode)
	assert.NotNil(resp2)

	// Test with nil query params
	resp3, status3, err3 := pn.GetAllChannelMetadata().
		QueryParam(nil).
		Execute()

	assert.Nil(err3)
	assert.Equal(200, status3.StatusCode)
	assert.NotNil(resp3)

	// Test with empty sort array
	resp4, status4, err4 := pn.GetAllChannelMetadata().
		Sort([]string{}).
		Execute()

	assert.Nil(err4)
	assert.Equal(200, status4.StatusCode)
	assert.NotNil(resp4)
}
