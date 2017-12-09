package admin

import (
    "github.com/billyninja/pgtools/connector"
    "github.com/billyninja/pgtools/scanner"
    "github.com/julienschmidt/httprouter"
    "log"
    "net/http"
    "strings"
)

var allTables []*scanner.Table
var baseroute_ string
var conn_ *connector.Connector

func AutoRoute(rt *httprouter.Router, conn *connector.Connector, baseroute string) {
    allTables = scanner.GetAllTables(conn)
    baseroute_ = baseroute
    conn_ = conn

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
        rt.GET(list_GET, ListHandler)

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

func getTable(r *http.Request) *scanner.Table {
    spl1 := strings.Split(r.URL.String(), baseroute_)
    table_name := scanner.TableName(strings.Split(spl1[1], "/")[1])
    println(">>>", table_name)

    for _, tb := range allTables {
        if tb.Name == table_name {
            return tb
        }
    }

    return nil
}

func StubHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.WriteHeader(http.StatusOK)
    return
}

func ListHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "text/html")

    tb := getTable(r)
    println("here", tb)
    q := QueryListAll(tb.Name)
    println("here", q)
    rows, err := conn_.Sel(q)
    if err != nil {
        log.Printf("Couldn't query table data!")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    lv := NewListView(tb, rows)
    err = lv.PartialHTML(w)
    if err != nil {
        log.Printf("Template Error!\n%s\n", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    return
}
