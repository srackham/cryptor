package fsx

import (
	"reflect"
	"testing"
	"time"
)

func TestDirExists(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		args   args
		wanted bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DirExists(tt.args.name); got != tt.wanted {
				t.Errorf("DirExists() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		args   args
		wanted bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.name); got != tt.wanted {
				t.Errorf("FileExists() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wanted  string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wanted {
				t.Errorf("ReadFile() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	type args struct {
		name string
		text string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteFile(tt.args.name, tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWritePath(t *testing.T) {
	type args struct {
		path string
		text string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WritePath(tt.args.path, tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("WritePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		args   args
		wanted string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileName(tt.args.name); got != tt.wanted {
				t.Errorf("FileName() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}

func TestReplaceExt(t *testing.T) {
	type args struct {
		name string
		ext  string
	}
	tests := []struct {
		name   string
		args   args
		wanted string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceExt(tt.args.name, tt.args.ext); got != tt.wanted {
				t.Errorf("ReplaceExt() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyFile(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMkMissingDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MkMissingDir(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("MkMissingDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPathIsInDir(t *testing.T) {
	type args struct {
		p   string
		dir string
	}
	tests := []struct {
		name   string
		args   args
		wanted bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PathIsInDir(tt.args.p, tt.args.dir); got != tt.wanted {
				t.Errorf("PathIsInDir() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}

func TestPathTranslate(t *testing.T) {
	type args struct {
		srcPath string
		srcRoot string
		dstRoot string
	}
	tests := []struct {
		name   string
		args   args
		wanted string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PathTranslate(tt.args.srcPath, tt.args.srcRoot, tt.args.dstRoot); got != tt.wanted {
				t.Errorf("PathTranslate() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}

func TestFileModTime(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name   string
		args   args
		wanted time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileModTime(tt.args.f); !reflect.DeepEqual(got, tt.wanted) {
				t.Errorf("FileModTime() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}

func TestDirCount(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name   string
		args   args
		wanted int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DirCount(tt.args.dir); got != tt.wanted {
				t.Errorf("DirCount() = %v, wanted %v", got, tt.wanted)
			}
		})
	}
}
