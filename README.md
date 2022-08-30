### Orderbook

The orderbook implements the required funcionality.

I used a self-balancing binary tree and hash table to store orders for efficient lookup.
Operations on the tree, e.g. lookup the spread by finding max and min orders, have a time complexity of O(log n).
Lookups, adding and removing orders in the hash table has a time complexity of O(1).

### Parsing

The input file gets parsed as a byte stream and an order book gets build from the snapshot.
l2updates get applied to the book as per task definition and the spread is printed after each update.
Note, the JSON structure in the sample stream/file differs from the example format in the task desciption.
I implemented support for the structure in the file.
The output format however, is as required in the task description.

### Data integrity

The parser listens for an interrupt signal and shuts down only before or after an update is fully processed so we don't have partial updates.
However, a synchronization with the orderbook processing would be required to achieve a real graceful shutdown.
The orderbook is idempotent so replaying a stream in case it got interrupted should not result in corrupted data.


#### Usage

```
go test ./...

# hit CTRL-C to shutdown the parser
go run main.go
```

#### Tests

Due to time constraints I only added very basic testing for the orderbook based on the example provided in the task definition.
