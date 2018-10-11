package crud

import (
	"code.vikunja.io/api/models"
	"github.com/labstack/echo"
	"net/http"
)

// UpdateWeb is the webhandler to update an object
func (c *WebHandler) UpdateWeb(ctx echo.Context) error {

	// Get our model
	currentStruct := c.EmptyStruct()

	// Get the object & bind params to struct
	if err := ParamBinder(currentStruct, ctx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "No or invalid model provided.")
	}

	// Check if the user has the right to do that
	currentUser, err := models.GetCurrentUser(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not determine the current user.")
	}
	if !currentStruct.CanUpdate(&currentUser) {
		models.Log.Noticef("%s [ID: %d] tried to update while not having the rights for it", currentUser.Username, currentUser.ID)
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// Do the update
	err = currentStruct.Update()
	if err != nil {
		return HandleHTTPError(err)
	}

	return ctx.JSON(http.StatusOK, currentStruct)
}
