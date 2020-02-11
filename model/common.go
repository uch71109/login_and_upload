package model

import (
	"github.com/jinzhu/gorm"
)

// New creates a DataOp
func New(db *gorm.DB) DataOp {
	return DataOp{
		User: newUserDao(db),
	}
}

// Init data op
func (o DataOp) Init() error {
	daos := []dao{
		o.User.(dao),
	}

	for _, d := range daos {
		if err := d.schema(); err != nil {
			return err
		}
		if err := d.index(); err != nil {
			return err
		}
	}
	return nil
}
