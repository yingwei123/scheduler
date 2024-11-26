
# Scheduler Package

The `scheduler` package provides a lightweight task scheduling framework for Go applications. It allows you to schedule tasks at specified intervals with options for repeating execution, graceful shutdown, and handling system signals.

## Features

- Schedule tasks at specified intervals.
- Support for repeating and one-time tasks.
- Graceful shutdown with cleanup using context.
- Signal handling for `SIGINT` and `SIGTERM` for smooth termination.
- Thread-safe task addition, removal, and execution.

---

## Installation

Install the package using:

```bash
go get github.com/yingwei123/scheduler
```

---

## Usage

### Create a Scheduler

Create a new instance of the scheduler using:

```go
import "github.com/yingwei123/scheduler"

s := scheduler.NewScheduler()
```

### Define a Task

Create a struct that implements the `Task` interface:

```go
type MyTask struct{} //add in values needed for the execute task

//Execute performs the task in the scheduler.
func (t *MyTask) Execute() error {
    fmt.Println("Executing MyTask!")
    return nil
}
```

### Add Tasks

Add tasks to the scheduler using the `AddTask` method:

```go
taskName := "task1"
myTask := &MyTask{}
scheduledTask := scheduler.ScheduledTask{
    Interval: 2 * time.Second, // Execute every 2 seconds
    Task:     myTask,
    Repeat:   true,            // Repeat execution
}

if err := s.AddTask(scheduledTask, taskName); err != nil {
    log.Fatalf("Failed to add task: %v", err)
}
```

### Start the Scheduler

Start the scheduler and begin executing tasks:

```go
go s.Start()
```

### Stop the Scheduler

Stop the scheduler gracefully:

```go
s.Stop()
```

### Remove a Task

Remove a specific task by name:

```go
s.RemoveTask("task1")
```

---

## Example

Below is a complete example that demonstrates how to use the `scheduler` package:

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/yingwei123/scheduler"
)

type MyTask struct{}

func (t *MyTask) Execute() error {
    fmt.Println("Executing MyTask!")
    return nil
}

func main() {
    s := scheduler.NewScheduler()

    // Define and add tasks
    myTask := &MyTask{}
    scheduledTask := scheduler.ScheduledTask{
        Interval: 2 * time.Second,
        Task:     myTask,
        Repeat:   true,
    }
    if err := s.AddTask(scheduledTask, "task1"); err != nil {
        log.Fatalf("Failed to add task: %v", err)
    }

    // Start the scheduler
    go s.Start()

    // Run for 10 seconds before stopping
    time.Sleep(10 * time.Second)
    s.Stop()
}
```

---

## Graceful Shutdown

The scheduler listens for OS signals (`SIGINT` and `SIGTERM`) and stops gracefully, ensuring all tasks are completed or stopped safely.

To trigger the graceful shutdown, use `Ctrl+C` or send a termination signal.

---

## Testing

Unit tests are provided for this package. To run the tests, use:

```bash
go test -v github.com/yingwei123/scheduler
```

---

## License

This project is open-source and available under the [MIT License](LICENSE.md).

---

## Contributions

Feel free to open issues or submit pull requests if you find any bugs or want to add features. Contributions are welcome!