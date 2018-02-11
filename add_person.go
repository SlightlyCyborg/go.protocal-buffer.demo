package main

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func promptForAddress(r io.Reader) (*Person, error) {
	var p *Person
	var pn *Person_PhoneNumber
	var rd *bufio.Reader
	var name, email string
	var err error

	p = &Person{}
	rd = bufio.NewReader(r)

	fmt.Print("Enter person ID number: ")

	if _, err := fmt.Fscanf(rd, "%d\n", &p.Id); err != nil {
		return p, err
	}

	fmt.Print("Enter name: ")
	name, err = rd.ReadString('\n')
	if err != nil {
		return p, err
	}

	p.Name = strings.TrimSpace(name)

	fmt.Print("Enter an email (leave blank for none): ")
	email, err = rd.ReadString('\n')
	if err != nil {
		return p, err
	}

	p.Email = strings.TrimSpace(email)

	for {
		fmt.Print("Enter a phone number (or leave blank to finish): ")
		phone, err := rd.ReadString('\n')
		if err != nil {
			return p, err
		}
		phone = strings.TrimSpace(phone)
		if phone == "" {
			break
		}

		pn = &Person_PhoneNumber{
			Number: phone,
		}

		fmt.Print("Is this a mobile, home, or work phone? ")
		ptype, err := rd.ReadString('\n')
		if err != nil {
			return p, err
		}
		ptype = strings.TrimSpace(ptype)

		switch ptype {
		case "mobile":
			pn.Type = Person_MOBILE
		case "home":
			pn.Type = Person_HOME
		case "work":
			pn.Type = Person_WORK
		default:
			fmt.Printf("Unknown phone type %q.  Using default.\n", ptype)
		}

		p.Phones = append(p.Phones, pn)
	}

	return p, err

}

// Main reads the entire address book from a file, adds one person based on
// user input, then writes it back out to the same file.
func main() {
	var book *AddressBook

	if len(os.Args) != 2 {
		log.Fatalf("Usage:  %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}
	fname := os.Args[1]

	// Read the existing address book.
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File not found.  Creating new file.\n", fname)
		} else {
			log.Fatalln("Error reading file:", err)
		}
	}

	// [START marshal_proto]
	book = &AddressBook{}
	// [START_EXCLUDE]
	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	// Add an address.
	addr, err := promptForAddress(os.Stdin)
	if err != nil {
		log.Fatalln("Error with address:", err)
	}
	book.People = append(book.People, addr)
	// [END_EXCLUDE]

	// Write the new address book back to disk.
	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
	// [END marshal_proto]
}
