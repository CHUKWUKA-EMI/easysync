#!/bin/bash
source ./.env

atlas migrate apply --url $DATABASE_URL