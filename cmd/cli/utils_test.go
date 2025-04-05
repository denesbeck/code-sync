package cli

import (
	"testing"
)

func Test_ParsePath(t *testing.T) {
	tests := []struct {
		fullPath         string
		expectedPath     string
		expectedFileName string
	}{
		{"/path/to/file.txt", "/path/to/", "file.txt"},
		{"file.txt", "", "file.txt"},
		{"/path/to/another/file.txt", "/path/to/another/", "file.txt"},
	}

	for _, test := range tests {
		path, fileName := ParsePath(test.fullPath)
		if path != test.expectedPath {
			t.Errorf("Expected path %s, but got %s", test.expectedPath, path)
		}
		if fileName != test.expectedFileName {
			t.Errorf("Expected file name %s, but got %s", test.expectedFileName, fileName)
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
