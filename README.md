# go-starter

`go-starter` is a public golang starter repo template. 
It is useful for getting started with a new project. 


## How to use

```bash
gh repo create <name> \
  --public \
  --template https://github.com/pavelpascari/go-starter
```

## Features

- [ ] GitHub Actions workflow for CI
  - [x] run tests
    - [ ] ...with coverage check
  - [x] vulnerability scanning with SARIF report thanks to [@Templum](https://github.com/Templum)'s `Templum/govulncheck-action`
  - [x] static check thanks to [@dominikh](https://github.com/dominikh)'s `dominikh/staticcheck-action`
- [x] Dependabot security scanning for `go`, `github-actions`

