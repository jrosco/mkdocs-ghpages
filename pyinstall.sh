#!/bin/bash

# Define the virtual environment directory
VENV_DIR=".venv"

# Check if the virtual environment exists
if [ ! -d "$VENV_DIR" ]; then
    echo "Virtual environment not found. Creating a new one..."

    # Check if the creation was successful
    if python -m venv "$VENV_DIR"; then
        echo "Virtual environment created successfully."
    else
        echo "Failed to create virtual environment."
        exit 1
    fi
else
    echo "Virtual environment already exists."
fi

# Activate the virtual environment
# For macOS/Linux
source "$VENV_DIR/bin/activate"

# If on Windows, use the following line instead
# source "$VENV_DIR/Scripts/activate"

# Install the package with optional dependencies
pip install .

echo "Installation complete."

echo "Executing jupyter nbconvert to notebook"

NOTEBOOKS=$(find docs/ -name "*.ipynb")

echo "Clearing outputs..."
for nb in $NOTEBOOKS; do
    echo " - Clearing: $nb"
    jupyter nbconvert --ClearOutputPreprocessor.enabled=True --inplace "$nb"
done

echo "Executing notebooks..."
for nb in $NOTEBOOKS; do
    echo " - Executing: $nb"
    jupyter nbconvert --to notebook --execute --inplace "$nb"
done

echo "Notebooks executed and saved."
