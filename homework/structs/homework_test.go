package main

import (
	"bytes"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	fourBits  = 4
	sixBits   = 6
	eightBits = 8
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		var p [42]byte
		buf := bytes.NewBufferString(name)
		for i, b := range buf.Bytes() {
			p[i] = b
		}
		person.personName = p
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndRespect[0] = byte(mana)
		person.manaAndRespect[1] = byte(mana >> eightBits)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthAndStrength[0] = byte(health)
		person.healthAndStrength[1] = byte(health >> eightBits)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndRespect[1] = byte(respect<<4) | person.manaAndRespect[1]
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthAndStrength[1] = byte(strength<<4) | person.healthAndStrength[1]
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelAndExperience = byte(experience) | person.levelAndExperience
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelAndExperience = byte(level<<4) | person.levelAndExperience
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.homeWeaponFamilyFraction = person.homeWeaponFamilyFraction | _WithHome
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.homeWeaponFamilyFraction = person.homeWeaponFamilyFraction | _WithWeapon
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.homeWeaponFamilyFraction = person.homeWeaponFamilyFraction | _WithFamily
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.homeWeaponFamilyFraction = person.homeWeaponFamilyFraction | byte(personType<<3)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
	_WithHome   = 1
	_WithWeapon = 2
	_WithFamily = 4
)

type GamePerson struct {
	x                        int32
	y                        int32
	z                        int32
	gold                     uint32
	personName               [42]byte
	healthAndStrength        [2]byte // byte0: xxxxxxxx, byte1: xxyyyyyy, x = health, y = strength
	manaAndRespect           [2]byte // byte0: xxxxxxxx, byte1: xxyyyyyy, x = mana, y = respect
	levelAndExperience       byte
	homeWeaponFamilyFraction byte
}

func NewGamePerson(options ...Option) GamePerson {
	g := GamePerson{}
	for _, option := range options {
		option(&g)
	}
	return g
}

func (p *GamePerson) Name() string {
	buf := bytes.NewBuffer([]byte{})
	buf.Grow(len(p.personName))
	for _, b := range p.personName {
		buf.WriteByte(b)
	}
	return buf.String()
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.manaAndRespect[0]) + int(p.manaAndRespect[1]<<sixBits)>>sixBits*256
}

func (p *GamePerson) Health() int {
	return int(p.healthAndStrength[0]) + int(p.healthAndStrength[1]<<sixBits)>>sixBits*256
}

func (p *GamePerson) Respect() int {
	return int(p.manaAndRespect[1] >> fourBits)
}

func (p *GamePerson) Strength() int {
	return int(p.healthAndStrength[1] >> fourBits)
}

func (p *GamePerson) Experience() int {
	return int((p.levelAndExperience << fourBits) >> fourBits)
}

func (p *GamePerson) Level() int {
	return int(p.levelAndExperience >> fourBits)
}

func (p *GamePerson) HasHouse() bool {
	return p.homeWeaponFamilyFraction&_WithHome > 0
}

func (p *GamePerson) HasGun() bool {
	return p.homeWeaponFamilyFraction&_WithWeapon > 0
}

func (p *GamePerson) HasFamily() bool {
	return p.homeWeaponFamilyFraction&_WithFamily > 0
}

func (p *GamePerson) Type() int {
	return int(p.homeWeaponFamilyFraction) >> 3
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
