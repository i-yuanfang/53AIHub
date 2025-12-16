package openaicompatible

import (
	"context"
	"fmt"

	"github.com/pkoukk/tiktoken-go"
	"github.com/songquanpeng/one-api/common/config"
	"github.com/songquanpeng/one-api/common/logger"
	"github.com/songquanpeng/one-api/relay/model"
)

// ErrorWrapper wraps an error with additional context
func ErrorWrapper(err error, code string, statusCode int) *model.ErrorWithStatusCode {
	logger.Error(context.TODO(), fmt.Sprintf("[%s]%+v", code, err))

	Error := model.Error{
		Message: err.Error(),
		Type:    "api_error",
		Code:    code,
	}
	return &model.ErrorWithStatusCode{
		Error:      Error,
		StatusCode: statusCode,
	}
}

var defaultTokenEncoder *tiktoken.Tiktoken

func init() {
	// Initialize default token encoder
	encoder, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
	if err != nil {
		logger.SysError(fmt.Sprintf("failed to get gpt-3.5-turbo token encoder: %s", err.Error()))
		// Use a fallback if initialization fails
		defaultTokenEncoder = nil
	} else {
		defaultTokenEncoder = encoder
	}
}

func getTokenEncoder(model string) *tiktoken.Tiktoken {
	// Try to get encoder for specific model
	encoder, err := tiktoken.EncodingForModel(model)
	if err != nil {
		// Fall back to default encoder
		if defaultTokenEncoder != nil {
			return defaultTokenEncoder
		}
		// If even default encoder is nil, try one more time
		encoder, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
		if err != nil {
			logger.SysError(fmt.Sprintf("failed to get token encoder: %s", err.Error()))
			return nil
		}
		return encoder
	}
	return encoder
}

func getTokenNum(tokenEncoder *tiktoken.Tiktoken, text string) int {
	if tokenEncoder == nil {
		// Approximate token count if encoder is not available
		return int(float64(len(text)) * 0.38)
	}
	if config.ApproximateTokenEnabled {
		return int(float64(len(text)) * 0.38)
	}
	return len(tokenEncoder.Encode(text, nil, nil))
}

// CountTokenText counts the number of tokens in a text string
func CountTokenText(text string, model string) int {
	tokenEncoder := getTokenEncoder(model)
	return getTokenNum(tokenEncoder, text)
}
