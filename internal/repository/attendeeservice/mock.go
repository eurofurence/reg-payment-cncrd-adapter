package attendeeservice

import "context"

type Mock interface {
	AttendeeService

	Reset()
	Recording() []string
	SimulateGetError(err error)
}

type MockImpl struct {
	recording        []string
	simulateGetError error
}

var (
	_ AttendeeService = (*MockImpl)(nil)
	_ Mock            = (*MockImpl)(nil)
)

func newMock() Mock {
	return &MockImpl{
		recording: make([]string, 0),
	}
}

func (m *MockImpl) GetAttendee(ctx context.Context, id uint) (AttendeeDto, error) {
	if m.simulateGetError != nil {
		return AttendeeDto{}, m.simulateGetError
	}

	attendee := AttendeeDto{
		Email: "jsquirrel_github_9a6d@packetloss.de",
	}

	return attendee, nil
}

// only used in tests

func (m *MockImpl) Reset() {
	m.recording = make([]string, 0)
	m.simulateGetError = nil
}

func (m *MockImpl) Recording() []string {
	return m.recording
}

func (m *MockImpl) SimulateGetError(err error) {
	m.simulateGetError = err
}
