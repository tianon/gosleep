# `gosleep`

`sleep(1)`, but with Go duration parsing (`1m`, `2h3m2s`) OR explicit time to sleep until (`13:04`, `2016-11-18 00:00:00`) and a progress bar.

# Example

![example animated GIF](example.gif)

```console
$ gosleep --help
Usage:
  gosleep --for <duration> OR --until <time>

Application Options:
      --for
      --until

Help Options:
  -h, --help   Show this help message

$ gosleep --for 10s
 10s / 10s [===========================================================] 100.00%
$ date
Thu Nov 17 13:21:56 PST 2016
$ gosleep --until 13:22
error: requested sleep time in the past: -1s
$ gosleep --until 13:23:05
 58s / 58s [===========================================================] 100.00%
$ gosleep --until '2016-11-17 00:00:00'
error: requested sleep time in the past: -13h23m43s
$ gosleep --until '2016-11-18 00:00:00'
 1m5s / 10h36m14s [>---------------------------------------------------]   0.17%
```
