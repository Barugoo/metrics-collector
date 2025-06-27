package main

import (
	"testing"

	"github.com/skdiver33/metrics-collector/internal/store"
)

func TestAgent_SendMetrics(t *testing.T) {
	type fields struct {
		metricStorage store.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "positive test",
			fields: fields{
				metricStorage: store.NewMemStorage(0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{
				metricStorage: tt.fields.metricStorage,
				config:        AgentConfig{serverAddress: "localhost:8080", pollInterval: 2, reportInterval: 10},
			}
			agent.UpdateMetrics()
			if err := agent.SendMetrics(); (err != nil) != tt.wantErr {
				t.Errorf("Agent.SendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
