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
	"time"
)

func ProcessStream(stream kodex.Stream, timeout time.Duration) error {

	// we get all the sources for the stream
	sourceMaps, _ := stream.Sources()

	// we create readers for all the sources
	sourceReaders := make([]SourceReader, 0)
	for _, sourceMap := range sourceMaps {
		sourceReader := MakeLocalSourceReader(1, []byte("test"))
		if err := sourceReader.Start(nil, sourceMap); err != nil {
			return err
		}
		sourceReaders = append(sourceReaders, sourceReader)
	}

	// we process the stream using a local stream executor
	streamExecutor := MakeLocalStreamExecutor(4, []byte("test"))

	if err := streamExecutor.Start(nil, stream); err != nil {
		return err
	}

	// we process all destinations using local destination writers
	destinationWriters := make([]DestinationWriter, 0)
	destinationMaps := make([]kodex.DestinationMap, 0)
	configs, err := stream.Configs()

	if err != nil {
		return err
	}

	for _, config := range configs {
		configDestinations, err := config.Destinations()
		if err != nil {
			return err
		}

		for _, configDestinationMaps := range configDestinations {
			destinationMaps = append(destinationMaps, configDestinationMaps...)
		}
	}

	for _, destinationMap := range destinationMaps {
		destinationWriter := MakeLocalDestinationWriter(4, []byte("test"))
		if err := destinationWriter.Start(nil, destinationMap); err != nil {
			return err
		}
		destinationWriters = append(destinationWriters, destinationWriter)
	}

	// we wait for all source readers to finish their work
	for {
		allStopped := true
		for _, sourceReader := range sourceReaders {
			if !sourceReader.Stopped() {
				allStopped = false
				break
			}
		}
		if allStopped {
			break
		}
		time.Sleep(time.Millisecond)
	}

	// we wait for the stream executor to finish its work
	for {
		if streamExecutor.Stopped() {
			break
		}
		time.Sleep(time.Millisecond)
	}

	// we wait for all destination writers to finish their work
	for {
		allStopped := true
		for _, destinationWriter := range destinationWriters {
			if !destinationWriter.Stopped() {
				allStopped = false
				break
			}
		}
		if allStopped {
			break
		}
		time.Sleep(time.Millisecond)
	}

	return nil

}
