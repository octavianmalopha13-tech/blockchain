package main

import (
    "crypto/sha256"
    "fmt"
    "strconv"
    "strings"
    "time"
)

// Block represents a single block in the chain
type Block struct {
    Index     int
    Timestamp int64
    Data      string
    PrevHash  string
    Hash      string
}

// Blockchain holds the chain of blocks
type Blockchain struct {
    Chain []Block
}

// calculateHash computes the SHA256 hash of a block's data
func calculateHash(b Block) string {
    // Concatenate all fields
    data := strconv.Itoa(b.Index) + strconv.FormatInt(b.Timestamp, 10) + b.Data + b.PrevHash
    hash := sha256.Sum256([]byte(data))
    return fmt.Sprintf("%x", hash[:])
}

// isValidHash checks if the hash starts with "00" (dummy proof-of-work)
func isValidHash(hash string) bool {
    return strings.HasPrefix(hash, "00")
}

// NewBlock creates a new block with a valid hash (loop until found)
func NewBlock(index int, data string, prevHash string) Block {
    newBlock := Block{
        Index:     index,
        Timestamp: time.Now().Unix(),
        Data:      data,
        PrevHash:  prevHash,
        Hash:      "",
    }

    // Keep re-calculating hash until it starts with "00"
    for {
        newBlock.Hash = calculateHash(newBlock)
        if isValidHash(newBlock.Hash) {
            break
        }
        // If hash doesn't start with "00", we increment timestamp slightly
        // to get a different hash (simulates mining)
        newBlock.Timestamp++
    }

    return newBlock
}

// AddBlock adds a new block to the blockchain (with validation)
func (bc *Blockchain) AddBlock(data string) error {
    // Get the last block
    last := bc.Chain[len(bc.Chain)-1]

    // Create new block
    newBlock := NewBlock(last.Index+1, data, last.Hash)

    // Validate the new block
    if newBlock.PrevHash != last.Hash {
        return fmt.Errorf("invalid block: previous hash mismatch")
    }
    if !isValidHash(newBlock.Hash) {
        return fmt.Errorf("invalid block: hash does not start with 00")
    }

    // Append to chain
    bc.Chain = append(bc.Chain, newBlock)
    return nil
}

// ValidateChain verifies the entire blockchain's integrity
func (bc *Blockchain) ValidateChain() bool {
    chain := bc.Chain

    // Validate genesis block
    if len(chain) > 0 {
        if !isValidHash(chain[0].Hash) {
            return false
        }
    }

    // Validate each block (except genesis)
    for i := 1; i < len(chain); i++ {
        current := chain[i]
        prev := chain[i-1]

        // Recalculate hash and compare
        if calculateHash(current) != current.Hash {
            return false
        }

        // Check PrevHash linkage
        if current.PrevHash != prev.Hash {
            return false
        }

        // Check proof-of-work
        if !isValidHash(current.Hash) {
            return false
        }
    }

    return true
}

// PrintChain prints all blocks in the chain
func (bc *Blockchain) PrintChain() {
    for _, block := range bc.Chain {
        fmt.Printf("Index: %d, Data: %s, Hash: %s, Prev: %s\n",
            block.Index, block.Data, block.Hash, block.PrevHash)
    }
}

func main() {
    // Create genesis block
    genesis := Block{
        Index:     0,
        Timestamp: time.Now().Unix(),
        Data:      "Genesis",
        PrevHash:  "",
        Hash:      "",
    }

    // Mine genesis block (find a valid hash)
    for {
        genesis.Hash = calculateHash(genesis)
        if isValidHash(genesis.Hash) {
            break
        }
        genesis.Timestamp++
    }

    // Initialize blockchain
    bc := &Blockchain{
        Chain: []Block{genesis},
    }

    // Add 3 blocks
    bc.AddBlock("Block 1")
    bc.AddBlock("Block 2")
    bc.AddBlock("Block 3")

    // Print chain
    fmt.Println("=== Blockchain ===")
    bc.PrintChain()

    // Validate
    fmt.Printf("\nChain is valid: %v\n", bc.ValidateChain())

    // TAMPER with Block 1 (index 1)
    if len(bc.Chain) > 1 {
        fmt.Println("\n--- TAMPERING with Block 1 ---")
        bc.Chain[1].Data = "Tampered!"
    }

    // Validate again
    fmt.Printf("After tampering, chain is valid: %v\n", bc.ValidateChain())
}
