

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
