package repository

import (
	"time"

	"github.com/ddytert/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailibilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
}
