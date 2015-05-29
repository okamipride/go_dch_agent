README
================================

## dch_agent 
#### Run 
	go run dch_agent.go [-option] 
	./dch_agent [-option] (Linux)

#### options 
	- -serv 

		relay server address or ip , no port included

		default : r0401.dch.dlink.com

	- -dev 

		number of devices want to connect to relay server


		default : 1

	- -concur

		concurrent send without delay

		default : 1


	- -delay 

		delay between concurrent send

		default : 10ms

	- -log

		turn log on(true) off(false) 

		default: false

	
#### examples 
	
	go run dch_agent.go -serv=r0401.dch.dlink.com -dev=100 -concur=1 -delay=100 -log=false

	go run dch_agent.go -serv=172.31.4.183 -dev=100 -concur=1 -delay=100 -log=true


## dch_app_http


## dch_app_socket