package dbrepo

import (
	"errors"
	"time"

	"github.com/ddytert/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if room id == 777 fail
	if res.RoomID == 777 {
		return 0, errors.New("new error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	// if room id == 777 fail
	if r.RoomID == 555 {
		return errors.New("new error")
	}
	return nil
}

// SearchAvailibilityByDatesByRoomID returns true if availibilty exists for room id, and false if no availibilty
func (m *testDBRepo) SearchAvailibilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 777 {
		return false, errors.New("new error")
	}
	return false, nil
}

// SearchAvailibilityForAllRooms returns a slice of available rooms, if any for given date range
func (m *testDBRepo) SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	rooms := []models.Room{}
	return rooms, nil
}

// GetRoomByID Gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 100 {
		return room, errors.New("no room")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

func (m *testDBRepo) UpdateUser(user models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// ReservationByID returns a reservation by id
func (m *testDBRepo) ReservationByID(reservationID int) (models.Reservation, error) {
	var reservation models.Reservation

	return reservation, nil
}
