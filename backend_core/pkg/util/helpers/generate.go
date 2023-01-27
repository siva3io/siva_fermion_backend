package helpers

import (
	"bytes"
	crypto_rand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"net/smtp"
	"strconv"
	"text/template"
	"time"

	"fermion/backend_core/db"
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
func GenerateSequence(moduleName string, userid string, modelName string) string {
	sequence := moduleName + "/" + userid + "/000"
	count := 0

	query := "SELECT COUNT(*) FROM " + modelName
	count = GetCount(query)
	count++
	sequence += strconv.Itoa(count)
	return sequence
}
func GetCount(query string) int {
	count := 0
	database := db.DbManager()
	database.Raw(query).Scan(&count)
	return count
}
func GenerateRandID(min int, max int, prefix string) string {
	rand.Seed(time.Now().UnixNano())
	if min == 0 {
		min = 1000
	}
	if max == 0 {
		max = 9999
	}
	randomNum := rand.Intn(max-min+1) + min
	randomId := strconv.Itoa(randomNum)
	return prefix + randomId
}
func VariantGeneration(array [][]map[string]interface{}) ([]interface{}, error) {

	/*
		{
			{{"id":1,"name":"S"},{"id":2,"name":"M"},{"id":3,"name":"L"}},
			{{"id":4,"name":"Red"},{"id":5,"name":"Blue"}},
			{{"id":6,"name":"Cotton"},{"id":7,"name":"Fabric"}},
			{{"id":8,"name":"Plain"},{"id":9,"name":"Checked"}}
		}
	*/

	length := len(array)
	var combinations []interface{}
	indices := make([]int, length)
	for i := 0; i < length; i++ {
		indices[i] = 0
	}
	for {
		var combination []interface{}
		for i := 0; i < length; i++ {
			j := indices[i]
			combination = append(combination, array[i][j])
		}
		combinations = append(combinations, map[string]interface{}{"combination": combination})
		next := length - 1

		for next >= 0 && (indices[next]+1) >= len(array[next]) {
			next = next - 1
		}
		if next < 0 {
			break
		}
		indices[next]++

		for i := next + 1; i < length; i++ {
			indices[i] = 0
		}
	}
	return combinations, nil
}
func GenerateOTP() (string, error) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	var max = 6
	buf := make([]byte, max)
	n, err := io.ReadAtLeast(crypto_rand.Reader, buf, max)
	if n != max {
		return "", err
	}
	for i := 0; i < len(buf); i++ {
		buf[i] = table[int(buf[i])%len(table)]
	}
	// return string(buf), nil
	return "666666", nil
}
func SendEmail(subject string, receiverEmail string, templateLocation string, data interface{}) error {

	// Sender data.
	from := "testhrmail1@gmail.com"
	password := "@Scientistcr7"

	// Receiver email address.
	to := []string{
		receiverEmail,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, err := template.ParseFiles(templateLocation)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body.Write([]byte(fmt.Sprintf("Subject: "+subject+" \n%s\n\n", mimeHeaders)))
	err = t.Execute(&body, data)

	if err != nil {
		return err
	}

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent!")
	return nil
}
