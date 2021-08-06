package main

import (
	"fmt"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	log "github.com/jeanphorn/log4go"
)

// UserResource is the REST layer to the User domain
type UserResource struct {
	// normally one would use DAO (data access object)
	users map[string]User
}

// WebService creates a new service that can handle REST requests for User resources.
func (u UserResource) WebService() *restful.WebService {
	//u.users = map[string]User{"1": User{ID: "1", Name: "admin", Age: 1}}
	u.users["1"] = User{ID: "1", Name: "admin", Age: 1}
	u.users["2"] = User{ID: "2", Name: "test", Age: 2}
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_JSON, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_JSON) // you can specify this per route as well

	tags := []string{"users"}

	ws.Route(ws.GET("/").To(u.findAllUsers).
		// docs
		Doc("get all users").
		Metadata("tags", tags).
		Writes([]User{}).
		Returns(200, "OK", []User{}))

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		// docs
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("integer").DefaultValue("1")).
		Metadata("tags", tags).
		Writes(User{}). // on the response
		Returns(200, "OK", User{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
		// docs
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Metadata("tags", tags).
		Reads(User{})) // from the request

	ws.Route(ws.PUT("").To(u.createUser).
		// docs
		Doc("create a user").
		Metadata("tags", tags).
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		// docs
		Doc("delete a user").
		Metadata("tags", tags).
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	return ws
}

// GET http://localhost:8080/users
//
func (u UserResource) findAllUsers(request *restful.Request, response *restful.Response) {
	list := []User{}
	// for each := range u.users {
	// 	fmt.Println(each, u.users[each])
	// 	list = append(list, u.users[each])
	// }

	for _, each := range u.users {
		fmt.Println(each)
		log.LOGGER("main").Info("users:", each)
		list = append(list, each)
	}
	response.WriteEntity(list)
}

// GET http://localhost:8080/users/1
//
func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr := u.users[id]
	if len(usr.ID) == 0 {
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
//
func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = *usr
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa</Name></User>
//
func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := User{ID: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = usr
		response.WriteHeaderAndEntity(http.StatusCreated, usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/users/1
//
func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func main() {
	log.LoadConfiguration("./log4go.json", "json")
	defer log.Close()
	log.LOGGER("main").Debug("test main")
	log.LOGGER("test").Debug("test")

	u := UserResource{map[string]User{}}
	restful.DefaultContainer.Add(u.WebService())

	log.LOGGER("main").Info("start listening on localhost:8080")
	log.LOGGER("main").Info(http.ListenAndServe(":8080", nil))
}

// User is just a sample type
type User struct {
	ID   string `json:"id" description:"identifier of the user"`
	Name string `json:"name" description:"name of the user" default:"admin"`
	Age  int    `json:"age" description:"age of the user" default:"18"`
}
