#! /usr/bin/python3
import sys
import itertools
import random

# Define the number of bytes for the initial combination
num_bytes = 8

# Loop through all 8-byte combinations (256^8 possibilities)

for combo in itertools.product(range(256), repeat=num_bytes):
    # Create the bytearray for the current 8-byte combination
    base_bytes = bytearray(combo)
    
    # Generate a random number of extra bytes to append (between 1 and 8)
    extra_size = random.randint(1, 10000)
    
    # Generate random extra bytes to append
    extra_bytes = bytearray(random.getrandbits(8) for _ in range(extra_size))
    
    # Combine the base bytes with the random extra bytes
    combined_bytes = base_bytes + extra_bytes
    
    # Write the combined bytes to stdout
    sys.stdout.buffer.write(combined_bytes)

    # Optionally add a newline at the end
    sys.stdout.buffer.write(b"\n")
