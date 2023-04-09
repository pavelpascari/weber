# Weber - A Go Command Line Tool for Web Developers

Weber is a command line interface (CLI) tool written in Go, designed to make web development hassle free by providing easy-to-use utilities that ensure you offer the best web experience to your users.

## Installation

To use Weber, you must first install it on your local machine. You can do this by running the following command in your terminal:

```bash
go get -u github.com/pavelpascari/weber
```

## Usage

Once you have installed Weber, you can use it to perform various tasks such as:

### 1. Checking resource caching

Weber provides a command that enables you to check for resources that are not cached and were downloaded from a particular domain. To do this, run the following command:

```bash
weber check cache --domain example.com session.har
```

