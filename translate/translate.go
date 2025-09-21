package translate

import (
	"slices"
	"strings"
)

type SplitData struct {
	ID          string
	Description string
}

func GetSplitDescriptionByID(id string) string {
	for _, split := range SplitsCache {
		if split.ID == id {
			return split.Description
		}
	}
	return ""
}

func GetIndexByID(id string) int {
	for i, split := range SplitsCache {
		if split.ID == id {
			return i
		}
	}
	return -1
}

func GetIDByDescription(desc string) string {
	for _, split := range SplitsCache {
		if split.Description == desc {
			return split.ID
		}
	}
	return ""
}

func init() {
	typeOrder := map[string]int{
		"（其它）": -100,
		"（开始）": -90,
		"（结束）": -80,
		"（旋律）": -70,
		"（道具）": -60,
	}
	slices.SortStableFunc(SplitsCache, func(a, b *SplitData) int {
		i := strings.LastIndex(a.Description, "（")
		j := strings.LastIndex(b.Description, "（")
		if i == -1 {
			i = len(a.Description)
		}
		if j == -1 {
			j = len(b.Description)
		}
		ad, bd := a.Description[i:], b.Description[j:]
		aScore, bScore := typeOrder[ad], typeOrder[bd]
		if aScore != bScore {
			return aScore - bScore
		}
		return strings.Compare(ad, bd)
	})
}

var SplitsCache = []*SplitData{
	{"ManualSplit", "手动分割（其它）"},
	{"StartNewGame", "开始新游戏（开始）"},
	{"EndingSplit", "任意结束（结束）"},
	{"EndingA", "苍白之母结局（结束）"},
	{"Menu", "主菜单（菜单）"},
	{"PlayerDeath", "死亡（事件）"},
	{"AnyTransition", "任意切图（切图）"},
	{"MossMother", "苔藓之母（Boss）"},
	{"MossMotherTrans", "苔藓之母（切图）"},
	{"SilkSpear", "丝之矛（技能）"},
	{"SilkSpearTrans", "丝之矛（切图）"},
	{"BellBeast", "钟道兽（Boss）"},
	{"BellBeastTrans", "钟道兽（切图）"},
	{"MarrowBell", "敲钟-髓骨洞窟（事件）"},
	{"SwiftStep", "冲刺（技能）"},
	{"SwiftStepTrans", "冲刺（切图）"},
	{"Lace1", "蕾丝（Boss）"},
	{"Lace1Trans", "蕾丝（切图）"},
	{"DeepDocksBell", "敲钟-深坞（事件）"},
	{"DriftersCloak", "斗篷（技能）"},
	{"DriftersCloakTrans", "斗篷（切图）"},
	{"FourthChorus", "第四圣咏团（Boss）"},
	{"EnterGreymoor", "进入灰沼（切图）"},
	{"GreymoorBell", "敲钟-灰沼（事件）"},
	{"Moorwing", "荒沼翼主（Boss）"},
	{"MoorwingTrans", "荒沼翼主（切图）"},
	{"ClingGrip", "爬墙（技能）"},
	{"ClingGripTrans", "爬墙（切图）"},
	{"ShellwoodBell", "敲钟-甲木林（事件）"},
	{"Widow", "黑寡妇（Boss）"},
	{"BellhartBell", "敲钟-钟心镇（事件）"},
	{"LastJudge", "末代裁决者（Boss）"},
	{"EnterMist", "进入迷雾（切图）"},
	{"LeaveMist", "离开迷雾（切图）"},
	{"Phantom", "幽影（Boss）"},
	{"Act2Started", "第二幕开始（事件）"},
	{"CogworkDancers", "机枢舞者（Boss）"},
	{"WhisperingVaultsGauntlet", "低语书库遭遇战（小Boss）"},
	{"Trobbio", "特罗比奥（Boss）"},
	{"TrobbioTrans", "特罗比奥（切图）"},
	{"Clawline", "飞针冲刺（技能）"},
	{"EnterHighHalls", "进入高庭（切图）"},
	{"EnterHighHallsGauntlet", "进入高庭遭遇战（切图）"},
	{"HighHallsGauntlet", "高庭遭遇战（小Boss）"},
	{"Lace2", "蕾丝2（Boss）"},
	{"VaultkeepersMelody", "管理员旋律（旋律）"},
	{"VaultkeepersMelodyTrans", "管理员旋律（切图）"},
	{"ArchitectsMelody", "建筑师旋律（旋律）"},
	{"ArchitectsMelodyTrans", "建筑师旋律（切图）"},
	{"ConductorsMelody", "指挥家旋律（旋律）"},
	{"ConductorsMelodyTrans", "指挥家旋律（切图）"},
	{"UnlockedMelodyLift", "解锁三重旋律电梯（事件）"},
	{"NeedleUpgrade1", "织针升级1（升级）"},
	{"NeedleUpgrade2", "织针升级2（升级）"},
	{"NeedleUpgrade3", "织针升级3（升级）"},
	{"NeedleUpgrade4", "织针升级4（升级）"},
	{"SavedFleaHuntersMarch", "救跳蚤-猎者小径（跳蚤）"},
	{"SavedFleaBellhart", "救跳蚤-钟心镇（跳蚤）"},
	{"SavedFleaMarrow", "救跳蚤-髓骨洞窟（跳蚤）"},
	{"SavedFleaDeepDocksSprint", "救跳蚤-深坞-冲刺（跳蚤）"},
	{"SavedFleaFarFieldsPilgrimsRest", "救跳蚤-远野-朝圣者憩所（跳蚤）"},
	{"SavedFleaFarFieldsTrap", "救跳蚤-远野-陷阱（跳蚤）"},
	{"SavedFleaSandsOfKarak", "救跳蚤-卡拉卡沙川（跳蚤）"},
	{"SavedFleaBlastedSteps", "救跳蚤-卡拉卡沙川-蚀阶（跳蚤）"},
	{"SavedFleaWormways", "救跳蚤-沙噬虫道（跳蚤）"},
	{"SavedFleaDeepDocksArena", "救跳蚤-深坞-遭遇战（跳蚤）"},
	{"SavedFleaDeepDocksBellway", "救跳蚤-深坞-钟道（跳蚤）"},
	{"SavedFleaBilewaterOrgan", "救跳蚤-废鸣管风琴（跳蚤）"},
	{"SavedFleaSinnersRoad", "救跳蚤-罪途（跳蚤）"},
	{"SavedFleaGreymoorRoof", "救跳蚤-灰沼-屋顶（跳蚤）"},
	{"SavedFleaGreymoorLake", "救跳蚤-灰沼-湖（跳蚤）"},
	{"SavedFleaWhisperingVaults", "救跳蚤-低语书库（跳蚤）"},
	{"SavedFleaSongclave", "救跳蚤-低语书库-圣歌盟地（跳蚤）"},
	{"SavedFleaMountFay", "救跳蚤-费耶山（跳蚤）"},
	{"SavedFleaBilewaterTrap", "救跳蚤-腐汁泽-陷阱（跳蚤）"},
	{"SavedFleaBilewaterThieves", "救跳蚤-腐汁泽-小偷（跳蚤）"},
	{"SavedFleaShellwood", "救跳蚤-甲木林（跳蚤）"},
	{"SavedFleaSlabBellway", "救跳蚤-罪石监狱-钟道（跳蚤）"},
	{"SavedFleaSlabCage", "救跳蚤-罪石监狱-笼子（跳蚤）"},
	{"SavedFleaChoralChambersWind", "救跳蚤-圣咏殿-风（跳蚤）"},
	{"SavedFleaChoralChambersCage", "救跳蚤-圣咏殿-笼子（跳蚤）"},
	{"SavedFleaUnderworksExplosions", "救跳蚤-圣堡工厂-爆炸怪（跳蚤）"},
	{"SavedFleaUnderworksWispThicket", "救跳蚤-圣堡工厂-火灵竹丛（跳蚤）"},
	{"SavedFleaGiantFlea", "击败大跳蚤（跳蚤）"},
	{"SavedFleaVog", "沃格（跳蚤）"},
	{"SavedFleaKratt", "救跳蚤-灰沼-克拉特（跳蚤）"},
	{"PutrifiedDuctsStation", "钟道-腐殖渠（钟道）"},
	{"BellhartStation", "钟道-钟心镇（钟道）"},
	{"FarFieldsStation", "钟道-远野（钟道）"},
	{"GrandBellwayStation", "钟道-圣堡钟道（钟道）"},
	{"BlastedStepsStation", "钟道-蚀阶（钟道）"},
	{"DeepDocksStation", "钟道-深坞（钟道）"},
	{"GreymoorStation", "钟道-灰沼（钟道）"},
	{"MountFayStation", "钟道-费耶山（钟道）"},
	{"BilewaterStation", "钟道-腐汁泽（钟道）"},
	{"ShellwoodStation", "钟道-甲木林（钟道）"},
	{"ChoralChambersTube", "圣脉枢管-圣咏殿（圣脉枢管）"},
	{"UnderworksTube", "圣脉枢管-圣堡工厂（圣脉枢管）"},
	{"CityBellwayTube", "圣脉枢管-圣堡钟道（圣脉枢管）"},
	{"HighHallsTube", "圣脉枢管-高庭（圣脉枢管）"},
	{"SongclaveTube", "圣脉枢管-始源钟殿（圣脉枢管）"},
	{"MemoriumTube", "圣脉枢管-忆廊（圣脉枢管）"},
	{"SeenShakraBonebottom", "制图师-骸底镇（NPC）"},
	{"SeenShakraMarrow", "制图师-髓骨洞窟（NPC）"},
	{"SeenShakraDeepDocks", "制图师-深坞（NPC）"},
	{"SeenShakraFarFields", "制图师-远野（NPC）"},
	{"SeenShakraWormways", "制图师-沙噬虫道（NPC）"},
	{"SeenShakraGreymoor", "制图师-灰沼（NPC）"},
	{"SeenShakraBellhart", "制图师-钟心镇（NPC）"},
	{"SeenShakraShellwood", "制图师-甲木林（NPC）"},
	{"SeenShakraHuntersMarch", "制图师-猎者小径（NPC）"},
	{"SeenShakraBlastedSteps", "制图师-蚀阶（NPC）"},
	{"SeenShakraSinnersRoad", "制图师-罪途（NPC）"},
	{"SeenShakraMountFay", "制图师-费耶山（NPC）"},
	{"SeenShakraBilewater", "制图师-腐汁泽（NPC）"},
	{"SeenShakraSandsOfKarak", "制图师-卡拉卡沙川（NPC）"},
	{"MetJubilanaEnclave", "朱比拉娜-圣歌盟地（NPC）"},
	{"MetShermaEnclave", "谢尔玛-圣歌盟地（NPC）"},
	{"UnlockedPrinceCage", "绿王子-罪途（事件）"},
	{"GreenPrinceInVerdania", "绿王子-圣堡工厂（事件）"},
	{"SeenFleatopiaEmpty", "蚤托邦（事件）"},
	{"FaydownCloak", "二段跳（技能）"},
	{"SilkSoar", "灵丝升腾（技能）"},
	{"CollectedHeartNyleth", "尼莱斯之心（道具）"},
	{"CollectedHeartKhann", "卡汗之心（道具）"},
	{"CollectedHeartKarmelita", "卡梅莉塔之心（道具）"},
	{"CollectedHeartClover", "双生之心（道具）"},
	{"CompletedRedMemory", "赤红忆境（事件）"},
	{"BellhouseKeyConversation", "钟居钥匙（NPC）"},
}
