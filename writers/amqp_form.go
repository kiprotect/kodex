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

package writers

import (
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
)

type TopicExchangeChosen struct{}

func (t TopicExchangeChosen) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {
	exchangeType := values["exchange_type"]
	if value.(bool) == true && exchangeType != "topic" {
		return nil, fmt.Errorf("exchange_type must be \"topic\" for this option")
	}
	return value, nil
}

var AMQPBaseForm = forms.Form{
	ErrorMsg: "invalid data encountered in the AMQP form",
	Fields: []forms.Field{
		{
			Name: "format",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "json"},
				forms.IsIn{Choices: []interface{}{"json"}},
			},
		},
		{
			Name: "compress",
			Validators: []forms.Validator{
				forms.IsOptional{Default: false},
				forms.IsBoolean{},
			},
		},
		{
			Name: "queue",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "routing_key",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "exchange",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "exchange_type",
			Validators: []forms.Validator{
				forms.IsIn{Choices: []interface{}{"fanout", "direct", "topic"}},
			},
		},
		{
			Name: "queue_expires_after_ms",
			Validators: []forms.Validator{
				forms.IsOptional{Default: int64(0)},
				forms.IsInteger{HasMin: true, Min: 0},
			},
		},
		{
			Name: "url",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
	},
}

var AMQPWriterForm = forms.Form{
	ErrorMsg: "invalid data encountered in the AMQP writer form",
	Fields: append([]forms.Field{
		{
			Name: "confirmation_timeout",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 1.0},
				forms.IsFloat{Convert: false},
			},
		},
	}, AMQPBaseForm.Fields...),
}
