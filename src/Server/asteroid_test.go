package main

import "testing"

//Movement a asteroid three times in a pase of xStep and yStep
func Movement(xStep int, yStep int) *asteroid{
	
	a := newAsteroid()
	a.x = 5
	a.y = 5

	a.xStep = xStep
	a.yStep = yStep

	a.move()
	a.move()
	a.move()

	return a
}

// TestMovementHorisontal tests to move a asteroid Horisontal
func TestMovementHorisontal(t *testing.T){

	var roid1 = Movement(1,0)
	
	if roid1.x != 8 || roid1.y!= 5 {
		t.Error("Expected position (x,y) = (8,5), got (",roid1.x,",",roid1.y,")")
	}

	var roid2 = Movement(-1,0)
	
	if roid2.x != 2 || roid2.y != 5 {
		t.Error("Expected position (x,y) = (2,5), got (",roid2.x,",",roid2.y,")")
	}
}


// TestMovementVertical test to move a asteroid vertical

func TestMovementVertical(t *testing.T){
	var roid1 = Movement(0,1)

	if roid1.x != 5 || roid1.y!= 8 {
		t.Error("Expected position (x,y) = (5,8), got (",roid1.x,",",roid1.y,")")
	}

	var roid2 = Movement(0,-1)

	if roid2.x != 5 || roid2.y!= 2 {
		t.Error("Expected position (x,y) = (5,2), got (",roid2.x,",",roid2.y,")")
	}	
}

// TestMovementDiagonal test to move a asteroid in a diagonal
func TestMovementDiagonal(t *testing.T){
	var roid1 = Movement(1,1)

	if roid1.x != 8 || roid1.y!= 8 {
		t.Error("Expected position (x,y) = (8,8), got (",roid1.x,",",roid1.y,")")
	}

	var roid2 = Movement(-1,-1)

	if roid2.x != 2 || roid2.y!= 2 {
		t.Error("Expected position (x,y) = (2,2), got (",roid2.x,",",roid2.y,")")
	}
}




