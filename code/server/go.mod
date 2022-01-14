module example.com/server

go 1.15

require (
	example.com/chat v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	google.golang.org/grpc v1.43.0
)

replace example.com/chat => ../chat
