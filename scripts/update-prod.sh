#!/bin/sh
export TT=$SMS_TOKEN

ssh $USER@$IP "export TOKEN=$TT && sudo -E docker compose down && sudo -E docker compose up"