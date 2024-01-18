#!/bin/bash

# Set the path variable
HOST_DIR="$(pwd)"
PARENT_DIR=$(dirname "$(pwd)")
VOLUME_NAME="gis_map_info"
VOLUME_NAME_DB="gis_map_info_db"
HOST_DIR_DB="$PARENT_DIR/$VOLUME_NAME_DB"
NETWORK_NAME="gis_map_info_network"

# Check if the volume already exists
if docker volume inspect "$VOLUME_NAME" >/dev/null 2>&1; then
  echo "Volume '$VOLUME_NAME' already exists. Skipping creation."
else
  # Create a Docker volume
  docker volume create --driver local \
    --opt type=none \
    --opt device="$HOST_DIR" \
    --opt o=bind \
    "$VOLUME_NAME"

  echo "Volume '$VOLUME_NAME' created."
fi

if docker volume inspect "$VOLUME_NAME_DB" >/dev/null 2>&1; then
    echo "Volume '$VOLUME_NAME_DB' already exists, Skipping creation."
else
  mkdir "$HOST_DIR_DB" || true
    # Create a Docker volume
  docker volume create --driver local \
    --opt type=none \
    --opt device="$HOST_DIR_DB" \
    --opt o=bind \
    "$VOLUME_NAME_DB"

  echo "Volume '$VOLUME_NAME_DB' created."
fi

VOLUME_NAME="gis_map_info_node"
HOST_DIR="$HOST_DIR/sub_app/node"
# Check if the volume already exists
if docker volume inspect "$VOLUME_NAME" >/dev/null 2>&1; then
  echo "Volume '$VOLUME_NAME' already exists. Skipping creation."
else
  # Create a Docker volume
  docker volume create --driver local \
    --opt type=none \
    --opt device="$HOST_DIR" \
    --opt o=bind \
    "$VOLUME_NAME"

  echo "Volume '$VOLUME_NAME' created."
fi

VOLUME_NAME="martin_config"
HOST_DIR="$(pwd)"
HOST_DIR="$HOST_DIR/sub_app/martin"
# Check if the volume already exists
if docker volume inspect "$VOLUME_NAME" >/dev/null 2>&1; then
  echo "Volume '$VOLUME_NAME' already exists. Skipping creation."
else
  # Create a Docker volume
  docker volume create --driver local \
    --opt type=none \
    --opt device="$HOST_DIR" \
    --opt o=bind \
    "$VOLUME_NAME"

  echo "Volume '$VOLUME_NAME' created."
fi

# Check if the network already exists
if docker network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
  echo "Network '$NETWORK_NAME' already exists. Skipping creation."
else
  # Create a Docker network
  docker network create "$NETWORK_NAME"

  echo "Network '$NETWORK_NAME' created."
fi

echo "Setup complete."
