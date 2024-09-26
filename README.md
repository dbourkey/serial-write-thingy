# serial-write-thingy
Code snippets showing how an api could deduplicate status updates so only the most recent update is written to the remote.

## Description

- Request handlers unmarshal the json body, call the "Serialiser" Write method to push the status update to a channel (Which may be buffered to implement a queue)
- The Serialiser runs a for-select loop so only one update may be processed at once.
- Serialiser runs a full "flush" on an interval, during which time it cannot read new updates.
- Handlers will block until serialiser is available again, potentially timing out so client can resent.
- Serialiser attempts to write to the remote database, and if this is not possible it buffers the latest status update.
- Both the buffer and the "lastupdate" maps will fill forever and eventually OOM, we could handle this by locking (via sync.Mutex) initialising a new map, and swapping the references.
- This implementation does not provide particularly strong guaruntees, and the serialiser is a bottleneck.

## Better Ideas

For larger scale system with thousands of updates per second, we could use a very similar idea to this, but run several serialisers in parallel.
The handler would hash the container-id to map the status update to a mapped serialiser.

Another method would be a "sharded map", whereby we hash on container-id, with the key being another hashmap. This lets us manage memory by clearing/re-initialising shards at a time rather than the entire datastructure.
