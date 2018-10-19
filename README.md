# alfred-chrome-history

Search Chrome history from Alfred and open in browser.

## Features

- Search Chrome history (title and URL)
- Support another Chrome profile

## Installation

Clone and `make dist` or just download [binary releases](https://github.com/pasela/alfred-chrome-history/releases).

```sh
git clone https://github.com/pasela/alfred-chrome-history.git
cd alfred-chrome-history
make dist
open alfred-chrome-history.alfredworkflow
```

## Usage

in Alfred:

```
ch {query}
```

## Use another Chrome profile

1. Open workflow `Chrome History` in Alfred Workflows tab.
2. Open Workflow Configuration dialog by upper right side button.
3. Set `CHROME_PROFILE` variable with your Chrome profile directory name or path such as `Profile 1`.

## License

MIT

## Author

Yuki (a.k.a. pasela)
