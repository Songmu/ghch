#!/bin/sh
set -e

ver=v$(gobump show -r)
make crossbuild
ghr v$ver dist/v$ver
