package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	testPath := "__test__/test.lock"
	
	// Clean up any existing lock file
	os.Remove(testPath)
	defer os.Remove(testPath)
	
	// Create a lock
	lock := NewLock("__test__/test")
	
	fmt.Println("Acquiring lock...")
	if err := lock.Acquire(5 * time.Second); err != nil {
		fmt.Printf("Failed to acquire lock: %v\n", err)
		return
	}
	
	fmt.Println("✅ Lock acquired!")
	
	// Read the lock file to see what's inside
	content, err := os.ReadFile(testPath)
	if err != nil {
		fmt.Printf("❌ Failed to read lock file: %v\n", err)
		return
	}
	
	fmt.Printf("\n📄 Lock file contents:\n")
	fmt.Printf("   Raw bytes: %v\n", content)
	fmt.Printf("   As string: %q\n", string(content))
	fmt.Printf("   Expected PID: %d\n", os.Getpid())
	
	// Verify it's the correct PID
	expectedContent := fmt.Sprintf("%d\n", os.Getpid())
	if string(content) == expectedContent {
		fmt.Println("\n✅ SUCCESS: Lock file contains correct PID!")
	} else {
		fmt.Printf("\n❌ FAIL: Expected %q but got %q\n", expectedContent, string(content))
	}
	
	// Release the lock
	fmt.Println("\nReleasing lock...")
	if err := lock.Release(); err != nil {
		fmt.Printf("Failed to release lock: %v\n", err)
		return
	}
	
	fmt.Println("✅ Lock released!")
}
