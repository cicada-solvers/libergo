#!/bin/bash

./processtextone.sh
./processtextpull.sh
./processtext.sh

mv *.filtered.txt /data/3301/chaffetext/filtered
mv *.txt /data/3301/chaffetext/full