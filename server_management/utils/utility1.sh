#!/bin/bash

# Define utility functions here

function count_lines() {
	echo "$(wc -l $1 | awk '{print $1}')"
}
