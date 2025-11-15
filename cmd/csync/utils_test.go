package main

import (
	"testing"
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
