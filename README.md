# action-averager

This project contains a package that exports an ActionAverager interface and
an ActionAverage struct that implements that interface. They have two exported
functions AddAction and GetStats.

AddAction takes in a json formatted string like:
`{"action":"jump","time":456}` and returns an error. GetStats returns a json
formatted string which contains the running average times for each of the
added actions like:
`[{"action":"crawl","avg":300},{"action":"jump","avg":289.5}]`

## Set up

This project expects that you are following the go development guide:
https://golang.org/doc/code.html. In order to set up this project use that
guide to set up your go environment then run the following commands:

* `cd $GOPATH/src/github.com`
* `git clone https://github.com/rtalleyman/action-averager.git`
* `cd action-averager`
* `make init`

## Building

Run `make build`

This will compile an executable that acts as an example for using the package.

## Running

Run `make run`

This will run the example executable for this project.

## Testing

Testing is done using the ginkgo BDD test framework with gomega as the matcher
library. More information can be found here: https://onsi.github.io/ginkgo/.
The repo needs to be initialized before running these tests so be sure to run
`make init` and follow any and all of the instructions in the Set up section.

Run `make test` to quickly run all of the tests.

Run `make test-debug` to debug tests and to provide more verbose test output.

## Other Make targets

Running `make all` will build, run, then delete the example executable.

Running `make clean` will delete the example executable.
