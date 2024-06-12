package commands

func Stat() ([]byte, error) {
	stat := Command{Name: "nncp-stat", Path: "/usr/bin/"}
	stat.load()

	dat, err := stat.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return dat, nil
}
