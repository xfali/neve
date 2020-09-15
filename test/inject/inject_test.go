// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package inject

import (
	"github.com/xfali/neve/container"
	"github.com/xfali/neve/injector"
	"io"
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
	i int
}

func (a *bImpl) Get() int {
	if a.i != 0 {
		return a.i
	}
	return 2
}

type dest struct {
	A a `inject:""`
	B a `inject:"b"`
	// Would not inject
	C io.Writer `inject:""`
}

func TestInjectInterface(t *testing.T) {
	t.Run("inject once", func(t *testing.T) {
		c := container.New()
		c.Register(&aImpl{})
		c.RegisterByName("b", &bImpl{})
		i := injector.New()

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
		i := injector.New()

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

	t.Run("inject twice with modify", func(t *testing.T) {
		c := container.New()
		c.Register(&aImpl{})
		b := &bImpl{}
		c.RegisterByName("b", b)
		i := injector.New()

		d := dest{}
		err := i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		//modify here
		b.i = 3
		if d.A == nil || d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B == nil || d.B.Get() != 3 {
			t.Fatal("inject B failed")
		}

		err = i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		b.i = 2
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
	B *bImpl `inject:"b"`
	B2 bImpl `inject:"b"`
	// Would not inject
	C dest `inject:""`
}

func TestInjectStruct(t *testing.T) {
	t.Run("inject once", func(t *testing.T) {
		c := container.New()
		c.Register(&aImpl{})
		c.RegisterByName("b", &bImpl{})
		i := injector.New()

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
		i := injector.New()

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

	t.Run("inject twice with modify", func(t *testing.T) {
		c := container.New()

		c.Register(&aImpl{})
		b := &bImpl{}
		c.RegisterByName("b", b)
		i := injector.New()

		d := dest2{}
		err := i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		b.i = 3
		if d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B.Get() != 3 {
			t.Fatal("inject B failed")
		}
		if d.B2.Get() != 2 {
			t.Fatal("inject B2 failed")
		}

		err = i.Inject(c, &d)
		if err != nil {
			t.Fatal(err)
		}

		b.i = 2
		if d.A.Get() != 1 {
			t.Fatal("inject A failed")
		}
		if d.B.Get() != 2 {
			t.Fatal("inject B failed")
		}
		if d.B2.Get() != 3 {
			t.Fatal("inject B2 failed")
		}
	})
}
