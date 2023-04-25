<img src="weber.svg" width="48">

# Weber - A Go Command Line Tool for Web Developers

Weber is a command line interface (CLI) tool written in Go, designed to make web development with a Chrome-based browser hassle free by providing 
easy-to-use utilities that ensure you offer the best web experience to your users.

## Use Cases
1. **Cache-Control Header Analysis** -
Weber can be used to analyze the `Cache-Control` header of a website. You can get a report of the resources your clients are loading to analyse at a glance.
This is useful for web developers who want to ensure that their website is optimized for caching.
2. **GDPR checks** -
Weber can output easily the hosts that your website is loading resources from. This is useful for web developers who want to ensure that their website is GDPR compliant.
3. etc.


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
$ weber -h
Usage: weber [OPTIONS] <url>

Options:
  -X <method>   Comma-separated list of HTTP methods to watch for (GET, POST, OPTIONS, PUT, DELETE). Default behavior is to consider all methods.
  -H <string>   Comma-separated list of hostname or IP address to watch for. Default behavior is to consider all hosts.
  -o <file>     Write the response to a file. CSV is the default and only supported format.
  -c <string>   Comma-separated list of columns to write to the output file. Default is "url,method,status". 
                Available columns are any valid response header, plus: url, method, status, hostname.
  -t <int>      Time after which the program will give up waiting for another request. 
                Every new processed request will restart the timer. Default is 15 seconds.
  -v            Enable verbose logging to observe all browser events.
  -q            Disable all logging.
  -h            Show this help message.

Examples:
      # Watch for all requests to https://example.com
      weber -o output.csv https://example.com

      # Watch for all requests to https://example.com> and https://example.org
      weber -H example.com,example.org -o output.csv https://example.com

      # Watch for GET requests on example.org and output the URL, request method, the status code, and cache-control header
      weber -X GET -H example.org -o output.csv -c "url,method,status,Cache-Control" https://example.com

      # Get all hostnames accessed for/by resources of a website
      weber -o output.csv -c "hostname" https://example.com
```

## Examples

### Cache-Control Header Analysis

```bash
$ weber -o google.csv -c Content-Type,Cache-Control https://google.com
.................
Giving up waiting...
Flushing writer...
Done.
```

```bash
$ cat google.csv
cat google.csv 
Content-Type,Cache-Control
text/html; charset=UTF-8,"private, max-age=0"
...
text/javascript; charset=UTF-8,"public, max-age=31536000"
```

### GDPR checks
To get a list of all domains that your website is loading resources from, you can run the following command:

```bash
$ weber -o focus.csv -c hostname https://www.focus.de/kultur/kino_tv/precht-beschimpft-baerbock-in-show-mit-lanz-unsagbar-zum-fremdschaemen_id_192056533.html
```

```bash
$ cat focus.csv | sort | uniq

5baf1288cf.dl8.me
a.bf-ad.net
a.bf-tools.net
acdn.adnxs-simple.com
...
web-vitals.bfops.io
widgets.opinary.com
www.focus.de
```

## Contributing

Weber is an open source project and we welcome contributions from the community. If you would like to contribute, please open an issue or submit a pull request.
