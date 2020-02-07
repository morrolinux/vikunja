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

package migration

import (
	"github.com/go-xorm/xorm"
	"src.techknowlogick.com/xormigrate"
)

type status20200120201756 struct {
	ID           int64  `xorm:"int(11) autoincr not null unique pk" json:"id"`
	UserID       int64  `xorm:"int(11) not null" json:"user_id"`
	MigratorName string `xorm:"varchar(255)" json:"migrator_name"`
	CreatedUnix  int64  `xorm:"created not null"`
}

func (s status20200120201756) TableName() string {
	return "migration_status"
}

func init() {
	migrations = append(migrations, &xormigrate.Migration{
		ID:          "20200120201756",
		Description: "Add migration status table",
		Migrate: func(tx *xorm.Engine) error {
			return tx.Sync2(status20200120201756{})
		},
		Rollback: func(tx *xorm.Engine) error {
			return dropTableColum(tx, "migration_status", "index")
		},
	})

}
