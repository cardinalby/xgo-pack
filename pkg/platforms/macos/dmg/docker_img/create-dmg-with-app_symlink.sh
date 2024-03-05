#!/bin/bash

# Create symlink to /Applications if APPLICATIONS_SYMLINK is set
if [ -n "$APPLICATIONS_SYMLINK" ]; then
  ln -s /Applications "$2/Applications" || exit 1
fi

# Call the original create-dmg.sh
/bin/bash create-dmg.sh "$1" "$2" "$3"