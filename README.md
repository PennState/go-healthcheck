# go-healthcheck

Implements the (draft) RFC for retrieving health status from Go servers.

## Protocol

This library is both wire and protocol compatible with the (draft) RFC
for HTTP API Health Checks.  This specification can be found at:

<https://inadarei.github.io/rfc-healthcheck>

## Project layout

The following packages are available:

- health - This package contains the data models that implement the
  RFC's protocol along with an ``http.Handler`` factory method and
  a client.

- checks - This package contains a set of health checks that are
  can be used in many environments.  Custom checks should be written
  to implement the ``Checker`` interface.
