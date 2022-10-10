package helpers

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
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
