package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

// Los IDs deben conservar su valor o rechazarse; nunca truncarse al convertir a uint.
func TestParseIDRespetaElRangoDeUint(t *testing.T) {
	cases := []struct {
		input string
		want  uint64
		valid bool
	}{
		{input: "1", want: 1, valid: true},
		{input: "4294967295", want: 4294967295, valid: true},
		{input: "4294967296", want: 4294967296, valid: strconv.IntSize == 64},
		{input: "4294967297", want: 4294967297, valid: strconv.IntSize == 64},
		{input: "18446744073709551615", want: 18446744073709551615, valid: strconv.IntSize == 64},
		{input: "18446744073709551616"},
		{input: "0"},
		{input: "-1"},
		{input: "no-es-numero"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			context.Params = gin.Params{{Key: "id", Value: tc.input}}

			id, ok := parseID(context)

			if ok != tc.valid {
				t.Fatalf("ID %q: válido=%t, se esperaba %t", tc.input, ok, tc.valid)
			}
			if tc.valid && uint64(id) != tc.want {
				t.Fatalf("ID alterado: obtenido=%d, se esperaba %d", id, tc.want)
			}
			if !tc.valid && (id != 0 || response.Code != http.StatusBadRequest) {
				t.Fatalf("ID inválido: id=%d, status=%d", id, response.Code)
			}
		})
	}
}
