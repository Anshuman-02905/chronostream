package chunker

//Fragment represents one ordered piece of chunked message

type Fragment struct {
	MessageID   string //deterministic hashderived ID of full message
	ChunkIndex  int    // zero based Index
	TotalChunks int    //total number of chunks for the message
	Payload     []byte // raw bytes for the fragment
}

func NewFragment(messageID string, index, totalChunks int, payload []byte) Fragment {
	return Fragment{
		MessageID:   messageID,
		ChunkIndex:  index,
		TotalChunks: totalChunks,
		Payload:     payload,
	}
}
