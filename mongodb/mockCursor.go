package mongodb

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type (
	DecodeStubResult struct {
		Error error
	}
)

type StubMongoCursor struct {
	mock.Mock
	DecodeStubResult func(val interface{}) DecodeStubResult
}

func (m *StubMongoCursor) Decode(val interface{}) error {
	return m.DecodeStubResult(val).Error
}

func (m *StubMongoCursor) Err() error {
	return nil
}

func (m *StubMongoCursor) Next(ctx context.Context) bool {
	return true
}

func (m *StubMongoCursor) Close(ctx context.Context) error {
	return nil
}

func (m *StubMongoCursor) ID() int64 {
	return 0
}
