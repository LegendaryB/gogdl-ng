package task

import "github.com/LegendaryB/gogdl-ng/models"

type Store interface {
	GetAll() ([]models.Task, error)
	Get(id int64) (models.Task, error)
	Create(task models.Task) (*models.Task, error)
	Update(task models.Task) error
	Delete(id int64) error
}
