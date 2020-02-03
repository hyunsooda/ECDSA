package crypto

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestInit(t *testing.T) {
	ok := Init()
	assert.Equal(t, ok, true)
}

func TestIsOnTheCurve(t *testing.T) {
	ok := Init()
	assert.Equal(t, ok, true)

	x := gx
	y := gy

	newx, newy := Mul(x, y, 5)
	ok = IsOnTheCurve(newx, newy)
	assert.Equal(t, ok, true)

	newx, newy = Mul(x, y, 57629)
	ok = IsOnTheCurve(newx, newy)
	assert.Equal(t, ok, true)

	newx, newy = Add(newx, newy, x, y)
	ok = IsOnTheCurve(newx, newy)
	assert.Equal(t, ok, true)

	privkey, err := Generate()
	assert.Equal(t, err, true)
	ok = IsOnTheCurve(privkey.X, privkey.Y)
	assert.Equal(t, ok, true)
}
