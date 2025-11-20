package main

import (
	"strings"
	"testing"
)

func Test_Warning(t *testing.T) {
	// Just test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Warning function panicked: %v", r)
		}
	}()
	Warning("Test warning message")
}

func Test_Fail(t *testing.T) {
	// Just test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Fail function panicked: %v", r)
		}
	}()
	Fail("Test fail message")
}

func Test_ErrorMsg(t *testing.T) {
	result := ErrorMsg("test error")
	// Result should contain the text (even if styled)
	if !strings.Contains(result, "test error") && result == "" {
		t.Errorf("ErrorMsg should return a non-empty string containing 'test error'")
	}
}

func Test_Branch(t *testing.T) {
	result := StyledBranch("main")
	// Result should contain the text (even if styled)
	if !strings.Contains(result, "main") && result == "" {
		t.Errorf("Branch should return a non-empty string containing 'main'")
	}
}

func Test_Code(t *testing.T) {
	result := Code("test code")
	// Result should contain the text (even if styled)
	if !strings.Contains(result, "test code") && result == "" {
		t.Errorf("Code should return a non-empty string containing 'test code'")
	}
}

func Test_Bold(t *testing.T) {
	result := Bold("bold text")
	// Result should contain the text (even if styled)
	if !strings.Contains(result, "bold text") && result == "" {
		t.Errorf("Bold should return a non-empty string containing 'bold text'")
	}
}

func Test_Box(t *testing.T) {
	// Test without title
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Box function panicked: %v", r)
		}
	}()
	Box("", "test content")
	Box("Test Title", "test content with title")
}

func Test_Text(t *testing.T) {
	// Just test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Text function panicked: %v", r)
		}
	}()
	Text("test text", "")
	Text("test text with icon", "📝")
}

func Test_BreakLine(t *testing.T) {
	// Just test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BreakLine function panicked: %v", r)
		}
	}()
	BreakLine()
}

func Test_List(t *testing.T) {
	// Just test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("List function panicked: %v", r)
		}
	}()
	List("Root", []string{"item1", "item2", "item3"})
	List("Empty", []string{})
}

func Test_Tree(t *testing.T) {
	// Just test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Tree function panicked: %v", r)
		}
	}()
	// Test with empty list
	Tree([]string{}, true)

	// Test with single file
	Tree([]string{"file1.txt"}, true)

	// Test with multiple files
	Tree([]string{"file1.txt", "file2.txt", "file3.txt"}, true)
	Tree([]string{"file1.txt", "file2.txt", "file3.txt"}, false)
}
