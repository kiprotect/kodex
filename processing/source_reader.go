// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package processing

import (
	"github.com/kiprotect/kodex"
)

// An source reader reads items from an external source and forwards them
// to internal queues for processing, performing a fanout operation to
// distribute the items to all streams that read from the given source

type SourceSupervisor interface {
	ReaderStopped(SourceReader, kodex.SourceMap)
}

type SourceReader interface {
	SourceMap() kodex.SourceMap
	Start(SourceSupervisor, kodex.SourceMap) error
	Stop(bool) error
	Stopped() bool
	ID() []byte
}
