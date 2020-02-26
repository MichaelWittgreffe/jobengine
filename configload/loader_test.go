package configload

import (
	"fmt"
	"testing"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
	"github.com/MichaelWittgreffe/jobengine/models"
)

//TestNewConfigLoader1 correctly constructs the object
func TestNewConfigLoader1(t *testing.T) {
	if result := NewConfigLoader("./test_file_1.yml", "os"); result == nil {
		t.Error("TestNewConfigLoader1: Returned Nil")
	}
}

//TestNewConfigLoader2 feeds invalid args to the constructor
func TestNewConfigLoader2(t *testing.T) {
	if result := NewConfigLoader("", "os"); result != nil {
		t.Error("NewConfigLoader Returned Valid Response Unexpectedly")
	}
}

//BenchmarkNewConfigLoader gets the average speed of constructing the struct
func BenchmarkNewConfigLoader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewConfigLoader("./test_file_1.yml", "os")
	}
}

//TestLoadFromFile1 tests loading a valid YAML file
func TestLoadFromFile1(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    []byte{118, 69, 45},
			ReadFileErrorResult:   nil,
			GetEnvResult:          "testsecret",
		},
		cfgParser: &MockYAMLHandler{
			UnmarshalResult: nil,
			UnmarshalOutput: map[interface{}]interface{}{
				"version":             1,
				"port":                6010,
				"job_keep_minutes":    30,
				"job_timeout_minutes": 30,
				"queues": map[interface{}]interface{}{
					"test_queue_1": map[interface{}]interface{}{
						"read":  []interface{}{"service1", "service2"},
						"write": []interface{}{"service3"},
					},
					"test_queue_2": map[interface{}]interface{}{
						"read":  []interface{}{"service1"},
						"write": []interface{}{"service2"},
					},
				},
			},
		},
		hasher: NewHashProcessor("md5"),
	}

	result, err := loader.LoadFromFile(1.0)
	if err != nil {
		t.Errorf("TestLoadFromFile1 Returned Error: %s", err.Error())
	}
	if result.Port != 6010 {
		t.Errorf("TestLoadFromFile1 Port Expected %d, Got %d", 6010, result.Port)
	}
	if result.Version != 1.0 {
		t.Errorf("TestLoadFromFile1 Version Expected %f Got %f", 1.0, result.Version)
	}
	if len(result.Queues) != 2 {
		t.Errorf("TestLoadFromFile1 Queues Expected Length %d, Got %d", 2, len(result.Queues))
	}
}

//TestLoadFromFile2 tests loading a valid YAML file, uses float for version
func TestLoadFromFile2(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    []byte{118, 69, 45},
			ReadFileErrorResult:   nil,
			GetEnvResult:          "testsecret",
		},
		cfgParser: &MockYAMLHandler{
			UnmarshalResult: nil,
			UnmarshalOutput: map[interface{}]interface{}{
				"version":             1.0,
				"port":                6010,
				"job_keep_minutes":    30,
				"job_timeout_minutes": 30,
				"queues": map[interface{}]interface{}{
					"test_queue_1": map[interface{}]interface{}{
						"read":  []interface{}{"service1", "service2"},
						"write": []interface{}{"service3"},
					},
					"test_queue_2": map[interface{}]interface{}{
						"read":  []interface{}{"service1"},
						"write": []interface{}{"service2"},
					},
				},
			},
		},
		hasher: NewHashProcessor("md5"),
	}

	result, err := loader.LoadFromFile(1.0)
	if err != nil {
		t.Errorf("TestLoadFromFile2 Returned Error: %s", err.Error())
	}
	if result.Port != 6010 {
		t.Errorf("TestLoadFromFile2 Port Expected %d, Got %d", 6010, result.Port)
	}
	if result.Version != 1.0 {
		t.Errorf("TestLoadFromFile2 Version Expected %f Got %f", 1.0, result.Version)
	}
	if len(result.Queues) != 2 {
		t.Errorf("TestLoadFromFile2 Queues Expected Length %d, Got %d", 2, len(result.Queues))
	}
}

//TestLoadFromFile3 has invalid []byte sent to yaml unmarshal func
func TestLoadFromFile3(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    nil,
			ReadFileErrorResult:   nil,
			GetEnvResult:          "testsecret",
		},
		cfgParser: &MockYAMLHandler{
			UnmarshalResult: fmt.Errorf("Test Error"),
		},
		hasher: NewHashProcessor("md5"),
	}
	_, err := loader.LoadFromFile(1.0)
	if err == nil {
		t.Error("TestLoadFromFile3 Did Not Return Error")
	}
}

//TestLoadFromFile4 has error reading the config file from disk
func TestLoadFromFile4(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    nil,
			ReadFileErrorResult:   fmt.Errorf("Test Error"),
			GetEnvResult:          "testsecret",
		},
		hasher: NewHashProcessor("md5"),
	}

	_, err := loader.LoadFromFile(1.0)
	if err == nil {
		t.Error("TestLoadFromFile4 Did Not Return Error")
	}
}

//TestLoadFromFile5 has error checking whether the config file exists
func TestLoadFromFile5(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  false,
			FileExistsErrorResult: fmt.Errorf("Test Error"),
			GetEnvResult:          "testsecret",
		},
		hasher: NewHashProcessor("md5"),
	}

	_, err := loader.LoadFromFile(1.0)
	if err == nil {
		t.Error("TestLoadFromFile5 Did Not Return Error")
	}
}

//TestLoadFromFile6 has invalid version type
func TestLoadFromFile6(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    []byte{118, 69, 45},
			ReadFileErrorResult:   nil,
			GetEnvResult:          "testsecret",
		},
		cfgParser: &MockYAMLHandler{
			UnmarshalResult: nil,
			UnmarshalOutput: map[interface{}]interface{}{
				"version":             "1.0",
				"port":                6010,
				"job_keep_minutes":    30,
				"job_timeout_minutes": 30,
				"queues": map[interface{}]interface{}{
					"test_queue_1": map[interface{}]interface{}{
						"read":  []interface{}{"service1", "service2"},
						"write": []interface{}{"service3"},
					},
					"test_queue_2": map[interface{}]interface{}{
						"read":  []interface{}{"service1"},
						"write": []interface{}{"service2"},
					},
				},
			},
		},
		hasher: NewHashProcessor("md5"),
	}

	_, err := loader.LoadFromFile(1.0)
	if err == nil {
		t.Error("TestLoadFromFile6 Did Not Return Error")
	}
}

//TestLoadFromFile7 has invalid version value
func TestLoadFromFile7(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    []byte{118, 69, 45},
			ReadFileErrorResult:   nil,
			GetEnvResult:          "testsecret",
		},
		cfgParser: &MockYAMLHandler{
			UnmarshalResult: nil,
			UnmarshalOutput: map[interface{}]interface{}{
				"version":             2.0,
				"port":                6010,
				"job_keep_minutes":    30,
				"job_timeout_minutes": 30,
				"queues": map[interface{}]interface{}{
					"test_queue_1": map[interface{}]interface{}{
						"read":  []interface{}{"service1", "service2"},
						"write": []interface{}{"service3"},
					},
					"test_queue_2": map[interface{}]interface{}{
						"read":  []interface{}{"service1"},
						"write": []interface{}{"service2"},
					},
				},
			},
		},
		hasher: NewHashProcessor("md5"),
	}

	_, err := loader.LoadFromFile(1.0)
	if err == nil {
		t.Error("TestLoadFromFile7 Did Not Return Error")
	}
}

//TestLoadFromFile8 unable to get SECRET env var
func TestLoadFromFile8(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    []byte{118, 69, 45},
			ReadFileErrorResult:   nil,
			GetEnvResult:          "",
		},
		cfgParser: &MockYAMLHandler{
			UnmarshalResult: nil,
			UnmarshalOutput: map[interface{}]interface{}{
				"version":             2.0,
				"port":                6010,
				"job_keep_minutes":    30,
				"job_timeout_minutes": 30,
				"queues": map[interface{}]interface{}{
					"test_queue_1": map[interface{}]interface{}{
						"read":  []interface{}{"service1", "service2"},
						"write": []interface{}{"service3"},
					},
					"test_queue_2": map[interface{}]interface{}{
						"read":  []interface{}{"service1"},
						"write": []interface{}{"service2"},
					},
				},
			},
		},
		hasher: NewHashProcessor("md5"),
	}

	_, err := loader.LoadFromFile(1.0)
	if err == nil {
		t.Error("TestLoadFromFile8 Did Not Return Error")
	}
}

//TestLoadFromFile9 unable to perform an MD5 hash on the secret
func TestLoadFromFile9(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    []byte{118, 69, 45},
			ReadFileErrorResult:   nil,
			GetEnvResult:          "testsecret",
		},
		cfgParser: &MockYAMLHandler{
			UnmarshalResult: nil,
			UnmarshalOutput: map[interface{}]interface{}{
				"version":             2.0,
				"port":                6010,
				"job_keep_minutes":    30,
				"job_timeout_minutes": 30,
				"queues": map[interface{}]interface{}{
					"test_queue_1": map[interface{}]interface{}{
						"read":  []interface{}{"service1", "service2"},
						"write": []interface{}{"service3"},
					},
					"test_queue_2": map[interface{}]interface{}{
						"read":  []interface{}{"service1"},
						"write": []interface{}{"service2"},
					},
				},
			},
		},
		hasher: &MockHashProcess{
			ProcessResult: "",
			ProcessError:  fmt.Errorf("Test Error"),
		},
	}

	_, err := loader.LoadFromFile(1.0)
	if err == nil {
		t.Error("TestLoadFromFile9 Did Not Return Error")
	}
}

/*BenchmarkLoadFromFile benchmarks loading a valid YAML config file - this gives an indication of time
taken to parse the YAML, it also uses an actual YAML file []byte rather than mocked */
func BenchmarkLoadFromFile(b *testing.B) {
	mockValidCfgFile := []byte(`
   version: 1
   port: 6010
   job_keep_minutes: 30
   job_timeout_minutes: 30
   queues:
      test_queue_1:
         read:
         - service1
         - service2
         write:
         - service3
      test_queue_2:
         read:
         - service3
         write:
         - service1
         - service2`)
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			ReadFileByteResult:    mockValidCfgFile,
			ReadFileErrorResult:   nil,
			GetEnvResult:          "testsecret",
		},
		cfgParser: NewConfigParser("yaml"),
		hasher:    NewHashProcessor("md5"),
	}

	for i := 0; i < b.N; i++ {
		_, _ = loader.LoadFromFile(1.0)
	}
}

//TestSaveToFile1 succesfully writes the given config to file
func TestSaveToFile1(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			DeleteFileResult:      nil,
			WiteFileResult:        nil,
		},
		cfgParser: &MockYAMLHandler{
			MarshalByteResult:  []byte{69, 45, 32},
			MarshalErrorResult: nil,
		},
		hasher: NewHashProcessor("md5"),
	}
	cfg := &models.Config{
		Version: 1.0,
		Port:    6006,
		Queues: map[string]*models.QueuePermissions{
			"test_queue_1": &models.QueuePermissions{
				Read:  []string{"service1", "service2"},
				Write: []string{"service3"},
			},
		},
	}

	if err := loader.SaveToFile(cfg); err != nil {
		t.Errorf("TestSaveToFile1 Failed: %s", err.Error())
	}
}

//TestSaveToFile2 is unable to delete the existing config before re-writing
func TestSaveToFile2(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  true,
			FileExistsErrorResult: nil,
			DeleteFileResult:      fmt.Errorf("Test Error"),
		},
		hasher: NewHashProcessor("md5"),
	}
	cfg := &models.Config{
		Version: 1.0,
		Port:    6006,
		Queues: map[string]*models.QueuePermissions{
			"test_queue_1": &models.QueuePermissions{
				Read:  []string{"service1", "service2"},
				Write: []string{"service3"},
			},
		},
	}

	if err := loader.SaveToFile(cfg); err == nil {
		t.Errorf("TestSaveToFile2 Did Not Return An Error")
	}
}

//TestSaveToFile3 encounters an error checking if the config exists before deciding to delete/write
func TestSaveToFile3(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  false,
			FileExistsErrorResult: fmt.Errorf("Test Error"),
		},
		hasher: NewHashProcessor("md5"),
	}
	cfg := &models.Config{
		Version: 1.0,
		Port:    6006,
		Queues: map[string]*models.QueuePermissions{
			"test_queue_1": &models.QueuePermissions{
				Read:  []string{"service1", "service2"},
				Write: []string{"service3"},
			},
		},
	}

	if err := loader.SaveToFile(cfg); err == nil {
		t.Errorf("TestSaveToFile3 Did Not Return An Error")
	}
}

//TestSaveToFile4 gives invalid argument to method
func TestSaveToFile4(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  false,
			FileExistsErrorResult: nil,
		},
		cfgParser: &MockYAMLHandler{
			MarshalByteResult:  nil,
			MarshalErrorResult: fmt.Errorf("Test Error"),
		},
		hasher: NewHashProcessor("md5"),
	}

	if err := loader.SaveToFile(nil); err == nil {
		t.Errorf("TestSaveToFile4 Did Not Return An Error")
	}
}

//TestSaveToFile5 has an error marshalling to YAML
func TestSaveToFile5(t *testing.T) {
	loader := ConfigLoad{
		filePath: "test_file_1.yml",
		fileHandler: &filesystem.MockFileSystem{
			FileExistsBoolResult:  false,
			FileExistsErrorResult: nil,
		},
		cfgParser: &MockYAMLHandler{
			MarshalByteResult:  nil,
			MarshalErrorResult: fmt.Errorf("Test Error"),
		},
		hasher: NewHashProcessor("md5"),
	}
	cfg := &models.Config{
		Version: 1.0,
		Port:    6006,
		Queues: map[string]*models.QueuePermissions{
			"test_queue_1": &models.QueuePermissions{
				Read:  []string{"service1", "service2"},
				Write: []string{"service3"},
			},
		},
	}

	if err := loader.SaveToFile(cfg); err == nil {
		t.Errorf("TestSaveToFile5 Did Not Return An Error")
	}
}
