package logging

import (
	"go.uber.org/zap"
)

type Field = zap.Field

var (
	StringType  = zap.String

	IntType     = zap.Int
	Int8Type    = zap.Int8
	Int16Type   = zap.Int16
	Int32Type   = zap.Int32
	Int64Type   = zap.Int64
	Uint8Type   = zap.Uint8
	Uint16Type  = zap.Uint16
	Uint32Type  = zap.Uint32
	Uint64Type  = zap.Uint64

	Float32Type = zap.Float32
	Float64Type = zap.Float64
	
	BoolType    = zap.Bool

	ErrorType	= zap.Error

	AnyType     = zap.Any
)
