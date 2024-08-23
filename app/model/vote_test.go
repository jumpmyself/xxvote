package model

import (
	"context"
	"fmt"
	"testing"
)

func TestDelVoteV1(t *testing.T) {
	NewMysql()
	ret, _ := GetVoteV5(3)
	fmt.Printf("ret:%+v\n", ret)
}

func TestGetVoteHistoryV1(t *testing.T) {
	NewMysql()
	NewRdb()
	GetVoteHistoryV1(context.TODO(), 1, 1)
}

func TestGetJwt(t *testing.T) {
	str, _ := GetJwt(1, "想想编程")
	fmt.Printf("str:%s", str)
}

func TestCheckJwt(t *testing.T) {
	str := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MSwiTmFtZSI6IuaDs-aDs-e8lueoiyIsImlzcyI6Iummmemmmee8lueoiyIsInN1YiI6IuWQjuWLpOmDqOenpuW4iOWChSIsImF1ZCI6WyJBbmRyb2lkIiwiSU9TIiwiSDUiXSwiZXhwIjoxNzA4NzQ4ODEzLCJuYmYiOjE3MDg3NDUyMjMsImlhdCI6MTcwODc0NTIxMywianRpIjoiVGVzdC0xIn0.5BV0TGnBgIqJWoDWDWxLHVlM8Hal7ff5DuZezXN6JeI"
	token, _ := CheckJwt(str)
	fmt.Printf("token:%+v\n", token)
}
