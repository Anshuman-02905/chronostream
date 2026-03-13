package chunker

import (
	"crypto/sha256"
	"encoding/hex"
)

// This package is a pure, deterministic utility for splitting raw message into fixed size byte fragments

//Design Constraints
//Its a pure function -No I/O  No randomness
//Deterministic O/P from the same I/P
//Byte based chunking not rune based - Byte based is taking into consideration number of bytes and rune based takes number of charecter
//Zero Dependency to external packages

//Note
//payload must be byte not string
//Chunkindex must be stable and sequential
// Fragment must be immutable once created

// Avoid chunker as an interface as multiple implementation of chunker is not expected here
// Mocking is an antipattern as its not an external dependency
// Interface is offen implemneted by structs which hold some configuration but this is pure function stateless

//Chunk wil split the message  into deterministic fixed size byte fragments
// 1) validate input if the chunkSize<=0 return empty slice and empty message  return empty slice
// 2) Convert message to []byte explicitly - must handle non ASCII explicitly Do not range over runes
// 3) Compute MessageId Deterministicly -Hash full Message bytes , Use SHA 256 , Truncate or format as string
// 4) Compute total Chunks
// 5) split byte into chunks preserve Order Last chunk might be smaller
//  6) Constryct  fragment Vaues- ChunkIndex starts from 0 , Total Chunks same for all fragments , Payload musbe a copy not a shared slice

// Prohibition - No logging, No randomness ,  no glonal state no external libraries
// Collect
func Chunk(message string, chunkSize int, setters ...Option) []Fragment { // this will be used for orchaestrating public facing API
	if chunkSize <= 0 || len(message) == 0 { // We will check Edge case here
		return nil
	}
	opts := defaultOptions() //Set the default options

	for _, setter := range setters { //if setters are available WithPadding() WithPayloadCopy() etc
		setter(&opts)
	}
	//
	msgBytes, msgID := prepareMessage(message)
	totalChunks := computeTotalChunks(len(msgBytes), chunkSize)
	return buildFragments(msgBytes, msgID, chunkSize, totalChunks, opts)

}

func prepareMessage(msg string) ([]byte, string) {
	messageBytes := []byte(msg)
	hash := sha256.Sum256(messageBytes)
	messageID := hex.EncodeToString(hash[:])
	return messageBytes, messageID
}

func computeTotalChunks(byteLen int, chunkSize int) int {
	return (byteLen + chunkSize - 1) / chunkSize
}
func padPayload(msgBytes []byte, chunkSize int) []byte {
	remainder := len(msgBytes) % chunkSize
	if remainder == 0 {
		return msgBytes
	}
	paddingNeeded := chunkSize - remainder
	paddedBytes := make([]byte, len(msgBytes)+paddingNeeded)
	copy(paddedBytes, msgBytes)
	return paddedBytes
}
