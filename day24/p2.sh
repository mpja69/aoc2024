#! /bin/bash

# z?? should be output of XOR except MSB
# get a list where: output is  "> z", but no "XOR"-ops, and not the last bit...get JUST the output-names
CANDIDATE_1=$(cat input.txt | grep '> z' | grep -v 'XOR' | grep -v 'z45$' | awk '{print $5}')

# all the XOR should have either x?? y?? for input, or z?? for output.
# get a list of all "XOR"-ops, but not the lines starting with "x" or "y", and not output "> z"...get JUST the output-names
CANDIDATE_2=$(cat input.txt | grep ' XOR ' | grep -v '^x' | grep -v '^y' | grep -v '> z' | awk '{print $5}')

# input of OR should be always output of AND except for LSB
# get a list of "OR"-ops, but get JUST the 2 inputs, and sort them and remove duplilcates
INPUT_OF_OR=$(cat input.txt | grep ' OR ' | awk '{ print $1; print $3 }' | sort -u)

# Exclude the first "input"-"AND", but include all other "AND"-ops, ...get JUST the output-names (sort and remove duplicates)
OUTPUT_OF_AND=$(cat input.txt | grep -v 'x00 AND y00' | grep ' AND ' | awk '{ print $5 }' | sort -u)

# From all "OR"s and "AND"s, (replace WHITE with NEWLINE), get just the unique ones
CANDIDATE_3=$(comm -3 <(echo $INPUT_OF_OR | tr ' ' '\n') <(echo $OUTPUT_OF_AND | tr ' ' '\n'))

# From all candidates: (replace WHITE with NL), sort and remove duplicates, then replace NL with ",", and remove last ","
echo $CANDIDATE_1 $CANDIDATE_2 $CANDIDATE_3 | tr ' ' '\n' | sort -u | tr '\n' ',' | sed -e 's/,$//'
