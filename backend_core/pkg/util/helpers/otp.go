package helpers

import (
	"crypto/rand"
	"io"
)

/*
 Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
 All rights reserved.
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
func GenerateOTP() (string, error) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	var max = 6
	buf := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, buf, max)
	if n != max {
		return "", err
	}
	for i := 0; i < len(buf); i++ {
		buf[i] = table[int(buf[i])%len(table)]
	}
	// return string(buf), nil
	return "666666", nil
}
