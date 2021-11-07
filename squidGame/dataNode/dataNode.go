package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

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
	fmt.Println("listo")
	return &pb.Empty{}, nil
}

func (s *Server) PedirDatos(ctx context.Context, jugadorEtapa *pb.JugadorEtapa) (*pb.HistorialJugadas, error) {
	nJug := strconv.Itoa(int(jugadorEtapa.NumJug))
	nj := int(jugadorEtapa.NumJug)
	nEta := strconv.Itoa(int(jugadorEtapa.Etapa))
	etapa1 := []int32{0, 0, 0, 0}
	etapa2 := int(0)
	etapa3 := int(0)
	i := int(1)
	nomArch := ("jugador_") + nJug + ("__ronda_") + nEta
	f, err := os.Create(nomArch + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	cont := int(0)
	for scanner.Scan() {
		jugada := scanner.Text()
		nume := strings.Fields(jugada)
		if nEta == "1" {
			i1, err := strconv.Atoi(nume[0])
			if err == nil {
				fmt.Println(i1)
			}
			etapa1[cont] = int32(i1)
			cont = cont + 1
		} else if nEta == "2" {
			i2, err := strconv.Atoi(nume[0])
			if err == nil {
				fmt.Println(i2)
			}
			etapa2 = i2
		} else {
			i3, err := strconv.Atoi(nume[0])
			if err == nil {
				fmt.Println(i3)
			}
			etapa3 = i3
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &pb.HistorialJugadas{Jugador: int32(nj), Ronda1: int32(etapa1[0]), Ronda2: int32(etapa1[1]), Ronda3: int32(etapa1[2]), Ronda4: int32(etapa1[3]), Etapa2: int32(etapa2), Etapa3: int32(etapa3)}, nil
}

func main() {

	//cantRondasJuego1 := 1

	fmt.Println("DataNode Iniciado")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9005))
	if err != nil {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", 9006))
		if err != nil {
			lis, err = net.Listen("tcp", fmt.Sprintf(":%d", 9007))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			} else {
				fmt.Println("Puerto: 9007")
			}
		} else {
			fmt.Println("Puerto: 9006")
		}
	} else {
		fmt.Println("Puerto: 9005")
	}

	s := grpc.NewServer()

	pb.RegisterDataNodeServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al Servir: %s", err)
	}
}
