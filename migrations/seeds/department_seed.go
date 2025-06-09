// Department and user seeder with transactional creation
package seeds

import (
	"encoding/json"
	"io"
	"os"

	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

type deptSeed struct {
	Name     string `json:"name"`
	Faculty  string `json:"faculty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DepartmentAndUserSeeder(db *gorm.DB) error {
	f, err := os.Open("./migrations/json/departments.json")
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	var seeds []deptSeed
	if err := json.Unmarshal(data, &seeds); err != nil {
		return err
	}
	for _, s := range seeds {
		tx := db.Begin()
		user := entity.User{Name: s.Name, Email: s.Email, Password: s.Password, Role: entity.RoleDepartemen}
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
		dept := entity.Department{Name: s.Name, Faculty: s.Faculty, UserID: user.ID}
		if err := tx.Create(&dept).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}
