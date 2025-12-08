package main

import (
	"fmt"
	"strings"
)

// ReportError represents different types of errors that can occur during report generation
type ReportError struct {
	Type    ErrorType
	Message string
	Err     error
}

// ErrorType categorizes different types of errors
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeAuth
	ErrorTypeNetwork
	ErrorTypeRateLimit
	ErrorTypePartialData
	ErrorTypeInvalidDateRange
)

// Error implements the error interface
func (re *ReportError) Error() string {
	if re.Err != nil {
		return fmt.Sprintf("%s: %v", re.Message, re.Err)
	}
	return re.Message
}

// NewReportError creates a new ReportError
func NewReportError(errType ErrorType, message string, err error) *ReportError {
	return &ReportError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// CategorizeError analyzes an error and returns a ReportError with appropriate type
func CategorizeError(err error) *ReportError {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	errMsgLower := strings.ToLower(errMsg)

	// Check for authentication errors
	if strings.Contains(errMsgLower, "401") ||
		strings.Contains(errMsgLower, "403") ||
		strings.Contains(errMsgLower, "unauthorized") ||
		strings.Contains(errMsgLower, "forbidden") ||
		strings.Contains(errMsgLower, "authentication") {
		return NewReportError(
			ErrorTypeAuth,
			"Authentication failed. Please check your API credentials in ~/.jirarc",
			err,
		)
	}

	// Check for rate limit errors
	if strings.Contains(errMsgLower, "429") ||
		strings.Contains(errMsgLower, "rate limit") ||
		strings.Contains(errMsgLower, "too many requests") {
		return NewReportError(
			ErrorTypeRateLimit,
			"Jira API rate limit exceeded. Please wait a few minutes and try again.",
			err,
		)
	}

	// Check for network errors
	if strings.Contains(errMsgLower, "connection") ||
		strings.Contains(errMsgLower, "network") ||
		strings.Contains(errMsgLower, "timeout") ||
		strings.Contains(errMsgLower, "dial") ||
		strings.Contains(errMsgLower, "no such host") {
		return NewReportError(
			ErrorTypeNetwork,
			"Network error: Unable to connect to Jira. Please check your internet connection.",
			err,
		)
	}

	// Check for date range errors
	if strings.Contains(errMsgLower, "date range") ||
		strings.Contains(errMsgLower, "invalid date") {
		return NewReportError(
			ErrorTypeInvalidDateRange,
			"Invalid date range: Start date must be before or equal to end date.",
			err,
		)
	}

	// Default to unknown error
	return NewReportError(
		ErrorTypeUnknown,
		errMsg,
		err,
	)
}

// GetUserFriendlyMessage returns a user-friendly error message
func (re *ReportError) GetUserFriendlyMessage() string {
	switch re.Type {
	case ErrorTypeAuth:
		return "Authentication failed. Please check your API credentials in ~/.jirarc"
	case ErrorTypeNetwork:
		return "Network error: Unable to connect to Jira. Please check your internet connection."
	case ErrorTypeRateLimit:
		return "Jira API rate limit exceeded. Please wait a few minutes and try again."
	case ErrorTypeInvalidDateRange:
		return "Invalid date range: Start date must be before or equal to end date."
	case ErrorTypePartialData:
		return "Warning: Some data could not be retrieved. Showing partial results."
	default:
		return re.Message
	}
}
