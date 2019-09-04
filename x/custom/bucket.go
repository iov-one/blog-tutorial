package custom

import (
	"github.com/iov-one/weave/migration"
	"github.com/iov-one/weave/orm"
)

type TimedStateBucket struct {
	orm.ModelBucket
}

func NewTimedStateBucket() *TimedStateBucket {
	b := orm.NewModelBucket("timedstate", &TimedState{})
	return &TimedStateBucket{
		ModelBucket: migration.NewModelBucket(packageName, b),
	}
}

type StateBucket struct {
	orm.ModelBucket
}

func NewStateBucket() *StateBucket {
	b := orm.NewModelBucket("state", &State{})
	return &StateBucket{
		ModelBucket: migration.NewModelBucket(packageName, b),
	}
}
