package main

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

var Zero = big.NewInt(0)
var One = big.NewInt(1)
var Two = big.NewInt(2)

// Finds numbers y and r for a given number x such that x = 2^y * r and r is odd
func TwoFactorize(x *big.Int) (*big.Int, *big.Int) {
	y := new(big.Int).Set(Zero)
	r := new(big.Int).Set(x)

	// Divide by 2 until r becomes odd
	for new(big.Int).Mod(r, Two).Sign() == 0 {
		r.Div(r, Two)
		y.Add(y, One)
	}

	return y, r
}

func Crack(N, E, D *big.Int) (*big.Int, *big.Int) {
	MinusOne := new(big.Int).Sub(N, One)
	P := big.NewInt(0)
	Q := new(big.Int)
	A := new(big.Int)
	B := new(big.Int)
	K := new(big.Int).Mul(E, D)
	K.Sub(K, One)
	_, R := TwoFactorize(K)

	for P.Mod(P, N).Cmp(Zero) == 0 || P.Cmp(One) == 0 {
		X, _ := rand.Int(rand.Reader, N)
		X.Exp(X, R, N)
		if X.Cmp(One) == 0 {
			continue
		}
		v0 := new(big.Int)
		for X.Mod(X, N).Cmp(One) != 1 {
			v0.Set(X)
			X.Exp(X, Two, N)
		}
		if v0.Cmp(MinusOne) == 0 {
			continue
		}
		P.GCD(A, B, X.Add(X, One), N)

	}

	return P, Q.Div(N, P)

}

func D2PrivKey(D *big.Int, pubK *rsa.PublicKey) (*rsa.PrivateKey, error) {
	P, Q := Crack(pubK.N, big.NewInt(int64(pubK.E)), D)
	prk1 := &rsa.PrivateKey{}
	prk1.N = pubK.N
	prk1.D = D
	prk1.PublicKey = *pubK // The assembled private key will always match public key, so doesnt makes sense to verify
	prk1.Primes = []*big.Int{P, Q}
	prk1.Precompute()
	return prk1, nil
}
