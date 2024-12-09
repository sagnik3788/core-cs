package main


import(
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request){
	// fmt.Println("<h1>http server started<h1>")
	w.Write([]byte("<h1>http server started</h1>"))
}

func main(){
	http.HandleFunc("/",handler)


	fmt.Println("server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil{
		fmt.Println("error in http server")
	}
}

