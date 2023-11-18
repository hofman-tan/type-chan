# Typechan

A TUI typing test powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Random quotes and passages are retrieved from [quotable](https://github.com/lukePeavey/quotable).

```shell
# Fetch dependencies
go mod download

# Compile to binary
go build .

# Launch the test!
./typechan sprint
```

---

typechan comes in 2 different modes:

## Sprint mode 🏃🏻‍♀️

Complete the test as fast as you can.

```shell
./typechan sprint
```

## Timed mode ⏱️

Type as much as you can before timeout.

```shell
# Starts a 5-minute test
./typechan timed

# Change the time limit to e.g. 30 seconds
./typechan timed -s 30s
```
