##
## Makefile
##
## Created by Frederic DELBOS - fred@hyperboloide.com on Oct 28 2014.
## This file is subject to the terms and conditions defined in
## file 'LICENSE', which is part of this source code package.
##

VERSION = $(shell cat .version)
NAME = github.com/hyperboloide/sprocess

test:
	ginkgo -r

fmt:
	go fmt ./...

travis:
	ginkgo \
	--skip="S3bucket|GoogleCloud" \
	-r \
	--race \
	--randomizeSuites \
	--trace

.PHONY: test fmt travis
