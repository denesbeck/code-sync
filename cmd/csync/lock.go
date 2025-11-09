package main

import (
	"errors"
	"os"
	"time"
)

const DefaultLockTimeout = 5 * time.Second

type Lock struct {
	path     string
	lockFile *os.File
}

func NewLock(path string) *Lock {
	return &Lock{
		path: path + ".lock",
	}
}

func (l *Lock) Acquire(timeout time.Duration) error {
	Debug("Attempting to acquire lock: %s", l.path)
	deadline := time.Now().Add(timeout)

	for {
		/*
			- os.O_CREATE - create the file if it doesn't exist
			- os.O_EXCL - exclusive creation - fail if file already exists
			- os.O_WRONLY - open for write-only access
			- Combined effect: atomically create the file ONLY if it doesn't exist
		*/
		lockFile, err := os.OpenFile(l.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			// Acquire lock, write PID to lock file (for debugging)
			l.lockFile = lockFile
			lockFile.WriteString(string(rune(os.Getpid())))
			Debug("Lock acquired: %s", l.path)
			return nil
		}

		if time.Now().After(deadline) {
			Debug("Lock acquisition timout: %s", l.path)
			return errors.New("Lock acquisition timeout")
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (l *Lock) Release() error {
	if l.lockFile == nil {
		Debug("No lock to release: %s", l.path)
		return nil
	}

	Debug("Releasing lock: %s", l.path)

	// Close lock file
	if err := l.lockFile.Close(); err != nil {
		Debug("Error closing lock file: %v", err)
		return err
	}

	// Remove lock file
	if err := os.Remove(l.path); err != nil {
		Debug("Error removing lock file: %v", err)
		return err
	}

	l.lockFile = nil
	Debug("Lock released: %s", l.path)
	return nil
}

func WithLock(lockPath string, timeout time.Duration, fn func() error) error {
	lock := NewLock(lockPath)

	if err := lock.Acquire(timeout); err != nil {
		return err
	}

	defer func() {
		if err := lock.Release(); err != nil {
			Debug("Warning: failed to release lock: %v", err)
		}
	}()

	return fn()
}
