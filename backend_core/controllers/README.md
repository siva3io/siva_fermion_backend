
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
## App


In this app folder have all the core modules.

### Folder structure for modules
 
### 1 dto.go file
Rows to model, apply some business logic, save changes, map models to DTOs, and send a response back to the client. 
As for the difference between DTO and model: DTO is a representation of model for view and has no behaviors (methods). 
Model is the abstraction of your business logic and has a lot of complex behaviours.

### 2 handler.go file
Go handlers can be any struct that has a method named ServeHTTP with two parameters: an HTTPResponseWriter interface and a pointer to a Request struct. 
Handler functions are functions that behave like handlers. 
Handler functions have the same signature as the ServeHTTP method and are used to process requests.

### 3 route.go file
A route associates an HTTP verb (such as GET, POST, PUT, DELETE) and a URL path to a handler function. A router is an object which creates routes; i.e. it maps an HTTP request to a handler. 
The bunrouter is a fast and flexible HTTP router for Go. 
It supports middlewares, grouping routes, and flexible error handling.

### 4 service.go file
In computing systems, a service is a program that executes specific tasks in response to events or requests, such as:

HTTP requests
Message or stream from message broker
Time event (every hour, each 5 minutes, etc)
A service is designed to run forever and, like everything that is good eventually dies, it must be prepared to stop robustly and intelligently when the time comes.

### 5 validation.go file
A test of a system to prove that it meets all its specified requirements at a particular stage of its development. 
An activity that ensures that an end product stakeholders true needs and expectations are met.



Eunimart
