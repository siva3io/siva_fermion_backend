package helpers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)
/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License v3.0 as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License v3.0 for more details.
You should have received a copy of the GNU Lesser General Public License v3.0
along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/
func BASE64optimization(imageOptionsArray interface{}) map[string]interface{} {
	var arr []map[string]interface{}
	dto, _ := json.Marshal(imageOptionsArray)
	json.Unmarshal(dto, &arr)

	var image image.Image
	var err error
	var value []byte
	fileName := "optimization.jpg"

	for i, j := range arr {
		data := j["data"].(string)

		// decode base64 string
		value, err = base64.StdEncoding.DecodeString(data)

		// generate image object
		image, err = png.Decode(bytes.NewReader(value))
		if err != nil {
			image, err = jpeg.Decode(bytes.NewReader(value))
		}

		// resize image
		dstImage := imaging.Resize(image, 300, 0, imaging.Linear)

		// create image in local system
		imgFile, err := os.Create(fileName)
		defer imgFile.Close()

		if err != nil {
			arr[i]["data"] = data
			return map[string]interface{}{"status": false, "message": "Cannot create file in local system", "error": err, "data": arr}
		}
		png.Encode(imgFile, dstImage.SubImage(dstImage.Rect))

		// generate base64 for local image
		imgFile, err = os.Open(fileName)

		if err != nil {
			arr[i]["data"] = data
			return map[string]interface{}{"status": false, "message": "Unable to open file in local system", "error": err, "data": arr}
		}

		defer imgFile.Close()

		// create a new buffer base on file size
		fInfo, _ := imgFile.Stat()
		var size int64 = fInfo.Size()
		buf := make([]byte, size)

		// read file content into buffer
		fReader := bufio.NewReader(imgFile)
		fReader.Read(buf)

		// convert the buffer bytes to base64 string
		imgBase64Str := base64.StdEncoding.EncodeToString(buf)
		arr[i]["data"] = imgBase64Str
		// deleting physical file
		os.Remove(fileName)
	}

	return map[string]interface{}{"status": true, "data": arr}
}
