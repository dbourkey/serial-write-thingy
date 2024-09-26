package status

import (
	"context"
	"errors"
	"time"
)

// Writer is an interface for some client object which can persist data.
type Writer interface {
	// Write takes an object key and a value, returns an error if the data was not written to the persistant store.
	Write(key string, value any) error
}

type Serialiser struct {
	dbClient Writer

	buffer     map[string]ContainerReport // Stores status updates which were not written to the remote.
	lastupdate map[string]ContainerReport // Records status updates which were successfully written to the remote.

	flushInterval time.Duration

	in chan ContainerReport
}

func NewSerialiser(remoteDBClient Writer, flushInterval time.Duration) *Serialiser {
	return &Serialiser{
		dbClient:      remoteDBClient,
		buffer:        make(map[string]ContainerReport),
		lastupdate:    make(map[string]ContainerReport),
		flushInterval: flushInterval,
		in:            make(chan ContainerReport),
	}
}

func (s *Serialiser) Run(ctx context.Context) error {
	tick := time.NewTicker(s.flushInterval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case statusUpdate := <-s.in:
			s.updateStatus(statusUpdate)
		case <-tick.C:
			s.flushBuffer()
		}
	}
}

// Write accepts a new status update to be peristed to the remote database, or buffered for later if the database is not available.
// The Serialiser provides update deduplication, so in general, the most recent update will be sent.
func (s *Serialiser) Write(ctx context.Context, newStatus ContainerReport) error {
	select {
	case <-ctx.Done():
		return errors.New("context cancelled before writing update")
	case s.in <- newStatus:
		return nil
	}
}

// updateStatus attempts to write the status update to the remote.
func (s *Serialiser) updateStatus(newUpdate ContainerReport) {
	// Check if a more recent status has already been written to the datastore.
	if lastUpdate, ok := s.lastupdate[newUpdate.ContainerID]; ok && !newUpdate.Time().After(lastUpdate.Time()) {
		return // Remote has a more recent update, so nothing else to do here.
	}

	// Attempt to update the remote with the new status.
	if err := s.dbClient.Write(newUpdate.ContainerID, newUpdate); err != nil {
		s.buffer[newUpdate.ContainerID] = newUpdate // Failed to update the remote, so buffer the update for later retries.
		return
	}

	s.lastupdate[newUpdate.ContainerID] = newUpdate // Record the last successful status update for the relevant container.
}

// flushBuffer attempts to write any unpersisted status updates.
func (s *Serialiser) flushBuffer() {
	for _, statusReport := range s.buffer {
		s.updateStatus(statusReport)
	}
}
