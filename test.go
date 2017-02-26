package main

import (
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
	"os"
)

type DBHandler struct {
	db *gorm.DB
}

type YanaName struct {
	Id   int `gorm:"AUTO_INCREMENT"`
	Name string `sql:"type:VARCHAR(50) CHARACTER SET utf8 COLLATE utf8_general_ci"`
}

func main() {
	connectionString := "testUser:" + os.Getenv("WEB_DB_PASSWORD") + "@tcp(" + os.Getenv("WEB_DB_HOST") + ")/obzyvalki?parseTime=true&charset=utf8"
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>"+connectionString)
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	db.LogMode(true) // This would be off in production.
	defer db.Close()
	db.DropTable(&YanaName{})
	db.AutoMigrate(&YanaName{}) // nice for development, but I would probably just write a SQL script to do this.
	db.Model(&YanaName{}).AddIndex("yana_names_id", "id")

	h := DBHandler{db: &db}
	for _, val := range []string{"попа", "сосіска", "малявка", "красуня", "золотце", "старушка", "вредіна", "умнічка", "казявка",
		"бузюська", "пісюська"} {
		newName := YanaName{Name: val}
		h.db.Create(&newName)
	}

	rand.Seed(time.Now().Unix())

	router := mux.NewRouter()
	router.HandleFunc("/whoIsYana", h.yanaRandomizer).Methods("GET")
	router.HandleFunc("/whoIsYana/all", h.getAllHandler).Methods("GET")
	router.HandleFunc("/whoIsYana/{name}", h.addValue).Methods("put")
	router.HandleFunc("/whoIsYana/isAlive", h.isAlive).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8081")
}

func (h *DBHandler) isAlive(rw http.ResponseWriter, req *http.Request){
	rw.Write([]byte("YES"))
}

func (h *DBHandler) yanaRandomizer(rw http.ResponseWriter, req *http.Request) {
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
	var yanaNames []YanaName
	h.db.Find(&yanaNames)
	var result string
	for _, res := range yanaNames {
		result += res.Name + " "
	}
	rw.Write([]byte(result))
}
