package main

import (
	"fmt"
	"log"

	pb "../comms"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func connectGRPC() pb.Juego1Client {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewJuego1Client(conn)
	return c
}

func QuieroJugarJugador(c pb.Juego1Client) int {
	fmt.Printf("Enviada solicitud de juego.")
	response, err := c.QuieroJugar(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	return int(response.NumJug)
}

func main() {

	c := connectGRPC()
	nJugador := QuieroJugarJugador(c)
	fmt.Printf("Juego iniciado. Jugador numero: %d", nJugador)

}
