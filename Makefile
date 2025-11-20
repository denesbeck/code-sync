.PHONY: install-hooks
install-hooks:
	@echo "> Installing Git hooks..."
	@cp scripts/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "> Pre-commit hook installed successfully!"
	@echo "\nINFO: The hook will run 'go vet' and 'gofmt' checks before each commit."

.PHONY: uninstall-hooks
uninstall-hooks:
	@echo "> Uninstalling Git hooks..."
	@rm -f .git/hooks/pre-commit
	@echo "> Pre-commit hook uninstalled."
