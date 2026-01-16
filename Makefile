#? help: ヘルプコマンド
help: Makefile
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^#?//p' $< | awk -F ':' '{ printf "  %-15s %s\n", $$1, $$2 }'
.PHONY: help

#? dev-up: 開発環境用サーバーを起動
dev-up:
	air
.PHONY: dev-up
