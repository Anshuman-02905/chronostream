# ChronoStream: Chunker Package Architecture

**Objective:**
The `chunker` package provides a pure, deterministic, uncoupled utility for splitting raw string messages into fixed-size byte fragments. It is designed to be entirely stateless, devoid of external I/O, and strictly adheres to the Single Responsibility Principle.

## Core Components
1. **Orchestrator (`chunker.go` / `Chunk`)**: The public-facing API. It acts purely as a coordinator. It validates inputs, reads configurations, and passes data down the pipeline.
2. **Configuration (`options.go` / `FragmentOptions`)**: Implements the Functional Options Pattern. It safely stores behaviour switches (like whether to pad the final chunk or copy memory) without crowding the public API signature.
3. **Preparation (`prep.go`)**: Handles all data derivation *before* any chunking occurs. This includes hashing the full message for a deterministic ID and calculating total chunks needed.
4. **Builder (`builder.go` / `buildFragments`)**: The core loop. It receives perfectly prepared data and configurations, slices the byte array, applies any padding/copy logic, and constructs the final array of Data Transfer Objects (`Fragment`).

## Architectural Diagram

```mermaid
flowchart TD
    %% Initial Entry
    A[Public API: Chunk(msg, size, ...opts)] --> B{Input Valid?}
    B -- No --> ERR[Return Empty/Error]
    
    %% Configuration
    B -- Yes --> C[Initialize defaultOptions()]
    C --> D[Apply setters ...Option]
    
    %% Preparation Phase
    D --> E[prepareMessage: Convert to Bytes & Hash]
    E --> F[computeTotalChunks: Calculate length / size ceiling]
    
    %% Builder Phase (buildFragments)
    F --> G[Enter Builder Pipeline]
    G --> H{opts.PadLastChunk?}
    H -- Yes --> I[Pad msgBytes with 0x00]
    H -- No --> J[Enter Loop: 0 to totalChunks]
    I --> J
    
    %% The Loop
    J --> K[Slice msgBytes start:end]
    K --> L{opts.CopyPayload?}
    L -- Yes --> M[Allocate new []byte & deep copy]
    L -- No --> N[Keep original slice reference]
    M --> O[Construct Fragment struct]
    N --> O[Construct Fragment struct]
    
    %% Output
    O --> P[Append to output slice]
    P --> Q{More chunks?}
    Q -- Yes --> J
    Q -- No --> R[Return []Fragment]

    style A fill:#4CAF50,stroke:#333,stroke-width:2px,color:#fff
    style R fill:#2196F3,stroke:#333,stroke-width:2px,color:#fff
    style J fill:#FF9800,stroke:#333,stroke-width:2px,color:#fff
```

## Knowledge Verification Questions (Self-Assessment)

**1. Slicing Memory vs. Copying Memory**
*Question:* By default, your builder slices the original `msgBytes` without copying it. What is the performance benefit of this? Conversely, what is the *danger* of doing this if the `[]Fragment` array is passed to an asynchronous Goroutine that lives for a long time? 

**2. The Option Pattern**
*Question:* If we had defined `Chunk(message string, chunkSize int, opts FragmentOptions)` instead of using `...Option`, what would happen if a user passed `nil`? Why does the `...Option` variadic approach combined with `defaultOptions()` physically prevent this specific class of runtime panics?

**3. Determinism**
*Question:* The architectural constraints state the package must be "pure" and "deterministic" with no randomness or time functions. Why is it so critical that the same string *always* produces the exact same `MessageID` and chunk boundaries every single time it passes through this package?

**4. Ceiling Math**
*Question:* In `prep.go`, you initially used `(byteLen+chunkSize+1)/chunkSize` to calculate total chunks, and we changed it to `(byteLen+chunkSize-1)/chunkSize`. Walk through the math: If a message is 16 bytes long, and the chunk size is 8, what does each formula output? *Why* does adding `chunkSize - 1` correctly behave like a mathematical ceiling function for Go's integer division?

**5. Capping vs. Padding**
*Question:* In the builder loop, we use `end := min(start+chunkSize, len(msgBytes))` to prevent a panic. If we didn't have `min()`, write out the `if/else` block you would need to write inside that loop to properly set the `end` variable for the final uneven snippet of data.
