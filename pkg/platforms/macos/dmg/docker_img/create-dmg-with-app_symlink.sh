#!/bin/bash

SRC_DIR=$2

# Create symlink to /Applications if APPLICATIONS_SYMLINK is set
if [ -n "$APPLICATIONS_SYMLINK" ]; then
  SRC_DIR="/tmp/src_copy"
  # copy source dir to temp dir in order to create a copy with current user permissions
  cp -r "$2" "$SRC_DIR"
  ln -s /Applications "$SRC_DIR/Applications" || exit 1
fi

# From the original create-dmg.sh
TMP_DMG="/tmp/tmp.dmg"
genisoimage -quiet -V "$1" -D -R -apple -no-pad -o "$TMP_DMG" "$SRC_DIR"
dmg dmg "$TMP_DMG" "$3"