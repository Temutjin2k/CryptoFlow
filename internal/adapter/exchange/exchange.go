package exchange

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"marketflow/pkg/logger"
	"net"
	"time"
)

// Exchange represent single exchange data source
type Exchange struct {
	name   types.Exchange
	Addr   string
	conn   net.Conn
	cancel context.CancelFunc

	log logger.Logger
}

// NewExchange creates new instance of Exchange
func NewExchange(name types.Exchange, connAddr string, log logger.Logger) *Exchange {
	return &Exchange{
		name: name,
		Addr: connAddr,

		log: log,
	}
}

// Start returns channel with data from given source and implements "Generator" pattern.
func (e *Exchange) Start(ctx context.Context) (<-chan *domain.PriceData, error) {
	// Using cancel func to cancel the goroutine when Stop() called
	ctx, cancel := context.WithCancel(ctx)
	e.cancel = cancel

	conn, err := net.Dial("tcp", e.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	e.conn = conn

	log := e.log.GetSlogLogger().With("exchange", e.Name(), "addr", e.Addr)
	log.InfoContext(ctx, "connected to exchange!")

	out := make(chan *domain.PriceData)

	go func() {
		defer conn.Close()
		defer close(out)

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Bytes()

			// Mapping to struct
			data := new(domain.PriceData)
			if err := json.Unmarshal(line, data); err != nil {
				log.ErrorContext(ctx, "failed to parse JSON", "error", err)
				continue
			}
			data.Exchange = e.name

			select {
			case out <- data:
			case <-ctx.Done():
				log.InfoContext(ctx, "context cancelled")
				return
			}
		}

		if err := scanner.Err(); err != nil {
			log.ErrorContext(ctx, "scanner error", "error", err)
		}
	}()

	return out, nil
}

// Stop closes the connection
func (e *Exchange) Stop() error {
	if e.cancel != nil {
		e.cancel() // canceling context to stop the gouroutine
	}

	ctx := context.Background()
	e.log.Info(ctx, "closing connection", "name", e.Name, "addres", e.Addr)

	if e.conn != nil {
		if err := e.conn.Close(); err != nil {
			e.log.Error(ctx, "failed to close connection", "exchange", e.Name, "error", err)
			return err
		}
	}
	return nil
}

func (e *Exchange) Name() string {
	return string(e.name)
}

// Health checks if the exchange service is available by attempting a connection
func (e *Exchange) Health(ctx context.Context) (bool, error) {
	conn, err := net.DialTimeout("tcp", e.Addr, 5*time.Second)
	if err != nil {
		return false, fmt.Errorf("health check failed for %s: %w", e.name, err)
	}
	defer conn.Close()

	// If we reached here, connection was successful
	return true, nil
}
