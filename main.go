package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()
	connection.DB_CONN()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", project_detail).Methods("GET")
	route.HandleFunc("/add-project", Add_Project).Methods("POST")
	route.HandleFunc("/delete-project/{index}", Delete_Project).Methods("GET")
	route.HandleFunc("/update-project/{id}", Form_Update_Project).Methods("GET")

	fmt.Println("Server berjalan pada port 4000")
	http.ListenAndServe("localhost:4000", route)

}

type Project struct {
	Project_name string
	Description  string
	Start_Date   string
	End_Date     string
}

var projects = []Project{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("view/index.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// dataProject := map[string]interface{}{
	// 	"Projects": projects,
	// }

	dataProject, errQuery := connection.Conn.Query(context.Background(), "SELECT project_name, description FROM tb_project")

	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	var result []Project

	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.Project_name, &each.Description)

		if err != nil {
			fmt.Println("Message : " + err.Error())
			return
		}

		result = append(result, each)

	}

	fmt.Println(result)
	resData := map[string]interface{}{
		"Projects": result,
	}
	tmpt.Execute(w, resData)
}

func Add_Project(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}
	project_name := r.PostForm.Get("project_name")
	description := r.PostForm.Get("description")
	start_date := r.PostForm.Get("start-date")
	end_date := r.PostForm.Get("end-date")

	var data = Project{
		Project_name: project_name,
		Description:  description,
		Start_Date:   start_date,
		End_Date:     end_date,
	}

	projects = append(projects, data)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("view/project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("view/Input_Form.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func project_detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("view/blog.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	dataProjectDetail := Project{}

	for index, data := range projects {
		if index == id {
			dataProjectDetail = Project{
				Project_name: data.Project_name,
				Description:  data.Description,
				Start_Date:   data.Start_Date,
				End_Date:     data.End_Date,
			}
		}
	}

	data := map[string]interface{}{
		"ProjectDetail": dataProjectDetail,
	}

	tmpt.Execute(w, data)
}

func Form_Update_Project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("view/update_project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	dataProjectDetail := Project{}
	for index, data := range projects {
		if index == id {
			dataProjectDetail = Project{
				Project_name: data.Project_name,
				Description:  data.Description,
				Start_Date:   data.Start_Date,
				End_Date:     data.End_Date,
			}
		}
	}

	data := map[string]interface{}{
		"ProjectDetail": dataProjectDetail,
	}
	tmpt.Execute(w, data)

}

func Delete_Project(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	projects = append(projects[:index], projects[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}
