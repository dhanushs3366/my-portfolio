package initializers

import (
	"database/sql"
	blog "dhanushs3366/my-portfolio/services/blog"
	"dhanushs3366/my-portfolio/services/db"
	"dhanushs3366/my-portfolio/services/logger"
	"dhanushs3366/my-portfolio/services/user"
)

func syncDB(db *sql.DB) error {
	userStore := user.NewUserStore(db)
	err := userStore.CreateUserTable()
	if err != nil {
		return err
	}
	logStore := logger.NewLogStore(db)
	err = logStore.CreateLogActivityTable()

	if err != nil {
		return err
	}

	blogStore := blog.NewBlogStore(db)
	err = blogStore.CreateBlogTable()
	if err != nil {
		return err
	}
	return nil
}

func Init() (*sql.DB, error) {
	db, err := db.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = syncDB(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
