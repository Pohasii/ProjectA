module mod

go 1.14

replace projecta.com/me/verification => ./verification/

require (
	github.com/go-redis/redis v6.15.8+incompatible // indirect
	projecta.com/me/verification v1.0.0
)
