module mod

replace projecta.com/me/chat => ../projecta/chat
replace projecta.com/me/client => ../projecta/client
//replace projecta.com/me/wsserver => ../projecta/wsserver

require (
	projecta.com/me/chat v1.0.0
    projecta.com/me/client v1.0.0
//    projecta.com/me/wsserver v1.0.0
)

go 1.14