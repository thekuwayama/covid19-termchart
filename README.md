# covid19-termchart

[![CI](https://github.com/thekuwayama/covid19-termchart/workflows/CI/badge.svg)](https://github.com/thekuwayama/covid19-termchart/actions?workflow=CI)
[![Go Report Card](https://goreportcard.com/badge/github.com/thekuwayama/covid19-termchart)](https://goreportcard.com/report/github.com/thekuwayama/covid19-termchart)
[![MIT licensed](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://raw.githubusercontent.com/thekuwayama/covid19-termchart/master/LICENSE.txt)

`covid19-termchart` is the CLI that print "日本国内の感染者数（NHKまとめ）"  to iTerm2.

- https://www3.nhk.or.jp/news/special/coronavirus/data-all/


## Install

```bash
$ go get -u -v github.com/thekuwayama/covid19-termchart
```


## Usage

```bash
$ covid19-termchart -help
Usage of covid19-termchart:
  -day int
    	period to aggregate (default 365)
```

![image](https://user-images.githubusercontent.com/42881635/104087803-c39f2380-52a5-11eb-9853-b45c6f276c9b.png)


## License

The CLI is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
