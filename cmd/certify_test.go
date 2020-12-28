/*
 * Copyright (C) 28/12/2020, 16:13 igors
 * This file is part of helmcertifier.
 *
 * helmcertifier is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * helmcertifier is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with helmcertifier.  If not, see <http://www.gnu.org/licenses/>.
 */

package cmd

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCertify(t *testing.T) {

	t.Run("uri is required", func(t *testing.T) {

		t.Run("flag -u not given", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			require.Error(t, cmd.Execute())
		})

		t.Run("flag -u is given but no value is informed", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			cmd.SetArgs([]string{"-u"})
			require.Error(t, cmd.Execute())
		})

		t.Run("flag -u and values are given", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			cmd.SetArgs([]string{"-u", "/tmp/chart.tgz"})
			require.NoError(t, cmd.Execute())
		})

	})

}
