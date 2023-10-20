package main

import "log"

type Server struct {
	host string
	port int
}

func New(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (s *Server) Start() error {
	return nil
}

func main() {
	svr := server.New("localhost", 8888)
	if err := svr.Start(); err != nil {
		log.Fatal(err)
	}
}
