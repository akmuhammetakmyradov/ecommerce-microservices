include makefiles/golang.mk
include makefiles/migration.mk
include makefiles/docker.mk

.PHONY: help

help:
	@echo "📘 Available commands:\n"
	@echo "🔹 Go Commands:"
	@grep -hE '^[a-zA-Z_-]+:.*?## ' makefiles/golang.mk Makefile | grep -E 'build|run|test|clean' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

	@echo "\n🐳 🔧 Dev Commands:"
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## ' makefiles/docker.mk Makefile | grep -E '^(create_network|docker_tag_push|dev_).*' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

	@echo "\n🐳 🚀 Production Commands:"
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## ' makefiles/docker.mk Makefile | grep -E '^(create_network|prod_).*' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

	@echo "\n🧬 Migration Commands:"
	@grep -hE '^[a-zA-Z_-]+:.*?## ' makefiles/migration.mk Makefile | grep -E 'migrate' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
