# package: chat_sync/use_case
mockgen -source chat_sync/use_case/reader.go -package use_case -destination chat_sync/use_case/mock_reader.go
mockgen -source chat_sync/use_case/writer.go -package use_case -destination chat_sync/use_case/mock_writer.go

# package: chat_sync/reader/pagination
mockgen -source chat_sync/reader/pagination/paginated_reader.go -package pagination -destination chat_sync/reader/pagination/mock_paginated_reader.go
mockgen -source chat_sync/reader/pagination/pagination_storage.go -package pagination -destination chat_sync/reader/pagination/mock_pagination_storage.go

# package: chat_sync/reader/buffer
mockgen -source chat_sync/reader/buffer/batch_reader.go -package buffer -destination chat_sync/reader/buffer/mock_batch_reader.go

# package: chat_sync/wecom
mockgen -source chat_sync/wecom/chat_record_service.go -package wecom -destination chat_sync/wecom_chat/mock_chat_record_service.go
mockgen -source chat_sync/wecom/open_api_service.go -package wecom -destination chat_sync/wecom_chat/mock_open_api_service.go
mockgen -source logger/logger.go -package wecom -destination chat_sync/wecom/mock_logger.go
mockgen -source chat_sync/wecom/chat_record_transformer.go -package wecom -destination chat_sync/wecom_chat/mock_chat_record_transformer.go

# package: chat_sync/http_controller
mockgen -source chat_sync/use_case/use_case.go -package http_controller -destination chat_sync/http_controller/mock_use_case.go
mockgen -package http_controller -destination chat_sync/http_controller/mock_response_writer.go net/http ResponseWriter
mockgen -source logger/logger.go -package http_controller -destination chat_sync/http_controller/mock_logger.go

# package: tencent_faas_adapter
mockgen -package tencent_faas_adapter -destination tencent_faas_adapter/mock_handler.go net/http Handler
mockgen -source logger/logger.go -package tencent_faas_adapter -destination tencent_faas_adapter/mock_logger.go

# package: chat_sync/chat_record_bitable_storage
mockgen -source logger/logger.go -package chat_record_bitable_storage -destination chat_sync/chat_record_bitable_storage/mock_logger.go

# package: chat_sync/wecom/chat_record_service
mockgen -source logger/logger.go -package chat_record_service -destination chat_sync/wecom/chat_record_service/mock_logger.go

# package: chat_sync/wecom_openapi_adapter
mockgen -source logger/logger.go -package wecom_openapi_adapter -destination chat_sync/wecom_openapi_adapter/mock_logger.go