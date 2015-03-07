# go-playbooks
Various hopefully useful playbooks for solving common problems with go.

Each playbook starts form a GO library and explains how to use it to solve a practival problem.


## Sitemap

```
.
├── encoding
│   └── json
│       ├── basetypes-marshalling
│       └── custom-marshaller
└── net
    └── http
        ├── handling-service-termination
        ├── redirecting-to-https
        │   ├── aws-ebs-nginx
        │   └── handle-redirect
        └── request-caching
```

## encoding/json

* <b>Basetype marshalling</b>: how to marshal basetypes and how to handle missing fields.
* <b>Custom marshalling</b>: how to implement custom marshalling and unmarshalling logic.

## net/http

* <b>Handling service termination</b>: how to gracefully terminate http handling.
* <b>Serving https with AWS EBS</b>: how to serve https requests with AWS EBS.
* <b>Handling https redirect</b>: how to programmatically redirecting from http to https.
* <b>Http request caching</b>: how to cache http requests.

# About

GO is a great language. I am writing these playbooks to help other developers to practice it.

## How to contribute

Pull requests are welcomed. Also if you want to request a Playbook just open a ticket.
