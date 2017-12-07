package admin

import "io"

type GenericView interface {
    GetTitle() string
    PartialHTML(io.Writer)
    CompleteHTML(io.Writer)
}
