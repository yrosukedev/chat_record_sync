# package: use_case
mockgen -source use_case/chat_record_reader.go -package use_case -destination use_case/mock_chat_record_reader.go
mockgen -source use_case/chat_record_writer.go -package use_case -destination use_case/mock_chat_record_writer.go