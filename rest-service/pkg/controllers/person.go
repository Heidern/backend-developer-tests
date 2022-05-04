package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	uuid "github.com/satori/go.uuid"

	"github.com/stackpath/backend-developer-tests/rest-service/pkg/errors"
	"github.com/stackpath/backend-developer-tests/rest-service/pkg/models"
)

type PersonController struct{}

func (c *PersonController) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", c.GetPeople)
	r.Get("/{id}", c.GetPersonByID)

	return r
}

func (c *PersonController) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	idAsString := chi.URLParam(r, "id")

	id, err := uuid.FromString(idAsString)

	if err != nil {
		errors.RenderError(w, r, &errors.InvalidArgumentError{Err: err, Message: "id is not a valid uuid"})
		return
	}

	person, err := models.FindPersonByID(id)

	if err != nil {
		errors.RenderError(w, r, &errors.ResourceNotFoundError{Err: err, Message: "person not found"})
		return
	}

	render.JSON(w, r, person)
}

func (c *PersonController) GetPeople(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	phoneNumber := r.URL.Query().Get("phone_number")

	firstAndLastNameFilterSpecified := len(firstName) > 0 || len(lastName) > 0
	phoneFilterSpecified := len(phoneNumber) > 0

	// if we want to allow these filters to be combined then we should combine the results from
	// lower layer calls

	if firstAndLastNameFilterSpecified && phoneFilterSpecified {
		errors.RenderError(w, r, &errors.InvalidArgumentError{Message: "cannot combine first_last, last_name, and phone_number filters"})
		return
	}

	if firstAndLastNameFilterSpecified {
		if len(firstName) == 0 {
			errors.RenderError(w, r, &errors.InvalidArgumentError{Message: "missing first_name"})
			return
		}
		if len(lastName) == 0 {
			errors.RenderError(w, r, &errors.InvalidArgumentError{Message: "missing last_name"})
			return
		}

		render.JSON(w, r, models.FindPeopleByName(firstName, lastName))
		return
	}

	if phoneFilterSpecified {
		render.JSON(w, r, models.FindPeopleByPhoneNumber(phoneNumber))
		return
	}

	render.JSON(w, r, models.AllPeople())
}
