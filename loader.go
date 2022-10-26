package env

import "os"

func loadFiles(files ...string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, v := range files {
		fileVars, err := loadFile(v)
		if err != nil {
			return nil, err
		}
		for k, v := range fileVars {
			vars[k] = v
		}
	}

	return vars, nil
}

func loadFile(file string) (map[string]string, error) {
	by, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	vars, err := parseFile(by)
	if err != nil {
		return nil, err
	}
	return vars, nil
}
