package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EupravaProjekat/police-citizens/Models"
	"github.com/EupravaProjekat/police-citizens/Repo"
	"github.com/google/uuid"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"
)

type Borderhendler struct {
	l    *log.Logger
	repo *Repo.Repo
}

func NewBorderhendler(l *log.Logger, r *Repo.Repo) *Borderhendler {
	return &Borderhendler{l, r}

}

func (h *Borderhendler) CheckIfUserExists(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	response, err := h.repo.GetByEmail(re.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if re.Email != response.Email {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

}

func (h *Borderhendler) NewUser(w http.ResponseWriter, r *http.Request) {

	res := ValidateJwt2(r, h.repo)

	rt, err := DecodeBodyUser(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	rt.Email = res
	rt.Role = "Operator"
	err = h.repo.NewUser(rt)
	if err != nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Profile not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Borderhendler) NewWeaponRequest(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := DecodeBodyWeapon(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var request Models.Request
	request.RequestState = "PROCESSING"

	var formatted = time.Time(time.Now()).Format("2006-01-02")
	request.RequestDate = formatted

	request.Weapon = *rt

	request.Email = res.Email

	newUUID := uuid.New().String()
	request.Uuid = newUUID

	userData := []byte(`{"jmbg":"` + res.JMBG + `"}`)

	apiUrl := "http://prosecution-service:9199/check-if-person-is-prosecuted"
	request2, err2 := http.NewRequest("POST", apiUrl, bytes.NewBuffer(userData))
	request2.Header.Set("Content-Type", "application/json; charset=utf-8")
	request2.Header.Set("jwt", r.Header.Get("jwt"))
	request2.Header.Set("intern", "police-service-secret-code")

	// send the request
	client := &http.Client{}
	response, err2 := client.Do(request2)
	if err2 != nil {
		fmt.Println(err2)
	}
	responseBody, err2 := io.ReadAll(response.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	var resp Models.Responsepros
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	request.Recorded = strconv.FormatBool(resp.Prosecuted)
	defer response.Body.Close()
	res.Requests = append(res.Requests, request)
	err = h.repo.UpdateUser(*res)
	if err != nil {
		log.Printf("Operation failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("couldn't add request"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("successfully added request"))
	if err != nil {
		return
	}
}

func (h *Borderhendler) GetallRequests(w http.ResponseWriter, r *http.Request) {

	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("user doesnt exist")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if res.Role != "Operator" {
		err := errors.New("role error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.repo.GetAllRequest()
	if err != nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Requests not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}
func (h *Borderhendler) CheckPlatesWanted(w http.ResponseWriter, r *http.Request) {

	_ = ValidateJwt2(r, h.repo)

	if r.Header.Get("intern") != "border-service-secret-code" {
		err := errors.New("you are not system user, incident will be reported")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	rt, err := DecodeBodyplates(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	response, err := h.repo.CheckPlatesWanted(*rt)
	if err != nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Requests not found"))
		if err != nil {
			return
		}
		return
	}
	var resp Models.Response
	resp.VehicleWanted = response
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, resp)
}
func (h *Borderhendler) GetAllPlatesWnated(w http.ResponseWriter, r *http.Request) {

	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("user doesnt exist")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if res.Role != "Operator" {
		err := errors.New("role error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.repo.GetAllPlatesWnated()
	if err != nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Requests not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}
func (h *Borderhendler) NewWantedPlates(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := DecodeBodyplates(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	err = h.repo.NewPlatesWanted(rt)
	if err != nil {
		log.Printf("Operation failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("couldn't add request"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("successfully added request"))
	if err != nil {
		return
	}
}
