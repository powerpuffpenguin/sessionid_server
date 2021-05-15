#!/bin/bash

set -eu;


exec gosu sessionid "$@"
