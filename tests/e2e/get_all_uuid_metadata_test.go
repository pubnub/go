package e2e

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v7"
	"github.com/stretchr/testify/assert"
)

// Helper function to create test UUID metadata
func createTestUUIDMetadata(t *testing.T, pn *pubnub.PubNub, id, name, email string, custom map[string]interface{}) {
	incl := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
		pubnub.PNUUIDMetadataIncludeStatus,
		pubnub.PNUUIDMetadataIncludeType,
	}
	_, _, err := pn.SetUUIDMetadata().
		Include(incl).
		UUID(id).
		Name(name).
		Email(email).
		Custom(custom).
		Status("active").
		Type("test").
		Execute()

	if err != nil {
		t.Logf("Warning: Failed to create test UUID metadata for %s: %v", id, err)
	}
}

// Helper function to clean up test UUID metadata
func cleanupTestUUIDMetadata(t *testing.T, pn *pubnub.PubNub, uuidIDs []string) {
	for _, id := range uuidIDs {
		_, _, err := pn.RemoveUUIDMetadata().UUID(id).Execute()
		if err != nil {
			t.Logf("Warning: Failed to cleanup UUID metadata for %s: %v", id, err)
		}
	}
}

// Helper function to clean up ALL test UUIDs with Go_Sdk_test prefix (failsafe cleanup)
func cleanupAllTestUUIDMetadata(t *testing.T, pn *pubnub.PubNub) {
	// Query for all UUIDs with our test prefix
	resp, _, err := pn.GetAllUUIDMetadata().
		Filter("id LIKE 'Go_Sdk_test*'").
		Limit(100).
		Execute()

	if err != nil {
		t.Logf("Warning: Failed to query test UUIDs for cleanup: %v", err)
		return
	}

	if resp != nil && resp.Data != nil {
		var testUUIDIDs []string
		for _, uuid := range resp.Data {
			testUUIDIDs = append(testUUIDIDs, uuid.ID)
		}

		if len(testUUIDIDs) > 0 {
			cleanupTestUUIDMetadata(t, pn, testUUIDIDs)
		}
	}
}

// Test basic happy path - get all UUID metadata
func TestGetAllUUIDMetadataBasic(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create some test UUIDs for this test
	testPrefix := "Go_Sdk_test_uuid_basic_"
	testUUIDs := []string{
		testPrefix + randomized("uuid_1"),
		testPrefix + randomized("uuid_2"),
	}
	defer cleanupTestUUIDMetadata(t, pn, testUUIDs)

	// Create test UUID metadata
	createTestUUIDMetadata(t, pn, testUUIDs[0], "Test User 1", "user1@test.com", map[string]interface{}{"type": "test1"})
	createTestUUIDMetadata(t, pn, testUUIDs[1], "Test User 2", "user2@test.com", map[string]interface{}{"type": "test2"})

	// Wait a bit for the data to propagate
	time.Sleep(1 * time.Second)

	// Test basic GetAllUUIDMetadata
	resp, status, err := pn.GetAllUUIDMetadata().Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotNil(resp.Data)
	// Should have at least our test UUIDs (we created 2)
	assert.True(len(resp.Data) >= 2)
}

// Test GetAllUUIDMetadata with Include parameter
func TestGetAllUUIDMetadataWithInclude(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	testUUID := "Go_Sdk_test_uuid_include_" + randomized("uuid")
	defer cleanupTestUUIDMetadata(t, pn, []string{testUUID})

	// Create test UUID with custom data
	custom := map[string]interface{}{
		"department": "engineering",
		"role":       "developer",
		"tags":       "test,metadata",
	}
	createTestUUIDMetadata(t, pn, testUUID, "Include Test User", "include@test.com", custom)

	time.Sleep(1 * time.Second)

	// Test with different include options
	testCases := []struct {
		name    string
		include []pubnub.PNUUIDMetadataInclude
	}{
		{
			name:    "Include Custom",
			include: []pubnub.PNUUIDMetadataInclude{pubnub.PNUUIDMetadataIncludeCustom},
		},
		{
			name:    "Multiple Includes",
			include: []pubnub.PNUUIDMetadataInclude{pubnub.PNUUIDMetadataIncludeCustom, pubnub.PNUUIDMetadataIncludeStatus, pubnub.PNUUIDMetadataIncludeType},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := pn.GetAllUUIDMetadata().
				Include(tc.include).
				Execute()

			assert.Nil(err)
			assert.Equal(200, status.StatusCode)
			assert.NotNil(resp)
			assert.NotNil(resp.Data)

			// Find our test UUID in the results
			for _, uuid := range resp.Data {
				if uuid.ID == testUUID {
					// When custom is included, custom fields should be present
					if len(tc.include) > 0 {
						assert.NotNil(uuid.Custom)
						assert.Equal("engineering", uuid.Custom["department"])
					}
					break
				}
			}
		})
	}
}

// Test GetAllUUIDMetadata with Limit parameter
func TestGetAllUUIDMetadataWithLimit(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create multiple test UUIDs with consistent prefix
	testPrefix := "Go_Sdk_test_uuid_limit_"
	testUUIDs := make([]string, 5)
	for i := 0; i < 5; i++ {
		testUUIDs[i] = testPrefix + randomized("uuid_"+strconv.Itoa(i))
		createTestUUIDMetadata(t, pn, testUUIDs[i], "Limit Test "+strconv.Itoa(i), "limit"+strconv.Itoa(i)+"@test.com", map[string]interface{}{"index": i})
	}
	defer cleanupTestUUIDMetadata(t, pn, testUUIDs)

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
			resp, status, err := pn.GetAllUUIDMetadata().
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

// Test GetAllUUIDMetadata with Count parameter
func TestGetAllUUIDMetadataWithCount(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Test with Count = true
	resp, status, err := pn.GetAllUUIDMetadata().
		Count(true).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.True(resp.TotalCount >= 0) // Should include total count

	// Test with Count = false
	resp2, status2, err2 := pn.GetAllUUIDMetadata().
		Count(false).
		Execute()

	assert.Nil(err2)
	assert.Equal(200, status2.StatusCode)
	assert.NotNil(resp2)
	// When count is false, TotalCount should be 0 or not meaningful
}

// Test GetAllUUIDMetadata with Filter parameter
func TestGetAllUUIDMetadataWithFilter(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create test UUIDs with specific names for filtering
	testPrefix := "Go_Sdk_test_uuid_filter_"
	testUUIDs := []string{
		testPrefix + randomized("alpha_uuid"),
		testPrefix + randomized("beta_uuid"),
		testPrefix + randomized("gamma_uuid"),
	}
	defer cleanupTestUUIDMetadata(t, pn, testUUIDs)

	createTestUUIDMetadata(t, pn, testUUIDs[0], "Alpha User", "alpha@test.com", map[string]interface{}{"category": "alpha"})
	createTestUUIDMetadata(t, pn, testUUIDs[1], "Beta User", "beta@test.com", map[string]interface{}{"category": "beta"})
	createTestUUIDMetadata(t, pn, testUUIDs[2], "Gamma User", "gamma@test.com", map[string]interface{}{"category": "gamma"})

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
			resp, status, err := pn.GetAllUUIDMetadata().
				Filter(tc.filter).
				Include([]pubnub.PNUUIDMetadataInclude{pubnub.PNUUIDMetadataIncludeCustom, pubnub.PNUUIDMetadataIncludeStatus, pubnub.PNUUIDMetadataIncludeType}).
				Execute()

			assert.Nil(err)
			assert.Equal(200, status.StatusCode)
			assert.NotNil(resp)
			assert.NotNil(resp.Data)
			// Results should be filtered based on the criteria
		})
	}
}

// Test GetAllUUIDMetadata with Sort parameter
func TestGetAllUUIDMetadataWithSort(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
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
			resp, status, err := pn.GetAllUUIDMetadata().
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

// Test GetAllUUIDMetadata with invalid sort fields
func TestGetAllUUIDMetadataWithInvalidSort(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
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
			_, status, err := pn.GetAllUUIDMetadata().
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

// Test GetAllUUIDMetadata with pagination (Start/End)
func TestGetAllUUIDMetadataPagination(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create several test UUIDs for pagination
	testPrefix := "Go_Sdk_test_uuid_pagination_"
	testUUIDs := make([]string, 10)
	for i := 0; i < 10; i++ {
		testUUIDs[i] = testPrefix + randomized("uuid_"+strconv.Itoa(i))
		createTestUUIDMetadata(t, pn, testUUIDs[i], "Pagination User "+strconv.Itoa(i), "pagination"+strconv.Itoa(i)+"@test.com", map[string]interface{}{"index": i})
	}
	defer cleanupTestUUIDMetadata(t, pn, testUUIDs)

	time.Sleep(1 * time.Second)

	// Get first page - filter to only our test UUIDs
	resp1, status1, err1 := pn.GetAllUUIDMetadata().
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
		resp2, status2, err2 := pn.GetAllUUIDMetadata().
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

// Test GetAllUUIDMetadata with QueryParam
func TestGetAllUUIDMetadataWithQueryParam(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	queryParams := map[string]string{
		"custom_param1": "value1",
		"custom_param2": "value2",
	}

	resp, status, err := pn.GetAllUUIDMetadata().
		QueryParam(queryParams).
		Limit(5).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotNil(resp.Data)
}

// Test GetAllUUIDMetadata with Context
func TestGetAllUUIDMetadataWithContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	resp, status, err := pn.GetAllUUIDMetadataWithContext(backgroundContext).
		Limit(5).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotNil(resp.Data)
}

// Test GetAllUUIDMetadata error scenarios
func TestGetAllUUIDMetadataErrorScenarios(t *testing.T) {
	assert := assert.New(t)

	// Test with missing Subscribe Key
	pn := pubnub.NewPubNub(&pubnub.Config{
		UUID: "test-uuid",
		// Missing SubscribeKey
	})

	_, _, err := pn.GetAllUUIDMetadata().Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

// Test GetAllUUIDMetadata with invalid parameters
func TestGetAllUUIDMetadataInvalidParameters(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
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
			builder := pn.GetAllUUIDMetadata()

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

// Test GetAllUUIDMetadata comprehensive scenario
func TestGetAllUUIDMetadataComprehensive(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Failsafe cleanup at the start - clean up any orphaned test UUIDs
	defer cleanupAllTestUUIDMetadata(t, pn)

	// Create comprehensive test data with consistent prefix for filtering
	testPrefix := "Go_Sdk_test_uuid_comprehensive_"
	testUUIDs := []string{
		testPrefix + randomized("alpha"),
		testPrefix + randomized("beta"),
		testPrefix + randomized("gamma"),
	}
	defer cleanupTestUUIDMetadata(t, pn, testUUIDs)

	// Create UUIDs with different metadata
	createTestUUIDMetadata(t, pn, testUUIDs[0], "Alpha Comprehensive", "alpha@comprehensive.com", map[string]interface{}{
		"role":       "admin",
		"department": "engineering",
		"tags":       "test,alpha",
	})
	createTestUUIDMetadata(t, pn, testUUIDs[1], "Beta Comprehensive", "beta@comprehensive.com", map[string]interface{}{
		"role":       "user",
		"department": "marketing",
		"tags":       "test,beta",
	})
	createTestUUIDMetadata(t, pn, testUUIDs[2], "Gamma Comprehensive", "gamma@comprehensive.com", map[string]interface{}{
		"role":       "viewer",
		"department": "support",
		"tags":       "test,gamma",
	})

	// Use retry mechanism to handle eventual consistency in CI/CD
	checkMetadataCall := func() error {
		resp, status, err := pn.GetAllUUIDMetadata().
			Include([]pubnub.PNUUIDMetadataInclude{pubnub.PNUUIDMetadataIncludeCustom, pubnub.PNUUIDMetadataIncludeStatus, pubnub.PNUUIDMetadataIncludeType}).
			Filter("id LIKE '" + testPrefix + "*'"). // Filter only our test UUIDs
			Count(true).
			Sort([]string{"name"}).
			QueryParam(map[string]string{"test": "comprehensive"}).
			Execute()

		if err != nil {
			return err
		}
		if status.StatusCode != 200 {
			return errors.New("status code not 200")
		}
		if resp == nil || resp.Data == nil {
			return errors.New("response or data is nil")
		}

		// Verify our test UUIDs are in the results
		foundUUIDs := make(map[string]bool)
		for _, uuid := range resp.Data {
			for _, testUUID := range testUUIDs {
				if uuid.ID == testUUID {
					foundUUIDs[testUUID] = true
					// Verify custom data is included
					assert.NotNil(uuid.Custom)
					assert.NotNil(uuid.Name)
					assert.NotNil(uuid.Email)
					break
				}
			}
		}

		// Check if all 3 test UUIDs are found
		if len(foundUUIDs) != 3 {
			return fmt.Errorf("expected 3 UUIDs, found %d", len(foundUUIDs))
		}
		for _, testUUID := range testUUIDs {
			if !foundUUIDs[testUUID] {
				return fmt.Errorf("test UUID not found: %s", testUUID)
			}
		}

		return nil
	}

	// Use retry mechanism with 10 seconds max timeout, 500ms intervals
	checkFor(assert, time.Second*10, time.Millisecond*500, checkMetadataCall)
}

// Test GetAllUUIDMetadata edge cases
func TestGetAllUUIDMetadataEdgeCases(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	defer pn.Destroy() // Cleanup to prevent goroutine leaks
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Test with limit 0
	resp1, status1, err1 := pn.GetAllUUIDMetadata().
		Limit(0).
		Execute()

	// This might return an error or handle gracefully
	if err1 == nil {
		assert.Equal(200, status1.StatusCode)
		assert.NotNil(resp1)
	}

	// Test with empty include array
	resp2, status2, err2 := pn.GetAllUUIDMetadata().
		Include([]pubnub.PNUUIDMetadataInclude{}).
		Execute()

	assert.Nil(err2)
	assert.Equal(200, status2.StatusCode)
	assert.NotNil(resp2)

	// Test with nil query params
	resp3, status3, err3 := pn.GetAllUUIDMetadata().
		QueryParam(nil).
		Execute()

	assert.Nil(err3)
	assert.Equal(200, status3.StatusCode)
	assert.NotNil(resp3)

	// Test with empty sort array
	resp4, status4, err4 := pn.GetAllUUIDMetadata().
		Sort([]string{}).
		Execute()

	assert.Nil(err4)
	assert.Equal(200, status4.StatusCode)
	assert.NotNil(resp4)
}
