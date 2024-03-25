package stringx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsafeToBytes(t *testing.T) {
	testCase := []struct {
		name string
		val  string
		want []byte
	}{
		{
			name: "normal conversion",
			val:  "hello",
			want: []byte("hello"),
		},
		{
			name: "emoji coversion",
			val:  "😀!hello world",
			want: []byte("😀!hello world"),
		},
		{
			name: "chinese coversion",
			val:  "你好 世界！",
			want: []byte("你好 世界！"),
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			val := UnsafeToBytes(tt.val)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestUnsafeToString(t *testing.T) {
	testCase := []struct {
		name string
		val  func(t *testing.T) []byte
		want string
	}{
		{
			name: "normal conversion",
			val: func(t *testing.T) []byte {
				return []byte("hello")
			},
			want: "hello",
		},
		{
			name: "emoji coversion",
			val: func(t *testing.T) []byte {
				return []byte("😀!hello world")
			},
			want: "😀!hello world",
		},
		{
			name: "chinese coversion",
			val: func(t *testing.T) []byte {
				return []byte("你好 世界！")
			},
			want: "你好 世界！",
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.val(t)
			val := UnsafeToString(b)
			assert.Equal(t, tt.want, val)
		})
	}
}
