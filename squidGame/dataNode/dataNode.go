package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "../comms"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDataNodeServer
}

func (s *Server) GuardarDatos(ctx context.Context, jugadaDataNode *pb.JugadaDataNode) (*pb.Empty, error) {
	nJug := strconv.Itoa(int(jugadaDataNode.Jugador))
	nRon := strconv.Itoa(int(jugadaDataNode.Ronda))
	juga := strconv.Itoa(int(jugadaDataNode.Jugada))
	//log.Printf("Recibido jugada, escogio el numero: %d", nEsc)
	frase := ("jugador_") + nJug + ("__ronda_") + nRon
	f, err := os.Create(frase + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(juga + "\n")

	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("done")
	return &pb.Empty{}, nil
}

func main() {

	//cantRondasJuego1 := 1

	fmt.Println("Soy el Lider!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9005))
	if err != nil {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", 9006))
		if err != nil {
			lis, err = net.Listen("tcp", fmt.Sprintf(":%d", 9007))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
		}
	}

	s := grpc.NewServer()

	pb.RegisterDataNodeServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
