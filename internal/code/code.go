package code

import (
	"fmt"

	"deutsch/internal/types"
)

// Code represents an API response code for unified error handling.
// Success is 0; errors are categorized by module (1000-9999).
// Use with Base { Code Code `json:"code"` Msg string `json:"msg"` }.
type Code int

const (
	// Success
	CodeSuccess Code = 0

	// Authentication & Authorization (1000-1999)
	CodeInvalidInviteCode  Code = 1001 // invalid or used invite code
	CodeInviteCodeRequired Code = 1002 // invite code required
	CodePhoneExists        Code = 1003 // phone already exists
	CodeUsernameExists     Code = 1004 // username already exists
	CodeEmailExists        Code = 1015 // email already exists
	CodeInvalidCredentials Code = 1005 // invalid credentials (phone/password error)
	CodeUserDisabled       Code = 1006 // user disabled
	CodeInvalidToken       Code = 1007 // invalid or expired token
	CodeMissingToken       Code = 1008 // missing token
	CodeTokenBlacklisted   Code = 1009 // token blacklisted (invalidated)
	CodeAdminRequired      Code = 1010 // admin role required
	CodeUnauthorized       Code = 1011 // unauthorized or insufficient permissions
	CodePasswordTooShort   Code = 1012 // password too short (minimum 8 characters)
	CodePhoneInvalid       Code = 1013 // invalid phone format
	CodeUsernameInvalid    Code = 1014 // invalid username format (6-50 characters)

	// User Management (2000-2999)
	CodeUserNotFound        Code = 2001 // user not found
	CodeProfileUpdateFailed Code = 2002 // failed to update user profile
	CodeNicknameTooLong     Code = 2003 // nickname too long (max 50 characters)
	CodeDescriptionTooLong  Code = 2004 // description too long
	CodePreferencesInvalid  Code = 2005 // invalid user preferences

	// Invite Code Management (2100-2199)
	CodeInviteGenerationFailed Code = 2101 // failed to generate invite codes
	CodeInviteCountExceeded    Code = 2102 // invite count exceeded
	CodeInviteNotFound         Code = 2103 // invite code not found
	CodeInviteCodeAlreadyUsed  Code = 2104 // invite code already used, cannot modify

	// Question Bank (3000-3999)
	CodeQuestionNotFound   Code = 3001 // question not found
	CodeCategoryInvalid    Code = 3002 // invalid category (e.g., history, law)
	CodeDifficultyInvalid  Code = 3003 // invalid difficulty (easy/medium/hard)
	CodeSearchQueryEmpty   Code = 3004 // search query empty
	CodeSearchQueryTooLong Code = 3005 // search query too long
	CodeLimitExceeded      Code = 3006 // limit exceeded (limit > 100)
	CodeOffsetInvalid      Code = 3007 // invalid offset (offset < 0)
	CodeStateNotFound      Code = 3008 // state id not found

	// Practice Tests (4000-4999)
	CodeTestNotFound          Code = 4001 // test not found
	CodeTestTypeInvalid       Code = 4002 // invalid test type (full/mini/category)
	CodeTestAlreadyStarted    Code = 4003 // test already started, cannot restart
	CodeTestExpired           Code = 4004 // test expired
	CodeAnswersInvalid        Code = 4005 // invalid answers (options out of range)
	CodeTimeLimitExceeded     Code = 4006 // time limit exceeded for test submission
	CodeTestSubmitFailed      Code = 4007 // failed to submit test
	CodeInsufficientQuestions Code = 4008 // insufficient questions (e.g., no questions in category)

	// Results & Analytics (5000-5999)
	CodeResultNotFound            Code = 5001 // test result not found
	CodeProgressNotFound          Code = 5002 // progress record not found
	CodeProgressCalculationFailed Code = 5003 // failed to calculate progress
	CodeWeakAreasEmpty            Code = 5004 // weak areas analysis empty
	CodeAverageScoreError         Code = 5005 // average score calculation error
	CodeLimitForResultsExceeded   Code = 5006 // results limit exceeded
	CodeExamRecordNotFound        Code = 5007 // exam record not found

	// System & General Errors (9000-9999)
	CodeInternalServerError Code = 9001 // internal server error
	CodeDatabaseError       Code = 9002 // database connection or query error
	CodeRedisError          Code = 9003 // Redis operation failed
	CodeConfigLoadError     Code = 9004 // config load failed
	CodeValidationError     Code = 9005 // parameter validation failed
	CodeRateLimitExceeded   Code = 9006 // rate limit exceeded
	CodeRequestTimeout      Code = 9007 // request timeout
	CodeResourceNotFound    Code = 9008 // resource not found (general)
	CodePayloadTooLarge     Code = 9009 // payload too large
	CodeMethodNotAllowed    Code = 9010 // HTTP method not allowed
)

// Messages maps Code to its English message for API responses.
var Messages = map[Code]string{
	CodeSuccess:                   "success",
	CodeInvalidInviteCode:         "invalid or used invite code",
	CodeInviteCodeRequired:        "invite code required",
	CodePhoneExists:               "phone already exists",
	CodeUsernameExists:            "username already exists",
	CodeEmailExists:               "email already exists",
	CodeInvalidCredentials:        "invalid credentials (phone/password error)",
	CodeUserDisabled:              "user disabled",
	CodeInvalidToken:              "invalid or expired token",
	CodeMissingToken:              "missing token",
	CodeTokenBlacklisted:          "token blacklisted (invalidated)",
	CodeAdminRequired:             "admin role required",
	CodeUnauthorized:              "unauthorized or insufficient permissions",
	CodePasswordTooShort:          "password too short (minimum 8 characters)",
	CodePhoneInvalid:              "invalid phone format",
	CodeUsernameInvalid:           "invalid username format (6-50 characters)",
	CodeUserNotFound:              "user not found",
	CodeProfileUpdateFailed:       "failed to update user profile",
	CodeNicknameTooLong:           "nickname too long (max 50 characters)",
	CodeDescriptionTooLong:        "description too long",
	CodePreferencesInvalid:        "invalid user preferences",
	CodeInviteGenerationFailed:    "failed to generate invite codes",
	CodeInviteCountExceeded:       "invite count exceeded",
	CodeInviteNotFound:            "invite code not found",
	CodeInviteCodeAlreadyUsed:     "invite code already used, cannot modify",
	CodeQuestionNotFound:          "question not found",
	CodeCategoryInvalid:           "invalid category (e.g., history, law)",
	CodeDifficultyInvalid:         "invalid difficulty (easy/medium/hard)",
	CodeSearchQueryEmpty:          "search query empty",
	CodeSearchQueryTooLong:        "search query too long",
	CodeLimitExceeded:             "limit exceeded (limit > 100)",
	CodeOffsetInvalid:             "invalid offset (offset < 0)",
	CodeStateNotFound:             "state id not found",
	CodeTestNotFound:              "test not found",
	CodeTestTypeInvalid:           "invalid test type (full/mini/category)",
	CodeTestAlreadyStarted:        "test already started, cannot restart",
	CodeTestExpired:               "test expired",
	CodeAnswersInvalid:            "invalid answers (options out of range)",
	CodeTimeLimitExceeded:         "time limit exceeded for test submission",
	CodeTestSubmitFailed:          "failed to submit test",
	CodeInsufficientQuestions:     "insufficient questions (e.g., no questions in category)",
	CodeResultNotFound:            "test result not found",
	CodeProgressNotFound:          "progress record not found",
	CodeProgressCalculationFailed: "failed to calculate progress",
	CodeWeakAreasEmpty:            "weak areas analysis empty",
	CodeAverageScoreError:         "average score calculation error",
	CodeLimitForResultsExceeded:   "results limit exceeded",
	CodeExamRecordNotFound:        "exam record not found",
	CodeInternalServerError:       "internal server error",
	CodeDatabaseError:             "database connection or query error",
	CodeRedisError:                "Redis operation failed",
	CodeConfigLoadError:           "config load failed",
	CodeValidationError:           "parameter validation failed",
	CodeRateLimitExceeded:         "rate limit exceeded",
	CodeRequestTimeout:            "request timeout",
	CodeResourceNotFound:          "resource not found (general)",
	CodePayloadTooLarge:           "payload too large",
	CodeMethodNotAllowed:          "HTTP method not allowed",
}

// GetMessage returns the message for the given code.
func GetMessage(c Code) string {
	if msg, ok := Messages[c]; ok {
		return msg
	}
	return "unknown error (code: " + fmt.Sprint(int(c)) + ")"
}

// Message returns the message for the code (method on Code).
func (c Code) Message() string {
	return GetMessage(c)
}

func BaseSuccessResp() *types.Base {
	return &types.Base{
		Code: int(CodeSuccess),
		Msg:  CodeSuccess.Message(),
	}
}

// CodeError 用于业务错误，携带 code 和 msg
type CodeError struct {
	Code Code
	Msg  string
}

func (e *CodeError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.Code.Message()
}

func NewCodeError(c Code) *CodeError {
	return &CodeError{Code: c, Msg: c.Message()}
}

func NewCodeErrorWithMsg(c Code, msg string) *CodeError {
	return &CodeError{Code: c, Msg: msg}
}
