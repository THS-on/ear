// Copyright 2022 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateCmd_unknown_argument(t *testing.T) {
	cmd := NewCreateCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CreateCmd_no_output_file(t *testing.T) {
	cmd := NewCreateCmd()

	cmd.SetArgs([]string{})

	err := cmd.Execute()
	assert.EqualError(t, err, "validating arguments: no output file supplied")
}

func Test_CreateCmd_skey_file_not_found(t *testing.T) {
	cmd := NewCreateCmd()

	files := []fileEntry{
		{"ear-claims.json", testMiniClaimsSet},
	}
	makeFS(t, files)

	args := []string{
		"--skey=non-existent-skey.json",
		"--claims=ear-claims.json",
		"--alg=ES256",
		"ear.jwt",
	}
	cmd.SetArgs(args)

	expectedErr := `loading signing key from "non-existent-skey.json": open non-existent-skey.json: file does not exist`

	err := cmd.Execute()
	assert.EqualError(t, err, expectedErr)
}

func Test_CreateCmd_skey_file_bad_format(t *testing.T) {
	cmd := NewCreateCmd()

	files := []fileEntry{
		{"ear-claims.json", testMiniClaimsSet},
		{"empty-skey.json", testEmptyKey},
	}
	makeFS(t, files)

	args := []string{
		"--skey=empty-skey.json",
		"--claims=ear-claims.json",
		"--alg=ES256",
		"ear.jwt",
	}
	cmd.SetArgs(args)

	expectedErr := `parsing signing key from "empty-skey.json": failed to unmarshal JSON into key hint: EOF`

	err := cmd.Execute()
	assert.EqualError(t, err, expectedErr)
}

func Test_CreateCmd_skey_not_ok_for_signing(t *testing.T) {
	cmd := NewCreateCmd()

	files := []fileEntry{
		{"ear-claims.json", testMiniClaimsSet},
		{"pkey.json", testPKey},
	}
	makeFS(t, files)

	args := []string{
		"--skey=pkey.json",
		"--claims=ear-claims.json",
		"--alg=ES256",
		"ear.jwt",
	}
	cmd.SetArgs(args)

	expectedErr := `failed to generate signature for signer #0 (alg=ES256): failed to sign payload: failed to retrieve ecdsa.PrivateKey out of *jwk.ecdsaPublicKey: failed to produce ecdsa.PrivateKey from *jwk.ecdsaPublicKey: argument to AssignIfCompatible() must be compatible with *ecdsa.PublicKey (was *ecdsa.PrivateKey)`

	err := cmd.Execute()
	assert.ErrorContains(t, err, expectedErr)
}

func Test_CreateCmd_input_file_not_found(t *testing.T) {
	cmd := NewCreateCmd()

	files := []fileEntry{}
	makeFS(t, files)

	args := []string{
		"--skey=ignored.json",
		"--claims=non-existent-input.json",
		"--alg=ES256",
		"ear.jwt",
	}
	cmd.SetArgs(args)

	expectedErr := `loading EAR claims-set from "non-existent-input.json": open non-existent-input.json: file does not exist`

	err := cmd.Execute()
	assert.EqualError(t, err, expectedErr)
}

func Test_CreateCmd_input_file_bad_format(t *testing.T) {
	cmd := NewCreateCmd()

	files := []fileEntry{
		{"ear-claims.json", testEmptyClaimsSet},
	}
	makeFS(t, files)

	args := []string{
		"--skey=ignored.json",
		"--claims=ear-claims.json",
		"--alg=ES256",
		"ear.jwt",
	}
	cmd.SetArgs(args)

	expectedErr := `decoding EAR claims-set from "ear-claims.json": missing mandatory 'ear.status', 'eat_profile', 'iat'`

	err := cmd.Execute()
	assert.EqualError(t, err, expectedErr)
}

func Test_CreateCmd_unknown_signing_alg(t *testing.T) {
	cmd := NewCreateCmd()

	files := []fileEntry{
		{"skey.json", testSKey},
		{"ear-claims.json", testMiniClaimsSet},
	}
	makeFS(t, files)

	args := []string{
		"--skey=skey.json",
		"--claims=ear-claims.json",
		"--alg=XYZ",
		"ear.jwt",
	}
	cmd.SetArgs(args)

	expectedErr := `expected algorithm to be of type jwa.SignatureAlgorithm but got ("XYZ", jwa.InvalidKeyAlgorithm)`

	err := cmd.Execute()
	assert.ErrorContains(t, err, expectedErr)
}

func Test_CreateCmd_ok(t *testing.T) {
	cmd := NewCreateCmd()

	files := []fileEntry{
		{"skey.json", testSKey},
		{"ear-claims.json", testMiniClaimsSet},
	}
	makeFS(t, files)

	args := []string{
		"--skey=skey.json",
		"--claims=ear-claims.json",
		"--alg=ES256",
		"ear.jwt",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("ear.jwt")
	assert.NoError(t, err)
}
