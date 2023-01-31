#!/bin/bash

# This script pulls from project gutenberg.
# Sorry, it's real basic because it works and fancier is for the future.

# FUTURE: Move into golang and trigger and track from there
# FUTURE: Allow non-managers to sync over IPFS
SAVEDIR=/storage-enc/small/gutenberg
rsync --info=progress2 -ha --delete aleph.gutenberg.org::gutenberg/ $SAVEDIR
