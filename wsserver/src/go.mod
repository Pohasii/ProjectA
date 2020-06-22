module wslib

replace projecta.com/me/redis => ../../redis/src

go 1.14

require (
    github.com/gorilla/websocket v1.4.2
    projecta.com/me/redis v1.0.0
)