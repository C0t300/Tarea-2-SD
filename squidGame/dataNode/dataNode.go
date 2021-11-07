package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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

/*Crear un archivo*/
func createFile(nombreArchivo string) {
	myfile, err := os.Create(nombreArchivo)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Se creó", nombreArchivo)
	myfile.Close()
}

/*Anexar contenido al archivo*/
func appendToFile(text string, nombreArchivo string) {
	f, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

/*Revisar si el archivo existe*/
func fileExists(nombreArchivo string) bool {
	if _, err := os.Stat(nombreArchivo); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		return false

	} else {
		fmt.Println("wtf just happened (err 1)")
		fmt.Println(err.Error())
		return false
	}
}

/*Lectura de Archivos*/
func readFile(nombreArchivo string) []string {
	if fileExists(nombreArchivo) {
		content, err := ioutil.ReadFile(nombreArchivo)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("err 2")
		}
		if len(content) == 0 {
			return []string{}
		}
		contentString := string(content)
		parsed := strings.Split(contentString, "\n")
		return parsed
	}
	fmt.Println("algo paso err 12389")
	return []string{}
}

/*GuardarDatos: Funcion que recibe un mensaje JugadaDataNote para luego crear y escribir el archivo respectivo
al jugador y la etapa con la información de sus jugadas.*/
func (s *Server) GuardarDatos(ctx context.Context, jugadaDataNode *pb.JugadaDataNode) (*pb.Empty, error) {
	nJug := strconv.Itoa(int(jugadaDataNode.Jugador))
	nRon := strconv.Itoa(int(jugadaDataNode.Ronda))
	juga := strconv.Itoa(int(jugadaDataNode.Jugada))
	//log.Printf("Recibido jugada, escogio el numero: %d", nEsc)
	frase := ("jugador_") + nJug + ("__ronda_") + nRon + ".txt"
	bol := fileExists(frase)
	linea := (juga + "\n")
	if !bol {
		createFile(frase)
		appendToFile(linea, frase)

	} else {

		appendToFile(linea, frase)
	}
	/*f, err := os.Create(frase + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(juga + "\n")

	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("listo")*/
	return &pb.Empty{}, nil
}

/*Pedir Datos: funcion que recibe un mensaje JugadorEtapa y lee el archivo relacionado a ese jugador y etapa
retorna un mensaje de HistorialJugadas con la info relacionada a esa etapa*/
func (s *Server) PedirDatos(ctx context.Context, jugadorEtapa *pb.JugadorEtapa) (*pb.HistorialJugadas, error) {
	fmt.Println("ENtrango pedir datos")
	nJug := strconv.Itoa(int(jugadorEtapa.Jugador))
	nj := int(jugadorEtapa.Jugador)
	nEta := strconv.Itoa(int(jugadorEtapa.Etapa))
	etapa1 := []int32{0, 0, 0, 0}
	etapa2 := int(0)
	etapa3 := int(0)

	nomArch := ("jugador_") + nJug + ("__ronda_") + nEta + (".txt")
	if !fileExists(nomArch) {
		fmt.Println("laskdjlkdj no existe el archivo")
	}
	if nEta == "1" {
		readFile(nomArch)
		i := 0
		for i < len(nomArch) {
			buf, err := strconv.Atoi(string(nomArch[i]))
			if err != nil {
				log.Fatalf("error añslkjdaslkj: %s", err)
			}
			etapa1[i] = int32(buf)
		}
	} else if nEta == "2" {
		buf2 := string(readFile(nomArch)[0])
		buf, err := strconv.Atoi(buf2)
		if err != nil {
			log.Fatalf("error añslkjdaslkj: %s", err)
		}
		etapa2 = int(buf)
	} else {
		buf3 := string(readFile(nomArch)[0])
		buf, err := strconv.Atoi(buf3)
		if err != nil {
			log.Fatalf("error añslkjdaslkj: %s", err)
		}
		etapa3 = int(buf)

	}
	/* f, err := os.Open(nomArch + ".txt")

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

	fmt.Println(nj) */

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
