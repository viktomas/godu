# godu
[![Build Status](https://travis-ci.org/viktomas/godu.svg?branch=master)](https://travis-ci.org/viktomas/godu)
[![codecov](https://codecov.io/gh/viktomas/godu/branch/master/graph/badge.svg)](https://codecov.io/gh/viktomas/godu)
[![Go Report Card](https://goreportcard.com/badge/github.com/viktomas/godu)](https://goreportcard.com/report/github.com/viktomas/godu)

Find the files that are taking up your space.

Tired of looking like a noob with [Disk Inventory X](http://www.derlien.com/) or SpaceMonger? Do you want something that
* can do the job
* works in terminal
* makes you look cool
* is written in Golang
* you can contribute to

??

Well then **look no more** and try out the godo.

## Installation
```
go get -u github.com/tomasvik/godu
```

## Configuration
You can specify names of ignored folders in `.goduignore` in your home directory:
```
> cat ~/.goduignore
node_modules
>
```
I found that I could reduce time it took to crawl throug the whole drive to 25% when I started ignoring all `node_modules` which cumulativelly contain gigabytes of small text files.

The `.goduignore` is currently only supporting whole folder names. PR that will make it work like `.gitignore` is welcomed.

## Usage
```
godo ~
godo folder-that-is-too-big
```

Once the folder is crawled (can take up to few minutes), you move around by selecting numbers (moving deeper in the structure) or you just press enter (moving up in the structure)

Exit with CTRL+D or kill of your own choice.
