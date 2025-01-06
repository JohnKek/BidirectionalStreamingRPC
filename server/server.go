package main

import (
	api "chat/api/grpc"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"math/rand"
	"net"
)

type Server struct {
	api.UnimplementedBattleshipServiceServer
}

func (s *Server) Game(stream api.BattleshipService_GameServer) error {
	msg, err := stream.Recv()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to receive: %v", err))
		return err
	}
	startGame := msg.GetStart()
	if startGame == nil {
		logger.Error(fmt.Sprintf("Failed to receive startGame: %v", err))
		return status.Errorf(codes.InvalidArgument, "You must provide a startGame")
	}
	if startGame.FieldSize <= 1 {
		logger.Error(fmt.Sprintf("Field size must be greater than 1"))
		return status.Errorf(codes.InvalidArgument, "You must provide a fieldSize greater than 1")
	}

	battleField := make([][]int, startGame.FieldSize)
	for i := range battleField {
		battleField[i] = make([]int, startGame.FieldSize)
	}
	battleField[rand.Intn(int(startGame.FieldSize))][rand.Intn(int(startGame.FieldSize))] = 1
	fmt.Println(battleField)

outerLoop: // Метка для внешнего цикла
	for {
		msgPack, errStream := stream.Recv()
		if errStream != nil {
			if errStream == io.EOF {
				logger.Error("Stream closed")
				break
			}
			logger.Error(errStream.Error())
			break
		}
		cord := msgPack.GetCoordinate()
		if cord == nil {
			continue
		}
		logger.Info(fmt.Sprintf("Received: %v", cord))
		if cord.X < 0 || cord.X >= startGame.FieldSize || cord.Y < 0 || cord.Y >= startGame.FieldSize {
			logger.Error(fmt.Sprintf("Invalid coordinate: %v", cord))
			return status.Errorf(codes.OutOfRange, "Invalid coordinate")
		}
		switch battleField[cord.X][cord.Y] {
		case 1:
			err = stream.Send(&api.AttackInformation{Status: "You win"})
			if err != nil {
				return err
			}
			break outerLoop
		case 0:
			battleField[cord.X][cord.Y] = 2
			err = stream.Send(&api.AttackInformation{Status: "You missed"})
			if err != nil {
				return err
			}
		case 2:
			err = stream.Send(&api.AttackInformation{Status: "You have already attacked this coordinate"})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func StartGrpcServer() error {
	lis, err := net.Listen("tcp",
		fmt.Sprintf("%s:%d", "localhost", 8080))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen: %v", err))
		return err
	}
	s := grpc.NewServer()
	api.RegisterBattleshipServiceServer(s, &Server{})
	logger.Info(fmt.Sprintf("Starting server on port %d", 8080))
	if err = s.Serve(lis); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		return err
	}
	return nil
}
