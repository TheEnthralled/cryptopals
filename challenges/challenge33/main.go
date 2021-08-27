package main

import (
	//"cryptopals/utils/"
	"log"
	//"math"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"crypto/sha1"
)

const (
	EXPONENT_LEN = /*1024*/128 // Make sure multiple of 8
	P_HEX = `ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec6f44c42e9a637ed6b0bff5cb6f406b7edee386bfb5a899fa5ae9f24117c4b1fe649286651ece45b3dc2007cb8a163bf0598da48361c55d39a69163fa8fd24cf5f83655d23dca3ad961c62f356208552bb9ed529077096966d670c354e4abc9804f1746c08ca237327ffffffffffffffff`
	G_HEX = `2`
)

func DiffieHellmanKeyExchange() {
	p := new(big.Int)
	p.SetString(P_HEX, 16)
	fmt.Println("p = ", p)

	g := new(big.Int)
	g.SetString(G_HEX, 16)
	fmt.Println("g = " , g)

	_a := make([]byte, EXPONENT_LEN)
	_, err := rand.Read(_a)
	if err != nil {
		log.Fatal(err)
	}
	a := new(big.Int)
	a.SetString(hex.EncodeToString(_a), 16)
	fmt.Println("Alice's private key a (randomly generated) = ", a)
	fmt.Println()
	fmt.Println()

	
	A := new(big.Int)
	A.Exp(g, a, p)
	fmt.Println("Alice's public key A = (g**a)%p = ", A)
	fmt.Println()
	fmt.Println()
	
	_b := make([]byte, EXPONENT_LEN)
	_, err = rand.Read(_b)
	if err != nil {
		log.Fatal(err)
	}
	b := new(big.Int)
	b.SetString(hex.EncodeToString(_b), 16)
	fmt.Println("Bob's private key b (randomly generated) = ", b)
	fmt.Println()
	fmt.Println()
	
	B := new(big.Int)
	B.Exp(g, b, p)
	fmt.Println("Bob's public key B = (g**b)%p = ", B)
	fmt.Println()
	fmt.Println()

	s_A := new(big.Int)
	s_A.Exp(B, a, p)
	fmt.Println("Alice's secret key s_A = (B**a)%p = ", s_A)
	fmt.Println()
	fmt.Println()

	s_B := new(big.Int)
	s_B.Exp(A, b, p)
	fmt.Println("Bob's secret key s_B = (A**b)%p = " , s_B)
	fmt.Println()
	fmt.Println()


	_S := sha1.New()
	_, err = _S.Write([]byte(s_B.String()))
	if err != nil {
		log.Fatal(err)
	}
	S := _S.Sum(nil)
	fmt.Println("Final secret S = SHA1(s_B) = SHA1(s_A) = ", hex.EncodeToString(S))
	
}

func main() {
	DiffieHellmanKeyExchange()
}

