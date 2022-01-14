module example.com/client

go 1.15

replace example.com/chat => ../chat

require (
	example.com/chat v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20220107192237-5cfca573fb4d
	google.golang.org/grpc v1.43.0
)
