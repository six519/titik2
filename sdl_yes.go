//go:build sdl
// +build sdl

package main

import (
	"errors"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"strconv"
)

const SDL_ENABLED bool = true

var SDL_INIT_TYPES = []uint32{
	sdl.INIT_EVERYTHING,
	sdl.INIT_TIMER,
	sdl.INIT_AUDIO,
	sdl.INIT_VIDEO,
	sdl.INIT_JOYSTICK,
	sdl.INIT_HAPTIC,
	sdl.INIT_GAMECONTROLLER,
	sdl.INIT_EVENTS,
	sdl.INIT_NOPARACHUTE,
	sdl.INIT_SENSOR,
}

var SDL_WPOSITIONS = []int32{
	sdl.WINDOWPOS_UNDEFINED,
	sdl.WINDOWPOS_CENTERED,
}

var SDL_WFLAGS = []uint32{
	sdl.WINDOW_FULLSCREEN,
	sdl.WINDOW_SHOWN,
	sdl.WINDOW_HIDDEN,
	sdl.WINDOW_BORDERLESS,
	sdl.WINDOW_RESIZABLE,
	sdl.WINDOW_MINIMIZED,
	sdl.WINDOW_MAXIMIZED,
	sdl.WINDOW_FULLSCREEN_DESKTOP,
	sdl.WINDOW_ALWAYS_ON_TOP,
	sdl.WINDOW_TOOLTIP,
	sdl.WINDOW_POPUP_MENU,
}

var SDL_EVENT_TYPES = map[uint32]int{
	sdl.FIRSTEVENT:   0,
	sdl.QUIT:         1,
	sdl.DISPLAYEVENT: 2,
	sdl.WINDOWEVENT:  3,
	sdl.SYSWMEVENT:   4,
}

var SDL_MIX_INIT_TYPES = []int{
	mix.INIT_FLAC,
	mix.INIT_MOD,
	mix.INIT_MP3,
	mix.INIT_OGG,
}

var SDL_MIX_DEFAULTS = []int{
	mix.DEFAULT_FREQUENCY,
	mix.DEFAULT_FORMAT,
	mix.DEFAULT_CHANNELS,
	mix.DEFAULT_CHUNKSIZE,
}

var SDL_RENDERER_FLAGS = []uint32{
	sdl.RENDERER_SOFTWARE,
	sdl.RENDERER_ACCELERATED,
	sdl.RENDERER_PRESENTVSYNC,
	sdl.RENDERER_TARGETTEXTURE,
}

func S_i_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: true}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if arguments[0].IntegerValue < 0 || (arguments[0].IntegerValue-1) > len(SDL_INIT_TYPES) {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := sdl.Init(SDL_INIT_TYPES[arguments[0].IntegerValue])

			if err != nil {
				ret.BooleanValue = false
			}
		}
	}

	return ret
}

func S_q_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	sdl.Quit()
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_mq_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	mix.Quit()
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_cw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 5, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 4, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if arguments[4].IntegerValue < 0 || (arguments[4].IntegerValue-1) > len(SDL_WPOSITIONS) ||
			arguments[3].IntegerValue < 0 || (arguments[3].IntegerValue-1) > len(SDL_WPOSITIONS) ||
			arguments[0].IntegerValue < 0 || (arguments[0].IntegerValue-1) > len(SDL_WFLAGS) {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {

			window, err := sdl.CreateWindow(arguments[5].StringValue, SDL_WPOSITIONS[arguments[4].IntegerValue], SDL_WPOSITIONS[arguments[3].IntegerValue], int32(arguments[2].IntegerValue), int32(arguments[1].IntegerValue), SDL_WFLAGS[arguments[0].IntegerValue])

			if err == nil {
				window_reference := "win_" + generateRandomNumbers()
				(*globalSettings).sdlWindow[window_reference] = window
				ret.Type = RET_TYPE_STRING
				ret.StringValue = window_reference
			}
		}

	}

	return ret
}

func S_dw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).sdlWindow[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized window on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlWindow[arguments[0].StringValue].Destroy()
			delete((*globalSettings).sdlWindow, arguments[0].StringValue)
		}
	}

	return ret
}

func S_usw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).sdlWindow[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized window on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlWindow[arguments[0].StringValue].UpdateSurface()
		}
	}

	return ret
}

func S_gsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).sdlWindow[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized window on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			surface, err := (*globalSettings).sdlWindow[arguments[0].StringValue].GetSurface()

			if err == nil {
				surface_reference := "surf_" + generateRandomNumbers()
				(*globalSettings).sdlSurface[surface_reference] = surface
				ret.Type = RET_TYPE_STRING
				ret.StringValue = surface_reference
			}
		}
	}

	return ret
}

func S_cr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		rect := sdl.Rect{int32(arguments[3].IntegerValue), int32(arguments[2].IntegerValue), int32(arguments[1].IntegerValue), int32(arguments[0].IntegerValue)}

		rect_reference := "rect_" + generateRandomNumbers()
		(*globalSettings).sdlRect[rect_reference] = rect
		ret.Type = RET_TYPE_STRING
		ret.StringValue = rect_reference
	}

	return ret
}

func S_frsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if (*globalSettings).sdlSurface[arguments[2].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized surface on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			if arguments[1].StringValue == "" {
				(*globalSettings).sdlSurface[arguments[2].StringValue].FillRect(nil, uint32(arguments[0].IntegerValue))
			} else {
				if val, ok := (*globalSettings).sdlRect[arguments[1].StringValue]; ok {
					(*globalSettings).sdlSurface[arguments[2].StringValue].FillRect(&val, uint32(arguments[0].IntegerValue))
				} else {
					*errMessage = errors.New("Error: Uninitialized rectangle on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
				}
			}
		}

	}

	return ret
}

func S_bsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {

		if (*globalSettings).sdlSurface[arguments[3].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized surface on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {

			ok_to_execute := true
			ok := false
			var val1 sdl.Rect
			var val2 sdl.Rect

			if arguments[2].StringValue != "" {
				if val1, ok = (*globalSettings).sdlRect[arguments[2].StringValue]; !ok {
					ok_to_execute = false
					*errMessage = errors.New("Error: Uninitialized rectangle on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
				}
			}

			if (*globalSettings).sdlSurface[arguments[1].StringValue] == nil {
				ok_to_execute = false
				*errMessage = errors.New("Error: Uninitialized surface on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}

			if val2, ok = (*globalSettings).sdlRect[arguments[0].StringValue]; !ok {
				ok_to_execute = false
				*errMessage = errors.New("Error: Uninitialized rectangle on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}

			if ok_to_execute {
				if arguments[2].StringValue == "" {
					(*globalSettings).sdlSurface[arguments[3].StringValue].Blit(nil, (*globalSettings).sdlSurface[arguments[1].StringValue], &val2)
				} else {
					(*globalSettings).sdlSurface[arguments[3].StringValue].Blit(&val1, (*globalSettings).sdlSurface[arguments[1].StringValue], &val2)
				}
			}

		}

	}

	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_gdsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ASSOCIATIVE_ARRAY}
	ret.AssociativeArrayValue = make(map[string]FunctionReturn)

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {

		if (*globalSettings).sdlSurface[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized surface on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			ret.AssociativeArrayValue["width"] = FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: int((*globalSettings).sdlSurface[arguments[0].StringValue].W)}
			ret.AssociativeArrayValue["height"] = FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: int((*globalSettings).sdlSurface[arguments[0].StringValue].H)}
		}

	}

	return ret
}

func S_fsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {

		if (*globalSettings).sdlSurface[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized surface on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlSurface[arguments[0].StringValue].Free()
			delete((*globalSettings).sdlSurface, arguments[0].StringValue)
		}

	}

	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_pe_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}
	event := sdl.PollEvent()

	evt_reference := "evt_" + generateRandomNumbers()
	(*globalSettings).sdlEvent[evt_reference] = event
	ret.StringValue = evt_reference

	return ret
}

func S_ce_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if _, ok := (*globalSettings).sdlEvent[arguments[0].StringValue]; ok {
			(*globalSettings).sdlEvent[arguments[0].StringValue] = nil
			delete((*globalSettings).sdlEvent, arguments[0].StringValue)
		}
	}

	return ret
}

func S_gte_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: -1}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {

		if val, ok := (*globalSettings).sdlEvent[arguments[0].StringValue]; ok {
			if val != nil {
				if val2, ok2 := SDL_EVENT_TYPES[val.GetType()]; ok2 {
					ret.IntegerValue = val2
				}
			}
		}

	}

	return ret
}

func S_d_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		sdl.Delay(uint32(arguments[0].IntegerValue))
	}

	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_it_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: true}

	err := ttf.Init()

	if err != nil {
		ret.BooleanValue = false
	}

	return ret
}

func S_qt_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ttf.Quit()
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_oft_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		font, err := ttf.OpenFont(arguments[1].StringValue, arguments[0].IntegerValue)
		if err == nil {
			font_reference := "fnt_" + generateRandomNumbers()
			(*globalSettings).sdlFont[font_reference] = font
			ret.Type = RET_TYPE_STRING
			ret.StringValue = font_reference
		}
	}

	return ret
}

func S_cft_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).sdlFont[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized font on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlFont[arguments[0].StringValue].Close()
			delete((*globalSettings).sdlFont, arguments[0].StringValue)
		}
	}

	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_rft_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 5, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 4, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if (*globalSettings).sdlFont[arguments[5].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized font on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {

			surface, err := (*globalSettings).sdlFont[arguments[5].StringValue].RenderUTF8Blended(arguments[4].StringValue, sdl.Color{R: uint8(arguments[3].IntegerValue), G: uint8(arguments[2].IntegerValue), B: uint8(arguments[1].IntegerValue), A: uint8(arguments[0].IntegerValue)})

			if err == nil {
				surface_reference := "surf_" + generateRandomNumbers()
				(*globalSettings).sdlSurface[surface_reference] = surface
				ret.Type = RET_TYPE_STRING
				ret.StringValue = surface_reference
			}

		}

	}

	return ret
}

func S_mi_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: true}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if arguments[0].IntegerValue < 0 || (arguments[0].IntegerValue-1) > len(SDL_MIX_INIT_TYPES) {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := mix.Init(SDL_MIX_INIT_TYPES[arguments[0].IntegerValue])

			if err != nil {
				ret.BooleanValue = false
			}
		}
	}

	return ret
}

func S_moa_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: true}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if arguments[3].IntegerValue < 0 || (arguments[3].IntegerValue-1) > len(SDL_MIX_DEFAULTS) || arguments[2].IntegerValue < 0 || (arguments[2].IntegerValue-1) > len(SDL_MIX_DEFAULTS) ||
			arguments[1].IntegerValue < 0 || (arguments[1].IntegerValue-1) > len(SDL_MIX_DEFAULTS) || arguments[0].IntegerValue < 0 || (arguments[0].IntegerValue-1) > len(SDL_MIX_DEFAULTS) {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := mix.OpenAudio(SDL_MIX_DEFAULTS[arguments[3].IntegerValue], uint16(SDL_MIX_DEFAULTS[arguments[2].IntegerValue]), SDL_MIX_DEFAULTS[arguments[1].IntegerValue], SDL_MIX_DEFAULTS[arguments[0].IntegerValue])

			if err != nil {
				ret.BooleanValue = false
			}
		}
	}

	return ret
}

func S_mca_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	mix.CloseAudio()
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_mlm_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {

		music, err := mix.LoadMUS(arguments[0].StringValue)
		if err == nil {
			music_reference := "msc_" + generateRandomNumbers()
			(*globalSettings).sdlMusic[music_reference] = music
			ret.Type = RET_TYPE_STRING
			ret.StringValue = music_reference
		}
	}

	return ret
}

func S_mfm_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {

		if (*globalSettings).sdlMusic[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized music on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlMusic[arguments[0].StringValue].Free()
			delete((*globalSettings).sdlMusic, arguments[0].StringValue)
		}

	}
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_mpm_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: true}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if (*globalSettings).sdlMusic[arguments[1].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized music on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := (*globalSettings).sdlMusic[arguments[1].StringValue].Play(arguments[0].IntegerValue)

			if err != nil {
				ret.BooleanValue = false
			}
		}
	}

	return ret
}

func S_mhm_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	mix.HaltMusic()
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_mlw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {

		wav, err := mix.LoadWAV(arguments[0].StringValue)
		if err == nil {
			wav_reference := "wav_" + generateRandomNumbers()
			(*globalSettings).sdlChunk[wav_reference] = wav
			ret.Type = RET_TYPE_STRING
			ret.StringValue = wav_reference
		}
	}

	return ret
}

func S_mfc_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {

		if (*globalSettings).sdlChunk[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized chunk on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlChunk[arguments[0].StringValue].Free()
			delete((*globalSettings).sdlChunk, arguments[0].StringValue)
		}

	}
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_mpc_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: true}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if (*globalSettings).sdlChunk[arguments[2].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized chunk on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			_, err := (*globalSettings).sdlChunk[arguments[2].StringValue].Play(arguments[1].IntegerValue, arguments[0].IntegerValue)

			if err != nil {
				ret.BooleanValue = false
			}
		}
	}

	return ret
}

func S_cre_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if (*globalSettings).sdlWindow[arguments[2].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized window on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {

			if arguments[0].IntegerValue < 0 || (arguments[0].IntegerValue-1) > len(SDL_RENDERER_FLAGS) {
				*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {

				ren, err := sdl.CreateRenderer((*globalSettings).sdlWindow[arguments[2].StringValue], arguments[1].IntegerValue, SDL_RENDERER_FLAGS[arguments[0].IntegerValue])
				if err == nil {
					ren_reference := "ren_" + generateRandomNumbers()
					(*globalSettings).sdlRenderer[ren_reference] = ren
					ret.Type = RET_TYPE_STRING
					ret.StringValue = ren_reference
				}

			}
		}

	}

	return ret
}

func S_dre_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).sdlRenderer[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized renderer on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlRenderer[arguments[0].StringValue].Destroy()
			delete((*globalSettings).sdlRenderer, arguments[0].StringValue)
		}
	}
	return FunctionReturn{Type: RET_TYPE_NONE}
}
