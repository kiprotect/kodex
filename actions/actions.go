// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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

package actions

import (
	"github.com/kiprotect/kodex"
)

var Actions = kodex.ActionDefinitions{
	"undo": kodex.ActionDefinition{
		Name:  "Undo",
		Maker: MakeUndoAction,
		// to do: add form
	},
	"pseudonymize": kodex.ActionDefinition{
		Name:  "Pseudonymize",
		Maker: MakePseudonymizeAction,
		Form:  &PseudonymizeConfigForm,
	},
	"quantize": kodex.ActionDefinition{
		Name:  "Quantize",
		Maker: MakeQuantizeAction,
	},
	"generalize": kodex.ActionDefinition{
		Name:  "Generalize",
		Maker: MakeGeneralizeAction,
	},
	"form": kodex.ActionDefinition{
		Name:  "Form Validation",
		Maker: MakeFormAction,
		Form:  &FormForm,
	},
	"anonymize": kodex.ActionDefinition{
		Name:  "Anonymize",
		Maker: MakeAnonymizeAction,
	},
	"transcode": kodex.ActionDefinition{
		Name:  "Transcode",
		Maker: MakeTranscodeAction,
		// to do: add form
	},
	"drop": kodex.ActionDefinition{
		Name:  "Drop",
		Maker: MakeDropAction,
		// to do: add form
	},
}
