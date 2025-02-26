// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/tq-systems/public-go-utils/v3/mqtt (interfaces: Client)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=../mocks/mqtt/mock_client.go -package=mqtt github.com/tq-systems/public-go-utils/v3/mqtt Client
//

// Package mqtt is a generated GoMock package.
package mqtt

import (
	reflect "reflect"

	mqtt "github.com/tq-systems/public-go-utils/v3/mqtt"
	gomock "go.uber.org/mock/gomock"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// Publish mocks base method.
func (m *MockClient) Publish(arg0 string, arg1 byte, arg2 bool, arg3 protoreflect.ProtoMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockClientMockRecorder) Publish(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockClient)(nil).Publish), arg0, arg1, arg2, arg3)
}

// PublishEmpty mocks base method.
func (m *MockClient) PublishEmpty(arg0 string, arg1 byte, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishEmpty", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishEmpty indicates an expected call of PublishEmpty.
func (mr *MockClientMockRecorder) PublishEmpty(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishEmpty", reflect.TypeOf((*MockClient)(nil).PublishEmpty), arg0, arg1, arg2)
}

// PublishRaw mocks base method.
func (m *MockClient) PublishRaw(arg0 string, arg1 byte, arg2 bool, arg3 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishRaw", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishRaw indicates an expected call of PublishRaw.
func (mr *MockClientMockRecorder) PublishRaw(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishRaw", reflect.TypeOf((*MockClient)(nil).PublishRaw), arg0, arg1, arg2, arg3)
}

// Subscribe mocks base method.
func (m *MockClient) Subscribe(arg0 string, arg1 mqtt.Callback) (mqtt.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1)
	ret0, _ := ret[0].(mqtt.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockClientMockRecorder) Subscribe(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockClient)(nil).Subscribe), arg0, arg1)
}
