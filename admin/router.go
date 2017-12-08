package admin

import (
    "github.com/billyninja/pgtools/connector"
    "github.com/billyninja/pgtools/scanner"
    "github.com/julienschmidt/httprouter"
    "net/http"
)

var actions = []string{
    "create_one", // GET, POST
    "list",       // GET, POST(bulk action handler)
    "edit_one",   // GET, POST
    "delete",     // DELETE
}

func AutoRoute(rt *httprouter.Router, conn *connector.Connector, baseroute string) {
    allTables := scanner.GetAllTables(conn)

    for _, tb := range allTables {

        create_GET := baseroute + "/" + string(tb.Name) + "/create"
        rt.GET(create_GET, StubHandler)

        create_POST := baseroute + "/" + string(tb.Name) + "/create"
        rt.POST(create_POST, StubHandler)

        edit_GET := baseroute + "/" + string(tb.Name) + "/edit"
        rt.GET(edit_GET, StubHandler)

        edit_POST := baseroute + "/" + string(tb.Name) + "/edit"
        rt.POST(edit_POST, StubHandler)

        list_GET := baseroute + "/" + string(tb.Name) + "/list"
        rt.GET(list_GET, StubHandler)

        list_POST := baseroute + "/" + string(tb.Name) + "/list"
        rt.POST(list_GET, StubHandler)

        delete_DELETE := baseroute + "/" + string(tb.Name) + "/delete"
        rt.DELETE(delete_DELETE, StubHandler)

        println("Routing for table", string(tb.Name), ":")
        println("\t", "GET", create_GET)
        println("\t", "POST", create_POST)
        println("\t", "GET", edit_GET)
        println("\t", "POST", edit_POST)
        println("\t", "GET", list_GET)
        println("\t", "POST", list_POST)
        println("\t", "DELETE", delete_DELETE)
        println("---")

    }
}

func StubHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.WriteHeader(http.StatusOK)
    return
}
