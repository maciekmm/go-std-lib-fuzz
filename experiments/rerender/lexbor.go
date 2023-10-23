package rerender

//#cgo CFLAGS: -g -Wall
//#cgo LDFLAGS: -llexbor -L/usr/local/lib
//
//#include "lexbor/core/str.h"
//#include "lexbor/html/html.h"
//
//char * parse_tree(unsigned char* html, size_t html_len)
//{
//    lxb_status_t status;
//
//    lxb_html_document_t *document = lxb_html_document_create();
//    if (document == NULL) {
//        return NULL;
//    }
//
//    status = lxb_html_document_parse(document, html, html_len);
//    if (status != LXB_STATUS_OK) {
//        return NULL;
//    }
//
//    lexbor_str_t* str = lexbor_str_create();
//
//    status = lxb_html_serialize_tree_str(lxb_dom_interface_node(document), str);
//    if (status != LXB_STATUS_OK) {
//        return NULL;
//    }
//
//    /* document_destroy seems to free the lexbor_str_t as well, hence we copy this to a char */
//    char* result = malloc(sizeof(unsigned char) * str->length + 1);
//    memset(result, '\0', sizeof(unsigned char) * str->length + 1);
//    strncpy(result, (char *) str->data, str->length);
//
//    /* Destroy document */
//    lxb_html_document_destroy(document);
//
//    return result;
//}
import "C"
import (
	"errors"
	"unicode/utf8"
	"unsafe"
)

func LexborParseRender(in string) (string, error) {
	input := C.CString(in)
	defer C.free(unsafe.Pointer(input))

	parsed := C.parse_tree((*C.uchar)(unsafe.Pointer(input)), C.ulong(len(in)))
	if parsed == nil {
		return "", errors.New("tree parsing error")
	}

	rendered := C.GoString(parsed)
	defer C.free(unsafe.Pointer(parsed))

	return removeInvalidUTF8(rendered), nil
}

func removeInvalidUTF8(s string) string {
	var result []rune
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r != utf8.RuneError {
			result = append(result, r)
		}
		s = s[size:]
	}
	return string(result)
}
