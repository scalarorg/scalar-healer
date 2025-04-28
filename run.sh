#!/bin/bash
make compose
#ssh btc4 -L 60001:127.0.0.1:60001 &
make daemon
