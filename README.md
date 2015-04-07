# go-playbooks
Go-playbooks is a collection of various hopefully useful playbooks and examples for solving common problems with go.

Go is a new and amazin language with a very short learning curve. However how any new language there is always a bit of confusion on how to address common tasks and how to best use use the standard libraries. Go-playbooks is my attempt to provide a reference to help developers adopt GO at its best.

[![Build Status](https://travis-ci.org/msama/go-playbooks.svg?branch=master)](https://travis-ci.org/msama/go-playbooks)

## Playbooks

```
.
├── algorithms
│   ├── algebraic
│   │   └── pascals-triangle
│   ├── constraints
│   │   └── n-queens
│   ├── data-streaming
│   │   └── moving-average
│   └── sorting
│       ├── bubble-sort
│       └── merge-sort
├── database
│   └── sql
│       └── driver
│           └── custom-types
├── encoding
│   └── json
│       ├── basetypes-marshalling
│       └── custom-marshaller
└── net
    └── http
        ├── graceful-service-termination
        ├── https-aws-ebs-nginx
        ├── programmatic-http-redirect
        └── request-caching
```

The above tree shows all the available playbooks. Playbooks are organised with a tree tructure following the `<package>/<book>` convention. Each playbook contains a `main.go` file which can be used to run it. Alse each playbook contains a `README.md` file with a description of what the playbook does and additional setup instruction.

To run a playbook from the root folder you can run the following:
```
cd <package>/<book>
go run main.go
```
Certain playbook may require additional setup. Please read the `README.md` first.

## algorithms/algebraic

* <b>[Pascal's triangle](https://github.com/msama/go-playbooks/tree/master/algorithms/algebraic/pascals-triangle)</b>: the Pascal's triangle both iterative and recursive.

## algorithms/constraints

* <b>[N-queens](https://github.com/msama/go-playbooks/tree/master/algorithms/constraints/n-queens)</b>: an iterative solution of the n-queens problem using constraints.

## algorithms/data-streaming

* <b>[Moving average](https://github.com/msama/go-playbooks/tree/master/algorithms/data-streaming/moving-average)</b>: a GO implementation of the moving average algorithm.

## algorithms/sorting
* <b>[Bubble sort](https://github.com/msama/go-playbooks/tree/master/algorithms/sorting/bubble-sort)</b>: implementation of bubble sort algorithm.
* <b>[Merge sort](https://github.com/msama/go-playbooks/tree/master/algorithms/sorting/merge-sort)</b>: recursive implementation of merge sort.

## database/sql/driver

* <b>[Custom types](https://github.com/msama/go-playbooks/tree/master/database/sql/driver/custom-types)</b>: how to write to specialized types to a sql driver.
* <b>[Struct serialization](https://github.com/msama/go-playbooks/tree/master/database/sql/driver/struct-serialization)</b>: how to serializa go structs in sql fields.

## encoding/json

* <b>[Basetype marshalling](https://github.com/msama/go-playbooks/tree/master/encoding/json/basetypes-marshalling)</b>: how to marshal basetypes and how to handle missing fields.
* <b>[Custom marshalling](https://github.com/msama/go-playbooks/tree/master/encoding/json/custom-marshaller)</b>: how to implement custom marshalling and unmarshalling logic.

## net/http

* <b>[Graceful service termination](https://github.com/msama/go-playbooks/tree/master/net/http/graceful-service-termination)</b>: how to gracefully terminate http handling.
* <b>[Serving https with AWS EBS](https://github.com/msama/go-playbooks/tree/master/net/http/https-aws-ebs-nginx)</b>: how to serve https requests with AWS EBS.
* <b>[Programmatic http redirect](https://github.com/msama/go-playbooks/tree/master/net/http/programmatic-http-redirect)</b>: how to programmatically redirecting from http to https.
* <b>[Http request caching](https://github.com/msama/go-playbooks/tree/master/net/http/request-caching)</b>: how to cache http requests.


## Contributing

### Submitting a new playbook

The best way to contribute is to follow the Github pull request workflow. Fork go-palybooks, create a functional branch and send a pull request to `msama/master`. 

Please respect the followign submission guidelines:
* <b>Should be useful</b>: Playbooks need to address a single well defined problem.
* <b>Should be as simple as possible but not simpler</b>: Playbooks have to be easy to read and self contained.
* <b>Must compile</b>: the pull request must compile. Please add new playbooks to Travis.
* <b>Must be well documented</b>: please update or create README.md files for each Playbook.

### Requesting a new playbook

If you are trying to solve a problem and you would like to see a playbook, please open a ticket in Github.

### Reporting bugs or improving a solution

If you see a bug or if you want to suggest a different solution please either send a pull request or open a ticket in Github.
