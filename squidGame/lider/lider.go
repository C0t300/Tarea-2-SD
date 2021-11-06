package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	pb "../comms"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedJuego1Server
}

func (s *Server) Jugada1(ctx context.Context, jugadorJuego1 *pb.JugadorJuego1) (*pb.EstadoJuego, error) {
	nEsc := int(jugadorJuego1.EscogidoJugador)
	log.Printf("Recibido jugada, escogio el numero: %d", nEsc)
	nLider := rand.Intn(5) + 6
	dead := false
	if nLider >= nEsc {
		dead = true
	}

	return &pb.EstadoJuego{Vivo: dead, EscogidoLider: int32(nLider), Win: false, Round: int32(1)}, nil
}

/* func juego1() {
	ronda := 1
	// jugadores :=
	for ronda < 5 {
		nLider := rand.Intn(5) + 6
		ronda = ronda + 1

	}
} */

func main() {

	fmt.Println("Soy el Lider!")
	//parte cliente Lider-nameNode
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se conecta con el nameNode: %s", err)
	}
	defer conn.Close()

	c := pb.NewMensajeDataLiderClient(conn)

	response, err := c.Jugada(context.Background(), &pb.JugadaDataNode{Jugador: 1, Ronda: 1, Jugada: 1})
	if err != nil {
		log.Fatalf("Error when calling nameNode: %s", err)
	}
	log.Printf("Respuesta desde nameNode: %d", response.EscogidoLider)
	//parte Servidor Lider-Jugadores
	//cantRondasJuego1 := 1
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterJuego1Server(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
