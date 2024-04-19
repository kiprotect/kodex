// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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

package readers

import (
	"github.com/kiprotect/kodex"
)

var Readers = kodex.ReaderDefinitions{
	"file": kodex.ReaderDefinition{
		Maker:    MakeFileReader,
		Form:     FileReaderForm,
		Internal: true,
	},
	"stdin": kodex.ReaderDefinition{
		Maker:    MakeStdinReader,
		Form:     StdinReaderForm,
		Internal: true,
	},
	"generate": kodex.ReaderDefinition{
		Maker:    MakeGenerateReader,
		Form:     GenerateForm,
		Internal: true,
	},
	"bytes": kodex.ReaderDefinition{
		Maker:    MakeBytesReader,
		Form:     BytesReaderForm,
		Internal: true,
	},
	"amqp": kodex.ReaderDefinition{
		Maker:    MakeAMQPReader,
		Form:     AMQPReaderForm,
		Internal: false,
	},
}
