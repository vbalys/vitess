/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package mysqlctlclient contains the generic client side of the remote
// mysqlctl protocol.
package mysqlctlclient

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"

	"vitess.io/vitess/go/vt/log"
	"vitess.io/vitess/go/vt/servenv"
)

var protocol = "grpc"

func init() {
	servenv.OnParseFor("mysqlctl", registerFlags)
}

func registerFlags(fs *pflag.FlagSet) {
	fs.StringVar(&protocol, "mysqlctl_client_protocol", protocol, "the protocol to use to talk to the mysqlctl server")
}

// MysqlctlClient defines the interface used to send remote mysqlctl commands
type MysqlctlClient interface {
	// Start calls Mysqld.Start remotely.
	Start(ctx context.Context, mysqldArgs ...string) error

	// Shutdown calls Mysqld.Shutdown remotely.
	Shutdown(ctx context.Context, waitForMysqld bool) error

	// RunMysqlUpgrade calls Mysqld.RunMysqlUpgrade remotely.
	RunMysqlUpgrade(ctx context.Context) error

	// ApplyBinlogFile calls Mysqld.ApplyBinlogFile remotely.
	ApplyBinlogFile(ctx context.Context, binlogFileName, binlogRestorePosition string) error

	// ReinitConfig calls Mysqld.ReinitConfig remotely.
	ReinitConfig(ctx context.Context) error

	// RefreshConfig calls Mysqld.RefreshConfig remotely.
	RefreshConfig(ctx context.Context) error

	// Close will terminate the connection. This object won't be used anymore.
	Close()
}

// Factory functions are registered by client implementations.
type Factory func(network, addr string) (MysqlctlClient, error)

var factories = make(map[string]Factory)

// RegisterFactory allows a client implementation to register itself
func RegisterFactory(name string, factory Factory) {
	if _, ok := factories[name]; ok {
		log.Fatalf("RegisterFactory %s already exists", name)
	}
	factories[name] = factory
}

// New creates a client implementation as specified by a flag.
func New(network, addr string) (MysqlctlClient, error) {
	factory, ok := factories[protocol]
	if !ok {
		return nil, fmt.Errorf("unknown mysqlctl client protocol: %v", protocol)
	}
	return factory(network, addr)
}
