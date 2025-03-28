package handler

import (
	"PoliSim/database"
	"PoliSim/helper"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func GetManagementPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !(loggedIn && acc.Role == database.RootAdmin) {
		GetNotFoundPage(writer, request)
		return
	}

	MakeFullPage(writer, acc, &AdminPage{})
}

func PostFileManagementPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !(loggedIn && acc.Role == database.RootAdmin) {
		PartialGetNotFoundPage(writer, request)
		return
	}

	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error(err.Error())
		MakeSpecialPagePartWithRedirect(writer, &MessageUpdate{IsError: true,
			Message: "error while trying to parse multipart form"})
		return
	}

	file, handler, err := request.FormFile("file")
	if errors.Is(err, http.ErrMissingFile) {
		slog.Debug("No file was send")
	} else if err != nil {
		slog.Error(err.Error())
		MakeSpecialPagePartWithRedirect(writer, &MessageUpdate{IsError: true,
			Message: "error while trying to open the send file"})
		return
	} else {
		defer file.Close()
		slog.Debug("File Info:", "name", handler.Filename, "size", handler.Size)

		var target *os.File
		target, err = os.Create("./public/" + handler.Filename)
		defer target.Close()
		if err != nil {
			slog.Error(err.Error())
			MakeSpecialPagePartWithRedirect(writer, &MessageUpdate{IsError: true,
				Message: "error while trying to open the file path on the server"})
			return
		}
		_, err = io.Copy(target, file)
		if err != nil {
			slog.Error(err.Error())
			MakeSpecialPagePartWithRedirect(writer, &MessageUpdate{IsError: true,
				Message: "error while trying to copy the file onto the server"})
			return
		}
	}

	MakeSpecialPagePartWithRedirect(writer, &MessageUpdate{IsError: false,
		Message: "file successfully uploaded"})
}

func PostDirectSQLQuery(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !(loggedIn && acc.Role == database.RootAdmin) {
		PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		MakeSpecialPagePart(writer, &AdminSQLQuery{Query: &database.AdministrationQuery{
			Rows:  nil,
			Error: errors.New("server error: could not parse values"),
		}})
		return
	}
	if os.Getenv("SQL_KEY") != values.GetString("key") {
		MakeSpecialPagePart(writer, &AdminSQLQuery{Query: &database.AdministrationQuery{
			Rows:  nil,
			Error: errors.New("user error: provided key for accessing direct SQL statements is not valid"),
		}})
		return
	}

	MakeSpecialPagePart(writer, &AdminSQLQuery{Query: database.ExecuteQueryString(values.GetString("query"))})
}
