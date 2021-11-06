# Tarea 2

## Instrucciones de conexion

`Estimado Grupo 23,`

`Las máquinas virtuales asignadas para su grupo son las siguientes:`

`Máquina= dist89`
`Clave= Mx%8YAK4RT]t<uV7`

`Máquina= dist90`
`Clave= fP2DU-xkA$e8.]dv`

`Máquina= dist91`
`Clave= HK9d@xL?aTR8fjw7`

`Máquina= dist92`
`Clave= yC7(/kJ:.TLmU]8z`

`El usuario que deben ocupar para conectarse a sus máquinas es: 'alumno'`

![img](https://cdn.discordapp.com/attachments/887425574195310594/906198590551121940/unknown.png)

## Definiciones conexiones

Jugadores <-> lider [GRPC]

nameNode <-> lider [grpc]

nameNode <-> datanode [gRPC]

lider <-> pozo [grpc] [rmq]

- 10.6.40.230
- 

## Maquinas

89: 10.6.40.230 [installed/tested gRPC and rMQ client ]

90: 10.6.40.231

91: 10.6.40.232 [installed/tested gRPC and rMQ Server (docker)]

92:  

89 <-> 91 gRPC

89 <-> 91 rMQ

## Compilar

```bash
protoc --go_out=. comms.proto
protoc --go-grpc_out=. comms.proto
```



## Puertos

- Pozo :9002 gRPC
- Pozo :5672 rMQ
- lider-namenode 9001
- datanodes 9005-7
- lider - jugador 9003

## Paginas útiles

- https://docs.docker.com/engine/install/centos/ [instalar docker]
- https://www.rabbitmq.com/download.html [instalar rabbitmq (yo lo hice en docker porque si)]
- https://grpc.io/docs/languages/go/quickstart/ [tutorial gRPC con instalacion de todo]
- https://www.rabbitmq.com/tutorials/tutorial-one-go.html [tutorial rMQ en Go]
- https://oracle-max.com/como-abrir-puertos-en-el-firewall-de-oracle-linux/ [abrir puertos en Linux]

