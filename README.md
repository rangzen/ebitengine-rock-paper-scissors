# Rock Paper Scissors War

This is a simple rock paper scissors game written in Go using the [Ebitengine](https://ebitengine.org/) game library.

I wanted my own version of https://9gag.com/gag/ap92GVn.

You can see an online version at https://rangzen.github.io/ebitengine-rock-paper-scissors/.

## Running

To run the game, you need to have Go installed.

Then, run the following command if you have cloned the repository:

```bash
go run main.go
```

or simply

```bash
go run github.com/rangzen/ebitengine-rock-paper-scissors@latest
```

## Building

To build the game, run the following command:

```bash
go build
```

## Building the WebAssembly version

To build the WebAssembly version, run the following command:

```bash
GOOS=js GOARCH=wasm go build -o ebitengine-rock-paper-scissors.wasm github.com/rangzen/ebitengine-rock-paper-scissors
```

## Resources

* <a target="_blank" href="https://icons8.com/icon/9FSQ5judlnAN/rock">Rock</a> icon by <a target="_blank" href="https://icons8.com">Icons8</a>
* <a target="_blank" href="https://icons8.com/icon/jDDj4ExfgPZV/page-facing-up">Page Facing Up</a> icon by <a target="_blank" href="https://icons8.com">Icons8</a>
* <a target="_blank" href="https://icons8.com/icon/A7egVNynrr0h/scissors">Scissors</a> icon by <a target="_blank" href="https://icons8.com">Icons8</a>
