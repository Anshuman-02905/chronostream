package chunker

//inernal stalefree pure Function Builder

// func buildFragments(msgBytes []byte, msgID string, chunkSize, totalChunks int, opts FragmentOptions) []Fragment {
// 	outputFragment := make([]Fragment, 0, totalChunks)

// 	//option1
// 	if opts.PadLastChunk{
// 		msgBytes=padPayload(msgBytes,chunkSize)
// 	}
// 	outputFragment:=make([]Fragment,0,totalChunks)
// 	for i:=0;i<totalChunks;i++{

// 	}

// 	for i := 0; i < totalChunks; i++ {
// 		start := i * chunkSize
// 		end := start + chunkSize
// 		fragment := NewFragment(msgID, i, totalChunks, start, end)
// 		outputFragment = append(outputFragment, fragment)
// 	}
// 	return outputFragment

// }

func buildFragments(msgBytes []byte, msgID string, chunkSize, totalChunks int, opts FragmentOptions) []Fragment {

	// Option 1: Apply Padding (Delegated to your new prep helper)
	if opts.PadLastChunk {
		msgBytes = padPayload(msgBytes, chunkSize)
	}

	outputFragment := make([]Fragment, 0, totalChunks)

	for i := 0; i < totalChunks; i++ {
		start := i * chunkSize

		// This protects the last loop from crashing if padding was false
		end := min(start+chunkSize, len(msgBytes))

		// Option 2: Copy Payload logic
		var payload []byte
		if opts.CopyPayload {
			// Allocate entirely new memory to hold the slice
			payload = make([]byte, end-start)
			copy(payload, msgBytes[start:end])
		} else {
			// Just use a reference slice (the default, zero-allocation behavior)
			payload = msgBytes[start:end]
		}

		// Option 3: Truncate ID logic
		finalMsgID := msgID
		if opts.TruncateMessageID > 0 && opts.TruncateMessageID < len(msgID) {
			finalMsgID = msgID[:opts.TruncateMessageID]
		}

		// At this exact point in time, all options have been handled.
		// You are left with pure data to plug into your Fragment.
		fragment := NewFragment(finalMsgID, i, totalChunks, payload)
		outputFragment = append(outputFragment, fragment)
	}

	return outputFragment
}
