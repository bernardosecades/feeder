package server_test

import (
	"context"
	"github.com/bernardosecades/feeder/pkg/server"
	"github.com/bernardosecades/feeder/pkg/service"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"testing"
	"time"
)

func TestServerDownAfterKeepAliveTime(t *testing.T) {
	// start server
	ctx := context.Background()
	cf := server.Config{
		Protocol:  "tcp",
		Host:      "",
		Port:      "5000",
		KeepAlive: time.Millisecond * 10,
		MaxConn:   1,
	}

	mockFeeder := &MockFeeder{}
	srv := server.NewServer(cf, mockFeeder)
	err := srv.Start(ctx)

	assert.Equal(t, context.DeadlineExceeded, err)

	// we ensure when server is down call to: Log, Report and Persist from feeder service
	assert.Equal(t, mockFeeder.CallsLog, 1)
	assert.Equal(t, mockFeeder.CallsReport, 1)
	assert.Equal(t, mockFeeder.CallsPersist, 1)
}

func TestServerDownByClient(t *testing.T) {
	// client send 'terminate' message to stop the server
	go func() {
		conn, err := net.Dial("tcp", "localhost:5005")
		assert.Nil(t, err)
		_, err = conn.Write([]byte("terminate\n"))
		assert.Nil(t, err)
	}()

	// start server
	ctx := context.Background()
	cf := server.Config{
		Protocol:  "tcp",
		Host:      "",
		Port:      "5005",
		KeepAlive: time.Millisecond * 10,
		MaxConn:   1,
	}

	mockFeeder := &MockFeeder{}
	srv := server.NewServer(cf, mockFeeder)

	err := srv.Start(ctx)

	assert.Equal(t, server.ErrClientIndicateTerminate, err)

	// we ensure when server is down call to: Log, Report and Persist from feeder service
	assert.Equal(t, mockFeeder.CallsLog, 1)
	assert.Equal(t, mockFeeder.CallsReport, 1)
	assert.Equal(t, mockFeeder.CallsPersist, 1)
}

func TestServerMaxConnectionsReached(t *testing.T) {
	go func() {
		conn1, err := net.Dial("tcp", "localhost:5010")
		assert.Nil(t, err)
		assert.NotNil(t, conn1)

		conn2, err := net.Dial("tcp", "localhost:5010")
		assert.Nil(t, err)
		assert.NotNil(t, conn2)

		// Server send a message with 'limit connections reached'
		// and disconnect him
		for {
			d := make([]byte, 120)
			_, err := conn2.Read(d)
			if err != nil {
				if err == io.EOF {
					assert.True(t, true) // Connection TCP was closed
					break
				} else {
					assert.Fail(t, "conn2 should have disconnected")
					break
				}
			}
		}
	}()

	// start server
	ctx := context.Background()
	cf := server.Config{
		Protocol:  "tcp",
		Host:      "",
		Port:      "5010",
		KeepAlive: time.Millisecond * 10,
		MaxConn:   1,
	}

	mockFeeder := &MockFeeder{}
	srv := server.NewServer(cf, mockFeeder)

	err := srv.Start(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)

	// we ensure when server is down call to: Log, Report and Persist from feeder service
	assert.Equal(t, mockFeeder.CallsLog, 1)
	assert.Equal(t, mockFeeder.CallsReport, 1)
	assert.Equal(t, mockFeeder.CallsPersist, 1)
}

type MockFeeder struct {
	CallsPersist int
	CallsReport  int
	CallsLog     int
}

func (m *MockFeeder) Persist() (service.SkusInserted, service.SkusInsertSkipped, error) {
	m.CallsPersist++
	return service.SkusInserted(0), service.SkusInsertSkipped(0), nil
}

func (m *MockFeeder) Report() (service.TotalUniqueSkus, service.TotalDuplicatedSkus, service.TotalInvalidSkus) {
	m.CallsReport++
	return service.TotalUniqueSkus(0), service.TotalDuplicatedSkus(0), service.TotalInvalidSkus(0)
}

func (m *MockFeeder) Log() {
	m.CallsLog++
}

func (m *MockFeeder) AddSku(sku string) {
}
