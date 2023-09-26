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

	GetUserByID(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(reservationID int) (models.Reservation, error)
	UpdateReservation(res models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id int, processed int) error
}
