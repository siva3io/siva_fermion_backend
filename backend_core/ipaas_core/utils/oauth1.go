package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"net/url"
	"strconv"
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
type OAuth1 struct {
	ConsumerKey     string
	ConsumerSecret  string
	AccessToken     string
	TokenSecret     string
	SignatureMethod string
	Version         string
}

// Params being any key-value url query parameter pairs
func (auth OAuth1) BuildOAuth1Header(method, path string) string {
	vals := url.Values{}
	vals.Add("oauth_nonce", generateNonce())
	vals.Add("oauth_consumer_key", auth.ConsumerKey)
	vals.Add("oauth_signature_method", auth.SignatureMethod)
	vals.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	vals.Add("oauth_token", auth.AccessToken)
	vals.Add("oauth_version", auth.Version)

	urlWithQuery := strings.Split(path, "?")
	if len(urlWithQuery) == 2 {
		query := strings.Split(path, "?")[1]
		for _, i := range strings.Split(query, "&") {
			x := strings.Split(i, "=")
			vals.Add(x[0], x[1])
		}
	}

	parameterString := strings.Replace(vals.Encode(), "+", "%20", -1)

	signatureBase := strings.ToUpper(method) + "&" + url.QueryEscape(strings.Split(path, "?")[0]) + "&" + url.QueryEscape(parameterString)
	signingKey := url.QueryEscape(auth.ConsumerSecret) + "&" + url.QueryEscape(auth.TokenSecret)

	signature := ""
	if auth.SignatureMethod == "HMAC-SHA1" {
		signature = calculateSignatureSHA1(signatureBase, signingKey)
	} else if auth.SignatureMethod == "HMAC-SHA256" {
		signature = calculateSignatureSHA256(signatureBase, signingKey)
	}

	vals.Add("oauth_signature", signature)

	return "OAuth oauth_consumer_key=\"" + url.QueryEscape(vals.Get("oauth_consumer_key")) + "\", oauth_nonce=\"" + url.QueryEscape(vals.Get("oauth_nonce")) +
		"\", oauth_signature=\"" + url.QueryEscape(signature) + "\", oauth_signature_method=\"" + url.QueryEscape(vals.Get("oauth_signature_method")) +
		"\", oauth_timestamp=\"" + url.QueryEscape(vals.Get("oauth_timestamp")) + "\", oauth_token=\"" + url.QueryEscape(vals.Get("oauth_token")) +
		"\", oauth_version=\"" + url.QueryEscape(vals.Get("oauth_version")) + "\""
}

func calculateSignatureSHA1(base, key string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(base))
	signature := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

func calculateSignatureSHA256(base, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(base))
	signature := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

func generateNonce() string {
	rand.Seed(time.Now().UnixNano())
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 48)
	for i := range b {
		b[i] = allowed[rand.Intn(len(allowed))]
	}
	return string(b)
}

func (f *Functions) Oauth(ConsumerKey, ConsumerSecret, AccessToken, TokenSecret, SignatureMethod, Version string, requestData map[string]interface{}, featureSessionVariables []model.KeyValuePair) string {

	auth := OAuth1{
		ConsumerKey:     ConsumerKey,
		ConsumerSecret:  ConsumerSecret,
		AccessToken:     AccessToken,
		TokenSecret:     TokenSecret,
		SignatureMethod: SignatureMethod,
		Version:         Version,
	}
	oauth := auth.BuildOAuth1Header(requestData["method"].(string), requestData["url"].(string))

	return oauth
}
