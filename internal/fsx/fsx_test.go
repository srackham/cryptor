package fsx

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDirExists(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fsx")
	if err != nil {
		t.Errorf("TempDir failed with error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dir := filepath.Join(tempDir, "test_dir_exists")
	err = os.MkdirAll(dir, 0775)
	if err != nil {
		t.Errorf("MkdirAll failed with error: %v", err)
	}

	exists := DirExists(dir)
	if !exists {
		t.Errorf("DirExists returned false for existing directory")
	}

	notExist := DirExists(filepath.Join(tempDir, "test_dir_not_exist"))
	if notExist {
		t.Errorf("DirExists returned true for non-existing directory")
	}
}

func TestFileExists(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fsx")
	if err != nil {
		t.Errorf("TempDir failed with error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fileName := filepath.Join(tempDir, "test_file_exists")
	err = ioutil.WriteFile(fileName, []byte("Test"), 0644)
	if err != nil {
		t.Errorf("WriteFile failed with error: %v", err)
	}

	exists := FileExists(fileName)
	if !exists {
		t.Errorf("FileExists returned false for existing file")
	}

	notExist := FileExists(filepath.Join(tempDir, "test_file_not_exist"))
	if notExist {
		t.Errorf("FileExists returned true for non-existing file")
	}
}

func TestReadAndWriteFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fsx")
	if err != nil {
		t.Errorf("TempDir failed with error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fileName := filepath.Join(tempDir, "test_read_and_write_file")
	text := "Test"

	err = WriteFile(fileName, text)
	if err != nil {
		t.Errorf("WriteFile failed with error: %v", err)
	}

	readText, err := ReadFile(fileName)
	if err != nil {
		t.Errorf("ReadFile failed with error: %v", err)
	}

	if readText != text {
		t.Errorf("ReadFile did not read the same text as written, got: %s, want: %s", readText, text)
	}
}

func TestWritePath(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fsx")
	if err != nil {
		t.Errorf("TempDir failed with error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dir := filepath.Join(tempDir, "test_write_path")
	err = os.MkdirAll(dir, 0775)
	if err != nil {
		t.Errorf("MkdirAll failed with error: %v", err)
	}

	fileName := filepath.Join(dir, "test_write_path.txt")
	text := "Test"

	err = WritePath(fileName, text)
	if err != nil {
		t.Errorf("WritePath failed with error: %v", err)
	}
}

func TestFileName(t *testing.T) {
	name := "/path/to/file/test_file.txt"
	fileName := FileName(name)
	expected := "test_file"
	if fileName != expected {
		t.Errorf("FileName did not return the expected file name, got: %s, want: %s", fileName, expected)
	}
}

func TestReplaceExt(t *testing.T) {
	name := "/path/to/file/test_file.txt"
	newExt := ".md"
	newName := ReplaceExt(name, newExt)
	expected := "/path/to/file/test_file.md"
	if newName != expected {
		t.Errorf("ReplaceExt did not return the expected file name, got: %s, want: %s", newName, expected)
	}
}

func TestCopyFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fsx")
	if err != nil {
		t.Errorf("TempDir failed with error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fromName := filepath.Join(tempDir, "test_copy_file_from")
	toName := filepath.Join(tempDir, "test_copy_file_to")
	text := "Test"

	err = WriteFile(fromName, text)
	if err != nil {
		t.Errorf("WriteFile failed with error: %v", err)
	}

	err = CopyFile(fromName, toName)
	if err != nil {
		t.Errorf("CopyFile failed with error: %v", err)
	}

	readText, err := ReadFile(toName)
	if err != nil {
		t.Errorf("ReadFile failed with error: %v", err)
	}

	if readText != text {
		t.Errorf("CopyFile did not copy the same text as written, got: %s, want: %s", readText, text)
	}
}

func TestMkMissingDir(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fsx")
	if err != nil {
		t.Errorf("TempDir failed with error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dir := filepath.Join(tempDir, "test_mk_missing_dir/test_subdir")
	err = MkMissingDir(dir)
	if err != nil {
		t.Errorf("MkMissingDir failed with error: %v", err)
	}

	if !DirExists(dir) {
		t.Errorf("MkMissingDir did not create the missing directory")
	}
}

func TestPathIsInDir(t *testing.T) {
	p := "/path/to/file/file.txt"
	dir := "/path/to"
	inDir := PathIsInDir(p, dir)
	if !inDir {
		t.Errorf("PathIsInDir returned false on path that is in directory")
	}

	outDir := PathIsInDir(p, "/path/out")
	if outDir {
		t.Errorf("PathIsInDir returned true on path that is not in directory")
	}
}

func TestPathTranslate(t *testing.T) {
	srcPath := "/path/to/src/file.txt"
	srcRoot := "/path/to"
	dstRoot := "/path/to/dst"
	expected := "/path/to/dst/src/file.txt"
	dstPath := PathTranslate(srcPath, srcRoot, dstRoot)
	if dstPath != expected {
		t.Errorf("PathTranslate did not return the expected destination path, got: %s, want: %s", dstPath, expected)
	}
}

func TestFileModTime(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fsx")
	if err != nil {
		t.Errorf("TempDir failed with error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fileName := filepath.Join(tempDir, "test_file_mod_time")
	err = ioutil.WriteFile(fileName, []byte("Test"), 0644)
	if err != nil {
		t.Errorf("WriteFile failed with error: %v", err)
	}

	modTime := FileModTime(fileName)
	if modTime.IsZero() {
		t.Errorf("FileModTime returned zero time for existing file")
	}

	notExistTime := FileModTime(filepath.Join(tempDir, "test_file_not_exist"))
	if !notExistTime.IsZero() {
		t.Errorf("FileModTime did not return zero time for non-existing file")
	}
}

func TestDirCount(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fsx")
	if err != nil {
		t.Errorf("TempDir failed with error: %v", err)
	}

	dir := filepath.Join(tempDir, "test_dir_count")
	err = os.MkdirAll(dir, 0775)
	if err != nil {
		t.Errorf("MkdirAll failed with error: %v", err)
	}

	fileName := filepath.Join(dir, "test_file")
	err = ioutil.WriteFile(fileName, []byte("Test"), 0644)
	if err != nil {
		t.Errorf("WriteFile failed with error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	count := DirCount(dir)
	expected := 1
	if count != expected {
		t.Errorf("DirCount did not return the expected number of files and folders, got: %d, want: %d", count, expected)
	}

	notExistCount := DirCount(filepath.Join(tempDir, "test_dir_not_exist"))
	if notExistCount != 0 {
		t.Errorf("DirCount did not return zero for non-existing directory, got: %d", notExistCount)
	}
}
