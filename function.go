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

	//tof(<anyvar>)
	defineFunction(globalFunctionArray, "tof", Tof_execute, 1, true)

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

	//sgv(<string>, <anyvar>)
	defineFunction(globalFunctionArray, "sgv", Sgv_execute, 2, true)

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

		//s_gsw(<string>) - sdl get window surface
		defineFunction(globalFunctionArray, "s_gsw", S_gsw_execute, 1, true)

		//s_cr(<integer>, <integer>, <integer>, <integer>) - sdl create rectangle
		defineFunction(globalFunctionArray, "s_cr", S_cr_execute, 4, true)

		//s_gvr(<string>) - sdl get rectangle values
		defineFunction(globalFunctionArray, "s_gvr", S_gvr_execute, 1, true)

		//s_svr(<string>) - sdl set rectangle values
		defineFunction(globalFunctionArray, "s_svr", S_svr_execute, 5, true)

		//s_clr(<string>) - sdl clear rectangle
		defineFunction(globalFunctionArray, "s_clr", S_clr_execute, 1, true)

		//s_frsw(<string>, <string>, <integer>) - sdl fill rect
		defineFunction(globalFunctionArray, "s_frsw", S_frsw_execute, 3, true)

		//s_gdsw(<string>) - sdl get surface dimension
		defineFunction(globalFunctionArray, "s_gdsw", S_gdsw_execute, 1, true)

		//s_bsw(<string>, <string>, <string> ,<string>) - sdl surface blit
		defineFunction(globalFunctionArray, "s_bsw", S_bsw_execute, 4, true)

		//s_fsw(<string>) - sdl free surface
		defineFunction(globalFunctionArray, "s_fsw", S_fsw_execute, 1, true)

		//s_lbsw(<string>) - sdl load bmp
		defineFunction(globalFunctionArray, "s_lbsw", S_lbsw_execute, 1, true)

		//s_pe() - sdl poll event
		defineFunction(globalFunctionArray, "s_pe", S_pe_execute, 0, true)

		//s_ce(<string>) - sdl clear event
		defineFunction(globalFunctionArray, "s_ce", S_ce_execute, 1, true)

		//s_gte(<string>) - sdl get event type
		defineFunction(globalFunctionArray, "s_gte", S_gte_execute, 1, true)

		//s_kre(<string>) - sdl get event keyboard repeat
		defineFunction(globalFunctionArray, "s_kre", S_kre_execute, 1, true)

		//s_ksce(<string>) - sdl get event keyboard scan code
		defineFunction(globalFunctionArray, "s_ksce", S_ksce_execute, 1, true)

		//s_d() - sdl delay
		defineFunction(globalFunctionArray, "s_d", S_d_execute, 1, true)

		//s_it() - sdl ttf init
		defineFunction(globalFunctionArray, "s_it", S_it_execute, 0, true)

		//s_qt() - sdl ttf quit
		defineFunction(globalFunctionArray, "s_qt", S_qt_execute, 0, true)

		//s_oft(<string>, <integer>) - sdl ttf open font
		defineFunction(globalFunctionArray, "s_oft", S_oft_execute, 2, true)

		//s_cft(<string>) - sdl ttf close font
		defineFunction(globalFunctionArray, "s_cft", S_cft_execute, 1, true)

		//s_rft(<string>, <string>, <integer>, <integer>, <integer>, <integer>) - sdl ttf close font
		defineFunction(globalFunctionArray, "s_rft", S_rft_execute, 6, true)

		//s_mi(<integer>) - sdl mix init
		defineFunction(globalFunctionArray, "s_mi", S_mi_execute, 1, true)

		//s_mq() - sdl mix quit
		defineFunction(globalFunctionArray, "s_mq", S_mq_execute, 0, true)

		//s_moa(<integer>, <integer>, <integer>, <integer>) - sdl mix openaudio
		defineFunction(globalFunctionArray, "s_moa", S_moa_execute, 4, true)

		//s_mca() - sdl mix close audio
		defineFunction(globalFunctionArray, "s_mca", S_mca_execute, 0, true)

		//s_mlm(<string>) - sdl mix load music
		defineFunction(globalFunctionArray, "s_mlm", S_mlm_execute, 1, true)

		//s_mfm(<string>) - sdl mix free music
		defineFunction(globalFunctionArray, "s_mfm", S_mfm_execute, 1, true)

		//s_mpm(<string>, <integer>) - sdl mix play music
		defineFunction(globalFunctionArray, "s_mpm", S_mpm_execute, 2, true)

		//s_mlw(<string>) - sdl mix load wav
		defineFunction(globalFunctionArray, "s_mlw", S_mlw_execute, 1, true)

		//s_mfc(<string>) - sdl mix free chunk
		defineFunction(globalFunctionArray, "s_mfc", S_mfc_execute, 1, true)

		//s_mpc(<string>, <integer>, <integer>) - sdl mix play chunk
		defineFunction(globalFunctionArray, "s_mpc", S_mpc_execute, 3, true)

		//s_mhm(<string>) - sdl mix halt music
		defineFunction(globalFunctionArray, "s_mhm", S_mhm_execute, 0, true)

		//s_cre(<string>, <integer>, <integer>) - sdl create renderer
		defineFunction(globalFunctionArray, "s_cre", S_cre_execute, 3, true)

		//s_slsre(<string>, <integer>, <integer>) - sdl set logical size renderer
		defineFunction(globalFunctionArray, "s_slsre", S_slsre_execute, 3, true)

		//s_dre(<string>) - sdl destroy renderer
		defineFunction(globalFunctionArray, "s_dre", S_dre_execute, 1, true)

		//s_ctfsre(<string>, <string>) - sdl renderer create texture from surface
		defineFunction(globalFunctionArray, "s_ctfsre", S_ctfsre_execute, 2, true)

		//s_dt(<string>) - sdl destroy texture
		defineFunction(globalFunctionArray, "s_dt", S_dt_execute, 1, true)

		//s_clsre(<string>) - sdl clear renderer
		defineFunction(globalFunctionArray, "s_clsre", S_clsre_execute, 1, true)

		//s_pre(<string>) - sdl present renderer
		defineFunction(globalFunctionArray, "s_pre", S_pre_execute, 1, true)

		//s_cpre(<string>, <string>, <string>, <string>) - sdl copy renderer
		defineFunction(globalFunctionArray, "s_cpre", S_cpre_execute, 4, true)

		//s_li(<string>) - sdl load image
		defineFunction(globalFunctionArray, "s_li", S_li_execute, 1, true)
	}

	//RAYLIB FUNCTIONALITIES
	if RAYLIB_ENABLED {
		//rl_iw(<integer>, <integer>, <string>) - raylib init window
		defineFunction(globalFunctionArray, "rl_iw", Rl_iw_execute, 3, true)

		//rl_scw() - raylib window should close
		defineFunction(globalFunctionArray, "rl_scw", Rl_scw_execute, 0, true)

		//rl_cw() - raylib close window
		defineFunction(globalFunctionArray, "rl_cw", Rl_cw_execute, 0, true)

		//rl_bd() - raylib begin drawing
		defineFunction(globalFunctionArray, "rl_bd", Rl_bd_execute, 0, true)

		//rl_ed() - raylib end drawing
		defineFunction(globalFunctionArray, "rl_ed", Rl_ed_execute, 0, true)

		//rl_cb(<integer>) - raylib clear background
		defineFunction(globalFunctionArray, "rl_cb", Rl_cb_execute, 1, true)

		//rl_stf(<integer>) - raylib set target fps
		defineFunction(globalFunctionArray, "rl_stf", Rl_stf_execute, 1, true)

		//rl_li(<string>) - raylib load image
		defineFunction(globalFunctionArray, "rl_li", Rl_li_execute, 1, true)

		//rl_ui(<string>) - raylib unload image
		defineFunction(globalFunctionArray, "rl_ui", Rl_ui_execute, 1, true)

		//rl_dt(<string>, <integer>, <integer>, <integer>, <integer>) - raylib draw text
		defineFunction(globalFunctionArray, "rl_dt", Rl_dt_execute, 5, true)

		//rl_ltfi(<string>) - raylib load texture from image
		defineFunction(globalFunctionArray, "rl_ltfi", Rl_ltfi_execute, 1, true)

		//rl_ut(<string>) - raylib unload texture
		defineFunction(globalFunctionArray, "rl_ut", Rl_ut_execute, 1, true)

		//rl_gvt(<string>) - raylib get texture dimensions
		defineFunction(globalFunctionArray, "rl_gvt", Rl_gvt_execute, 1, true)

		//rl_drt(<string>, <integer>, <integer>, <integer>) - raylib draw texture
		defineFunction(globalFunctionArray, "rl_drt", Rl_drt_execute, 4, true)

		//rl_drtp(<string>, <float>, <float>, <float>, <float>, <float>, <float>, <float>, <float>, <float>, <float>, <float>, <integer>) - raylib draw texture pro
		defineFunction(globalFunctionArray, "rl_drtp", Rl_drtp_execute, 13, true)

		//rl_iad() - raylib init audio device
		defineFunction(globalFunctionArray, "rl_iad", Rl_iad_execute, 0, true)

		//rl_lms(<string>) - raylib load music stream
		defineFunction(globalFunctionArray, "rl_lms", Rl_lms_execute, 1, true)

		//rl_pms(<string>) - raylib play music stream
		defineFunction(globalFunctionArray, "rl_pms", Rl_pms_execute, 1, true)

		//rl_ums(<string>) - raylib update music stream
		defineFunction(globalFunctionArray, "rl_ums", Rl_ums_execute, 1, true)

		//rl_unms(<string>) - raylib unload music stream
		defineFunction(globalFunctionArray, "rl_unms", Rl_unms_execute, 1, true)

		//rl_lt(<string>) - raylib load texture
		defineFunction(globalFunctionArray, "rl_lt", Rl_lt_execute, 1, true)

		//rl_ikd(<integer>) - raylib is key down
		defineFunction(globalFunctionArray, "rl_ikd", Rl_ikd_execute, 1, true)

		//rl_lrt(<integer>, <integer>) - raylib load render texture
		defineFunction(globalFunctionArray, "rl_lrt", Rl_lrt_execute, 2, true)

		//rl_urt(<string>) - raylib unload render texture
		defineFunction(globalFunctionArray, "rl_urt", Rl_urt_execute, 1, true)

		//rl_btm(<string>) - raylib begin texture mode
		defineFunction(globalFunctionArray, "rl_btm", Rl_btm_execute, 1, true)

		//rl_etm() - raylib end texture mode
		defineFunction(globalFunctionArray, "rl_etm", Rl_etm_execute, 0, true)

		//rl_gtfrt(<string>) - raylib get texture from render texture
		defineFunction(globalFunctionArray, "rl_gtfrt", Rl_gtfrt_execute, 1, true)

	}
}
