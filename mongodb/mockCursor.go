package mongodb

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type (
	DecodeMockResult struct {
		Error error
	}
)

type MongoCursorMock struct {
	mock.Mock
	DecodeMockResult DecodeMockResult
}

func (m *MongoCursorMock) Decode(val interface{}) error {
	return m.DecodeMockResult.Error
}

func (m *MongoCursorMock) Err() error {
	return nil
}

func (m *MongoCursorMock) Next(ctx context.Context) bool {
	return true
}

func (m *MongoCursorMock) Close(ctx context.Context) error {
	return nil
}

func (m *MongoCursorMock) ID() int64 {
	return 0
}
