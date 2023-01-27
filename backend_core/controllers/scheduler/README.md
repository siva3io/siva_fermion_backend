
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
## Scheduler

The Task Scheduler service allows you to perform automated tasks on a chosen computer. With this service, you can schedule any program to run at a convenient time for you or when a specific event occurs.

Go Task scheduler is a small library that you can use within your application that enables you to execute callbacks (goroutines) after a pre-defined amount of time. GTS also provides task storage which is used to invoke callbacks for tasks which couldn’t be executed during down-time as well as maintaining a history of the callbacks that got executed.

Execute tasks based after a specific duration or at a specific point in time.

Instantiate a scheduler as follows:
s := scheduler.New(storage)

Scheduling tasks can be done in 3 ways:
Execute a task after 5 seconds.
Execute a task at a specific time.
Execute a task every 1 minute.




Eunimart
