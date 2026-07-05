include .env
export

# .SILENT
# .ONESHELL
# .NOTPARALLEL
# .DEFAULT_GOAL := target-name

export PROJECT_ROOT := $(CURDIR)
export LOGS_DIR := $(PROJECT_ROOT)/out/logs

# ===========================App===========================
.PHONY: run
run:
	@docker compose up --build $(APP_SERVICE_NAME)



# ===========================Logs===========================
.PHONY: cleanup-logs
cleanup-logs:
	@read -p "Would you like to cleanup all logs? DANGER OPERATION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf $(LOGS_DIR); \
	fi;

