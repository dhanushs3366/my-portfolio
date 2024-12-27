package db

import "log"

func (s *Store) sync() error {
	err := s.CreateUserTable()
	if err != nil {
		log.Println("couldnt create user table")
		return err
	}
	log.Println("user table created succesfully")

	err = s.CreateBlogTable()
	if err != nil {
		log.Println("couldnt create blog table")
		return err
	}
	log.Println("blog table created succesfully")

	err = s.CreateLogActivityTable()
	if err != nil {
		log.Println("couldnt create log activity table")
		return err
	}
	log.Println("log activity table created succesfully")

	return nil
}

func Init() (*Store, error) {
	store, err := ConnectToDB()
	if err != nil {
		return nil, err
	}

	err = store.sync()
	if err != nil {
		return nil, err
	}

	return store, err

}
