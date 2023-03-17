# package: use_case
mockgen -source use_case/chat_record_reader.go -package use_case -destination use_case/mock_chat_record_reader.go
mockgen -source use_case/chat_record_writer.go -package use_case -destination use_case/mock_chat_record_writer.go

# package: paginated_reader
mockgen -source paginated_reader/paginated_buffered_reader.go -package paginated_reader -destination paginated_reader/mock_paginated_buffered_reader.go
mockgen -source paginated_reader/pagination_storage.go -package paginated_reader -destination paginated_reader/mock_pagination_storage.go