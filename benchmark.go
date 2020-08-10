package BenchmarkDB

type storage interface {
	Insert()
	Select()
	Delete()
}

func insertDB(s storage) {
	s.Insert()
}

func selectDB(s storage) {
	s.Select()
}

func deleteDB(s storage) {
	s.Delete()
}
