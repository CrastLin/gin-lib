package upload
/**
 @auth CrastGin
 @date 2020-10
 */
type Ftp struct {
	Driver
}

// init ftp driver
func InitFtp(opt *Option) Uploader {
	ftp := &Ftp{Driver{config: opt}}
	return ftp
}

// check root path
func (*Ftp) CheckPath(path ...string) bool {
	return false
}


// make path directory
func (*Ftp) MakeDir(path string) bool {
	return false
}

// save file to path
func (*Ftp) Save(file interface{}, replace bool) bool {
	return false
}
