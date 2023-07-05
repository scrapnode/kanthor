#!/bin/sh
set -e

# split args string to args param and pass it to our application command
# /app/kanthor "service start dataplane" -> NOT WORK because the app receives string as 1 args
# /app/kanthor service start dataplane -> WORK because the app receives list of args
set -- /app/kanthor $@

# then execute the command with args
exec $@