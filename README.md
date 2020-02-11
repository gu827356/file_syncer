## file_syncer
file_syncer is a simple tool used to synchronize directory content between machines.  

## Build and usage
By executing `build.sh`, you can find two executable files in `bin/` directory.  

#### file_syncer_server
options:    
`--root`: specify a directory that will be pulled by client.    
`--port`: the server port.   

#### file_syncer_client
options:  
`--addr`: the server socket address, example: `192.168.1.2:3333`  
`--root`: the root directory in client machine, that will used as the root directory.  