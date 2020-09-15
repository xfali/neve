// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package inject

import (
	"github.com/xfali/neve/container"
	"github.com/xfali/neve/injector"
	"github.com/xfali/neve/log"
	"testing"
)

type a interface {
	Get() int
}

type aImpl struct {
}

func (a *aImpl) Get() int {
	return 1
}

type bImpl struct {
}

func (a *bImpl) Get() int {
	return 2
}

type dest struct {
	A a `inject:""`
	B a `inject:"b"`
}

func TestInjectInterface(t *testing.T) {
	t.Run("inject once", func(t *testing.T) {
		c := container.New()
		c.Register(&aImpl{})
		c.RegisterByName("b", &bImpl{})
		i := injector.New(log.GetLogger())

		d := dest{}
		err := i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		if d.A == nil || d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B == nil || d.B.Get() != 2 {
			t.Fatal("inject B failed")
		}
	})

	t.Run("inject twice", func(t *testing.T) {
		c := container.New()
		c.Register(&aImpl{})
		c.RegisterByName("b", &bImpl{})
		i := injector.New(log.GetLogger())

		d := dest{}
		err := i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		if d.A == nil || d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B == nil || d.B.Get() != 2 {
			t.Fatal("inject B failed")
		}

		err = i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		if d.A == nil || d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B == nil || d.B.Get() != 2 {
			t.Fatal("inject B failed")
		}
	})
}

type dest2 struct {
	A aImpl `inject:""`
	B bImpl `inject:"b"`
}

func TestInjectStruct(t *testing.T) {
	t.Run("inject once", func(t *testing.T) {
		c := container.New()
		c.Register(&aImpl{})
		c.RegisterByName("b", &bImpl{})
		i := injector.New(log.GetLogger())

		d := dest2{}
		err := i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		if d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B.Get() != 2 {
			t.Fatal("inject B failed")
		}
	})

	t.Run("inject twice", func(t *testing.T) {
		c := container.New()
		c.Register(&aImpl{})
		c.RegisterByName("b", &bImpl{})
		i := injector.New(log.GetLogger())

		d := dest2{}
		err := i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		if d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B.Get() != 2 {
			t.Fatal("inject B failed")
		}

		err = i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		if d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B.Get() != 2 {
			t.Fatal("inject B failed")
		}
	})
}
