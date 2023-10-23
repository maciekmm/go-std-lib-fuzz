package main

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -llexbor -L/usr/local/lib
// #include <stdlib.h>
// #include <lexbor/html/tokenizer.h>
// #include <inttypes.h>
//
// static lxb_html_token_t *
// token_callback(lxb_html_tokenizer_t *tkz, lxb_html_token_t *token, void *ctx) {
//     lexbor_hash_t *tags = lxb_html_tokenizer_tags(tkz);
//
//     bool is_close = token->type & LXB_HTML_TOKEN_TYPE_CLOSE;
//     bool is_self_close = token->type & LXB_HTML_TOKEN_TYPE_CLOSE_SELF;
//
//     // I can't get this to work consistently.
//     //if (!is_close) {
//     //   lxb_html_tokenizer_set_state_by_tag(tkz, false, token->tag_id, LXB_NS_HTML);
//     //}
//
//     TokenCallback(tkz, token->begin, (int)(token->end - token->begin), token->tag_id, token->type);
//
//     return token;
// }
//
// static void register_token_callback(lxb_html_tokenizer_t *tkz) {
// 	lxb_html_tokenizer_callback_token_done_set(tkz, token_callback, NULL);
// };
//
import "C"
import (
	"errors"
	"sync"
	"unsafe"

	"golang.org/x/net/html"
)

var tokenizerTokens map[unsafe.Pointer][]Token = make(map[unsafe.Pointer][]Token)
var tokensMutex *sync.RWMutex = &sync.RWMutex{}

type LexborTagID uint64

const (
	LexborTagEndOfFile uint64 = 0x0001
	LexborTagText      uint64 = 0x0002
	LexborTagComment   uint64 = 0x0004
	LexborTagDoctype   uint64 = 0x0005
)

type LexborType int

const (
	LexborTypeOpen        int = 0x0000
	LexborTypeClose       int = 0x0001
	LexborTypeCloseSelf   int = 0x0002
	LexborTypeForceQuirks int = 0x0004
	LexborTypeDone        int = 0x0008
)

//export TokenCallback
func TokenCallback(tokenizerPtr unsafe.Pointer, cName *C.char, cNameLen C.int, cTypeId C.uintptr_t, cType C.int) {
	name := C.GoStringN(cName, cNameLen)
	// fmt.Printf("%x Tag name: %s, tag id: %x, type: %d\n", tokenizerPtr, name, cTypeId, cType)

	tagId := uint64(cTypeId)
	typ := int(cType)

	var tokenType html.TokenType
	switch {
	case tagId == LexborTagDoctype:
		tokenType = html.DoctypeToken
		name = "DOCTYPE"
	case tagId == LexborTagEndOfFile:
		tokenType = html.ErrorToken
	case tagId == LexborTagComment:
		tokenType = html.TextToken
		name = "#text"
	case tagId == LexborTagText:
		tokenType = html.TextToken
		name = "#text"
	case typ == LexborTypeOpen:
		tokenType = html.StartTagToken
		name = lower([]byte(name))
	case typ == LexborTypeClose:
		tokenType = html.EndTagToken
		name = lower([]byte(name))
	case typ == LexborTypeCloseSelf:
		tokenType = html.StartTagToken
		name = lower([]byte(name))
	default:
		tokenType = html.TextToken
		name = "#text"
	}

	tokensMutex.Lock()
	defer tokensMutex.Unlock()
	tokenizerTokens[tokenizerPtr] = append(tokenizerTokens[tokenizerPtr], Token{
		Name: name,
		Type: tokenType,
	})
}

func lower(b []byte) string {
	for i, c := range b {
		if 'A' <= c && c <= 'Z' {
			b[i] = c + 'a' - 'A'
		}
	}
	return string(b)
}

func LexborTokenize(data string) ([]Token, error) {
	input := unsafe.Pointer(C.CString(data))
	inputSize := C.ulong(len(data))
	defer C.free(unsafe.Pointer(input))

	tkz := C.lxb_html_tokenizer_create()
	defer C.lxb_html_tokenizer_destroy(tkz)

	tkzPtr := unsafe.Pointer(tkz)

	tokensMutex.Lock()
	tokenizerTokens[tkzPtr] = []Token{}
	tokensMutex.Unlock()

	defer func() {
		tokensMutex.Lock()
		delete(tokenizerTokens, tkzPtr)
		tokensMutex.Unlock()
	}()

	status := C.lxb_html_tokenizer_init(tkz)
	if status != 0 {
		return nil, errors.New("tonizer not initialized successfully")
	}

	C.register_token_callback(tkz)

	status = C.lxb_html_tokenizer_begin(tkz)
	if status != 0 {
		return nil, errors.New("failed to prepare tokenizer object for parsing")
	}

	status = C.lxb_html_tokenizer_chunk(tkz, (*C.uchar)(input), inputSize)
	if status != 0 {
		return nil, errors.New("failed to parse the html data")
	}

	status = C.lxb_html_tokenizer_end(tkz)
	if status != 0 {
		return nil, errors.New("failed to ending of parsing the html data")
	}

	return tokenizerTokens[tkzPtr], nil
}
