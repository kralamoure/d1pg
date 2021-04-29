// WARNING: Extremely hacky code ahead — seriously.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
	"github.com/kralamoure/d1proto"
	"github.com/kralamoure/d1util"

	"github.com/kralamoure/d1pg"
)

var errInvalidAssertion = errors.New("invalid assertion")

var (
	ctx  = context.Background()
	pool *pgxpool.Pool
	repo *d1pg.Repo
)

type ItemTemplate struct {
	Name              string        `json:"n"`
	Description       string        `json:"d"`
	Type              int           `json:"t"`
	Enhanceable       bool          `json:"fm"`
	TwoHands          bool          `json:"tw"`
	Ethereal          bool          `json:"et"`
	Hidden            bool          `json:"h"`
	ItemSetId         int           `json:"s"`
	CanUse            bool          `json:"u"`
	CanTarget         bool          `json:"ut"`
	Level             int           `json:"l"`
	GFX               int           `json:"g"`
	Price             int           `json:"p"`
	Weight            int           `json:"w"`
	Cursed            bool          `json:"m"`
	Conditions        string        `json:"c"`
	WeaponEffects     []interface{} `json:"e"`
	WeaponEffectsReal WeaponEffects
}

type WeaponEffects struct {
	CriticalHitBonus int
	APCost           int
	RangeMin         int
	RangeMax         int
	CriticalHit      int
	CriticalFailure  int
	LineOnly         bool
	LineOfSight      bool
}

type D2Level struct {
	Id                         int
	SpellId                    int
	SpellBreed                 int
	APCost                     int
	MinRange                   int
	MaxRange                   int
	CastInLine                 bool
	CastTestLOS                bool
	CriticalHitProbability     int
	CriticalFailureProbability int
	NeedsFreeCell              bool
	RangeCanBeBoosted          bool
	MaxCastPerTurn             int
	MaxCastPerTarget           int
	MaxCastPerInterval         int
	MinPlayerLevel             int
	CriticalFailureEndsTurn    bool
	StatesRequired             []int
	StatesForbidden            []int
	Effects                    []D2Effect
	CriticalEffects            []D2Effect
}

type D2Effect struct {
	Id        int
	Duration  int
	Param1    int
	Param2    int
	Param3    int
	ZoneSize  int
	ZoneShape int
	TargetId  int
}

func main() {
	err := run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	tmp, err := pgxpool.Connect(ctx, "postgresql://postgres:password@localhost/dofus")
	if err != nil {
		return err
	}
	pool = tmp
	defer pool.Close()

	repo, err = d1pg.NewRepo(pool)
	if err != nil {
		return err
	}

	// return orderItemTemplateEffects()
	// return orderItemSetEffects()
	// return decryptGameMapData()
	// return createNPCDialogsAndResponses()
	// return createEffectTemplates()
	// return createSystemMarketItems()
	// return createSystemMarketWeaponElements()
	// return createClasses()
	// return createSpells()
	// return createMountTemplates()
	// return addTargetIdToSpells()

	return nil
}

func createMountTemplates() error {
	p, err := ioutil.ReadFile("C:/Users/raul/AppData/Roaming/airtest/Local Store/f.json")
	if err != nil {
		return err
	}

	var raws map[int]struct {
		Name   string `json:"n"`
		GFXId  string `json:"g"`
		Color1 string `json:"c1"`
		Color2 string `json:"c2"`
		Color3 string `json:"c3"`
	}

	err = json.Unmarshal(p, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		color1N, err := strconv.Atoi(v.Color1)
		if err != nil {
			return err
		}
		color1 := ""
		if color1N != -1 {
			color1 = fmt.Sprintf("%x", color1N)
		}

		color2N, err := strconv.Atoi(v.Color2)
		if err != nil {
			return err
		}
		color2 := ""
		if color2N != -1 {
			color2 = fmt.Sprintf("%x", color2N)
		}

		color3N, err := strconv.Atoi(v.Color3)
		if err != nil {
			return err
		}
		color3 := ""
		if color3N != -1 {
			color3 = fmt.Sprintf("%x", color3N)
		}

		query := "INSERT INTO d1_static.mounts (id, name, gfx_id, color_1, color_2, color_3)" +
			" VALUES ($1, $2, $3, $4, $5, $6);"

		_, err = pool.Exec(ctx, query,
			k, v.Name, v.GFXId, color1, color2, color3,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func createSpells() error {
	p, err := ioutil.ReadFile("C:/Users/raul/AppData/Roaming/airtest/Local Store/f.json")
	if err != nil {
		return err
	}

	var raws map[int]struct {
		Name        string        `json:"n"`
		Description string        `json:"d"`
		Level1      []interface{} `json:"l1"`
		Level2      []interface{} `json:"l2"`
		Level3      []interface{} `json:"l3"`
		Level4      []interface{} `json:"l4"`
		Level5      []interface{} `json:"l5"`
		Level6      []interface{} `json:"l6"`
	}

	err = json.Unmarshal(p, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		var levels []string

		p, err := json.Marshal(v.Level1)
		if err != nil {
			return err
		}
		levels = append(levels, string(p))

		p, err = json.Marshal(v.Level2)
		if err != nil {
			return err
		}
		levels = append(levels, string(p))

		p, err = json.Marshal(v.Level3)
		if err != nil {
			return err
		}
		levels = append(levels, string(p))

		p, err = json.Marshal(v.Level4)
		if err != nil {
			return err
		}
		levels = append(levels, string(p))

		p, err = json.Marshal(v.Level5)
		if err != nil {
			return err
		}
		levels = append(levels, string(p))

		if v.Level6 != nil {
			p, err = json.Marshal(v.Level6)
			if err != nil {
				return err
			}
			levels = append(levels, string(p))
		}

		query := "INSERT INTO d1_static.spells (id, name, description, levels)" +
			" VALUES ($1, $2, $3, $4);"

		_, err = pool.Exec(ctx, query,
			k, v.Name, v.Description, levels,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func createClasses() error {
	p, err := ioutil.ReadFile("C:/Users/raul/AppData/Roaming/airtest/Local Store/f.json")
	if err != nil {
		return err
	}

	var raws map[int]struct {
		Name                  string  `json:"sn"`
		Label                 string  `json:"ln"`
		ShortDescription      string  `json:"sd"`
		Description           string  `json:"d"`
		Spells                []int   `json:"s"`
		BoostCostVitality     [][]int `json:"b11"`
		BoostCostWisdom       [][]int `json:"b12"`
		BoostCostStrength     [][]int `json:"b10"`
		BoostCostIntelligence [][]int `json:"b15"`
		BoostCostChance       [][]int `json:"b13"`
		BoostCostAgility      [][]int `json:"b14"`
	}

	err = json.Unmarshal(p, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		p, err := json.Marshal([][][]int{
			v.BoostCostVitality,
			v.BoostCostWisdom,
			v.BoostCostStrength,
			v.BoostCostIntelligence,
			v.BoostCostChance,
			v.BoostCostAgility,
		})
		if err != nil {
			return err
		}

		query := "INSERT INTO d1_static.classes (id, name, label, short_description, description, spells, boost_costs)" +
			" VALUES ($1, $2, $3, $4, $5, $6, $7);"

		_, err = pool.Exec(ctx, query,
			k, v.Name, v.Label, v.ShortDescription, v.Description, v.Spells, string(p),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func createSystemMarketWeaponElements() error {
	// Damage rates are 50%, 68% and 85%

	marketId := "31dbf930-0cfc-475e-989e-8c9ca171e1f3"
	etherealMarketId := "f1c956f3-1486-4ffa-b29e-8925f4c77964"

	elementIds := []int{97, 99, 96, 98}

	for i := 0; i < 2; i++ {
		marketId := marketId
		if i == 1 {
			marketId = etherealMarketId
		}
		marketItems, err := repo.MarketItemsByMarketId(ctx, marketId)
		if err != nil {
			return err
		}

		for _, v := range marketItems {
			heals := false
			for _, effect := range v.Effects {
				if effect.Id == 108 {
					heals = true
					break
				}
			}
			for i := 0; i < 2; i++ {
				rate := 85
				if i == 1 {
					if !heals {
						continue
					}
					rate = 50
				}

				for _, elementId := range elementIds {
					effects := make([]d1typ.Effect, len(v.Effects))

					changed := false
					for i, effect := range v.Effects {
						if effect.Id == 100 {
							changed = true

							effect.Id = elementId

							if effect.DiceSide != 0 {
								effect.DiceSide = int(float32(math.Floor(float64(float32(effect.DiceNum-1)*float32(rate)/100))) + float32(math.Floor(float64(float32(effect.DiceSide-effect.DiceNum+1)*float32(rate)/100))))
							}
							effect.DiceNum = int(math.Floor(float64((float32(effect.DiceNum)-1)*float32(rate)/100)) + 1)
							if effect.DiceSide <= effect.DiceNum {
								effect.DiceSide = 0
							}

							effect.Param = d1.EffectDiceParam(effect)
						}
						effects[i] = effect
					}

					if changed {
						_, err := repo.CreateMarketItem(ctx, d1.MarketItem{
							Item: d1.Item{
								TemplateId: v.TemplateId,
								Quantity:   1,
								Effects:    effects,
							},
							Price:    1,
							MarketId: marketId,
						})
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func createSystemMarketItems() error {
	gameServerId := 2
	marketId := "31dbf930-0cfc-475e-989e-8c9ca171e1f3"
	etherealMarketId := "f1c956f3-1486-4ffa-b29e-8925f4c77964"
	effectIdsToRemove := []int{205, 208, 209, 601, 740, 785, 983, 800, 806, 807, 808, 811}
	disallowedItemIds := []int{9031, 9202, 9919, 9396, 6894, 6895, 7913, 7920, 2154, 2155, 2156, 6713, 8575, 8854, 10076, 10073, 9627, 2170, 8627, 1505, 6971, 6975, 8574, 7043, 7112, 8098, 10125, 10126, 10127, 10133, 1944, 1628, 1629, 1630, 1631, 1632, 1633, 684, 1710, 958, 1099, 10846, 10677, 9641, 9635, 9472, 8952, 8949, 992, 991, 990, 989, 7865, 7864, 7809, 7807, 7806}

	markets, err := repo.Markets(ctx, gameServerId)
	if err != nil {
		return err
	}
	market, ok := markets[marketId]
	if !ok {
		return errors.New("market not found")
	}
	etherealMarket, ok := markets[etherealMarketId]
	if !ok {
		return errors.New("ethereal market not found")
	}

	effectTemplates, err := repo.EffectTemplates(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		itemTemplates, err := repo.ItemTemplates(ctx)
		if err != nil {
			return err
		}

		market := market
		if i == 1 {
			market = etherealMarket
		}

		var items []d1.MarketItem

	itemTemplates:
		for k, v := range itemTemplates {
			allowed := true
			for _, id := range disallowedItemIds {
				if id == k {
					allowed = false
					break
				}
			}

			contains := false
			for _, itemType := range market.Types {
				if itemType == v.Type {
					contains = true
					break
				}
			}

			if !allowed || !contains ||
				(i == 0 && v.Ethereal) || (i == 1 && !v.Ethereal) ||
				v.Name == "" || v.Name == "x" || v.Hidden || v.Cursed || v.GFX == 0 || v.Level > 200 ||
				(v.Id >= 10954 && v.Id <= 11464) ||
				v.ItemSetId >= 150 && v.ItemSetId != 151 ||
				v.Name == "Ecaflip Paw" ||
				v.ItemSetId == 45 || v.ItemSetId == 59 || v.ItemSetId == 130 ||
				(strings.Contains(v.Name, "Trophy ") && strings.Contains(v.Name, "Shield")) ||
				(strings.Contains(v.Name, "Master ") && strings.Contains(v.Name, "Shield")) ||
				strings.Contains(v.Conditions, "Ps=3") ||
				strings.Contains(v.Conditions, "BI=") ||
				(v.Type == d1typ.ItemTypeCandy && !strings.Contains(v.Name, "Shigekax")) ||
				(v.Type == d1typ.ItemTypeUsableItem && !(v.Id == 7799 || v.Id == 8626)) {
				continue
			}

			var effects []d1typ.Effect
		effects:
			for _, effect := range v.Effects {
				// Shushumi weapon
				if effect.Id == 740 {
					continue itemTemplates
				}

				for _, v := range effectIdsToRemove {
					if v == effect.Id {
						continue effects
					}
				}

				t := effectTemplates[effect.Id]
				if t.Dice && effect.DiceSide > 0 {
					if t.Operator == d1typ.EffectOperatorAdd {
						effect.DiceNum = effect.DiceSide
						effect.DiceSide = 0
						effect.Param = d1.EffectDiceParam(effect)
					} else if t.Operator == d1typ.EffectOperatorSub {
						effect.DiceSide = 0
						effect.Param = d1.EffectDiceParam(effect)
					}
				}

				// durability
				if effect.Id == 812 {
					effect.DiceNum = effect.Value
					effect.DiceSide = effect.Value
				}

				effects = append(effects, effect)
			}

			if v.Type == d1typ.ItemTypePet {
				var petItems []d1.MarketItem

				// pets with multiple variable stats
				if v.Id == 8154 || v.Id == 7891 || v.Id == 7714 || v.Id == 7520 || v.Id == 2074 || v.Id == 2075 || v.Id == 2076 || v.Id == 2077 || v.Id == 1748 || v.Id == 1728 {
					for _, effect := range effects {
						petItems = append(petItems, d1.MarketItem{
							Item: d1.Item{
								TemplateId: v.Id,
								Quantity:   1,
								Effects:    []d1typ.Effect{effect},
							},
							Price:    1,
							MarketId: market.Id,
						})
					}
				} else {
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects:    effects,
						},
						Price:    1,
						MarketId: market.Id,
					})
				}

				// El Scarador and Croum
				if v.Id == 8154 || v.Id == 7520 {
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      214,
									DiceNum: 6,
								},
								{
									Id:      210,
									DiceNum: 5,
								},
								{
									Id:      213,
									DiceNum: 5,
								},
								{
									Id:      211,
									DiceNum: 5,
								},
								{
									Id:      212,
									DiceNum: 5,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      214,
									DiceNum: 13,
								},
								{
									Id:      210,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      214,
									DiceNum: 13,
								},
								{
									Id:      213,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      214,
									DiceNum: 13,
								},
								{
									Id:      211,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      214,
									DiceNum: 13,
								},
								{
									Id:      212,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      210,
									DiceNum: 13,
								},
								{
									Id:      213,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      210,
									DiceNum: 13,
								},
								{
									Id:      211,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      210,
									DiceNum: 13,
								},
								{
									Id:      212,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      213,
									DiceNum: 13,
								},
								{
									Id:      211,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      213,
									DiceNum: 13,
								},
								{
									Id:      212,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
					petItems = append(petItems, d1.MarketItem{
						Item: d1.Item{
							TemplateId: v.Id,
							Quantity:   1,
							Effects: []d1typ.Effect{
								{
									Id:      211,
									DiceNum: 13,
								},
								{
									Id:      212,
									DiceNum: 13,
								},
							},
						},
						Price:    1,
						MarketId: market.Id,
					})
				}

				for _, item := range petItems {
					if v.Id != 7523 && v.Id != 6895 && v.Id != 6894 && v.Id != 6718 && v.Id != 6717 {
						otherPetEffects := []d1typ.Effect{
							{
								Id:       800,
								DiceNum:  0,
								DiceSide: 0,
								Value:    10,
							},
						}

						if v.Id != 6604 && v.Id != 1711 {
							otherPetEffects = append(otherPetEffects, d1typ.Effect{
								Id:       940,
								DiceNum:  0,
								DiceSide: 0,
								Value:    1,
							})
						}

						item.Effects = append(otherPetEffects, item.Effects...)
					}
					items = append(items, item)
				}
			} else {
				items = append(items, d1.MarketItem{
					Item: d1.Item{
						TemplateId: v.Id,
						Quantity:   1,
						Effects:    effects,
					},
					Price:    1,
					MarketId: market.Id,
				})
			}
		}

		for _, item := range items {
			_, err := repo.CreateMarketItem(ctx, item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createEffectTemplates() error {
	p, err := ioutil.ReadFile("C:/Users/raul/AppData/Roaming/airtest/Local Store/f.json")
	if err != nil {
		return err
	}

	var raws map[string]struct {
		Description    string `json:"d"`
		Dice           bool   `json:"j"`
		Operator       string `json:"o"`
		Characteristic int    `json:"c"`
		Element        string `json:"e"`
	}

	err = json.Unmarshal(p, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		id, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		query := "INSERT INTO d1_static.effects (id, description, dice, operator, characteristic, element)" +
			" VALUES ($1, $2, $3, $4, $5, $6);"

		_, err = pool.Exec(ctx, query, id, v.Description, v.Dice, v.Operator, v.Characteristic, v.Element)
		if err != nil {
			return err
		}
	}

	return nil
}

func createNPCDialogsAndResponses() error {
	p, err := ioutil.ReadFile("C:/Users/raul/AppData/Roaming/airtest/Local Store/f.json")
	if err != nil {
		return err
	}

	var raws struct {
		Q map[string]string `json:"q"`
		A map[string]string `json:"a"`
	}

	err = json.Unmarshal(p, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws.Q {
		id, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		query := "INSERT INTO d1_static.npc_dialogs (id, text)" +
			" VALUES ($1, $2);"

		_, err = pool.Exec(ctx, query, id, v)
		if err != nil {
			return err
		}
	}

	for k, v := range raws.A {
		id, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		query := "INSERT INTO d1_static.npc_responses (id, text)" +
			" VALUES ($1, $2);"

		_, err = pool.Exec(ctx, query, id, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func decryptGameMapData() error {
	gameMaps, err := repo.GameMaps(ctx)
	if err != nil {
		return err
	}
	for _, gameMap := range gameMaps {
		data, err := d1util.DecipherGameMap(gameMap.EncryptedData, gameMap.Key)
		if err != nil {
			return err
		}
		gameMap.Data = data

		query := "UPDATE d1_static.maps" +
			" SET data = $2" +
			" WHERE id = $1;"

		_, err = pool.Exec(ctx, query, gameMap.Id, gameMap.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func orderItemTemplateEffects() error {
	itemTemplates, err := repo.ItemTemplates(ctx)
	if err != nil {
		return err
	}

	orderIds := []int{
		100,
		97,
		99,
		96,
		98,

		95,
		92,
		94,
		91,
		93,

		82,
		81,
		108,
		143,
		407,
		646,
		90,
		101,
		84,
		77,

		130,

		111,
		168,

		128,
		169,

		117,
		116,

		182,

		125,
		153,

		124,
		156,

		118,
		157,

		126,
		155,

		123,
		152,

		119,
		154,

		138,
		186,

		112,
		145,

		226,
		225,

		115,
		171,
		122,

		174,
		175,

		178,
		179,

		176,
		177,

		214,
		219,

		210,
		215,

		213,
		218,

		211,
		216,

		212,
		216,

		244,
		249,

		240,
		245,

		243,
		248,

		241,
		246,

		242,
		247,

		254,
		259,

		250,
		255,

		253,
		258,

		251,
		256,

		252,
		257,

		264,
		260,
		263,
		261,
		262,
	}
	for k, v := range itemTemplates {
		sort.SliceStable(v.Effects, func(i, j int) bool {
			e := v.Effects[i]
			e2 := v.Effects[j]

			ix := math.MaxInt32
			ix2 := math.MaxInt32

			for i2, id := range orderIds {
				if id == e.Id {
					ix = i2
				}
				if id == e2.Id {
					ix2 = i2
				}
			}

			return ix < ix2
		})

		effects := d1.EncodeItemEffects(v.Effects)

		query := "UPDATE d1_static.items" +
			" SET effects = $2" +
			" WHERE id = $1;"

		_, err = pool.Exec(ctx, query, k, strings.Join(effects, ","))
		if err != nil {
			return err
		}
	}

	return nil
}

func orderItemSetEffects() error {
	itemSets, err := repo.ItemSets(ctx)
	if err != nil {
		return err
	}

	orderIds := []int{
		100,
		97,
		99,
		96,
		98,

		95,
		92,
		94,
		91,
		93,

		82,
		81,
		108,
		143,
		407,
		646,
		90,
		101,
		84,
		77,

		130,

		111,
		168,

		128,
		169,

		117,
		116,

		182,

		125,
		153,

		124,
		156,

		118,
		157,

		126,
		155,

		123,
		152,

		119,
		154,

		138,
		186,

		112,
		145,

		226,
		225,

		115,
		171,
		122,

		174,
		175,

		178,
		179,

		176,
		177,

		214,
		219,

		210,
		215,

		213,
		218,

		211,
		216,

		212,
		216,

		244,
		249,

		240,
		245,

		243,
		248,

		241,
		246,

		242,
		247,

		254,
		259,

		250,
		255,

		253,
		258,

		251,
		256,

		252,
		257,

		264,
		260,
		263,
		261,
		262,
	}
	for k, itemSet := range itemSets {
		for _, effects := range itemSet.Bonus {
			sort.SliceStable(effects, func(i, j int) bool {
				e := effects[i]
				e2 := effects[j]

				ix := math.MaxInt32
				ix2 := math.MaxInt32

				for i2, id := range orderIds {
					if id == e.Id {
						ix = i2
					}
					if id == e2.Id {
						ix2 = i2
					}
				}

				return ix < ix2
			})
		}

		bonus := make([]string, len(itemSet.Bonus))
		for i, effects := range itemSet.Bonus {
			effectsStr := make([]string, len(effects))
			for i, effect := range effects {
				effectsStr[i] = fmt.Sprintf("%d:%d", effect.Id, effect.DiceNum)
			}
			bonus[i] = strings.Join(effectsStr, ",")
		}

		query := "UPDATE d1_static.itemsets" +
			" SET bonus = $2" +
			" WHERE id = $1;"

		_, err = pool.Exec(ctx, query, k, strings.Join(bonus, ";"))
		if err != nil {
			return err
		}
	}

	return nil
}

func updateItemTemplateTexts() error {
	p, err := ioutil.ReadFile("C:/Users/raul/AppData/Roaming/airtest/Local Store/f.json")
	if err != nil {
		return err
	}

	var raws map[string]ItemTemplate

	err = json.Unmarshal(p, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		id, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		query := "UPDATE d1_static.items" +
			" SET name = $2, description = $3" +
			" WHERE id = $1;"

		_, err = pool.Exec(ctx, query, id, v.Name, v.Description)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateItemTemplateEffects() error {
	p, err := ioutil.ReadFile("C:/Users/raul/AppData/Roaming/airtest/Local Store/f.json")
	if err != nil {
		return err
	}

	var raws map[string]string

	err = json.Unmarshal(p, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		v = strings.Trim(v, ",")
		v = strings.ReplaceAll(v, ",,", ",")
		if v == "" {
			continue
		}

		id, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		query := "UPDATE d1_static.itemts" +
			" SET effects = $2" +
			" WHERE id = $1;"

		_, err = pool.Exec(ctx, query, id, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func createItemTemplates() error {
	p, err := ioutil.ReadFile("C:/Users/raul/AppData/Roaming/airtest/Local Store/f.json")
	if err != nil {
		return err
	}

	var raws map[string]ItemTemplate

	err = json.Unmarshal(p, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		id, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		itemTemplate := d1.ItemTemplate{
			Id:            id,
			Name:          v.Name,
			Description:   v.Description,
			Type:          d1typ.ItemType(v.Type),
			Enhanceable:   v.Enhanceable,
			TwoHands:      v.TwoHands,
			Ethereal:      v.Ethereal,
			Hidden:        v.Hidden,
			ItemSetId:     v.ItemSetId,
			CanUse:        v.CanUse,
			CanTarget:     v.CanTarget,
			Level:         v.Level,
			GFX:           v.GFX,
			Price:         v.Price,
			Weight:        v.Weight,
			Cursed:        v.Cursed,
			Conditions:    v.Conditions,
			WeaponEffects: d1.WeaponEffects{},
		}

		var weaponEffects string
		if len(v.WeaponEffects) == 8 {
			criticalHitBonus, ok := v.WeaponEffects[0].(float64)
			if !ok {
				return errInvalidAssertion
			}
			itemTemplate.WeaponEffects.CriticalHitBonus = int(criticalHitBonus)

			apCost, ok := v.WeaponEffects[1].(float64)
			if !ok {
				return errInvalidAssertion
			}
			itemTemplate.WeaponEffects.APCost = int(apCost)

			rangeMin, ok := v.WeaponEffects[2].(float64)
			if !ok {
				return errInvalidAssertion
			}
			itemTemplate.WeaponEffects.RangeMin = int(rangeMin)

			rangeMax, ok := v.WeaponEffects[3].(float64)
			if !ok {
				return errInvalidAssertion
			}
			itemTemplate.WeaponEffects.RangeMax = int(rangeMax)

			criticalHit, ok := v.WeaponEffects[4].(float64)
			if !ok {
				return errInvalidAssertion
			}
			itemTemplate.WeaponEffects.CriticalHit = int(criticalHit)

			criticalFailure, ok := v.WeaponEffects[5].(float64)
			if !ok {
				return errInvalidAssertion
			}
			itemTemplate.WeaponEffects.CriticalFailure = int(criticalFailure)

			itemTemplate.WeaponEffects.LineOnly, ok = v.WeaponEffects[6].(bool)
			if !ok {
				return errInvalidAssertion
			}

			itemTemplate.WeaponEffects.LineOfSight, ok = v.WeaponEffects[7].(bool)
			if !ok {
				return errInvalidAssertion
			}

			weaponEffects = fmt.Sprintf("%d,%d,%d,%d,%d,%d,%t,%t",
				itemTemplate.WeaponEffects.CriticalHitBonus,
				itemTemplate.WeaponEffects.APCost,
				itemTemplate.WeaponEffects.RangeMin,
				itemTemplate.WeaponEffects.RangeMax,
				itemTemplate.WeaponEffects.CriticalHit,
				itemTemplate.WeaponEffects.CriticalFailure,
				itemTemplate.WeaponEffects.LineOnly,
				itemTemplate.WeaponEffects.LineOfSight,
			)
		}

		query := "INSERT INTO d1_static.items (id, name, description, type, enhanceable, two_hands, ethereal, hidden, itemset_id, can_use, can_target, level, gfx, price, weight, cursed, conditions, weapon_effects, effects)" +
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19);"

		_, err = pool.Exec(ctx, query,
			itemTemplate.Id, itemTemplate.Name, itemTemplate.Description, itemTemplate.Type, itemTemplate.Enhanceable,
			itemTemplate.TwoHands, itemTemplate.Ethereal, itemTemplate.Hidden, itemTemplate.ItemSetId,
			itemTemplate.CanUse, itemTemplate.CanTarget, itemTemplate.Level, itemTemplate.GFX, itemTemplate.Price,
			itemTemplate.Weight, itemTemplate.Cursed, itemTemplate.Conditions, weaponEffects, "",
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func addTargetIdToSpells() error {
	p, err := ioutil.ReadFile("assets/spell_levels.json")
	if err != nil {
		return err
	}

	var levels map[int]D2Level
	err = json.Unmarshal(p, &levels)
	if err != nil {
		return err
	}

	d2spells := make(map[int][]D2Level)
	for _, level := range levels {
		d2spells[level.SpellId] = append(d2spells[level.SpellId], level)
	}

	for _, d2spell := range d2spells {
		sort.Slice(d2spell, func(i, j int) bool {
			return d2spell[i].Id < d2spell[j].Id
		})
	}

	spells, err := repo.Spells(ctx)
	if err != nil {
		return err
	}

	for _, spell := range spells {
		d2spell, ok := d2spells[spell.Id]
		if !ok {
			continue
		}

		for i, level := range spell.Levels {
			if len(d2spell) < i+1 {
				continue
			}
			d2level := d2spell[i]

			usedIndexes := make(map[int]struct{})
			for i, effect := range level.Effects {
				for i, d2effect := range d2level.Effects {
					_, ok := usedIndexes[i]
					if ok {
						continue
					}

					if d2effect.Id != effect.Id {
						continue
					}
					if d2effect.Param1 != effect.DiceNum {
						continue
					}
					if d2effect.Param2 != effect.DiceSide {
						continue
					}

					effect.TargetId = d2effect.TargetId

					usedIndexes[i] = struct{}{}
					break
				}
				level.Effects[i] = effect
			}

			usedIndexesCritical := make(map[int]struct{})
			for i, effect := range level.EffectsCritical {
				for i, d2effect := range d2level.CriticalEffects {
					_, ok := usedIndexesCritical[i]
					if ok {
						continue
					}

					if d2effect.Id != effect.Id {
						continue
					}
					if d2effect.Param1 != effect.DiceNum {
						continue
					}
					if d2effect.Param2 != effect.DiceSide {
						continue
					}

					effect.TargetId = d2effect.TargetId

					usedIndexesCritical[i] = struct{}{}
					break
				}
				level.EffectsCritical[i] = effect
			}
		}
	}

	for _, spell := range spells {
		sli, err := encodeSpellLevels(spell.Levels)
		if err != nil {
			return err
		}

		query := "UPDATE d1_static.spells" +
			" SET levels = $2" +
			" WHERE id = $1;"

		_, err = pool.Exec(ctx, query, spell.Id, sli)
		if err != nil {
			return err
		}
	}

	return nil
}

func encodeSpellLevels(levels []d1typ.SpellLevel) ([]string, error) {
	if len(levels) == 0 {
		return nil, nil
	}

	sli := make([]string, len(levels))
	for i, level := range levels {
		encoded, err := encodeSpellLevel(level)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(encoded)
		if err != nil {
			return nil, err
		}

		sli[i] = string(data)
	}
	return sli, nil
}

func encodeSpellLevel(level d1typ.SpellLevel) ([]interface{}, error) {
	effects := encodeSpellLevelEffects(level.Effects)
	effectsCritical := encodeSpellLevelEffects(level.EffectsCritical)

	// Empty effects (and not critical effects) need to be `[]` instead of `null`. See original spell 514, for example.
	if effects == nil {
		effects = []interface{}{}
	}

	effectZones := make([]string, len(level.Effects)+len(level.EffectsCritical))
	for i, effect := range level.Effects {
		size, err := d1proto.Encode64(effect.ZoneSize)
		if err != nil {
			return nil, err
		}
		effectZones[i] = fmt.Sprintf("%s%s", string(effect.ZoneShape), string(size))
	}
	for i, effect := range level.EffectsCritical {
		size, err := d1proto.Encode64(effect.ZoneSize)
		if err != nil {
			return nil, err
		}
		effectZones[i+len(level.Effects)] = fmt.Sprintf("%s%s", string(effect.ZoneShape), string(size))
	}

	data := []interface{}{effects, effectsCritical, level.APCost, level.Range, level.RangeMax,
		level.CriticalHitProbability, level.CriticalFailureProbability, level.Linear, level.RequiresLineOfSight,
		level.RequiresFreeCell, level.AdjustableRange, level.ClassId, level.MaxCastsPerTurn,
		level.MaxCastsPerTarget, level.MinCastInterval, strings.Join(effectZones, ""), level.StatesRequired,
		level.StatesForbidden, level.MinPlayerLevel, level.CriticalFailureEndsTurn,
	}

	return data, nil
}

func encodeSpellLevelEffects(effects []d1typ.Effect) []interface{} {
	if len(effects) == 0 {
		return nil
	}

	sli := make([]interface{}, len(effects))
	for i, effect := range effects {
		sli[i] = encodeSpellLevelEffect(effect)
	}
	return sli
}

func encodeSpellLevelEffect(effect d1typ.Effect) []interface{} {
	var diceNum *int
	if effect.DiceNum != -1 {
		diceNum = &effect.DiceNum
	}

	var diceSide *int
	if effect.DiceSide != -1 {
		diceSide = &effect.DiceSide
	}

	var value *int
	if effect.Value != -1 {
		value = &effect.Value
	}

	data := []interface{}{effect.Id, diceNum, diceSide, value, effect.Duration, effect.Random, effect.TargetId}

	if effect.Param != "" {
		data = append(data, effect.Param)
	}

	return data
}
