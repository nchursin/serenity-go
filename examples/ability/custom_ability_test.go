package examples

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nchursin/serenity-go/serenity/abilities"
	"github.com/nchursin/serenity-go/serenity/abilities/api"
	"github.com/nchursin/serenity-go/serenity/core"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// FileSystemAbility enables an actor to interact with the file system
type FileSystemAbility interface {
	abilities.Ability
	// Core file operations
	ReadFile(path string) (string, error)
	WriteFile(path string, content string) error
	DeleteFile(path string) error
	CreateDirectory(path string) error
	Exists(path string) bool

	// Directory operations
	ListFiles(dir string) ([]string, error)
	WorkingDirectory() string
	SetWorkingDirectory(dir string) error

	// State tracking
	LastOperation() string
	OperationCount() int
}

// fileSystemAbility implements FileSystemAbility
type fileSystemAbility struct {
	workingDir    string
	lastOperation string
	opCount       int
	mutex         sync.RWMutex
}

// ManageFiles creates a new FileSystemAbility with default working directory
func ManageFiles() FileSystemAbility {
	return &fileSystemAbility{
		workingDir:    ".",
		lastOperation: "none",
		opCount:       0,
	}
}

// ManageFilesIn creates a new FileSystemAbility with specified working directory
func ManageFilesIn(directory string) FileSystemAbility {
	if !filepath.IsAbs(directory) {
		abs, err := filepath.Abs(directory)
		if err != nil {
			directory = "."
		} else {
			directory = abs
		}
	}

	return &fileSystemAbility{
		workingDir:    directory,
		lastOperation: "none",
		opCount:       0,
	}
}

func (f *fileSystemAbility) ReadFile(path string) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	fullPath := filepath.Join(f.workingDir, path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		f.lastOperation = fmt.Sprintf("read error: %s", path)
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	f.lastOperation = fmt.Sprintf("read: %s", path)
	f.opCount++
	return string(content), nil
}

func (f *fileSystemAbility) WriteFile(path string, content string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	fullPath := filepath.Join(f.workingDir, path)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		f.lastOperation = fmt.Sprintf("write error (mkdir): %s", path)
		return fmt.Errorf("failed to create directory for %s: %w", path, err)
	}

	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		f.lastOperation = fmt.Sprintf("write error: %s", path)
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	f.lastOperation = fmt.Sprintf("write: %s", path)
	f.opCount++
	return nil
}

func (f *fileSystemAbility) DeleteFile(path string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	fullPath := filepath.Join(f.workingDir, path)
	err := os.Remove(fullPath)
	if err != nil {
		f.lastOperation = fmt.Sprintf("delete error: %s", path)
		return fmt.Errorf("failed to delete file %s: %w", path, err)
	}

	f.lastOperation = fmt.Sprintf("delete: %s", path)
	f.opCount++
	return nil
}

func (f *fileSystemAbility) CreateDirectory(path string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	fullPath := filepath.Join(f.workingDir, path)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		f.lastOperation = fmt.Sprintf("mkdir error: %s", path)
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	f.lastOperation = fmt.Sprintf("mkdir: %s", path)
	f.opCount++
	return nil
}

func (f *fileSystemAbility) Exists(path string) bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	fullPath := filepath.Join(f.workingDir, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (f *fileSystemAbility) ListFiles(dir string) ([]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	fullPath := filepath.Join(f.workingDir, dir)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		f.lastOperation = fmt.Sprintf("list error: %s", dir)
		return nil, fmt.Errorf("failed to list directory %s: %w", dir, err)
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	f.lastOperation = fmt.Sprintf("list: %s (%d files)", dir, len(files))
	f.opCount++
	return files, nil
}

func (f *fileSystemAbility) WorkingDirectory() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.workingDir
}

func (f *fileSystemAbility) SetWorkingDirectory(dir string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !filepath.IsAbs(dir) {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		dir = abs
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	f.workingDir = dir
	return nil
}

func (f *fileSystemAbility) LastOperation() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.lastOperation
}

func (f *fileSystemAbility) OperationCount() int {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.opCount
}

// Activities for the FileSystemAbility

// ReadFileActivity represents an activity to read a file
type ReadFileActivity struct {
	path string
}

func ReadFile(path string) *ReadFileActivity {
	return &ReadFileActivity{path: path}
}

func (r *ReadFileActivity) PerformAs(actor core.Actor) error {
	ability, err := actor.AbilityTo(&fileSystemAbility{})
	if err != nil {
		return fmt.Errorf("actor does not have file system ability: %w", err)
	}

	fileManager := ability.(FileSystemAbility)
	_, err = fileManager.ReadFile(r.path)
	return err
}

func (r *ReadFileActivity) Description() string {
	return fmt.Sprintf("reads file: %s", r.path)
}

// WriteFileActivity represents an activity to write a file
type WriteFileActivity struct {
	path    string
	content string
}

func WriteFile(path string, content string) *WriteFileActivity {
	return &WriteFileActivity{path: path, content: content}
}

func (w *WriteFileActivity) PerformAs(actor core.Actor) error {
	ability, err := actor.AbilityTo(&fileSystemAbility{})
	if err != nil {
		return fmt.Errorf("actor does not have file system ability: %w", err)
	}

	fileManager := ability.(FileSystemAbility)
	return fileManager.WriteFile(w.path, w.content)
}

// FailureMode returns the failure mode for send requests (default: FailFast)
func (s *WriteFileActivity) FailureMode() core.FailureMode {
	return core.FailFast
}

func (w *WriteFileActivity) Description() string {
	return fmt.Sprintf("writes file: %s", w.path)
}

// DeleteFileActivity represents an activity to delete a file
type DeleteFileActivity struct {
	path string
}

func DeleteFile(path string) *DeleteFileActivity {
	return &DeleteFileActivity{path: path}
}

func (d *DeleteFileActivity) PerformAs(actor core.Actor) error {
	ability, err := actor.AbilityTo(&fileSystemAbility{})
	if err != nil {
		return fmt.Errorf("actor does not have file system ability: %w", err)
	}

	fileManager := ability.(FileSystemAbility)
	return fileManager.DeleteFile(d.path)
}

// FailureMode returns the failure mode for send requests (default: FailFast)
func (d *DeleteFileActivity) FailureMode() core.FailureMode {
	return core.FailFast
}

func (d *DeleteFileActivity) Description() string {
	return fmt.Sprintf("deletes file: %s", d.path)
}

// Questions for the FileSystemAbility

// FileContentQuestion asks about the content of a file
type FileContentQuestion struct {
	path string
}

func FileContent(path string) *FileContentQuestion {
	return &FileContentQuestion{path: path}
}

func (f *FileContentQuestion) AnsweredBy(actor core.Actor) (string, error) {
	ability, err := actor.AbilityTo(&fileSystemAbility{})
	if err != nil {
		return "", fmt.Errorf("actor does not have file system ability: %w", err)
	}

	fileManager := ability.(FileSystemAbility)
	return fileManager.ReadFile(f.path)
}

func (f *FileContentQuestion) Description() string {
	return fmt.Sprintf("content of file: %s", f.path)
}

// FileExistsQuestion asks whether a file exists
type FileExistsQuestion struct {
	path string
}

func FileExists(path string) *FileExistsQuestion {
	return &FileExistsQuestion{path: path}
}

func (f *FileExistsQuestion) AnsweredBy(actor core.Actor) (bool, error) {
	ability, err := actor.AbilityTo(&fileSystemAbility{})
	if err != nil {
		return false, fmt.Errorf("actor does not have file system ability: %w", err)
	}

	fileManager := ability.(FileSystemAbility)
	return fileManager.Exists(f.path), nil
}

func (f *FileExistsQuestion) Description() string {
	return fmt.Sprintf("existence of file: %s", f.path)
}

// Tests for FileSystemAbility

func TestFileSystemAbility_BasicOperations(t *testing.T) {
	test := serenity.NewSerenityTest(t)
	defer test.Shutdown()
	tempDir := t.TempDir()
	actor := test.ActorCalled("FileTester").WhoCan(ManageFilesIn(tempDir))

	// Test writing and reading a file
	testContent := "Hello, World!"
	actor.AttemptsTo(
		WriteFile("test.txt", testContent),
	)

	content, err := FileContent("test.txt").AnsweredBy(actor)
	require.NoError(t, err)
	require.Equal(t, testContent, content)

	// Test file existence
	exists, err := FileExists("test.txt").AnsweredBy(actor)
	require.NoError(t, err)
	require.True(t, exists)

	// Test non-existent file
	exists, err = FileExists("nonexistent.txt").AnsweredBy(actor)
	require.NoError(t, err)
	require.False(t, exists)

	// Test deleting file
	actor.AttemptsTo(
		DeleteFile("test.txt"),
	)
	require.NoError(t, err)

	exists, err = FileExists("test.txt").AnsweredBy(actor)
	require.NoError(t, err)
	require.False(t, exists)
}

func TestFileSystemAbility_DirectoryOperations(t *testing.T) {
	test := serenity.NewSerenityTest(t)
	defer test.Shutdown()
	tempDir := t.TempDir()
	actor := test.ActorCalled("DirectoryTester").WhoCan(ManageFilesIn(tempDir))

	// Test creating directory
	ability, err := actor.AbilityTo(&fileSystemAbility{})
	require.NoError(t, err)
	fileManager := ability.(FileSystemAbility)
	err = fileManager.CreateDirectory("testdir")
	require.NoError(t, err)
	require.NoError(t, err)

	// Check directory exists by creating a file inside it
	actor.AttemptsTo(
		WriteFile("testdir/nested.txt", "nested content"),
	)
	require.NoError(t, err)

	exists, err := FileExists("testdir/nested.txt").AnsweredBy(actor)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestFileSystemAbility_WorkingDirectory(t *testing.T) {
	tempDir := t.TempDir()
	ability := ManageFilesIn(tempDir)

	// Test initial working directory
	require.Equal(t, tempDir, ability.WorkingDirectory())

	// Test setting working directory to existing directory
	subDir := filepath.Join(tempDir, "subdir")
	require.NoError(t, os.Mkdir(subDir, 0755))
	require.NoError(t, ability.SetWorkingDirectory(subDir))
	require.Equal(t, subDir, ability.WorkingDirectory())

	// Test operations in new working directory
	err := ability.WriteFile("test.txt", "test content")
	require.NoError(t, err)
	require.True(t, ability.Exists("test.txt"))
}

func TestFileSystemAbility_StateTracking(t *testing.T) {
	tempDir := t.TempDir()
	ability := ManageFilesIn(tempDir)

	// Initial state
	require.Equal(t, "none", ability.LastOperation())
	require.Equal(t, 0, ability.OperationCount())

	// After write operation
	require.NoError(t, ability.WriteFile("test.txt", "content"))
	require.Contains(t, ability.LastOperation(), "write: test.txt")
	require.Equal(t, 1, ability.OperationCount())

	// After read operation
	_, err := ability.ReadFile("test.txt")
	require.NoError(t, err)
	require.Contains(t, ability.LastOperation(), "read: test.txt")
	require.Equal(t, 2, ability.OperationCount())

	// After delete operation
	require.NoError(t, ability.DeleteFile("test.txt"))
	require.Contains(t, ability.LastOperation(), "delete: test.txt")
	require.Equal(t, 3, ability.OperationCount())
}

func TestFileSystemAbility_ErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	ability := ManageFilesIn(tempDir)

	// Test reading non-existent file
	_, err := ability.ReadFile("nonexistent.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read file")
	require.Contains(t, ability.LastOperation(), "read error: nonexistent.txt")

	// Test deleting non-existent file
	err = ability.DeleteFile("nonexistent.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to delete file")
	require.Contains(t, ability.LastOperation(), "delete error: nonexistent.txt")

	// Test setting non-existent working directory
	err = ability.SetWorkingDirectory("/nonexistent/directory")
	require.Error(t, err)
	require.Contains(t, err.Error(), "directory does not exist")
}

func TestFileSystemAbility_ListFiles(t *testing.T) {
	tempDir := t.TempDir()
	ability := ManageFilesIn(tempDir)

	// Create some files
	require.NoError(t, ability.WriteFile("file1.txt", "content1"))
	require.NoError(t, ability.WriteFile("file2.txt", "content2"))
	require.NoError(t, ability.CreateDirectory("subdir"))
	require.NoError(t, ability.WriteFile("subdir/file3.txt", "content3"))

	// List root directory
	files, err := ability.ListFiles(".")
	require.NoError(t, err)
	require.Len(t, files, 3) // file1.txt, file2.txt, subdir
	require.Contains(t, files, "file1.txt")
	require.Contains(t, files, "file2.txt")
	require.Contains(t, files, "subdir")

	// List subdirectory
	files, err = ability.ListFiles("subdir")
	require.NoError(t, err)
	require.Len(t, files, 1)
	require.Contains(t, files, "file3.txt")
}

func TestFileSystemAbility_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	ability := ManageFilesIn(tempDir)

	// Test concurrent access to the ability
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			filename := fmt.Sprintf("file%d.txt", i)
			content := fmt.Sprintf("content%d", i)

			// Write file
			require.NoError(t, ability.WriteFile(filename, content))

			// Read file back
			readContent, err := ability.ReadFile(filename)
			require.NoError(t, err)
			require.Equal(t, content, readContent)

			// Check file exists
			require.True(t, ability.Exists(filename))
		}(i)
	}

	wg.Wait()

	// Verify all files exist and count operations
	require.Equal(t, numGoroutines*2, ability.OperationCount()) // write + read per goroutine
}

// Integration test showing FileSystemAbility working with other abilities
func TestFileSystemAbility_WithAPIIntegration(t *testing.T) {
	test := serenity.NewSerenityTest(t)
	defer test.Shutdown()
	tempDir := t.TempDir()
	actor := test.ActorCalled("IntegrationTester").WhoCan(
		ManageFilesIn(tempDir),
		api.CallAnApiAt("https://jsonplaceholder.typicode.com"),
	)

	// Get data from API
	actor.AttemptsTo(
		api.SendGetRequest("/posts/1"),
	)

	// Save API response to file
	responseBody, err := api.LastResponseBody{}.AnsweredBy(actor)
	require.NoError(t, err)

	actor.AttemptsTo(
		WriteFile("post.json", responseBody),
	)

	// Verify file was created and contains expected data
	fileContent, err := FileContent("post.json").AnsweredBy(actor)
	require.NoError(t, err)
	require.Contains(t, fileContent, "sunt aut facere")
	require.Contains(t, fileContent, "quia et suscipit")
}
