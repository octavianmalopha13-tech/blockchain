package main

import (
    "crypto/sha256"
    "fmt"
    "strings"
)

type Block struct {
    Index    int
    Data     string
    PrevHash string
    Nonce    int      // BUG 3: was string — the nonce is a NUMBER, and calculateHash expects int
    Hash     string
}

func calculateHash(index int, data, prevHash string, nonce int) string { // BUG 1: missing return type
    input := fmt.Sprintf("%d%s%s%d", index, data, prevHash, nonce)
    sum := sha256.Sum256([]byte(input))
    return fmt.Sprintf("%x", sum)
}

func mineBlock(index int, data, prevHash string, difficulty int) Block {
    prefix := strings.Repeat("0", difficulty)
    for nonce := 0; ; nonce++ {   // BUG 2: was `o` (letter) — must be `0` (zero)
        hash := calculateHash(index, data, prevHash, nonce)
        if strings.HasPrefix(hash, prefix) {
            return Block{Index: index, Data: data, PrevHash: prevHash, Nonce: nonce, Hash: hash}
        }
    }
}

func isValidChain(chain []Block, difficulty int) bool {
    prefix := strings.Repeat("0", difficulty)
    for i := 1; i < len(chain); i++ {   // BUG 4: was i:=0 — at i=0, chain[i-1] = chain[-1] = panic!
        if chain[i].PrevHash != chain[i-1].Hash {
            fmt.Println("INVALID: block", i, "does not point to block", i-1)
            return false
        }
        if chain[i].Hash != calculateHash(chain[i].Index, chain[i].Data, chain[i].PrevHash, chain[i].Nonce) {
            fmt.Println("INVALID: block", i, "was tampered with")
            return false
        }
        if !strings.HasPrefix(chain[i].Hash, prefix) {   // BUG 5: was `Prefix` — Go is case-sensitive
            fmt.Println("INVALID: block", i, "did not do the work")
            return false
        }
    }
    fmt.Println("chain is VALID")
    return true
}

func main() {
    difficulty := 4

    genesis := Block{Index: 0, Data: "genesis block", PrevHash: ""}
    genesis.Hash = calculateHash(genesis.Index, genesis.Data, genesis.PrevHash, 0)

    b1 := mineBlock(1, "tx: Alice->Bob 10", genesis.Hash, difficulty)
    b2 := mineBlock(2, "tx: Bob->Carol 5", b1.Hash, difficulty)
    b3 := mineBlock(3, "tx: Carol->Emma 2", b2.Hash, difficulty)   // BUG 6: was b2:= again — redeclaration error + b3 never existed

    chain := []Block{genesis, b1, b2, b3}
isValidChain(chain, difficulty)         // VALID first

    chain[1].Data = "tx: Mallory->Bob 10"   // tamper AFTER mining — the data changes but the hash stays
    isValidChain(chain, difficulty)         // INVALID: block 1 was tampered with

    for _, b := range chain {
        fmt.Printf("Block %d: %s (nonce %d)\n  hash: %s\n  prev: %s\n", b.Index, b.Data, b.Nonce, b.Hash, b.PrevHash)   // BUG 7: 4 verbs, 5 args
    }

    isValidChain(chain, difficulty)
}
