package rnd

import (
//    "math/rand"
    "time"
)


// varying char
func r_char() string {
    out := "out r_char"
    return out
}


// varying char
func r_var_char(min, max int) string {
    out := "out r_var_char"
    return out
}


// text
func r_text() string {
    return r_var_char(1024, 2048)
}


// SMALLINT, INTEGER, AND BIGINT
func r_int(max int64) int64 {
    out := int64(123456)
    return out
}


// NUMERIC with varying precision
func r_numeric(max float64, places uint8) float64 {
    out := 3.14
    return out
}


// timestamp etc.
func r_datetime(relative int8) time.Time {
    return time.Now()
}


// timestamp       8 bytes WITHOUT timezone
// timestamp       8 bytes WITH timezone
// date            4 bytes (no time of day)
// time            8 bytes WITHOUT timezone
// time           12 bytes WITH timezone
// interval       16bytes
func r_interval(relative int8) time.Duration {
    return 10 * time.Second
}


// https://www.postgresql.org/docs/9.3/static/datatype-binary.html
func r_byte_array(min, max int64) []byte {
    out := "out r_byte_array"
    return []byte(out)
}
