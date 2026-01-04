
swagger: ## 生成统一Swagger文档
	@echo "📚 生成统一Swagger文档..."
	make swagger-api 
	make swagger-admin 

swagger-api: ## 生成API Swagger文档 (使用Api标签)
	@echo "📚 生成API Swagger文档..."
	swag init -g main/main.go \
		--instanceName api \
		--tags "Api" \
		--parseDependency --parseInternal \
		-o ./docs/api

swagger-admin: ## 生成Admin Swagger文档 (使用Admin标签)
	@echo "📚 生成Admin Swagger文档..."
	swag init -g main/main.go \
		--instanceName admin \
		--tags "Admin" \
		--parseDependency --parseInternal \
		-o ./docs/admin
