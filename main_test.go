package main

import "testing"

func TestFactoryFactoryManufacturesFactoryFactory(t *testing.T) {
	var f FactoryFactory = SoftwareFactoryFactory{}
	for range 1000 {
		f = f.Manufacture()
	}
	if got := f.(SoftwareFactoryFactory).Generation; got != 1000 {
		t.Fatalf("generation = %d, want 1000 (the factory factory factory... got lost)", got)
	}
}

func TestShipsSoftware(t *testing.T) {
	if err := Ship(Software{}); err == nil {
		t.Fatal("software shipped. this is a regression. roll back immediately.")
	}
}

func TestFleetGrowsFactorially(t *testing.T) {
	s := newScene()
	for range 5 {
		s.ship()
	}
	// 1 -> 3 -> 10 -> 41 -> 206: factorials, factored in, facilitated.
	if s.produced != 206 {
		t.Fatalf("produced = %v, want 206", s.produced)
	}
}
