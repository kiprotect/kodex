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

package readers

import (
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex/writers"
)

var AMQPReaderForm = forms.Form{
	ErrorMsg: "invalid data encountered in the AMQP reader form",
	Fields: append([]forms.Field{
		{
			Name: "consumer",
			Validators: []forms.Validator{
				// empty consumer means the name will be auto-generated
				forms.IsOptional{Default: ""},
				forms.IsString{},
			},
		},
	}, writers.AMQPBaseForm.Fields...),
}
