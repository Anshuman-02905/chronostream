package chunker

import (
	"testing"
)

func TestChunk_Basic(t *testing.T) {
	msg := "Hello World"
	chunkSize := 5
	fragments := Chunk(msg, chunkSize)

	// "Hello" (5), " Worl" (5), "d" (1)
	expectedChunks := 3
	if len(fragments) != expectedChunks {
		t.Fatalf("expected %d fragments, got %d", expectedChunks, len(fragments))
	}

	if fragments[0].TotalChunks != expectedChunks {
		t.Errorf("expected TotalChunks %d, got %d", expectedChunks, fragments[0].TotalChunks)
	}

	for i, f := range fragments {
		if f.ChunkIndex != i {
			t.Errorf("expected ChunkIndex %d, got %d", i, f.ChunkIndex)
		}
	}

	if string(fragments[0].Payload) != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", string(fragments[0].Payload))
	}
	if string(fragments[2].Payload) != "d" {
		t.Errorf("expected 'd', got '%s'", string(fragments[2].Payload))
	}
}

func TestChunk_Padding(t *testing.T) {
	msg := "Hello"
	chunkSize := 3
	// "Hel" (3), "lo\x00" (3)
	fragments := Chunk(msg, chunkSize, WithPadding())

	if len(fragments) != 2 {
		t.Fatalf("expected 2 fragments, got %d", len(fragments))
	}

	if len(fragments[1].Payload) != 3 {
		t.Errorf("expected payload length 3, got %d", len(fragments[1].Payload))
	}

	if fragments[1].Payload[2] != 0 {
		t.Errorf("expected null padding, got %d", fragments[1].Payload[2])
	}
}

func TestChunk_PayloadCopy(t *testing.T) {
	msg := "Hello"
	chunkSize := 2
	fragments := Chunk(msg, chunkSize, WithPayloadCopy())

	// Verify it still works correctly
	if len(fragments) != 3 {
		t.Fatalf("expected 3 fragments, got %d", len(fragments))
	}

	// We can't easily verify the "copy" part from outside,
	// but we ensure the content is correct.
	if string(fragments[0].Payload) != "He" {
		t.Errorf("expected 'He', got '%s'", string(fragments[0].Payload))
	}
}

func TestChunk_TruncateID(t *testing.T) {
	msg := "Hello"
	chunkSize := 2
	fragments := Chunk(msg, chunkSize, WithTruncateID(8))

	if len(fragments[0].MessageID) != 8 {
		t.Errorf("expected MessageID length 8, got %d", len(fragments[0].MessageID))
	}
}

func TestChunk_EdgeCases(t *testing.T) {
	t.Run("EmptyMessage", func(t *testing.T) {
		fragments := Chunk("", 5)
		if fragments != nil {
			t.Errorf("expected ni fragments for empty message")
		}
	})

	t.Run("InvalidChunkSize", func(t *testing.T) {
		fragments := Chunk("hello", 0)
		if fragments != nil {
			t.Errorf("expected nil fragments for zero chunk size")
		}
	})
}
