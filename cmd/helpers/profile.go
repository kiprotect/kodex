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

package helpers

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"os"
	"runtime"
	"runtime/pprof"
)

func runWithProfiler(name string, runner func() error) error {

	kodex.Log.Info("Running with profiler...")

	fc, err := os.Create(fmt.Sprintf("%s-cpu.pprof", name))

	if err != nil {
		return err
	}

	if err := pprof.StartCPUProfile(fc); err != nil {
		return err
	}

	defer pprof.StopCPUProfile()

	runnerErr := runner()

	fm, err := os.Create(fmt.Sprintf("%s-mem.pprof", name))

	if err != nil {
		return err
	}

	runtime.GC() // get up-to-date statistics

	if err := pprof.WriteHeapProfile(fm); err != nil {
		return err
	}

	fm.Close()

	return runnerErr
}
