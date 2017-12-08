package admin

import (
    "fmt"
    "github.com/billyninja/pgtools/scanner"
    "html/template"
    "log"
    "time"
)

type col2html func(cl *scanner.Column) template.HTML
type val2html func(cl *scanner.Column, value interface{}) template.HTML

func InputHTML(cl *scanner.Column, value interface{}) template.HTML {
    archtype, fieldtype := field_type_translation(cl.Type)

    required := ""
    if cl.Nullable == "NO" {
        required = `required="required"`
    }

    max_length := ""
    if cl.CharMaxLength != nil {
        max_length = fmt.Sprintf(`max_length="%d"`, *cl.CharMaxLength)
    }

    value_str := ""
    if value != nil {
        switch v := value.(type) {
        case bool:
            if v == true {
                value_str = `checked="checked"`
            } else {
                value_str = `checked=""`
            }
            break
        case string:
            if archtype == "textarea" {
                value_str = v
            } else {
                value_str = `value="` + v + `"`
            }
            break
        case float64, float32:
            value_str = fmt.Sprintf(`value="%.2f"`, v)
            break
        case int, uint8, int8, uint16, int16, uint32, int32, int64:
            value_str = fmt.Sprintf(`value="%d"`, v)
            break
        case time.Time:
            value_str = fmt.Sprintf(`value="%s"`, v.Format("2006-01-02T15:04:05"))
            break
        default:
            value_str = fmt.Sprintf(`value="%s"`, v)
        }
    }

    var input template.HTML = ""
    if archtype == "textarea" {
        input = template.HTML(fmt.Sprintf(`<textarea name="%s" %s>%s</textarea>`, cl.Name, required, value_str))
    } else {
        input = template.HTML(fmt.Sprintf(`<input type="%s" name="%s" %s %s %s/>`, fieldtype, cl.Name, max_length, required, value_str))
    }

    return input
}

func ThHTML(cl *scanner.Column) template.HTML {
    return template.HTML("<th>" + cl.Name + "</th>")
}

func LabelHTML(cl *scanner.Column) template.HTML {
    return template.HTML("<label>" + cl.Name + "</label>")
}

func LabelAndInputHTML(cl *scanner.Column, value interface{}) template.HTML {

    label := LabelHTML(cl)
    input := InputHTML(cl, value)

    // TODO-improvement: configurable label-class, input class
    return `<div class="form-group"><div class="">` + label + `</div><div class="">` + input + `</div></div>`
}

func TdHTML(cl *scanner.Column, value interface{}) template.HTML {
    return template.HTML("<td>" + format_value(value) + "</td>")
}

func field_type_translation(column_type string) (string, string) {

    switch column_type {
    case "character varying":
        return "input", "text"
    case "timestamp without time zone":
        return "input", "date"
    case "timestamp with time zone":
        return "input", "date"
    case "numeric":
        return "input", "number"
    case "text":
        return "textarea", "textarea"
    case "boolean":
        return "input", "checkbox"
    default:
        log.Printf("\n\nUnmapped PSQL field type %s\n\n", column_type)
        return "input", "text"
    }

    return "err", "err"
}

func format_value(value interface{}) string {

    switch v := value.(type) {
    case bool:
        if v == true {
            return "true"
        } else {
            return "false"
        }
    case string:
        return v
    case float64, float32:
        return fmt.Sprintf(`%.2f`, v)
    case int, uint8, uint16, int64:
        return fmt.Sprintf(`%d`, v)
    case time.Time:
        return fmt.Sprintf(`%s`, v.Format("2006-01-02T15:04:05"))
    default:
        return fmt.Sprintf(`%s`, v)
    }

    return "[---]"
}
