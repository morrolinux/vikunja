// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-2020 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package v1

import (
	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/models"
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/web/handler"
	"github.com/labstack/echo/v4"
	"net/http"
)

// RegisterUser is the register handler
// @Summary Register
// @Description Creates a new user account.
// @tags user
// @Accept json
// @Produce json
// @Param credentials body user.APIUserPassword true "The user credentials"
// @Success 200 {object} user.User
// @Failure 400 {object} code.vikunja.io/web.HTTPError "No or invalid user register object provided / User already exists."
// @Failure 500 {object} models.Message "Internal error"
// @Router /register [post]
func RegisterUser(c echo.Context) error {
	if !config.ServiceEnableRegistration.GetBool() {
		return echo.ErrNotFound
	}
	// Check for Request Content
	var datUser *user.APIUserPassword
	if err := c.Bind(&datUser); err != nil {
		return c.JSON(http.StatusBadRequest, models.Message{"No or invalid user model provided."})
	}

	// Insert the user
	newUser, err := user.CreateUser(datUser.APIFormat())
	if err != nil {
		return handler.HandleHTTPError(err, c)
	}

	// Add its namespace
	newN := &models.Namespace{Name: newUser.Username, Description: newUser.Username + "'s namespace.", Owner: newUser}
	err = newN.Create(newUser)
	if err != nil {
		return handler.HandleHTTPError(err, c)
	}

	return c.JSON(http.StatusOK, newUser)
}
