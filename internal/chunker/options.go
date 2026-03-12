package chunker

// This struct wont be exported to the user directly
type FragmentOptions struct {
	CopyPayload       bool
	TruncateMessageID int
	PadLastChunk      bool
}

// Option is a functions which modifies FragmentOptions
type Option func(*FragmentOptions)

func defaultOptions() FragmentOptions {
	return FragmentOptions{
		CopyPayload:       false,
		TruncateMessageID: 0,
	}
}

func WithPayloadCopy() Option {
	return func(o *FragmentOptions) {
		o.CopyPayload = true
	}
}

func WithTruncateID(length int) Option {
	return func(o *FragmentOptions) {
		o.TruncateMessageID = length
	}
}

func WithPadding() Option {
	return func(o *FragmentOptions) {
		o.PadLastChunk = true
	}
}
