# Day 21 - Keypad Conundrum

## Go learnings
- Writer and Reader
### Modules, Packages and repos
Ex.

- I have a a folder: `~/Code/go/aoc2024` with a git repo
    - I have a remote repo: `github.com/mpja69/aoc2024`
    - (but since there is no "general/common" code, i have _not_ run `go mod init`, so the repo it _not_ a module)
- Every day I create a new sub-folder: `dayXX`
    - and I do a `go mod init aoc2024/dayXX`, (but I should probably have done `go mod init github.com/mpja69/a0c2024/dayXX`)
    - but i do _not_ run `git init`! Every day is still part och the same `aoc2024` git repo.
- Sometime I create a subfolder, e.g. `./day21/keypad`
    - where I did a `go mod init github.com/mpja69/aoc2024/day21/keypad`, and got a `go.mod` file
    - _but_ I removed that one! Because that package/folder is part of the `day21`-module
    - when I develop in the working tree, I do a `go mod edit -replace github.com/mpja69/aoc2024/day21/keypad=./keypad/`


