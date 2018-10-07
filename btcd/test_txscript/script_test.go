package script_test

import (
	"github.com/btcsuite/btcd/chaincfg"
	"fmt"
	"encoding/hex"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

// This func exemplifies creating P2SH scripts.
func ExampleP2SH () {
	addrStr := "1A2keVVxUs6Ww7QSBgorbXFW9ii96ycxf3"
	// generate Address type from address string
	addr, err := btcutil.DecodeAddress(addrStr, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a pukeyscript that pays to the address.
	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Script Hex: %x\n", script)
	

	disasm, err := txscript.DisasmString(script)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Disassembly:", disasm)
	// Output:
	// Script Hex: 76a914630dfb715583d8f307b9439ef8da42352e7a1a3788ac
	// Script Disassembly: OP_DUP OP_HASH160 630dfb715583d8f307b9439ef8da42352e7a1a37 OP_EQUALVERIFY OP_CHECKSIG
}

func ExampleExtractPkScriptAddrs() {
	scriptHex := "76a914630dfb715583d8f307b9439ef8da42352e7a1a3788ac"
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract and print details from the script.
	scriptClass, addrs, reqSigs, err := txscript.ExtractPkScriptAddrs(script, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Class:", scriptClass)
	fmt.Println("Addresses:", addrs)
	fmt.Println("Required Signatures:", reqSigs)
	// Output:
	// Script Class: pubkeyhash
	// Addresses: [1A2keVVxUs6Ww7QSBgorbXFW9ii96ycxf3]
	// Required Signatures: 1
}

func ExampleSignTxOutput() {
	privKeyBytes, err := hex.DecodeString("5820319847201928" + "d8c8a9a9a9a9a9d8")
	if err != nil {
		fmt.Println(err)
		return
	}

}
