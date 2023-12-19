package database

import (
	"context"

	"github.com/qiniu/qmgo"
)

func Connect(db, coll string, ctx context.Context) (*qmgo.QmgoClient, error) {

	return qmgo.Open(ctx, &qmgo.Config{
		Uri:      "mongodb://10.6.1.6:27017",
		Database: db,
		Coll:     coll,
	})

}
