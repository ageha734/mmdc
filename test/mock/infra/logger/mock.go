package logger

type MockLogger struct {
	InfoFunc      func(msg string)
	SuccessFunc   func(msg string)
	WarnFunc      func(msg string)
	ErrorExitFunc func(msg string)
}

func (m *MockLogger) Info(msg string) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg)
	}
}

func (m *MockLogger) Success(msg string) {
	if m.SuccessFunc != nil {
		m.SuccessFunc(msg)
	}
}

func (m *MockLogger) Warn(msg string) {
	if m.WarnFunc != nil {
		m.WarnFunc(msg)
	}
}

func (m *MockLogger) ErrorExit(msg string) {
	if m.ErrorExitFunc != nil {
		m.ErrorExitFunc(msg)
	}
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}
