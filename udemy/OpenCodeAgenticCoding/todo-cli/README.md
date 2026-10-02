# todo-cli

A todo list for the terminal. Pure Go standard library — no dependencies — with
subtasks that nest as deep as you like, a tree view, and colorized output when
stdout is a terminal.

## Install

Requires Go 1.26 or newer.

```bash
go build -o todo .
```

Then run it as `./todo`, or copy the binary somewhere on your `PATH` and run
`todo`. Every example below uses `todo`.

## Quick start

```bash
todo add "ship v1.0"
todo add -p 1 "write tests"
todo add -p 2 "unit tests"
todo add -p 2 "integration tests"
todo add -p 1 "cut release"
todo add "buy milk"
todo done 3
todo list -all
```

```
ID  TITLE                       STATUS   CREATED
 1  ship v1.0                   pending  2026-10-02
 2  │  ├─ write tests           pending  2026-10-02
 3  │  │  ├─ unit tests         done     2026-10-02
 4  │  │  └─ integration tests  pending  2026-10-02
 5  │  └─ cut release           pending  2026-10-02
 6  buy milk                    pending  2026-10-02
```

On a terminal the same table is colorized: bold header, cyan IDs, yellow
`pending`, green `done`, dim dates, and completed titles struck through. The
tree connectors are dim. Piped or redirected output is always plain text.

## Commands

| Command | Description |
| --- | --- |
| `todo add <title>` | Add a todo. `-p, -parent <id>` nests it under another todo. |
| `todo list`, `ls` | Show pending todos. `-all` includes completed ones. |
| `todo done <id>...` | Mark todos done. Ids may be comma-separated: `todo done 1,2`. |
| `todo reopen <id>...` | Mark completed todos pending again. |
| `todo rm <id>...`, `remove` | Delete todos **and their subtasks**. |
| `todo clear` | Delete completed todos, promoting unfinished subtasks up a level. `-all` deletes everything. |
| `todo help`, `h` | Show usage. |

`-file <path>` is accepted by every command and overrides the data file;
`-color` is specific to `todo list`.

## Subtasks

Any todo can be the parent of another, to any depth:

```bash
todo add -p 2 "integration tests"   # nest under #2
todo add -p 3 "fixtures"            # nest under #3
```

Deleting a parent takes its whole subtree with it:

```bash
todo rm 1
removed 5 todos    # #1 and the four subtasks under it
```

`clear` is the exception: it deletes only completed todos and promotes their
subtasks so unfinished work is never lost.

```bash
todo list -all
ID  TITLE        STATUS   CREATED
 1  alpha        done     2026-10-02
 2     └─ bravo  pending  2026-10-02

todo clear
cleared 1 todo

todo list -all
ID  TITLE  STATUS   CREATED
 2  bravo  pending  2026-10-02
```

Filtering never hides a subtree: a completed parent stays visible while it still
has pending children, so the hierarchy is always readable. Completing a parent
does not complete its subtasks.

## Where the data lives

A single JSON file, chosen in this order:

1. `-file <path>`
2. `$TODO_FILE`
3. `<user config dir>/todo-cli/todos.json` — `~/Library/Application Support` on
   macOS, `~/.config` on Linux, `%AppData%` on Windows

```json
{
  "next_id": 3,
  "todos": [
    {
      "id": 1,
      "title": "alpha",
      "done": false,
      "created_at": "2026-10-02T11:14:17.04893-04:00"
    },
    {
      "id": 2,
      "title": "bravo",
      "done": true,
      "parent_id": 1,
      "created_at": "2026-10-02T11:14:17.052029-04:00",
      "completed_at": "2026-10-02T11:14:17.054837-04:00"
    }
  ]
}
```

The file is written atomically (temp file plus rename), so an interrupted write
cannot corrupt it. Ids are never reused. On load the tree is validated: a
hand-edited file with a missing parent or a parent cycle is rejected with an
error rather than silently mangled. `parent_id` is omitted for top-level todos.

## Colors

`todo list` accepts `-color` with `auto` (default), `always`, or `never`. In
`auto` mode output is colorized only when stdout is a terminal, and `NO_COLOR`
or `TERM=dumb` turn it off. Force colors through a pipe with:

```bash
todo list -all -color=always | less -R
```

## Exit codes and errors

Errors go to stderr with a `todo:` prefix and exit status 1. Nothing is written
unless the whole operation succeeds: `done 1,99` leaves todo 1 untouched.

## Development

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
staticcheck ./...   # optional
```

Layout:

```
main.go              subcommand dispatch, flags, command implementations
table.go             aligned, ANSI-aware column renderer and tree connectors
internal/todo/       Todo and List types, tree operations, JSON store
internal/color/      ANSI palette, terminal detection, color modes
```