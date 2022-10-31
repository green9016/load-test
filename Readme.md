# Load test assessment #

Build a load test for the endpoint `https://jsonplaceholder.typicode.com/todos/1`
And store the count of response code like as following
```
{
  200: 5000,
  429: 1234,
  501: 1
}
```

## How implemented ##

Implemented using N requesters, 1 writer, 1 reporter.
- Requester: request the endpoint and inform the statusCode of the response to the writer using resultCh channel.
- Writer: get the statusCode from the resultCh channel and update the result based on it.
- Reporter: read the result and print it every 1 second.

Where N is the count of cpu processors.
And you can stop running by pressing CTRL+C.

### File structure ###
- load-test.go : main file
- workers.go : workers file to manager requesters, reporter and writer thread
- worker.go : worker file included the implementation of requester, reporter and writer

### Diagram ###
https://docs.google.com/drawings/d/1HF7L3WJMzS94nUBZX2Lxs-ifKGEfBoMLNZgN-JAExu4/edit?usp=sharing


## How to run ##

```
go run *.go
```

### Testing Result ###
```
Loading test:
Loading test: 200: 1
Loading test: 200: 10
```