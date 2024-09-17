.PHONY: all clean

all: stop_mosquitto kill_ports

stop_mosquitto:
	@echo "Stopping mosquitto service..."
	sudo systemctl stop mosquitto
	@echo "Checking mosquitto service status:"
	sudo systemctl status mosquitto --no-pager | head -n 3
	@echo "Verifying if port 1883 is closed:"
	if ! sudo lsof -i:1883 | grep LISTEN > /dev/null; then \
		echo "Port 1883 is closed"; \
	else \
		echo "Port 1883 is still open"; \
	fi

kill_ports:
	@echo "Killing processes on ports 5432 and 1883..."
	@for PORT in 5432 1883; do \
		echo "Processing port $$PORT"; \
		PID=$$(sudo lsof -nP -i:$$PORT | awk '/LISTEN/{print $2}' | head -n 1); \
		if [ -n "$$PID" ]; then \
			echo "Found PID $$PID for port $$PORT"; \
			sudo kill -9 $$PID; \
			if [ $? -eq 0 ]; then \
				echo "Successfully killed process with PID $$PID listening on port $$PORT"; \
			else \
				echo "Error killing process with PID $$PID"; \
			fi; \
		else \
			echo "No process found listening on port $$PORT"; \
		fi; \
	done

clean:
	@echo "Cleaning up..."
	# Add any cleanup commands here if needed

