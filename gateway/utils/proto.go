package utils

import (
	"fmt"
	dto "gateway/DTO"
	"gateway/gen"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MarshalBookTicketRequest(tranID, userID, ticketNumber uint) ([]byte, error) {
	bookTicketRequest := &gen.BookTicketRequest{
		TrainId:      uint32(tranID),
		UserId:       uint32(userID),
		TicketNumber: uint32(ticketNumber),
	}

	data, err := proto.Marshal(bookTicketRequest)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnmarshalTickets(data []byte) ([]dto.TicketDTO, error) {
	protoTickets := &gen.ListTickets{}

	err := proto.Unmarshal(data, protoTickets)
	if err != nil {
		return nil, err
	}

	ticketsDTO := []dto.TicketDTO{}

	for _, protoTicket := range protoTickets.Tickets {
		ticketDTO := *convertProtoTicketToDTOTicket(protoTicket)

		ticketsDTO = append(ticketsDTO, ticketDTO)
	}

	return ticketsDTO, nil
}

func UnmarshalTicket(data []byte) (*dto.TicketDTO, error) {
	protoTicket := &gen.Ticket{}

	err := proto.Unmarshal(data, protoTicket)
	if err != nil {
		return nil, err
	}

	ticketDTO := convertProtoTicketToDTOTicket(protoTicket)

	return ticketDTO, nil
}

func UnmarshalListUsersReplay(data []byte) ([]dto.UserDTO, error) {
	protoUsers := &gen.ListUsersReplay{}

	err := proto.Unmarshal(data, protoUsers)
	if err != nil {
		return nil, err
	}

	userDTOs := []dto.UserDTO{}
	for _, protoUser := range protoUsers.Users {
		userDTO := dto.UserDTO{
			ID:        uint(protoUser.ID),
			FirstName: protoUser.FirstName,
			LastName:  protoUser.LastName,
			Email:     protoUser.Email,
		}

		userDTOs = append(userDTOs, userDTO)
	}

	return userDTOs, nil
}

func UnmarshalUser(data []byte) (*dto.UserDTO, error) {
	protoUser := &gen.User{}

	err := proto.Unmarshal(data, protoUser)
	if err != nil {
		return nil, err
	}

	userDTO := &dto.UserDTO{
		ID:        uint(protoUser.ID),
		FirstName: protoUser.FirstName,
		LastName:  protoUser.LastName,
		Email:     protoUser.Email,
	}

	return userDTO, nil
}

func MarshalCreateUserRequest(firstName, lastName, email string) ([]byte, error) {
	protoUser := &gen.CreateUserRequest{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	data, err := proto.Marshal(protoUser)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalUpdateUserRequest(userID uint, firstName, lastName string) ([]byte, error) {
	protoUser := &gen.UpdateUserRequest{
		ID:        uint32(userID),
		FirstName: firstName,
		LastName:  lastName,
	}

	data, err := proto.Marshal(protoUser)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnmarshalTrain(data []byte) (*dto.TrainDTO, error) {
	protoTrain := &gen.Train{}

	err := proto.Unmarshal(data, protoTrain)
	if err != nil {
		return nil, err
	}

	DTOTrain := convertTrainProtoToTrainDTO(protoTrain)

	return DTOTrain, nil
}

func UnmarshalListTrainsReplay(data []byte) ([]dto.TrainDTO, error) {
	protoTrains := &gen.ListTrainsReplay{}

	err := proto.Unmarshal(data, protoTrains)
	if err != nil {
		return nil, err
	}

	trainDTOs := []dto.TrainDTO{}

	for _, protoTrain := range protoTrains.Trains {
		trainDTO := convertTrainProtoToTrainDTO(protoTrain)

		trainDTOs = append(trainDTOs, *trainDTO)
	}

	return trainDTOs, nil
}

func MarshalListTrainsFilteredRequest(trainFilter map[string]string) ([]byte, error) {
	availableSeats, _ := strconv.Atoi(trainFilter["available_seats"])

	departureTimeParsed, err := time.Parse(time.DateTime, trainFilter["departure_time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse departure time: %w", err)
	}

	arrivalTimeParsed, err := time.Parse(time.DateTime, trainFilter["arrival_time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse arrival time: %w", err)
	}

	protoFilter := &gen.TrainFilter{
		Name:           trainFilter["name"],
		AvailableSeats: uint32(availableSeats),
		Origin:         trainFilter["origin"],
		Destination:    trainFilter["destination"],
		DepartureTime:  timestamppb.New(departureTimeParsed),
		ArrivalTime:    timestamppb.New(arrivalTimeParsed),
	}

	data, err := proto.Marshal(protoFilter)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalCreateTrainRequest(name string, capacity uint) ([]byte, error) {
	protoTrain := &gen.CreateTrainRequest{
		Name:     name,
		Capacity: uint32(capacity),
	}

	data, err := proto.Marshal(protoTrain)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalUpdateTrainNameRequest(trainID uint, name string) ([]byte, error) {
	protoRequest := &gen.UpdateTrainNameRequest{
		ID:   uint32(trainID),
		Name: name,
	}

	data, err := proto.Marshal(protoRequest)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalUpdateTrainTravelDetailsRequest(trainID uint, travelDetails map[string]string) ([]byte, error) {
	departureTimeParsed, err := ParseTime(travelDetails["departure_time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse departure time: %w", err)
	}

	arrivalTimeParsed, err := ParseTime(travelDetails["arrival_time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse arrival time: %w", err)
	}

	protoRequest := &gen.UpdateTrainTravelDetailsRequest{
		ID:            uint32(trainID),
		Origin:        travelDetails["origin"],
		Destination:   travelDetails["destination"],
		DepartureTime: timestamppb.New(*departureTimeParsed),
		ArrivalTime:   timestamppb.New(*arrivalTimeParsed),
	}

	data, err := proto.Marshal(protoRequest)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnmarshalSeat(data []byte) (*dto.SeatDTO, error) {
	protoSeat := &gen.Seat{}

	err := proto.Unmarshal(data, protoSeat)
	if err != nil {
		return nil, err
	}

	seatDTO := convertProtoSeatToDTOSeat(protoSeat)

	return seatDTO, nil
}

func UnmarshalListSeatsByTrainIDReplay(data []byte) ([]dto.SeatDTO, error) {
	protoSeats := &gen.ListSeatsReplay{}

	err := proto.Unmarshal(data, protoSeats)
	if err != nil {
		return nil, err
	}

	seatDTOs := convertProtoSeatsToDTOSeats(protoSeats.Seats)

	return seatDTOs, nil
}

func MarshalCreateSeatRequest(trainID, seatNumber uint) ([]byte, error) {
	protoSeat := &gen.CreateSeatRequest{
		TrainId: uint32(trainID),
		Number:  uint32(seatNumber),
	}

	data, err := proto.Marshal(protoSeat)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalUpdateSeatNumberRequest(seatID, seatNumber uint) ([]byte, error) {
	protoSeat := &gen.UpdateSeatNumberRequest{
		ID:     uint32(seatID),
		Number: uint32(seatNumber),
	}

	data, err := proto.Marshal(protoSeat)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func convertTrainProtoToTrainDTO(protoTrain *gen.Train) *dto.TrainDTO {
	travelDetailDTO := dto.TravelDetailDTO{
		Origin:        protoTrain.Origin,
		Destination:   protoTrain.Destination,
		DepartureTime: protoTrain.DepartureTime.AsTime(),
		ArrivalTime:   protoTrain.ArrivalTime.AsTime(),
	}

	DTOSeats := convertProtoSeatsToDTOSeats(protoTrain.Seats)

	return &dto.TrainDTO{
		ID:              uint(protoTrain.ID),
		Name:            protoTrain.Name,
		Capacity:        uint(protoTrain.Capacity),
		Seats:           DTOSeats,
		IsFull:          protoTrain.IsFull,
		TravelDetailDTO: travelDetailDTO,
	}
}

func convertProtoSeatsToDTOSeats(protoSeats []*gen.Seat) []dto.SeatDTO {
	dtoSeats := make([]dto.SeatDTO, 0, len(protoSeats))

	for _, protoSeat := range protoSeats {
		dtoSeats = append(dtoSeats, *convertProtoSeatToDTOSeat(protoSeat))
	}

	return dtoSeats
}

func convertProtoSeatToDTOSeat(protoSeat *gen.Seat) *dto.SeatDTO {
	return &dto.SeatDTO{
		ID:      uint(protoSeat.ID),
		Number:  uint(protoSeat.SeatNumber),
		TrainID: uint(protoSeat.TrainId),
		UserID:  uint(protoSeat.UserId),
	}
}

func ParseTime(str string) (*time.Time, error) {
	time, err := time.Parse(time.DateTime, str)
	if err != nil {
		return nil, err
	}

	return &time, nil
}

func convertProtoTicketToDTOTicket(protoTicket *gen.Ticket) *dto.TicketDTO {
	travelDetailDTO := dto.TravelDetailDTO{
		Origin:        protoTicket.Origin,
		Destination:   protoTicket.Destination,
		DepartureTime: protoTicket.DepartureTime.AsTime(),
		ArrivalTime:   protoTicket.ArrivalTime.AsTime(),
	}

	ticketDTO := &dto.TicketDTO{
		ID:              uint(protoTicket.ID),
		SeatNumber:      uint(protoTicket.SeatNumber),
		UserID:          uint(protoTicket.UserId),
		TrainID:         uint(protoTicket.TrainId),
		ExpiresAt:       protoTicket.ExpiresAt.AsTime(),
		TravelDetailDTO: travelDetailDTO,
	}

	canceledAt := protoTicket.CanceledAt.AsTime()
	if !canceledAt.IsZero() {
		ticketDTO.CanceledAt = &canceledAt
	}

	return ticketDTO
}
