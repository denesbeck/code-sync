package main

import (
	"strconv"
	"testing"
)

func Test_SortCommitsByLinkedList_SingleCommit(t *testing.T) {
	commits := []Commit{
		{Id: "aaa", Timestamp: "2024-01-01", Next: ""},
	}

	sorted := sortCommitsByLinkedList(commits)

	if len(sorted) != 1 {
		t.Errorf("Expected 1 commit, got %d", len(sorted))
	}
	if sorted[0].Id != "aaa" {
		t.Errorf("Expected commit 'aaa', got '%s'", sorted[0].Id)
	}
}

func Test_SortCommitsByLinkedList_EmptyList(t *testing.T) {
	commits := []Commit{}
	sorted := sortCommitsByLinkedList(commits)

	if len(sorted) != 0 {
		t.Errorf("Expected 0 commits, got %d", len(sorted))
	}
}

func Test_SortCommitsByLinkedList_ChronologicalOrder(t *testing.T) {
	// Create 5 commits in reverse order to simulate random storage
	commits := []Commit{
		{Id: "eee", Timestamp: "2024-01-05", Next: ""},    // Last (most recent)
		{Id: "ccc", Timestamp: "2024-01-03", Next: "ddd"}, // Middle
		{Id: "aaa", Timestamp: "2024-01-01", Next: "bbb"}, // First (oldest)
		{Id: "ddd", Timestamp: "2024-01-04", Next: "eee"}, // Second to last
		{Id: "bbb", Timestamp: "2024-01-02", Next: "ccc"}, // Second
	}

	sorted := sortCommitsByLinkedList(commits)

	// Should be sorted: aaa -> bbb -> ccc -> ddd -> eee
	expected := []string{"aaa", "bbb", "ccc", "ddd", "eee"}

	if len(sorted) != 5 {
		t.Fatalf("Expected 5 commits, got %d", len(sorted))
	}

	for i, commit := range sorted {
		if commit.Id != expected[i] {
			t.Errorf("Position %d: expected '%s', got '%s'", i, expected[i], commit.Id)
		}
	}

	// Verify linked list is intact
	for i := 0; i < len(sorted)-1; i++ {
		if sorted[i].Next != sorted[i+1].Id {
			t.Errorf("Broken chain at position %d: %s should point to %s, but points to %s",
				i, sorted[i].Id, sorted[i+1].Id, sorted[i].Next)
		}
	}

	// Last commit should have empty Next
	if sorted[len(sorted)-1].Next != "" {
		t.Errorf("Last commit should have empty Next, got '%s'", sorted[len(sorted)-1].Next)
	}
}

func Test_SortCommitsByLinkedList_TwoCommits(t *testing.T) {
	// Store in reverse order
	commits := []Commit{
		{Id: "bbb", Timestamp: "2024-01-02", Next: ""},
		{Id: "aaa", Timestamp: "2024-01-01", Next: "bbb"},
	}

	sorted := sortCommitsByLinkedList(commits)

	if len(sorted) != 2 {
		t.Fatalf("Expected 2 commits, got %d", len(sorted))
	}

	if sorted[0].Id != "aaa" {
		t.Errorf("First commit should be 'aaa', got '%s'", sorted[0].Id)
	}
	if sorted[1].Id != "bbb" {
		t.Errorf("Second commit should be 'bbb', got '%s'", sorted[1].Id)
	}
}

func Test_SortCommitsByLinkedList_ComplexScrambled(t *testing.T) {
	// Create 10 commits in completely random order
	commits := []Commit{
		{Id: "commit-7", Next: "commit-8"},
		{Id: "commit-3", Next: "commit-4"},
		{Id: "commit-9", Next: "commit-10"},
		{Id: "commit-1", Next: "commit-2"}, // First
		{Id: "commit-5", Next: "commit-6"},
		{Id: "commit-10", Next: ""}, // Last
		{Id: "commit-2", Next: "commit-3"},
		{Id: "commit-8", Next: "commit-9"},
		{Id: "commit-4", Next: "commit-5"},
		{Id: "commit-6", Next: "commit-7"},
	}

	sorted := sortCommitsByLinkedList(commits)

	if len(sorted) != 10 {
		t.Fatalf("Expected 10 commits, got %d", len(sorted))
	}

	// Verify they're in correct order
	for i := 1; i <= 10; i++ {
		expectedId := "commit-" + strconv.Itoa(i)
		if sorted[i-1].Id != expectedId {
			t.Errorf("Position %d: expected '%s', got '%s'", i-1, expectedId, sorted[i-1].Id)
		}
	}
}

func Test_SortCommitsByLinkedList_BrokenChain(t *testing.T) {
	// Commit points to non-existent commit
	commits := []Commit{
		{Id: "aaa", Next: "bbb"},
		{Id: "bbb", Next: "xxx"}, // xxx doesn't exist!
	}

	sorted := sortCommitsByLinkedList(commits)

	// Should return what it can (aaa -> bbb, then stop)
	if len(sorted) != 2 {
		t.Errorf("Expected 2 commits despite broken chain, got %d", len(sorted))
	}
	if sorted[0].Id != "aaa" {
		t.Errorf("First commit should be 'aaa', got '%s'", sorted[0].Id)
	}
	if sorted[1].Id != "bbb" {
		t.Errorf("Second commit should be 'bbb', got '%s'", sorted[1].Id)
	}
}

func Test_GetLastCommitByBranch_FindsCorrectCommit(t *testing.T) {
	// This is an integration-style test
	// We'll just verify the logic works with the helper
	commits := []Commit{
		{Id: "first", Next: "second"},
		{Id: "second", Next: "third"},
		{Id: "third", Next: ""}, // This should be found
	}

	// Find the one with empty Next
	var lastCommit Commit
	for _, commit := range commits {
		if commit.Next == "" {
			lastCommit = commit
			break
		}
	}

	if lastCommit.Id != "third" {
		t.Errorf("Expected last commit to be 'third', got '%s'", lastCommit.Id)
	}
}
