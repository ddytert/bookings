package repository

import "github.com/ddytert/bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRetriction(r models.RoomRestriction) error
}
