package finite_fields

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearEquation(t *testing.T) {
	// given the linear equation
	// x + 2y = 1
	// 7x + 3y = 2
	// for real numbers the solution is
	// x = 1/11 and y = 5/11
	// but in a finite field of modulo 11, there is no solution
	// re-writing the equation as
	// y = 1/2 - x/2 and x = 2/3 - 7y/3
	// we see that y is the multiplicative inverse of 2 minus the multiplicative inverse of 2 time x
	// and x is the multiplicative inverse of 3 times 2 minus the multiplicative inverse of 3 time 7y
	// solving for y we get 
	a := new(big.Int).SetUint64(0).ModInverse(big.NewInt(2), big.NewInt(11))
	assert.True(t, a.Cmp(big.NewInt(6)) == 0) // 2 * 6 = 1 mod 11
	// hence y = 6 - 6x because the mul_inv of 2 is 6
	// the slope of the line is 6 note that slope is the coefficient of x of an equation in the form of y = mx + b
	// and for x
	b := new(big.Int).SetUint64(0).ModInverse(big.NewInt(3), big.NewInt(11))
	assert.True(t, b.Cmp(big.NewInt(4)) == 0) // 3 * 4 = 1 mod 11
	// hence x = 4 * 2 - 4 * 7y
	c := new(big.Int).SetUint64(0).Mul(b, big.NewInt(7))
	assert.True(t, c.Mod(c, big.NewInt(11)).Cmp(big.NewInt(6)) == 0) // 4 * 7 = 6 mod 11
	// we see that the slope of the line is 6. Both equations have the same slope
}