#!/bin/bash

# Directory where the ekycapp-randomno files are stored
WATCH_DIR="/home/path_toapp"

# Symlink name
SYMLINK="$WATCH_DIR/ekycapp.symlink"
TEMP_SYMLINK="$WATCH_DIR/ekycapp.symlink.tmp"

# Watch for specific file pattern in the directory
inotifywait -m -e delete --format "%f" --exclude 'storage/.*' $WATCH_DIR | while read FILENAME
do
    if [[ "$FILENAME" == ekycapp-* ]]; then

        echo "======= New deployment: $(date '+%Y-%m-%d %H:%M:%S') ======="
        echo "New Deployment Initated :)"
        echo "Removed old build $FILENAME"

        # Find the most recent valid ekycapp-* file to update the symlink
        NEW_TARGET=$(ls -1tr $WATCH_DIR/ekycapp-* 2>/dev/null | tail -n 1)

        if [[ -z "$NEW_TARGET" ]]; then
            echo "No valid ekycapp-* files left in the directory."
            sudo supervisorctl stop apx_fiber_ekycapp
            echo "Stopped apx_fiber_ekycapp due to lack of valid target."
        else
            echo "Settingup new Build..."
            sudo chmod +x "$NEW_TARGET"
            ln -sf "$NEW_TARGET" "$SYMLINK"
            echo "Updated symlink to $NEW_TARGET"
            echo "Restarting supervisor..."
            # Start the process using Supervisor
            if ! sudo supervisorctl restart apx_fiber_ekycapp; then
                echo "Failed to start apx_fiber_ekycapp" >&2
            else
                echo "Restarted apx_fiber_ekycapp successfully after updating symlink."
                echo "New deploment completed!"
            fi
        fi
        echo "=========== Deployment end ==========" 
    fi
done