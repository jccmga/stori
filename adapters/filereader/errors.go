package filereader

import "errors"

var ErrFileNotFound = errors.New("file not found")
var ErrInvalidFile = errors.New("invalid file")
var ErrInvalidHeader = errors.New("invalid header")
var ErrFileIsEmpty = errors.New("file is empty")
var ErrInvalidAmount = errors.New("invalid amount")
var ErrInvalidDateFormat = errors.New("invalid date format")
var ErrInvalidURI = errors.New("invalid URI")
var ErrS3Connection = errors.New("error connecting to S3")
var ErrFileCreation = errors.New("error creating temporary file")
var ErrDownloadFile = errors.New("error downloading file")
