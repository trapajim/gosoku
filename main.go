package main

import "gosoku/system/router"

func main() {
	/* dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	dbConn, err := sql.Open(`postgres`, connection)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	err = dbConn.Ping()
	userRepo := repository.NewPsqlUserRepository(dbConn) */
	router := router.NewRouter()
	router.RegisterRoutes()

}
