package dal

type DatabaseAdaptor struct {
	bufferedDatabaseChannel    chan string
	bufferedHealthCheckChannel chan string
	envVars                    string
}

func (dbAdaptor *DatabaseAdaptor) DatabaseLoop() {
	for {

	}
}
