

## Routs Description

Implementation of router in Golang, the standard package net/http has a routing function. That feature is called multiplexer in Golang. 
However, the standard package does not support path parameters, 
so if you want to use it, you need to prepare an external package or extend the standard multiplexer

A route associates an HTTP verb (such as GET, POST, PUT, DELETE) and a URL path to a handler function. A router is an object which creates routes; i.e. it maps an HTTP request to a handler. 
The bunrouter is a fast and flexible HTTP router for Go. 
It supports middlewares, grouping routes, and flexible error handling.


Eunimart
