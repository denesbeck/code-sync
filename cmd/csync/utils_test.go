package main

import (
	"testing"
	"time"
)

func Test_ParsePath(t *testing.T) {
	tests := []struct {
		fullPath         string
		expectedPath     string
		expectedFileName string
	}{
		// Existing tests
		{"/path/to/file.txt", "/path/to/", "file.txt"},
		{"file.txt", "", "file.txt"},
		{"/path/to/another/file.txt", "/path/to/another/", "file.txt"},

		// Edge case: empty path (prevents panic)
		{"", "", ""},

		// Trailing slashes
		{"path/to/file.txt/", "path/to/", "file.txt"},
		{"/path/to/file.txt/", "/path/to/", "file.txt"},
		{"file.txt/", "", "file.txt"},

		// Windows paths (backslash separators)
		{"path\\to\\file.txt", "path/to/", "file.txt"},
		{"C:\\path\\to\\file.txt", "C:/path/to/", "file.txt"},
		{"\\path\\to\\file.txt", "/path/to/", "file.txt"},

		// Mixed separators
		{"path/to\\file.txt", "path/to/", "file.txt"},

		// Relative paths
		{"./file.txt", "", "file.txt"},
		{"./path/to/file.txt", "path/to/", "file.txt"},

		// Multiple trailing slashes
		{"path/to/file.txt///", "path/to/", "file.txt"},
	}
	for _, test := range tests {
		path, fileName := ParsePath(test.fullPath)
		if path != test.expectedPath {
			t.Errorf("Input: '%s' - Expected path '%s', but got '%s'", test.fullPath, test.expectedPath, path)
		}
		if fileName != test.expectedFileName {
			t.Errorf("Input: '%s' - Expected file name '%s', but got '%s'", test.fullPath, test.expectedFileName, fileName)
		}
	}
}

func Test_FindIndex(t *testing.T) {
	tests := []struct {
		arr      []string
		val      string
		expected int
	}{
		{[]string{"a", "b", "c"}, "b", 1},
		{[]string{"x", "y", "z"}, "a", -1},
		{[]string{"1", "2", "3"}, "3", 2},
	}

	for _, test := range tests {
		index := FindIndex(test.arr, test.val)
		if index != test.expected {
			t.Errorf("Expected index %d, but got %d", test.expected, index)
		}
	}
}

func Test_IsValidBranchName(t *testing.T) {
	tests := []struct {
		branchName string
		expected   bool
	}{
		{"feature/awesome-feature", true},
		{"bugfix-issue-123", true},
		{"hotfix_urgent_fix", true},
		{"invalid branch name", false},
		{"", false},
		{"-invalid-start", false},
		{"_invalid_start", false},
		{"invalid--name--", false},
		{"invalid__name__", false},
		{"invalid//name", false},
		{"invalid@%#&", false},
	}

	for _, test := range tests {
		result := IsValidBranchName(test.branchName)
		if result != test.expected {
			t.Errorf("Expected %v for branch name '%s', but got %v", test.expected, test.branchName, result)
		}
	}
}

func Test_Capitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"WORLD", "World"},
		{"gO", "Go"},
		{"test", "Test"},
		{"USERNAME", "Username"},
		{"email", "Email"},
		{"a", "A"},
		{"Z", "Z"},
		{"mIxEd", "Mixed"},
	}
	for _, test := range tests {
		result := Capitalize(test.input)
		if result != test.expected {
			t.Errorf("Expected '%s' for input '%s', but got '%s'", test.expected, test.input, result)
		}
	}
}

func Test_TimeAgo(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		timestamp string
		expected  string
	}{
		// Just now (< 60 seconds)
		{
			name:      "just now - 30 seconds ago",
			timestamp: now.Add(-30 * time.Second).Format(time.RFC3339),
			expected:  "just now",
		},
		{
			name:      "just now - 59 seconds ago",
			timestamp: now.Add(-59 * time.Second).Format(time.RFC3339),
			expected:  "just now",
		},

		// Minutes
		{
			name:      "1 minute ago",
			timestamp: now.Add(-1 * time.Minute).Format(time.RFC3339),
			expected:  "1 minute ago",
		},
		{
			name:      "5 minutes ago",
			timestamp: now.Add(-5 * time.Minute).Format(time.RFC3339),
			expected:  "5 minutes ago",
		},
		{
			name:      "45 minutes ago",
			timestamp: now.Add(-45 * time.Minute).Format(time.RFC3339),
			expected:  "45 minutes ago",
		},

		// Hours
		{
			name:      "1 hour ago",
			timestamp: now.Add(-1 * time.Hour).Format(time.RFC3339),
			expected:  "1 hour ago",
		},
		{
			name:      "3 hours ago",
			timestamp: now.Add(-3 * time.Hour).Format(time.RFC3339),
			expected:  "3 hours ago",
		},
		{
			name:      "12 hours ago",
			timestamp: now.Add(-12 * time.Hour).Format(time.RFC3339),
			expected:  "12 hours ago",
		},

		// Days
		{
			name:      "1 day ago",
			timestamp: now.Add(-24 * time.Hour).Format(time.RFC3339),
			expected:  "1 day ago",
		},
		{
			name:      "4 days ago",
			timestamp: now.Add(-4 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "4 days ago",
		},
		{
			name:      "6 days ago",
			timestamp: now.Add(-6 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "6 days ago",
		},

		// Weeks
		{
			name:      "1 week ago",
			timestamp: now.Add(-7 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "1 week ago",
		},
		{
			name:      "2 weeks ago",
			timestamp: now.Add(-14 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "2 weeks ago",
		},
		{
			name:      "3 weeks ago",
			timestamp: now.Add(-21 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "3 weeks ago",
		},

		// Months
		{
			name:      "1 month ago",
			timestamp: now.Add(-30 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "1 month ago",
		},
		{
			name:      "2 months ago",
			timestamp: now.Add(-60 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "2 months ago",
		},
		{
			name:      "6 months ago",
			timestamp: now.Add(-180 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "6 months ago",
		},

		// Years
		{
			name:      "1 year ago",
			timestamp: now.Add(-365 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "1 year ago",
		},
		{
			name:      "2 years ago",
			timestamp: now.Add(-730 * 24 * time.Hour).Format(time.RFC3339),
			expected:  "2 years ago",
		},

		// Invalid timestamp
		{
			name:      "invalid timestamp",
			timestamp: "invalid-timestamp",
			expected:  "N/A",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TimeAgo(test.timestamp)
			if result != test.expected {
				t.Errorf("Expected '%s', but got '%s'", test.expected, result)
			}
		})
	}
}
