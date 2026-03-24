package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

// ExitCode defines standard exit codes per CLI spec
type ExitCode int

const (
	ExitSuccess             ExitCode = 0
	ExitGeneralFailure      ExitCode = 1
	ExitUsageError          ExitCode = 2
	ExitNetworkError        ExitCode = 3
	ExitChecksumFailed      ExitCode = 4
	ExitInstallFailed       ExitCode = 5
	ExitManifestValidation  ExitCode = 6
	ExitPackageNotFound     ExitCode = 7
	ExitPkgdbError          ExitCode = 8
	ExitDryRunSuccess       ExitCode = 9
)

// ErrorResponse is the JSON error format per spec
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// GPkgError wraps an error with an exit code
type GPkgError struct {
	Code    ExitCode
	Message string
	Detail  string
	Err     error
}

// Error implements the error interface
func (e *GPkgError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

// NewGPkgError creates a new GPkgError
func NewGPkgError(code ExitCode, message string, err error) *GPkgError {
	return &GPkgError{
		Code:    code,
		Message: message,
		Err:     err,
		Detail:  err.Error(),
	}
}

// ExitWithError handles error scenarios with proper exit codes
func ExitWithError(err error) {
	if err == nil {
		os.Exit(int(ExitSuccess))
	}

	var gpkgErr *GPkgError
	if e, ok := err.(*GPkgError); ok {
		gpkgErr = e
	} else {
		gpkgErr = NewGPkgError(ExitGeneralFailure, "Error", err)
	}

	if jsonOutput {
		errResp := ErrorResponse{
			Error: ErrorDetail{
				Code:    int(gpkgErr.Code),
				Message: gpkgErr.Message,
				Detail:  gpkgErr.Detail,
			},
		}
		json.NewEncoder(os.Stderr).Encode(errResp)
	} else {
		if !quiet {
			fmt.Fprintf(os.Stderr, "error: %v\n", gpkgErr)
		}
	}

	os.Exit(int(gpkgErr.Code))
}

// ErrorChecksumMismatch creates a checksum failure error
func ErrorChecksumMismatch(detail string) *GPkgError {
	return &GPkgError{
		Code:    ExitChecksumFailed,
		Message: "Checksum validation failed",
		Detail:  detail,
	}
}

// ErrorNetworkFailure creates a network error
func ErrorNetworkFailure(err error) *GPkgError {
	return &GPkgError{
		Code:    ExitNetworkError,
		Message: "Network error",
		Err:     err,
		Detail:  err.Error(),
	}
}

// ErrorManifestInvalid creates a manifest validation error
func ErrorManifestInvalid(detail string) *GPkgError {
	return &GPkgError{
		Code:    ExitManifestValidation,
		Message: "Manifest validation failed",
		Detail:  detail,
	}
}

// ErrorPackageNotFound creates a package not found error
func ErrorPackageNotFound(pkg string) *GPkgError {
	return &GPkgError{
		Code:    ExitPackageNotFound,
		Message: "Package not found",
		Detail:  pkg,
	}
}

// ErrorInstallFailed creates an install failure error
func ErrorInstallFailed(detail string) *GPkgError {
	return &GPkgError{
		Code:    ExitInstallFailed,
		Message: "Installation failed",
		Detail:  detail,
	}
}

// ErrorPkgdbFailure creates a pkgdb error
func ErrorPkgdbFailure(err error) *GPkgError {
	return &GPkgError{
		Code:    ExitPkgdbError,
		Message: "Package database error",
		Err:     err,
		Detail:  err.Error(),
	}
}
