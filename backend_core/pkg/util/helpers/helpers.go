package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"fermion/backend_core/db"
	model_core "fermion/backend_core/internal/model/core"
	cache "fermion/backend_core/pkg/util/cache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-mail/mail"
)

const (
	AWS_S3_REGION = "ap-southeast-1"
	AWS_S3_BUCKET = "dev-api-files"
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

func Contains[T comparable](inputArray []T, searchKey T) bool {
	for _, intputKey := range inputArray {
		if intputKey == searchKey {
			return true
		}
	}
	return false
}
func BasicAuth(username string, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
func TokenBearerType(token string) (string, error) {
	if token == "" {
		return "", errors.New("token must not be empty")
	}
	output := "Bearer " + token
	return output, nil
}
func GetDateTime(isUtc bool, format string, delay float64) string {

	dateTime := time.Now()

	if isUtc {
		dateTime = dateTime.UTC()
	}

	dateTime = dateTime.Add(time.Hour*time.Duration(0) +
		time.Minute*time.Duration(delay) +
		time.Second*time.Duration(0))

	formatedDateTime := dateTime.Format(format)

	return formatedDateTime
}
func Mod(input1 int, input2 int) int {
	return input1 % input2
}
func ReturnURLEncodeString(input interface{}) (string, error) {
	data := url.Values{}
	var inputconversion map[string]interface{}

	jsondata, err := json.Marshal(input)
	if err != nil {
		return "", errors.New("invalid conversion(map[string]interface{})")
	}
	err = json.Unmarshal(jsondata, &inputconversion)
	if err != nil {
		return "", errors.New("invalid conversion(map[string]interface{})")
	}

	for i := range inputconversion {
		data.Set(i, fmt.Sprintf("%v", inputconversion[i]))
	}
	output := data.Encode()

	fmt.Println("ReturnURLEncodeString --->>>>", output)
	return output, nil
}
func PrettyPrint(key string, data interface{}) {
	byteData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("=============> " + key + " <============= ")
	fmt.Println(string(byteData))
}
func JsonMarshaller(input interface{}, output interface{}) error {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, output)
	if err != nil {
		return err
	}
	return nil
}
func ReadFile(path string) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	fileDataByte, ioErr := ioutil.ReadFile(path)
	if ioErr != nil {
		fmt.Println("--------> Read File error <----------------", ioErr)
		return data, ioErr
	}

	if err := json.Unmarshal([]byte(fileDataByte), &data); err != nil {
		fmt.Println("--------> Error while parsing Readfile <----------------" + path)
		return nil, err
	}
	// fmt.Println("ReadFile --->>>>", data)
	return data, nil
}

func GetStructMethods(structName interface{}) []string {

	var structMethods []string
	structType := reflect.TypeOf(structName)
	for i := 0; i < structType.NumMethod(); i++ {
		structMethods = append(structMethods, structType.Method(i).Name)
	}
	return structMethods
}

func GetAppFromCache(appCode string) (map[string]interface{}, error) {
	var appData map[string]interface{}
	cachedAppData := cache.GetCacheKey(appCode)
	if cachedAppData != nil {
		err := JsonMarshaller(cachedAppData, &appData)
		if err != nil {
			return nil, err
		}
		fmt.Println("-------------FROM CACHE--------------")
		return appData, nil
	}
	db := db.DbManager()
	err := db.Model(&model_core.InstalledApps{}).Where("code", appCode).First(&appData).Error
	if err != nil {
		return nil, err
	}
	JsonMarshaller(appData, &appData)
	fmt.Println("-------------FROM DB--------------")
	return appData, nil
}

func UpsertAppAuthToken(appData map[string]interface{}, payload map[string]interface{}) error {

	appAuthToken, _ := TokenBearerType(appData["access_token"].(string))

	fieldsName := `[`
	for key := range payload {
		fieldsName += `"` + key + `",`
	}
	fieldsName = fieldsName[:len(fieldsName)-1] + `]`

	apiResponse, err := MakeRequest(Request{
		Method: "GET",
		Scheme: "http",
		Host:   "localhost:3031",
		Path:   "/api/v1/omnichannel_fields/app_data/filter",
		Header: map[string]string{
			"Authorization": appAuthToken,
		},
		Params: map[string]string{
			"fields": fieldsName,
			"app_id": fmt.Sprint(appData["id"]),
		},
	})

	if err != nil {
		return err
	}

	responseObj := apiResponse.(map[string]interface{})

	if responseObj["data"] == nil {
		return errors.New("API ERROR : /api/v1/omnichannel_fields/app_data/filter")
	}

	fields := make([]interface{}, len(payload))

	var ind int
	for pKey, pValue := range payload {
		for rKey, rValue := range responseObj["data"].(map[string]interface{}) {
			if pKey == rKey {
				fields[ind] = map[string]interface{}{
					"omnichannel_field_id": rValue.(map[string]interface{})["id"],
					"data":                 pValue,
				}
			}
		}
		ind++
	}

	upsertFieldsPayload := map[string]interface{}{
		"app_code": appData["code"],
		"fields":   fields,
	}

	apiResponse, err = MakeRequest(Request{
		Method: "POST",
		Scheme: "http",
		Host:   "localhost:3031",
		Path:   "/api/v1/omnichannel_fields/upsert_data",
		Header: map[string]string{
			"Authorization": appAuthToken,
			"Content-Type":  "application/json",
		},
		Body: upsertFieldsPayload,
		Params: map[string]string{
			"register": "false",
		},
	})

	if err != nil {
		return err
	}

	PrettyPrint("payload", upsertFieldsPayload)

	PrettyPrint("response", apiResponse)

	return nil
}

func GetAppAuthFields(appData map[string]interface{}, fields []string) (map[string]interface{}, error) {

	appAuthToken, _ := TokenBearerType(appData["access_token"].(string))

	fieldsName := `["` + strings.Join(fields, `","`) + `"]`

	apiResponse, err := MakeRequest(Request{
		Method: "GET",
		Scheme: "http",
		Host:   "localhost:3031",
		Path:   "/api/v1/omnichannel_fields/app_data/filter",
		Header: map[string]string{
			"Authorization": appAuthToken,
		},
		Params: map[string]string{
			"fields": fieldsName,
			"app_id": fmt.Sprint(appData["id"]),
		},
	})

	if err != nil {
		return nil, err
	}

	responseObj := apiResponse.(map[string]interface{})

	if responseObj["data"] == nil {
		return nil, errors.New("API ERROR : /api/v1/omnichannel_fields/app_data/filter")
	}

	return responseObj["data"].(map[string]interface{}), nil
}

func LogOutput(logFile string) func() {

	f, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	out := os.Stdout
	mw := io.MultiWriter(out, f)
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	log.SetOutput(mw)

	exit := make(chan bool)
	go func() {
		_, _ = io.Copy(mw, r)
		exit <- true
	}()
	return func() {
		_ = w.Close()
		<-exit
		_ = f.Close()
	}
}

func ShallowCopy[T any](inputArray []T) []T {
	shallowCopy := make([]T, len(inputArray))
	copy(shallowCopy, inputArray)
	return shallowCopy
}

func TimeIt(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("Time taken to run %v : %vs\n", name, time.Since(start).Seconds())
	}
}

func UploadFile(uploadFileDir string, path string) error {
	session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		log.Fatal(err)
	}
	upFile, err := os.Open(uploadFileDir)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AWS_S3_BUCKET),
		Key:                  aws.String(path + uploadFileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func SendMail(emailtype string, orderid string, companyname string, reciveremail string, message string, data string) error {
	resp, err := http.Get(data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("backend_core/pkg/util/helpers/file.pdf", imageData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	m := mail.NewMessage()

	m.SetHeader("From", "noreply@eunimart.com")

	m.SetHeader("To", reciveremail)
	m.SetBody("text/html", message)
	m.SetHeader("Subject", emailtype+" : "+orderid+" | "+companyname)

	m.Attach("backend_core/pkg/util/helpers/file.pdf")
	d := mail.NewDialer("smtp.gmail.com", 587, "noreply@eunimart.com", "kyhooykulwnoakdy")

	if err := d.DialAndSend(m); err != nil {

		return err

	}

	defer os.Remove("backend_core/pkg/util/helpers/file.pdf")

	return nil

}
