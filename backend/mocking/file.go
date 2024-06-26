package mocking

import (
	"errors"
	"io"
	"math"
	"os"
	"time"
)

type File interface {
	Close() error
	io.Closer
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Seek(offset int64, _ int) (int64, error)
	Stat() (os.FileInfo, error)
	Copy(source File, destination File) (int64, error)
}

// FileWrapper is a wrapper for the os.File type that implements the File interface.
//
// Fields:
//   - file (*os.File): The underlying os.File that the FileWrapper wraps.
//   - Buffer ([]byte): A byte slice that serves as the buffer for reading and writing data.
//   - Err (error): An error that can be set to simulate an error condition.
//   - Writer (*os.File): A pointer to the underlying os.File that is used for writing data.
//   - Reader (*os.File): A pointer to the underlying os.File that is used for reading data.
type FileWrapper struct {
	file   *os.File
	Buffer []byte
	Err    error
	Writer *os.File
	Reader *os.File
}

// Wrap creates a new FileWrapper struct
//
// Parameters:
//   - file (*os.File): The underlying os.File to wrap.
//
// Returns:
//   - A pointer to a new FileWrapper struct.
func Wrap(file *os.File) *FileWrapper {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil
	}
	fileSize := fileInfo.Size()
	return &FileWrapper{
		file:   file,
		Buffer: make([]byte, int64(math.Round(float64(fileSize)*1.1))),
		Writer: file,
		Reader: file,
	}
}

// Close is a method of the FileWrapper struct that closes the underlying os.File.
//
// It calls the Close method of the os.File that the FileWrapper is wrapping.
//
// Parameters: None.
//
// Returns:
//   - An error if the underlying os.File cannot be closed. If the file is closed successfully, it returns nil.
func (f *FileWrapper) Close() error {
	return f.file.Close()
}

// Read is a method of the FileWrapper struct that reads from the underlying os.File.
//
// It calls the Read method of the os.File that the FileWrapper is wrapping.
//
// Parameters:
//   - p ([]byte): A byte slice that serves as the buffer into which the data is read.
//
// Returns:
//   - The number of bytes read. If the file is at the end, it returns 0.
//   - An error if the underlying os.File cannot be read. If the file is read successfully, it returns nil.
func (f *FileWrapper) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

// Seek is a method of the FileWrapper struct that sets the offset for the next Read or Write operation on the underlying os.File.
//
// It calls the Seek method of the os.File that the FileWrapper is wrapping.
//
// Parameters:
//   - offset (int64): The offset for the next Read or Write operation. The interpretation of the offset is determined by the whence parameter.
//   - whence (int): The point relative to which the offset is interpreted. It can be 0 (relative to the start of the file), 1 (relative to the current offset), or 2 (relative to the end of the file).
//
// Returns:
//   - The new offset relative to the start of the file.
//   - An error if the underlying os.File cannot seek to the specified offset. If the seek operation is successful, it returns nil.func (f *FileWrapper) Seek(offset int64, whence int) (int64, error) {
func (f *FileWrapper) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

// Write is a method of the FileWrapper struct that writes to the underlying os.File.
//
// It calls the Write method of the os.File that the FileWrapper is wrapping.
//
// Parameters:
//   - p ([]byte): A byte slice that serves as the buffer from which the data is written.
//
// Returns:
//   - The number of bytes written. If the file is at the end, it returns 0.
//   - An error if the underlying os.File cannot be written to. If the file is written successfully, it returns nil.
func (f *FileWrapper) Write(p []byte) (int, error) {
	return f.file.Write(p)
}

// Copy is a method of the FileWrapper struct that copies data from a source File to a destination File.
//
// It reads from the source File into the FileWrapper's Buffer, then writes from the Buffer to the destination File.
//
// Parameters:
//   - source (File): The source File to copy data from.
//   - destination (File): The destination File to copy data to.
//
// Returns:
//   - The number of bytes written to the destination File.
//   - An error if the source File cannot be read or the destination File cannot be written to. If the copy operation is successful, it returns nil.
func (f *FileWrapper) Copy(source File, destination File) (int64, error) {
	// Read from the source file
	bytesRead, err := source.Read(f.Buffer)
	if err != nil {
		return 0, err
	}

	// Write to the destination file
	bytesWritten, err := destination.Write(f.Buffer[:bytesRead])
	if err != nil {
		return 0, err
	}

	// Return the number of bytes written
	return int64(bytesWritten), nil
}

// Stat is a method of the FileWrapper struct that retrieves the file descriptor's metadata.
//
// It calls the Stat method of the os.File that the FileWrapper is wrapping.
//
// Parameters: None.
//
// Returns:
//   - An os.FileInfo object that describes the file. If the method is successful, it returns this object and a nil error.
//   - An error if the underlying os.File's metadata cannot be retrieved. If the method is unsuccessful, it returns a nil os.FileInfo and the error.
func (f *FileWrapper) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}

// FileInfoMock is a struct that mocks the os.FileInfo interface for testing purposes.
//
// It contains a single field, file, which is a pointer to a FileMock. This allows the FileInfoMock to return
// information about the mocked file when its methods are called.
//
// Fields:
//   - file (*FileMock): A pointer to a FileMock struct.
type FileInfoMock struct {
	file *FileMock
}

// Size returns the length of the buffer in the FileMock that FileInfoMock is associated with.
func (f *FileInfoMock) Size() int64 {
	return int64(len(f.file.Buffer))
}

// Name returns an empty string as it's not relevant in this mock implementation.
func (f *FileInfoMock) Name() string { return "" }

// Mode returns 0 as it's not relevant in this mock implementation.
func (f *FileInfoMock) Mode() os.FileMode { return 0 }

// ModTime returns the zero value for time.Time as it's not relevant in this mock implementation.
func (f *FileInfoMock) ModTime() time.Time { return time.Time{} }

// IsDir returns false as it's not relevant in this mock implementation.
func (f *FileInfoMock) IsDir() bool { return false }

// Sys returns nil as it's not relevant in this mock implementation.
func (f *FileInfoMock) Sys() interface{} { return nil }

// FileMock is a struct that mocks the File interface for testing purposes.
//
// Fields:
//   - FileName (string): The name of the file.
//   - IsOpen (bool): A boolean indicating whether the file is open or not.
//   - Buffer ([]byte): A byte slice that serves as the buffer for the file data.
//   - Bytes (int): The number of bytes in the buffer.
//   - Err (error): An error that can be set to simulate an error condition.
//   - FileInfo (*FileInfoMock): A pointer to a FileInfoMock struct. This allows the FileMock to return information about the mocked file when its methods are called.
type FileMock struct {
	FileName string
	IsOpen   bool
	Buffer   []byte
	Bytes    int
	Err      error
	FileInfo *FileInfoMock
}

// ReadDir is a method of the FileMock struct that simulates the behavior of reading a directory.
// It doesn't actually read a directory, but instead returns a predefined result that can be set for testing purposes.
//
// Parameters:
//   - _ (string): This method ignores its input parameter. The underscore character is a convention in Go for discarding a variable.
//
// Returns:
//   - A slice of os.DirEntry: In this mock implementation, it always returns nil.
//   - An error: The error that was previously set in the FileMock struct. If no error was set, it returns nil.
func (f *FileMock) ReadDir(_ string) ([]os.DirEntry, error) {
	return nil, f.Err
}

// Close is a method of the FileMock struct that simulates the behavior of closing a file.
//
// It checks if the FileMock is nil or if the file is already closed, and returns an appropriate error in each case.
// If the file is open, it sets the IsOpen field to false, simulating the closing of the file.
//
// Parameters: None.
//
// Returns:
//   - os.ErrInvalid if the FileMock is nil.
//   - os.ErrClosed if the file is already closed.
//   - The error that was previously set in the FileMock struct. If no error was set and the file is closed successfully, it returns nil.
func (f *FileMock) Close() error {
	if f == nil {
		return os.ErrInvalid
	}
	if f.IsOpen {
		f.IsOpen = false
	} else {
		return os.ErrClosed
	}
	return f.Err
}

// Read is a method of the FileMock struct that simulates the behavior of reading from a file.
//
// It checks if there is an error set in the FileMock. If there is, it returns 0 and the error.
// If there is no error, it copies data from the FileMock's Buffer into the provided byte slice and updates the Buffer.
// If the Buffer is empty, it returns 0 and io.EOF to simulate the end of the file.
//
// Parameters:
//   - p ([]byte): A byte slice that serves as the buffer into which the data is read.
//
// Returns:
//   - The number of bytes read. If the Buffer is empty, it returns 0.
//   - An error if one was set in the FileMock. If the Buffer is empty, it returns io.EOF. If the read operation is successful, it returns nil.
func (f *FileMock) Read(p []byte) (int, error) {
	if f.Err != nil {
		return 0, f.Err
	}

	if len(f.Buffer) == 0 {
		return 0, io.EOF
	}

	n := copy(p, f.Buffer)
	f.Buffer = f.Buffer[n:]
	return n, nil
}

// Write is a method of the FileMock struct that simulates the behavior of writing to a file.
//
// It ignores its input parameter and instead returns the number of bytes that were previously set in the FileMock struct.
// This allows you to control the behavior of the Write method for testing purposes.
//
// Parameters:
//   - _ ([]byte): A byte slice that serves as the buffer from which the data is written. In this mock implementation, this parameter is ignored.
//
// Returns:
//   - The number of bytes that were previously set in the FileMock struct. If no number was set, it returns 0.
//   - An error if one was set in the FileMock. If no error was set, it returns nil.
func (f *FileMock) Write(_ []byte) (int, error) {
	return f.Bytes, f.Err
}

// Seek is a method of the FileMock struct that simulates the behavior of setting the offset for the next Read or Write operation on a file.
//
// It checks if there is an error set in the FileMock. If there is, it returns 0 and the error.
// If there is no error, it adjusts the Buffer according to the offset and whence parameters.
// If the offset is out of range, it returns 0 and io.EOF to simulate the end of the file.
//
// Parameters:
//   - offset (int64): The offset for the next Read or Write operation. The interpretation of the offset is determined by the whence parameter.
//   - whence (int): The point relative to which the offset is interpreted. It can be 0 (relative to the start of the file), 1 (relative to the current offset), or 2 (relative to the end of the file).
//
// Returns:
//   - The new offset relative to the start of the file.
//   - An error if one was set in the FileMock. If the offset is out of range, it returns io.EOF. If the seek operation is successful, it returns nil.
func (f *FileMock) Seek(offset int64, whence int) (int64, error) {
	if f.Err != nil {
		return 0, f.Err
	}

	switch whence {
	case 0: // relative to the origin of the file
		if offset < 0 || offset > int64(len(f.Buffer)) {
			return 0, io.EOF
		}
		f.Buffer = f.Buffer[offset:]
	case 1: // relative to the current offset
		if offset < 0 || offset > int64(len(f.Buffer)) {
			return 0, io.EOF
		}
		f.Buffer = f.Buffer[offset:]
	case 2: // relative to the end
		if offset > 0 || -offset > int64(len(f.Buffer)) {
			return 0, io.EOF
		}
		f.Buffer = f.Buffer[:len(f.Buffer)+int(offset)]
	default:
		return 0, errors.New("invalid whence")
	}

	return int64(len(f.Buffer)), nil
}

// Copy is a method of the FileMock struct that simulates the behavior of copying data from a source File to a destination File.
//
// It reads from the source File and writes to the destination File. The actual data transfer is not simulated in this mock implementation.
// Instead, it returns the number of bytes that were previously set in the FileMock struct and the error that was previously set, if any.
// This allows you to control the behavior of the Copy method for testing purposes.
//
// Parameters:
//   - source (File): The source File to copy data from. In this mock implementation, this parameter is used but no data is actually read from it.
//   - destination (File): The destination File to copy data to. In this mock implementation, this parameter is used but no data is actually written to it.
//
// Returns:
//   - The number of bytes that were previously set in the FileMock struct. If no number was set, it returns 0.
//   - An error if one was set in the FileMock. If no error was set, it returns nil.
func (f *FileMock) Copy(source File, destination File) (int64, error) {
	var err error
	_, err = source.Read([]byte{})
	if err != nil {
		return 0, err
	}
	_, err = destination.Read([]byte{})
	if err != nil {
		return 0, err
	}
	return int64(f.Bytes), f.Err
}

// Stat is a method of the FileMock struct that simulates the behavior of retrieving the file descriptor's metadata.
//
// It returns the FileInfo and error that were previously set in the FileMock struct. This allows you to control the behavior of the Stat method for testing purposes.
//
// Parameters: None.
//
// Returns:
//   - An os.FileInfo object that describes the file. If no error was set and the method is successful, it returns this object and a nil error.
//   - An error if one was set in the FileMock. If no error was set, it returns nil.
func (f *FileMock) Stat() (os.FileInfo, error) {
	return f.FileInfo, f.Err
}
