package upload

type Binary struct {
	Driver
	config *Option
}



func (*Binary) InitBinary(opt *Option) Uploader {
	bin := &Binary{config: opt}
	return bin
}

// check root path
func (*Binary) CheckPath(path ...string) bool {
	return false
}


// make path directory
func (*Binary) MakeDir(path string) bool {
	return false
}

// save file to path
func (*Binary) Save(file interface{}, replace bool) bool {

	return false
}