package juusmenu

/*

juusmenu, golang: menu maker for go programs running in terminal.
Supports unlimited SubMenu nesting and runtime menu manipulation.

Copyright (C) November 2020 Julius Heinrich Ludwig Schön / Ronald Michael Spicer
Foto.TimePirate.org / TimePirate.org / PaganToday.TimePirate.org
Julius@TimePirate.org

This code is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, version 3 of the License.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You may have received a copy of the GNU General Public License
along with this code.  If not, see <http://www.gnu.org/licenses/>.


There is a lot of validation in this unit. The validations was written
for the express purpose to help anyone using this package during
the development phase of your project. For the most part, adding
menus is elementary.

But. There is the ability to dynamically at run-time manipulate menus
in almost any way you can imagine.

Algorithms are not as smart as programmers <ref?> and so the ability to
mess up the Menus relies on algorithmic bugs (or sleepy programmers).

Use fmt.Println(<unit>.MenuOptions.InfoMenuOptions()) for an explanation
of settings that can help you debug dynamic menus.

*/

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func init() {
	if menuScanner == nil {
		menuScanner = bufio.NewScanner(os.Stdin)
	}

	//No method to change aligner alignment-side during run-time, odd.
	//So have to make two...
	if alignerLeft == nil {
		alignerLeft = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	}
	if alignerRight == nil {
		alignerRight = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	}

	alignerLocal = alignerLeft

	MenuOptions = &menuOptions{
		funcBracketTop:         funcBracketTopStr,
		funcBracketBottom:      funcBracketBottomStr,
		idFuncRunner:           true,
		killPhrase:             killPhraseStr,
		menuPrompt:             menuPromptStr,
		menuSeparator:          menuSeparatorStr,
		pauseOnOutput:          true,
		runTimeErrMsgs_Display: true,
		runTimeErrMsgs_Pause:   true,
	}
}

type (
	menuEntries map[string]*menuEntry
	validateKey map[string]int
	menuList    []*Menu
)

var ( //many of these are "constants", but there is a need to pass &addresses
	MenuOptions *menuOptions
	//private
	additionalString          = "\n^^^ Important information above, please read..."
	alignerLeft               *tabwriter.Writer
	alignerLocal              *tabwriter.Writer
	alignerRight              *tabwriter.Writer
	allMenus                  menuList
	defBreakh                 = "Quit this Menu"
	defBreakv                 = "QQ.QQ"
	defidFuncRunner           = true
	defpauseOnOutput          = true
	defrunTimeErrMsgs_Display = true
	defrunTimeErrMsgs_Pause   = true
	emptyHint                 = "Menu hint not specified"
	emptyString               = ""
	funcBracketBottomStr      = "..............*"
	funcBracketTopStr         = "*.............."
	killPhraseStr             = "Bye!"
	killSwitch                = false
	menuID                    = 0
	menuPromptStr             = ">>: "
	menuScanner               *bufio.Scanner //all Menu's will use this
	menuSeparatorStr          = ":"
	unNamedMenuTitle          = "UnNamedMenu"
)

const (
	warn = "\n--- Attention --------------------\n"
	//breakIndicator: menuEntries "value" that self-manages by carrying the loop "break" key
	breakIndicator = "^,^BrEaK^*^" //unlikely to be typed as a map Key by anybody
	killTemplate   = "===============  '%s' immediately exits all Menus  ========="
	menuFormat     = "%s\t: \t%s"
	trimString     = " \t\r\n"
)

//dropDownStruct : for managing Menu Entries that Start() an already open Menu
type dropDownStruct struct {
	doDropDown, isLastExit bool
	id                     int
}

var dropDown = &dropDownStruct{
	doDropDown: false,
	isLastExit: false,
	id:         0,
}

//trying out various ways to have an enumerated type
type bracketSwitch int

const (
	bsTop bracketSwitch = iota
	bsBottom
	bsPartial
	bsCOUNT //handy, gives a count to use for iterating
)

//menuOptions : Settings to manage menu behaviour. Run
//fmt.Println(<unit>.MenuOptions.InfoMenuOptions()), or
//see const's at bottom of unit, for an explanation
//of the meaning of these fields
type menuOptions struct {
	funcBracketTop         string
	funcBracketBottom      string
	killPhrase             string
	idFuncRunner           bool
	menuPrompt             string
	menuSeparator          string
	pauseOnOutput          bool
	runTimeErrMsgs_Display bool
	runTimeErrMsgs_Pause   bool
}

//Stringer for menuOptions
func (mo *menuOptions) String() string {
	const f = "%s:\n  current value = '%v' (default is: '%v')"
	return "Menu Options" + "\n" +
		fmt.Sprintf(f, "idFuncRunner", mo.idFuncRunner, defidFuncRunner) + "\n" +
		fmt.Sprintf(f, "killPhrase", mo.killPhrase, killPhraseStr) + "\n" +
		fmt.Sprintf(f, "menuPrompt", mo.menuPrompt, menuPromptStr) + "\n" +
		fmt.Sprintf(f, "menuSeparator", mo.menuSeparator, menuSeparatorStr) + "\n" +
		fmt.Sprintf(f, "pauseOnOutput", mo.pauseOnOutput, defpauseOnOutput) + "\n" +
		fmt.Sprintf(f, "runTimeErrMsgs_Display", mo.runTimeErrMsgs_Display, defrunTimeErrMsgs_Display) + "\n" +
		fmt.Sprintf(f, "runTimeErrMsgs_Pause", mo.runTimeErrMsgs_Pause, defrunTimeErrMsgs_Pause) + "\n" +
		fmt.Sprintf(f, "funcBracketTop", mo.funcBracketTop, funcBracketTopStr) + "\n" +
		fmt.Sprintf(f, "funcBracketBottom", mo.funcBracketBottom, funcBracketBottomStr) + "\n" +
		fmt.Sprint("\n>> execute fmt.Println(<unit>.MenuOptions.InfoMenuOptions()) for field information") + "\n\n"
}

//InfoMenuOptions : Displays explanations for the menuOptions fields
func (mo *menuOptions) InfoMenuOptions() string {
	return "menuOptions fields:" + "\n\n" +
		idFuncRunnerInfo + "\n\n" +
		killPhraseInfo + "\n\n" +
		menuPromptInfo + "\n\n" +
		menuSeparatorInfo + "\n\n" +
		pauseOnOutputInfo + "\n\n" +
		runTimeErrMsgs_DisplayInfo + "\n\n" +
		runTimeErrMsgs_PauseInfo + "\n\n" +
		funcBracketTopInfo + "\n\n" +
		funcBracketBottomInfo + "\n\n" +
		fmt.Sprint(">> execute fmt.Println(<unit>.MenuOptions()) for field values") + "\n\n"
}

//AlignRight : align the Menu values to the right
func (mo *menuOptions) AlignRight() {
	alignerLocal = alignerRight
}

//AlignLeft : align the Menu values to the left
func (mo *menuOptions) AlignLeft() {
	alignerLocal = alignerLeft
}

//SetMenuPrompt : Set the user prompt. Even " " is allowed, but "" uses default value
func (mo *menuOptions) SetMenuPrompt(val string) {
	//don't want space trimmed prompts, allows " "
	var Val string
	if Val = strings.Trim(val, "\t\n\r"); Val != "" {
		Val = menuPromptStr
	}
	mo.menuPrompt = Val
}

//SetMenuSeparator : Set the separator between Menu Title breadcrumbs,
//only works if submenus are declared and added via Menu.AddSubMenu(xxx)
func (mo *menuOptions) SetMenuSeparator(val string) {
	valueClean(&val, &menuSeparatorStr, vIgnore, func() {})
	mo.menuSeparator = val
}

//SetfuncBracketTop : set function bracket entry indicator string
//this string will be printed before the function runs.
func (mo *menuOptions) SetfuncBracketTop(val string) {
	valueClean(&val, &emptyString, vIgnore, func() {})
	mo.funcBracketTop = val
}

//SetfuncBracketBottom : set function bracket exit indicator string
//this string will be printed after the function returns.
func (mo *menuOptions) SetfuncBracketBottom(val string) {
	valueClean(&val, &emptyString, vIgnore, func() {})
	mo.funcBracketBottom = val
}

//SetKillPhrase : Set the kill phrase. This is used to exit the menu system at
//any time, no matter how deeply nested. Empty string "" disables this ability.
func (mo *menuOptions) SetKillPhrase(val string) {
	valueClean(&val, &killPhraseStr, vIgnore, func() {})
	mo.killPhrase = val
}

//SetPauseOnOutput : If true user will be asked to press return after the Menu Item's
//associated func() returns.
func (mo *menuOptions) SetPauseOnOutput(val bool) {
	mo.pauseOnOutput = val
}

//SetIdFuncRunner : If true function output will be bracketed by the calling Menu Title
//and the Key the user typed. Useful for debugging, but nice to have
//in general also maybe.
func (mo *menuOptions) SetIdFuncRunner(val bool) {
	mo.idFuncRunner = val
}

//SetRunTimeErrMsgs_Display : if true error conditions will print at run time.
//Intended for the devlopment phase of your project, because no particular need
//to check error results, they will always be printed. If your project is finished
//I, personally, would still leave this on. It may be of use to your users too.
//If you set this to false, then you will need to handle errors and decide if you
//want them displayed. There are, really, no fatal errors in system except for
//Menu.Start() and menuSystem.StartMenuSystem()
func (mo *menuOptions) SetRunTimeErrMsgs_Display(val bool) {
	mo.runTimeErrMsgs_Display = val
}

//SetRunTimeErrMsgs_Pause : if true error conditions printed at run time will pause
//for user acknowledgement. Intended for the devlopment phase of your project.
func (mo *menuOptions) SetRunTimeErrMsgs_Pause(val bool) {
	mo.runTimeErrMsgs_Pause = val
}

type Menu struct {
	Title string //use this or GetID for use in switch statements
	//private
	//skipFunctionNotification : if a Menu Entry func() starts a menu directly with <menuvar>.Start()
	//the menu is treated as a function and not a submenu (one should be Using
	//<menuvar>.AddSubMenu to get that functionality). This means that the menu
	//system will pause after the menu is quit, depending on the MenuOptions
	//setting. This interferes with the expected smooth flow from menu to menu.
	//The skipFunctionNotification var (set through func SkipFunctionNofication)
	//bypasses this possible interruption. It is niche use and once turned on
	//it will be reset to false after the exit from the function, or the menu is
	//normally closed. No need to use if isChooseOne is set, those menus auto-manage.
	skipFunctionNotification bool
	entries                  menuEntries
	finalized                bool
	//id : each NewMenu is automatically numbered starting at -1 and
	//decreasing. One can call GetID() to retrieve them or use
	//SetID() to set them to values of your own choice.
	//Careful: use positive values if you set id's, so that there
	//is no accidental id collision.
	id int
	//isChooseOne : set this option of an individual Menu with SetChooseOne()
	//causes any single value to exit the menu scan loop. Meant for menus
	//that are intended as "yes/no/maybe/cancel" types.
	//The menu MUST still have a break indicator, which is then a good choice for
	//the "cancel" option. func() can, of course, be added to break indicators too.
	//It will function to end the menu if isChooseOne value is dynamically re-set
	isChooseOne, isMainMenu bool
	isModified, isRunning   bool
	killThisMenu            bool
	//quitValue : super important. The key to be used to quit the menu.
	quitValue   string
	parent      *Menu
	reverseSort bool
	//sortkeys - quote:	When iterating over a map with a range loop, the iteration order is not specified
	//and is not guaranteed to be the same from one iteration to the next.
	//For a stable iteration order one must maintain a separate data structure that specifies that order.
	sortKeys []string
	//performs validations on menu's Keys. Mostly to warn programmer
	//that duplicate Keys were sent in and the menu may not function as designed.
	validateKeys validateKey
}

//Stringers for the Menu
func (menu *Menu) String() string {
	const (
		f  = "  %s : %v"
		ef = "  '%s'=%s\n"
	)
	parentName := "<no parent>"
	if menu.parent != nil {
		parentName = fmt.Sprintf("'%s', ID: %d", menu.parent.Title, menu.parent.id)
	}
	unsortedEntries := ">> menu entries (unsorted):\n"
	for _, e := range menu.entries {
		unsortedEntries = unsortedEntries + fmt.Sprintf(ef, e.value, e.hint)
	}

	return fmt.Sprintf("Menu '%s':", menu.Title) + "\n" +
		fmt.Sprintf(f, "ID", menu.id) + "\n" +
		fmt.Sprintf(f, "Break Value", "'"+menu.quitValue+"'") + "\n" +
		fmt.Sprintf(f, "ChooseOne Menu", menu.isChooseOne) + "\n" +
		fmt.Sprintf(f, "Sort Ascending", menu.reverseSort) + "\n" +
		fmt.Sprintf(f, "Parent Menu", parentName) + "\n" +
		unsortedEntries
}

//SetChooseOne : default false. If true any valid Key runs the associated func()
//and then quits the menu. Usefule for "yes/no/maybe/cancel" menus.
func (menu *Menu) SetChooseOne(val bool) {
	menu.isChooseOne = val
}

//SortAscending : Default. Menu entries will be sorted in ascending order
func (menu *Menu) SortAscending() {
	menu.reverseSort = false
}

//SortDescending : Menu entries will be sorted in descending order
func (menu *Menu) SortDescending() {
	menu.reverseSort = true
}

//SkipFunctionNotification : Niche use, but very useful when needed. if a Menu Entry func()
//starts a menu directly with <menuvar>.Start() this function turns
//on a func() pause bypass to ensure smooth menu flow. If SetChooseOne
//is used on a menu, there is no need to use this function as those
//menus are auto-managed.
func (menu *Menu) SkipFunctionNotification() {
	menu.skipFunctionNotification = true
}

//getmenuID : internal use to assign default ID's to New Menu's
func getmenuID() int {
	menuID--
	return menuID
}

//NewMenu : Returns a bright and shiny new *Menu
func NewMenu(title string) (result *Menu) {
	valueClean(&title, &unNamedMenuTitle, vIgnore, func() {})
	result = &Menu{
		Title: title,
		//private
		skipFunctionNotification: false,
		entries:                  make(menuEntries),
		finalized:                false,
		id:                       getmenuID(),
		isChooseOne:              false,
		isMainMenu:               false,
		isModified:               false,
		isRunning:                false,
		killThisMenu:             false,
		parent:                   nil,
		quitValue:                "",
		reverseSort:              false,
		validateKeys:             make(validateKey),
	}
	if len(allMenus) == 0 {
		result.isMainMenu = true
	}
	allMenus = append(allMenus, result)
	return
}

//menuEntry : holds the Menu Key, description, and entry's func()
type menuEntry struct {
	value string
	hint  string
	//isSubMenuEntry : internal use for indicating submenus
	isSubMenuEntry bool
	//doRun : any func can be put in here, simply type the func()
	//in the body. Obviates the need for varieties of func declarations.
	doRun func()
}

//makeMenuEntry : Internal use only, return a *menuEntry to be attached to a Menu.
//Used to add the Menu breakIndicator in a special way
func makeMenuEntry(aval string, ahint string, afunc func()) *menuEntry {
	//internal function, param validation should be done before calling this func()
	return &menuEntry{
		value:          aval,
		hint:           ahint,
		isSubMenuEntry: false,
		doRun:          afunc,
	}
}

//vClean : testing "enumerated types"; allowing placing of func() running in func valueClean
type vClean uint

const (
	vIgnore vClean = iota
	vInBlock
	vAfterBlock
	vCOUNT //handy, gives a count to use for iterating
)

func (vc vClean) String() string {
	//echo constant name when printed as %s/%v
	//After working with this I wonder if there isn't a better way with a map?
	return [...]string{"vIgnore", "vInBlock", "vAfterBlock"}[vc]
}

//vClean.Display : playing with possible enumerated types, internal use.
func (vc vClean) Display() {
	fmt.Println(fmt.Sprintf("Enumerated type: %T", vc))
	for i := 0; i < int(vCOUNT); i++ {
		fmt.Println(fmt.Sprintf(" %d \t %v", vClean(i), vClean(i)))
	}
	fmt.Printf("Current value: %d -- %v\n", vc, vc)
}

//MenuSystemInterface : creates a generic publication of abstract "MenuSystem" func()'s'
type MenuSystemInterface interface {
	StartMenuSystem() error
	WasKilled() bool
	Unkill()
}

//menuSystem : a place to park a couple of system funcs, is type MenuSystemInterface
type menuSystem struct{}

var MenuSystem menuSystem

//WasKilled : Adds ability to check if the user killed the menu system using the killPhrase
func (ms menuSystem) WasKilled() bool {
	return killSwitch
}

//UnKill : Re-Sets a killed menu system, of use only if you want to query the
//user for confirmation of closing the menu system.
func (ms menuSystem) UnKill() {
	killSwitch = false
}

//StartMenuSystem : One of the ways to start the menu system. One can also simply use <menuvar>.Start()
//Essentially this func() makes it easier for a novice programmer to use.
func (ms menuSystem) StartMenuSystem() error {
	idx := -1
	for i := 0; i < len(allMenus); i++ {
		if allMenus[i].isMainMenu {
			idx = i
			break
		}
	}
	if idx < 0 {
		panic(warn + fmt.Sprint("StartMenuSystem: Can not Start Menu System because setting MAIN MENU failed!"))
	}
	return allMenus[idx].Start()
}

//setRunning : internal function used to control and manipulate Menus
func (menu *Menu) setRunning(val bool) {
	menu.isRunning = val
	menu.skipFunctionNotification = false
}

//SetID : Set a Menu's id. By default each NewMenu() gets a negative id.
//You can set a Menu to your own id's with this function.
func (menu *Menu) SetID(id int) error {
	errmsg := ""
	if menu.isRunning {
		errmsg = warn + fmt.Sprintf("SetID(): Menu '%s', ID '%d' (requesting ID '%d') is currently in scan loop, not allowed to change its id.", menu.Title, menu.id, id)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	for i := 0; i < len(allMenus); i++ {
		if allMenus[i].id == id {
			errmsg = warn + "SetID: Menu '%s' requested id '%d' but it is already assigned to Menu " +
				fmt.Sprintf("%s, menu.id not changed", allMenus[i].Title)
			break
		}
	}

	if errmsg != "" {
		errmsg = fmt.Sprintf(errmsg, menu.Title, id)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}
	menu.id = id
	return nil
}

//GetID : Function that retrieves a Menu's id.
func (menu *Menu) GetID() int {
	return menu.id
}

//AddSubMenu : val is the string to type to run the entry. This ties a *Menu to
//another *Menu as a sub-menu. Using this then the menu system
//can "breadcrumb" your menu nesting. If you run submenus yourself through a func()
//you will lose the breadcrumbing. This also performs a lot of validations.
//recommended to use this for subMenus.
func (menu *Menu) AddSubMenu(subMenu *Menu, val string, hint string) error {
	const (
		methodName = "AddSubMenu method: "
		endMsg     = "SubMenu not added.\n"
	)
	var errmsg string

	if subMenu == nil {
		errmsg = warn + fmt.Sprintf(methodName+"subMenu was nil for val '%s', hint '%s'. %s", val, hint, endMsg)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	if subMenu.isMainMenu {
		errmsg = warn + fmt.Sprintf(methodName+"Menu '%s' requested subMenu '%s' which is the MAIN MENU. Not allowed.", menu.Title, subMenu.Title)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	if menu == subMenu {
		errmsg = warn + fmt.Sprintf(methodName+"Menu '%s' & SubMenu '%s' are the same, can not add a menu onto itself.", menu.Title, subMenu.Title)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	//prevent circular references of submenus
	if subMenu.parent != nil {
		errmsg = warn + fmt.Sprintf(methodName+"subMenu '%s' is already assigned to menu '%s', ignored.", subMenu.Title, subMenu.parent.Title)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	valueClean(&val, &emptyString, vInBlock, func() {
		errmsg = warn + fmt.Sprintf(methodName+"Empty Entry Value sent in. Menu '%s', SubMenu '%s'. %s", menu.Title, subMenu.Title, endMsg)
		alertUser(&errmsg)
	})
	if val == "" {
		return errors.New(errmsg)
	}
	valueClean(&hint, &emptyHint, vIgnore, func() {})

	//the func() telling the menu to start is what makes this a submenu
	menu.AddMenuEntry(val, hint, func() { subMenu.Start() })
	menu.entries[val].isSubMenuEntry = true

	subMenu.parent = menu
	//a Menu is "modified" only if has already been Start()'d and Finalized'
	menu.isModified = menu.finalized
	return nil
}

//SetMenuBreakItem : Quite important, sets the menu value string
//that the user will type to quit the  menu's scan loop
func (menu *Menu) SetMenuBreakItem(val string, hint string, afunc func()) error {
	//breakIndicator is special, the Start method can figure out what
	//string will break the infinite Scan loop
	errmsg := warn + fmt.Sprintf("SetMenuBreakItem method: Menu '%s', Empty Break Value sent in, Break value set to '%s'.\n", menu.Title, defBreakv)

	valueClean(&val, &defBreakv, vInBlock, func() {
		alertUser(&errmsg)
	})

	valueClean(&hint, &defBreakh, vIgnore, func() {})

	if _, ok := menu.entries[breakIndicator]; ok {
		//they've accidently sent in the breakindicator again, or it is dynamically changed,
		//Take the newest values
		menu.entries[breakIndicator].value = val
		menu.entries[breakIndicator].hint = hint
		menu.entries[breakIndicator].doRun = afunc
		menu.isModified = menu.finalized
	} else {
		menu.entries[breakIndicator] = makeMenuEntry(val, hint, afunc)
	}
	menu.quitValue = val

	//error can be returned only AFTER the default is set
	if val == defBreakv {
		return errors.New(errmsg)
	}

	return nil
}

//valueClean : simple helper function to trim and (semi)validate parameters sent to func()s
//afunc will be run if it is a non-empty func().
func valueClean(toClean, theDefault *string, where vClean, afunc func()) {
	//this is mostly for programmer errors and/or dynamic menu errors.
	//if menus are properly constructed, one will never see this displayed
	//except when vAfterBlock is specified.
	*toClean = strings.Trim(*toClean, trimString)
	if *toClean == "" {
		*toClean = *theDefault
		if where == vInBlock {
			where.Display()
			afunc()
		}
	}
	if where == vAfterBlock {
		afunc()
	}
}

//AddMenuEntry : Method that adds a MenuEntry to a Menu.
//aval is the string the user must type to initiate that menu entry.
func (menu *Menu) AddMenuEntry(aval string, ahint string, afunc func()) error {

	errmsg := warn + fmt.Sprintf("AddMenuEntry method: Menu '%s', Empty Entry Value sent in (hint was '%s'), Entry not added.\n", menu.Title, ahint)

	valueClean(&aval, &emptyString, vInBlock, func() {
		alertUser(&errmsg)
	})
	//the return and alert need to be separate here
	if aval == "" {
		return errors.New(errmsg)
	}
	valueClean(&ahint, &emptyHint, vIgnore, func() {})
	menu.entries[aval] = &menuEntry{
		value:          aval,
		hint:           ahint,
		isSubMenuEntry: false,
		doRun:          afunc,
	}
	//count the number of times the key has been added
	menu.validateKeys[aval]++
	menu.isModified = menu.finalized
	return nil
}

//RemoveMenuEntry : Remove a menu entry from a *Menu, only useful for dynamic menus.
//key is the menuEntry's value the user would type to run it.
func (menu *Menu) RemoveMenuEntry(key string) error {
	const (
		methodName = "RemoveMenuEntry method: "
	)
	var errmsg string
	valueClean(&key, &emptyString, vInBlock, func() {
		errmsg = warn + fmt.Sprintf(methodName+"Menu '%s', empty key value sent in, can't  remove.", menu.Title)
		alertUser(&errmsg)
	})
	//alert and error return must be separate here
	if key == "" {
		return errors.New(errmsg)
	}

	if key == menu.entries[breakIndicator].value {
		errmsg = warn + fmt.Sprintf(methodName+"Menu '%s', key '%s' is the menu break key, removal not allowed.", menu.Title, key)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	if len(menu.entries) < 2 {
		//this will probably never happen because the menu break key should always exist
		errmsg = warn + fmt.Sprintf(methodName+"Menu '%s', Only 1 remaining entries, not allowed to completely empty Menu.\n", menu.Title)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	if _, ok := menu.entries[key]; !ok {
		errmsg = warn + fmt.Sprintf(methodName+"Menu '%s', key '%s' does not exist, nothing to remove.\n", menu.Title, key)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}
	//adjust the validateKeys list to reflect the change
	menu.validateKeys[key]--
	delete(menu.entries, key)
	menu.isModified = menu.finalized
	return nil
}

//alertUser : internal function. All most all (or mabye all) error conditions
//call this func in addition to returning an error. Very useful for
//debugging, but also just generally a good idea to leave on unless
//you always check error values. These displays can be turned off/on
//with MenuOptions.
func alertUser(errmsg *string) {
	if !MenuOptions.runTimeErrMsgs_Display {
		return
	}
	fmt.Printf("\n" + *errmsg + "\n")
	if MenuOptions.runTimeErrMsgs_Pause {
		WaitForInput(&additionalString)
	}
}

//ChangeMenuEntry : Change MenuEntry value and/or hint, only useful for dynamic menus.
//enpty strings indicate ignoring that setting. If either of oldkey
//or newkey is empty ("") the key value will not be changed
//valid patterns:
//"newhint", "oldkey", ""
//"", "oldkey", "newkey"
//"newhint", "oldkey", "newkey"
func (menu *Menu) ChangeMenuEntry(newhint, oldkey, newkey string) error {
	const (
		methodName = "ChangeMenuEntry method: "
		endMsg     = " Entry not changed."
	)
	var errmsg string

	valueClean(&oldkey, &emptyString, vIgnore, func() {})
	if oldkey == "" {
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s': Param 'oldkey' was invalid or blank.", menu.Title) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	if _, ok := menu.entries[oldkey]; !ok { //old key doesn't exist
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s': Entry oldkey, value '%s', does not exist.", menu.Title, oldkey) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	valueClean(&newkey, &emptyString, vIgnore, func() {})
	valueClean(&newhint, &emptyString, vIgnore, func() {})
	if oldkey == newkey {
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s': oldkey and newkey are the same value '%s'.", menu.Title, oldkey) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}
	if oldkey == breakIndicator || newkey == breakIndicator {
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s': oldkey '%s' or newkey '%s' attempting to change BREAK value.", menu.Title, oldkey, newkey) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	//ready to apply changes...
	if newkey != "" {
		var ahint string
		if newhint == "" {
			ahint = menu.entries[oldkey].hint
		} else {
			ahint = newhint
		}
		menu.AddMenuEntry(newkey, ahint, menu.entries[oldkey].doRun)
		menu.transferEntryFields(newkey, oldkey)
		menu.RemoveMenuEntry(oldkey)
		menu.isModified = menu.finalized
		return nil
	}

	if newhint == "" {
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s': newhint param was blank.", menu.Title) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}
	menu.entries[oldkey].hint = newhint
	menu.isModified = menu.finalized
	return nil
}

//ChangeMenuEntryFunc : send in a new func() for this menu entry. You are not allowed to change
//the func() for a registered subMenu. If you want to change the Break Item's
//func() just call <menuvar>.SetMenuBreakItem again. key is the Value of the menu entry that
//the user would type.
func (menu *Menu) ChangeMenuEntryFunc(key string, afunc func()) error {
	const (
		methodName = "ChangeMenuEntryFunc method: "
		endMsg     = " Entry not changed."
	)
	var errmsg string

	valueClean(&key, &emptyString, vIgnore, func() {})
	if key == "" {
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s': Param 'key' was blank.", menu.Title) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	if key == breakIndicator {
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s': key '%s'! Attempting to change BREAK value, not allowed.", menu.Title, key) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	entry, ok := menu.entries[key]
	if !ok { //old key doesn't exist
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s': Entry key, value '%s', does not exist.", menu.Title, key) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	if entry.isSubMenuEntry {
		errmsg = warn + methodName + fmt.Sprintf("Menu '%s', key '%s': Attempting to change a SubMenu func(), not allowed.", menu.Title, key) + endMsg
		alertUser(&errmsg)
		return errors.New(errmsg)
	}

	//ready to apply changes...
	entry.doRun = afunc
	menu.isModified = menu.finalized
	return nil
}

//transferEntryFields : internal func() to pass any important values when changing a Menu Entry
func (menu *Menu) transferEntryFields(newkey, oldkey string) {
	//right now only 1 "property" on Menu Entries, make sure they are transferred
	menu.entries[newkey].isSubMenuEntry = menu.entries[oldkey].isSubMenuEntry
}

//ChangeMenuTitle : Change a Menu's Title, really only useful for dynamic menus. You also have
//access to <menuvar>.Title, but this validates so recommended to use this func().
func (menu *Menu) ChangeMenuTitle(newtitle string) error {
	var errmsg string
	valueClean(&newtitle, &emptyString, vInBlock, func() {
		errmsg = warn + fmt.Sprintf("ChangeMenuTitle method: Menu '%s', empty Title sent in, can't  change title.", menu.Title)
		alertUser(&errmsg)
	})
	if newtitle == "" {
		return errors.New(errmsg)
	}
	menu.Title = newtitle
	return nil
}

//reSet : internal only. When using dynamic Menus (i.e., changing keys, hints, titles, funcs)
//this function MUST be called to actualize the changes and
//reset the display orders. If not called you can be assured of a panic!
func (menu *Menu) reSet() error {
	var errmsg string
	if !menu.finalized {
		errmsg = warn + fmt.Sprintf("reSet method: Menu '%s' cannot be reSet(), it has never been Start'ed", menu.Title)
		alertUser(&errmsg)
		return errors.New(errmsg)
	}
	menu.sortKeys = nil
	menu.finalized = false
	err := menu.finalize()

	if err != nil {
		menu.killThisMenu = true
		return err
	}
	return nil
}

//doValidate : Validation of Menu Entries. Errors should be added at top of func()
//because the Menu system may not be started in some cases. Warnings go at bottome of func()
//so the user will be made aware of declaration problems, but the Menu will run.
func (menu *Menu) doValidate() (string, error) {
	const (
		breakStr     = "breakIndicator"
		menuStr      = "Menu Entry"
		breakErr     = "doValidate(): Menu '%s' has no %s to quit the menu loop."
		killErr      = "doValidate(): Menu '%s' has a %s '%s' which conflicts with Menu system kill phrase '%s'."
		keysMsg      = "entry '%s' was added %d times\n"
		keysBreak    = ">> Menu '%s' has a Value '%s' which conflicts with Menu Break Value, Entry ignored\n"
		keysWarning  = ">> Menu '%s' was declared with duplicate Key values:\n"
		keysFunction = "The menu may not function as you intended.\n"
	)

	var errmsg string

	item, ok := menu.entries[breakIndicator]
	if !ok {
		errmsg = warn + fmt.Sprintf(breakErr, menu.Title, breakStr)
		alertUser(&errmsg)
		return "", errors.New(errmsg)
	}

	if item.value == MenuOptions.killPhrase {
		errmsg = warn + fmt.Sprintf(killErr, menu.Title, breakStr, menu.entries[breakIndicator].value, MenuOptions.killPhrase)
		alertUser(&errmsg)
		return "", errors.New(errmsg)
	}

	if _, ok2 := menu.validateKeys[MenuOptions.killPhrase]; ok2 {
		errmsg = warn + fmt.Sprintf(killErr, menu.Title, menuStr, menu.entries[MenuOptions.killPhrase].value, MenuOptions.killPhrase)
		alertUser(&errmsg)
		return "", errors.New(errmsg)
	}
	//errors go above this comment and message(s) go below
	//basically error conditions flag the menu as faulty
	//while messages are important information but do not stop menu functionality
	var report = ""

	if _, ok3 := menu.validateKeys[menu.entries[breakIndicator].value]; ok3 {
		report = report + fmt.Sprintf(keysBreak, menu.Title, menu.entries[breakIndicator].value)
	}

	var reportKeys = ""
	for k := range menu.validateKeys {
		if menu.validateKeys[k] > 1 {
			reportKeys = reportKeys + fmt.Sprintf(keysMsg, k, menu.validateKeys[k])
		}
	}

	if reportKeys != "" {
		reportKeys = fmt.Sprintf(keysWarning, menu.Title) + reportKeys + keysFunction
		report = report + reportKeys + "\n"
	}
	//reset the validateKeys for the case of dynamic menu manipulation
	//during the running of the menu system
	for k := range menu.validateKeys {
		delete(menu.validateKeys, k)
	}

	//other "report"s can be added, as needed

	return report, nil
}

//finalize : Sorts the menu entries and ensures the BREAK value is last, or
//first, depending on Menu.reverseSort, by setting Menu.sortKeys
//also performs some validations, and validates dynamically manipulated Menus
func (menu *Menu) finalize() error {

	if menu.finalized {
		return nil
	}

	msg, err := menu.doValidate()
	if err != nil {
		return err
	}
	//error and alert must be separated here. Above is a faulty menu,
	//below the menu has an irregularity but will run
	if msg != "" {
		alertUser(&msg)
	}

	menu.sortKeys = make([]string, 0, len(menu.entries))

	skipKey := menu.entries[breakIndicator].value

	for _, k := range menu.entries {
		if k.value == skipKey {
			continue
		}
		menu.sortKeys = append(menu.sortKeys, k.value)
	}
	sort.Strings(menu.sortKeys)
	//ensure quit key is at bottom (or top depending)
	menu.sortKeys = append(menu.sortKeys, breakIndicator)

	if menu.reverseSort {
		//probably a better alogorithm somewhere
		cpy := make([]string, len(menu.sortKeys))
		copy(cpy, menu.sortKeys)
		ln := len(cpy) - 1
		for i := ln; i > -1; i-- {
			menu.sortKeys[ln-i] = cpy[i]
		}
	}

	menu.finalized = true
	menu.killThisMenu = false
	menu.isModified = false

	return nil
}

//displayMenu : Prints the menu to screen, used internally
func (menu *Menu) displayMenu() {

	var killMsg string
	if MenuOptions.killPhrase != "" {
		killMsg = fmt.Sprintf(killTemplate, MenuOptions.killPhrase)
	} else {
		killMsg = "=============================="
	}

	fmt.Println("")

	//get and print menu breadcrumbs, this is the most dangerous part
	//of this code because here one could go circular. All the current
	//validation code prevents this from happening.
	menuSep := fmt.Sprintf(" %s ", MenuOptions.menuSeparator)
	parentsTitles := ""
	menuParent := menu.parent
	for menuParent != nil {
		parentsTitles = menuParent.Title + menuSep + parentsTitles
		menuParent = menuParent.parent
	}
	fmt.Println(parentsTitles + menu.Title)

	fmt.Println("------------------------------")

	//print the menu
	for _, k := range menu.sortKeys {
		fmt.Fprintln(alignerLocal, fmt.Sprintf(menuFormat, menu.entries[k].value, menu.entries[k].hint))
	}

	alignerLocal.Flush()
	fmt.Println(killMsg)
	fmt.Print(MenuOptions.menuPrompt)
}

//droppingDown : internal use. If a Menu Entry Start()-s an already open menu
//this function tells the menu to close itself if it is not the
//target menu to drop down to.
func (menu *Menu) droppingDown() bool {
	if dropDown.doDropDown {
		if menu.id != dropDown.id {
			return true
		} else {
			dropDown.doDropDown = false
			return false
		}
	}
	return false
}

//getSwitch : internal use. in the case of dropping down to a menu that is already open
//We want to bypass user input so that the transition is smooth
func getSwitch() bracketSwitch {
	switch dropDown.isLastExit {
	case true:
		dropDown.isLastExit = false
		return bsPartial
	default:
		return bsBottom
	}
}

//printFuncBrackets : internal use. Management of bracketing wrappers around a menu entry's func() results
func (menu *Menu) printFuncBrackets(brType bracketSwitch, choice string) {
	const (
		funcRunnerID = "  %s Menu: '%s' - choice: '%s'\n"
		isBegin      = ">>"
		isEnd        = "<<"
	)

	getfuncRunnerStr := func(indicator string) string {
		if MenuOptions.idFuncRunner {
			return fmt.Sprintf(funcRunnerID, indicator, menu.Title, choice)
		} else {
			return ""
		}
	}

	switch brType {
	case bsTop:
		fmt.Println("\n\n\n" + MenuOptions.funcBracketTop + getfuncRunnerStr(isBegin))
	case bsBottom:
		fmt.Println(MenuOptions.funcBracketBottom + getfuncRunnerStr(isEnd))
		if MenuOptions.pauseOnOutput {
			WaitForInput(&emptyString)
		}
	case bsPartial:
		fmt.Println(MenuOptions.funcBracketBottom + getfuncRunnerStr(isEnd))
	}
}

//Start() : Start the calling Menu. Can be called on any menu at any time.
//Only the main (root) menu needs to be Start'ed to run the menu system.  Or, one
//can use the StartMenuSystem() call. This can be called on ANY menu so that
//the programmer can do as she likes. But for a managed system it is recommended
//to use AddSubMenu() function for subMenus rather than calling their Start().
func (menu *Menu) Start() error {

	var errmsg string

	if menu.isRunning {
		dropDown.doDropDown = true
		dropDown.id = menu.id
		dropDown.isLastExit = true
		return nil
	}

	//menu was dynamically changed in the background; i.e., it was not running at the time.
	if menu.finalized && menu.isModified {
		menu.reSet()
	}

	if !menu.finalized {
		startErr := menu.finalize()
		if startErr != nil {
			errmsg = warn + fmt.Sprintf("func Start(), not Start()'ing Menu '%s'\n>>> %s", menu.Title, startErr)
			alertUser(&errmsg)
			return errors.New(errmsg)
		}
	}

	menu.displayMenu()
	menu.setRunning(true)
	defer menu.setRunning(false)

	var input string
	for menuScanner.Scan() {

		input = strings.Trim(menuScanner.Text(), " \t")

		if input != "" && input == MenuOptions.killPhrase {
			//master killPhrase was used, stop the menu system
			fmt.Println("Stopping Menu system...")
			killSwitch = true
			break
		}

		if input == menu.quitValue {
			//The break indicator can also have a func(), so we run it.
			menu.entries[breakIndicator].doRun()
			break
		}

		if elem, ok := menu.entries[input]; ok {

			if !elem.isSubMenuEntry && !menu.isChooseOne {
				menu.printFuncBrackets(bsTop, input)
			}

			//run the associated menu entry's func()
			elem.doRun()

			//menu was dynamically changed while the menu was running
			if menu.isModified {
				menu.reSet()
			}

			if menu.killThisMenu || killSwitch || menu.isChooseOne || menu.droppingDown() {
				break
			}

			//this helps directly Start()-ed menus behave similarly to a
			//properly defined SubMenu. Niche use, but useful.
			if menu.skipFunctionNotification {
				menu.skipFunctionNotification = false
				menu.printFuncBrackets(bsPartial, input)
				if MenuOptions.pauseOnOutput {
					fmt.Println(fmt.Sprintf("< Function pause bypassed by  %s.skipFunctionNotification >", menu.Title))
				}
				menu.displayMenu()
				continue
			}

			if !elem.isSubMenuEntry {
				menu.printFuncBrackets(getSwitch(), input)
			}

			menu.displayMenu()

		} else {
			//No.... menu.displayMenu() too annoying to do this in the case of a mis-typed Menu Key
			fmt.Println("???? " + MenuOptions.menuPrompt + " '" + input + "' is not a valid menu choice...")
			fmt.Print(MenuOptions.menuPrompt)
		}
	}
	return nil
}

//WaitForInput : Gives the user opportunity to read output
//before the menu system continues on. You can use this if you
//want to make sure the user reads somthing. additional is
//optional but would give addiational information in the message.
func WaitForInput(additional *string) {
	fmt.Println(*additional + "\nPress <RET> to continue...")
	menuScanner.Scan()
	menuScanner.Text()
}

//GetUserInput : Good for basic input, returns the string the User enters.
//The calling func must deal with typecasting.
//Looping input is possible but probably (?) needs to be made for
//each type, and would need a passed function to ship each
//input out to?? I have not tested this func() under all circumstances.
func GetUserInput(prompt string) string {
	valueClean(&prompt, &emptyString, vIgnore, func() {})
	if prompt == "" {
		prompt = "Enter data: "
	}
	fmt.Println(fmt.Sprintf("%s  [%s]:", prompt, "<RET> cancels"))
	fmt.Print("=> ")
	menuScanner.Scan()
	input := strings.Trim(menuScanner.Text(), trimString)
	if input == "" {
		fmt.Println("<canceled>")
	}
	return input
}

//descriptions of menuOptions settings.
const (
	idFuncRunnerInfo = `idFuncRunner: If true then when a Menu Entry runs 
its function it will be noted at function start and
function end, which Menu (by Title) and which Menu
choice started that function. Useful for debugging.
Works with, but is separate from funcBracketTop and
funcBracketBottom.`
	killPhraseInfo = `killPhrase: Input at the prompt that causes exit 
from menu system no matter how deeply nested. 
Empty string ("") disables this ability.`
	menuPromptInfo = `menuPrompt: Prompt for the menu user. This can be 
set even to " ", but an empty "" will default to 
prompt default value.`
	menuSeparatorInfo = `menuSeparator: Separator between Menu Title breadcrumbs.`
	pauseOnOutputInfo = `pauseOnOutput: if true will wait for user to press 
<RET> after outputting func() results. Otherwise 
immediately returns menu display`
	runTimeErrMsgs_DisplayInfo = `runTimeErrMsgs_Display: Certain methods print out 
error messages directly to the screen, as well as 
returning errors. The printed messages are useful 
during development, particularly for dynamically
manipulated menus. Once you are happy with your 
code, this can be set to false. But then you will need 
to check returned error values manually to handle 
and/or print out the error messages. A false value 
here implies a false value for the _Pause variant`
	runTimeErrMsgs_PauseInfo = `runTimeErrMsgs_Pause: If this value is true, runtime 
err msgs will pause for user input acknowledgement 
so they can be read. If false the err msgs will display 
if the _Display variant is true, but the menu system 
will not pause after each msg. There is a lot of 
validation and error conditions checking but almost 
all of them are not fatal and Menus will run in almost 
any case.`
	funcBracketTopInfo = `funcBracketTop: func Brackets are strings printed before and 
after a Menu\'s func() to make the output easier to differentiate
from the Menu output. This string will be printed before the 
func() starts running. Works with, but is separate from 
idFuncRunnerInfo.`
	funcBracketBottomInfo = `funcBracketBottom: This func Bracket string will be printed 
after the func() finishes. Works with, but is separate from 
idFuncRunnerInfo.`
)
