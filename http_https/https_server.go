package main


import(
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request){
	// fmt.Println("<h1>http server started<h1>")
	w.Write([]byte("<h1>https server started</h1>"))
}

func main(){
	http.HandleFunc("/",handler)


	fmt.Println("server is running on port 8443")
	err := http.ListenAndServeTLS(":8443","server.crt", "server.key", nil)

	if err != nil{
		fmt.Println("error in https server")
	}
}

