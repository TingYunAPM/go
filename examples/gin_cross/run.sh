#!/bin/bash
./gin_cross tingyun.json :8080 http://www.tingyun.com/ &
./gin_cross t1.json :8081 http://127.0.0.1:8080/extern &

