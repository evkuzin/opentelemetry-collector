// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testcomponents

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

// ExampleReceiver is for testing purposes. We are defining an example config and factory
// for "examplereceiver" receiver type.
type ExampleReceiver struct {
	configmodels.ReceiverSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	// Configures the receiver server protocol.
	confignet.TCPAddr `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	ExtraSetting     string            `mapstructure:"extra"`
	ExtraMapSetting  map[string]string `mapstructure:"extra_map"`
	ExtraListSetting []string          `mapstructure:"extra_list"`
}

const recvType = "examplereceiver"

// ExampleReceiverFactory is factory for ExampleReceiver.
var ExampleReceiverFactory = receiverhelper.NewFactory(
	recvType,
	createReceiverDefaultConfig,
	receiverhelper.WithTraces(createTracesReceiver),
	receiverhelper.WithMetrics(createMetricsReceiver),
	receiverhelper.WithLogs(createLogsReceiver))

func createReceiverDefaultConfig() configmodels.Receiver {
	return &ExampleReceiver{
		ReceiverSettings: configmodels.ReceiverSettings{
			TypeVal: recvType,
			NameVal: recvType,
		},
		TCPAddr: confignet.TCPAddr{
			Endpoint: "localhost:1000",
		},
		ExtraSetting:     "some string",
		ExtraMapSetting:  nil,
		ExtraListSetting: nil,
	}
}

// CreateTraceReceiver creates a trace receiver based on this config.
func createTracesReceiver(
	_ context.Context,
	_ component.ReceiverCreateParams,
	cfg configmodels.Receiver,
	nextConsumer consumer.Traces,
) (component.TracesReceiver, error) {
	receiver := createReceiver(cfg)
	receiver.Traces = nextConsumer
	return receiver, nil
}

// CreateMetricsReceiver creates a metrics receiver based on this config.
func createMetricsReceiver(
	_ context.Context,
	_ component.ReceiverCreateParams,
	cfg configmodels.Receiver,
	nextConsumer consumer.Metrics,
) (component.MetricsReceiver, error) {
	receiver := createReceiver(cfg)
	receiver.Metrics = nextConsumer
	return receiver, nil
}

func createLogsReceiver(
	_ context.Context,
	_ component.ReceiverCreateParams,
	cfg configmodels.Receiver,
	nextConsumer consumer.Logs,
) (component.LogsReceiver, error) {
	receiver := createReceiver(cfg)
	receiver.Logs = nextConsumer

	return receiver, nil
}

func createReceiver(cfg configmodels.Receiver) *ExampleReceiverProducer {
	// There must be one receiver for all data types. We maintain a map of
	// receivers per config.

	// Check to see if there is already a receiver for this config.
	receiver, ok := exampleReceivers[cfg]
	if !ok {
		receiver = &ExampleReceiverProducer{}
		// Remember the receiver in the map
		exampleReceivers[cfg] = receiver
	}

	return receiver
}

// ExampleReceiverProducer allows producing traces and metrics for testing purposes.
type ExampleReceiverProducer struct {
	Started bool
	Stopped bool
	consumer.Traces
	consumer.Metrics
	consumer.Logs
}

// Start tells the receiver to start its processing.
func (erp *ExampleReceiverProducer) Start(_ context.Context, _ component.Host) error {
	erp.Started = true
	return nil
}

// Shutdown tells the receiver that should stop reception,
func (erp *ExampleReceiverProducer) Shutdown(context.Context) error {
	erp.Stopped = true
	return nil
}

// This is the map of already created example receivers for particular configurations.
// We maintain this map because the ReceiverFactory is asked trace and metric receivers separately
// when it gets CreateTracesReceiver() and CreateMetricsReceiver() but they must not
// create separate objects, they must use one Receiver object per configuration.
var exampleReceivers = map[configmodels.Receiver]*ExampleReceiverProducer{}
