package prt

func BuildString(fn func(*Printer)) string {
	p := NewStringPrinter()
	fn(p)
	return p.String()
}

func BuildFile(filename string, fn func(*Printer)) error {
	p, err := NewFilePrinter(filename)
	if err != nil {
		return err
	}
	defer p.Close()
	
	fn(p)
	
	return p.Err()
}

