package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReqID(t *testing.T) {
	ctx := CreateReqIDCtx(context.TODO(), "request-id")
	assert.Equal(t, "request-id", ctx.Value(reqCtx), "Request must be set")
}
