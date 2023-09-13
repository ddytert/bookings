package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
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
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation.RoomID = 200

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusTemporaryRedirect)
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
	reqBody := "first_name=Dan"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Deiers")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=dan@deiers.de")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0049157344658392")

	body := strings.NewReader(reqBody)

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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid form data
	reqBody = "first_name=Dan"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Deiers")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0049157344658392")

	body = strings.NewReader(reqBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test failing of insertion of reservation into db
	reservation.RoomID = 777

	reqBody = "first_name=Dan"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Deiers")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=dan@deiers.de")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0049157344658392")

	body = strings.NewReader(reqBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test failing of insertion of restriction into db
	reservation.RoomID = 555

	body.Reset(reqBody)
	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
