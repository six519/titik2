package main

import (
	"errors"
	"fmt"
	"strconv"
)

// function return type
const (
	RET_TYPE_NONE = iota
	RET_TYPE_STRING
	RET_TYPE_INTEGER
	RET_TYPE_FLOAT
	RET_TYPE_ARRAY
	RET_TYPE_ASSOCIATIVE_ARRAY
	RET_TYPE_BOOLEAN
)

// function argument types
const (
	ARG_TYPE_NONE = iota
	ARG_TYPE_STRING
	ARG_TYPE_INTEGER
	ARG_TYPE_FLOAT
	ARG_TYPE_ARRAY
	ARG_TYPE_ASSOCIATIVE_ARRAY
	ARG_TYPE_BOOLEAN
)

type FunctionReturn struct {
	Type                  int
	StringValue           string
	IntegerValue          int
	FloatValue            float64
	BooleanValue          bool
	ArrayValue            []FunctionReturn
	AssociativeArrayValue map[string]FunctionReturn
}

type FunctionArgument struct {
	Type                  int
	StringValue           string
	IntegerValue          int
	FloatValue            float64
	BooleanValue          bool
	ArrayValue            []FunctionArgument
	AssociativeArrayValue map[string]FunctionArgument
}

type Execute func([]FunctionArgument, *error, *[]Variable, *[]Function, string, *[]string, *GlobalSettingsObject, int, int, string) FunctionReturn

type Function struct {
	Name          string
	IsNative      bool
	Tokens        []Token
	Run           Execute
	Arguments     []Token
	ArgumentCount int
}

func DumpFunction(functions []Function) {
	fmt.Printf("====================================\n")

	for x := 0; x < len(functions); x++ {
		fmt.Printf("Function Name: %s\n", functions[x].Name)
		fmt.Printf("Argument Count: %d\n", functions[x].ArgumentCount)

		if functions[x].IsNative {
			fmt.Println("Is Native: Yes")
		} else {
			fmt.Println("Is Native: No")
		}

		fmt.Printf("====================================\n")
	}
}

func isFunctionExists(token Token, globalFunctionArray []Function) (bool, int) {

	for x := 0; x < len(globalFunctionArray); x++ {
		if globalFunctionArray[x].Name == token.Value {
			return true, x
		}
	}

	return false, 0
}

func isParamExists(token Token, functionParams []Token) bool {

	for x := 0; x < len(functionParams); x++ {
		if functionParams[x].Value == token.Value {
			return true
		}
	}

	return false
}

func validateParameters(arguments []FunctionArgument, errMessage *error, line_number int, column_number int, file_name string, param_index int, param_expected int) bool {
	var ret bool = true
	var err_msg string = "a Nil"

	if arguments[param_index].Type != param_expected {
		ret = false
	}

	if !ret {
		switch param_expected {
		case ARG_TYPE_STRING:
			err_msg = "a string"
		case ARG_TYPE_INTEGER:
			err_msg = "an integer"
		case ARG_TYPE_FLOAT:
			err_msg = "a float"
		case ARG_TYPE_ARRAY:
			err_msg = "a lineup"
		case ARG_TYPE_ASSOCIATIVE_ARRAY:
			err_msg = "a glossary"
		case ARG_TYPE_BOOLEAN:
			err_msg = "a boolean"
		default:
		}
		*errMessage = errors.New("Error: Parameter " + strconv.Itoa(len(arguments)-param_index) + " must be " + err_msg + " type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	}

	return ret
}

func defineFunction(globalFunctionArray *[]Function, funcName string, funcExec Execute, argumentCount int, isNative bool) {
	function := Function{Name: funcName, IsNative: isNative, Run: funcExec, ArgumentCount: argumentCount}
	//append to global functions
	*globalFunctionArray = append(*globalFunctionArray, function)
}

func initNativeFunctions(globalFunctionArray *[]Function) {

	//p(<anyvar>)
	defineFunction(globalFunctionArray, "p", P_execute, 1, true)

	//ex(<integer>)
	defineFunction(globalFunctionArray, "ex", Ex_execute, 1, true)

	//exe(<string>)
	defineFunction(globalFunctionArray, "exe", Exe_execute, 1, true)

	//abt(<string>)
	defineFunction(globalFunctionArray, "abt", Abt_execute, 1, true)

	//!(<bool>)
	defineFunction(globalFunctionArray, "!", ReverseBoolean_execute, 1, true)

	//zzz(<integer>)
	defineFunction(globalFunctionArray, "zzz", Zzz_execute, 1, true)

	//r(<string>)
	defineFunction(globalFunctionArray, "r", R_execute, 1, true)

	//rp(<string>)
	defineFunction(globalFunctionArray, "rp", Rp_execute, 1, true)

	//toi(<anyvar>)
	defineFunction(globalFunctionArray, "toi", Toi_execute, 1, true)

	//tos(<anyvar>)
	defineFunction(globalFunctionArray, "tos", Tos_execute, 1, true)

	//len(<lineup>)
	defineFunction(globalFunctionArray, "len", Len_execute, 1, true)

	//sav()
	defineFunction(globalFunctionArray, "sav", Sav_execute, 0, true)

	//sc(<integer>)
	defineFunction(globalFunctionArray, "sc", Sc_execute, 1, true)

	//i(<string>)
	defineFunction(globalFunctionArray, "i", I_execute, 1, true)

	//gcp()
	defineFunction(globalFunctionArray, "gcp", Gcp_execute, 0, true)

	//in(<anyvar>)
	defineFunction(globalFunctionArray, "in", In_execute, 1, true)

	//la(<lineup>, <anyvar>)
	defineFunction(globalFunctionArray, "la", La_execute, 2, true)

	//lp(<lineup>, <integer>)
	defineFunction(globalFunctionArray, "lp", Lp_execute, 2, true)

	//WEB FUNCTIONALITIES
	//http_au(<string>, <string>) - Add URL
	defineFunction(globalFunctionArray, "http_au", Http_au_execute, 2, true)

	//http_run(<string>) - Run server
	defineFunction(globalFunctionArray, "http_run", Http_run_execute, 1, true)

	//http_p(<string>) - Print
	defineFunction(globalFunctionArray, "http_p", Http_p_execute, 1, true)

	//http_h(<string>, <string>) - Set header
	defineFunction(globalFunctionArray, "http_h", Http_h_execute, 2, true)

	//http_gm() - Get method
	defineFunction(globalFunctionArray, "http_gm", Http_gm_execute, 0, true)

	//http_su(<string>, <string>) - Static URL
	defineFunction(globalFunctionArray, "http_su", Http_su_execute, 2, true)

	//http_gq(<string>) - Get query
	defineFunction(globalFunctionArray, "http_gq", Http_gq_execute, 1, true)

	//http_gfp(<string>) - Get form POST
	defineFunction(globalFunctionArray, "http_gfp", Http_gfp_execute, 1, true)

	//http_lt(<string>) - Load template
	defineFunction(globalFunctionArray, "http_lt", Http_lt_execute, 2, true)

	//http_gp() - Get path
	defineFunction(globalFunctionArray, "http_gp", Http_gp_execute, 0, true)

	//http_cr() - HTTP client request
	defineFunction(globalFunctionArray, "http_cr", Http_cr_execute, 4, true)

	//MYSQL FUNCTIONALITIES
	//mysql_set(<string>, <string>, <string>, <string>) - Set connection
	defineFunction(globalFunctionArray, "mysql_set", Mysql_set_execute, 4, true)

	//mysql_q(<string>, <string>) - Query
	defineFunction(globalFunctionArray, "mysql_q", Mysql_q_execute, 2, true)

	//mysql_cr(<string>) - Clear resources
	defineFunction(globalFunctionArray, "mysql_cr", Mysql_cr_execute, 1, true)

	//mysql_fa(<string>, <string>) - Fetch all
	defineFunction(globalFunctionArray, "mysql_fa", Mysql_fa_execute, 2, true)

	//STRING FUNCTIONALITIES
	//str_rpl(<string>, <string>, <string>) - String replace
	defineFunction(globalFunctionArray, "str_rpl", Str_rpl_execute, 3, true)

	//str_spl(<string>, <string>) - String split
	defineFunction(globalFunctionArray, "str_spl", Str_spl_execute, 2, true)

	//str_l(<string>) - String to lower
	defineFunction(globalFunctionArray, "str_l", Str_l_execute, 1, true)

	//str_u(<string>) - String to upper
	defineFunction(globalFunctionArray, "str_u", Str_u_execute, 1, true)

	//str_t(<string>) - String trim
	defineFunction(globalFunctionArray, "str_t", Str_t_execute, 1, true)

	//str_chr(<string>) - Integer to character string
	defineFunction(globalFunctionArray, "str_chr", Str_chr_execute, 1, true)

	//str_ord(<string>) - Character to integer code point
	defineFunction(globalFunctionArray, "str_ord", Str_ord_execute, 1, true)

	//str_sub(<string>, <integer>, <integer>) - Substring
	defineFunction(globalFunctionArray, "str_sub", Str_sub_execute, 3, true)

	//str_ind(<string>, <string>) - Get index of specified substring
	defineFunction(globalFunctionArray, "str_ind", Str_ind_execute, 2, true)

	//MATH FUNCTIONALITIES
	//abs(<float/integer>) - Absolute
	defineFunction(globalFunctionArray, "abs", Abs_execute, 1, true)

	//acs(<float/integer>) - Arccosine
	defineFunction(globalFunctionArray, "acs", Acs_execute, 1, true)

	//acsh(<float/integer>) - Inverse hyperbolic cosine of x
	defineFunction(globalFunctionArray, "acsh", Acsh_execute, 1, true)

	//asn(<float/integer>) - arc sine of a number
	defineFunction(globalFunctionArray, "asn", Asn_execute, 1, true)

	//asnh(<float/integer>) - returns the inverse hyperbolic sine of x
	defineFunction(globalFunctionArray, "asnh", Asnh_execute, 1, true)

	//atn(<float/integer>) - returns the arctangent, in radians, of x
	defineFunction(globalFunctionArray, "atn", Atn_execute, 1, true)

	//SQLITE FUNCTIONALITIES
	//sqlite_set(<string>) - Set file
	defineFunction(globalFunctionArray, "sqlite_set", Sqlite_set_execute, 1, true)

	//sqlite_q(<string>) - Query
	defineFunction(globalFunctionArray, "sqlite_q", Sqlite_q_execute, 1, true)

	//sqlite_cr() - Clear resources
	defineFunction(globalFunctionArray, "sqlite_cr", Sqlite_cr_execute, 0, true)

	//sqlite_fa(<string>) - Fetch all
	defineFunction(globalFunctionArray, "sqlite_fa", Sqlite_fa_execute, 1, true)

	//CRYPTOGRAPHIC FUNCTIONALITIES
	//m5(<string>) - md5
	defineFunction(globalFunctionArray, "m5", M5_execute, 1, true)

	//s1(<string>) - sha1
	defineFunction(globalFunctionArray, "s1", S1_execute, 1, true)

	//s256(<string>) - sha256
	defineFunction(globalFunctionArray, "s256", S256_execute, 1, true)

	//s512(<string>) - sha512
	defineFunction(globalFunctionArray, "s512", S512_execute, 1, true)

	//b64e(<string>) - base64 encode
	defineFunction(globalFunctionArray, "b64e", B64e_execute, 1, true)

	//b64d(<string>) - base64 decode
	defineFunction(globalFunctionArray, "b64d", B64d_execute, 1, true)

	//rsae(<string>, <string>) - RSA encrypt
	defineFunction(globalFunctionArray, "rsae", Rsae_execute, 2, true)

	//SOCKET FUNCTIONALITIES
	//netc(<string>, <string>) - socket connect
	defineFunction(globalFunctionArray, "netc", Netc_execute, 2, true)

	//netl(<string>, <string>) - socket listen
	defineFunction(globalFunctionArray, "netl", Netl_execute, 2, true)

	//netul(<string>, <string>, <integer>) - UDP socket listen
	defineFunction(globalFunctionArray, "netul", Netul_execute, 3, true)

	//netulf(<string>, <string>, <integer>, <string>) - UDP socket listen then call a function handler
	defineFunction(globalFunctionArray, "netulf", Netulf_execute, 4, true)

	//netla(<string>) - listener socket accept connection
	defineFunction(globalFunctionArray, "netla", Netla_execute, 1, true)

	//netlaf(<string>, <string>) - accept socket connection then call a function handler
	defineFunction(globalFunctionArray, "netlaf", Netlaf_execute, 2, true)

	//netlx(<string>) - listener socket close
	defineFunction(globalFunctionArray, "netlx", Netlx_execute, 1, true)

	//netx(<string>) - socket close
	defineFunction(globalFunctionArray, "netx", Netx_execute, 1, true)

	//netw(<string>, <string>) - socket write
	defineFunction(globalFunctionArray, "netw", Netw_execute, 2, true)

	//netr(<string>, <integer>) - socket read
	defineFunction(globalFunctionArray, "netr", Netr_execute, 2, true)

	//netur(<string>, <integer>) - UDP socket read when server is created using netul
	defineFunction(globalFunctionArray, "netur", Netur_execute, 2, true)

	//netus(<string>, <string>, <string>, <integer>) - send message via UDP
	defineFunction(globalFunctionArray, "netus", Netus_execute, 4, true)

	//FILE FUNCTIONALITIES
	//flrm(string)
	defineFunction(globalFunctionArray, "flrm", Flrm_execute, 1, true)

	//flmv(string, string)
	defineFunction(globalFunctionArray, "flmv", Flmv_execute, 2, true)

	//flcp(string, string)
	defineFunction(globalFunctionArray, "flcp", Flcp_execute, 2, true)

	//fo(string, string)
	defineFunction(globalFunctionArray, "fo", Fo_execute, 2, true)

	//fc(string)
	defineFunction(globalFunctionArray, "fc", Fc_execute, 1, true)

	//fw(string)
	defineFunction(globalFunctionArray, "fw", Fw_execute, 2, true)

	//fr(string)
	defineFunction(globalFunctionArray, "fr", Fr_execute, 2, true)

	//SDL FUNCTIONALITIES
	if SDL_ENABLED {
		//s_i(<integer>) - sdl init
		defineFunction(globalFunctionArray, "s_i", S_i_execute, 1, true)

		//s_q() - sdl quit
		defineFunction(globalFunctionArray, "s_q", S_q_execute, 0, true)

		//s_cw(<string>, <integer>, <integer>, <integer>, <integer>, <integer>) - sdl create window
		defineFunction(globalFunctionArray, "s_cw", S_cw_execute, 6, true)

		//s_dw(<string>) - sdl destroy window
		defineFunction(globalFunctionArray, "s_dw", S_dw_execute, 1, true)

		//s_usw(<string>) - sdl update surface window
		defineFunction(globalFunctionArray, "s_usw", S_usw_execute, 1, true)
	}
}
