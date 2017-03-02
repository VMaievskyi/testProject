package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type DBHandler struct {
	db               *gorm.DB
	connectionString string
}

type YanaName struct {
	Id   int    `gorm:"AUTO_INCREMENT"`
	Name string `sql:"type:VARCHAR(50) CHARACTER SET utf8 COLLATE utf8_general_ci"`
}

func main() {

	h := DBHandler{
		connectionString: fmt.Sprintf(
			"%s:%s@tcp(%s)/obzyvalki?parseTime=true&charset=utf8",
			os.Getenv("MYSQL_USER"),
			os.Getenv("WEB_DB_PASSWORD"),
			os.Getenv("WEB_DB_HOST"))}

	h.connect()
	defer h.close()

	h.initDevConfigs()
	h.insertTestData()

	router := mux.NewRouter()
	router.HandleFunc("/whoIsYana", h.yanaRandomizer).Methods("GET")
	router.HandleFunc("/whoIsYana/all", h.getAllHandler).Methods("GET")
	router.HandleFunc("/whoIsYana/{name}", h.addValue).Methods("put")
	router.HandleFunc("/whoIsYana/isAlive", h.isAlive).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8081")
}

func (h *DBHandler) isAlive(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("YES"))
}

func (h *DBHandler) yanaRandomizer(rw http.ResponseWriter, req *http.Request) {
	rand.Seed(time.Now().Unix())
	var count int
	h.db.Model(&YanaName{}).Count(&count)
	returnObject := YanaName{}
	h.db.First(&returnObject, rand.Intn(count-1)+1)
	log.Printf("loaded value is %s", returnObject.Name)
	rw.Write([]byte(returnObject.Name))
}

func (h *DBHandler) addValue(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	insertValue := YanaName{Name: vars["name"]}
	h.db.Save(&insertValue)

}

func (h *DBHandler) getAllHandler(rw http.ResponseWriter, req *http.Request) {
	var yanaNames []*YanaName
	h.db.Find(&yanaNames)
	result := ToString(",", yanaNames)
	fmt.Printf("%v", yanaNames)
	rw.Write([]byte(result))
}

func (h *DBHandler) connect() {
	db, err := gorm.Open("mysql", h.connectionString)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	h.db = &db
	h.waitDataBase()
}

func (h *DBHandler) close() {
	h.db.Close()
}

func (h *DBHandler) waitDataBase() {
	log.Println("<<<<<<<<<<<<<<ping >>>>>>>>>>>>>>>>>>")
	pingResult := h.db.Exec("SELECT 1")
	if pingResult.Error != nil {
		log.Println("<<<<<<<<<<<<<< failed >>>>>>>>>>>>>>>>>>")
		time.Sleep(100 * time.Millisecond)
		h.waitDataBase()
	}
}

func (h *DBHandler) insertTestData() {
	for _, val := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9",
		"10", "11"} {
		newName := YanaName{Name: val}
		h.db.Create(&newName)
	}
}

func (h *DBHandler) initDevConfigs() {
	h.db.LogMode(true) // This would be off in production.
	h.db.DropTable(&YanaName{})
	h.db.AutoMigrate(&YanaName{}) // nice for development, but I would probably just write a SQL script to do this.
	h.db.Model(&YanaName{}).AddIndex("yana_names_id", "id")
}

func (n *YanaName) String() string {
	return n.Name
}
