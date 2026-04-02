#!/usr/bin/env bash

# This script acts like a "tail -f" for your Kinesis stream.
# It initializes iterators for all shards at "LATEST" and continuously polls them.

STREAM_NAME="chronostream-events"
REGION="ap-south-1"

echo "Initializing Live Stream Tail for $STREAM_NAME..."

# 1. Get all Shard IDs into an array
SHARD_IDS=($(aws kinesis describe-stream \
  --stream-name $STREAM_NAME \
  --region $REGION \
  --query "StreamDescription.Shards[*].ShardId" \
  --output text))

echo "Found ${#SHARD_IDS[@]} Shards. Starting to listen for new events..."
echo "----------------------------------------"

# 2. Initialize a standard indexed array to hold the current iterator for each shard
# (macOS uses older bash 3.2 by default, which does not support associative arrays)
declare -a iterators

for i in "${!SHARD_IDS[@]}"; do
  SHARD_ID="${SHARD_IDS[$i]}"
  # Use LATEST so we only see new events arriving from this moment on
  ITERATOR=$(aws kinesis get-shard-iterator \
    --stream-name $STREAM_NAME \
    --shard-id "$SHARD_ID" \
    --shard-iterator-type LATEST \
    --region $REGION \
    --query "ShardIterator" \
    --output text)
    
  iterators[$i]=$ITERATOR
done

# 3. Infinite loop to continuously poll all shards
while true; do
  for i in "${!SHARD_IDS[@]}"; do
    SHARD_ID="${SHARD_IDS[$i]}"
    CURRENT_ITERATOR="${iterators[$i]}"
    
    # If the iterator is empty (e.g., shard closed), skip
    if [ -z "$CURRENT_ITERATOR" ] || [ "$CURRENT_ITERATOR" == "null" ]; then
        continue
    fi

    # Fetch records
    RESULT=$(aws kinesis get-records \
      --shard-iterator "$CURRENT_ITERATOR" \
      --region $REGION \
      2>/dev/null) # Suppress errors (like expired iterators temporarily)

    if [ -z "$RESULT" ]; then
        continue
    fi

    RECORD_COUNT=$(echo "$RESULT" | jq '.Records | length')

    if [ "$RECORD_COUNT" -gt 0 ]; then
      echo "[$(date '+%Y-%m-%d %H:%M:%S')] ⚡ RECEIVED $RECORD_COUNT RECORDS FROM $SHARD_ID"
      echo "$RESULT" | jq -r '.Records[].Data | @base64d | fromjson'
      echo "----------------------------------------"
    fi

    # Update the iterator for this shard for the next loop iteration
    NEXT_ITERATOR=$(echo "$RESULT" | jq -r '.NextShardIterator')
    iterators[$i]=$NEXT_ITERATOR
  done
  
  # Pause briefly to avoid hitting AWS API rate limits
  sleep 1
done

