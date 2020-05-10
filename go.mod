module mod

replace projecta.com/me/chat => ../projecta/chat
replace projecta.com/me/client => ../projecta/client
replace projecta.com/me/verification => ../projecta/verification
replace projecta.com/me/wsserver => ../projecta/wsserver
replace projecta.com/me/setenv => ../projecta/setenv

require (
	projecta.com/me/chat v1.0.0
    projecta.com/me/client v1.0.0
    projecta.com/me/verification v1.0.0
    projecta.com/me/wsserver v1.0.0
    projecta.com/me/setenv v1.0.0
)

go 1.14