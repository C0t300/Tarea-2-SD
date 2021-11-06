package main

import (
	"log"

	"../comms"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := comms.NewJuego1Client(conn)

	response, err := c.Jugada1(context.Background(), &comms.JugadorJuego1{EscogidoJugador: 1})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %d", response.EscogidoLider)

}
