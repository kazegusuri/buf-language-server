// Copyright 2022 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buflscli

import (
	"context"
	"io"

	"github.com/bufbuild/buf-language-server/internal/bufls"
	"github.com/bufbuild/buf/private/buf/bufcli"
	"github.com/bufbuild/buf/private/buf/bufctl"
	"github.com/bufbuild/buf/private/pkg/app/appext"

	"go.lsp.dev/jsonrpc2"
	"go.uber.org/multierr"
)

// NewEngine returns a new bufls.Engine.
func NewEngine(
	ctx context.Context,
	container appext.Container,
	disableSymlinks bool,
) (bufls.Engine, error) {
	controller, err := bufcli.NewController(
		container,
		bufctl.WithDisableSymlinks(disableSymlinks),
	)
	if err != nil {
		return nil, err
	}

	return bufls.NewEngine(
		container.Logger(),
		container,
		controller,
	), nil
}

// NewConn returns a new jsonrpc2.Conn backed by the given io.{Read,Write}Closer
// (which is usually os.Stdin and os.Stdout).
func NewConn(readCloser io.ReadCloser, writeCloser io.WriteCloser) jsonrpc2.Conn {
	return jsonrpc2.NewConn(
		jsonrpc2.NewStream(
			&readWriteCloser{
				readCloser:  readCloser,
				writeCloser: writeCloser,
			},
		),
	)
}

type readWriteCloser struct {
	readCloser  io.ReadCloser
	writeCloser io.WriteCloser
}

func (r *readWriteCloser) Read(b []byte) (int, error) {
	return r.readCloser.Read(b)
}

func (r *readWriteCloser) Write(b []byte) (int, error) {
	return r.writeCloser.Write(b)
}

func (r *readWriteCloser) Close() error {
	return multierr.Append(r.readCloser.Close(), r.writeCloser.Close())
}
