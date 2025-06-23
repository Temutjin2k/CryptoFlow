package exchange

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"marketflow/internal/domain"
	"marketflow/pkg/logger"
	"net"
)

// Exchange represent single exchange data source
type Exchange struct {
	Name   string
	Addr   string
	conn   net.Conn
	cancel context.CancelFunc

	log logger.Logger
}

// NewExchange creates new instance of Exchange
func NewExchange(name, connAddr string, log logger.Logger) *Exchange {
	return &Exchange{
		Name: name,
		Addr: connAddr,

		log: log,
	}
}

// Start returns channel with data from given source and implements "Generator" pattern.
func (e *Exchange) Start(ctx context.Context) (<-chan domain.PriceData, error) {
	// Using cancel func to cancel the goroutine when Stop() called
	ctx, cancel := context.WithCancel(ctx)
	e.cancel = cancel

	conn, err := net.Dial("tcp", e.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	e.conn = conn
	e.log.Info(ctx, "connected", "exchange", e.Name, "addr", e.Addr)

	out := make(chan domain.PriceData)

	go func() {
		defer conn.Close()
		defer close(out)

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Bytes()

			// Mapping to struct
			var data domain.PriceData
			if err := json.Unmarshal(line, &data); err != nil {
				e.log.Error(ctx, "failed to parse JSON", "error", err)
				continue
			}
			data.Exchange = e.Name

			select {
			case out <- data:
			case <-ctx.Done():
				e.log.Info(ctx, "context cancelled", "exchange", e.Name)
				return
			}
		}

		if err := scanner.Err(); err != nil {
			e.log.Error(ctx, "scanner error", "error", err)
		}
	}()

	return out, nil
}

// Stop closes the connection
func (e *Exchange) Stop() error {
	if e.cancel != nil {
		e.cancel() // Контекст отменяется — горутина завершится
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
