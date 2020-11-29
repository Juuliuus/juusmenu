package juusmenu

//testing is currently under development so tests are minimal and crude.
//I need to study how to test a unit that requires user input.

import (
	"fmt"
	"testing"
)

func TestUnInitializedMenu(t *testing.T) {
	MenuOptions.SetRunTimeErrMsgs_Display(false)
	tmpMenu := NewMenu("NewMenu")
	err := tmpMenu.Start()
	if err == nil {
		t.Error("Failed: Menu should return non nil error if it is not initialized.")
	}
}

func TestDuplicateIDs(t *testing.T) {
	//todo ok so here is a good example of what to do about different errors.
	//SetID will return an error FIRST if the menu.isRunning=true.
	//so this test will give a false positive if tmpMenu2 were running...
	MenuOptions.SetRunTimeErrMsgs_Display(false)
	tmpMenu, tmpMenu2 := NewMenu("TmpMenu"), NewMenu("TmpMenu2")
	tmpMenu.SetID(42)
	err := tmpMenu2.SetID(42)
	if err == nil {
		t.Error("Failed: Duplicate ID's are not allowed.")
	}
}

func TestAddSubMenuStub(t *testing.T) {
	//todo so that func has lots of validations, so I would need a test
	//for each validation that I make, and other random tests that test things
	//I didn't think of. Need to think about an efficient way to do this.
	MenuOptions.SetRunTimeErrMsgs_Display(false)
	tmpMenu, tmpMenu2 := NewMenu("TmpMenu"), NewMenu("TmpMenu2")
	tmpMenu.SetMenuBreakItem("q", "Quit", func() { fmt.Println("byebye") })
	tmpMenu2.SetMenuBreakItem("b", "b", func() { fmt.Println("back") })
	tmpMenu2.parent = tmpMenu
	//test the parent nil validation
	err := tmpMenu.AddSubMenu(tmpMenu2, "1", "a submenu")
	if err == nil {
		t.Error("Failed: func should have complained that the menu already has a parent.")
	}
}

func TestManyMenus(t *testing.T) {
	//needs (much) work to make this "interactive".
	//todo also do I need to fix source code to not display the menu under testing???
	menuNum := 0
	menus := make([]*Menu, 0, 3)
	var tmpMenu *Menu
	MenuOptions.SetRunTimeErrMsgs_Display(false)

	for i := 0; i < 3; i++ {
		menuNum++
		key := fmt.Sprintf("Q%d", menuNum)
		tmpMenu = NewMenu(fmt.Sprintf("Menu%d", menuNum))
		tmpMenu.SetMenuBreakItem(key, "Quit", func() { fmt.Println("byebye") })
		menus = append(menus, tmpMenu)
	}

	for i := 0; i < len(menus); i++ {
		menus[i].Start()
	}
}
