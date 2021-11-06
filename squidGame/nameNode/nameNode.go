package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "../comms"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedMensajeDataLiderServer
}

func (s *Server) Jugada(ctx context.Context, jugadaDataNode *pb.JugadaDataNode) (*pb.Empty, error) {
	nJug := strconv.Itoa(int(jugadaDataNode.Jugador))
	nRon := strconv.Itoa(int(jugadaDataNode.Ronda))
	//log.Printf("Recibido jugada, escogio el numero: %d", nEsc)
	frase := ("Jugador_") + nJug + (" Ronda_") + nRon
	f, err := os.Create("NameNodeData.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(frase + "\n")

	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("done")
	return &pb.Empty{}, nil
}

func main() {

	fmt.Println("nameNode arriba")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9001))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterMensajeDataLiderServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
