<!--
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
-->
## About Middleware

A middleware handler is simply an http. Handler that wraps another http. 
Handler to do some pre- and/or post-processing of the request. 
It's called "middleware" because it sits in the middle between the Go web server and the actual handler

Because Go comes with and expects this Handler interface in it's most crucial functions, 
it's very easy and reliable to build packages around it and extend 
the interface with wrappers, such as middleware.


## How to write Go middleware

Reading the Request. All of the middlewares in our examples will accept an http.Handler as an argument, and return an http.Handler.
Modifying the Request. 
Let's say we want to add a header to the request, or otherwise modify it.
Writing Response Headers.
Create Middleware and Implement Logic In Middleware and Route and Method In Controller

## Authorization with Golang

Adding authorization will allow you to protect your API. 
Since your app deals with projects that are in active development, you don’t want any data to be publicly available.

You've already accomplished the first step, which requires the user to sign in to your application. The next step is to pull the data from the Go application, but only if the user has a valid access token.

## Custom Errors in Golang

In Go, we can create custom errors by implementing an error interface in a struct. Here, the Error() method returns an error message in string form if there is an error. Otherwise, it returns nil . Now, to create a custom error, we have to implement the Error() method on a Go struct

Eunimart
