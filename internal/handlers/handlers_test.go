package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ddytert/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{"home", "/", http.StatusOK},
	{"about", "/about", http.StatusOK},
	{"gq", "/generals-quarters", http.StatusOK},
	{"mj", "/majors-suite", http.StatusOK},
	{"sa", "/search-availability", http.StatusOK},
	{"contact", "/contact", http.StatusOK},
	{"non-existent route", "/green/red/blue", http.StatusNotFound},
	// New routes
	{"login", "/user/login", http.StatusOK},
	{"logout", "/user/logout", http.StatusOK},
	{"dashboard", "/admin/dashboard", http.StatusOK},
	{"new res", "/admin/reservations-new", http.StatusOK},
	{"all res", "/admin/reservations-all", http.StatusOK},
	{"show res", "/admin/reservations/new/1/show", http.StatusOK},

	// {"psa", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "01.01.2023"},
	// 	{key: "end", value: "02.01.2023"},
	// }, http.StatusOK},
	// {"psaj", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "01.01.2023"},
	// 	{key: "end", value: "02.01.2023"},
	// }, http.StatusOK},
	// {"pmr", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "Dan"},
	// 	{key: "last_name", value: "Deiers"},
	// 	{key: "email", value: "dan@deiers.de"},
	// 	{key: "phone", value: "0049157344658392"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		t.Logf("Testing handler %s", e.name)
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Generals Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusOK)
	}

	// Test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusSeeOther)
	}

	// Test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation.RoomID = 200

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	// prepare reservation
	layout := "02.01.2006"
	startDate, _ := time.Parse(layout, "27.10.2050")
	endDate, _ := time.Parse(layout, "31.10.2050")
	reservation := models.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
	}
	// prepare form post data
	postedData := url.Values{}
	postedData.Add("first_name", "Dan")
	postedData.Add("last_name", "Deiers")
	postedData.Add("email", "dan@deiers.de")
	postedData.Add("phone", "0049157344658392")

	body := strings.NewReader(postedData.Encode())

	req, _ := http.NewRequest("POST", "/make-reservation", body)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusSeeOther)
	}

	//test for missing reservation in session
	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", nil)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid form data
	postedData = url.Values{}
	postedData.Add("first_name", "Dan")
	postedData.Add("last_name", "Deiers")
	postedData.Add("email", "invalid")
	postedData.Add("phone", "0049157344658392")

	body = strings.NewReader(postedData.Encode())

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusSeeOther)
	}

	// Test failing of insertion of reservation into db
	reservation.RoomID = 777

	postedData = url.Values{}
	postedData.Add("first_name", "Dan")
	postedData.Add("last_name", "Deiers")
	postedData.Add("email", "dan@deiers.de")
	postedData.Add("phone", "0049157344658392")

	body = strings.NewReader(postedData.Encode())

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusSeeOther)
	}

	// Test failing of insertion of restriction into db
	reservation.RoomID = 555

	body.Reset(postedData.Encode())
	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// first case - rooms are not available
	postedData := url.Values{}
	postedData.Add("start", "27.10.2050")
	postedData.Add("end", "31.10.2050")
	postedData.Add("room_id", "1")

	// create request
	req, _ := http.NewRequest("POST", "search-availability-json", strings.NewReader((postedData.Encode())))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &j)

	if err != nil {
		t.Error("failed to parse json")
	}
	// room should not be available otherwise throw an error
	if j.OK {
		t.Error("room should not be available")
	} else if j.Message != "" {
		t.Error("message should be empty")
	}

	// test invalid body
	postedData = url.Values{}

	req, _ = http.NewRequest("POST", "search-availability-json", strings.NewReader((postedData.Encode())))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Error("failed to parse json")
	}
	// body should be unreadable, so an internal server error should be given back
	if j.Message != "Internal server error" {
		t.Error("Response should conatain Internal server error")
	}

	// test failing database request
	postedData = url.Values{}
	postedData.Add("start", "27.10.2050")
	postedData.Add("end", "31.10.2050")
	postedData.Add("room_id", "777")

	req, _ = http.NewRequest("POST", "search-availability-json", strings.NewReader((postedData.Encode())))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Error("failed to parse json")
	}
	// a database connection error should be given back
	if j.Message != "Error connecting to database" {
		t.Error("Response should conatain Error connecting to database")
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid ",
		"fred@furz.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestRepository_Login(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			acutalLoc, _ := rr.Result().Location()
			if acutalLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, acutalLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
