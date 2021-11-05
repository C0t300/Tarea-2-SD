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

func juego1() {
	ronda := 1
	// jugadores :=
	for ronda < 5 {
		nLider := rand.Intn(5) + 6
		ronda = ronda + 1

	}
}

func main() {

	cantRondasJuego1 := 1

	fmt.Println("Soy el Lider!")

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
