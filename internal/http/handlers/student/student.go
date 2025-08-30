package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/nkchakradhari780/students-api/internal/storage"
	"github.com/nkchakradhari780/students-api/internal/types"
	"github.com/nkchakradhari780/students-api/internal/utils/response"
)



func CreateNewStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid json body: %w", err)))
			return 
		}
		// request validation 
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.String("userId", fmt.Sprintf("%d", lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		response.WriteJson(w, http.StatusOK, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id",id))

		intId, err := strconv.ParseInt(id,10,64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}
		student, err := storage.GetStudentById( int64(intId) )

		if err != nil {
			slog.Error("error getting user", slog.String("Id",id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		students, err := storage.GetStudentList()
		if err != nil {
			fmt.Println("Error getting students: ", err)
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func DeleteById(storage storage.Storage) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting student")

		id := r.PathValue("id")
		intId, err := strconv.ParseInt(id,10,64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		err = storage.DeleteStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"result": "success"})
	}
}