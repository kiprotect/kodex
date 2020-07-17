// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
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

package writers

import (
	"github.com/kiprotect/kiprotect"
)

var Writers = kiprotect.WriterDefinitions{
	"file": kiprotect.WriterDefinition{
		Maker:    MakeFileWriter,
		Form:     FileWriterForm,
		Internal: true,
	},
	"http": kiprotect.WriterDefinition{
		Maker: MakeHTTPWriter,
		Form:  HTTPWriterForm,
	},
	"bytes": kiprotect.WriterDefinition{
		Maker:    MakeBytesWriter,
		Form:     BytesWriterForm,
		Internal: true,
	},
	"in-memory": kiprotect.WriterDefinition{
		Maker:    MakeInMemoryWriter,
		Internal: true,
	},
	"stdout": kiprotect.WriterDefinition{
		Maker:    MakeStdoutWriter,
		Internal: true,
	},
	"amqp": kiprotect.WriterDefinition{
		Maker:    MakeAMQPWriter,
		Form:     AMQPWriterForm,
		Internal: false,
	},
	"count": kiprotect.WriterDefinition{
		Maker:    MakeCountWriter,
		Form:     CountWriterForm,
		Internal: true,
	},
}
