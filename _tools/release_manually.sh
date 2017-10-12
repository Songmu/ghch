#!/bin/sh
set -e

ver=v$(gobump show -r)
make crossbuild
ghr $ver dist/$ver
