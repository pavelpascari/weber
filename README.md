<img src="weber.svg" width="48">

# Weber - A Go Command Line Tool for Web Developers

Weber is a command line interface (CLI) tool written in Go, designed to make web development with a Chrome-based browser hassle free by providing 
easy-to-use utilities that ensure you offer the best web experience to your users.

## Use Cases
1. **Cache-Control Header Analysis** -
Weber can be used to analyze the `Cache-Control` header of a website. You can get a report of the resources your clients are loading to analyse at a glance.
This is useful for web developers who want to ensure that their website is optimized for caching.
2. 


## Installation

To use Weber, you must first install it on your local machine. You can do this by running the following command in your terminal:

```bash
go get -u github.com/pavelpascari/weber
```

`weber` relies on Chrome browser to perform its tasks. You can install it by following the instructions [here](https://www.google.com/chrome/).

`weber` relies on `github.com/chromedp/chromedp` to leverage Chrome DevTools Protocol.

## Usage

Once you have installed Weber, you can use it to perform various tasks such as:

```bash
$ weber                                                               
Usage: weber -o output.csv [OPTIONS] <url>

Options:
  -X <method>   Comma-separated list of HTTP methods to watch for (GET, POST, OPTIONS, PUT, DELETE). Default behavior is to consider all methods.
  -H <string>   Comma-separated list of hostname or IP address to watch for. Default behavior is to consider all hosts.
  -o <file>     Write the response to a file. CSV is the default and only supported format.
  -c <string>   Comma-separated list of columns to write to the output file. Default is "url,method,status". Available columns are:
                    status, url, method, Content-Type, Cache-Control, Content-Length
  -v            Enable verbose logging to observe all browser events.
  -q            Disable all logging.

Examples:
      # Watch for all requests to https://example.com
      weber -o output.csv https://example.com

      # Watch for all requests to https://example.com> and https://example.org
      weber -H example.com,example.org -o output.csv https://example.com

      # Watch for GET requests on example.org and output the URL, request method, the status code, and cache-control header
      weber -X GET -H example.org -o output.csv -c "url,method,status,Cache-Control" https://example.com
```

## Contributing

Weber is an open source project and we welcome contributions from the community. If you would like to contribute, please open an issue or submit a pull request.
