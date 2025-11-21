package main

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_CopyFile_PreservesPermissions(t *testing.T) {
	// Setup: Create temporary directory
	tmpDir := namespace + "test_copy_permissions"
	defer os.RemoveAll(tmpDir)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	tests := []struct {
		name        string
		permissions os.FileMode
		content     string
	}{
		{
			name:        "executable_script",
			permissions: 0755,
			content:     "#!/bin/bash\necho 'Hello World'\n",
		},
		{
			name:        "regular_file",
			permissions: 0644,
			content:     "Just a regular file\n",
		},
		{
			name:        "restricted_file",
			permissions: 0600,
			content:     "Secret data\n",
		},
		{
			name:        "owner_readwrite",
			permissions: 0700,
			content:     "Owner only data\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create source file with specific permissions
			srcPath := filepath.Join(tmpDir, "src_"+test.name)
			if err := os.WriteFile(srcPath, []byte(test.content), test.permissions); err != nil {
				t.Fatalf("Failed to create source file: %v", err)
			}

			// Verify source file has correct permissions
			srcInfo, err := os.Stat(srcPath)
			if err != nil {
				t.Fatalf("Failed to stat source file: %v", err)
			}
			if srcInfo.Mode().Perm() != test.permissions {
				t.Fatalf("Source file permissions not set correctly. Expected %o, got %o",
					test.permissions, srcInfo.Mode().Perm())
			}

			// Copy file
			dstPath := filepath.Join(tmpDir, "dst_"+test.name)
			if err := CopyFile(srcPath, dstPath); err != nil {
				t.Fatalf("CopyFile failed: %v", err)
			}

			// Verify destination file permissions match source
			dstInfo, err := os.Stat(dstPath)
			if err != nil {
				t.Fatalf("Failed to stat destination file: %v", err)
			}

			if dstInfo.Mode().Perm() != test.permissions {
				t.Errorf("Permissions not preserved. Expected %o, got %o",
					test.permissions, dstInfo.Mode().Perm())
			}

			// Verify content is also correct
			dstContent, err := os.ReadFile(dstPath)
			if err != nil {
				t.Fatalf("Failed to read destination file: %v", err)
			}
			if string(dstContent) != test.content {
				t.Errorf("Content not preserved. Expected %q, got %q", test.content, string(dstContent))
			}

			// Cleanup individual test files
			os.Remove(srcPath)
			os.Remove(dstPath)
		})
	}
}

func Test_CopyFile_PreservesContent(t *testing.T) {
	tmpDir := namespace + "test_copy_content"
	defer os.RemoveAll(tmpDir)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "empty_file",
			content: "",
		},
		{
			name:    "small_file",
			content: "Hello, World!\n",
		},
		{
			name:    "multiline_file",
			content: "Line 1\nLine 2\nLine 3\n",
		},
		{
			name:    "binary_content",
			content: "\x00\x01\x02\x03\xFF\xFE\xFD",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srcPath := filepath.Join(tmpDir, "src_"+test.name)
			dstPath := filepath.Join(tmpDir, "dst_"+test.name)

			// Create source file
			if err := os.WriteFile(srcPath, []byte(test.content), 0644); err != nil {
				t.Fatalf("Failed to create source file: %v", err)
			}

			// Copy file
			if err := CopyFile(srcPath, dstPath); err != nil {
				t.Fatalf("CopyFile failed: %v", err)
			}

			// Verify content
			dstContent, err := os.ReadFile(dstPath)
			if err != nil {
				t.Fatalf("Failed to read destination file: %v", err)
			}

			if string(dstContent) != test.content {
				t.Errorf("Content mismatch. Expected %q, got %q", test.content, string(dstContent))
			}

			// Cleanup
			os.Remove(srcPath)
			os.Remove(dstPath)
		})
	}
}

func Test_CopyFile_CreatesDirectories(t *testing.T) {
	tmpDir := namespace + "test_copy_dirs"
	defer os.RemoveAll(tmpDir)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(srcPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy to nested directory that doesn't exist
	dstPath := filepath.Join(tmpDir, "nested", "dirs", "destination.txt")
	if err := CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify file exists and has correct content
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		t.Errorf("Destination file was not created")
	}

	content, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("Content mismatch. Expected 'test content', got %q", string(content))
	}
}

func Test_CopyFile_OverwritesExisting(t *testing.T) {
	tmpDir := namespace + "test_copy_overwrite"
	defer os.RemoveAll(tmpDir)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "destination.txt")

	// Create source file
	if err := os.WriteFile(srcPath, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create destination file with old content
	if err := os.WriteFile(dstPath, []byte("old content"), 0644); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}

	// Copy (should overwrite)
	if err := CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify content is updated
	content, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(content) != "new content" {
		t.Errorf("Content not overwritten. Expected 'new content', got %q", string(content))
	}
}

func Test_CopyFile_ErrorCases(t *testing.T) {
	tmpDir := namespace + "test_copy_errors"
	defer os.RemoveAll(tmpDir)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	t.Run("source_does_not_exist", func(t *testing.T) {
		srcPath := filepath.Join(tmpDir, "nonexistent.txt")
		dstPath := filepath.Join(tmpDir, "destination.txt")

		err := CopyFile(srcPath, dstPath)
		if err == nil {
			t.Errorf("Expected error when source doesn't exist, got nil")
		}
	})

	t.Run("source_is_directory", func(t *testing.T) {
		srcPath := filepath.Join(tmpDir, "source_dir")
		dstPath := filepath.Join(tmpDir, "destination.txt")

		if err := os.MkdirAll(srcPath, 0755); err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		err := CopyFile(srcPath, dstPath)
		if err == nil {
			t.Errorf("Expected error when source is directory, got nil")
		}
		if err != os.ErrInvalid {
			t.Errorf("Expected os.ErrInvalid, got %v", err)
		}
	})
}

func Test_CopyFile_PreservesExecutableBit(t *testing.T) {
	tmpDir := namespace + "test_executable_bit"
	defer os.RemoveAll(tmpDir)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	srcPath := filepath.Join(tmpDir, "script.sh")
	dstPath := filepath.Join(tmpDir, "copied_script.sh")

	// Create executable script
	content := "#!/bin/bash\necho 'Hello'\n"
	if err := os.WriteFile(srcPath, []byte(content), 0755); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy file
	if err := CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Check destination is executable
	dstInfo, err := os.Stat(dstPath)
	if err != nil {
		t.Fatalf("Failed to stat destination: %v", err)
	}

	// Check if owner execute bit is set
	if dstInfo.Mode().Perm()&0100 == 0 {
		t.Errorf("Executable bit not preserved. Destination permissions: %o", dstInfo.Mode().Perm())
	}

	// Verify full permissions match (0755)
	if dstInfo.Mode().Perm() != 0755 {
		t.Errorf("Expected permissions 0755, got %o", dstInfo.Mode().Perm())
	}
}

func Test_EmptyDir(t *testing.T) {
	tmpDir := namespace + "test_empty_dir"
	defer os.RemoveAll(tmpDir)

	t.Run("empty_existing_directory", func(t *testing.T) {
		testDir := filepath.Join(tmpDir, "test1")
		os.MkdirAll(testDir, 0755)

		// Create some files and directories
		os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content1"), 0644)
		os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("content2"), 0644)
		os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
		os.WriteFile(filepath.Join(testDir, "subdir", "file3.txt"), []byte("content3"), 0644)

		// Empty the directory
		err := EmptyDir(testDir)
		if err != nil {
			t.Errorf("EmptyDir failed: %v", err)
		}

		// Verify directory still exists but is empty
		entries, err := os.ReadDir(testDir)
		if err != nil {
			t.Errorf("Failed to read directory: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("Expected directory to be empty, but found %d entries", len(entries))
		}
	})

	t.Run("create_nonexistent_directory", func(t *testing.T) {
		testDir := filepath.Join(tmpDir, "nonexistent")

		// Try to empty a non-existent directory - should create it
		err := EmptyDir(testDir)
		if err != nil {
			t.Errorf("EmptyDir failed: %v", err)
		}

		// Verify directory was created
		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			t.Errorf("Directory should have been created")
		}
	})
}

func Test_RemoveFile(t *testing.T) {
	tmpDir := namespace + "test_remove"
	defer os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, 0755)

	t.Run("remove_file", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(testFile, []byte("test"), 0644)

		RemoveFile(testFile)

		if _, err := os.Stat(testFile); !os.IsNotExist(err) {
			t.Errorf("File should have been removed")
		}
	})

	t.Run("remove_directory", func(t *testing.T) {
		testDir := filepath.Join(tmpDir, "testdir")
		os.MkdirAll(testDir, 0755)
		os.WriteFile(filepath.Join(testDir, "file.txt"), []byte("test"), 0644)

		RemoveFile(testDir)

		if _, err := os.Stat(testDir); !os.IsNotExist(err) {
			t.Errorf("Directory should have been removed")
		}
	})
}

func Test_IsModified(t *testing.T) {
	tmpDir := namespace + "test_modified"
	defer os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, 0755)

	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")

	// Test with identical files
	os.WriteFile(file1, []byte("same content"), 0644)
	os.WriteFile(file2, []byte("same content"), 0644)

	modified, err := IsModified(file1, file2)
	if err != nil {
		t.Errorf("IsModified failed: %v", err)
	}
	if modified {
		t.Errorf("Expected files to be identical")
	}

	// Test with different files
	os.WriteFile(file2, []byte("different content"), 0644)
	modified, err = IsModified(file1, file2)
	if err != nil {
		t.Errorf("IsModified failed: %v", err)
	}
	if !modified {
		t.Errorf("Expected files to be different")
	}

	// Test with non-existent file
	_, err = IsModified(file1, filepath.Join(tmpDir, "nonexistent.txt"))
	if err == nil {
		t.Errorf("Expected error for non-existent file")
	}
}
