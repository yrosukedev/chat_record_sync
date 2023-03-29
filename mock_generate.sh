# package: use_case
mockgen -source use_case/chat_record_reader.go -package use_case -destination use_case/mock_chat_record_reader.go
mockgen -source use_case/chat_record_writer.go -package use_case -destination use_case/mock_chat_record_writer.go

# package: paginated_reader
mockgen -source paginated_reader/paginated_buffered_reader.go -package paginated_reader -destination paginated_reader/mock_paginated_buffered_reader.go
mockgen -source paginated_reader/pagination_storage.go -package paginated_reader -destination paginated_reader/mock_pagination_storage.go

# package: buffer_reader
mockgen -source buffer_reader/chat_record_buffered_reader.go -package buffer_reader -destination buffer_reader/mock_chat_record_buffered_reader.go

# package: wecom_chat
mockgen -source wecom_chat/chat_record_service.go -package wecom_chat -destination wecom_chat/mock_chat_record_service.go
mockgen -source wecom_chat/open_api_service.go -package wecom_chat -destination wecom_chat/mock_open_api_service.go
mockgen -source wecom_chat/chat_record_transformer.go -package wecom_chat -destination wecom_chat/mock_chat_record_transformer.go
mockgen -source logger/logger.go -package wecom_chat -destination wecom_chat/mock_logger.go

# package: http_controller
mockgen -source use_case/use_case.go -package http_controller -destination http_controller/mock_use_case.go
mockgen -package http_controller -destination http_controller/mock_response_writer.go net/http ResponseWriter
mockgen -source logger/logger.go -package http_controller -destination http_controller/mock_logger.go

# package: tencent_faas_adapter
mockgen -package tencent_faas_adapter -destination tencent_faas_adapter/mock_handler.go net/http Handler
mockgen -source logger/logger.go -package tencent_faas_adapter -destination tencent_faas_adapter/mock_logger.go

# package: chat_record_bitable_storage
mockgen -source logger/logger.go -package chat_record_bitable_storage -destination chat_record_bitable_storage/mock_logger.go

# package: wecom_chat_adapter
mockgen -source logger/logger.go -package wecom_chat_adapter -destination wecom_chat_adapter/mock_logger.go

# package: wecom_openapi_adapter
mockgen -source logger/logger.go -package wecom_openapi_adapter -destination wecom_openapi_adapter/mock_logger.go