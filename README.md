# godu
[![Build Status](https://travis-ci.org/viktomas/godu.svg?branch=master)](https://travis-ci.org/viktomas/godu)
[![codecov](https://codecov.io/gh/viktomas/godu/branch/master/graph/badge.svg)](https://codecov.io/gh/viktomas/godu)
[![Go Report Card](https://goreportcard.com/badge/github.com/viktomas/godu)](https://goreportcard.com/report/github.com/viktomas/godu)

Find the files that are taking up your space.

<img src="https://media.giphy.com/media/OGAm19SP0enD2/giphy.gif" width="100%" />

Tired of looking like a noob with [Disk Inventory X](http://www.derlien.com/) or SpaceMonger? Do you want something that
* can do the job
* works in terminal
* makes you look cool
* is written in Golang
* you can contribute to

??

Well then **look no more** and try out the godu.

## Installation
```
go get -u github.com/viktomas/godu
```

## Configuration
You can specify names of ignored folders in `.goduignore` in your home directory:
```
> cat ~/.goduignore
node_modules
>
```
I found that I could reduce time it took to crawl through the whole drive to 25% when I started ignoring all `node_modules` which cumulatively contain gigabytes of small text files.

The `.goduignore` is currently only supporting whole folder names. PR that will make it work like `.gitignore` is welcomed.

## Usage
```
godu ~
godu -l 100 / # walks the whole root but shows only files larger than 100MB
```

Mind you `-l  <size_limit_in_mb>` option is not speeding up the walking process, it just allows you to filter small files you are not interested in from the output. **The default limit is 10MB**.

Use arrows to move around, ESC or CTRL+C to quit
