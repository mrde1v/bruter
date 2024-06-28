# Variables
BINARY_NAME=bruter.exe
BUILD_DIR=build
MAIN_FILE=./main.go

# Phony targets
.PHONY: build clean run

# Default target to build the application
build:
	@echo "Building the application..."
	@if not exist $(BUILD_DIR) (mkdir $(BUILD_DIR))
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Successfully built the application."

# Clean the build directory
clean:
	@echo "Cleaning the build directory..."
	@if exist $(BUILD_DIR) (rmdir /s /q $(BUILD_DIR))
	@echo "Successfully cleaned the build directory."

# Run the application
run: build
	@echo "Running the application..."
	@$(BUILD_DIR)/$(BINARY_NAME)
	@echo "Successfully ran the application."