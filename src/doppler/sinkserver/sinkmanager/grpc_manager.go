package sinkmanager

import (
	"log"
	"plumbing"
	"sync"
)

type grpcRegistry struct {
	mu       sync.RWMutex
	registry map[string][]GRPCSender
}

func newGRPCRegistry() *grpcRegistry {
	return &grpcRegistry{
		registry: make(map[string][]GRPCSender),
	}
}

type GRPCSender interface {
	Send(resp *plumbing.Response) (err error)
}

func (m *SinkManager) Stream(req *plumbing.StreamRequest, d plumbing.Doppler_StreamServer) error {
	log.Printf("Connection accepted for app %s", req.AppID)
	m.RegisterStream(req, d)
	log.Printf("Stream registered")
	<-d.Context().Done()
	return nil
}

func (m *SinkManager) RegisterStream(req *plumbing.StreamRequest, sender GRPCSender) {
	m.grpcStreams.mu.Lock()
	defer m.grpcStreams.mu.Unlock()
	m.grpcStreams.registry[req.AppID] = append(m.grpcStreams.registry[req.AppID], sender)
}

func (m *SinkManager) Firehose(req *plumbing.FirehoseRequest, d plumbing.Doppler_FirehoseServer) error {
	m.RegisterFirehose(req, d)
	<-d.Context().Done()
	return nil
}

func (m *SinkManager) RegisterFirehose(req *plumbing.FirehoseRequest, sender GRPCSender) {
	m.grpcFirehoses.mu.Lock()
	defer m.grpcFirehoses.mu.Unlock()
	m.grpcFirehoses.registry[req.SubID] = append(m.grpcFirehoses.registry[req.SubID], sender)
}
