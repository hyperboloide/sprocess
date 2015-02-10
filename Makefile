## 
## Makefile
## 
## Created by Frederic DELBOS - fred@hyperboloide.com on Oct 28 2014.
## This file is subject to the terms and conditions defined in
## file 'LICENSE', which is part of this source code package.
## 

VERSION = $(shell cat data/version)
NAME = github.com/hyperboloide/sprocess

test:
	go test

fmt:
	go fmt ./...
