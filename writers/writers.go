// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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

package writers

import (
	"github.com/kiprotect/kodex"
)

var Writers = kodex.WriterDefinitions{
	"file": kodex.WriterDefinition{
		Maker:    MakeFileWriter,
		Form:     FileWriterForm,
		Internal: true,
	},
	"http": kodex.WriterDefinition{
		Maker: MakeHTTPWriter,
		Form:  HTTPWriterForm,
	},
	"bytes": kodex.WriterDefinition{
		Maker:    MakeBytesWriter,
		Form:     BytesWriterForm,
		Internal: true,
	},
	"in-memory": kodex.WriterDefinition{
		Maker:    MakeInMemoryWriter,
		Internal: true,
	},
	"stdout": kodex.WriterDefinition{
		Maker:    MakeStdoutWriter,
		Internal: true,
	},
	"amqp": kodex.WriterDefinition{
		Maker:    MakeAMQPWriter,
		Form:     AMQPWriterForm,
		Internal: false,
	},
	"count": kodex.WriterDefinition{
		Maker:    MakeCountWriter,
		Form:     CountWriterForm,
		Internal: true,
	},
}
