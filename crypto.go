package crypto

import (
	"crypto/rand"
	"math/big"
)

var (
	gx   *big.Int
	gy   *big.Int
	p    *big.Int
	n    *big.Int
	big7 *big.Int
	big3 *big.Int
	big2 *big.Int
)

type PrivateKey struct {
	X      *big.Int
	Y      *big.Int
	text32 string
}

func Init() bool {
	var ok bool
	big7, ok = new(big.Int).SetString("7", 10)
	if ok != true {
		goto out
	}
	big2, ok = new(big.Int).SetString("2", 10)
	if ok != true {
		goto out
	}
	big3, ok = new(big.Int).SetString("3", 10)
	if ok != true {
		goto out
	}
	n, ok = new(big.Int).SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)
	if ok != true {
		goto out
	}

	gx, ok = new(big.Int).SetString("55066263022277343669578718895168534326250603453777594175500187360389116729240", 10)
	if ok != true {
		goto out
	}
	gy, ok = new(big.Int).SetString("32670510020758816978083085130507043184471273380659243275938904335757337482424", 10)
	if ok != true {
		goto out
	}
	p, ok = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10)
	if ok != true {
		goto out
	}
	return true

out:
	panic("big.Int.SetString Error")
	return false
}
func Doubling(x *big.Int, y *big.Int) (*big.Int, *big.Int) {
	// s = (3x^2) / (2y)
	xSqaure := new(big.Int).Mul(x, x)
	xSqaure = new(big.Int).Mod(xSqaure, p)
	xSqaure3 := new(big.Int).Mul(big3, xSqaure)
	xSqaure3 = new(big.Int).Mod(xSqaure3, p)
	y2 := new(big.Int).Mul(y, big2)
	y2 = new(big.Int).Mod(y2, p)
	pminus2 := new(big.Int).Sub(p, big2)
	y2 = new(big.Int).Exp(y2, pminus2, p) // 2y^-1 = 1/2y
	s := new(big.Int).Mul(xSqaure3, y2)
	s = new(big.Int).Mod(s, p)

	// newx = s^2 - 2x
	sSqaure := new(big.Int).Mul(s, s)
	sSqaure = new(big.Int).Mod(sSqaure, p)
	x2 := new(big.Int).Mul(x, big2)
	x2 = new(big.Int).Mod(x2, p)
	newx := new(big.Int).Sub(sSqaure, x2)
	newx = new(big.Int).Mod(newx, p)

	// newy = s(x1-x3)-y
	x1minusx3 := new(big.Int).Sub(x, newx)
	x1minusx3 = new(big.Int).Mod(x1minusx3, p)
	sx1x3 := new(big.Int).Mul(s, x1minusx3)
	sx1x3 = new(big.Int).Mod(sx1x3, p)
	newy := new(big.Int).Sub(sx1x3, y)
	newy = new(big.Int).Mod(newy, p)
	return newx, newy
}

func Add(x0, y0, x1, y1 *big.Int) (*big.Int, *big.Int) {
	if x0.Cmp(x1) == 0 && y0.Cmp(y1) == 0 {
		return Doubling(x0, y0)
	}
	if x0.Cmp(x1) == 0 && y0.Cmp(y1) != 0 {
		return nil, nil
	}
	// s = (y1-y0)/(x1-x0)
	y1minusy0 := new(big.Int).Sub(y1, y0)
	y1minusy0 = new(big.Int).Mod(y1minusy0, p)
	x1minusx0 := new(big.Int).Sub(x1, x0)
	x1minusx0 = new(big.Int).Mod(x1minusx0, p)
	pminus2 := new(big.Int).Sub(p, big2)
	x1minusx0 = new(big.Int).Exp(x1minusx0, pminus2, p)
	s := new(big.Int).Mul(y1minusy0, x1minusx0)
	s = new(big.Int).Mod(s, p)

	// newx = s^2-x0-x1
	sSquare := new(big.Int).Mul(s, s)
	sSquare = new(big.Int).Mod(sSquare, p)
	sSquarex0 := new(big.Int).Sub(sSquare, x0)
	sSquarex0 = new(big.Int).Mod(sSquarex0, p)
	newx := new(big.Int).Sub(sSquarex0, x1)
	newx = new(big.Int).Mod(newx, p)

	// newy = s(x0-x2)-y0
	x0minusx2 := new(big.Int).Sub(x0, newx)
	x0minusx2 = new(big.Int).Mod(x0minusx2, p)
	sx0x2 := new(big.Int).Mul(s, x0minusx2)
	sx0x2 = new(big.Int).Mod(sx0x2, p)
	newy := new(big.Int).Sub(sx0x2, y0)
	newy = new(big.Int).Mod(newy, p)
	return newx, newy
}

func Mul(x *big.Int, y *big.Int, N int) (*big.Int, *big.Int) {
	xinit, yinit := new(big.Int).Set(gx), new(big.Int).Set(gy)
	newx, newy := new(big.Int).Set(gx), new(big.Int).Set(gy)
	for i := 0; i < N-1; i++ {
		newx, newy = Add(newx, newy, xinit, yinit)
	}
	return newx, newy

}

func IsOnTheCurve(x *big.Int, y *big.Int) bool {
	if x == nil || y == nil {
		return false
	}
	// y**2 = x**3 + 7 mod p
	eccgxp := new(big.Int).Mul(x, x)
	eccgxp = new(big.Int).Mul(eccgxp, x)
	eccgxp = new(big.Int).Add(eccgxp, big7)
	eccgxp = new(big.Int).Mod(eccgxp, p)

	eccgyp := new(big.Int).Mul(y, y)
	eccgyp = new(big.Int).Mod(eccgyp, p)
	if eccgxp.Cmp(eccgyp) == 0 {
		return true
	} else {
		return false
	}
}

func Generate() (*PrivateKey, bool) {
	k := genRandomK()
	if k == nil {
		return nil, false
	}
	xinit, yinit := new(big.Int).Set(gx), new(big.Int).Set(gy)
	privx, privy := new(big.Int).Set(gx), new(big.Int).Set(gy)
	one := big.NewInt(1)
	for cnt := big.NewInt(1); cnt.Cmp(k) == -1; cnt = new(big.Int).Add(cnt, one) {
		privx, privy = Add(privx, privy, xinit, yinit)
	}
	return &PrivateKey{
		X:      privx,
		Y:      privy,
		text32: k.Text(32),
	}, true
}

func genRandomK() *big.Int {
	//Max random value, 2^16 % order_number - 1
	max := new(big.Int)
	max.Exp(big2, big.NewInt(16), n).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	k, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil
	}
	return k
}
