package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"fermion/backend_core/ipaas_core/model"
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

func (f *Functions) AmazonHmacSignatureV4(inputs map[string]interface{}, request_data map[string]interface{}, featureSessionVariables []model.KeyValuePair) string {

	// fmt.Println("\n\n\n--------------------featureSessionVariables--------------------", featureSessionVariables)

	dateTime := time.Now().UTC()

	credentialScope := dateTime.Format("20060102") + "/" + inputs["region"].(string) + "/execute-api/aws4_request"

	fmt.Println("UTC Datetimes ------>", dateTime.Format("20060102T150405"))

	authString := generateAuthString(dateTime, inputs["aws_access_key"].(string), inputs["region"].(string))

	fmt.Println("\nAuth String generated --------->", authString)

	hash := generateHMAC(inputs["aws_secret"].(string), dateTime, inputs["region"].(string))

	fmt.Println("\nHashed data ----------------->", hash)

	stringPayload := ""

	if request_data["payload"] != nil {
		jsonPayload, _ := json.Marshal(request_data["payload"].(map[string]interface{}))
		stringPayload = string(jsonPayload)
	}

	if stringPayload == "{}" {
		stringPayload = ""
	}

	fmt.Println("\n stringBytePayload ------------>", stringPayload)

	hashedPayload := getHash(stringPayload)

	fmt.Println("\nhashed payload ------------->", hashedPayload)

	canonicalString := getHashedCanonicalString(request_data, hashedPayload)

	fmt.Println("\ncanonicalSignature -------->", canonicalString)

	stringToSign := inputs["algorithm"].(string) + "\n" + dateTime.Format("20060102T150405Z") + "\n" + credentialScope + "\n" + canonicalString

	fmt.Println("\n finalStringTosign ----------------->", stringToSign)

	signatureWithStringToSign := getHMAC(hash, []byte(stringToSign))

	fmt.Println("\n hashedStringToSign ------------------>", hex.EncodeToString(signatureWithStringToSign))

	authString += ", Signature=" + hex.EncodeToString(signatureWithStringToSign)

	fmt.Println("\n authorizationSignature -------------------->", authString)

	return authString

}

func extractPath(url string) (string, map[string]interface{}) {

	queryParams := make(map[string]interface{}, 0)

	urlWithQuery := strings.Split(url, "?")
	if len(urlWithQuery) == 2 {
		query := urlWithQuery[1]
		for _, i := range strings.Split(query, "&") {
			x := strings.Split(i, "=")
			queryParams[x[0]] = x[1]
		}
	}

	url = urlWithQuery[0]
	urlList := strings.Split(url, "/")
	path := "/" + strings.Join(urlList[3:], "/")

	return path, queryParams
}

func getHashedCanonicalString(requestData map[string]interface{}, hashedPayload string) string {

	jsonData, _ := json.MarshalIndent(requestData, "", "\t")
	fmt.Println("\n\n-------------------------------------", string(jsonData))

	canonicalString := requestData["method"].(string) + "\n"
	path, queryParams := extractPath(requestData["url"].(string))
	canonicalString += path + "\n"

	if len(queryParams) == 0 {
		canonicalString += "\n"
	} else {
		queryString := urlEncode(queryParams)
		canonicalString += queryString + "\n"
	}

	for _, header := range requestData["headers"].([]model.KeyValuePair) {
		canonicalString += header.Key + ":" + header.Value.(string) + "\n"
	}

	canonicalString += "\n"

	for index, header := range requestData["headers"].([]model.KeyValuePair) {
		canonicalString += header.Key
		if len(requestData["headers"].([]model.KeyValuePair))-1 != index {
			canonicalString += ";"
		}
	}

	canonicalString += "\n" + hashedPayload
	fmt.Println("\ncanonicalStringToSign ----------------->", canonicalString)

	hashedCanonicalString := getHash(canonicalString)
	return hashedCanonicalString
}

func urlEncode(data map[string]interface{}) string {

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fmt.Println("\nSorted Object keys ------------>", keys)
	queryString := ""
	keysLength := len(data)
	for index, paramKey := range keys {
		queryString += paramKey + "=" + data[paramKey].(string)
		if keysLength-1 != index {
			queryString += "&"
		}
	}
	fmt.Println("\nUrl Encoded QueryParams ------->", queryString)
	return queryString

}

func generateHMAC(secret string, dateTime time.Time, region string) []uint8 {

	hash := getHMAC([]byte("AWS4"+secret), []byte(dateTime.Format("20060102")))
	hash = getHMAC(hash, []byte(region))
	hash = getHMAC(hash, []byte("execute-api"))
	hash = getHMAC(hash, []byte("aws4_request"))
	return hash

}

func generateAuthString(dateTime time.Time, credential string, region string) string {

	authString := "AWS4-HMAC-SHA256 Credential=" + credential + "/"
	authString += dateTime.Format("20060102" /*T150405Z"*/) + "/"
	authString += region + "/"
	authString += "execute-api/"
	authString += "aws4_request,"
	authString += " SignedHeaders=content-type;host;x-amz-access-token;x-amz-content-sha256;x-amz-date"
	return authString
}

func getHMAC(key []byte, data []byte) []byte {

	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	hashedValue := hash.Sum(nil)
	return hashedValue

}

func getHash(key string) string {

	h := sha256.New()
	h.Write([]byte(key))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash

}

func (f *Functions) GenerateSignatureHMAC(inputs map[string]interface{}, request_data map[string]interface{}, featureSessionVariables []model.KeyValuePair) string {

	api, keyValues := extractPath(request_data["url"].(string))
	secret := inputs["secret"].(string)
	api = api[5:]
	fmt.Println(keyValues)
	fmt.Println(secret)

	keys := make([]string, 0, len(keyValues))
	for k := range keyValues {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	concatString := ""
	for _, k := range keys {
		fmt.Println(k, keyValues[k].(string))
		val, _ := url.QueryUnescape(keyValues[k].(string))
		concatString += k + val
	}
	concatString = api + concatString
	fmt.Println(concatString)

	// hmacConcatString := hmac.New(sha256.New, []byte(secret))
	// hmacConcatString.Write([]byte(concatString))
	// hmac := hex.EncodeToString(hmacConcatString.Sum(nil))

	hmac := hex.EncodeToString(getHMAC([]byte(secret), []byte(concatString)))
	fmt.Println(strings.ToUpper(hmac))
	return strings.ToUpper(hmac)
}
