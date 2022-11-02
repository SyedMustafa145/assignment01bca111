package assignment01bca
import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

type MerkelNode struct {
	left        *MerkelNode
	right       *MerkelNode
	transaction string
}

type MerkelTree struct {
	Root *MerkelNode
}

func addnode(transac string, arr *[10]MerkelTree, merkelindex *int) {

	var ind int
	ind = *merkelindex
	var tree *MerkelNode = arr[ind].Root

	var newitem MerkelNode

	newitem.transaction = transac
	newitem.right = nil
	newitem.left = nil

	if tree == nil {
		fmt.Println("Inserting nil ", transac)
		arr[ind].Root = &newitem
		ind++

	} else {

		fmt.Println("Inserting ", transac)
		insert_item(tree, &newitem)
		ind++
	}

	*merkelindex = ind
}

func insert_item(tree *MerkelNode, item *MerkelNode) {

	if len(item.transaction) <= len(tree.transaction) {

		if tree.left == nil {
			tree.left = item
			return
		} else {
			insert_item(tree.left, item)
			return
		}
	} else if len(item.transaction) > len(tree.transaction) {

		if tree.right == nil {
			tree.right = item
			return
		} else {
			insert_item(tree.right, item)
			return
		}

	}
}

func display_tree(tree *MerkelNode) {

	if tree == nil {
		return
	}
	fmt.Println(tree.transaction)

	if tree.left != nil {
		display_tree(tree.left)
	}

	if tree.right != nil {
		display_tree(tree.right)
	}
}

func update2(tree *MerkelNode, prev string, now string) {

	if tree == nil {
		return
	}
	if prev == tree.transaction {
		tree.transaction = now
	}

	if tree.left != nil {
		update2(tree.left, prev, now)
	}

	if tree.right != nil {
		update2(tree.right, prev, now)
	}

}

func update(arr *[10]MerkelTree, merkelindex *int, prev string, now string) {

	var ind int
	ind = *merkelindex
	var i = 0
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		update2(tree, prev, now)

	}

	*merkelindex = ind

}

func traversal2(tree *MerkelNode, tt *string) {

	if tree == nil {
		return
	}
	*tt += tree.transaction

	if tree.left != nil {
		traversal2(tree.left, tt)
	}

	if tree.right != nil {
		traversal2(tree.right, tt)
	}

}

func traversal(arr *[10]MerkelTree, merkelindex *int) string {

	var ind int
	ind = *merkelindex
	var alltransactions string
	var i = 0
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		traversal2(tree, &alltransactions)

	}

	*merkelindex = ind

	return alltransactions

}

func displayMerkelTree(arr *[10]MerkelTree, merkelindex *int) {
	var ind int
	ind = *merkelindex
	var i = 0

	fmt.Println("------------------Merkel Tree elements-------------------------------------")
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		if tree == nil {
			fmt.Println("arr[", i, "] has no elements")

		} else {
			fmt.Println("arr[", i, "] has elements")
			display_tree(tree)
		}
	}

	*merkelindex = ind

}

func createMerkelTree(arr *[10]MerkelTree, merkelindex *int) {

	for i := 0; i < 10; i++ {
		arr[i].Root = nil
	}

	addnode("transaction1", arr, merkelindex)
	addnode("transaction2", arr, merkelindex)
	addnode("transaction3", arr, merkelindex)
	addnode("transaction4", arr, merkelindex)
	addnode("transaction5", arr, merkelindex)
	addnode("transaction6", arr, merkelindex)

}

type block struct {
	arr           [10]MerkelTree
	merkelindex   int
	id            int
	nonce         string
	previous_hash string
	current_hash  string
}

type blockchain struct {
	list []*block
}

func newBlock(x int) *block {
	//fmt.Println("------------------------fdsfdsfdsfds-----------------------------------")
	tempblock := new(block)
	tempblock.id = x
	tempblock.nonce = "0"
	tempblock.merkelindex = 0
	createMerkelTree(&tempblock.arr, &tempblock.merkelindex)
	return tempblock
}

func verifyChain(chain *blockchain) bool {
	var temp = ""
	var check = true
	for i := 0; i < len(chain.list); i++ {
		tt := traversal(&chain.list[i].arr, &chain.list[i].merkelindex)

		var attributes string
		attributes += strconv.Itoa(chain.list[i].id)
		attributes += tt + chain.list[i].previous_hash
		total_sum := sha256.Sum256([]byte(attributes))
		temp = fmt.Sprintf("%x", total_sum)

		if temp != chain.list[i].current_hash {
			check = false
			fmt.Printf("Previous block has been tampered, i.e. Block # %d\n", i)
			break

		}
	}

	if check == false {
		fmt.Println("error occured")
	} else {
		fmt.Printf("Blocks verified. No tampering\n")
	}
	return check
}

func Mineblock(blocklist *blockchain) {

	for j := 0; j < len(blocklist.list); j++ {
		print("to match:", blocklist.list[j].current_hash, "\n")
		for i := 0; ; i++ {
			temp := sha256.Sum256([]byte(strconv.Itoa(i)))
			noncex := fmt.Sprintf("%x", temp)
			dum := noncex[:3]
			fmt.Println("dum:", dum)
			fmt.Println(strings.Contains(blocklist.list[j].current_hash, dum))

			if strings.Contains(blocklist.list[j].current_hash, dum) == true {
				blocklist.list[j].nonce = dum
				break

			}

		}

	}

}

func CalculateHash(chain *blockchain) {

	for i := 0; i < len(chain.list); i++ {
		tt := traversal(&chain.list[i].arr, &chain.list[i].merkelindex)
		var attributes string
		attributes += strconv.Itoa(chain.list[i].id)
		attributes += tt + chain.list[i].previous_hash
		total_sum := sha256.Sum256([]byte(attributes))
		chain.list[i].current_hash = fmt.Sprintf("%x", total_sum) // formating to string
		if i < len(chain.list)-1 {
			chain.list[i+1].previous_hash = fmt.Sprintf("%x", total_sum) //storing current block hash to next block in its previous hash var
		}

	}
}

func (blocklist *blockchain) addblock(x int) *block {
	tempblock := newBlock(x)

	if verifyChain(blocklist) {
		blocklist.list = append(blocklist.list, tempblock)
		CalculateHash(blocklist)

		fmt.Printf("block addition in chain successful\n")
	} else {
		fmt.Printf(" error. block addition unsuccessful.\n")
		return nil
	}
	return tempblock
}

func DisplayBlocks(blocklist *blockchain) {
	fmt.Println("")

	for i := 0; i < len(blocklist.list); i++ {
		fmt.Printf("Block id:%d\n\n", blocklist.list[i].id)
		displayMerkelTree(&blocklist.list[i].arr, &blocklist.list[i].merkelindex)
		fmt.Println("nonce value : \n", blocklist.list[i].nonce)
		fmt.Println("current hash: \n", blocklist.list[i].current_hash)
		fmt.Println("previous hash: \n", blocklist.list[i].previous_hash)

	}

	fmt.Println("")

}

func changeBlock(chain *blockchain, x int) { // updating on basis of id value as identifier

	found := false
	for i := 0; i < len(chain.list); i++ {

		if x == chain.list[i].id {

			var now string
			var prev string
			fmt.Println("Enter transaction to change\n")

			fmt.Scanln(&prev)

			fmt.Println("Enter updated value\n")

			fmt.Scanln(&now)
			fmt.Println("updated successfully\n")
			update(&chain.list[i].arr, &chain.list[i].merkelindex, prev, now)
			found = true
		}
	}
	if found == false {
		fmt.Println("error. Couldnt update. block not found")
	}
	return
}

func Main() {

	chain := new(blockchain)
	var x = 50
	chain.addblock(x)
	chain.addblock(x + 10)
	chain.addblock(x + 20)
	chain.addblock(x + 30)
	chain.addblock(x + 40)
	chain.addblock(x + 50)

	Mineblock(chain)
	DisplayBlocks(chain)

	//fmt.Println("updating block with id value :", x+30)
	//changeBlock(chain, x+30)

	//DisplayBlocks(chain)

}
