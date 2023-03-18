# Typechan 

A TUI typing test powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Quotes/sentences for typing are retrieved from [quotable](https://github.com/lukePeavey/quotable)

```go
# Fetch dependencies
go mod download 

# Compile to binary
go build .

# Launch the test!
./typechan sprint
```
---

You can launch the test in 2 different modes:

## Sprint mode 🏃🏻‍♀️
Complete the sentence as fast as you can.
```go
./typechan sprint
```

## Timed mode ⏱️
Type as far as you can within the time limit.
```go
# Starts a 5-minute test
./typechan timed

# Specify a different time limit e.g. 30 seconds
./typechan timed -s 30
```
